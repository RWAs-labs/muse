package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/nullify"
	"github.com/RWAs-labs/muse/x/crosschain/keeper"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// Keeper Tests
func createNOutboundTracker(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.OutboundTracker {
	items := make([]types.OutboundTracker, n)
	for i := range items {
		items[i].ChainId = int64(i)
		items[i].Nonce = uint64(i)
		items[i].Index = fmt.Sprintf("%d-%d", items[i].ChainId, items[i].Nonce)

		keeper.SetOutboundTracker(ctx, items[i])
	}
	return items
}

func TestOutboundTrackerGet(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNOutboundTracker(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetOutboundTracker(ctx,
			item.ChainId,
			item.Nonce,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestOutboundTrackerRemove(t *testing.T) {
	t.Run("Remove tracker if it exists", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		items := createNOutboundTracker(keeper, ctx, 10)
		for _, item := range items {
			keeper.RemoveOutboundTrackerFromStore(ctx,
				item.ChainId,
				item.Nonce,
			)
			_, found := keeper.GetOutboundTracker(ctx,
				item.ChainId,
				item.Nonce,
			)
			require.False(t, found)
		}
	})

	t.Run("Do nothing if tracker doesn't exist", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		require.NotPanics(t, func() {
			keeper.RemoveOutboundTrackerFromStore(ctx, 1, 1)
		})
	})

}

func TestOutboundTrackerGetAll(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNOutboundTracker(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllOutboundTracker(ctx)),
	)
}
