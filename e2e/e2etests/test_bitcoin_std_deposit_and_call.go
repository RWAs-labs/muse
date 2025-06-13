package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	testcontract "github.com/RWAs-labs/muse/e2e/contracts/example"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	musebitcoin "github.com/RWAs-labs/muse/museclient/chains/bitcoin/common"
	"github.com/RWAs-labs/muse/pkg/memo"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestBitcoinStdMemoDepositAndCall(r *runner.E2ERunner, args []string) {
	// parse amount to deposit
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])

	// deploy an example contract in MEVM
	contractAddr, _, contract, err := testcontract.DeployExample(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)

	// create standard memo with [receiver, payload]
	memo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeDepositAndCall,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: contractAddr,
			Payload:  []byte("hello satoshi"),
		},
	}

	// deposit BTC with standard memo
	txHash := r.DepositBTCWithExactAmount(amount, memo)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_std_memo_deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check if example contract has been called, 'bar' value should be set to amount
	amountSats, err := musebitcoin.GetSatoshis(amount)
	require.NoError(r, err)
	utils.MustHaveCalledExampleContract(
		r,
		contract,
		big.NewInt(amountSats),
		[]byte(r.GetBtcAddress().EncodeAddress()),
	)
}
