package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/testmrc20"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/testutil/sample"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

// TestUpdateBytecodeMRC20 tests updating the bytecode of a mrc20 and interact with it
func TestUpdateBytecodeMRC20(r *runner.E2ERunner, _ []string) {
	// Random approval
	approved := sample.EthAddress()
	tx, err := r.ETHMRC20.Approve(r.MEVMAuth, approved, big.NewInt(1e10))
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Deploy the TestMRC20 contract
	r.Logger.Info("Deploying contract with new bytecode")
	newMRC20Address, tx, newMRC20Contract, err := testmrc20.DeployTestMRC20(
		r.MEVMAuth,
		r.MEVMClient,
		big.NewInt(5),
		// #nosec G115 test - always in range
		uint8(coin.CoinType_Gas),
	)
	require.NoError(r, err)

	// Wait for the contract to be deployed
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Get the code hash of the new contract
	codeHashRes, err := r.FungibleClient.CodeHash(r.Ctx, &fungibletypes.QueryCodeHashRequest{
		Address: newMRC20Address.String(),
	})
	require.NoError(r, err)

	r.Logger.Info("New contract code hash: %s", codeHashRes.CodeHash)

	// Get current info of the MRC20
	name, err := r.ETHMRC20.Name(&bind.CallOpts{})
	require.NoError(r, err)

	symbol, err := r.ETHMRC20.Symbol(&bind.CallOpts{})
	require.NoError(r, err)

	decimals, err := r.ETHMRC20.Decimals(&bind.CallOpts{})
	require.NoError(r, err)

	totalSupply, err := r.ETHMRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)

	balance, err := r.ETHMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	approval, err := r.ETHMRC20.Allowance(&bind.CallOpts{}, r.EVMAddress(), approved)
	require.NoError(r, err)

	r.Logger.Info("Updating the bytecode of the MRC20")
	msg := fungibletypes.NewMsgUpdateContractBytecode(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		r.ETHMRC20Addr.Hex(),
		codeHashRes.CodeHash,
	)
	res, err := r.MuseTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	r.Logger.Info("Update mrc20 bytecode tx hash: %s", res.TxHash)

	// Get new info of the MRC20
	r.Logger.Info("Checking the state of the MRC20 remains the same")
	newName, err := r.ETHMRC20.Name(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, name, newName)

	newSymbol, err := r.ETHMRC20.Symbol(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, symbol, newSymbol)

	newDecimals, err := r.ETHMRC20.Decimals(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, decimals, newDecimals)

	newTotalSupply, err := r.ETHMRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, 0, totalSupply.Cmp(newTotalSupply))

	newBalance, err := r.ETHMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	require.Equal(r, 0, balance.Cmp(newBalance))

	newApproval, err := r.ETHMRC20.Allowance(&bind.CallOpts{}, r.EVMAddress(), approved)
	require.NoError(r, err)
	require.Equal(r, 0, approval.Cmp(newApproval))

	r.Logger.Info("Can interact with the new code of the contract")

	testMRC20Contract, err := testmrc20.NewTestMRC20(r.ETHMRC20Addr, r.MEVMClient)
	require.NoError(r, err)

	tx, err = testMRC20Contract.UpdateNewField(r.MEVMAuth, big.NewInt(1e10))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	newField, err := testMRC20Contract.NewField(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, 0, newField.Cmp(big.NewInt(1e10)))

	r.Logger.Info("Interacting with the bytecode contract doesn't disrupt the mrc20 contract")
	tx, err = newMRC20Contract.UpdateNewField(r.MEVMAuth, big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	newField, err = newMRC20Contract.NewField(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, 0, newField.Cmp(big.NewInt(1e5)), "new field value mismatch on bytecode contract")

	newField, err = testMRC20Contract.NewField(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, 0, newField.Cmp(big.NewInt(1e10)), "new field value mismatch on mrc20 contract")

	// can continue to operate the MRC20
	r.Logger.Info("Checking the MRC20 can continue to operate after state change")
	tx, err = r.ETHMRC20.Transfer(r.MEVMAuth, approved, big.NewInt(1e14))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	newBalance, err = r.ETHMRC20.BalanceOf(&bind.CallOpts{}, approved)
	require.NoError(r, err)
	require.Equal(r, 0, newBalance.Cmp(big.NewInt(1e14)))
}
