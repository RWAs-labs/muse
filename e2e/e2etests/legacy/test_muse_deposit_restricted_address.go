package legacy

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/testutil/sample"
)

func TestMuseDepositRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the deposit amount
	amount := utils.ParseBigInt(r, args[0])

	// Deposit amount to restricted address
	txHash := r.LegacyDepositMuseWithAmount(ethcommon.HexToAddress(sample.RestrictedEVMAddressTest), amount)

	// wait for 5 muse blocks
	r.WaitForBlocks(5)

	// no cctx should be created
	utils.EnsureNoCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient)
}
