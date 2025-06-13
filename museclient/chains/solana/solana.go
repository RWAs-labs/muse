package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/solana/observer"
	"github.com/RWAs-labs/muse/museclient/chains/solana/signer"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/scheduler"
	"github.com/RWAs-labs/muse/pkg/ticker"
)

const (
	// outboundLookbackFactor is the factor to determine how many nonces to look back for pending cctxs
	// For example, give OutboundScheduleLookahead of 120, pending NonceLow of 1000 and factor of 1.0,
	// the scheduler need to be able to pick up and schedule any pending cctx with nonce < 880 (1000 - 120 * 1.0)
	// NOTE: 1.0 means look back the same number of cctxs as we look ahead
	outboundLookbackFactor = 1.0
)

// Solana represents Solana observer-signer.
type Solana struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
}

// New Solana constructor.
func New(scheduler *scheduler.Scheduler, observer *observer.Observer, signer *signer.Signer) *Solana {
	return &Solana{
		scheduler: scheduler,
		observer:  observer,
		signer:    signer,
	}
}

// Chain returns chain
func (s *Solana) Chain() chains.Chain {
	return s.observer.Chain()
}

// Start starts observer-signer for
// processing inbound & outbound cross-chain transactions.
func (s *Solana) Start(ctx context.Context) error {
	if ok := s.observer.Observer.Start(); !ok {
		return errors.New("observer is already started")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get app from context")
	}

	newBlockChan, err := s.observer.MusecoreClient().NewBlockSubscriber(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to create new block subscriber")
	}

	optInboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().InboundTicker)
	})

	optGasInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().GasPriceTicker)
	})

	optOutboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().OutboundTicker)
	})

	optInboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsInboundObservationEnabled()
	})

	optOutboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsOutboundObservationEnabled()
	})

	optGenericSkipper := scheduler.Skipper(func() bool {
		return !s.observer.ChainParams().IsSupported
	})

	register := func(exec scheduler.Executable, name string, opts ...scheduler.Opt) {
		opts = append([]scheduler.Opt{
			scheduler.GroupName(s.group()),
			scheduler.Name(name),
		}, opts...)

		s.scheduler.Register(ctx, exec, opts...)
	}

	register(s.observer.ObserveInbound, "observe_inbound", optInboundInterval, optInboundSkipper)
	register(s.observer.ProcessInboundTrackers, "process_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(s.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)
	register(s.observer.CheckRPCStatus, "check_rpc_status")
	register(s.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)
	register(s.observer.ProcessOutboundTrackers, "process_outbound_trackers", optOutboundInterval, optOutboundSkipper)

	// CCTX scheduler (every musechain block)
	register(s.scheduleCCTX, "schedule_cctx", scheduler.BlockTicker(newBlockChan), optOutboundSkipper)

	return nil
}

// Stop stops all relevant tasks.
func (s *Solana) Stop() {
	s.observer.Logger().Chain.Info().Msg("stopping observer-signer")
	s.scheduler.StopGroup(s.group())
}

func (s *Solana) group() scheduler.Group {
	return scheduler.Group(
		fmt.Sprintf("sol:%d", s.observer.Chain().ChainId),
	)
}

// scheduleCCTX schedules solana outbound keysign
func (s *Solana) scheduleCCTX(ctx context.Context) error {
	if err := s.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	museBlock, delay, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get muse block from context")
	}

	time.Sleep(delay)

	var (
		chain   = s.observer.Chain()
		chainID = chain.ChainId

		// #nosec G115 positive
		museHeight = uint64(museBlock.Block.Height)

		// #nosec G115 positive
		interval           = uint64(s.observer.ChainParams().OutboundScheduleInterval)
		scheduleLookahead  = s.observer.ChainParams().OutboundScheduleLookahead
		scheduleLookback   = uint64(float64(scheduleLookahead) * outboundLookbackFactor)
		needsProcessingCtr = 0
	)

	cctxList, _, err := s.observer.MusecoreClient().ListPendingCCTX(ctx, chain)
	if err != nil {
		return errors.Wrap(err, "unable to list pending cctx")
	}

	// schedule keysign for each pending cctx
	for _, cctx := range cctxList {
		var (
			params        = cctx.GetCurrentOutboundParam()
			inboundParams = cctx.GetInboundParams()
			nonce         = params.TssNonce
			outboundID    = base.OutboundIDFromCCTX(cctx)
		)

		switch {
		case params.ReceiverChainId != chainID:
			s.outboundLogger(outboundID).Error().Msg("chain id mismatch")
			continue
		case params.TssNonce > cctxList[0].GetCurrentOutboundParam().TssNonce+scheduleLookback:
			return fmt.Errorf(
				"nonce %d is too high (%s). Earliest nonce %d",
				params.TssNonce,
				outboundID,
				cctxList[0].GetCurrentOutboundParam().TssNonce,
			)
		}

		// schedule newly created cctx right away, no need to wait for next interval
		// 1. schedule the very first cctx (there can be multiple) created in the last Muse block.
		// 2. schedule new cctx only when there is no other older cctx to process
		isCCTXNewlyCreated := inboundParams.ObservedExternalHeight == museHeight
		shouldProcessCCTXImmedately := isCCTXNewlyCreated && needsProcessingCtr == 0

		// even if the outbound is currently active, we should increment this counter
		// to avoid immediate processing logic
		needsProcessingCtr++

		if s.signer.IsOutboundActive(outboundID) {
			continue
		}

		// vote outbound if it's already confirmed
		continueKeysign, err := s.observer.VoteOutboundIfConfirmed(ctx, cctx)
		switch {
		case err != nil:
			s.outboundLogger(outboundID).Error().Err(err).Msg("Schedule CCTX: VoteOutboundIfConfirmed failed")
			continue
		case !continueKeysign:
			s.outboundLogger(outboundID).Info().Msg("Schedule CCTX: outbound already processed")
			continue
		}

		shouldScheduleProcess := nonce%interval == museHeight%interval

		// schedule a TSS keysign
		if shouldProcessCCTXImmedately || shouldScheduleProcess {
			go s.signer.TryProcessOutbound(
				ctx,
				cctx,
				s.observer.MusecoreClient(),
				museHeight,
			)
		}
	}

	return nil
}

func (s *Solana) updateChainParams(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	chain, err := app.GetChain(s.observer.Chain().ChainId)
	if err != nil {
		return err
	}

	params := chain.Params()

	s.observer.SetChainParams(*params)
	s.signer.SetGatewayAddress(params.GatewayAddress)

	return nil
}

func (s *Solana) outboundLogger(id string) *zerolog.Logger {
	l := s.observer.Logger().Outbound.With().Str("outbound.id", id).Logger()

	return &l
}
