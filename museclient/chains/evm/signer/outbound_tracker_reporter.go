// Package signer implements the ChainSigner interface for EVM chains
package signer

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/chains/evm/common"
	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/metrics"
	"github.com/RWAs-labs/muse/pkg/bg"
	crosschainkeeper "github.com/RWAs-labs/muse/x/crosschain/keeper"
)

// reportToOutboundTracker reports outboundHash to tracker only when tx receipt is available
func (signer *Signer) reportToOutboundTracker(
	ctx context.Context,
	musecoreClient interfaces.MusecoreClient,
	chainID int64,
	nonce uint64,
	outboundHash string,
	logger zerolog.Logger,
) {
	// prepare logger
	logger = logger.With().
		Str(logs.FieldMethod, "reportToOutboundTracker").
		Int64(logs.FieldChain, chainID).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, outboundHash).
		Logger()

	// set being reported flag to avoid duplicate reporting
	alreadySet := signer.SetBeingReportedFlag(outboundHash)
	if alreadySet {
		logger.Info().Msg("outbound is being reported to tracker")
		return
	}

	// launch a goroutine to monitor tx confirmation status
	bg.Work(ctx, func(ctx context.Context) error {
		metrics.NumTrackerReporters.WithLabelValues(signer.Chain().Name).Inc()

		defer func() {
			metrics.NumTrackerReporters.WithLabelValues(signer.Chain().Name).Dec()
			signer.ClearBeingReportedFlag(outboundHash)
		}()

		// try monitoring tx inclusion status for 20 minutes
		tStart := time.Now()
		for {
			// take a rest between each check
			time.Sleep(10 * time.Second)

			// give up (forget about the tx) after 20 minutes of monitoring, there are 2 reasons:
			// 1. the gas stability pool should have kicked in and replaced the tx by then.
			// 2. even if there is a chance that the tx is included later, most likely it's going to be a false tx hash (either replaced or dropped).
			// 3. we prefer missed tx hash over potentially invalid txhash.
			if time.Since(tStart) > common.OutboundInclusionTimeout {
				logger.Info().Msgf("timeout waiting outbound inclusion")
				return nil
			}

			// stop if the CCTX is already finalized for optimization purposes:
			// 1. all monitoring goroutines should stop and release resources if the CCTX is finalized
			// 2. especially reduces the lifetime of goroutines that monitor "nonce too low" tx hashes
			cctx, err := musecoreClient.GetCctxByNonce(ctx, chainID, nonce)
			if err != nil {
				logger.Err(err).Msg("unable to query cctx from musecore")
			} else if !crosschainkeeper.IsPending(cctx) {
				logger.Info().Msg("cctx is already finalized")
				return nil
			}

			// check tx confirmation status
			confirmed, err := signer.client.IsTxConfirmed(ctx, outboundHash, common.ReorgProtectBlockCount)
			if err != nil {
				logger.Err(err).Msg("unable to check confirmation status of outbound")
				continue
			}
			if !confirmed {
				continue
			}

			// report outbound hash to tracker
			museHash, err := musecoreClient.PostOutboundTracker(ctx, chainID, nonce, outboundHash)
			if err != nil {
				logger.Err(err).Msg("error adding outbound to tracker")
			} else if museHash != "" {
				logger.Info().Msgf("added outbound to tracker; muse txhash %s", museHash)
			} else {
				// exit goroutine until the tracker contains the hash (reported by either this or other signers)
				logger.Info().Msg("outbound now exists in tracker")
				return nil
			}
		}
	}, bg.WithName("TrackerReporterEVM"), bg.WithLogger(logger))
}
