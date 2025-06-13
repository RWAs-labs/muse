package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
)

func TestBitcoinWithdrawToInvalidAddress(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawalAmount := utils.ParseFloat(r, args[0])
	amount := utils.BTCAmountFromFloat64(r, withdrawalAmount)

	withdrawToInvalidAddress(r, amount)
}

func withdrawToInvalidAddress(r *runner.E2ERunner, amount *big.Int) {
	approvalAmount := 1000000000000000000
	// approve the MRC20 contract to spend approvalAmount BTC from the deployer address.
	// the actual amount transferred is provided as test arg BTC, but we approve more to cover withdraw fee
	tx, err := r.BTCMRC20.Approve(r.MEVMAuth, r.BTCMRC20Addr, big.NewInt(int64(approvalAmount)))
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// mine blocks if testing on regnet
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// withdraw amount provided as test arg BTC from MRC20 to BTC legacy address
	// the address "1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3" is for mainnet, not regtest
	tx, err = r.BTCMRC20.Withdraw(r.MEVMAuth, []byte("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3"), amount)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)
}
