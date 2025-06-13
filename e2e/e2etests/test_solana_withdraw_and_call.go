package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	solanacontract "github.com/RWAs-labs/muse/pkg/contracts/solana"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestSolanaWithdrawAndCall executes withdrawAndCall on mevm and calls connected program on solana
// message and mevm sender are stored in connected program pda, and withdrawn lamports are stored
// in connected program pda and account provided in remaining accounts to demonstrate that lamports
// can be moved to accounts in connected program as well as gateway program
func TestSolanaWithdrawAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	withdrawAmount := utils.ParseBigInt(r, args[0])

	// get ERC20 SOL balance before withdraw
	balanceBefore, err := r.SOLMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL before withdraw: %d", balanceBefore)

	require.Equal(r, 1, balanceBefore.Cmp(withdrawAmount), "Insufficient balance for withdrawal")

	// parse withdraw amount (in lamports), approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)
	require.Equal(
		r,
		-1,
		withdrawAmount.Cmp(approvedAmount),
		"Withdrawal amount must be less than the approved amount: %v",
		approvedAmount,
	)

	// load deployer private key
	privkey := r.GetSolanaPrivKey()

	// check balances before withdraw
	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(runner.ConnectedProgramID)
	require.NoError(r, err)

	connectedPdaInfoBefore, err := r.SolanaClient.GetAccountInfo(r.Ctx, connectedPda)
	require.NoError(r, err)

	senderBefore, err := r.SolanaClient.GetAccountInfo(r.Ctx, privkey.PublicKey())
	require.NoError(r, err)

	// encode msg
	msg := solanacontract.ExecuteMsg{
		Accounts: []solanacontract.AccountMeta{
			{PublicKey: [32]byte(connectedPda.Bytes()), IsWritable: true},
			{PublicKey: [32]byte(r.ComputePdaAddress().Bytes()), IsWritable: false},
			{PublicKey: [32]byte(r.GetSolanaPrivKey().PublicKey().Bytes()), IsWritable: true},
			{PublicKey: [32]byte(solana.SystemProgramID.Bytes()), IsWritable: false},
		},
		Data: []byte("hello"),
	}

	msgEncoded, err := msg.Encode()
	require.NoError(r, err)

	// withdraw and call
	tx := r.WithdrawAndCallSOLMRC20(
		runner.ConnectedProgramID,
		withdrawAmount,
		approvedAmount,
		msgEncoded,
		gatewaymevm.RevertOptions{
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// get ERC20 SOL balance after withdraw
	balanceAfter, err := r.SOLMRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Info("runner balance of SOL after withdraw: %d", balanceAfter)

	// check if the balance is reduced correctly
	amountReduced := new(big.Int).Sub(balanceBefore, balanceAfter)
	require.True(r, amountReduced.Cmp(withdrawAmount) >= 0, "balance is not reduced correctly")

	// check pda account info of connected program
	connectedPdaInfo, err := r.SolanaClient.GetAccountInfo(r.Ctx, connectedPda)
	require.NoError(r, err)

	sender, err := r.SolanaClient.GetAccountInfo(r.Ctx, privkey.PublicKey())
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

	// connected program splits amount between account provided in remaining accounts, and its own pda
	require.Equal(r, connectedPdaInfoBefore.Value.Lamports+withdrawAmount.Uint64()/2, connectedPdaInfo.Value.Lamports)
	require.Equal(r, senderBefore.Value.Lamports+withdrawAmount.Uint64()/2, sender.Value.Lamports)
}
