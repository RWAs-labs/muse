package keeper_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestKeeper_AddMuseAccounting(t *testing.T) {
	t.Run("should add aborted muse amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		originalAmount := sdkmath.NewUint(rand.Uint64())
		k.SetMuseAccounting(ctx, types.MuseAccounting{
			AbortedMuseAmount: originalAmount,
		})
		val, found := k.GetMuseAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.AbortedMuseAmount)
		addAmount := sdkmath.NewUint(rand.Uint64())
		k.AddMuseAbortedAmount(ctx, addAmount)
		val, found = k.GetMuseAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount.Add(addAmount), val.AbortedMuseAmount)
	})

	t.Run("should add aborted muse amount if accounting not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		originalAmount := sdkmath.NewUint(0)
		_, found := k.GetMuseAccounting(ctx)
		require.False(t, found)
		addAmount := sdkmath.NewUint(rand.Uint64())
		k.AddMuseAbortedAmount(ctx, addAmount)
		val, found := k.GetMuseAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount.Add(addAmount), val.AbortedMuseAmount)
	})

	t.Run("cant find aborted amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		val, found := k.GetMuseAccounting(ctx)
		require.False(t, found)
		require.Equal(t, types.MuseAccounting{}, val)
	})

	t.Run("add very high muse amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		highAmount := sdkmath.NewUintFromString("100000000000000000000000000000000000000000000000")
		k.SetMuseAccounting(ctx, types.MuseAccounting{
			AbortedMuseAmount: highAmount,
		})
		val, found := k.GetMuseAccounting(ctx)
		require.True(t, found)
		require.Equal(t, highAmount, val.AbortedMuseAmount)
	})

}

func TestKeeper_RemoveMuseAbortedAmount(t *testing.T) {
	t.Run("should remove aborted muse amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		originalAmount := sdkmath.NewUintFromString("100000000000000000000000000000000000000000000000")
		k.SetMuseAccounting(ctx, types.MuseAccounting{
			AbortedMuseAmount: originalAmount,
		})
		val, found := k.GetMuseAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.AbortedMuseAmount)
		removeAmount := originalAmount.Sub(sdkmath.NewUintFromString("10000000000000000000000000000000000000000000000"))
		err := k.RemoveMuseAbortedAmount(ctx, removeAmount)
		require.NoError(t, err)
		val, found = k.GetMuseAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount.Sub(removeAmount), val.AbortedMuseAmount)
	})
	t.Run("fail remove aborted muse amount if accounting not set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		err := k.RemoveMuseAbortedAmount(ctx, sdkmath.OneUint())
		require.ErrorIs(t, err, types.ErrUnableToFindMuseAccounting)
	})
	t.Run("fail remove aborted muse amount if insufficient amount", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		originalAmount := sdkmath.NewUint(100)
		k.SetMuseAccounting(ctx, types.MuseAccounting{
			AbortedMuseAmount: originalAmount,
		})
		val, found := k.GetMuseAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.AbortedMuseAmount)
		removeAmount := originalAmount.Add(sdkmath.NewUint(500))
		err := k.RemoveMuseAbortedAmount(ctx, removeAmount)
		require.ErrorIs(t, err, types.ErrInsufficientMuseAmount)
		val, found = k.GetMuseAccounting(ctx)
		require.True(t, found)
		require.Equal(t, originalAmount, val.AbortedMuseAmount)
	})
}
