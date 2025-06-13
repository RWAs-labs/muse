// Package context provides global app context for MuseClient
package context

import (
	"fmt"
	"slices"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"

	"github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/pkg/chains"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// AppContext represents application (museclient) context.
type AppContext struct {
	// config is the config of the app
	config config.Config

	// chainRegistry is a registry of supported chains
	chainRegistry *ChainRegistry

	// crosschainFlags is the current crosschain flags state
	crosschainFlags  observertypes.CrosschainFlags
	operationalFlags observertypes.OperationalFlags

	// logger is the logger of the app
	logger zerolog.Logger

	mu sync.RWMutex
}

// New creates and returns new empty AppContext
func New(cfg config.Config, relayerKeyPasswords map[string]string, logger zerolog.Logger) *AppContext {
	return &AppContext{
		config:          cfg,
		chainRegistry:   NewChainRegistry(relayerKeyPasswords),
		crosschainFlags: observertypes.CrosschainFlags{},
		logger:          logger.With().Str("module", "appcontext").Logger(),
	}
}

// Config returns the config of the app
func (a *AppContext) Config() config.Config {
	return a.config
}

// GetChain returns the chain by ID.
func (a *AppContext) GetChain(chainID int64) (Chain, error) {
	return a.chainRegistry.Get(chainID)
}

// ListChainIDs returns the list of existing chain ids in the registry.
func (a *AppContext) ListChainIDs() []int64 {
	return a.chainRegistry.ChainIDs()
}

// ListChains returns the list of existing chains in the registry.
func (a *AppContext) ListChains() []Chain {
	return a.chainRegistry.All()
}

// FilterChains returns the list of chains that satisfy the filter
func (a *AppContext) FilterChains(filter func(Chain) bool) []Chain {
	var (
		all = a.ListChains()
		out = make([]Chain, 0, len(all))
	)

	for _, chain := range all {
		if filter(chain) {
			out = append(out, chain)
		}
	}

	return out
}

// IsOutboundObservationEnabled returns true if outbound flag is enabled
func (a *AppContext) IsOutboundObservationEnabled() bool {
	return a.GetCrossChainFlags().IsOutboundEnabled
}

// IsInboundObservationEnabled returns true if inbound flag is enabled
func (a *AppContext) IsInboundObservationEnabled() bool {
	return a.GetCrossChainFlags().IsInboundEnabled
}

// GetCrossChainFlags returns crosschain flags
func (a *AppContext) GetCrossChainFlags() observertypes.CrosschainFlags {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.crosschainFlags
}

// GetOperationalFlags returns operational flags
func (a *AppContext) GetOperationalFlags() observertypes.OperationalFlags {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.operationalFlags
}

// Update updates AppContext and params for all chains
// this must be the ONLY function that writes to AppContext
func (a *AppContext) Update(
	freshChains, additionalChains []chains.Chain,
	freshChainParams map[int64]*observertypes.ChainParams,
	crosschainFlags observertypes.CrosschainFlags,
	operationalFlags observertypes.OperationalFlags,
) error {
	// some sanity checks
	switch {
	case len(freshChains) == 0:
		return fmt.Errorf("no chains present")
	case len(freshChainParams) == 0:
		return fmt.Errorf("no chain params present")
	case len(additionalChains) > 0:
		for _, c := range additionalChains {
			if !c.IsExternal {
				return fmt.Errorf("additional chain %d is not external", c.ChainId)
			}
		}
	}

	err := a.updateChainRegistry(freshChains, additionalChains, freshChainParams)
	if err != nil {
		return errors.Wrap(err, "unable to update chain registry")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.crosschainFlags = crosschainFlags
	a.operationalFlags = operationalFlags

	return nil
}

// updateChainRegistry updates the chain registry with fresh chains and chain params.
// Note that there's an edge-case for MuseChain itself because we WANT to have it in chains list,
// but it doesn't have chain params.
func (a *AppContext) updateChainRegistry(
	freshChains []chains.Chain,
	additionalChains []chains.Chain,
	freshChainParams map[int64]*observertypes.ChainParams,
) error {
	var museChainID int64

	// 1. build map[chainId]Chain
	freshChainsByID := make(map[int64]chains.Chain, len(freshChains)+len(additionalChains))
	for _, c := range freshChains {
		freshChainsByID[c.ChainId] = c

		if isMuse(c.ChainId) && museChainID == 0 {
			museChainID = c.ChainId
		}
	}

	for _, c := range additionalChains {
		// shouldn't happen, but just in case
		if _, found := freshChainsByID[c.ChainId]; found {
			continue
		}

		freshChainsByID[c.ChainId] = c
	}

	var (
		freshChainIDs    = maps.Keys(freshChainsByID)
		existingChainIDs = a.chainRegistry.ChainIDs()
	)

	slices.Sort(freshChainIDs)
	slices.Sort(existingChainIDs)

	// 2. Compare existing chains with fresh ones
	if len(existingChainIDs) > 0 && !slicesEqual(existingChainIDs, freshChainIDs) {
		a.logger.Warn().
			Ints64("chains.current", existingChainIDs).
			Ints64("chains.new", freshChainIDs).
			Msg("Chain list changed at the runtime!")
	}

	// Log warn if somehow chain doesn't chainParam
	for _, chainID := range freshChainIDs {
		if _, ok := freshChainParams[chainID]; !ok && !isMuse(chainID) {
			a.logger.Warn().
				Int64("chain.id", chainID).
				Msg("Chain doesn't have according ChainParams present. Skipping.")
		}
	}

	// 3. If we have muse chain, we want to force "fake" chainParams for it
	if museChainID != 0 {
		freshChainParams[museChainID] = museObserverChainParams(museChainID)
	}

	// 3. Update chain registry
	// okay, let's update the chains.
	// Set() ensures that chain, chainID, and params are consistent and chain is not muse + chain is supported
	for chainID, params := range freshChainParams {
		chain, ok := freshChainsByID[chainID]
		if !ok {
			return fmt.Errorf("unable to locate fresh chain %d based on chain params", chainID)
		}

		if !isMuse(chainID) {
			if err := params.Validate(); err != nil {
				return errors.Wrapf(err, "invalid chain params for chain %d", chainID)
			}
		}

		if err := a.chainRegistry.Set(chainID, &chain, params); err != nil {
			return errors.Wrap(err, "unable to set chain in the registry")
		}
	}

	a.chainRegistry.SetAdditionalChains(additionalChains)

	toBeDeleted, _ := lo.Difference(existingChainIDs, freshChainIDs)
	if len(toBeDeleted) > 0 {
		a.logger.Warn().
			Ints64("chains.deleted", toBeDeleted).
			Msg("Deleting chains that are no longer relevant")

		a.chainRegistry.Delete(toBeDeleted...)
	}

	return nil
}

func isMuse(chainID int64) bool {
	return chains.IsMuseChain(chainID, nil)
}

// museObserverChainParams returns "fake" chain params because
// actually chainParams is a concept of observer
func museObserverChainParams(chainID int64) *observertypes.ChainParams {
	return &observertypes.ChainParams{ChainId: chainID, IsSupported: true}
}

func slicesEqual[T constraints.Ordered](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
