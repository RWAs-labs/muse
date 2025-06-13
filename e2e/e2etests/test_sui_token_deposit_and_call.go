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

func TestSuiTokenDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	oldBalance, err := r.SuiTokenMRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2MEVMAddr)
	require.NoError(r, err)

	payload := randomPayload(r)

	// make the deposit transaction
	resp := r.SuiFungibleTokenDepositAndCall(r.TestDAppV2MEVMAddr, math.NewUintFromBigInt(amount), []byte(payload))

	r.Logger.Info("Sui deposit and call tx: %s", resp.Digest)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, resp.Digest, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.EqualValues(r, coin.CoinType_ERC20, cctx.InboundParams.CoinType)
	require.EqualValues(r, amount.Uint64(), cctx.InboundParams.Amount.Uint64())

	// check the payload was received on the contract
	r.AssertTestDAppMEVMCalled(true, payload, amount)

	newBalance, err := r.SuiTokenMRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2MEVMAddr)
	require.NoError(r, err)
	require.EqualValues(r, oldBalance.Add(oldBalance, amount).Uint64(), newBalance.Uint64())
}
