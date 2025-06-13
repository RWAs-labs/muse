package keeper_test

import (
	"github.com/RWAs-labs/muse/e2e/contracts/testabort"
	"math/big"
	"testing"

	"cosmossdk.io/math"
	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/e2e/contracts/example"
	"github.com/RWAs-labs/muse/e2e/contracts/reverter"
	"github.com/RWAs-labs/muse/e2e/contracts/testdappv2"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	fungiblekeeper "github.com/RWAs-labs/muse/x/fungible/keeper"
	"github.com/RWAs-labs/muse/x/fungible/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// getTestDAppNoMessageIndex queries the no message index of the test dapp v2 contract
func getTestDAppNoMessageIndex(
	t *testing.T,
	ctx sdk.Context,
	k fungiblekeeper.Keeper,
	contract,
	account common.Address,
) string {
	testDAppABI, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(t, err)
	res, err := k.CallEVM(
		ctx,
		*testDAppABI,
		types.ModuleAddressEVM,
		contract,
		fungiblekeeper.BigIntZero,
		nil,
		false,
		false,
		"getNoMessageIndex",
		account,
	)
	require.NoError(t, err)

	unpacked, err := testDAppABI.Unpack("getNoMessageIndex", res.Ret)
	require.NoError(t, err)
	require.Len(t, unpacked, 1)

	index, ok := unpacked[0].(string)
	require.True(t, ok)

	return index
}

// deployTestDAppV2 deploys the test dapp v2 contract and returns its address
func deployTestDAppV2(t *testing.T, ctx sdk.Context, k *fungiblekeeper.Keeper, evmk types.EVMKeeper) common.Address {
	testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData, true, sample.EthAddress())
	require.NoError(t, err)
	require.NotEmpty(t, testDAppV2)
	assertContractDeployment(t, evmk, ctx, testDAppV2)

	return testDAppV2
}

// deployTestAbort deploys the test abort contract and returns its address
func deployTestAbort(t *testing.T, ctx sdk.Context, k *fungiblekeeper.Keeper, evmk types.EVMKeeper) common.Address {
	testAbort, err := k.DeployContract(ctx, testabort.TestAbortMetaData)
	require.NoError(t, err)
	require.NotEmpty(t, testAbort)
	assertContractDeployment(t, evmk, ctx, testAbort)

	return testAbort
}

// assertTestDAppV2MessageAndAmount asserts the message and amount of the test dapp v2 contract
func assertTestDAppV2MessageAndAmount(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	contract common.Address,
	expectedMessage string,
	expectedAmount int64,
) {
	testDAppABI, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(t, err)

	// message
	res, err := k.CallEVM(
		ctx,
		*testDAppABI,
		types.ModuleAddressEVM,
		contract,
		fungiblekeeper.BigIntZero,
		nil,
		false,
		false,
		"getCalledWithMessage",
		expectedMessage,
	)
	require.NoError(t, err)

	unpacked, err := testDAppABI.Unpack("getCalledWithMessage", res.Ret)
	require.NoError(t, err)
	require.Len(t, unpacked, 1)
	found, ok := unpacked[0].(bool)
	require.True(t, ok)
	require.True(t, found)

	// amount
	res, err = k.CallEVM(
		ctx,
		*testDAppABI,
		types.ModuleAddressEVM,
		contract,
		fungiblekeeper.BigIntZero,
		nil,
		false,
		false,
		"getAmountWithMessage",
		expectedMessage,
	)
	require.NoError(t, err)

	unpacked, err = testDAppABI.Unpack("getAmountWithMessage", res.Ret)
	require.NoError(t, err)
	require.Len(t, unpacked, 1)
	amount, ok := unpacked[0].(*big.Int)
	require.True(t, ok)
	require.Equal(t, expectedAmount, amount.Int64())
}

