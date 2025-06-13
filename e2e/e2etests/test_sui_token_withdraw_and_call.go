package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestSuiTokenWithdrawAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// Given target package ID (example package) and a token amount
	targetPackageID := r.SuiExample.PackageID.String()
	amount := utils.ParseBigInt(r, args[0])

	// use the deployer address as on_call payload message
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	suiAddress := signer.Address()

	// Given initial balance and called_count
	balanceBefore := r.SuiGetFungibleTokenBalance(suiAddress)
	calledCountBefore := r.SuiGetConnectedCalledCount()

	// create the on_call payload
	payloadOnCall, err := r.SuiCreateExampleWACPayload(suiAddress)
	require.NoError(r, err)

	// ACT
	// approve both SUI gas budget token and fungible token MRC20
	r.ApproveSUIMRC20(r.GatewayMEVMAddr)
	r.ApproveFungibleTokenMRC20(r.GatewayMEVMAddr)

	// perform the fungible token withdraw and call
	tx := r.SuiWithdrawAndCallFungibleToken(
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

	// check the balance after the withdraw
	balanceAfter := r.SuiGetFungibleTokenBalance(signer.Address())
	require.EqualValues(r, balanceBefore+amount.Uint64(), balanceAfter)

	// verify the called_count increased by 1
	calledCountAfter := r.SuiGetConnectedCalledCount()
	require.Equal(r, calledCountBefore+1, calledCountAfter)
}
