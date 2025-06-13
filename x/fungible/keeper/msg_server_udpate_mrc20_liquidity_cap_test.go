package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/keeper"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestMsgServer_UpdateMRC20LiquidityCap(t *testing.T) {
	t.Run("can update the liquidity cap of mrc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		coinAddress := sample.EthAddress().String()

		foreignCoin := sample.ForeignCoins(t, coinAddress)
		foreignCoin.LiquidityCap = math.Uint{}
		k.SetForeignCoins(ctx, foreignCoin)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		// can update liquidity cap
		msg := types.NewMsgUpdateMRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(42),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateMRC20LiquidityCap(ctx, msg)
		require.NoError(t, err)

		coin, found := k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(t, coin.LiquidityCap.Equal(math.NewUint(42)), "invalid liquidity cap", coin.LiquidityCap.String())

		// can update liquidity cap again
		msg = types.NewMsgUpdateMRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(4200000),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateMRC20LiquidityCap(ctx, msg)
		require.NoError(t, err)

		coin, found = k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(
			t,
			coin.LiquidityCap.Equal(math.NewUint(4200000)),
			"invalid liquidity cap",
			coin.LiquidityCap.String(),
		)

		// can set liquidity cap to 0
		msg = types.NewMsgUpdateMRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(0),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateMRC20LiquidityCap(ctx, msg)
		require.NoError(t, err)

		coin, found = k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(t, coin.LiquidityCap.Equal(math.ZeroUint()), "invalid liquidity cap", coin.LiquidityCap.String())

		// can set liquidity cap to nil
		msg = types.NewMsgUpdateMRC20LiquidityCap(
			admin,
			coinAddress,
			math.Uint{},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateMRC20LiquidityCap(ctx, msg)
		require.NoError(t, err)

		coin, found = k.GetForeignCoins(ctx, coinAddress)
		require.True(t, found)
		require.True(t, coin.LiquidityCap.Equal(math.ZeroUint()), "invalid liquidity cap", coin.LiquidityCap.String())
	})

	t.Run("should fail if not admin", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		coinAddress := sample.EthAddress().String()

		foreignCoin := sample.ForeignCoins(t, coinAddress)
		foreignCoin.LiquidityCap = math.Uint{}
		k.SetForeignCoins(ctx, foreignCoin)
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateMRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(42),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.UpdateMRC20LiquidityCap(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if mrc20 does not exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		coinAddress := sample.EthAddress().String()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateMRC20LiquidityCap(
			admin,
			coinAddress,
			math.NewUint(42),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateMRC20LiquidityCap(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
