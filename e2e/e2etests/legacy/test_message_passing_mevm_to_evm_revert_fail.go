package legacy

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/testdappnorevert"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	cctxtypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestMessagePassingMEVMtoEVMRevertFail(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := utils.ParseBigInt(r, args[0])

	// Deploying a test contract not containing a logic for reverting the cctx
	testDappNoRevertAddr, tx, testDappNoRevert, err := testdappnorevert.DeployTestDAppNoRevert(
		r.MEVMAuth,
		r.MEVMClient,
		r.ConnectorMEVMAddr,
		r.WMuseAddr,
	)
	require.NoError(r, err)

	r.Logger.Info("TestDAppNoRevert deployed at: %s", testDappNoRevertAddr.Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "deploy TestDAppNoRevert")
	utils.RequireTxSuccessful(r, receipt)

	// Set destination details
	EVMChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	destinationAddress := r.EvmTestDAppAddr

	// Contract call originates from MEVM chain
	r.MEVMAuth.Value = amount
	tx, err = r.WMuse.Deposit(r.MEVMAuth)
	require.NoError(r, err)

	r.MEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wmuse deposit tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wmuse deposit")
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.WMuse.Approve(r.MEVMAuth, testDappNoRevertAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("wmuse approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wmuse approve")
	utils.RequireTxSuccessful(r, receipt)

	// Get previous balances to check funds are not minted anywhere when aborted
	previousBalanceMEVM, err := r.WMuse.BalanceOf(&bind.CallOpts{}, testDappNoRevertAddr)
	require.NoError(r, err)

	// Send message with doRevert
	tx, err = testDappNoRevert.SendHelloWorld(r.MEVMAuth, destinationAddress, EVMChainID, amount, true)
	require.NoError(r, err)

	r.Logger.Info("TestDAppNoRevert.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// The revert tx will fail, the cctx state should be aborted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_Aborted)

	// Check the funds are not minted to the contract as the cctx has been aborted
	newBalanceMEVM, err := r.WMuse.BalanceOf(&bind.CallOpts{}, testDappNoRevertAddr)
	require.NoError(r, err)
	require.Equal(r,
		0,
		newBalanceMEVM.Cmp(previousBalanceMEVM),
		"expected new balance to be %s, got %s",
		previousBalanceMEVM.String(),
		newBalanceMEVM.String(),
	)
}
