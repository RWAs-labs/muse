package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestSuiWithdrawAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// Given target package ID (example package) and a SUI amount
	targetPackageID := r.SuiExample.PackageID.String()
	amount := utils.ParseBigInt(r, args[0])

	// use the deployer address as on_call payload message
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	suiAddress := signer.Address()

	// Given initial balance and called_count
	balanceBefore := r.SuiGetSUIBalance(suiAddress)
	calledCountBefore := r.SuiGetConnectedCalledCount()

	// create the on_call payload
	payloadOnCall, err := r.SuiCreateExampleWACPayload(suiAddress)
	require.NoError(r, err)

	// ACT
	// approve SUI MRC20 token
	r.ApproveSUIMRC20(r.GatewayMEVMAddr)

	// perform the withdraw and call
	tx := r.SuiWithdrawAndCallSUI(
		targetPackageID,
		amount,
		payloadOnCall,
		gatewaymevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)
	r.Logger.EVMTransaction(*tx, "withdraw_and_call")

	// ASSERT
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// balance after
	balanceAfter := r.SuiGetSUIBalance(suiAddress)
	require.Equal(r, balanceBefore+amount.Uint64(), balanceAfter)

	// verify the called_count increased by 1
	calledCountAfter := r.SuiGetConnectedCalledCount()
	require.Equal(r, calledCountBefore+1, calledCountAfter)
}
