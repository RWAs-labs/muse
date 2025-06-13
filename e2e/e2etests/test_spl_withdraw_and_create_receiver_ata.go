package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestSPLWithdrawAndCreateReceiverAta withdraws SPL, but letting gateway program to create receiver ATA
// using PDA funds, instead of providing already created ATA
func TestSPLWithdrawAndCreateReceiverAta(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawAmount := utils.ParseBigInt(r, args[0])

	// get SPL MRC20 balance before withdraw
	mrc20BalanceBefore, err := r.SPLMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SPL before withdraw: %d", mrc20BalanceBefore)

	require.Equal(r, 1, mrc20BalanceBefore.Cmp(withdrawAmount), "Insufficient balance for withdrawal")

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the %v",
		approvedAmount,
	)

	// create new priv key, with empty ata
	receiverPrivKey, err := solana.NewRandomPrivateKey()
	require.NoError(r, err)

	// verify receiver ata account doesn't exist
	receiverAta, _, err := solana.FindAssociatedTokenAddress(receiverPrivKey.PublicKey(), r.SPLAddr)
	require.NoError(r, err)

	receiverAtaAcc, err := r.SolanaClient.GetAccountInfoWithOpts(
		r.Ctx,
		receiverAta,
		&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed},
	)
	require.Error(r, err)
	require.Nil(r, receiverAtaAcc)

	// withdraw
	tx := r.WithdrawSPLMRC20(receiverPrivKey.PublicKey(), withdrawAmount, approvedAmount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// get SPL MRC20 balance after withdraw
	mrc20BalanceAfter, err := r.SPLMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SPL after withdraw: %d", mrc20BalanceAfter)

	// verify receiver ata was created
	receiverAtaAcc, err = r.SolanaClient.GetAccountInfoWithOpts(
		r.Ctx,
		receiverAta,
		&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed},
	)
	require.NoError(r, err)
	require.NotNil(r, receiverAtaAcc)

	// verify balances are updated
	receiverBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, receiverAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	r.Logger.Info("receiver balance of SPL after withdraw: %s", receiverBalanceAfter.Value.Amount)

	// verify amount is added to receiver ata
	require.EqualValues(r, withdrawAmount.String(), utils.ParseBigInt(r, receiverBalanceAfter.Value.Amount).String())

	// verify amount is subtracted on mrc20
	require.EqualValues(r, new(big.Int).Sub(mrc20BalanceBefore, withdrawAmount).String(), mrc20BalanceAfter.String())
}
