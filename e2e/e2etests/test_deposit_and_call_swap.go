package e2etests

import (
	"math/big"
	"time"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TODO: This test is similar to TestCrosschainSwap
// purpose is to test similar scenario with v2 contracts where there is swap + withdraw in onCall
// to showcase that it's not reverting with gas limit issues
// this test should be removed when this issue is completed: https://github.com/RWAs-labs/muse/issues/2711
func TestDepositAndCallSwap(r *runner.E2ERunner, _ []string) {
	// create tokens pair (erc20 and eth)
	tx, err := r.UniswapV2Factory.CreatePair(r.MEVMAuth, r.ERC20MRC20Addr, r.ETHMRC20Addr)
	if err != nil {
		r.Logger.Print("ℹ️ create pair error %s", err.Error())
	} else {
		utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	}

	// approve router to spend tokens being swapped
	tx, err = r.ERC20MRC20.Approve(r.MEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)

	tx, err = r.ETHMRC20.Approve(r.MEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// fund MEVMSwapApp with gas MRC20s for withdraw
	tx, err = r.ETHMRC20.Transfer(r.MEVMAuth, r.MEVMSwapAppAddr, big.NewInt(1e10))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)

	tx, err = r.ERC20MRC20.Transfer(r.MEVMAuth, r.MEVMSwapAppAddr, big.NewInt(1e6))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// temporarily increase gas limit to 400000
	previousGasLimit := r.MEVMAuth.GasLimit
	defer func() {
		r.MEVMAuth.GasLimit = previousGasLimit
	}()

	// add liquidity for swap
	r.MEVMAuth.GasLimit = 400000
	tx, err = r.UniswapV2Router.AddLiquidity(
		r.MEVMAuth,
		r.ERC20MRC20Addr,
		r.ETHMRC20Addr,
		big.NewInt(1e8),
		big.NewInt(1e8),
		big.NewInt(1e8),
		big.NewInt(1e5),
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// memobytes is dApp specific; see the contracts/MEVMSwapApp.sol for details
	// it is [targetMRC20, receiver]
	memobytes, err := r.MEVMSwapApp.EncodeMemo(
		&bind.CallOpts{},
		r.ETHMRC20Addr,
		r.EVMAddress().Bytes(),
	)
	require.NoError(r, err)

	// perform the deposit and call
	r.ApproveERC20OnEVM(r.GatewayEVMAddr)
	tx = r.ERC20DepositAndCall(
		r.MEVMSwapAppAddr,
		big.NewInt(8e7),
		memobytes,
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit_and_call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
}
