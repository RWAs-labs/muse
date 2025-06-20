package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestKeeper_GetSystemContract(t *testing.T) {
	t.Run("should get and remove system contract", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.SetSystemContract(ctx, types.SystemContract{SystemContract: "test"})
		val, found := k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, types.SystemContract{SystemContract: "test"}, val)

		// can remove contract
		k.RemoveSystemContract(ctx)
		_, found = k.GetSystemContract(ctx)
		require.False(t, found)
	})
}

func TestKeeper_GetSystemContractAddress(t *testing.T) {
	t.Run("should fail to get system contract address if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, err := k.GetSystemContractAddress(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should get system contract address if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, _, _, _, systemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		found, err := k.GetSystemContractAddress(ctx)
		require.NoError(t, err)
		require.Equal(t, systemContract, found)
	})
}

func TestKeeper_GetWMuseContractAddress(t *testing.T) {
	t.Run("should fail to get wmuse contract address if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, err := k.GetWMuseContractAddress(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should get wmuse contract address if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		wmuse, _, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		found, err := k.GetWMuseContractAddress(ctx)
		require.NoError(t, err)
		require.Equal(t, wmuse, found)
	})

	t.Run("should fail if wmuse not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            false,
			DeployUniswapV2Router:  true,
			DeployUniswapV2Factory: true,
		})

		_, err := k.GetWMuseContractAddress(ctx)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()

		_, err := k.GetWMuseContractAddress(ctx)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if abi unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMSuccessCallOnce()

		_, err := k.GetWMuseContractAddress(ctx)
		require.ErrorIs(t, err, types.ErrABIUnpack)
	})
}

func TestKeeper_GetUniswapV2FactoryAddress(t *testing.T) {
	t.Run(
		"should fail to get uniswapfactory contract address if system contracts are not deployed",
		func(t *testing.T) {
			k, ctx, _, _ := keepertest.FungibleKeeper(t)
			k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

			_, err := k.GetUniswapV2FactoryAddress(ctx)
			require.Error(t, err)
			require.ErrorIs(t, err, types.ErrStateVariableNotFound)
		},
	)

	t.Run("should get uniswapfactory contract address if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, factory, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		found, err := k.GetUniswapV2FactoryAddress(ctx)
		require.NoError(t, err)
		require.Equal(t, factory, found)
	})

	t.Run("should fail in factory not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            true,
			DeployUniswapV2Router:  true,
			DeployUniswapV2Factory: false,
		})

		_, err := k.GetUniswapV2FactoryAddress(ctx)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()

		_, err := k.GetUniswapV2FactoryAddress(ctx)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if abi unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMSuccessCallOnce()

		_, err := k.GetUniswapV2FactoryAddress(ctx)
		require.ErrorIs(t, err, types.ErrABIUnpack)
	})
}

func TestKeeper_GetUniswapV2Router02Address(t *testing.T) {
	t.Run("should fail to get uniswaprouter contract address if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, err := k.GetUniswapV2Router02Address(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should get uniswaprouter contract address if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, _, router, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		found, err := k.GetUniswapV2Router02Address(ctx)
		require.NoError(t, err)
		require.Equal(t, router, found)
	})

	t.Run("should fail in router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            true,
			DeployUniswapV2Router:  false,
			DeployUniswapV2Factory: true,
		})

		_, err := k.GetUniswapV2Router02Address(ctx)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()

		_, err := k.GetUniswapV2Router02Address(ctx)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if abi unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMSuccessCallOnce()

		_, err := k.GetUniswapV2Router02Address(ctx)
		require.ErrorIs(t, err, types.ErrABIUnpack)
	})
}

func TestKeeper_CallWMuseDeposit(t *testing.T) {
	t.Run("should fail to deposit if system contracts are not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// mint tokens
		addr := sample.Bech32AccAddress()
		ethAddr := common.BytesToAddress(addr.Bytes())
		coins := sample.Coins()
		err := sdkk.BankKeeper.MintCoins(ctx, types.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkk.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
		require.NoError(t, err)

		// fail if no system contract
		err = k.CallWMuseDeposit(ctx, ethAddr, big.NewInt(42))
		require.Error(t, err)
	})

	t.Run("should deposit if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// mint tokens
		addr := sample.Bech32AccAddress()
		ethAddr := common.BytesToAddress(addr.Bytes())
		coins := sample.Coins()
		err := sdkk.BankKeeper.MintCoins(ctx, types.ModuleName, sample.Coins())
		require.NoError(t, err)
		err = sdkk.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
		require.NoError(t, err)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// deposit
		err = k.CallWMuseDeposit(ctx, ethAddr, big.NewInt(42))
		require.NoError(t, err)

		balance, err := k.QueryWMuseBalanceOf(ctx, ethAddr)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})
}

func TestKeeper_QueryWMuseBalanceOf(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.QueryWMuseBalanceOf(ctx, sample.EthAddress())
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})
}

