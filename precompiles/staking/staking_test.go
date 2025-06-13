package staking

import (
	"encoding/json"
	"testing"

	"github.com/RWAs-labs/muse/e2e/contracts/erc1967proxy"

	"math/big"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/holiman/uint256"

	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	ethermint "github.com/RWAs-labs/ethermint/types"
	evmkeeper "github.com/RWAs-labs/ethermint/x/evm/keeper"
	"github.com/RWAs-labs/ethermint/x/evm/statedb"
	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/ptr"
	"github.com/RWAs-labs/muse/precompiles/prototype"
	"github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	fungiblekeeper "github.com/RWAs-labs/muse/x/fungible/keeper"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func Test_IStakingContract(t *testing.T) {
	s := newTestSuite(t)
	gasConfig := storetypes.TransientGasConfig()

	t.Run("should check methods are present in ABI", func(t *testing.T) {
		require.NotNil(t, s.stkContractABI.Methods[StakeMethodName], "stake method should be present in the ABI")
		require.NotNil(t, s.stkContractABI.Methods[UnstakeMethodName], "unstake method should be present in the ABI")
		require.NotNil(
			t,
			s.stkContractABI.Methods[MoveStakeMethodName],
			"moveStake method should be present in the ABI",
		)

		require.NotNil(
			t,
			s.stkContractABI.Methods[GetAllValidatorsMethodName],
			"getAllValidators method should be present in the ABI",
		)
		require.NotNil(
			t,
			s.stkContractABI.Methods[GetSharesMethodName],
			"getShares method should be present in the ABI",
		)
	})

	t.Run("should check gas requirements for methods", func(t *testing.T) {
		var method [4]byte

		t.Run("stake", func(t *testing.T) {
			// ACT
			stake := s.stkContract.RequiredGas(s.stkContractABI.Methods[StakeMethodName].ID)
			// ASSERT
			copy(method[:], s.stkContractABI.Methods[StakeMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				stake,
				"stake method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				stake,
			)
		})

		t.Run("unstake", func(t *testing.T) {
			// ACT
			unstake := s.stkContract.RequiredGas(s.stkContractABI.Methods[UnstakeMethodName].ID)
			// ASSERT
			copy(method[:], s.stkContractABI.Methods[UnstakeMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				unstake,
				"unstake method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				unstake,
			)
		})

		t.Run("moveStake", func(t *testing.T) {
			// ACT
			moveStake := s.stkContract.RequiredGas(s.stkContractABI.Methods[MoveStakeMethodName].ID)
			// ASSERT
			copy(method[:], s.stkContractABI.Methods[MoveStakeMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				moveStake,
				"moveStake method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				moveStake,
			)
		})

		t.Run("getAllValidators", func(t *testing.T) {
			// ACT
			getAllValidators := s.stkContract.RequiredGas(s.stkContractABI.Methods[GetAllValidatorsMethodName].ID)
			// ASSERT
			copy(method[:], s.stkContractABI.Methods[GetAllValidatorsMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.ReadCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				getAllValidators,
				"getAllValidators method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				getAllValidators,
			)
		})

		t.Run("getShares", func(t *testing.T) {
			// ACT
			getShares := s.stkContract.RequiredGas(s.stkContractABI.Methods[GetSharesMethodName].ID)
			// ASSERT
			copy(method[:], s.stkContractABI.Methods[GetSharesMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.ReadCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				getShares,
				"getShares method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				getShares,
			)
		})

		t.Run("invalid method", func(t *testing.T) {
			// ARRANGE
			invalidMethodBytes := []byte("invalidMethod")
			// ACT
			gasInvalidMethod := s.stkContract.RequiredGas(invalidMethodBytes)
			// ASSERT
			require.Equal(
				t,
				uint64(0),
				gasInvalidMethod,
				"invalid method should require %d gas, got %d",
				uint64(0),
				gasInvalidMethod,
			)
		})
	})
}

func Test_InvalidMethod(t *testing.T) {
	s := newTestSuite(t)

	_, doNotExist := s.stkContractABI.Methods["invalidMethod"]
	require.False(t, doNotExist, "invalidMethod should not be present in the ABI")
}

func Test_InvalidABI(t *testing.T) {
	IStakingMetaData.ABI = "invalid json"
	defer func() {
		if r := recover(); r != nil {
			require.IsType(t, &json.SyntaxError{}, r, "expected error type: json.SyntaxError, got: %T", r)
		}
	}()

	initABI()
}

func Test_RunInvalidMethod(t *testing.T) {
	// ARRANGE
	s := newTestSuite(t)

	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	gasConfig := storetypes.TransientGasConfig()

	prototype := prototype.NewIPrototypeContract(s.fungibleKeeper, appCodec, gasConfig)

	prototypeAbi := prototype.Abi()
	methodID := prototypeAbi.Methods["bech32ToHexAddr"]
	args := []interface{}{"123"}
	s.mockVMContract.Input = packInputArgs(t, methodID, args...)

	// ACT
	_, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)

	// ASSERT
	require.Error(t, err)
}

func setup(t *testing.T) (sdk.Context, *Contract, abi.ABI, keeper.SDKKeepers, *vm.EVM, *vm.Contract) {
	// Initialize state.
	// Get sdk keepers initialized with this state and the context.
	cdc := keeper.NewCodec()
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	keys, memKeys, tkeys, allKeys := keeper.StoreKeys()

	sdkKeepers := keeper.NewSDKKeepersWithKeys(cdc, keys, memKeys, tkeys, allKeys)

	for _, key := range keys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	}
	for _, key := range tkeys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeTransient, nil)
	}
	for _, key := range memKeys {
		stateStore.MountStoreWithDB(key, storetypes.StoreTypeMemory, nil)
	}

	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := keeper.NewContext(stateStore)

	// Intiliaze codecs and gas config.
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	gasConfig := storetypes.TransientGasConfig()

	stakingGenesisState := stakingtypes.DefaultGenesisState()
	stakingGenesisState.Params.BondDenom = config.BaseDenom
	sdkKeepers.StakingKeeper.InitGenesis(ctx, stakingGenesisState)

	// Get the fungible keeper.
	fungibleKeeper, _, _, _ := keeper.FungibleKeeper(t)

	accAddress := sdk.AccAddress(ContractAddress.Bytes())
	acc := fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ctx, accAddress)
	fungibleKeeper.GetAuthKeeper().SetAccount(ctx, acc)

	// Initialize staking contract.
	contract := NewIStakingContract(
		ctx,
		&sdkKeepers.StakingKeeper,
		*fungibleKeeper,
		sdkKeepers.BankKeeper,
		sdkKeepers.DistributionKeeper,
		appCodec,
		gasConfig,
	)
	require.NotNil(t, contract, "NewIStakingContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	address := contract.Address()
	require.NotNil(t, address, "contract address should not be nil")

	mockEVM := vm.NewEVM(
		vm.BlockContext{},
		vm.TxContext{},
		statedb.New(ctx, sdkKeepers.EvmKeeper, statedb.TxConfig{}),
		&params.ChainConfig{},
		vm.Config{},
	)

	mockVMContract := vm.NewContract(
		contractRef{address: common.Address{}},
		contractRef{address: ContractAddress},
		uint256.NewInt(0),
		0,
	)

	return ctx, contract, abi, sdkKeepers, mockEVM, mockVMContract
}

/*
	Complete Test Suite
	TODO: Migrate all staking tests to this suite.
*/

type testSuite struct {
	ctx            sdk.Context
	stkContract    *Contract
	stkContractABI *abi.ABI
	fungibleKeeper *fungiblekeeper.Keeper
	sdkKeepers     keeper.SDKKeepers
	mockEVM        *vm.EVM
	mockVMContract *vm.Contract
	defaultCaller  common.Address
	defaultLocker  common.Address
	mrc20Address   common.Address
	mrc20ABI       *abi.ABI
}

func newTestSuite(t *testing.T) testSuite {
	// Initialize basic parameters to mock the chain.
	fungibleKeeper, ctx, sdkKeepers, _ := keeper.FungibleKeeper(t)
	chainID := getValidChainID(t)

	// Make sure the account store is initialized.
	// This is completely needed for accounts to be created in the state.
	fungibleKeeper.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

	// Deploy system contracts in order to deploy a MRC20 token.
	deploySystemContracts(t, ctx, fungibleKeeper, *sdkKeepers.EvmKeeper)
	mrc20Address := setupGasCoin(t, ctx, fungibleKeeper, sdkKeepers.EvmKeeper, chainID, "MRC20", "MRC20")

	// Keepers and chain configuration.
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	gasConfig := storetypes.TransientGasConfig()

	// Create the staking contract.
	contract := NewIStakingContract(
		ctx,
		&sdkKeepers.StakingKeeper,
		*fungibleKeeper,
		sdkKeepers.BankKeeper,
		sdkKeepers.DistributionKeeper,
		appCodec,
		gasConfig,
	)
	require.NotNil(t, contract, "NewIStakingContract() should not return a nil contract")

	accAddress := sdk.AccAddress(ContractAddress.Bytes())
	acc := fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ctx, accAddress)
	fungibleKeeper.GetAuthKeeper().SetAccount(ctx, acc)

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	address := contract.Address()
	require.NotNil(t, address, "contract address should not be nil")

	mockEVM := vm.NewEVM(
		vm.BlockContext{},
		vm.TxContext{},
		statedb.New(ctx, sdkKeepers.EvmKeeper, statedb.TxConfig{}),
		&params.ChainConfig{},
		vm.Config{},
	)

	mockVMContract := vm.NewContract(
		contractRef{address: common.Address{}},
		contractRef{address: ContractAddress},
		uint256.NewInt(0),
		0,
	)

	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	require.NoError(t, err)

	// Default locker is the bank address.
	locker := common.HexToAddress("0x0000000000000000000000000000000000000067")

	// Set default caller.
	caller := fungibletypes.ModuleAddressEVM
	mockVMContract.CallerAddress = caller
	mockEVM.Origin = caller

	return testSuite{
		ctx,
		contract,
		&abi,
		fungibleKeeper,
		sdkKeepers,
		mockEVM,
		mockVMContract,
		caller,
		locker,
		mrc20Address,
		mrc20ABI,
	}
}

