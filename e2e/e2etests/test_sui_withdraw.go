package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestSuiWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")

	balanceBefore := r.SuiGetSUIBalance(signer.Address())
	tssBalanceBefore := r.SuiGetSUIBalance(r.SuiTSSAddress)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveSUIMRC20(r.GatewayMEVMAddr)

	// perform the withdraw
	tx := r.SuiWithdrawSUI(signer.Address(), amount, gatewaymevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})
	r.Logger.EVMTransaction(*tx, "withdraw")

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the balance after the withdraw
	balanceAfter := r.SuiGetSUIBalance(signer.Address())
	require.EqualValues(r, balanceBefore+amount.Uint64(), balanceAfter)

	// check the TSS balance after transaction is higher or equal to the balance before
	// reason is that the max budget is refunded to the TSS
	tssBalanceAfter := r.SuiGetSUIBalance(r.SuiTSSAddress)
	require.GreaterOrEqual(r, tssBalanceAfter, tssBalanceBefore)
}
