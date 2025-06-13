package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/x/observer/keeper"
	"github.com/RWAs-labs/muse/x/observer/types"
)

// Keeper Tests
func createTestKeygen(keeper *keeper.Keeper, ctx sdk.Context) types.Keygen {
	item := types.Keygen{
		BlockNumber: 10,
	}
	keeper.SetKeygen(ctx, item)
	return item
}

func TestKeygenGet(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	item := createTestKeygen(k, ctx)
	rst, found := k.GetKeygen(ctx)
	require.True(t, found)
	require.Equal(t, item, rst)
}
func TestKeygenRemove(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	createTestKeygen(k, ctx)
	k.RemoveKeygen(ctx)
	_, found := k.GetKeygen(ctx)
	require.False(t, found)
}
