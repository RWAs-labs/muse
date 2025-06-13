package e2etests

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/common"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestBitcoinDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	depositAmount := utils.ParseFloat(r, args[0])
	// MRC20 BTC amounts have 8 decimals
	depositAmountMRC20 := uint64(depositAmount * btcutil.SatoshiPerBitcoin)

	startingBalance, err := r.BTCMRC20.BalanceOf(&bind.CallOpts{}, r.MEVMAuth.From)
	require.NoError(r, err)

	txHash := r.DepositBTCWithAmount(depositAmount, nil)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// calculate fee
	tx, err := r.BtcRPCClient.GetTransaction(r.Ctx, txHash)
	require.NoError(r, err)
	rawTx, err := r.BtcRPCClient.GetRawTransactionResult(r.Ctx, txHash, tx)
	require.NoError(r, err)
	fee, err := common.CalcDepositorFee(r.Ctx, r.BtcRPCClient, &rawTx, r.BitcoinParams)
	require.NoError(r, err)
	feeSatoshis := uint64(fee * btcutil.SatoshiPerBitcoin)

	expectedAmount := depositAmountMRC20 - feeSatoshis

	// assert that the inbound amount is expected
	require.InDelta(r, expectedAmount, cctx.InboundParams.Amount.Uint64(), 100)

	// assert that the balance increases by the expected amount
	endingBalance, err := r.BTCMRC20.BalanceOf(&bind.CallOpts{}, r.MEVMAuth.From)
	require.NoError(r, err)
	balanceDiff := bigSub(endingBalance, startingBalance)
	require.InDelta(r, expectedAmount, balanceDiff.Uint64(), 100)
}
