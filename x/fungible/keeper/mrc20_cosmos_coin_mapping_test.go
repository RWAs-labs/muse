package keeper_test

import (
	"math/big"
	"testing"

	"github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	fungiblekeeper "github.com/RWAs-labs/muse/x/fungible/keeper"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_LockMRC20(t *testing.T) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	require.NoError(t, err)

	ts := setupChain(t)

	owner := fungibletypes.ModuleAddressEVM
	locker := sample.EthAddress()
	depositTotal := big.NewInt(1000)
	allowanceTotal := big.NewInt(100)
	higherThanAllowance := big.NewInt(101)
	smallerThanAllowance := big.NewInt(99)

	// Make sure locker account exists in state.
	accAddress := sdk.AccAddress(locker.Bytes())
	acc := ts.fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ts.ctx, accAddress)
	ts.fungibleKeeper.GetAuthKeeper().SetAccount(ts.ctx, acc)

	// Deposit 1000 MRC20 tokens into the fungible.
	ts.fungibleKeeper.DepositMRC20(ts.ctx, ts.mrc20Address, owner, depositTotal)

	t.Run("should fail when trying to lock zero amount", func(t *testing.T) {
		// Check lock with zero amount.
		err = ts.fungibleKeeper.LockMRC20(ts.ctx, ts.mrc20Address, locker, owner, locker, big.NewInt(0))
		require.Error(t, err)
		require.ErrorIs(t, err, fungibletypes.ErrInvalidAmount)
	})

	t.Run("should fail when trying to lock a zero address MRC20", func(t *testing.T) {
		// Check lock with MRC20 zero address.
		err = ts.fungibleKeeper.LockMRC20(ts.ctx, common.Address{}, locker, owner, locker, big.NewInt(10))
		require.Error(t, err)
		require.ErrorIs(t, err, fungibletypes.ErrMRC20ZeroAddress)
	})

	t.Run("should fail when trying to lock a non whitelisted MRC20", func(t *testing.T) {
		// Check lock with non whitelisted MRC20.
		err = ts.fungibleKeeper.LockMRC20(ts.ctx, sample.EthAddress(), locker, owner, locker, big.NewInt(10))
		require.Error(t, err)
		require.ErrorIs(t, err, fungibletypes.ErrMRC20NotWhiteListed)
	})

	t.Run("should fail when trying to lock a higher amount than totalSupply", func(t *testing.T) {
		approveAllowance(t, ts, mrc20ABI, owner, locker, big.NewInt(1000000000000000))

		// Check lock with higher amount than totalSupply.
		err = ts.fungibleKeeper.LockMRC20(
			ts.ctx,
			ts.mrc20Address,
			locker,
			owner,
			locker,
			big.NewInt(1000000000000000),
		)
		require.Error(t, err)
		require.ErrorIs(t, err, fungibletypes.ErrInvalidAmount)
	})

	t.Run("should fail when trying to lock a higher amount than owned balance", func(t *testing.T) {
		approveAllowance(t, ts, mrc20ABI, owner, locker, big.NewInt(1001))

		// Check allowance smaller, equal and bigger than the amount.
		err = ts.fungibleKeeper.LockMRC20(ts.ctx, ts.mrc20Address, locker, owner, locker, big.NewInt(1001))
		require.Error(t, err)

		// We do not check in LockMRC20 explicitly if the amount is bigger than the balance.
		// Instead, the ERC20 transferFrom function will revert the transaction if the amount is bigger than the balance.
		require.Contains(t, err.Error(), "execution reverted")
	})

	t.Run("should fail when trying to lock an amount higher than approved", func(t *testing.T) {
		approveAllowance(t, ts, mrc20ABI, owner, locker, allowanceTotal)

		// Check allowance smaller, equal and bigger than the amount.
		err = ts.fungibleKeeper.LockMRC20(ts.ctx, ts.mrc20Address, locker, owner, locker, higherThanAllowance)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid allowance, got 100")
	})

	t.Run("should pass when trying to lock a valid approved amount", func(t *testing.T) {
		approveAllowance(t, ts, mrc20ABI, owner, locker, allowanceTotal)

		err = ts.fungibleKeeper.LockMRC20(ts.ctx, ts.mrc20Address, locker, owner, locker, allowanceTotal)
		require.NoError(t, err)

		ownerBalance, err := ts.fungibleKeeper.MRC20BalanceOf(ts.ctx, ts.mrc20Address, owner)
		require.NoError(t, err)
		require.Equal(t, uint64(900), ownerBalance.Uint64())

		lockerBalance, err := ts.fungibleKeeper.MRC20BalanceOf(ts.ctx, ts.mrc20Address, locker)
		require.NoError(t, err)
		require.Equal(t, uint64(100), lockerBalance.Uint64())
	})

	t.Run("should pass when trying to lock an amount smaller than approved", func(t *testing.T) {
		approveAllowance(t, ts, mrc20ABI, owner, locker, allowanceTotal)

		err = ts.fungibleKeeper.LockMRC20(
			ts.ctx,
			ts.mrc20Address,
			locker,
			owner,
			locker,
			smallerThanAllowance,
		)
		require.NoError(t, err)

		// Note that balances are cumulative for all tests. That's why we check 801 and 199 here.
		ownerBalance, err := ts.fungibleKeeper.MRC20BalanceOf(ts.ctx, ts.mrc20Address, owner)
		require.NoError(t, err)
		require.Equal(t, uint64(801), ownerBalance.Uint64())

		lockerBalance, err := ts.fungibleKeeper.MRC20BalanceOf(ts.ctx, ts.mrc20Address, locker)
		require.NoError(t, err)
		require.Equal(t, uint64(199), lockerBalance.Uint64())
	})
}

