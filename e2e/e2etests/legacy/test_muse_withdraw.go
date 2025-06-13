package legacy

import (
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestMuseWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse withdraw amount
	amount := utils.ParseBigInt(r, args[0])

	r.LegacyDepositAndApproveWMuse(amount)
	tx := r.LegacyWithdrawMuse(amount, true)

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "muse withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
}
