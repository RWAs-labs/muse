package orchestrator

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/pkg/chains"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

type Musecore interface {
	GetBlockHeight(ctx context.Context) (int64, error)
	GetUpgradePlan(ctx context.Context) (*upgradetypes.Plan, error)
	GetSupportedChains(ctx context.Context) ([]chains.Chain, error)
	GetAdditionalChains(ctx context.Context) ([]chains.Chain, error)
	GetCrosschainFlags(ctx context.Context) (observertypes.CrosschainFlags, error)
	GetChainParams(ctx context.Context) ([]*observertypes.ChainParams, error)
	GetTSS(ctx context.Context) (observertypes.TSS, error)
	GetKeyGen(ctx context.Context) (observertypes.Keygen, error)
	GetOperationalFlags(ctx context.Context) (observertypes.OperationalFlags, error)
}

var ErrUpgradeRequired = errors.New("upgrade required")

// UpdateAppContext fetches latest data from Musecore and updates the AppContext.
// Also detects if an upgrade is required. If an upgrade is required, it returns ErrUpgradeRequired.
func UpdateAppContext(ctx context.Context, app *zctx.AppContext, zc Musecore, logger zerolog.Logger) error {
	bn, err := zc.GetBlockHeight(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get muse block height")
	}

	if err = checkForMusecoreUpgrade(ctx, bn, zc); err != nil {
		return err
	}

	supportedChains, err := zc.GetSupportedChains(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch supported chains")
	}

	additionalChains, err := zc.GetAdditionalChains(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch additional chains")
	}

	chainParams, err := zc.GetChainParams(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch chain params")
	}

	crosschainFlags, err := zc.GetCrosschainFlags(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch crosschain flags")
	}

	operationalFlags, err := zc.GetOperationalFlags(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to fetch operational flags")
	}

	freshParams := make(map[int64]*observertypes.ChainParams, len(chainParams))

	// check and update chain params for each chain
	// Note that we are EXCLUDING MuseChain from the chainParams if it's present
	for i := range chainParams {
		cp := chainParams[i]

		if !cp.IsSupported {
			logger.Warn().Int64("chain.id", cp.ChainId).Msg("Skipping unsupported chain")
			continue
		}

		if chains.IsMuseChain(cp.ChainId, nil) {
			continue
		}

		if err := cp.Validate(); err != nil {
			logger.Warn().Err(err).Int64("chain.id", cp.ChainId).Msg("Skipping invalid chain params")
			continue
		}

		freshParams[cp.ChainId] = cp
	}

	return app.Update(
		supportedChains,
		additionalChains,
		freshParams,
		crosschainFlags,
		operationalFlags,
	)
}

// returns an error if an upgrade is required
func checkForMusecoreUpgrade(ctx context.Context, museHeight int64, zc Musecore) error {
	plan, err := zc.GetUpgradePlan(ctx)
	switch {
	case err != nil:
		return errors.Wrap(err, "unable to get upgrade plan")
	case plan == nil:
		// no upgrade planned
		return nil
	}

	upgradeHeight := plan.Height

	// We can return an error in a few blocks ahead.
	// It's okay because the ticker might have an interval longer than 1 block (~5s).
	//
	// Example: if an upgrade is on block #102, we can return an error on block #100, #101, #102, ...
	// Note that tha plan is deleted from musecore after the upgrade block.
	const upgradeBlockBuffer = 2

	if (upgradeHeight - museHeight) <= upgradeBlockBuffer {
		return errors.Wrapf(ErrUpgradeRequired, "current height: %d, upgrade height: %d", museHeight, upgradeHeight)
	}

	return nil
}
