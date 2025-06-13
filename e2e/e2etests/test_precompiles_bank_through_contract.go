package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/testbank"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/precompiles/bank"
)

func TestPrecompilesBankThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	var (
		spender        = r.EVMAddress()
		bankAddress    = bank.ContractAddress
		mrc20Address   = r.ERC20MRC20Addr
		oneThousand    = big.NewInt(1e3)
		oneThousandOne = big.NewInt(1001)
		fiveHundred    = big.NewInt(500)
		fiveHundredOne = big.NewInt(501)
		zero           = big.NewInt(0)
	)

	// Get ERC20MRC20.
	txHash := r.LegacyDepositERC20WithAmountAndMessage(r.EVMAddress(), oneThousand, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	bankPrecompileCaller, err := bank.NewIBank(bank.ContractAddress, r.MEVMClient)
	require.NoError(r, err, "Failed to create bank precompile caller")

	// Deploy the TestBank. Ensure the transaction is successful.
	_, tx, testBank, err := testbank.DeployTestBank(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Deployment of TestBank contract failed")

	previousGasLimit := r.MEVMAuth.GasLimit
	r.MEVMAuth.GasLimit = 10_000_000
	defer func() {
		r.MEVMAuth.GasLimit = previousGasLimit

		// Reset the allowance to 0; this is needed when running upgrade tests where this test runs twice.
		approveAllowance(r, bank.ContractAddress, zero)

		// Reset balance to 0; this is needed when running upgrade tests where this test runs twice.
		tx, err = r.ERC20MRC20.Transfer(
			r.MEVMAuth,
			common.HexToAddress("0x000000000000000000000000000000000000dEaD"),
			oneThousand,
		)
		require.NoError(r, err)
		receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt, "Resetting balance failed")
	}()

	// always ensure allowance is set to zero before test starts
	approveAllowance(r, bank.ContractAddress, zero)

	// get starting balances
	startSpenderCosmosBalance := checkCosmosBalanceThroughBank(r, testBank, mrc20Address, spender)
	startSpenderMRC20Balance := checkMRC20Balance(r, spender)
	startBankMRC20Balance := checkMRC20Balance(r, bankAddress)

	// Deposit without previous alllowance should fail.
	receipt = depositThroughTestBank(r, testBank, mrc20Address, oneThousand)
	utils.RequiredTxFailed(r, receipt, "Deposit ERC20MRC20 without allowance should fail")

	// Check balances, should be the same.
	balanceShouldBe(r, startSpenderCosmosBalance, checkCosmosBalanceThroughBank(r, testBank, mrc20Address, spender))
	balanceShouldBe(r, startSpenderMRC20Balance, checkMRC20Balance(r, spender))
	balanceShouldBe(r, startBankMRC20Balance, checkMRC20Balance(r, bankAddress))

	// Allow 500 MRC20 to bank precompile.
	approveAllowance(r, bankAddress, fiveHundred)

	// Deposit 501 ERC20MRC20 tokens to the bank contract, through TestBank.
	// It's higher than allowance but lower than balance, should fail.
	receipt = depositThroughTestBank(r, testBank, mrc20Address, fiveHundredOne)
	utils.RequiredTxFailed(r, receipt, "Depositting an amount higher than allowed should fail")

	// Balances shouldn't change.
	balanceShouldBe(r, startSpenderCosmosBalance, checkCosmosBalanceThroughBank(r, testBank, mrc20Address, spender))
	balanceShouldBe(r, startSpenderMRC20Balance, checkMRC20Balance(r, spender))
	balanceShouldBe(r, startBankMRC20Balance, checkMRC20Balance(r, bankAddress))

	// Allow 1000 MRC20 to bank precompile.
	approveAllowance(r, bankAddress, oneThousand)

	// Deposit 1001 ERC20MRC20 tokens to the bank contract.
	// It's higher than spender balance but within approved allowance, should fail.
	receipt = depositThroughTestBank(r, testBank, mrc20Address, oneThousandOne)
	utils.RequiredTxFailed(r, receipt, "Depositting an amount higher than balance should fail")

	// Balances shouldn't change.
	balanceShouldBe(r, startSpenderCosmosBalance, checkCosmosBalanceThroughBank(r, testBank, mrc20Address, spender))
	balanceShouldBe(r, startSpenderMRC20Balance, checkMRC20Balance(r, spender))
	balanceShouldBe(r, startBankMRC20Balance, checkMRC20Balance(r, bankAddress))

	// Deposit 500 ERC20MRC20 tokens to the bank contract, it's within allowance and balance. Should pass.
	receipt = depositThroughTestBank(r, testBank, mrc20Address, fiveHundred)
	utils.RequireTxSuccessful(r, receipt, "Depositting a correct amount should pass")

	// Balances should be transferred. Bank now locks 500 MRC20 tokens.
	balanceShouldBe(
		r,
		bigAdd(startSpenderCosmosBalance, fiveHundred),
		checkCosmosBalanceThroughBank(r, testBank, mrc20Address, spender),
	)
	balanceShouldBe(r, bigSub(startSpenderMRC20Balance, fiveHundred), checkMRC20Balance(r, spender))
	balanceShouldBe(r, bigAdd(startBankMRC20Balance, fiveHundred), checkMRC20Balance(r, bankAddress))

	// Check the deposit event.
	eventDeposit, err := bankPrecompileCaller.ParseDeposit(*receipt.Logs[0])
	require.NoError(r, err, "Parse Deposit event")
	require.Equal(r, r.EVMAddress(), eventDeposit.Mrc20Depositor, "Deposit event token should be r.EVMAddress()")
	require.Equal(r, r.ERC20MRC20Addr, eventDeposit.Mrc20Token, "Deposit event token should be ERC20MRC20Addr")
	require.Equal(r, fiveHundred, eventDeposit.Amount, "Deposit event amount should be 500")

	// Should faild to withdraw more than cosmos balance.
	receipt = withdrawThroughTestBank(r, testBank, mrc20Address, bigAdd(startSpenderCosmosBalance, fiveHundredOne))
	utils.RequiredTxFailed(r, receipt, "Withdrawing an amount higher than balance should fail")

	// Balances shouldn't change.
	balanceShouldBe(
		r,
		bigAdd(startSpenderCosmosBalance, fiveHundred),
		checkCosmosBalanceThroughBank(r, testBank, mrc20Address, spender),
	)
	balanceShouldBe(r, bigSub(startSpenderMRC20Balance, fiveHundred), checkMRC20Balance(r, spender))
	balanceShouldBe(r, bigAdd(startBankMRC20Balance, fiveHundred), checkMRC20Balance(r, bankAddress))

	// Try to withdraw 500 ERC20MRC20 tokens. Should pass.
	receipt = withdrawThroughTestBank(r, testBank, mrc20Address, fiveHundred)
	utils.RequireTxSuccessful(r, receipt, "Withdraw correct amount should pass")

	// Balances should be reverted to initial state.
	balanceShouldBe(r, startSpenderCosmosBalance, checkCosmosBalanceThroughBank(r, testBank, mrc20Address, spender))
	balanceShouldBe(r, startSpenderMRC20Balance, checkMRC20Balance(r, spender))
	balanceShouldBe(r, startBankMRC20Balance, checkMRC20Balance(r, bankAddress))

	// Check the withdraw event.
	eventWithdraw, err := bankPrecompileCaller.ParseWithdraw(*receipt.Logs[0])
	require.NoError(r, err, "Parse Withdraw event")
	require.Equal(r, r.EVMAddress(), eventWithdraw.Mrc20Withdrawer, "Withdrawer should be r.EVMAddress()")
	require.Equal(r, r.ERC20MRC20Addr, eventWithdraw.Mrc20Token, "Withdraw event token should be ERC20MRC20Addr")
	require.Equal(r, fiveHundred, eventWithdraw.Amount, "Withdraw event amount should be 500")
}