func TestKeeper_QuerySystemContractGasCoinMRC20(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		_, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should query if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		_, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		found, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID))
		require.NoError(t, err)
		require.Equal(t, mrc20, found)
	})

	t.Run("should fail if gas coin not setup", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		_, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID))
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		_, err = k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID))
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()

		_, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(1))
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if abi unpack fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMSuccessCallOnce()

		_, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(1))
		require.ErrorIs(t, err, types.ErrABIUnpack)
	})
}

func TestKeeper_CallUniswapV2RouterSwapExactETHForToken(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should swap if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// deploy system contracts and swap exact eth for 1 token
		tokenAmount := big.NewInt(1)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetMuseAmountsIn(ctx, tokenAmount, mrc20)
		require.NoError(t, err)
		err = sdkk.BankKeeper.MintCoins(
			ctx,
			types.ModuleName,
			sdk.NewCoins(sdk.NewCoin("amuse", sdkmath.NewIntFromBigInt(amountToSwap))),
		)
		require.NoError(t, err)

		amounts, err := k.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			mrc20,
			true,
		)
		require.NoError(t, err)

		require.Equal(t, 2, len(amounts))
		require.Equal(t, tokenAmount, amounts[1])
	})

	t.Run("should fail if missing muse balance", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// deploy system contracts and swap 1 token fails because of missing wrapped balance
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetMuseAmountsIn(ctx, big.NewInt(1), mrc20)
		require.NoError(t, err)

		_, err = k.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			mrc20,
			true,
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wmuse not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})

		_, err := k.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})

		_, err := k.CallUniswapV2RouterSwapExactETHForToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_CallUniswapV2RouterSwapEthForExactToken(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactETHForToken(
			ctx, types.ModuleAddressEVM, types.ModuleAddressEVM, big.NewInt(1), sample.EthAddress(), true)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should swap if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// deploy system contracts and swap exact 1 token
		tokenAmount := big.NewInt(1)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetMuseAmountsIn(ctx, tokenAmount, mrc20)
		require.NoError(t, err)
		err = sdkk.BankKeeper.MintCoins(
			ctx,
			types.ModuleName,
			sdk.NewCoins(sdk.NewCoin("amuse", sdkmath.NewIntFromBigInt(amountToSwap))),
		)
		require.NoError(t, err)

		amounts, err := k.CallUniswapV2RouterSwapEthForExactToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			tokenAmount,
			mrc20,
		)
		require.NoError(t, err)

		require.Equal(t, 2, len(amounts))
		require.Equal(t, big.NewInt(1), amounts[1])
	})

	t.Run("should fail if missing muse balance", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// deploy system contracts and swap 1 token fails because of missing wrapped balance
		tokenAmount := big.NewInt(1)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetMuseAmountsIn(ctx, tokenAmount, mrc20)
		require.NoError(t, err)

		_, err = k.CallUniswapV2RouterSwapEthForExactToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			tokenAmount,
			mrc20,
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wmuse not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})

		_, err := k.CallUniswapV2RouterSwapEthForExactToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			big.NewInt(1),
			sample.EthAddress(),
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})

		_, err := k.CallUniswapV2RouterSwapEthForExactToken(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			big.NewInt(1),
			sample.EthAddress(),
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_CallUniswapV2RouterSwapExactTokensForETH(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should swap if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.Error(t, err)

		// deploy system contracts and swap exact eth for 1 token
		ethAmount := big.NewInt(1)
		_, _, router, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetMRC4AmountsIn(ctx, ethAmount, mrc20)
		require.NoError(t, err)

		_, err = k.DepositMRC20(ctx, mrc20, types.ModuleAddressEVM, amountToSwap)
		require.NoError(t, err)
		k.CallMRC20Approve(
			ctx,
			types.ModuleAddressEVM,
			mrc20,
			router,
			amountToSwap,
			false,
		)

		amounts, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			mrc20,
			true,
		)
		require.NoError(t, err)

		require.Equal(t, 2, len(amounts))
		require.Equal(t, ethAmount, amounts[0])
	})

	t.Run("should fail if missing tokens balance", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.Error(t, err)

		// deploy system contracts and swap fails because of missing balance
		ethAmount := big.NewInt(1)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		amountToSwap, err := k.QueryUniswapV2RouterGetMRC4AmountsIn(ctx, ethAmount, mrc20)
		require.NoError(t, err)

		_, err = k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx, types.ModuleAddressEVM, types.ModuleAddressEVM, amountToSwap, mrc20, true)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wmuse not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})
		_, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})
		_, err := k.CallUniswapV2RouterSwapExactTokensForETH(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_CallUniswapV2RouterSwapExactTokensForTokens(t *testing.T) {
	t.Run("should fail if system contracts are not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrStateVariableNotFound)
	})

	t.Run("should swap if system contracts are deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.Error(t, err)

		// deploy system contracts and swap exact token for 1 token
		tokenAmount := big.NewInt(1)
		_, _, router, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		inmrc20 := deployMRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", sample.EthAddress().String(), "foo")
		outmrc20 := deployMRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "bar", sample.EthAddress().String(), "bar")
		setupMRC20Pool(t, ctx, k, sdkk.BankKeeper, inmrc20)
		setupMRC20Pool(t, ctx, k, sdkk.BankKeeper, outmrc20)

		amountToSwap, err := k.QueryUniswapV2RouterGetMRC4ToMRC4AmountsIn(ctx, tokenAmount, inmrc20, outmrc20)
		require.NoError(t, err)

		_, err = k.DepositMRC20(ctx, inmrc20, types.ModuleAddressEVM, amountToSwap)
		require.NoError(t, err)
		k.CallMRC20Approve(
			ctx,
			types.ModuleAddressEVM,
			inmrc20,
			router,
			amountToSwap,
			false,
		)

		amounts, err := k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx, types.ModuleAddressEVM, types.ModuleAddressEVM, amountToSwap, inmrc20, outmrc20, true)
		require.NoError(t, err)
		require.Equal(t, 3, len(amounts))
		require.Equal(t, amounts[2], tokenAmount)
	})

	t.Run("should fail if missing tokens balance", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		chainID := getValidChainID(t)

		// deploy system contracts and swap fails because of missing balance
		tokenAmount := big.NewInt(1)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		inmrc20 := deployMRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "foo", sample.EthAddress().String(), "foo")
		outmrc20 := deployMRC20(t, ctx, k, sdkk.EvmKeeper, chainID, "bar", sample.EthAddress().String(), "bar")
		setupMRC20Pool(t, ctx, k, sdkk.BankKeeper, inmrc20)
		setupMRC20Pool(t, ctx, k, sdkk.BankKeeper, outmrc20)

		amountToSwap, err := k.QueryUniswapV2RouterGetMRC4ToMRC4AmountsIn(ctx, tokenAmount, inmrc20, outmrc20)
		require.NoError(t, err)

		_, err = k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			amountToSwap,
			inmrc20,
			outmrc20,
			true,
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wmuse not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.Error(t, err)

		// deploy system contracts except router
		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployUniswapV2Router:  true,
			DeployWMuse:            false,
			DeployUniswapV2Factory: true,
		})

		_, err = k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// fail if no system contract
		_, err := k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.Error(t, err)

		// deploy system contracts except router
		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployUniswapV2Router:  false,
			DeployWMuse:            true,
			DeployUniswapV2Factory: true,
		})

		_, err = k.CallUniswapV2RouterSwapExactTokensForTokens(
			ctx,
			types.ModuleAddressEVM,
			types.ModuleAddressEVM,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
			true,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

}