func Test_UnlockMRC20(t *testing.T) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	require.NoError(t, err)

	ts := setupChain(t)

	owner := fungibletypes.ModuleAddressEVM
	locker := sample.EthAddress()
	depositTotal := big.NewInt(1000)
	allowanceTotal := big.NewInt(100)

	// Make sure locker account exists in state.
	accAddress := sdk.AccAddress(locker.Bytes())
	acc := ts.fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ts.ctx, accAddress)
	ts.fungibleKeeper.GetAuthKeeper().SetAccount(ts.ctx, acc)

	// Deposit 1000 MRC20 tokens into the fungible.
	ts.fungibleKeeper.DepositMRC20(ts.ctx, ts.mrc20Address, owner, depositTotal)

	// Approve allowance for locker to spend owner's MRC20 tokens.
	approveAllowance(t, ts, mrc20ABI, owner, locker, allowanceTotal)

	// Lock 100 MRC20.
	err = ts.fungibleKeeper.LockMRC20(ts.ctx, ts.mrc20Address, locker, owner, locker, allowanceTotal)
	require.NoError(t, err)

	t.Run("should fail when trying to unlock zero amount", func(t *testing.T) {
		err = ts.fungibleKeeper.UnlockMRC20(ts.ctx, ts.mrc20Address, owner, locker, big.NewInt(0))
		require.Error(t, err)
		require.ErrorIs(t, err, fungibletypes.ErrInvalidAmount)
	})

	t.Run("should fail when trying to unlock a zero address MRC20", func(t *testing.T) {
		err = ts.fungibleKeeper.UnlockMRC20(ts.ctx, common.Address{}, owner, locker, big.NewInt(10))
		require.Error(t, err)
		require.ErrorIs(t, err, fungibletypes.ErrMRC20ZeroAddress)
	})

	t.Run("should fail when trying to unlock a non whitelisted MRC20", func(t *testing.T) {
		err = ts.fungibleKeeper.UnlockMRC20(ts.ctx, sample.EthAddress(), owner, locker, big.NewInt(10))
		require.Error(t, err)
		require.ErrorIs(t, err, fungibletypes.ErrMRC20NotWhiteListed)
	})

	t.Run("should fail when trying to unlock an amount bigger than locker's balance", func(t *testing.T) {
		err = ts.fungibleKeeper.UnlockMRC20(ts.ctx, ts.mrc20Address, owner, locker, big.NewInt(1001))
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid balance, got 100")
	})

	t.Run("should pass when trying to unlock a correct amount", func(t *testing.T) {
		err = ts.fungibleKeeper.UnlockMRC20(ts.ctx, ts.mrc20Address, owner, locker, allowanceTotal)
		require.NoError(t, err)

		ownerBalance, err := ts.fungibleKeeper.MRC20BalanceOf(ts.ctx, ts.mrc20Address, owner)
		require.NoError(t, err)
		require.Equal(t, uint64(1000), ownerBalance.Uint64())

		lockerBalance, err := ts.fungibleKeeper.MRC20BalanceOf(ts.ctx, ts.mrc20Address, locker)
		require.NoError(t, err)
		require.Equal(t, uint64(0), lockerBalance.Uint64())
	})
}

func Test_CheckMRC20Allowance(t *testing.T) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	require.NoError(t, err)

	ts := setupChain(t)

	owner := fungibletypes.ModuleAddressEVM
	spender := sample.EthAddress()
	depositTotal := big.NewInt(1000)
	allowanceTotal := big.NewInt(100)
	higherThanAllowance := big.NewInt(101)
	smallerThanAllowance := big.NewInt(99)

	// Make sure locker account exists in state.
	accAddress := sdk.AccAddress(spender.Bytes())
	acc := ts.fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ts.ctx, accAddress)
	ts.fungibleKeeper.GetAuthKeeper().SetAccount(ts.ctx, acc)

	// Deposit MRC20 tokens into the fungible.
	ts.fungibleKeeper.DepositMRC20(ts.ctx, ts.mrc20Address, fungibletypes.ModuleAddressEVM, depositTotal)

	t.Run("should fail when checking zero amount", func(t *testing.T) {
		err = ts.fungibleKeeper.CheckMRC20Allowance(ts.ctx, owner, spender, ts.mrc20Address, big.NewInt(0))
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrInvalidAmount)
	})

	t.Run("should fail when allowance is not approved", func(t *testing.T) {
		err = ts.fungibleKeeper.CheckMRC20Allowance(ts.ctx, owner, spender, ts.mrc20Address, big.NewInt(10))
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid allowance, got 0")
	})

	t.Run("should fail when checking a higher amount than approved", func(t *testing.T) {
		approveAllowance(t, ts, mrc20ABI, owner, spender, allowanceTotal)

		err = ts.fungibleKeeper.CheckMRC20Allowance(
			ts.ctx,
			owner,
			spender,
			ts.mrc20Address,
			higherThanAllowance,
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid allowance, got 100")
	})

	t.Run("should pass when checking the same amount as approved", func(t *testing.T) {
		approveAllowance(t, ts, mrc20ABI, owner, spender, allowanceTotal)

		err = ts.fungibleKeeper.CheckMRC20Allowance(ts.ctx, owner, spender, ts.mrc20Address, allowanceTotal)
		require.NoError(t, err)
	})

	t.Run("should pass when checking a lower amount than approved", func(t *testing.T) {
		approveAllowance(t, ts, mrc20ABI, owner, spender, allowanceTotal)

		err = ts.fungibleKeeper.CheckMRC20Allowance(
			ts.ctx,
			owner,
			spender,
			ts.mrc20Address,
			smallerThanAllowance,
		)
		require.NoError(t, err)
	})
}

