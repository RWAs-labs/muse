package runner

import (
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/contextapp"
	"github.com/RWAs-labs/muse/e2e/contracts/mevmswap"
	"github.com/RWAs-labs/muse/e2e/contracts/testdapp"
	e2eutils "github.com/RWAs-labs/muse/e2e/utils"
)

// SetupLegacyMEVMContracts sets up the legacy contracts on MEVM
// In particular it deploys test contracts used with the protocol contracts v1
func (r *E2ERunner) SetupLegacyMEVMContracts() {
	// deploy TestDApp contract on mEVM
	appAddr, txApp, _, err := testdapp.DeployTestDApp(
		r.MEVMAuth,
		r.MEVMClient,
		r.ConnectorMEVMAddr,
		r.WMuseAddr,
	)
	require.NoError(r, err)

	r.MevmTestDAppAddr = appAddr
	r.Logger.Info("TestDApp Mevm contract address: %s, tx hash: %s", appAddr.Hex(), txApp.Hash().Hex())

	// deploy MEVMSwapApp and ContextApp
	mevmSwapAppAddr, txMEVMSwapApp, mevmSwapApp, err := mevmswap.DeployMEVMSwapApp(
		r.MEVMAuth,
		r.MEVMClient,
		r.UniswapV2RouterAddr,
		r.SystemContractAddr,
	)
	require.NoError(r, err)

	contextAppAddr, txContextApp, contextApp, err := contextapp.DeployContextApp(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)

	receipt := e2eutils.MustWaitForTxReceipt(
		r.Ctx,
		r.MEVMClient,
		txMEVMSwapApp,
		r.Logger,
		r.ReceiptTimeout,
	)
	r.requireTxSuccessful(receipt, "MEVMSwapApp deployment failed")

	r.MEVMSwapAppAddr = mevmSwapAppAddr
	r.MEVMSwapApp = mevmSwapApp

	receipt = e2eutils.MustWaitForTxReceipt(
		r.Ctx,
		r.MEVMClient,
		txContextApp,
		r.Logger,
		r.ReceiptTimeout,
	)
	r.requireTxSuccessful(receipt, "ContextApp deployment failed")

	r.ContextAppAddr = contextAppAddr
	r.ContextApp = contextApp
}
