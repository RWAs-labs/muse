package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/testabort"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/testutil/sample"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestETHWithdrawRevertAndAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	amount := utils.ParseBigInt(r, args[0])
	gasLimit := utils.ParseBigInt(r, args[1])

	r.ApproveETHMRC20(r.GatewayMEVMAddr)

	// deploy testabort contract
	testAbortAddr, _, testAbort, err := testabort.DeployTestAbort(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)

	// perform the withdraw
	tx := r.ETHWithdrawAndCall(
		sample.EthAddress(), // non-existing address
		amount,
		[]byte("revert"),
		gatewaymevm.RevertOptions{
			RevertAddress: sample.EthAddress(), // non-existing address
			CallOnRevert:  true,
			RevertMessage: []byte(
				"withdraw",
			), // withdraw is passed as message to create a withdraw in onAbort and test cctx can be created
			OnRevertGasLimit: big.NewInt(200000),
			AbortAddress:     testAbortAddr,
		},
		gasLimit,
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	require.Equal(r, crosschaintypes.CctxStatus_Aborted, cctx.CctxStatus.Status)

	// check onAbort was called
	aborted, err := testAbort.IsAborted(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, aborted)

	// check abort context was passed
	abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, "withdraw")
	require.NoError(r, err)
	require.EqualValues(r, r.ETHMRC20Addr.Hex(), abortContext.Asset.Hex())

	// check the create withdraw get mined
	cctxWithdrawFromAbort := utils.WaitCctxMinedByInboundHash(
		r.Ctx,
		cctx.Index,
		r.CctxClient,
		r.Logger,
		r.CctxTimeout,
	)

	// check the cctx status
	utils.RequireCCTXStatus(r, cctxWithdrawFromAbort, crosschaintypes.CctxStatus_OutboundMined)
}
