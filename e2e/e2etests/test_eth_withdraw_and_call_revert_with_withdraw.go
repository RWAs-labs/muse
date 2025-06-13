package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestETHWithdrawAndCallRevertWithWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveETHMRC20(r.GatewayMEVMAddr)

	// perform the withdraw
	tx := r.ETHWithdrawAndArbitraryCall(
		r.TestDAppV2EVMAddr,
		amount,
		r.EncodeGasCall("revert"),
		gatewaymevm.RevertOptions{
			RevertAddress:    r.TestDAppV2MEVMAddr,
			CallOnRevert:     true,
			RevertMessage:    []byte("withdraw"), // call withdraw in the onRevert hook
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the cctx to be mined
	cctxRevert := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctxRevert, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_Reverted, cctxRevert.CctxStatus.Status)

	cctxWithdrawFromRevert := utils.WaitCctxMinedByInboundHash(
		r.Ctx,
		cctxRevert.Index,
		r.CctxClient,
		r.Logger,
		r.CctxTimeout,
	)

	// check the cctx status
	utils.RequireCCTXStatus(r, cctxWithdrawFromRevert, crosschaintypes.CctxStatus_OutboundMined)
}
