package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/gatewaymevmcaller"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestETHWithdrawAndCallThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	previousGasLimit := r.MEVMAuth.GasLimit
	r.MEVMAuth.GasLimit = 10000000
	defer func() {
		r.MEVMAuth.GasLimit = previousGasLimit
	}()

	amount := utils.ParseBigInt(r, args[0])

	// deploy caller contract and send it gas mrc20 to pay gas fee
	gatewayCallerAddr, tx, gatewayCaller, err := gatewaymevmcaller.DeployGatewayMEVMCaller(
		r.MEVMAuth,
		r.MEVMClient,
		r.GatewayMEVMAddr,
		r.WMuseAddr,
	)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)

	tx, err = r.ETHMRC20.Transfer(r.MEVMAuth, gatewayCallerAddr, big.NewInt(100000000000000000))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)

	payload := randomPayload(r)

	// perform the authenticated call
	tx = r.ETHWithdrawAndCallThroughContract(gatewayCaller, r.TestDAppV2EVMAddr,
		amount,
		[]byte(payload),
		gatewaymevmcaller.RevertOptions{OnRevertGasLimit: big.NewInt(0)})

	utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	r.AssertTestDAppEVMCalled(true, payload, amount)

	// check expected sender was used
	senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payload),
	)
	require.NoError(r, err)
	require.Equal(r, gatewayCallerAddr, senderForMsg)
}
