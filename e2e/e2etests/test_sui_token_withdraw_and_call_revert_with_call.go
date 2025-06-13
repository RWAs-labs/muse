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

// TestSuiTokenWithdrawAndCallRevertWithCall executes withdrawAndCall on mevm gateway with fungible token.
// The outbound is rejected by the connected module due to the special payload message "revert" and the
// 'onRevert' method is called in the MEVM to handle the revert.
func TestSuiTokenWithdrawAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// Given target package ID (example package) and a token amount
	targetPackageID := r.SuiExample.PackageID.String()
	amount := utils.ParseBigInt(r, args[0])

	// create the special revert payload for 'on_call'
	revertPayloadOnCall, err := r.SuiCreateExampleWACPayloadForRevert()
	require.NoError(r, err)

	// given MEVM revert address (the dApp)
	dAppAddress := r.TestDAppV2MEVMAddr
	dAppBalanceBefore, err := r.SuiTokenMRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)

	// given random payload for 'onRevert'
	payloadOnRevert := randomPayload(r)
	r.AssertTestDAppEVMCalled(false, payloadOnRevert, amount)

	// ACT
	// approve both SUI gas budget token and fungible token MRC20
	r.ApproveSUIMRC20(r.GatewayMEVMAddr)
	r.ApproveFungibleTokenMRC20(r.GatewayMEVMAddr)

	// perform the withdraw and call with revert options
	tx := r.SuiWithdrawAndCallFungibleToken(
		targetPackageID,
		amount,
		revertPayloadOnCall,
		gatewaymevm.RevertOptions{
			CallOnRevert:     true,
			RevertAddress:    dAppAddress,
			RevertMessage:    []byte(payloadOnRevert),
			OnRevertGasLimit: big.NewInt(0),
		},
	)
	r.Logger.EVMTransaction(*tx, "withdraw_and_call")

	// ASSERT
	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// should have called 'onRevert'
	r.AssertTestDAppMEVMCalled(true, payloadOnRevert, big.NewInt(0))

	// sender and message should match
	sender, err := r.TestDAppV2MEVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payloadOnRevert),
	)
	require.NoError(r, err)
	require.Equal(r, r.MEVMAuth.From, sender)

	// the dApp address should get reverted amount
	dAppBalanceAfter, err := r.SuiTokenMRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)
	require.Equal(r, amount.Int64(), dAppBalanceAfter.Int64()-dAppBalanceBefore.Int64())
}
