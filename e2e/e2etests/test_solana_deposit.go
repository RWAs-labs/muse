package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestSolanaDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// get ERC20 SOL balance before deposit
	balanceBefore, err := r.SOLMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL before deposit: %d", balanceBefore)

	// parse deposit amount (in lamports)
	depositAmount := utils.ParseBigInt(r, args[0])

	// execute the deposit transaction
	sig := r.SOLDepositAndCall(nil, r.EVMAddress(), depositAmount, nil, nil)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, r.EVMAddress().Hex())

	// get ERC20 SOL balance after deposit
	balanceAfter, err := r.SOLMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL after deposit: %d", balanceAfter)

	// the runner balance should be increased by the deposit amount
	amountIncreased := new(big.Int).Sub(balanceAfter, balanceBefore)
	require.Equal(r, depositAmount.String(), amountIncreased.String())
}
