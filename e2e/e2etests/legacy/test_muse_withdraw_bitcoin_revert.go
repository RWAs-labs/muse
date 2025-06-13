package legacy

import (
	"math/big"

	connectormevm "github.com/RWAs-labs/protocol-contracts/pkg/museconnectormevm.sol"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
)

func TestMuseWithdrawBTCRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse withdraw amount
	amount := utils.ParseBigInt(r, args[0])

	r.MEVMAuth.Value = amount
	tx, err := r.WMuse.Deposit(r.MEVMAuth)
	require.NoError(r, err)

	r.MEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("Deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "Deposit")
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.WMuse.Approve(r.MEVMAuth, r.ConnectorMEVMAddr, big.NewInt(1e18))
	require.NoError(r, err)

	r.Logger.Info("wmuse.approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.EVMReceipt(*receipt, "Approve")

	lessThanAmount := amount.Div(amount, big.NewInt(10)) // 1/10 of amount
	tx, err = r.ConnectorMEVM.Send(r.MEVMAuth, connectormevm.MuseInterfacesSendInput{
		DestinationChainId:  big.NewInt(chains.BitcoinRegtest.ChainId),
		DestinationAddress:  r.EVMAddress().Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		MuseValueAndGas:     lessThanAmount,
		MuseParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)

	r.Logger.EVMReceipt(*receipt, "send")
}