func packInputArgs(t *testing.T, methodID abi.Method, args ...interface{}) []byte {
	input, err := methodID.Inputs.Pack(args...)
	require.NoError(t, err)
	return append(methodID.ID, input...)
}

func allowStaking(t *testing.T, ts testSuite, amount *big.Int) {
	resAllowance, err := callEVM(
		t,
		ts.ctx,
		ts.fungibleKeeper,
		ts.mrc20ABI,
		fungibletypes.ModuleAddressEVM,
		ts.mrc20Address,
		"approve",
		[]interface{}{ts.stkContract.Address(), amount},
	)
	require.NoError(t, err, "error allowing staking to spend MRC20 tokens")

	allowed, ok := resAllowance[0].(bool)
	require.True(t, ok)
	require.True(t, allowed)
}

func stakeThroughCosmosAPI(
	t *testing.T,
	ctx sdk.Context,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	validator stakingtypes.Validator,
	staker sdk.AccAddress,
	amount math.Int,
) {
	// Coins to stake with default cosmos denom.
	coins := sdk.NewCoins(sdk.NewCoin("stake", amount))

	err := bankKeeper.MintCoins(ctx, fungibletypes.ModuleName, coins)
	require.NoError(t, err)

	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	require.NoError(t, err)

	shares, err := stakingKeeper.Delegate(
		ctx,
		staker,
		coins.AmountOf(coins.Denoms()[0]),
		validator.Status,
		validator,
		true,
	)
	require.NoError(t, err)
	require.Equal(t, amount.Uint64(), shares.TruncateInt().Uint64())
}