func TestKeeper_MRC20DepositAndCallContract(t *testing.T) {
	t.Run("can deposit gas coin for transfers", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// deposit
		to := sample.EthAddress()
		_, contractCall, err := k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
			crosschaintypes.ProtocolContractVersion_V1,
			false,
		)
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfMRC4(ctx, mrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("can deposit non-gas coin for transfers", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId
		assetAddress := sample.EthAddress().String()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := deployMRC20(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", assetAddress, "foobar")

		// deposit
		to := sample.EthAddress()
		_, contractCall, err := k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_ERC20,
			assetAddress,
			crosschaintypes.ProtocolContractVersion_V1,
			false,
		)
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfMRC4(ctx, mrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("should fail if trying to call a contract with data to a EOC", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId
		assetAddress := sample.EthAddress().String()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		deployMRC20(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", assetAddress, "foobar")

		// deposit
		to := sample.EthAddress()
		_, _, err := k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte("DEADBEEF"),
			coin.CoinType_ERC20,
			assetAddress,
			crosschaintypes.ProtocolContractVersion_V1,
			false,
		)
		require.ErrorIs(t, err, types.ErrCallNonContract)
	})

	t.Run("can deposit coin for transfers with liquidity cap not reached", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// there is an initial total supply minted during gas pool setup
		initialTotalSupply, err := k.TotalSupplyMRC4(ctx, mrc20)
		require.NoError(t, err)

		// set a liquidity cap
		foreignCoin, found := k.GetForeignCoins(ctx, mrc20.String())
		require.True(t, found)
		foreignCoin.LiquidityCap = math.NewUint(initialTotalSupply.Uint64() + 1000)
		k.SetForeignCoins(ctx, foreignCoin)

		// increase total supply
		_, err = k.DepositMRC20(ctx, mrc20, sample.EthAddress(), big.NewInt(500))
		require.NoError(t, err)

		// deposit
		to := sample.EthAddress()
		_, contractCall, err := k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(500),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
			crosschaintypes.ProtocolContractVersion_V1,
			false,
		)
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfMRC4(ctx, mrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(500), balance)
	})

	t.Run("should fail if coin paused", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// pause the coin
		foreignCoin, found := k.GetForeignCoins(ctx, mrc20.String())
		require.True(t, found)
		foreignCoin.Paused = true
		k.SetForeignCoins(ctx, foreignCoin)

		to := sample.EthAddress()
		_, _, err := k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
			crosschaintypes.ProtocolContractVersion_V1,
			false,
		)
		require.ErrorIs(t, err, types.ErrPausedMRC20)
	})

	t.Run("should fail if liquidity cap reached", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// there is an initial total supply minted during gas pool setup
		initialTotalSupply, err := k.TotalSupplyMRC4(ctx, mrc20)
		require.NoError(t, err)

		// set a liquidity cap
		foreignCoin, found := k.GetForeignCoins(ctx, mrc20.String())
		require.True(t, found)
		foreignCoin.LiquidityCap = math.NewUint(initialTotalSupply.Uint64() + 1000)
		k.SetForeignCoins(ctx, foreignCoin)

		// increase total supply
		_, err = k.DepositMRC20(ctx, mrc20, sample.EthAddress(), big.NewInt(500))
		require.NoError(t, err)

		// deposit (500 + 501 > 1000)
		to := sample.EthAddress()
		_, _, err = k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(501),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
			crosschaintypes.ProtocolContractVersion_V1,
			false,
		)
		require.ErrorIs(t, err, types.ErrForeignCoinCapReached)
	})

	t.Run("should fail if gas coin not found", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// deposit
		to := sample.EthAddress()
		_, _, err := k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
			crosschaintypes.ProtocolContractVersion_V1,
			false,
		)
		require.ErrorIs(t, err, crosschaintypes.ErrGasCoinNotFound)
	})

	t.Run("should fail if mrc20 not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId
		assetAddress := sample.EthAddress().String()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// deposit
		to := sample.EthAddress()
		_, _, err := k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_ERC20,
			assetAddress,
			crosschaintypes.ProtocolContractVersion_V1,
			false,
		)
		require.ErrorIs(t, err, crosschaintypes.ErrForeignCoinNotFound)
	})

	t.Run("should return contract call if receiver is a contract", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		example, err := k.DeployContract(ctx, example.ExampleMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, example)

		// deposit
		_, contractCall, err := k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			example,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
			crosschaintypes.ProtocolContractVersion_V1,
			false,
		)
		require.NoError(t, err)
		require.True(t, contractCall)

		balance, err := k.BalanceOfMRC4(ctx, mrc20, example)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)

		// check onCrossChainCall() hook was called
		assertExampleBarValue(t, ctx, k, example, 42)
	})

	t.Run("should fail if call contract fails", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		reverter, err := k.DeployContract(ctx, reverter.ReverterMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, reverter)

		// deposit
		_, contractCall, err := k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			reverter,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
			crosschaintypes.ProtocolContractVersion_V1,
			false,
		)
		require.Error(t, err)
		require.True(t, contractCall)

		balance, err := k.BalanceOfMRC4(ctx, mrc20, reverter)
		require.NoError(t, err)
		require.EqualValues(t, int64(0), balance.Int64())
	})

	t.Run("can deposit using V2", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// deposit
		to := sample.EthAddress()
		_, contractCall, err := k.MRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
			crosschaintypes.ProtocolContractVersion_V2,
			false,
		)
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfMRC4(ctx, mrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})
}

