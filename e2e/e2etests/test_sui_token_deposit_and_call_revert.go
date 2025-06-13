package e2etests

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/coin"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestSuiTokenDepositAndCallRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "get deployer signer")
	balanceBefore := r.SuiGetFungibleTokenBalance(signer.Address())
	tssBalanceBefore := r.SuiGetSUIBalance(r.SuiTSSAddress)

	// add liquidity in pool to allow revert fee to be paid
	museAmount := big.NewInt(1e18)
	mrc20Amount := big.NewInt(9000000000)
	r.AddLiquiditySUI(museAmount, mrc20Amount)
	r.AddLiquiditySuiFungibleToken(museAmount, mrc20Amount)

	// make the deposit transaction
	resp := r.SuiFungibleTokenDepositAndCall(r.TestDAppV2MEVMAddr, math.NewUintFromBigInt(amount), []byte("revert"))

	r.Logger.Info("Sui deposit and call tx: %s", resp.Digest)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, resp.Digest, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.EqualValues(r, crosschaintypes.CctxStatus_Reverted, cctx.CctxStatus.Status)
	require.EqualValues(r, coin.CoinType_ERC20, cctx.InboundParams.CoinType)
	require.EqualValues(r, amount.Uint64(), cctx.InboundParams.Amount.Uint64())

	// check the balance after the failed deposit is higher than balance before - amount
	// reason it's not equal is because of the gas fee for revert
	balanceAfter := r.SuiGetFungibleTokenBalance(signer.Address())
	require.Greater(r, balanceAfter, balanceBefore-amount.Uint64())

	// check the TSS balance after transaction is higher or equal to the balance before
	// reason is that the max budget is refunded to the TSS
	tssBalanceAfter := r.SuiGetSUIBalance(r.SuiTSSAddress)
	require.GreaterOrEqual(r, tssBalanceAfter, tssBalanceBefore)
}
