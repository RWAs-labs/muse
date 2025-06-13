package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	musebitcoin "github.com/RWAs-labs/muse/museclient/chains/bitcoin/common"
	"github.com/RWAs-labs/muse/testutil/sample"
)

func TestBitcoinDepositAndCallRevert(r *runner.E2ERunner, args []string) {
	// Given "Live" BTC network
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// Given amount to send
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])
	amount += musebitcoin.DefaultDepositorFee

	// ACT
	// Send BTC to TSS address with a dummy memo
	// musecore should revert cctx if call is made on a non-existing address
	nonExistReceiver := sample.EthAddress()
	badMemo := append(nonExistReceiver.Bytes(), []byte("gibberish-memo")...)
	txHash, err := r.SendToTSSWithMemo(amount, badMemo)
	require.NoError(r, err)
	require.NotEmpty(r, txHash)

	// ASSERT
	// Now we want to make sure refund TX is completed.
	cctx := utils.WaitCctxRevertedByInboundHash(r.Ctx, r, txHash.String(), r.CctxClient)

	// Check revert tx receiver address and amount
	receiver, value := r.QueryOutboundReceiverAndAmount(cctx.OutboundParams[1].Hash)
	require.Equal(r, r.GetBtcAddress().EncodeAddress(), receiver)
	require.True(r, value > 0)

	r.Logger.Info("Sent %f BTC to TSS with invalid memo, got refund of %d satoshis", amount, value)
}
