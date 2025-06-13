package e2etests

import (
	"time"

	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	musebitcoin "github.com/RWAs-labs/muse/museclient/chains/bitcoin/common"
	"github.com/RWAs-labs/muse/pkg/constant"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestBitcoinDonation(r *runner.E2ERunner, args []string) {
	// Given amount to send
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])
	amountTotal := amount + musebitcoin.DefaultDepositorFee

	// ACT
	// Send BTC to TSS address with donation message
	memo := []byte(constant.DonationMessage)
	txHash, err := r.SendToTSSWithMemo(amountTotal, memo)
	require.NoError(r, err)

	// ASSERT after 6 Muse blocks
	time.Sleep(constant.MuseBlockTime * 6)
	req := &crosschaintypes.QueryInboundHashToCctxDataRequest{InboundHash: txHash.String()}
	_, err = r.CctxClient.InTxHashToCctxData(r.Ctx, req)
	require.Error(r, err)
}
