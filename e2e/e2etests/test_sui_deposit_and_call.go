package e2etests

import (
	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/contracts/sui"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestSuiDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	oldBalance, err := r.SUIMRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2MEVMAddr)
	require.NoError(r, err)

	payload := randomPayload(r)

	// make the deposit transaction
	resp := r.SuiDepositAndCallSUI(r.TestDAppV2MEVMAddr, math.NewUintFromBigInt(amount), []byte(payload))

	r.Logger.Info("Sui deposit and call tx: %s", resp.Digest)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, resp.Digest, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.EqualValues(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	require.EqualValues(r, coin.CoinType_Gas, cctx.InboundParams.CoinType)
	require.EqualValues(r, amount.Uint64(), cctx.InboundParams.Amount.Uint64())
	require.True(r, cctx.InboundParams.IsCrossChainCall)

	newBalance, err := r.SUIMRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2MEVMAddr)
	require.NoError(r, err)
	require.EqualValues(r, oldBalance.Add(oldBalance, amount).Uint64(), newBalance.Uint64())

	// check sender passed in the call
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err)

	sender, err := sui.EncodeAddress(signer.Address())
	require.NoError(r, err)

	actualSender, err := r.TestDAppV2MEVM.GetSenderWithMessage(&bind.CallOpts{}, payload)
	require.NoError(r, err)
	require.EqualValues(r, sender, actualSender)

	// check the payload was received on the contract
	r.AssertTestDAppMEVMCalled(true, payload, amount)
}