func approveAllowance(r *runner.E2ERunner, target common.Address, amount *big.Int) {
	tx, err := r.ERC20MRC20.Approve(r.MEVMAuth, target, amount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Approve ERC20MRC20 allowance tx failed")
}

func balanceShouldBe(r *runner.E2ERunner, expected *big.Int, balance *big.Int) {
	require.Equal(r, expected.Uint64(), balance.Uint64(), "Balance should be %d, got: %d", expected, balance.Uint64())
}

func checkMRC20Balance(r *runner.E2ERunner, target common.Address) *big.Int {
	bankMRC20Balance, err := r.ERC20MRC20.BalanceOf(&bind.CallOpts{Context: r.Ctx}, target)
	require.NoError(r, err, "Call ERC20MRC20.BalanceOf")
	return bankMRC20Balance
}

func checkCosmosBalanceThroughBank(
	r *runner.E2ERunner,
	bank *testbank.TestBank,
	mrc20, target common.Address,
) *big.Int {
	balance, err := bank.BalanceOf(&bind.CallOpts{Context: r.Ctx, From: r.MEVMAuth.From}, mrc20, target)
	require.NoError(r, err)
	return balance
}

func depositThroughTestBank(
	r *runner.E2ERunner,
	bank *testbank.TestBank,
	mrc20Address common.Address,
	amount *big.Int,
) *types.Receipt {
	tx, err := bank.Deposit(r.MEVMAuth, mrc20Address, amount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	return receipt
}

func withdrawThroughTestBank(
	r *runner.E2ERunner,
	bank *testbank.TestBank,
	mrc20Address common.Address,
	amount *big.Int,
) *types.Receipt {
	tx, err := bank.Withdraw(r.MEVMAuth, mrc20Address, amount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	return receipt
}
