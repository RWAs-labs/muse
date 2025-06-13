package orchestrator

import (
	"testing"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	"github.com/RWAs-labs/muse/museclient/testutils/testlog"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/ptr"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
	"github.com/stretchr/testify/require"
)

func Test_UpdateAppContext(t *testing.T) {
	var (
		eth       = chains.Ethereum
		ethParams = mocks.MockChainParams(eth.ChainId, 100)

		btc       = chains.BitcoinMainnet
		btcParams = mocks.MockChainParams(btc.ChainId, 100)
	)

	t.Run("Updates app context", func(t *testing.T) {
		var (
			logger                 = testlog.New(t).Logger
			chainList, chainParams = parseChainsWithParams(t, eth, ethParams)
			ctx, app               = newAppContext(t, logger, chainList, chainParams)
			musecore               = mocks.NewMusecoreClient(t)
		)

		// Given musecore client that has eth and btc chains
		newChains := []chains.Chain{eth, btc}
		newParams := []*observertypes.ChainParams{&ethParams, &btcParams}
		ccFlags := observertypes.CrosschainFlags{
			IsInboundEnabled:  true,
			IsOutboundEnabled: true,
		}
		opFlags := observertypes.OperationalFlags{
			RestartHeight:         123,
			SignerBlockTimeOffset: ptr.Ptr(time.Second),
			MinimumVersion:        "",
		}

		on(musecore, "GetBlockHeight", 1).Return(int64(123), nil)
		on(musecore, "GetUpgradePlan", 1).Return(nil, nil)
		on(musecore, "GetSupportedChains", 1).Return(newChains, nil)
		on(musecore, "GetAdditionalChains", 1).Return(nil, nil)
		on(musecore, "GetChainParams", 1).Return(newParams, nil)
		on(musecore, "GetCrosschainFlags", 1).Return(ccFlags, nil)
		on(musecore, "GetOperationalFlags", 1).Return(opFlags, nil)

		// ACT
		err := UpdateAppContext(ctx, app, musecore, logger)

		// ASSERT
		require.NoError(t, err)

		// New chains should be added
		_, err = app.GetChain(btc.ChainId)
		require.NoError(t, err)

		// Check OP flags
		require.Equal(t, opFlags.RestartHeight, app.GetOperationalFlags().RestartHeight)
	})

	t.Run("Upgrade plan detected", func(t *testing.T) {
		// ARRANGE
		var (
			logger                 = testlog.New(t).Logger
			chainList, chainParams = parseChainsWithParams(t, eth, ethParams)
			ctx, app               = newAppContext(t, logger, chainList, chainParams)
			musecore               = mocks.NewMusecoreClient(t)
		)

		on(musecore, "GetBlockHeight", 1).Return(int64(123), nil)
		on(musecore, "GetUpgradePlan", 1).Return(&upgradetypes.Plan{
			Name:   "hello",
			Height: 124,
		}, nil)

		// ACT
		err := UpdateAppContext(ctx, app, musecore, logger)

		// ASSERT
		require.ErrorIs(t, err, ErrUpgradeRequired)
	})
}
