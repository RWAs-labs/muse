package ton

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/ton/observer"
	"github.com/RWAs-labs/muse/museclient/chains/ton/signer"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/scheduler"
	"github.com/RWAs-labs/muse/pkg/ticker"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// TON represents TON observer-signer components that is responsible
// for processing and scheduling inbound and outbound TON transactions.
type TON struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
}

// New TON constructor.
func New(scheduler *scheduler.Scheduler, observer *observer.Observer, signer *signer.Signer) *TON {
	return &TON{
		scheduler: scheduler,
		observer:  observer,
		signer:    signer,
	}
}

// Chain returns the chain struct
func (t *TON) Chain() chains.Chain {
	return t.observer.Chain()
}

// Start starts the observer-signer and schedules various regular background tasks e.g. inbound observation.
func (t *TON) Start(ctx context.Context) error {
	if ok := t.observer.Observer.Start(); !ok {
		return errors.Errorf("observer is already started")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get app from context")
	}

	newBlockChan, err := t.observer.MusecoreClient().NewBlockSubscriber(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to create new block subscriber")
	}

	optInboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(t.observer.ChainParams().InboundTicker)
	})

	optGasInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(t.observer.ChainParams().GasPriceTicker)
	})

	optOutboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(t.observer.ChainParams().OutboundTicker)
	})

	optInboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsInboundObservationEnabled()
	})

	optOutboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsOutboundObservationEnabled()
	})

	optGenericSkipper := scheduler.Skipper(func() bool {
		return !t.observer.ChainParams().IsSupported
	})

	register := func(exec scheduler.Executable, name string, opts ...scheduler.Opt) {
		opts = append([]scheduler.Opt{
			scheduler.GroupName(t.group()),
			scheduler.Name(name),
		}, opts...)

		t.scheduler.Register(ctx, exec, opts...)
	}

	register(t.observer.ObserveInbound, "observe_inbound", optInboundInterval, optInboundSkipper)
	register(t.observer.ProcessInboundTrackers, "process_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(t.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)
	register(t.observer.CheckRPCStatus, "check_rpc_status")
	register(t.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)
	register(t.observer.ProcessOutboundTrackers, "process_outbound_trackers", optOutboundInterval, optOutboundSkipper)

	// CCTX Scheduler
	register(t.scheduleCCTX, "schedule_cctx", scheduler.BlockTicker(newBlockChan), optOutboundSkipper)

	return nil
}

// Stop stops the observer-signer.
func (t *TON) Stop() {
	t.observer.Logger().Chain.Info().Msg("stopping observer-signer")
	t.scheduler.StopGroup(t.group())
}

func (t *TON) group() scheduler.Group {
	return scheduler.Group(
		fmt.Sprintf("ton:%d", t.observer.Chain().ChainId),
	)
}

// scheduleCCTX schedules cross-chain tx processing.
// It loads pending cctx from musecore, then tries to sign and broadcast them.
func (t *TON) scheduleCCTX(ctx context.Context) error {
	if err := t.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	museBlock, delay, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get muse block from context")
	}

	time.Sleep(delay)

	// #nosec G115 always in range
	museHeight := uint64(museBlock.Block.Height)
	chain := t.observer.Chain()

	cctxList, _, err := t.observer.MusecoreClient().ListPendingCCTX(ctx, chain)
	if err != nil {
		return errors.Wrap(err, "unable to list pending cctx")
	}

	for i := range cctxList {
		cctx := cctxList[i]
		outboundID := base.OutboundIDFromCCTX(cctx)

		if err := t.processCCTX(ctx, outboundID, cctx, museHeight); err != nil {
			t.outboundLogger(outboundID).Error().Err(err).Msg("Schedule CCTX failed")
		}
	}

	return nil
}

func (t *TON) processCCTX(ctx context.Context, outboundID string, cctx *types.CrossChainTx, museHeight uint64) error {
	switch {
	case t.signer.IsOutboundActive(outboundID):
		//noop
		return nil
	case cctx.GetCurrentOutboundParam().ReceiverChainId != t.observer.Chain().ChainId:
		return errors.New("chain id mismatch")
	}

	// vote outbound if it's already confirmed
	continueKeySign, err := t.observer.VoteOutboundIfConfirmed(ctx, cctx)
	switch {
	case err != nil:
		return errors.Wrap(err, "failed to VoteOutboundIfConfirmed")
	case !continueKeySign:
		t.outboundLogger(outboundID).Info().Msg("Schedule CCTX: outbound already processed")
		return nil
	}

	go t.signer.TryProcessOutbound(
		ctx,
		cctx,
		t.observer.MusecoreClient(),
		museHeight,
	)

	return nil
}

func (t *TON) updateChainParams(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	chain, err := app.GetChain(t.observer.Chain().ChainId)
	if err != nil {
		return err
	}

	t.signer.SetGatewayAddress(chain.Params().GatewayAddress)
	t.observer.SetChainParams(*chain.Params())

	return nil
}

func (t *TON) outboundLogger(id string) *zerolog.Logger {
	l := t.observer.Logger().Outbound.With().Str("outbound.id", id).Logger()

	return &l
}
