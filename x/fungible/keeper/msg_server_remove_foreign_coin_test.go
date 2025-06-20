package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/keeper"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestMsgServer_RemoveForeignCoin(t *testing.T) {
	t.Run("can remove a foreign coin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		chainID := getValidChainID(t)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", "foo")

		_, found := k.GetForeignCoins(ctx, mrc20.Hex())
		require.True(t, found)

		msg := types.NewMsgRemoveForeignCoin(admin, mrc20.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.RemoveForeignCoin(ctx, msg)
		require.NoError(t, err)
		_, found = k.GetForeignCoins(ctx, mrc20.Hex())
		require.False(t, found)
	})

	t.Run("should fail if not admin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", "foo")

		msg := types.NewMsgRemoveForeignCoin(admin, mrc20.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.RemoveForeignCoin(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgRemoveForeignCoin(admin, sample.EthAddress().Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.RemoveForeignCoin(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})
}
