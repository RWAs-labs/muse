package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/withdrawer"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestBitcoinDepositAndWithdrawWithDust deposits Bitcoin and call a smart contract that withdraw dust amount
// It tests the edge case where during a cross-chain call, a invaild withdraw is initiated (processLogs fails)
func TestBitcoinDepositAndWithdrawWithDust(r *runner.E2ERunner, args []string) {
	// Given "Live" BTC network
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	require.Len(r, args, 0)

	// ARRANGE
	// Deploy the withdrawer contract on MuseChain with a withdraw amount of 100 satoshis (dust amount is 1000 satoshis)
	withdrawerAddr, tx, _, err := withdrawer.DeployWithdrawer(r.MEVMAuth, r.MEVMClient, big.NewInt(100))
	require.NoError(r, err)

	// Wait for the transaction to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.Equal(r, receipt.Status, uint64(1))

	// ACT
	// Deposit 0.01 BTC to withdrawer, this is an arbitrary amount, must be greater than dust amount
	txHash, err := r.SendToTSSWithMemo(
		0.01,
		append(withdrawerAddr.Bytes(), []byte("payload")...),
	)
	require.NoError(r, err)
	require.NotEmpty(r, txHash)

	// ASSERT
	// Now we want to make sure the cctx is reverted with expected error message
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.Equal(r, crosschaintypes.CctxStatus_Reverted, cctx.CctxStatus.Status)
	require.Contains(r, cctx.CctxStatus.ErrorMessage, crosschaintypes.ErrCannotProcessWithdrawal.Error())

	// check the contract has no BTC balance, this would mean the contract call state transition is not reverted
	// get BTC MRC20 balance of the withdrawer contract
	bal, err := r.BTCMRC20.BalanceOf(&bind.CallOpts{}, withdrawerAddr)
	require.NoError(r, err)
	require.Zero(r, bal.Uint64())
}
