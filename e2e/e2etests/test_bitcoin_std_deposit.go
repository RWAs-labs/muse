package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/common"
	"github.com/RWAs-labs/muse/pkg/memo"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestBitcoinStdMemoDeposit(r *runner.E2ERunner, args []string) {
	// parse amount to deposit
	require.Len(r, args, 1)
	amount := utils.ParseFloat(r, args[0])

	// get ERC20 BTC balance before deposit
	balanceBefore, err := r.BTCMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of BTC before deposit: %d satoshis", balanceBefore)

	// create standard memo with receiver address
	memo := &memo.InboundMemo{
		Header: memo.Header{
			Version:     0,
			EncodingFmt: memo.EncodingFmtCompactShort,
			OpCode:      memo.OpCodeDeposit,
		},
		FieldsV0: memo.FieldsV0{
			Receiver: r.EVMAddress(), // to deployer self
		},
	}

	// deposit BTC with standard memo
	txHash := r.DepositBTCWithExactAmount(amount, memo)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "bitcoin_std_memo_deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// get ERC20 BTC balance after deposit
	balanceAfter, err := r.BTCMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of BTC after deposit: %d satoshis", balanceAfter)

	// the runner balance should be increased by the deposit amount
	amountIncreased := new(big.Int).Sub(balanceAfter, balanceBefore)
	amountSatoshis, err := common.GetSatoshis(amount)
	require.NoError(r, err)
	require.Positive(r, amountSatoshis)
	// #nosec G115 always positive
	require.Equal(r, uint64(amountSatoshis), amountIncreased.Uint64())
}
