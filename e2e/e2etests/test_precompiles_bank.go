package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/precompiles/bank"
)

func TestPrecompilesBank(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	totalAmount := big.NewInt(1e3)
	depositAmount := big.NewInt(500)
	higherBalanceAmount := big.NewInt(1001)
	higherAllowanceAmount := big.NewInt(501)
	spender := r.EVMAddress()
	bankAddress := bank.ContractAddress

	// Increase the gasLimit. It's required because of the gas consumed by precompiled functions.
	previousGasLimit := r.MEVMAuth.GasLimit
	r.MEVMAuth.GasLimit = 10_000_000
	defer func() {
		r.MEVMAuth.GasLimit = previousGasLimit

		// Reset the allowance to 0; this is needed when running upgrade tests where
		// this test runs twice.
		tx, err := r.ERC20MRC20.Approve(r.MEVMAuth, bankAddress, big.NewInt(0))
		require.NoError(r, err)
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt, "Resetting allowance failed")

		// Reset balance to 0; this is needed when running upgrade tests where this test runs twice.
		tx, err = r.ERC20MRC20.Transfer(
			r.MEVMAuth,
			common.HexToAddress("0x000000000000000000000000000000000000dEaD"),
			totalAmount,
		)
		require.NoError(r, err)
		receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt, "Resetting balance failed")
	}()

	// Ensure starting allowance is zero; this is needed when running the tests multiple times
	tx, err := r.ERC20MRC20.Approve(r.MEVMAuth, bankAddress, big.NewInt(0))
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Resetting allowance failed")

	// Get ERC20MRC20.
	txHash := r.LegacyDepositERC20WithAmountAndMessage(r.EVMAddress(), totalAmount, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	// Create a bank contract caller.
	bankContract, err := bank.NewIBank(bank.ContractAddress, r.MEVMClient)
	require.NoError(r, err, "Failed to create bank contract caller")

	// get starting balances
	startSpenderCosmosBalance, err := bankContract.BalanceOf(&bind.CallOpts{Context: r.Ctx}, r.ERC20MRC20Addr, spender)
	require.NoError(r, err, "Call bank.BalanceOf()")
	startSpenderMRC20Balance, err := r.ERC20MRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, spender)
	require.NoError(r, err, "Call bank.BalanceOf()")
	startBankMRC20Balance, err := r.ERC20MRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, bankAddress)
	require.NoError(r, err, "Call ERC20MRC20.BalanceOf")

	// Approve allowance of 500 ERC20MRC20 tokens for the bank contract. Should pass.
	tx, err = r.ERC20MRC20.Approve(r.MEVMAuth, bankAddress, depositAmount)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Approve ETHMRC20 bank allowance tx failed")

	// Deposit 501 ERC20MRC20 tokens to the bank contract.
	// It's higher than allowance but lower than balance, should fail.
	tx, err = bankContract.Deposit(r.MEVMAuth, r.ERC20MRC20Addr, higherAllowanceAmount)
	require.NoError(r, err, "Call bank.Deposit() with amout higher than allowance")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "Depositting an amount higher than allowed should fail")

	// Approve allowance of 1000 ERC20MRC20 tokens.
	tx, err = r.ERC20MRC20.Approve(r.MEVMAuth, bankAddress, big.NewInt(1e3))
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Approve ETHMRC20 bank allowance tx failed")

	// Deposit 1001 ERC20MRC20 tokens to the bank contract.
	// It's higher than spender balance but within approved allowance, should fail.
	tx, err = bankContract.Deposit(r.MEVMAuth, r.ERC20MRC20Addr, higherBalanceAmount)
	require.NoError(r, err, "Call bank.Deposit() with amout higher than balance")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "Depositting an amount higher than balance should fail")

	// Deposit 500 ERC20MRC20 tokens to the bank contract. Should pass.
	tx, err = bankContract.Deposit(r.MEVMAuth, r.ERC20MRC20Addr, depositAmount)
	require.NoError(r, err, "Call bank.Deposit() with correct amount")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Depositting a correct amount should pass")

	// Check the deposit event.
	eventDeposit, err := bankContract.ParseDeposit(*receipt.Logs[0])
	require.NoError(r, err, "Parse Deposit event")
	require.Equal(r, r.EVMAddress(), eventDeposit.Mrc20Depositor, "Deposit event token should be r.EVMAddress()")
	require.Equal(r, r.ERC20MRC20Addr, eventDeposit.Mrc20Token, "Deposit event token should be ERC20MRC20Addr")
	require.Equal(r, depositAmount, eventDeposit.Amount, "Deposit event amount should be 500")

	// Spender: cosmos coin balance should be 500 at this point.
	spenderCosmosBalance, err := bankContract.BalanceOf(&bind.CallOpts{Context: r.Ctx}, r.ERC20MRC20Addr, spender)
	require.NoError(r, err, "Call bank.BalanceOf()")
	require.Equal(
		r,
		startSpenderCosmosBalance.Int64()+500,
		spenderCosmosBalance.Int64(),
		"spender cosmos coin balance should be +500",
	)

	// Bank: ERC20MRC20 balance should be 500 tokens locked.
	bankMRC20Balance, err := r.ERC20MRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, bankAddress)
	require.NoError(r, err, "Call ERC20MRC20.BalanceOf")
	require.Equal(
		r,
		startBankMRC20Balance.Int64()+500,
		bankMRC20Balance.Int64(),
		"bank ERC20MRC20 balance should be +500",
	)

	// Try to withdraw one more than current balance. Should fail.
	tx, err = bankContract.Withdraw(r.MEVMAuth, r.ERC20MRC20Addr, new(big.Int).Add(spenderCosmosBalance, big.NewInt(1)))
	require.NoError(r, err, "Error calling bank.withdraw()")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "Withdrawing more than cosmos coin balance amount should fail")

	// Bank: ERC20MRC20 balance should be 500 tokens locked after a failed withdraw.
	// No tokens should be unlocked with a failed withdraw.
	bankMRC20Balance, err = r.ERC20MRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, bankAddress)
	require.NoError(r, err, "Call ERC20MRC20.BalanceOf")
	require.Equal(
		r,
		startBankMRC20Balance.Int64()+500,
		bankMRC20Balance.Int64(),
		"bank ERC20MRC20 balance should be +500",
	)

	// Try to withdraw 500 ERC20MRC20 tokens. Should pass.
	tx, err = bankContract.Withdraw(r.MEVMAuth, r.ERC20MRC20Addr, depositAmount)
	require.NoError(r, err, "Error calling bank.withdraw()")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Withdraw correct amount should pass")

	// Check the withdraw event.
	eventWithdraw, err := bankContract.ParseWithdraw(*receipt.Logs[0])
	require.NoError(r, err, "Parse Withdraw event")
	require.Equal(r, r.EVMAddress(), eventWithdraw.Mrc20Withdrawer, "Withdrawer should be r.EVMAddress()")
	require.Equal(r, r.ERC20MRC20Addr, eventWithdraw.Mrc20Token, "Withdraw event token should be ERC20MRC20Addr")
	require.Equal(r, depositAmount, eventWithdraw.Amount, "Withdraw event amount should be 500")

	// Spender: cosmos coin balance should be +0 at this point.
	spenderCosmosBalance, err = bankContract.BalanceOf(&bind.CallOpts{Context: r.Ctx}, r.ERC20MRC20Addr, spender)
	require.NoError(r, err, "Call bank.BalanceOf()")
	require.Equal(
		r,
		startSpenderCosmosBalance.Int64(),
		spenderCosmosBalance.Int64(),
		"spender cosmos coin balance should match starting balance",
	)

	// Spender: ERC20MRC20 balance should be +0 at this point.
	spenderMRC20Balance, err := r.ERC20MRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, spender)
	require.NoError(r, err, "Call ERC20MRC20.BalanceOf")
	require.Equal(
		r,
		startSpenderMRC20Balance.Int64(),
		spenderMRC20Balance.Int64(),
		"spender ERC20MRC20 balance should match starting balance",
	)

	// Bank: ERC20MRC20 balance should be +0 tokens locked.
	bankMRC20Balance, err = r.ERC20MRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, bankAddress)
	require.NoError(r, err, "Call ERC20MRC20.BalanceOf")
	require.Equal(
		r,
		startBankMRC20Balance.Int64(),
		bankMRC20Balance.Int64(),
		"bank ERC20MRC20 balance should match starting balance",
	)
}

