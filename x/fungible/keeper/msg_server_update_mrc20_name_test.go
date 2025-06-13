package keeper_test

import (
	"github.com/RWAs-labs/muse/pkg/chains"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/keeper"
	"github.com/RWAs-labs/muse/x/fungible/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgServer_UpdateMRC20Name(t *testing.T) {
	t.Run("should fail if not admin", func(t *testing.T) {
		// arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		authorityMock.On("GetAdditionalChainList", mock.Anything).Return([]chains.Chain{})

		admin := sample.AccAddress()
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20Address := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "MRC20", "MRC20")

		msg := types.NewMsgUpdateMRC20Name(
			admin,
			mrc20Address.Hex(),
			"foo",
			"bar",
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)

		// act
		_, err := msgServer.UpdateMRC20Name(ctx, msg)

		// assert
		require.Error(t, err)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if invalid mrc20 address", func(t *testing.T) {
		// arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		authorityMock.On("GetAdditionalChainList", mock.Anything).Return([]chains.Chain{})

		admin := sample.AccAddress()
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "MRC20", "MRC20")

		msg := types.NewMsgUpdateMRC20Name(
			admin,
			"invalid",
			"foo",
			"bar",
		)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// act
		_, err := msgServer.UpdateMRC20Name(ctx, msg)

		// assert
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("should fail if non existent mrc20", func(t *testing.T) {
		// arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		authorityMock.On("GetAdditionalChainList", mock.Anything).Return([]chains.Chain{})

		admin := sample.AccAddress()
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20Address := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "MRC20", "MRC20")
		k.RemoveForeignCoins(ctx, mrc20Address.Hex())

		msg := types.NewMsgUpdateMRC20Name(
			admin,
			mrc20Address.Hex(),
			"foo",
			"bar",
		)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// act
		_, err := msgServer.UpdateMRC20Name(ctx, msg)

		// assert
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("can update name and symbol", func(t *testing.T) {
		// arrange
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		authorityMock.On("GetAdditionalChainList", mock.Anything).Return([]chains.Chain{})

		admin := sample.AccAddress()
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20Address := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "MRC20", "MRC20")

		msg := types.NewMsgUpdateMRC20Name(
			admin,
			mrc20Address.Hex(),
			"foo",
			"bar",
		)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// act
		_, err := msgServer.UpdateMRC20Name(ctx, msg)

		// assert
		require.NoError(t, err)

		// check the name and symbol
		name, err := k.MRC20Name(ctx, mrc20Address)
		require.NoError(t, err)
		require.Equal(t, "foo", name)

		symbol, err := k.MRC20Symbol(ctx, mrc20Address)
		require.NoError(t, err)
		require.Equal(t, "bar", symbol)

		// check object
		fc, found := k.GetForeignCoins(ctx, mrc20Address.Hex())
		require.True(t, found)
		require.Equal(t, "foo", fc.Name)
		require.Equal(t, "bar", fc.Symbol)

		// can update name only
		// arrange
		msg = types.NewMsgUpdateMRC20Name(
			admin,
			mrc20Address.Hex(),
			"foo2",
			"",
		)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// act
		_, err = msgServer.UpdateMRC20Name(ctx, msg)

		// assert
		require.NoError(t, err)

		name, err = k.MRC20Name(ctx, mrc20Address)
		require.NoError(t, err)
		require.Equal(t, "foo2", name)

		symbol, err = k.MRC20Symbol(ctx, mrc20Address)
		require.NoError(t, err)
		require.Equal(t, "bar", symbol)

		// check object
		fc, found = k.GetForeignCoins(ctx, mrc20Address.Hex())
		require.True(t, found)
		require.Equal(t, "foo2", fc.Name)
		require.Equal(t, "bar", fc.Symbol)

		// can update symbol only
		// arrange
		msg = types.NewMsgUpdateMRC20Name(
			admin,
			mrc20Address.Hex(),
			"",
			"bar2",
		)

		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// act
		_, err = msgServer.UpdateMRC20Name(ctx, msg)

		// assert
		require.NoError(t, err)

		name, err = k.MRC20Name(ctx, mrc20Address)
		require.NoError(t, err)
		require.Equal(t, "foo2", name)

		symbol, err = k.MRC20Symbol(ctx, mrc20Address)
		require.NoError(t, err)
		require.Equal(t, "bar2", symbol)

		// check object
		fc, found = k.GetForeignCoins(ctx, mrc20Address.Hex())
		require.True(t, found)
		require.Equal(t, "foo2", fc.Name)
		require.Equal(t, "bar2", fc.Symbol)
	})
}
