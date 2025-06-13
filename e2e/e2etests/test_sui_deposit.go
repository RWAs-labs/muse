package e2etests

import (
	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/coin"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestSuiDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	oldBalance, err := r.SUIMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// make the deposit transaction
	resp := r.SuiDepositSUI(r.EVMAddress(), math.NewUintFromBigInt(amount))

	r.Logger.Info("Sui deposit tx: %s", resp.Digest)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, resp.Digest, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.EqualValues(r, coin.CoinType_Gas, cctx.InboundParams.CoinType)
	require.EqualValues(r, amount.Uint64(), cctx.InboundParams.Amount.Uint64())

	newBalance, err := r.SUIMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	require.EqualValues(r, oldBalance.Add(oldBalance, amount).Uint64(), newBalance.Uint64())
}