func Test_IsValidMRC20(t *testing.T) {
	ts := setupChain(t)

	t.Run("should fail when mrc20 address is zero", func(t *testing.T) {
		err := ts.fungibleKeeper.IsValidMRC20(ts.ctx, common.Address{})
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when mrc20 is not whitelisted", func(t *testing.T) {
		err := ts.fungibleKeeper.IsValidMRC20(ts.ctx, sample.EthAddress())
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrMRC20NotWhiteListed)
	})

	t.Run("should pass when mrc20 is a valid whitelisted token", func(t *testing.T) {
		err := ts.fungibleKeeper.IsValidMRC20(ts.ctx, ts.mrc20Address)
		require.NoError(t, err)
	})
}

func Test_IsValidDepositAmount(t *testing.T) {
	ts := setupChain(t)

	t.Run("should fail when any input is nil", func(t *testing.T) {
		isValid := ts.fungibleKeeper.IsValidDepositAmount(nil, big.NewInt(0), big.NewInt(0))
		require.False(t, isValid)

		isValid = ts.fungibleKeeper.IsValidDepositAmount(big.NewInt(0), nil, big.NewInt(0))
		require.False(t, isValid)

		isValid = ts.fungibleKeeper.IsValidDepositAmount(big.NewInt(0), big.NewInt(0), nil)
		require.False(t, isValid)
	})

	t.Run("should fail when alreadyLocked + amountToDeposit > totalSupply", func(t *testing.T) {
		isValid := ts.fungibleKeeper.IsValidDepositAmount(big.NewInt(1000), big.NewInt(500), big.NewInt(501))
		require.False(t, isValid)
	})

	t.Run("should pass when alreadyLocked + amountToDeposit = totalSupply", func(t *testing.T) {
		isValid := ts.fungibleKeeper.IsValidDepositAmount(big.NewInt(1000), big.NewInt(500), big.NewInt(500))
		require.True(t, isValid)
	})

	t.Run("should pass when alreadyLocked + amountToDeposit < totalSupply", func(t *testing.T) {
		isValid := ts.fungibleKeeper.IsValidDepositAmount(big.NewInt(1000), big.NewInt(500), big.NewInt(499))
		require.True(t, isValid)
	})
}

/*
	Test utils.
*/

type testSuite struct {
	ctx            sdk.Context
	fungibleKeeper *fungiblekeeper.Keeper
	sdkKeepers     keeper.SDKKeepers
	mrc20Address   common.Address
}

func setupChain(t *testing.T) testSuite {
	// Initialize basic parameters to mock the chain.
	fungibleKeeper, ctx, sdkKeepers, _ := keeper.FungibleKeeper(t)
	chainID := getValidChainID(t)

	// Make sure the account store is initialized.
	// This is completely needed for accounts to be created in the state.
	fungibleKeeper.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

	// Deploy system contracts in order to deploy a MRC20 token.
	deploySystemContracts(t, ctx, fungibleKeeper, sdkKeepers.EvmKeeper)

	mrc20Address := setupGasCoin(t, ctx, fungibleKeeper, sdkKeepers.EvmKeeper, chainID, "MRC20", "MRC20")

	return testSuite{
		ctx,
		fungibleKeeper,
		sdkKeepers,
		mrc20Address,
	}
}

func approveAllowance(t *testing.T, ts testSuite, mrc20ABI *abi.ABI, owner, spender common.Address, amount *big.Int) {
	resAllowance, err := callEVM(
		t,
		ts.ctx,
		ts.fungibleKeeper,
		mrc20ABI,
		owner,
		ts.mrc20Address,
		"approve",
		[]interface{}{spender, amount},
	)
	require.NoError(t, err, "error allowing bank to spend MRC20 tokens")

	allowed, ok := resAllowance[0].(bool)
	require.True(t, ok)
	require.True(t, allowed)
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
