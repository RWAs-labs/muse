package keeper_test

import (
	"math/big"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/systemcontract.sol"

	"github.com/RWAs-labs/muse/pkg/chains"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/keeper"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestKeeper_UpdateSystemContract(t *testing.T) {
	t.Run("can update the system contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		queryMRC20SystemContract := func(contract common.Address) string {
			abi, err := mrc20.MRC20MetaData.GetAbi()
			require.NoError(t, err)
			res, err := k.CallEVM(
				ctx,
				*abi,
				types.ModuleAddressEVM,
				contract,
				keeper.BigIntZero,
				nil,
				false,
				false,
				"SYSTEM_CONTRACT_ADDRESS",
			)
			require.NoError(t, err)
			unpacked, err := abi.Unpack("SYSTEM_CONTRACT_ADDRESS", res.Ret)
			require.NoError(t, err)
			address, ok := unpacked[0].(common.Address)
			require.True(t, ok)
			return address.Hex()
		}

		chains := chains.DefaultChainsList()
		require.True(t, len(chains) > 1)
		require.NotNil(t, chains[0])
		require.NotNil(t, chains[1])
		chainID1 := chains[0].ChainId
		chainID2 := chains[1].ChainId

		wmuse, factory, router, _, oldSystemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		gas1 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID1, "foo", "foo")
		gas2 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID2, "bar", "bar")
		// this one should be skipped and not impact update
		fc := types.ForeignCoins{
			Mrc20ContractAddress: "0x",
		}
		k.SetForeignCoins(ctx, fc)

		// deploy a new system contracts
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wmuse, factory, router)
		require.NoError(t, err)
		require.NotEqual(t, oldSystemContract, newSystemContract)

		// can update the system contract
		msg := types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateSystemContract(ctx, msg)
		require.NoError(t, err)

		// can retrieve the system contract
		sc, found := k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, newSystemContract.Hex(), sc.SystemContract)

		// check gas updated
		foundGas1, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID1))
		require.NoError(t, err)
		require.Equal(t, gas1, foundGas1)
		foundGas2, err := k.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chainID2))
		require.NoError(t, err)
		require.Equal(t, gas2, foundGas2)

		require.Equal(t, newSystemContract.Hex(), queryMRC20SystemContract(gas1))
		require.Equal(t, newSystemContract.Hex(), queryMRC20SystemContract(gas2))
	})

	t.Run("can update and overwrite the system contract if system contract not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		wmuse, err := k.DeployWMUSE(ctx)
		require.NoError(t, err)

		factory, err := k.DeployUniswapV2Factory(ctx)
		require.NoError(t, err)

		router, err := k.DeployUniswapV2Router02(ctx, factory, wmuse)
		require.NoError(t, err)

		// deploy a new system contracts
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wmuse, factory, router)
		require.NoError(t, err)

		// can update the system contract
		msg := types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateSystemContract(ctx, msg)
		require.NoError(t, err)

		// can retrieve the system contract
		sc, found := k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, newSystemContract.Hex(), sc.SystemContract)

		// deploy a new system contracts
		newSystemContract, err = k.DeployContract(ctx, systemcontract.SystemContractMetaData, wmuse, factory, router)
		require.NoError(t, err)

		// can overwrite the previous system contract
		msg = types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateSystemContract(ctx, msg)
		require.NoError(t, err)

		// can retrieve the system contract
		sc, found = k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, newSystemContract.Hex(), sc.SystemContract)
	})

	t.Run("should not update the system contract if not admin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		// deploy a new system contracts
		wmuse, factory, router, _, oldSystemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wmuse, factory, router)
		require.NoError(t, err)
		require.NotEqual(t, oldSystemContract, newSystemContract)

		// should not update the system contract if not admin
		msg := types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err = msgServer.UpdateSystemContract(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should not update the system contract if invalid address", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		// deploy a new system contracts
		wmuse, factory, router, _, oldSystemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wmuse, factory, router)
		require.NoError(t, err)
		require.NotEqual(t, oldSystemContract, newSystemContract)

		// should not update the system contract if invalid address
		msg := types.NewMsgUpdateSystemContract(admin, "invalid")
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateSystemContract(ctx, msg)
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("should not update if any of 3 evm calls for foreign coin fail", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		chains := chains.DefaultChainsList()
		require.True(t, len(chains) > 1)
		require.NotNil(t, chains[0])
		chainID1 := chains[0].ChainId

		wmuse, factory, router, _, _ := deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)
		// setup mocks and setup gas coin
		var encodedAddress [32]byte
		copy(encodedAddress[12:], router[:])
		uniswapMock := &evmtypes.MsgEthereumTxResponse{
			Ret: encodedAddress[:],
		}
		mockEVMKeeper.MockEVMSuccessCallTimes(4)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(uniswapMock)
		mockEVMKeeper.MockEVMSuccessCallOnce()

		addLiqMockReturn := &evmtypes.MsgEthereumTxResponse{
			Ret: make([]byte, 3*32),
		}
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(addLiqMockReturn)

		setupGasCoin(t, ctx, k, mockEVMKeeper, chainID1, "foo", "foo")

		// deploy a new system contracts
		mockEVMKeeper.MockEVMSuccessCallOnce()
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wmuse, factory, router)
		require.NoError(t, err)

		// fail on first evm call
		mockEVMKeeper.MockEVMFailCallOnce()

		// can't update the system contract
		msg := types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateSystemContract(ctx, msg)
		require.ErrorIs(t, err, types.ErrContractCall)

		// fail on second evm call
		mockEVMKeeper.MockEVMSuccessCallOnce()
		mockEVMKeeper.MockEVMFailCallOnce()

		// can't update the system contract
		msg = types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err = msgServer.UpdateSystemContract(ctx, msg)
		require.ErrorIs(t, err, types.ErrContractCall)

		// fail on third evm call
		mockEVMKeeper.MockEVMSuccessCallTimes(2)
		mockEVMKeeper.MockEVMFailCallOnce()

		// can't update the system contract
		msg = types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err = msgServer.UpdateSystemContract(ctx, msg)
		require.ErrorIs(t, err, types.ErrContractCall)
	})
}
