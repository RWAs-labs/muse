package e2etests

import (
	"math/big"
	"time"

	"cosmossdk.io/math"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestSuiWithdrawRevertWithCall executes withdraw on mevm gateway.
// The outbound is rejected by Sui network, and 'nonce_increase' is called instead to cancel the tx.
func TestSuiWithdrawRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)
	amount := utils.ParseBigInt(r, args[0])

	// ARRANGE
	// given signer
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	signerBalanceBefore := r.SuiGetSUIBalance(signer.Address())

	// given MEVM revert address (the dApp)
	dAppAddress := r.TestDAppV2MEVMAddr
	dAppBalanceBefore, err := r.SUIMRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)

	// given random payload
	payload := randomPayload(r)
	r.AssertTestDAppEVMCalled(false, payload, amount)

	// retrieve current mrc20 gas limit
	oldGasLimit, err := r.SUIMRC20.GASLIMIT(&bind.CallOpts{})
	require.NoError(r, err)
	r.Logger.Info("current gas limit: %s", oldGasLimit.String())

	// set a low MRC20 gas limit so gasBudget will be low: "1000000"
	// withdraw tx will be rejected due to execution error "InsufficientGas"
	lowGasLimit := math.NewUintFromBigInt(big.NewInt(1000))
	_, err = r.MuseTxServer.UpdateMRC20GasLimit(r.SUIMRC20Addr, lowGasLimit)
	require.NoError(r, err)

	// wait for the new gas limit to take effect
	utils.WaitForMuseBlocks(r.Ctx, r, r.MEVMClient, 1, 10*time.Second)

	// ACT
	// approve the MRC20
	r.ApproveSUIMRC20(r.GatewayMEVMAddr)

	// perform the withdraw with revert options
	tx := r.SuiWithdrawSUI(
		signer.Address(),
		amount,
		gatewaymevm.RevertOptions{
			CallOnRevert:     true,
			RevertAddress:    dAppAddress,
			RevertMessage:    []byte(payload),
			OnRevertGasLimit: big.NewInt(0),
		},
	)
	r.Logger.EVMTransaction(*tx, "withdraw")

	// ASSERT
	// wait for the CCTX to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// should have called 'onRevert'
	r.AssertTestDAppMEVMCalled(true, payload, big.NewInt(0))

	// sender and message should match
	sender, err := r.TestDAppV2MEVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payload),
	)
	require.NoError(r, err)
	require.Equal(r, r.MEVMAuth.From, sender)

	// signer balance should remain unchanged in Sui chain
	signerBalanceAfter := r.SuiGetSUIBalance(signer.Address())
	require.Equal(r, signerBalanceBefore, signerBalanceAfter)

	// the dApp address should get reverted amount
	dAppBalanceAfter, err := r.SUIMRC20.BalanceOf(&bind.CallOpts{}, dAppAddress)
	require.NoError(r, err)
	require.Equal(r, amount.Int64(), dAppBalanceAfter.Int64()-dAppBalanceBefore.Int64())

	// TEARDOWN
	// restore old gas limit
	_, err = r.MuseTxServer.UpdateMRC20GasLimit(r.SUIMRC20Addr, math.NewUintFromBigInt(oldGasLimit))
	require.NoError(r, err, "failed to restore gas limit")
}
