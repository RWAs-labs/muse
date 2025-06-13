package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestSolanaDepositThroughProgram tests triggering gateway deposit through another solana program
// it is same as TestSolanaDeposit, but instead inbound is inside inner instructions
func TestSolanaDepositThroughProgram(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// get ERC20 SOL balance before deposit
	balanceBefore, err := r.SOLMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL before deposit: %d", balanceBefore)

	// parse deposit amount (in lamports)
	depositAmount := utils.ParseBigInt(r, args[0])

	// execute the deposit transaction through connected program
	sig := r.SOLDepositAndCallThroughProgram(nil, r.EVMAddress(), depositAmount, nil)

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
