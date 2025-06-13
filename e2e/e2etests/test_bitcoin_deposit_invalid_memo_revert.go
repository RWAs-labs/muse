package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestBitcoinDepositInvalidMemoRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// mine blocks at normal speed
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// CASE 1
	// make a deposit without memo output
	txHash, err := r.SendToTSSWithMemo(0.1, nil)
	require.NoError(r, err)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit without memo")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)

	// CASE 2
	// make a deposit with a empty memo
	txHash, err = r.SendToTSSWithMemo(0.1, []byte{})
	require.NoError(r, err)

	// wait for the cctx to be mined
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit empty memo")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)

	// CASE 3
	// make a deposit with an invalid memo
	txHash, err = r.SendToTSSWithMemo(0.1, []byte("invalid memo"))
	require.NoError(r, err)

	// wait for the cctx to be mined
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit invalid memo")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.EqualValues(r, crosschaintypes.InboundStatus_INVALID_MEMO, cctx.InboundParams.Status)
}
