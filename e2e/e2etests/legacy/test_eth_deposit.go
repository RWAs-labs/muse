package legacy

import (
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
)

// TestEtherDeposit tests deposit of ethers
func TestEtherDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the deposit amount
	amount := utils.ParseBigInt(r, args[0])

	hash := r.LegacyDepositEtherWithAmount(amount) // in wei
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
}
