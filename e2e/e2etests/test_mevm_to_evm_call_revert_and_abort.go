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

func TestMEVMToEVMCallRevertAndAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	r.ApproveETHMRC20(r.GatewayMEVMAddr)

	// deploy testabort contract
	testAbortAddr, _, testAbort, err := testabort.DeployTestAbort(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)

	// perform the withdraw
	tx := r.MEVMToEMVCall(
		sample.EthAddress(), // non-existing address
		[]byte("revert"),
		gatewaymevm.RevertOptions{
			RevertAddress:    sample.EthAddress(), // non-existing address
			CallOnRevert:     true,
			RevertMessage:    []byte("revert"),
			OnRevertGasLimit: big.NewInt(200000),
			AbortAddress:     testAbortAddr,
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "call")
	require.Equal(r, crosschaintypes.CctxStatus_Aborted, cctx.CctxStatus.Status)

	// check onAbort was called
	aborted, err := testAbort.IsAborted(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, aborted)

	// check abort context was passed
	abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, "revert")
	require.NoError(r, err)
	require.EqualValues(r, r.ETHMRC20Addr.Hex(), abortContext.Asset.Hex())
	require.EqualValues(r, int64(0), abortContext.Amount.Int64())
}