func TestKeeper_DepositCoinMuse(t *testing.T) {
	t.Run("successfully deposit coin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		to := sample.EthAddress()
		amount := big.NewInt(1)
		museToAddress := sdk.AccAddress(to.Bytes())

		b := sdkk.BankKeeper.GetBalance(ctx, museToAddress, config.BaseDenom)
		require.Equal(t, int64(0), b.Amount.Int64())

		err := k.DepositCoinMuse(ctx, to, amount)
		require.NoError(t, err)
		b = sdkk.BankKeeper.GetBalance(ctx, museToAddress, config.BaseDenom)
		require.Equal(t, amount.Int64(), b.Amount.Int64())
	})

	t.Run("should fail if MintMuseToEVMAccount fails", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{UseBankMock: true})
		bankMock := keepertest.GetFungibleBankMock(t, k)
		to := sample.EthAddress()
		amount := big.NewInt(1)
		museToAddress := sdk.AccAddress(to.Bytes())

		b := sdkk.BankKeeper.GetBalance(ctx, museToAddress, config.BaseDenom)
		require.Equal(t, int64(0), b.Amount.Int64())
		errorMint := errors.New("error minting coins")

		bankMock.On("GetSupply", ctx, mock.Anything, mock.Anything).
			Return(sdk.NewCoin(config.BaseDenom, math.NewInt(0))).
			Once()
		bankMock.On("MintCoins", ctx, types.ModuleName, mock.Anything).Return(errorMint).Once()
		err := k.DepositCoinMuse(ctx, to, amount)
		require.ErrorIs(t, err, errorMint)

	})
}

