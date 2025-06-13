package evm

import (
	"context"
	"fmt"
	"time"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/evm/observer"
	"github.com/RWAs-labs/muse/museclient/chains/evm/signer"
	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/scheduler"
	"github.com/RWAs-labs/muse/pkg/ticker"
)

type EVM struct {
	scheduler *scheduler.Scheduler
	observer  *observer.Observer
	signer    *signer.Signer
}

const (
	// outboundLookBackFactor is the factor to determine how many nonces to look back for pending cctxs
	// For example, give OutboundScheduleLookahead of 120, pending NonceLow of 1000 and factor of 1.0,
	// the scheduler need to be able to pick up and schedule any pending cctx with nonce < 880 (1000 - 120 * 1.0)
	// NOTE: 1.0 means look back the same number of cctxs as we look ahead
	outboundLookBackFactor = 1.0
)

func New(scheduler *scheduler.Scheduler, observer *observer.Observer, signer *signer.Signer) *EVM {
	return &EVM{
		scheduler: scheduler,
		observer:  observer,
		signer:    signer,
	}
}

func (e *EVM) Chain() chains.Chain {
	return e.observer.Chain()
}

func (e *EVM) Start(ctx context.Context) error {
	if ok := e.observer.Observer.Start(); !ok {
		return errors.New("observer is already started")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get app from context")
	}

	newBlockChan, err := e.observer.MusecoreClient().NewBlockSubscriber(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to create new block subscriber")
	}

	optInboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(e.observer.ChainParams().InboundTicker)
	})

	optGasInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(e.observer.ChainParams().GasPriceTicker)
	})

	optOutboundInterval := scheduler.IntervalUpdater(func() time.Duration {
		return ticker.DurationFromUint64Seconds(e.observer.ChainParams().OutboundTicker)
	})

	optInboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsInboundObservationEnabled()
	})

	optOutboundSkipper := scheduler.Skipper(func() bool {
		return !app.IsOutboundObservationEnabled()
	})

	optGenericSkipper := scheduler.Skipper(func() bool {
		return !e.observer.ChainParams().IsSupported
	})

	register := func(exec scheduler.Executable, name string, opts ...scheduler.Opt) {
		opts = append([]scheduler.Opt{
			scheduler.GroupName(e.group()),
			scheduler.Name(name),
		}, opts...)

		e.scheduler.Register(ctx, exec, opts...)
	}

	// Observers
	register(e.observer.ObserveInbound, "observe_inbound", optInboundInterval, optInboundSkipper)
	register(e.observer.ProcessInboundTrackers, "process_inbound_trackers", optInboundInterval, optInboundSkipper)
	register(e.observer.PostGasPrice, "post_gas_price", optGasInterval, optGenericSkipper)
	register(e.observer.CheckRPCStatus, "check_rpc_status")
	register(e.observer.ProcessOutboundTrackers, "process_outbound_trackers", optOutboundInterval, optOutboundSkipper)

	// CCTX Scheduler
	register(e.scheduleCCTX, "schedule_cctx", scheduler.BlockTicker(newBlockChan), optOutboundSkipper)

	return nil
}

func (e *EVM) Stop() {
	e.observer.Logger().Chain.Info().Msg("stopping observer-signer")
	e.scheduler.StopGroup(e.group())
}

func (e *EVM) group() scheduler.Group {
	return scheduler.Group(
		fmt.Sprintf("evm:%d", e.observer.Chain().ChainId),
	)
}

