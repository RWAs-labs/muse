// This file contains helper functions for testing the crosschain module
package keeper_test

import (
	"math/big"
	"testing"

	"github.com/RWAs-labs/muse/e2e/contracts/erc1967proxy"

	sdkmath "cosmossdk.io/math"

	evmkeeper "github.com/RWAs-labs/ethermint/x/evm/keeper"
	"github.com/RWAs-labs/muse/pkg/contracts/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/RWAs-labs/muse/pkg/ptr"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/pkg/chains"
	testkeeper "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	fungiblekeeper "github.com/RWAs-labs/muse/x/fungible/keeper"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// getValidEthChainID get a valid eth chain id
func getValidEthChainID() int64 {
	return getValidEthChain().ChainId
}

// getValidEthChain() get a valid eth chain
func getValidEthChain() chains.Chain {
	goerli := chains.GoerliLocalnet
	return goerli
}

func getValidBTCChain() chains.Chain {
	btcRegNet := chains.BitcoinRegtest
	return btcRegNet
}

func getValidBtcChainID() int64 {
	return getValidBTCChain().ChainId
}

func getValidSolanaChainID() int64 {
	return chains.SolanaLocalnet.ChainId
}

// getValidEthChainIDWithIndex get a valid eth chain id with index
func getValidEthChainIDWithIndex(t *testing.T, index int) int64 {
	switch index {
	case 0:
		return chains.GoerliLocalnet.ChainId
	case 1:
		return chains.Goerli.ChainId
	default:
		require.Fail(t, "invalid index")
	}
	return 0
}

// require that a contract has been deployed by checking stored code is non-empty.
func assertContractDeployment(t *testing.T, k *evmkeeper.Keeper, ctx sdk.Context, contractAddress common.Address) {
	acc := k.GetAccount(ctx, contractAddress)
	require.NotNil(t, acc)

	code := k.GetCode(ctx, common.BytesToHash(acc.CodeHash))
	require.NotEmpty(t, code)
}

// deploy upgradable gateway contract and return its address
func deployGatewayContract(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk *evmkeeper.Keeper,
	wmuse, admin common.Address,
) common.Address {
	// Deploy the gateway contract
	implAddr, err := k.DeployContract(ctx, gatewaymevm.GatewayMEVMMetaData)
	require.NoError(t, err)
	require.NotEmpty(t, implAddr)
	assertContractDeployment(t, evmk, ctx, implAddr)

	// Deploy the proxy contract
	gatewayABI, err := gatewaymevm.GatewayMEVMMetaData.GetAbi()
	require.NoError(t, err)

	// Encode the initializer data
	initializerData, err := gatewayABI.Pack("initialize", wmuse, admin)
	require.NoError(t, err)

	gatewayContract, err := k.DeployContract(ctx, erc1967proxy.ERC1967ProxyMetaData, implAddr, initializerData)
	require.NoError(t, err)
	require.NotEmpty(t, gatewayContract)
	assertContractDeployment(t, evmk, ctx, gatewayContract)

	// store the gateway in the system contract object
	sys, found := k.GetSystemContract(ctx)
	if !found {
		sys = fungibletypes.SystemContract{}
	}
	sys.Gateway = gatewayContract.Hex()
	k.SetSystemContract(ctx, sys)

	return gatewayContract
}

// deploySystemContracts deploys the system contracts and returns their addresses.
func deploySystemContracts(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk *evmkeeper.Keeper,
) (wmuse, uniswapV2Factory, uniswapV2Router, connector, systemContract common.Address) {
	var err error

	wmuse, err = k.DeployWMUSE(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, wmuse)
	assertContractDeployment(t, evmk, ctx, wmuse)

	uniswapV2Factory, err = k.DeployUniswapV2Factory(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, uniswapV2Factory)
	assertContractDeployment(t, evmk, ctx, uniswapV2Factory)

	uniswapV2Router, err = k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wmuse)
	require.NoError(t, err)
	require.NotEmpty(t, uniswapV2Router)
	assertContractDeployment(t, evmk, ctx, uniswapV2Router)

	connector, err = k.DeployConnectorMEVM(ctx, wmuse)
	require.NoError(t, err)
	require.NotEmpty(t, connector)
	assertContractDeployment(t, evmk, ctx, connector)

	systemContract, err = k.DeploySystemContract(ctx, wmuse, uniswapV2Factory, uniswapV2Router)
	require.NoError(t, err)
	require.NotEmpty(t, systemContract)
	assertContractDeployment(t, evmk, ctx, systemContract)

	// deploy the gateway contract
	contract := deployGatewayContract(t, ctx, k, evmk, wmuse, sample.EthAddress())
	require.NotEmpty(t, contract)

	return
}

