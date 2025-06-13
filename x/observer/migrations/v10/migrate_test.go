package v10_test

import (
	"testing"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	v10 "github.com/RWAs-labs/muse/x/observer/migrations/v10"
	"github.com/RWAs-labs/muse/x/observer/types"
	"github.com/stretchr/testify/require"
)

func TestMigrateStore(t *testing.T) {
	t.Run("can migrate confirmation count", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain params
		testChainParams := getTestChainParams()
		k.SetChainParamsList(ctx, testChainParams)

		// ensure the chain params are set correctly
		oldChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, testChainParams, oldChainParams)

		// migrate the store
		err := v10.MigrateStore(ctx, *k)
		require.NoError(t, err)

		// ensure we still have same number of chain params after migration
		newChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, newChainParams.ChainParams, len(oldChainParams.ChainParams))

		// compare the old and new chain params
		for i, newParam := range newChainParams.ChainParams {
			oldParam := oldChainParams.ChainParams[i]

			// ensure the confirmation fields are set correctly
			require.Equal(t, oldParam.ConfirmationCount, newParam.ConfirmationParams.SafeInboundCount)
			require.Equal(t, oldParam.ConfirmationCount, newParam.ConfirmationParams.FastInboundCount)
			require.Equal(t, oldParam.ConfirmationCount, newParam.ConfirmationParams.SafeOutboundCount)
			require.Equal(t, oldParam.ConfirmationCount, newParam.ConfirmationParams.FastOutboundCount)

			// ensure nothing else has changed except the confirmation
			oldParam.ConfirmationParams = &types.ConfirmationParams{
				SafeInboundCount:  oldParam.ConfirmationCount,
				FastInboundCount:  oldParam.ConfirmationCount,
				SafeOutboundCount: oldParam.ConfirmationCount,
				FastOutboundCount: oldParam.ConfirmationCount,
			}
			require.Equal(t, newParam, oldParam)
		}
	})

	t.Run("migrate nothing if chain params not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// ensure no chain params are set
		allChainParams, found := k.GetChainParamsList(ctx)
		require.False(t, found)
		require.Empty(t, allChainParams.ChainParams)

		// migrate the store
		err := v10.MigrateStore(ctx, *k)
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)

		// ensure nothing has changed
		allChainParams, found = k.GetChainParamsList(ctx)
		require.False(t, found)
		require.Empty(t, allChainParams.ChainParams)
	})

	t.Run("migrate nothing if chain params list validation fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// get test chain params
		testChainParams := getTestChainParams()

		// make the first chain params invalid
		testChainParams.ChainParams[0].InboundTicker = 0

		// set chain params
		k.SetChainParamsList(ctx, testChainParams)

		// ensure the chain params are set correctly
		oldChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, testChainParams, oldChainParams)

		// migrate the store
		err := v10.MigrateStore(ctx, *k)
		require.ErrorIs(t, err, types.ErrInvalidChainParams)

		// ensure nothing has changed
		newChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, oldChainParams, newChainParams)
	})
}

// makeChainParamsEmptyConfirmation creates a sample chain params with empty confirmation
func makeChainParamsEmptyConfirmation(chainID int64, confirmationCount uint64) *types.ChainParams {
	chainParams := sample.ChainParams(chainID)
	chainParams.ConfirmationCount = confirmationCount
	chainParams.ConfirmationParams = nil
	return chainParams
}

// getTestChainParams returns a list of chain params for testing
func getTestChainParams() types.ChainParamsList {
	return types.ChainParamsList{
		ChainParams: []*types.ChainParams{
			makeChainParamsEmptyConfirmation(1, 14),
			makeChainParamsEmptyConfirmation(56, 20),
			makeChainParamsEmptyConfirmation(8332, 3),
			makeChainParamsEmptyConfirmation(7000, 0),
			makeChainParamsEmptyConfirmation(137, 200),
			makeChainParamsEmptyConfirmation(8453, 90),
			makeChainParamsEmptyConfirmation(900, 32),
		},
	}
}