// scheduleCCTX schedules outbound transactions on each muse block.
func (e *EVM) scheduleCCTX(ctx context.Context) error {
	if err := e.updateChainParams(ctx); err != nil {
		return errors.Wrap(err, "unable to update chain params")
	}

	museBlock, delay, err := scheduler.BlockFromContextWithDelay(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get muse block from context")
	}

	time.Sleep(delay)

	var (
		chain   = e.observer.Chain()
		chainID = chain.ChainId

		// #nosec G115 always in range
		museHeight = uint64(museBlock.Block.Height)

		lookahead = e.observer.ChainParams().OutboundScheduleLookahead

		// #nosec G115 positive
		scheduleInterval = uint64(e.observer.ChainParams().OutboundScheduleInterval)

		// for critical pending outbound we reduce re-try interval
		criticalInterval = uint64(10)

		// for non-critical pending outbound we increase re-try interval
		nonCriticalInterval = scheduleInterval * 2

		// #nosec G115 always in range
		outboundScheduleLookBack = uint64(float64(lookahead) * outboundLookBackFactor)
	)

	cctxList, _, err := e.observer.MusecoreClient().ListPendingCCTX(ctx, chain)
	if err != nil {
		return errors.Wrap(err, "unable to list pending cctx")
	}

	trackerSet, err := e.getTrackerSet(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get tracker set")
	}

	for idx, cctx := range cctxList {
		var (
			params     = cctx.GetCurrentOutboundParam()
			nonce      = params.TssNonce
			outboundID = base.OutboundIDFromCCTX(cctx)
		)

		switch {
		case params.ReceiverChainId != chainID:
			e.outboundLogger(outboundID).Error().Msg("chain id mismatch")
			continue
		case params.TssNonce > cctxList[0].GetCurrentOutboundParam().TssNonce+outboundScheduleLookBack:
			return fmt.Errorf(
				"nonce %d is too high (%s). Earliest nonce %d",
				params.TssNonce,
				outboundID,
				cctxList[0].GetCurrentOutboundParam().TssNonce,
			)
		}

		// vote outbound if it's already confirmed
		continueKeysign, err := e.observer.VoteOutboundIfConfirmed(ctx, cctx)
		switch {
		case err != nil:
			e.outboundLogger(outboundID).Error().Err(err).Msg("Schedule CCTX: VoteOutboundIfConfirmed failed")
			continue
		case !continueKeysign:
			e.outboundLogger(outboundID).Info().Msg("Schedule CCTX: outbound already processed")
			continue
		case e.signer.IsOutboundActive(outboundID):
			// outbound is already being processed
			continue
		}

		// determining critical outbound; if it satisfies following criteria
		// 1. it's the first pending outbound for this chain
		// 2. the following 5 nonces have been in tracker
		if nonce%criticalInterval == museHeight%criticalInterval {
			count := 0
			for i := nonce + 1; i <= nonce+10; i++ {
				if _, found := trackerSet[i]; found {
					count++
				}
			}
			if count >= 5 {
				scheduleInterval = criticalInterval
			}
		}

		// if it's already in tracker, we increase re-try interval
		if _, ok := trackerSet[nonce]; ok {
			scheduleInterval = nonCriticalInterval
		}

		// otherwise, the normal interval is used
		if nonce%scheduleInterval == museHeight%scheduleInterval {
			go e.signer.TryProcessOutbound(
				ctx,
				cctx,
				e.observer.MusecoreClient(),
				museHeight,
			)
		}

		// #nosec G115 always in range
		// only look at 'lookahead' cctxs per chain
		if int64(idx) >= lookahead-1 {
			break
		}
	}

	return nil
}

func (e *EVM) updateChainParams(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	chain, err := app.GetChain(e.observer.Chain().ChainId)
	if err != nil {
		return errors.Wrap(err, "unable to get chain")
	}

	params := chain.Params()

	// Update chain params
	e.observer.SetChainParams(*params)

	// Update muse connector, ERC20 custody, and gateway addresses
	e.signer.SetMuseConnectorAddress(eth.HexToAddress(params.ConnectorContractAddress))
	e.signer.SetERC20CustodyAddress(eth.HexToAddress(params.Erc20CustodyContractAddress))
	e.signer.SetGatewayAddress(params.GatewayAddress)

	return nil
}

func (e *EVM) getTrackerSet(ctx context.Context) (map[uint64]struct{}, error) {
	chainID := e.observer.Chain().ChainId

	trackers, err := e.observer.MusecoreClient().GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get outbound trackers")
	}

	set := make(map[uint64]struct{})

	for _, tracker := range trackers {
		set[tracker.Nonce] = struct{}{}
	}

	return set, nil
}

func (e *EVM) outboundLogger(id string) *zerolog.Logger {
	l := e.observer.Logger().Outbound.With().Str("outbound.id", id).Logger()

	return &l
}
