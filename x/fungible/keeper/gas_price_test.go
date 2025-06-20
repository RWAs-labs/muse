package keeper_test

import (
	"math/big"
	"testing"

	"github.com/RWAs-labs/protocol-contracts/pkg/systemcontract.sol"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/keeper"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestKeeper_SetGasPrice(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, _, _, _, system := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

	queryGasPrice := func(chainID *big.Int) *big.Int {
		abi, err := systemcontract.SystemContractMetaData.GetAbi()
		require.NoError(t, err)
		res, err := k.CallEVM(
			ctx,
			*abi,
			types.ModuleAddressEVM,
			system,
			keeper.BigIntZero,
			nil,
			false,
			false,
			"gasPriceByChainId",
			chainID,
		)
		require.NoError(t, err)
		unpacked, err := abi.Unpack("gasPriceByChainId", res.Ret)
		require.NoError(t, err)
		gasPrice, ok := unpacked[0].(*big.Int)
		require.True(t, ok)
		return gasPrice
	}

	_, err := k.SetGasPrice(ctx, big.NewInt(1), big.NewInt(42))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(42), queryGasPrice(big.NewInt(1)))
}

func TestKeeper_SetGasPriceContractNotFound(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.SetGasPrice(ctx, big.NewInt(1), big.NewInt(42))
	require.ErrorIs(t, err, types.ErrContractNotFound)
}

func TestKeeper_SetNilGasPrice(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, err := k.SetGasPrice(ctx, big.NewInt(1), nil)
	require.ErrorIs(t, err, types.ErrNilGasPrice)
}

func TestKeeper_SetGasPriceContractIs0x0(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	k.SetSystemContract(ctx, types.SystemContract{})

	_, err := k.SetGasPrice(ctx, big.NewInt(1), big.NewInt(42))
	require.ErrorIs(t, err, types.ErrContractNotFound)
}

func TestKeeper_SetGasPriceReverts(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMFailCallOnce()
	_, err := k.SetGasPrice(ctx, big.NewInt(1), big.NewInt(1))
	require.ErrorIs(t, err, types.ErrContractCall)
}

func TestKeeper_SetGasCoin(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	gas := sample.EthAddress()

	deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	err := k.SetGasCoin(ctx, big.NewInt(1), gas)
	require.NoError(t, err)

	found, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(1))
	require.NoError(t, err)
	require.Equal(t, gas.Hex(), found.Hex())
}

func TestKeeper_SetGasCoinContractNotFound(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	gas := sample.EthAddress()

	err := k.SetGasCoin(ctx, big.NewInt(1), gas)
	require.ErrorIs(t, err, types.ErrContractNotFound)
}

func TestKeeper_SetGasCoinContractIs0x0(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	gas := sample.EthAddress()

	deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	k.SetSystemContract(ctx, types.SystemContract{})

	err := k.SetGasCoin(ctx, big.NewInt(1), gas)
	require.ErrorIs(t, err, types.ErrContractNotFound)
}

func TestKeeper_SetGasCoinReverts(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMFailCallOnce()
	err := k.SetGasCoin(ctx, big.NewInt(1), sample.EthAddress())
	require.ErrorIs(t, err, types.ErrContractCall)
}

func TestKeeper_SetGasMusePool(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	mrc20 := sample.EthAddress()

	_, _, _, _, system := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

	queryMusePool := func(chainID *big.Int) ethcommon.Address {
		abi, err := systemcontract.SystemContractMetaData.GetAbi()
		require.NoError(t, err)
		res, err := k.CallEVM(
			ctx,
			*abi,
			types.ModuleAddressEVM,
			system,
			keeper.BigIntZero,
			nil,
			false,
			false,
			"gasMusePoolByChainId",
			chainID,
		)
		require.NoError(t, err)
		unpacked, err := abi.Unpack("gasMusePoolByChainId", res.Ret)
		require.NoError(t, err)
		pool, ok := unpacked[0].(ethcommon.Address)
		require.True(t, ok)
		return pool
	}

	err := k.SetGasMusePool(ctx, big.NewInt(1), mrc20)
	require.NoError(t, err)
	require.NotEqual(t, ethcommon.Address{}, queryMusePool(big.NewInt(1)))
}

func TestKeeper_SetGasMusePoolContractNotFound(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	mrc20 := sample.EthAddress()

	err := k.SetGasMusePool(ctx, big.NewInt(1), mrc20)
	require.ErrorIs(t, err, types.ErrContractNotFound)
}

func TestKeeper_SetGasMusePoolContractIs0x0(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	mrc20 := sample.EthAddress()

	deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	k.SetSystemContract(ctx, types.SystemContract{})

	err := k.SetGasMusePool(ctx, big.NewInt(1), mrc20)
	require.ErrorIs(t, err, types.ErrContractNotFound)
}

func TestKeeper_SetGasMusePoolReverts(t *testing.T) {
	k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
		UseEVMMock: true,
	})
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)

	mockEVMKeeper.MockEVMFailCallOnce()
	err := k.SetGasMusePool(ctx, big.NewInt(1), sample.EthAddress())
	require.ErrorIs(t, err, types.ErrContractCall)
}
