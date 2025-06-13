package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestETHWithdrawAndCallNoMessage(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	previousGasLimit := r.MEVMAuth.GasLimit
	r.MEVMAuth.GasLimit = 10000000
	defer func() {
		r.MEVMAuth.GasLimit = previousGasLimit
	}()

	amount := utils.ParseBigInt(r, args[0])
	gasLimit := utils.ParseBigInt(r, args[1])

	r.ApproveETHMRC20(r.GatewayMEVMAddr)

	// perform the withdraw
	tx := r.ETHWithdrawAndCall(
		r.TestDAppV2EVMAddr,
		amount,
		[]byte{},
		gatewaymevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		gasLimit,
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw_and_call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check called
	messageIndex, err := r.TestDAppV2EVM.GetNoMessageIndex(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.AssertTestDAppEVMCalled(true, messageIndex, amount)

	// check expected sender was used
	senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(messageIndex),
	)
	require.NoError(r, err)
	require.Equal(r, r.MEVMAuth.From, senderForMsg)
}
