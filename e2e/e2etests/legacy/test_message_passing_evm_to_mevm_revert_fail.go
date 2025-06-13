package legacy

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/testdappnorevert"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	cctxtypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestMessagePassingEVMtoMEVMRevertFail(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := utils.ParseBigInt(r, args[0])

	// Deploying a test contract not containing a logic for reverting the cctx
	testDappNoRevertEVMAddr, tx, testDappNoRevertEVM, err := testdappnorevert.DeployTestDAppNoRevert(
		r.EVMAuth,
		r.EVMClient,
		r.ConnectorEthAddr,
		r.MuseEthAddr,
	)
	require.NoError(r, err)

	r.Logger.Info("TestDAppNoRevertEVM deployed at: %s", testDappNoRevertEVMAddr.Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.EVMReceipt(*receipt, "deploy TestDAppNoRevert")

	// Set destination details
	mEVMChainID, err := r.MEVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	destinationAddress := r.MevmTestDAppAddr

	// Contract call originates from EVM chain
	tx, err = r.MuseEth.Approve(r.EVMAuth, testDappNoRevertEVMAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("Approve tx receipt: %d", receipt.Status)

	// Get MUSE balance before test
	previousBalanceMEVM, err := r.WMuse.BalanceOf(&bind.CallOpts{}, r.MevmTestDAppAddr)
	require.NoError(r, err)

	previousBalanceEVM, err := r.MuseEth.BalanceOf(&bind.CallOpts{}, testDappNoRevertEVMAddr)
	require.NoError(r, err)

	// Send message with doRevert
	tx, err = testDappNoRevertEVM.SendHelloWorld(r.EVMAuth, destinationAddress, mEVMChainID, amount, true)
	require.NoError(r, err)

	r.Logger.Info("TestDAppNoRevert.SendHello tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// New inbound message picked up by muse-clients and voted on by observers to initiate a contract call on mEVM which would revert the transaction
	// A revert transaction is created and gets finalized on the original sender chain.
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_Aborted)

	// Check MUSE balance on MEVM TestDApp and check new balance is previous balance
	newBalanceMEVM, err := r.WMuse.BalanceOf(&bind.CallOpts{}, r.MevmTestDAppAddr)
	require.NoError(r, err)
	require.Equal(
		r,
		0,
		newBalanceMEVM.Cmp(previousBalanceMEVM),
		"expected new balance to be %s, got %s",
		previousBalanceMEVM.String(),
		newBalanceMEVM.String(),
	)

	// Check MUSE balance on EVM TestDApp and check new balance is previous balance
	newBalanceEVM, err := r.MuseEth.BalanceOf(&bind.CallOpts{}, testDappNoRevertEVMAddr)
	require.NoError(r, err)
	require.Equal(r, 0, newBalanceEVM.Cmp(previousBalanceEVM))
}