func distributeMRC20(
	t *testing.T,
	s testSuite,
	amount *big.Int,
) {
	distributeMethod := s.stkContractABI.Methods[DistributeMethodName]

	_, err := s.fungibleKeeper.DepositMRC20(s.ctx, s.mrc20Address, s.defaultCaller, amount)
	require.NoError(t, err)
	allowStaking(t, s, amount)

	// Setup method input.
	s.mockVMContract.Input = packInputArgs(
		t,
		distributeMethod,
		[]interface{}{s.mrc20Address, amount}...,
	)

	// Call distribute method.
	success, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
	require.NoError(t, err)

	res, err := distributeMethod.Outputs.Unpack(success)
	require.NoError(t, err)

	ok := res[0].(bool)
	require.True(t, ok)
}

func callEVM(
	t *testing.T,
	ctx sdk.Context,
	fungibleKeeper *fungiblekeeper.Keeper,
	abi *abi.ABI,
	from common.Address,
	dst common.Address,
	method string,
	args []interface{},
) ([]interface{}, error) {
	res, err := fungibleKeeper.CallEVM(
		ctx,           // ctx
		*abi,          // abi
		from,          // from
		dst,           // to
		big.NewInt(0), // value
		nil,           // gasLimit
		true,          // commit
		true,          // noEthereumTxEvent
		method,        // method
		args...,       // args
	)
	require.NoError(t, err, "CallEVM error")
	require.Equal(t, "", res.VmError, "res.VmError should be empty")

	ret, err := abi.Methods[method].Outputs.Unpack(res.Ret)
	require.NoError(t, err, "Unpack error")

	return ret, nil
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
		ptr.Ptr(math.NewUintFromString("100000000000000000000000000")),
	)
	require.NoError(t, err)
	assertContractDeployment(t, *evmk, ctx, addr)
	return addr
}

