package e2etests

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	testcontract "github.com/RWAs-labs/muse/e2e/contracts/reverter"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	solanacontracts "github.com/RWAs-labs/muse/pkg/contracts/solana"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestSPLDepositAndCallRevert tests deposit of SPL tokens calling a example contract that reverts.
func TestSPLDepositAndCallRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)
	amount := utils.ParseInt(r, args[0])

	// add liquidity in pool to allow revert fee to be paid
	museAmount := big.NewInt(1e18)
	mrc20Amount := big.NewInt(100000)
	r.AddLiquiditySOL(museAmount, mrc20Amount)
	r.AddLiquiditySPL(museAmount, mrc20Amount)

	// deploy a reverter contract in MEVM
	reverterAddr, _, _, err := testcontract.DeployReverter(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Reverter contract deployed at: %s", reverterAddr.String())

	// load deployer private key
	privKey := r.GetSolanaPrivKey()
	r.ResolveSolanaATA(privKey, privKey.PublicKey(), r.SPLAddr)

	revertAddressPrivateKey, err := solana.NewRandomPrivateKey()
	require.NoError(r, err)
	revertAddressAta := r.ResolveSolanaATA(privKey, revertAddressPrivateKey.PublicKey(), r.SPLAddr)

	// execute the deposit transaction
	data := []byte("hello reverter")
	// #nosec G115 e2eTest - always in range
	sig := r.SPLDepositAndCall(&privKey, uint64(amount), r.SPLAddr, reverterAddr, data, &solanacontracts.RevertOptions{
		RevertAddress: revertAddressPrivateKey.PublicKey(),
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit_spl_and_call_revert")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, revertAddressPrivateKey.PublicKey().String())

	require.Contains(r, cctx.CctxStatus.ErrorMessage, utils.ErrHashRevertFoo)

	// verify balances are updated
	reverterBalance, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, revertAddressAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	require.Greater(r, utils.ParseUint(r, reverterBalance.Value.Amount).Uint64(), uint64(0))
}
