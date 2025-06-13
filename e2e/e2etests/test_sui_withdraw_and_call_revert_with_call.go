package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/testutil/sample"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestSuiWithdrawAndCallRevertWithCall executes withdrawAndCall on mevm gateway with SUI token.
// The outbound is rejected by the connected module due to invalid payload (invalid address),
// and the 'onRevert' method is called in the MEVM to handle the revert.
func TestSuiWithdrawAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// ARRANGE
	// Given target package ID (example package) and a SUI amount
	targetPackageID := r.SuiExample.PackageID.String()
	amount := utils.ParseBigInt(r, args[0])

	// create the payload for 'on_call' with invalid address
	// taking the first 10 letters to form an invalid payload
	invalidAddress := sample.SuiAddress(r)[:10]
	invalidPayloadOnCall, err := r.SuiCreateExampleWACPayload(invalidAddress)
	require.NoError(r, err)

	// given MEVM revert address (the dApp)
	dAppAddress := r.TestDAppV2MEVMAddr
	dAppBalanceBefore, err := r.SUIMRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)

	// given random payload for 'onRevert'
	payloadOnRevert := randomPayload(r)
	r.AssertTestDAppEVMCalled(false, payloadOnRevert, amount)

	// ACT
	// approve SUI MRC20 token
	r.ApproveSUIMRC20(r.GatewayMEVMAddr)

	// perform the withdraw and call with revert options
	tx := r.SuiWithdrawAndCallSUI(
		targetPackageID,
		amount,
		invalidPayloadOnCall,
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
	dAppBalanceAfter, err := r.SUIMRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)
	require.Equal(r, amount.Int64(), dAppBalanceAfter.Int64()-dAppBalanceBefore.Int64())
}
