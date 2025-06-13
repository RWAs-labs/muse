package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	solanacontract "github.com/RWAs-labs/muse/pkg/contracts/solana"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestSPLWithdrawAndCall executes withdrawAndCall on mevm and calls connected program on solana
// message and mevm sender are stored in connected program pda, and withdrawn spl tokens are stored
// in connected program pda and account provided in remaining accounts to demonstrate that spl tokens
// can be moved to accounts in connected program as well as gateway program
func TestSPLWithdrawAndCall(r *runner.E2ERunner, args []string) {
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

	// load deployer private key
	privkey := r.GetSolanaPrivKey()

	// get receiver ata balance before withdraw
	receiverAta := r.ResolveSolanaATA(privkey, privkey.PublicKey(), r.SPLAddr)
	receiverBalanceBefore, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, receiverAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	r.Logger.Info("receiver balance of SPL before withdraw: %s", receiverBalanceBefore.Value.Amount)

	connected := solana.MustPublicKeyFromBase58(runner.ConnectedSPLProgramID.String())
	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(connected)
	require.NoError(r, err)

	connectedPdaAta := r.ResolveSolanaATA(privkey, connectedPda, r.SPLAddr)
	connectedPdaBalanceBefore, err := r.SolanaClient.GetTokenAccountBalance(
		r.Ctx,
		connectedPdaAta,
		rpc.CommitmentConfirmed,
	)
	require.NoError(r, err)
	r.Logger.Info("connected pda balance of SPL before withdraw: %s", connectedPdaBalanceBefore.Value.Amount)

	// withdraw
	tx := r.WithdrawAndCallSPLMRC20(
		runner.ConnectedSPLProgramID,
		withdrawAmount,
		approvedAmount,
		[]byte("hello"),
		gatewaymevm.RevertOptions{
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// get SPL MRC20 balance after withdraw
	mrc20BalanceAfter, err := r.SPLMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SPL after withdraw: %d", mrc20BalanceAfter)

	// check pda account info of connected program
	connectedPdaInfo, err := r.SolanaClient.GetAccountInfo(r.Ctx, connectedPda)
	require.NoError(r, err)

	type ConnectedPdaInfo struct {
		Discriminator     [8]byte
		LastSender        common.Address
		LastMessage       string
		LastRevertSender  solana.PublicKey
		LastRevertMessage string
	}
	pda := ConnectedPdaInfo{}
	err = borsh.Deserialize(&pda, connectedPdaInfo.Bytes())
	require.NoError(r, err)

	require.Equal(r, "hello", pda.LastMessage)
	require.Equal(r, r.MEVMAuth.From.String(), common.BytesToAddress(pda.LastSender[:]).String())

	// verify balances are updated
	receiverBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(r.Ctx, receiverAta, rpc.CommitmentConfirmed)
	require.NoError(r, err)
	r.Logger.Info("receiver balance of SPL after withdraw: %s", receiverBalanceAfter.Value.Amount)

	connectedPdaBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(
		r.Ctx,
		connectedPdaAta,
		rpc.CommitmentConfirmed,
	)
	require.NoError(r, err)
	r.Logger.Info("connected pda balance of SPL after withdraw: %s", connectedPdaBalanceAfter.Value.Amount)

	// verify half of amount is added to receiver ata and half to connected pda ata
	halfWithdrawAmount := new(big.Int).Div(withdrawAmount, big.NewInt(2))
	require.EqualValues(
		r,
		new(big.Int).Add(halfWithdrawAmount, utils.ParseBigInt(r, receiverBalanceBefore.Value.Amount)).String(),
		utils.ParseBigInt(r, receiverBalanceAfter.Value.Amount).String(),
	)

	require.EqualValues(
		r,
		new(big.Int).Add(halfWithdrawAmount, utils.ParseBigInt(r, connectedPdaBalanceBefore.Value.Amount)).String(),
		utils.ParseBigInt(r, connectedPdaBalanceAfter.Value.Amount).String(),
	)

	// verify amount is subtracted on mrc20
	require.EqualValues(r, new(big.Int).Sub(mrc20BalanceBefore, withdrawAmount).String(), mrc20BalanceAfter.String())
}