func TestPrecompilesBankNonMRC20(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	// Increase the gasLimit. It's required because of the gas consumed by precompiled functions.
	previousGasLimit := r.MEVMAuth.GasLimit
	r.MEVMAuth.GasLimit = 10_000_000
	defer func() {
		r.MEVMAuth.GasLimit = previousGasLimit
	}()

	spender, bankAddress := r.EVMAddress(), bank.ContractAddress

	// Create a bank contract caller.
	bankContract, err := bank.NewIBank(bank.ContractAddress, r.MEVMClient)
	require.NoError(r, err, "Failed to create bank contract caller")

	// Deposit and approve 50 WMUSE for the test.
	approveAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(50))
	r.LegacyDepositAndApproveWMuse(approveAmount)

	// Non MRC20 balanceOf check should fail.
	_, err = bankContract.BalanceOf(&bind.CallOpts{Context: r.Ctx}, r.WMuseAddr, spender)
	require.Error(r, err, "bank.balanceOf() should error out when checking for non MRC20 balance")
	require.Contains(
		r,
		err.Error(),
		"invalid token 0x5F0b1a82749cb4E2278EC87F8BF6B618dC71a8bf: token is not a whitelisted MRC20",
		"Error should be 'token is not a whitelisted MRC20'",
	)

	// Allow the bank contract to spend 25 WMuse tokens.
	tx, err := r.WMuse.Approve(r.MEVMAuth, bankAddress, big.NewInt(25))
	require.NoError(r, err, "Error approving allowance for bank contract")
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.EqualValues(r, uint64(1), receipt.Status, "approve allowance tx failed")

	// Check the allowance of the bank in WMuse tokens. Should be 25.
	allowance, err := r.WMuse.Allowance(&bind.CallOpts{Context: r.Ctx}, spender, bankAddress)
	require.NoError(r, err, "Error retrieving bank allowance")
	require.EqualValues(r, uint64(25), allowance.Uint64(), "Error allowance for bank contract")

	// Call Deposit with 25 Non MRC20 tokens. Should fail.
	tx, err = bankContract.Deposit(r.MEVMAuth, r.WMuseAddr, big.NewInt(25))
	require.NoError(r, err, "Error calling bank.deposit()")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.Equal(r, uint64(0), receipt.Status, "Non MRC20 deposit should fail")

	// Call Withdraw with 25 on MRC20 tokens. Should fail.
	tx, err = bankContract.Withdraw(r.MEVMAuth, r.WMuseAddr, big.NewInt(25))
	require.NoError(r, err, "Error calling bank.withdraw()")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.Equal(r, uint64(0), receipt.Status, "Non MRC20 withdraw should fail")
}