// get a valid chain id independently of the build flag
func getValidChainID(t *testing.T) int64 {
	list := chains.DefaultChainsList()
	require.NotEmpty(t, list)
	require.NotNil(t, list[0])
	return list[0].ChainId
}

// require that a contract has been deployed by checking stored code is non-empty.
func assertContractDeployment(t *testing.T, k evmkeeper.Keeper, ctx sdk.Context, contractAddress common.Address) {
	acc := k.GetAccount(ctx, contractAddress)
	require.NotNil(t, acc)
	code := k.GetCode(ctx, common.BytesToHash(acc.CodeHash))
	require.NotEmpty(t, code)
}

// deploySystemContracts deploys the system contracts and returns their addresses.
func deploySystemContracts(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk evmkeeper.Keeper,
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
	contract := deployGatewayContract(t, ctx, k, &evmk, wmuse, sample.EthAddress())
	require.NotEmpty(t, contract)

	return
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
	assertContractDeployment(t, *evmk, ctx, implAddr)

	// Deploy the proxy contract
	gatewayABI, err := gatewaymevm.GatewayMEVMMetaData.GetAbi()
	require.NoError(t, err)

	// Encode the initializer data
	initializerData, err := gatewayABI.Pack("initialize", wmuse, admin)
	require.NoError(t, err)

	gatewayContract, err := k.DeployContract(ctx, erc1967proxy.ERC1967ProxyMetaData, implAddr, initializerData)
	require.NoError(t, err)
	require.NotEmpty(t, gatewayContract)
	assertContractDeployment(t, *evmk, ctx, gatewayContract)

	// store the gateway in the system contract object
	sys, found := k.GetSystemContract(ctx)
	if !found {
		sys = fungibletypes.SystemContract{}
	}
	sys.Gateway = gatewayContract.Hex()
	k.SetSystemContract(ctx, sys)

	return gatewayContract
}

type contractRef struct {
	address common.Address
}

func (c contractRef) Address() common.Address {
	return c.address
}
