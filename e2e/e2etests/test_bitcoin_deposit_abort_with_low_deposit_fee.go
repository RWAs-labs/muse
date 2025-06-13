package e2etests

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	musebitcoin "github.com/RWAs-labs/muse/museclient/chains/bitcoin/common"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestBitcoinDepositAndAbortWithLowDepositFee(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0)

	// ARRANGE
	// Given small amount
	depositAmount := musebitcoin.DefaultDepositorFee - float64(1)/btcutil.SatoshiPerBitcoin

	// ACT
	txHash := r.DepositBTCWithAmount(depositAmount, nil)

	// ASSERT
	// cctx status should be aborted
	cctx := utils.WaitCctxAbortedByInboundHash(r.Ctx, r, txHash.String(), r.CctxClient)
	r.Logger.CCTX(cctx, "deposit aborted")

	// check cctx details
	require.Equal(r, cctx.InboundParams.Amount.Uint64(), uint64(0))
	require.Equal(r, cctx.GetCurrentOutboundParam().Amount.Uint64(), uint64(0))

	// check cctx error
	require.EqualValues(r, crosschaintypes.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE, cctx.InboundParams.Status)
}
