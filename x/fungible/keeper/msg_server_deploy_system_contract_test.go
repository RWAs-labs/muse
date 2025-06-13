package keeper_test

import (
	"context"
	"math/big"
	"testing"

	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/keeper"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestMsgServer_DeploySystemContracts(t *testing.T) {
	t.Run("admin can deploy system contracts", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgDeploySystemContracts(admin)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		res, err := msgServer.DeploySystemContracts(ctx, msg)

		require.NoError(t, err)
		require.NotNil(t, res)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.UniswapV2Factory))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.Wmuse))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.UniswapV2Router))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.ConnectorMEVM))
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, ethcommon.HexToAddress(res.SystemContract))
	})

	t.Run("non-admin cannot deploy system contracts", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		nonadmin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgDeploySystemContracts(nonadmin)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.DeploySystemContracts(ctx, msg)

		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if uniswapv2factory contract deployment fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		// mock failed uniswapv2factory deployment
		mockFailedContractDeployment(ctx, t, k)

		msg := types.NewMsgDeploySystemContracts(admin)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.DeploySystemContracts(ctx, msg)
		require.ErrorContains(t, err, "failed to deploy")
	})

	t.Run("should fail if wmuse contract deployment fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		// mock successful uniswapv2factory deployment
		mockSuccessfulContractDeployment(ctx, t, k)
		// mock failed wmuse deployment deployment
		mockFailedContractDeployment(ctx, t, k)

		msg := types.NewMsgDeploySystemContracts(admin)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.DeploySystemContracts(ctx, msg)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to deploy")
	})

	t.Run("should fail if uniswapv2router deployment fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		// mock successful uniswapv2factory and wmuse deployments
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		// mock failed uniswapv2router deployment
		mockFailedContractDeployment(ctx, t, k)

		msg := types.NewMsgDeploySystemContracts(admin)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.DeploySystemContracts(ctx, msg)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to deploy")
	})

	t.Run("should fail if connectormevm deployment fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		// mock successful uniswapv2factory, wmuse and uniswapv2router deployments
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		// mock failed connectormevm deployment
		mockFailedContractDeployment(ctx, t, k)

		msg := types.NewMsgDeploySystemContracts(admin)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		_, err := msgServer.DeploySystemContracts(ctx, msg)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to deploy")
	})

	t.Run("should fail if system contract deployment fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		// mock successful uniswapv2factory, wmuse, uniswapv2router and connectormevm deployments
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		mockSuccessfulContractDeployment(ctx, t, k)
		// mock failed system contract deployment
		mockFailedContractDeployment(ctx, t, k)

		msg := types.NewMsgDeploySystemContracts(admin)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.DeploySystemContracts(ctx, msg)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to deploy")
	})
}

func mockSuccessfulContractDeployment(ctx context.Context, t *testing.T, k *keeper.Keeper) {
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
	mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
	mockEVMKeeper.On(
		"EstimateGas",
		mock.Anything,
		mock.Anything,
	).Return(&evmtypes.EstimateGasResponse{Gas: 5}, nil)
	mockEVMKeeper.MockEVMSuccessCallOnce()
}

func mockFailedContractDeployment(ctx context.Context, t *testing.T, k *keeper.Keeper) {
	mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
	mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
	mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))
	mockEVMKeeper.On(
		"EstimateGas",
		mock.Anything,
		mock.Anything,
	).Return(&evmtypes.EstimateGasResponse{Gas: 5}, nil)
	mockEVMKeeper.MockEVMFailCallOnce()
}