func TestKeeper_ProcessDeposit(t *testing.T) {
	t.Run("should process no-call deposit", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		receiver := sample.EthAddress()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		_, contractCall, err := k.ProcessDeposit(
			ctx,
			sample.EthAddress().Bytes(),
			chainID,
			mrc20,
			receiver,
			big.NewInt(42),
			[]byte{},
			coin.CoinType_Gas,
			false,
		)

		// ASSERT
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfMRC4(ctx, mrc20, receiver)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("should process no-call deposit, message should be ignored", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		receiver := sample.EthAddress()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		_, contractCall, err := k.ProcessDeposit(
			ctx,
			sample.EthAddress().Bytes(),
			chainID,
			mrc20,
			receiver,
			big.NewInt(42),
			[]byte("foo"),
			coin.CoinType_Gas,
			false,
		)

		// ASSERT
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfMRC4(ctx, mrc20, receiver)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("should process deposit and call", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId

		// deploy test dapp
		testDapp := deployTestDAppV2(t, ctx, k, sdkk.EvmKeeper)

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		_, contractCall, err := k.ProcessDeposit(
			ctx,
			sample.EthAddress().Bytes(),
			chainID,
			mrc20,
			testDapp,
			big.NewInt(82),
			[]byte("foo"),
			coin.CoinType_Gas,
			true,
		)

		// ASSERT
		require.NoError(t, err)
		require.True(t, contractCall)
		balance, err := k.BalanceOfMRC4(ctx, mrc20, testDapp)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(82), balance)
		assertTestDAppV2MessageAndAmount(t, ctx, k, testDapp, "foo", 82)
	})

	t.Run("should process deposit and call with no message", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId

		// deploy test dapp
		testDapp := deployTestDAppV2(t, ctx, k, sdkk.EvmKeeper)

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		sender := sample.EthAddress()

		// ACT
		_, contractCall, err := k.ProcessDeposit(
			ctx,
			sender.Bytes(),
			chainID,
			mrc20,
			testDapp,
			big.NewInt(82),
			[]byte{},
			coin.CoinType_Gas,
			true,
		)

		// ASSERT
		require.NoError(t, err)
		require.True(t, contractCall)
		balance, err := k.BalanceOfMRC4(ctx, mrc20, testDapp)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(82), balance)

		messageIndex := getTestDAppNoMessageIndex(t, ctx, *k, testDapp, sender)

		assertTestDAppV2MessageAndAmount(
			t,
			ctx,
			k,
			testDapp,
			messageIndex,
			82,
		)
	})
}

func TestKeeper_ProcessAbort(t *testing.T) {
	t.Run("should process abort with onAbort call", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId

		// deploy test dapp
		testAbort := deployTestAbort(t, ctx, k, sdkk.EvmKeeper)

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		_, err := k.ProcessAbort(
			ctx,
			sample.EthAddress().String(),
			big.NewInt(82),
			false,
			chainID,
			coin.CoinType_Gas,
			"",
			testAbort,
			[]byte("foo"),
		)

		// ASSERT
		require.NoError(t, err)
		balance, err := k.BalanceOfMRC4(ctx, mrc20, testAbort)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(82), balance)
	})

	t.Run("should return a onAbortFailError if onAbortFailed", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		mrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// onAbort will fail because the testAbort contract is not a valid contract
		abortAddress := sample.EthAddress()

		// ACT
		_, err := k.ProcessAbort(
			ctx,
			sample.EthAddress().String(),
			big.NewInt(82),
			false,
			chainID,
			coin.CoinType_Gas,
			"",
			abortAddress,
			[]byte("foo"),
		)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrOnAbortFailed)

		// account still founded
		balance, err := k.BalanceOfMRC4(ctx, mrc20, abortAddress)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(82), balance)
	})

	t.Run("can't process abort for MUSE token", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId

		// deploy test dapp
		testAbort := deployTestAbort(t, ctx, k, sdkk.EvmKeeper)

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// ACT
		_, err := k.ProcessAbort(
			ctx,
			sample.EthAddress().String(),
			big.NewInt(82),
			false,
			chainID,
			coin.CoinType_Muse,
			"",
			testAbort,
			[]byte("foo"),
		)

		// ASSERT
		require.Error(t, err)
		require.NotErrorIs(t, err, types.ErrOnAbortFailed)
	})

	t.Run("can't process abort for invalid chain ID", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// deploy test dapp
		testAbort := deployTestAbort(t, ctx, k, sdkk.EvmKeeper)

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// ACT
		_, err := k.ProcessAbort(
			ctx,
			sample.EthAddress().String(),
			big.NewInt(82),
			false,
			919191,
			coin.CoinType_Gas,
			"",
			testAbort,
			[]byte("foo"),
		)

		// ASSERT
		require.Error(t, err)
		require.NotErrorIs(t, err, types.ErrOnAbortFailed)
	})
}
