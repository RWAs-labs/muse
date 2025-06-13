package base

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/compliance"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// Signer is the base structure for grouping the common logic between chain signers.
// The common logic includes: chain, chainParams, contexts, tss, metrics, loggers etc.
type Signer struct {
	// chain contains static information about the external chain
	chain chains.Chain

	// tss is the TSS signer
	tss interfaces.TSSSigner

	// logger contains the loggers used by signer
	logger Logger

	// outboundBeingReported is a map of outbound being reported to tracker
	outboundBeingReported map[string]bool

	activeOutbounds map[string]time.Time

	// mu protects fields from concurrent access
	// Note: base signer simply provides the mutex. It's the sub-struct's responsibility to use it to be thread-safe
	mu sync.Mutex
}

// NewSigner creates a new base signer.
func NewSigner(chain chains.Chain, tss interfaces.TSSSigner, logger Logger) *Signer {
	withLogFields := func(log zerolog.Logger) zerolog.Logger {
		return log.With().
			Int64(logs.FieldChain, chain.ChainId).
			Str(logs.FieldModule, "signer").
			Logger()
	}

	return &Signer{
		chain:                 chain,
		tss:                   tss,
		outboundBeingReported: make(map[string]bool),
		activeOutbounds:       make(map[string]time.Time),
		logger: Logger{
			Std:        withLogFields(logger.Std),
			Compliance: withLogFields(logger.Compliance),
		},
	}
}

// Chain returns the chain for the signer.
func (s *Signer) Chain() chains.Chain {
	return s.chain
}

// TSS returns the tss signer for the signer.
func (s *Signer) TSS() interfaces.TSSSigner {
	return s.tss
}

// Logger returns the logger for the signer.
func (s *Signer) Logger() *Logger {
	return &s.logger
}

// SetBeingReportedFlag sets the outbound as being reported if not already set.
// Returns true if the outbound is already being reported.
// This method is used by outbound tracker reporter to avoid repeated reporting of same hash.
func (s *Signer) SetBeingReportedFlag(hash string) (alreadySet bool) {
	s.Lock()
	defer s.Unlock()

	alreadySet = s.outboundBeingReported[hash]
	if !alreadySet {
		// mark as being reported
		s.outboundBeingReported[hash] = true
	}
	return
}

// ClearBeingReportedFlag clears the being reported flag for the outbound.
func (s *Signer) ClearBeingReportedFlag(hash string) {
	s.Lock()
	defer s.Unlock()
	delete(s.outboundBeingReported, hash)
}

// Exported for unit tests

// GetReportedTxList returns a list of outboundHash being reported.
// TODO: investigate pointer usage
// https://github.com/RWAs-labs/muse/issues/2084
func (s *Signer) GetReportedTxList() *map[string]bool {
	return &s.outboundBeingReported
}

// Lock locks the signer.
func (s *Signer) Lock() {
	s.mu.Lock()
}

// Unlock unlocks the signer.
func (s *Signer) Unlock() {
	s.mu.Unlock()
}

// MarkOutbound marks the outbound as active.
func (s *Signer) MarkOutbound(outboundID string, active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	startedAt, found := s.activeOutbounds[outboundID]

	switch {
	case active == found:
		// noop
	case active:
		now := time.Now().UTC()
		s.activeOutbounds[outboundID] = now

		s.logger.Std.Info().
			Bool("outbound.active", active).
			Str("outbound.id", outboundID).
			Time("outbound.timestamp", now).
			Int("outbound.total", len(s.activeOutbounds)).
			Msg("MarkOutbound")
	default:
		timeTaken := time.Since(startedAt)

		s.logger.Std.Info().
			Bool("outbound.active", active).
			Str("outbound.id", outboundID).
			Float64("outbound.time_taken", timeTaken.Seconds()).
			Int("outbound.total", len(s.activeOutbounds)).
			Msg("MarkOutbound")

		delete(s.activeOutbounds, outboundID)
	}
}

func (s *Signer) IsOutboundActive(outboundID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.activeOutbounds[outboundID]
	return found
}

// PassesCompliance checks if the cctx passes the compliance check and prints compliance log.
func (s *Signer) PassesCompliance(cctx *types.CrossChainTx) bool {
	if !compliance.IsCCTXRestricted(cctx) {
		return true
	}

	params := cctx.GetCurrentOutboundParam()

	compliance.PrintComplianceLog(
		s.Logger().Std,
		s.Logger().Compliance,
		true,
		s.Chain().ChainId,
		cctx.Index,
		cctx.InboundParams.Sender,
		params.Receiver,
		params.CoinType.String(),
	)

	return false
}

// OutboundID returns the outbound ID.
func OutboundID(index string, receiverChainID int64, nonce uint64) string {
	return fmt.Sprintf("%s-%d-%d", index, receiverChainID, nonce)
}

// OutboundIDFromCCTX returns the outbound ID from the cctx.
func OutboundIDFromCCTX(cctx *types.CrossChainTx) string {
	index, params := cctx.GetIndex(), cctx.GetCurrentOutboundParam()
	return OutboundID(index, params.ReceiverChainId, params.TssNonce)
}
