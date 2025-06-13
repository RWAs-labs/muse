package e2etests

import (
	"github.com/stretchr/testify/require"

	testcontract "github.com/RWAs-labs/muse/e2e/contracts/reverter"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	toncontracts "github.com/RWAs-labs/muse/pkg/contracts/ton"
	cctypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestTONDepositAndCallRefund(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// Given amount and arbitrary call data
	var (
		amount = utils.ParseUint(r, args[0])
		data   = []byte("hello reverter")
	)

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given deployer mock revert contract
	// deploy a reverter contract in MEVM
	reverterAddr, _, _, err := testcontract.DeployReverter(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Reverter contract deployed at: %s", reverterAddr.String())

	// Given a sender
	_, sender, err := r.Account.AsTONWallet(r.Clients.TON)
	require.NoError(r, err)

	// ACT
	// Send a deposit and call transaction from the deployer (faucet)
	// to the reverter contract
	cctx, err := r.TONDepositAndCall(
		gw,
		sender,
		amount,
		reverterAddr,
		data,
		runner.TONExpectStatus(cctypes.CctxStatus_Reverted),
	)

	// ASSERT
	require.NoError(r, err)
	r.Logger.CCTX(*cctx, "ton_deposit_and_refund")

	require.Contains(r, cctx.CctxStatus.ErrorMessage, utils.ErrHashRevertFoo)
}
