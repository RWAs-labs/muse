package legacy

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	testcontract "github.com/RWAs-labs/muse/e2e/contracts/withdrawerv2"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
)

func TestMultipleERC20Withdraws(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	approvedAmount := big.NewInt(1e18)

	// parse the withdrawal amount and number of withdrawals
	withdrawalAmount := utils.ParseBigInt(r, args[0])
	require.Equal(
		r,
		-1,
		withdrawalAmount.Cmp(approvedAmount),
		"Invalid withdrawal amount specified for TestMultipleWithdraws.",
	)
	numberOfWithdrawals := utils.ParseBigInt(r, args[1])
	require.NotEmpty(r, numberOfWithdrawals.Int64())

	// calculate total withdrawal to ensure it doesn't exceed approved amount.
	totalWithdrawal := big.NewInt(0).Mul(withdrawalAmount, numberOfWithdrawals)
	require.Equal(r, -1, totalWithdrawal.Cmp(approvedAmount), "Total withdrawal amount exceeds approved limit.")

	// deploy withdrawer
	withdrawerAddr, _, withdrawer, err := testcontract.DeployWithdrawer(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)

	// approve
	tx, err := r.ERC20MRC20.Approve(r.MEVMAuth, withdrawerAddr, approvedAmount)
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("ERC20 MRC20 approve receipt: status %d", receipt.Status)

	// approve gas token
	tx, err = r.ETHMRC20.Approve(r.MEVMAuth, withdrawerAddr, approvedAmount)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("eth mrc20 approve receipt: status %d", receipt.Status)

	// check the balance
	bal, err := r.ERC20MRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("balance of deployer on ERC20 MRC20: %d", bal)

	require.Less(r, totalWithdrawal.Int64(), bal.Int64(), "not enough ERC20 MRC20 balance!")

	// withdraw
	tx, err = withdrawer.RunWithdraws(
		r.MEVMAuth,
		r.EVMAddress().Bytes(),
		r.ERC20MRC20Addr,
		withdrawalAmount,
		numberOfWithdrawals,
	)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	cctxs := utils.WaitCctxsMinedByInboundHash(
		r.Ctx,
		tx.Hash().Hex(),
		r.CctxClient,
		int(numberOfWithdrawals.Int64()),
		r.Logger,
		r.CctxTimeout,
	)
	require.Len(r, cctxs, 3)

	// verify the withdraw value
	for _, cctx := range cctxs {
		r.EVMVerifyOutboundTransferAmount(cctx.GetCurrentOutboundParam().Hash, withdrawalAmount.Int64())
	}
}
