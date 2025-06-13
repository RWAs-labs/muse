package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/testconnectormevm"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

// TestUpdateBytecodeConnector tests updating the bytecode of a connector and interact with it
func TestUpdateBytecodeConnector(r *runner.E2ERunner, _ []string) {
	// Can withdraw 10MUSE
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(10))
	r.LegacyDepositAndApproveWMuse(amount)

	tx := r.LegacyWithdrawMuse(amount, true)
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	r.Logger.CCTX(*cctx, "muse withdraw")

	// Deploy the test contract
	newTestConnectorAddr, tx, _, err := testconnectormevm.DeployTestMuseConnectorMEVM(
		r.MEVMAuth,
		r.MEVMClient,
		r.WMuseAddr,
	)
	require.NoError(r, err)

	// Wait for the contract to be deployed
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Get the code hash of the new contract
	codeHashRes, err := r.FungibleClient.CodeHash(r.Ctx, &fungibletypes.QueryCodeHashRequest{
		Address: newTestConnectorAddr.String(),
	})
	require.NoError(r, err)
	r.Logger.Info("New contract code hash: %s", codeHashRes.CodeHash)

	r.Logger.Info("Updating the bytecode of the Connector")
	msg := fungibletypes.NewMsgUpdateContractBytecode(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		r.ConnectorMEVMAddr.Hex(),
		codeHashRes.CodeHash,
	)
	res, err := r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)
	r.Logger.Info("Update connector bytecode tx hash: %s", res.TxHash)

	r.Logger.Info("Can interact with the new code of the contract")
	testConnectorContract, err := testconnectormevm.NewTestMuseConnectorMEVM(r.ConnectorMEVMAddr, r.MEVMClient)
	require.NoError(r, err)

	response, err := testConnectorContract.Foo(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, "foo", response)

	// Can continue to interact with the connector: withdraw 10MUSE
	r.LegacyDepositAndApproveWMuse(amount)
	tx = r.LegacyWithdrawMuse(amount, true)
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "muse withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
}