func TestKeeper_QueryUniswapV2RouterGetMRC4AmountsIn(t *testing.T) {
	t.Run("should fail if no amounts out", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		_, err := k.QueryUniswapV2RouterGetMRC4AmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wmuse not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})

		_, err := k.QueryUniswapV2RouterGetMRC4AmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})

		_, err := k.QueryUniswapV2RouterGetMRC4AmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_QueryUniswapV2RouterGetMuseAmountsIn(t *testing.T) {
	t.Run("should fail if no amounts out", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		_, err := k.QueryUniswapV2RouterGetMuseAmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wmuse not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})

		_, err := k.QueryUniswapV2RouterGetMuseAmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})

		_, err := k.QueryUniswapV2RouterGetMuseAmountsIn(ctx, big.NewInt(1), sample.EthAddress())
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_QueryUniswapV2RouterGetMRC4ToMRC4AmountsIn(t *testing.T) {
	t.Run("should fail if no amounts out", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		_, err := k.QueryUniswapV2RouterGetMRC4ToMRC4AmountsIn(
			ctx,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if wmuse not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            false,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  true,
		})

		_, err := k.QueryUniswapV2RouterGetMRC4ToMRC4AmountsIn(
			ctx,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})

	t.Run("should fail if router not deployed", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsConfigurable(t, ctx, k, sdkk.EvmKeeper, &SystemContractDeployConfig{
			DeployWMuse:            true,
			DeployUniswapV2Factory: true,
			DeployUniswapV2Router:  false,
		})

		_, err := k.QueryUniswapV2RouterGetMRC4ToMRC4AmountsIn(
			ctx,
			big.NewInt(1),
			sample.EthAddress(),
			sample.EthAddress(),
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
	})
}

func TestKeeper_CallMRC20Burn(t *testing.T) {
	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()
		err := k.CallMRC20Burn(ctx, types.ModuleAddressEVM, sample.EthAddress(), big.NewInt(1), false)
		require.ErrorIs(t, err, types.ErrContractCall)
	})
}

func TestKeeper_CallMRC20Approve(t *testing.T) {
	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()
		err := k.CallMRC20Approve(
			ctx,
			types.ModuleAddressEVM,
			sample.EthAddress(),
			types.ModuleAddressEVM,
			big.NewInt(1),
			false,
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})
}

func TestKeeper_CallMRC20Deposit(t *testing.T) {
	t.Run("should fail if evm call fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

		mockEVMKeeper.MockEVMFailCallOnce()
		err := k.CallMRC20Deposit(
			ctx,
			types.ModuleAddressEVM,
			sample.EthAddress(),
			types.ModuleAddressEVM,
			big.NewInt(1),
		)
		require.ErrorIs(t, err, types.ErrContractCall)
	})
}
