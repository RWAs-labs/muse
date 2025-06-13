package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/vault"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

func TestPauseMRC20(r *runner.E2ERunner, _ []string) {
	// Setup vault used to test mrc20 interactions
	r.Logger.Info("Deploying vault")
	vaultAddr, _, vaultContract, err := vault.DeployVault(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)

	// Approving vault to spend MRC20
	tx, err := r.ETHMRC20.Approve(r.MEVMAuth, vaultAddr, big.NewInt(1e18))
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.ERC20MRC20.Approve(r.MEVMAuth, vaultAddr, big.NewInt(1e18))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Pause ETH MRC20
	r.Logger.Info("Pausing ETH")
	msgPause := fungibletypes.NewMsgPauseMRC20(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		[]string{r.ETHMRC20Addr.Hex()},
	)
	res, err := r.MuseTxServer.BroadcastTx(utils.EmergencyPolicyName, msgPause)
	require.NoError(r, err)
	r.Logger.Info("pause mrc20 tx hash: %s", res.TxHash)

	// Fetch and check pause status
	fcRes, err := r.FungibleClient.ForeignCoins(r.Ctx, &fungibletypes.QueryGetForeignCoinsRequest{
		Index: r.ETHMRC20Addr.Hex(),
	})
	require.NoError(r, err)
	require.True(r, fcRes.GetForeignCoins().Paused, "ETH should be paused")

	r.Logger.Info("ETH is paused")

	// Try operations with ETH MRC20
	r.Logger.Info("Can no longer do operations on ETH MRC20")

	tx, err = r.ETHMRC20.Transfer(r.MEVMAuth, sample.EthAddress(), big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)

	tx, err = r.ETHMRC20.Burn(r.MEVMAuth, big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)

	// Operation on a contract that interact with ETH MRC20 should fail
	r.Logger.Info("Vault contract can no longer interact with ETH MRC20: %s", r.ETHMRC20Addr.Hex())
	tx, err = vaultContract.Deposit(r.MEVMAuth, r.ETHMRC20Addr, big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)

	r.Logger.Info("Operations all failed")

	// Check we can still interact with ERC20 MRC20
	r.Logger.Info("Check other MRC20 can still be operated")

	tx, err = r.ERC20MRC20.Transfer(r.MEVMAuth, sample.EthAddress(), big.NewInt(1e3))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tx, err = vaultContract.Deposit(r.MEVMAuth, r.ERC20MRC20Addr, big.NewInt(1e3))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Check deposit revert when paused
	depositHash := r.DepositEtherDeployer()

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, depositHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_Reverted)

	r.Logger.Info("CCTX has been reverted")

	// Unpause ETH MRC20
	r.Logger.Info("Unpausing ETH")
	msgUnpause := fungibletypes.NewMsgUnpauseMRC20(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		[]string{r.ETHMRC20Addr.Hex()},
	)
	res, err = r.MuseTxServer.BroadcastTx(utils.OperationalPolicyName, msgUnpause)
	require.NoError(r, err)

	r.Logger.Info("unpause mrc20 tx hash: %s", res.TxHash)

	// Fetch and check pause status
	fcRes, err = r.FungibleClient.ForeignCoins(r.Ctx, &fungibletypes.QueryGetForeignCoinsRequest{
		Index: r.ETHMRC20Addr.Hex(),
	})
	require.NoError(r, err)
	require.False(r, fcRes.GetForeignCoins().Paused, "ETH should be unpaused")

	r.Logger.Info("ETH is unpaused")

	// Try operations with ETH MRC20
	r.Logger.Info("Can do operations on ETH MRC20 again")

	tx, err = r.ETHMRC20.Transfer(r.MEVMAuth, sample.EthAddress(), big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.ETHMRC20.Burn(r.MEVMAuth, big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Can deposit tokens into the vault again
	tx, err = vaultContract.Deposit(r.MEVMAuth, r.ETHMRC20Addr, big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("Operations all succeeded")
}
