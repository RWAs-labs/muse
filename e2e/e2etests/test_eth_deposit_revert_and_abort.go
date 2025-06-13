package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/testabort"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/testutil/sample"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestETHDepositRevertAndAbort(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveERC20OnEVM(r.GatewayEVMAddr)

	// deploy testabort contract
	testAbortAddr, _, testAbort, err := testabort.DeployTestAbort(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)

	// perform the deposit
	tx := r.ETHDepositAndCall(
		r.TestDAppV2MEVMAddr,
		amount,
		[]byte("revert"),
		gatewayevm.RevertOptions{
			RevertAddress:    r.TestDAppV2EVMAddr,
			CallOnRevert:     true,
			RevertMessage:    []byte("revert"),
			OnRevertGasLimit: big.NewInt(200000),
			AbortAddress:     testAbortAddr,
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.Equal(r, crosschaintypes.CctxStatus_Aborted, cctx.CctxStatus.Status)

	// check onAbort was called
	aborted, err := testAbort.IsAborted(&bind.CallOpts{})
	require.NoError(r, err)
	require.True(r, aborted)

	// check abort context was passed
	abortContext, err := testAbort.GetAbortedWithMessage(&bind.CallOpts{}, "revert")
	require.NoError(r, err)
	require.EqualValues(r, r.ETHMRC20Addr.Hex(), abortContext.Asset.Hex())

	// check abort contract received the tokens
	balance, err := r.ETHMRC20.BalanceOf(&bind.CallOpts{}, testAbortAddr)
	require.NoError(r, err)
	require.True(r, balance.Uint64() > 0)

	// Test 2: no contract for abort

	// check that funds are still received if onAbort is not called or fails
	eoaAddress := sample.EthAddress()

	tx = r.ETHDepositAndCall(
		r.TestDAppV2MEVMAddr,
		amount,
		[]byte("revert"),
		gatewayevm.RevertOptions{
			RevertAddress:    r.TestDAppV2EVMAddr,
			CallOnRevert:     true,
			RevertMessage:    []byte("revert"),
			OnRevertGasLimit: big.NewInt(200000),
			AbortAddress:     eoaAddress,
		},
	)

	// wait for the cctx to be mined
	cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.Equal(r, crosschaintypes.CctxStatus_Aborted, cctx.CctxStatus.Status)

	// check abort contract received the tokens
	balance, err = r.ETHMRC20.BalanceOf(&bind.CallOpts{}, eoaAddress)
	require.NoError(r, err)
	require.True(r, balance.Uint64() > 0)
}
