package sui

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/sui/observer"
	"github.com/RWAs-labs/muse/museclient/chains/sui/signer"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/pkg/bg"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/scheduler"
	"github.com/RWAs-labs/muse/pkg/ticker"
)

// Sui observer-signer.
type Sui struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
}

const (
	// outboundLookbackFactor is the factor to determine how many nonces to look back for pending cctxs
	// For example, give OutboundScheduleLookahead of 120, pending NonceLow of 1000 and factor of 1.0,
	// the scheduler need to be able to pick up and schedule any pending cctx with nonce < 880 (1000 - 120 * 1.0)
	// NOTE: 1.0 means look back the same number of cctxs as we look ahead
	outboundLookbackFactor = 1.0
)

// New Sui observer-signer constructor.
func New(scheduler *scheduler.Scheduler, observer *observer.Observer, signer *signer.Signer) *Sui {
	return &Sui{scheduler, observer, signer}
}

// Chain returns chain
func (s *Sui) Chain() chains.Chain {
	return s.observer.Chain()
}

// Start starts observer-signer for processing inbound & outbound cross-chain transactions.
func (s *Sui) Start(ctx context.Context) error {
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

	optOutboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsOutboundObservationEnabled()
	})

	register := func(exec scheduler.Executable, name string, opts ...scheduler.Opt) {
		opts = append([]scheduler.Opt{
			scheduler.GroupName(s.group()),
			scheduler.Name(name),
		}, opts...)

		s.scheduler.Register(ctx, exec, opts...)
	}

	optInboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().InboundTicker)
	})

	optOutboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().OutboundTicker)
	})

	optGasInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(s.observer.ChainParams().GasPriceTicker)
	})

	optInboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsInboundObservationEnabled()
	})

	optGenericSkipper := scheduler.Skipper(func() bool {
		return !s.observer.ChainParams().IsSupported
	})

	register(s.observer.ObserveInbound, "observer_inbound", optInboundInterval, optInboundSkipper)
	register(s.observer.ProcessInboundTrackers, "process_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(s.observer.CheckRPCStatus, "check_rpc_status")
	register(s.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)
	register(s.observer.ProcessOutboundTrackers, "process_outbound_trackers", optOutboundInterval, optOutboundSkipper)

	// CCTX scheduler (every musechain block)
	register(s.scheduleCCTX, "schedule_cctx", scheduler.BlockTicker(newBlockChan), optOutboundSkipper)

	return nil
}

// Stop stops all relevant tasks.
func (s *Sui) Stop() {
	s.observer.Logger().Chain.Info().Msg("stopping observer-signer")
	s.scheduler.StopGroup(s.group())
}

func (s *Sui) group() scheduler.Group {
	return scheduler.Group(fmt.Sprintf("sui:%d", s.Chain().ChainId))
}

// scheduleCCTX schedules outbound cross-chain transactions.
func (s *Sui) scheduleCCTX(ctx context.Context) error {
	if err := s.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	museBlock, delay, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get muse block from context")
	}

	time.Sleep(delay)

	cctxList, _, err := s.observer.MusecoreClient().ListPendingCCTX(ctx, s.observer.Chain())
	if err != nil {
		return errors.Wrap(err, "unable to list pending cctx")
	}

	// noop
	if len(cctxList) == 0 {
		return nil
	}

	var (
		// #nosec G115 always in range
		museHeight = uint64(museBlock.Block.Height)
		chainID    = s.observer.Chain().ChainId

		lookahead = s.observer.ChainParams().OutboundScheduleLookahead
		// #nosec G115 always in range
		lookback = uint64(float64(lookahead) * outboundLookbackFactor)

		firstNonce = cctxList[0].GetCurrentOutboundParam().TssNonce
		maxNonce   = firstNonce + lookback
	)

	for i := range cctxList {
		var (
			cctx           = cctxList[i]
			outboundID     = base.OutboundIDFromCCTX(cctx)
			outboundParams = cctx.GetCurrentOutboundParam()
			nonce          = outboundParams.TssNonce
		)

		switch {
		case int64(i) == lookahead:
			// take only first N cctxs
			return nil
		case outboundParams.ReceiverChainId != chainID:
			// should not happen
			s.outboundLogger(outboundID).Error().Msg("chain id mismatch")
			continue
		case nonce >= maxNonce:
			return fmt.Errorf("nonce %d is too high (%s). Earliest nonce %d", nonce, outboundID, firstNonce)
		case s.signer.IsOutboundActive(outboundID):
			// cctx is already being processed & broadcasted by signer
			continue
		case s.observer.OutboundCreated(nonce):
			// ProcessOutboundTrackers HAS fetched existing Sui outbound,
			// Let's report this by voting to musecore
			if err := s.observer.VoteOutbound(ctx, cctx); err != nil {
				s.outboundLogger(outboundID).Error().Err(err).Msg("VoteOutbound failed")
			}
			continue
		}

		// Here we have a cctx that needs to be scheduled. Let's invoke async operation.
		// - Signer will build, sign & broadcast the tx.
		// - It will also monitor Sui to report outbound tracker
		//   so we'd have a pair of (tss_nonce -> sui tx hash)
		// - Then this pair will be handled by ProcessOutboundTrackers -> OutboundCreated -> VoteOutbound
		bg.Work(ctx, func(ctx context.Context) error {
			if err := s.signer.ProcessCCTX(ctx, cctx, museHeight); err != nil {
				s.outboundLogger(outboundID).Error().Err(err).Msg("ProcessCCTX failed")
			}

			return nil
		})
	}

	return nil
}

func (s *Sui) updateChainParams(ctx context.Context) error {
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

	// note that address should be in format of `$packageID,$gatewayObjectID`
	if err := s.observer.Gateway().UpdateIDs(params.GatewayAddress); err != nil {
		return errors.Wrap(err, "unable to update gateway ids")
	}

	return nil
}

func (s *Sui) outboundLogger(id string) *zerolog.Logger {
	l := s.observer.Logger().Outbound.With().Str("outbound.id", id).Logger()

	return &l
}