// setupGasCoin is a helper function to setup the gas coin for testing
func setupGasCoin(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk *evmkeeper.Keeper,
	chainID int64,
	assetName string,
	symbol string,
) (mrc20 common.Address) {
	addr, err := k.SetupChainGasCoinAndPool(
		ctx,
		chainID,
		assetName,
		symbol,
		8,
		nil,
		ptr.Ptr(sdkmath.NewUint(1000)),
	)
	require.NoError(t, err)
	assertContractDeployment(t, evmk, ctx, addr)
	return addr
}

// deployMRC20 deploys a MRC20 contract and returns its address
func deployMRC20(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk *evmkeeper.Keeper,
	chainID int64,
	assetName string,
	assetAddress string,
	symbol string,
) (mrc20 common.Address) {
	addr, err := k.DeployMRC20Contract(
		ctx,
		assetName,
		symbol,
		8,
		chainID,
		0,
		assetAddress,
		big.NewInt(21_000),
		ptr.Ptr(sdkmath.NewUint(1000)),
	)
	require.NoError(t, err)
	assertContractDeployment(t, evmk, ctx, addr)
	return addr
}

// setupMRC20Pool setup a Uniswap pool with liquidity for the pair muse/asset
func setupMRC20Pool(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	mrc20Addr common.Address,
) {
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	require.NoError(t, err)
	routerABI, err := uniswapv2router02.UniswapV2Router02MetaData.GetAbi()
	require.NoError(t, err)

	// enough for the small numbers used in test
	liquidityAmount := big.NewInt(1e17)

	// mint some mrc20 and muse
	_, err = k.DepositMRC20(ctx, mrc20Addr, types.ModuleAddressEVM, liquidityAmount)
	require.NoError(t, err)
	err = bankKeeper.MintCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewIntFromBigInt(liquidityAmount))),
	)
	require.NoError(t, err)

	// approve the router to spend the mrc20
	err = k.CallMRC20Approve(
		ctx,
		types.ModuleAddressEVM,
		mrc20Addr,
		routerAddress,
		liquidityAmount,
		false,
	)
	require.NoError(t, err)

	// add the liquidity
	//function addLiquidityETH(
	//	address token,
	//	uint amountTokenDesired,
	//	uint amountTokenMin,
	//	uint amountETHMin,
	//	address to,
	//	uint deadline
	//)
	_, err = k.CallEVM(
		ctx,
		*routerABI,
		types.ModuleAddressEVM,
		routerAddress,
		liquidityAmount,
		big.NewInt(5_000_000),
		true,
		false,
		"addLiquidityETH",
		mrc20Addr,
		liquidityAmount,
		fungiblekeeper.BigIntZero,
		fungiblekeeper.BigIntZero,
		types.ModuleAddressEVM,
		liquidityAmount,
	)
	require.NoError(t, err)
}

// setSupportedChain sets the supported chains for the observer module
func setSupportedChain(ctx sdk.Context, zk testkeeper.MuseKeepers, chainIDs ...int64) {
	chainParamsList := make([]*observertypes.ChainParams, len(chainIDs))
	for i, chainID := range chainIDs {
		chainParams := sample.ChainParams(chainID)
		chainParams.IsSupported = true
		chainParamsList[i] = chainParams
	}
	zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{
		ChainParams: chainParamsList,
	})
}
