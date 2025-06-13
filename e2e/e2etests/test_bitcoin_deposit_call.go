package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	testcontract "github.com/RWAs-labs/muse/e2e/contracts/example"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/common"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestBitcoinDepositAndCall(r *runner.E2ERunner, args []string) {
	// Given amount to send
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])
	amountTotal := amount + common.DefaultDepositorFee

	// deploy an example contract in MEVM
	contractAddr, _, contract, err := testcontract.DeployExample(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Bitcoin: Example contract deployed at: %s", contractAddr.String())

	// ACT
	// Send BTC to TSS address with a dummy memo
	data := []byte("hello satoshi")
	memo := append(contractAddr.Bytes(), data...)
	txHash, err := r.SendToTSSWithMemo(amountTotal, memo)
	require.NoError(r, err)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check if example contract has been called, 'bar' value should be set to amount
	amountSats, err := common.GetSatoshis(amount)
	require.NoError(r, err)
	utils.MustHaveCalledExampleContract(
		r,
		contract,
		big.NewInt(amountSats),
		[]byte(r.GetBtcAddress().EncodeAddress()),
	)
}
