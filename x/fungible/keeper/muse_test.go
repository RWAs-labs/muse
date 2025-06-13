package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/mock"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	testkeeper "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/keeper"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestKeeper_MintMuseToEVMAccount(t *testing.T) {
	t.Run("should mint the token in the specified balance", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		acc := sample.Bech32AccAddress()
		bal := sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.IsZero())

		err := k.MintMuseToEVMAccount(ctx, acc, big.NewInt(42))
		require.NoError(t, err)
		bal = sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.Amount.Equal(sdkmath.NewInt(42)))
	})

	t.Run("mint the token to reach max supply", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		acc := sample.Bech32AccAddress()
		bal := sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.IsZero())

		museMaxSupply, ok := sdkmath.NewIntFromString(keeper.MUSEMaxSupplyStr)
		require.True(t, ok)

		supply := sdkk.BankKeeper.GetSupply(ctx, config.BaseDenom).Amount

		newAmount := museMaxSupply.Sub(supply)

		err := k.MintMuseToEVMAccount(ctx, acc, newAmount.BigInt())
		require.NoError(t, err)
		bal = sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.Amount.Equal(newAmount))
	})

	t.Run("can't mint more than max supply", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		acc := sample.Bech32AccAddress()
		bal := sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.IsZero())

		museMaxSupply, ok := sdkmath.NewIntFromString(keeper.MUSEMaxSupplyStr)
		require.True(t, ok)

		supply := sdkk.BankKeeper.GetSupply(ctx, config.BaseDenom).Amount

		newAmount := museMaxSupply.Sub(supply).Add(sdkmath.NewInt(1))

		err := k.MintMuseToEVMAccount(ctx, acc, newAmount.BigInt())
		require.ErrorIs(t, err, types.ErrMaxSupplyReached)
	})

	coins42 := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewInt(42)))

	t.Run("should fail if minting fail", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)

		mockBankKeeper := testkeeper.GetFungibleBankMock(t, k)

		mockBankKeeper.On("GetSupply", ctx, mock.Anything, mock.Anything).
			Return(sdk.NewCoin(config.BaseDenom, sdkmath.NewInt(0))).
			Once()
		mockBankKeeper.On(
			"MintCoins",
			ctx,
			types.ModuleName,
			coins42,
		).Return(errors.New("error"))

		err := k.MintMuseToEVMAccount(ctx, sample.Bech32AccAddress(), big.NewInt(42))
		require.Error(t, err)

		mockBankKeeper.AssertExpectations(t)
	})

	t.Run("should fail if sending coins fail", func(t *testing.T) {
		k, ctx := testkeeper.FungibleKeeperAllMocks(t)
		acc := sample.Bech32AccAddress()

		mockBankKeeper := testkeeper.GetFungibleBankMock(t, k)

		mockBankKeeper.On("GetSupply", ctx, mock.Anything, mock.Anything).
			Return(sdk.NewCoin(config.BaseDenom, sdkmath.NewInt(0))).
			Once()
		mockBankKeeper.On(
			"MintCoins",
			ctx,
			types.ModuleName,
			coins42,
		).Return(nil)

		mockBankKeeper.On(
			"SendCoinsFromModuleToAccount",
			ctx,
			types.ModuleName,
			acc,
			coins42,
		).Return(errors.New("error"))

		err := k.MintMuseToEVMAccount(ctx, acc, big.NewInt(42))
		require.Error(t, err)

		mockBankKeeper.AssertExpectations(t)
	})
}

func TestKeeper_MintMuseToFungibleModule(t *testing.T) {
	t.Run("should mint the token in the specified balance", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		acc := k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName).GetAddress()

		bal := sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.IsZero())

		err := k.MintMuseToEVMAccount(ctx, acc, big.NewInt(42))
		require.NoError(t, err)
		bal = sdkk.BankKeeper.GetBalance(ctx, acc, config.BaseDenom)
		require.True(t, bal.Amount.Equal(sdkmath.NewInt(42)))
	})

	t.Run("can't mint more than max supply", func(t *testing.T) {
		k, ctx, sdkk, _ := testkeeper.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		museMaxSupply, ok := sdkmath.NewIntFromString(keeper.MUSEMaxSupplyStr)
		require.True(t, ok)

		supply := sdkk.BankKeeper.GetSupply(ctx, config.BaseDenom).Amount

		newAmount := museMaxSupply.Sub(supply).Add(sdkmath.NewInt(1))

		err := k.MintMuseToFungibleModule(ctx, newAmount.BigInt())
		require.ErrorIs(t, err, types.ErrMaxSupplyReached)
	})
}
