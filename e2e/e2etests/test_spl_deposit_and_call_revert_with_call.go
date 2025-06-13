package e2etests

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	testcontract "github.com/RWAs-labs/muse/e2e/contracts/reverter"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	solanacontracts "github.com/RWAs-labs/muse/pkg/contracts/solana"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestSPLDepositAndCallRevertWithCall tests deposit of SPL tokens with revert options
func TestSPLDepositAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)
	amount := utils.ParseInt(r, args[0])

	// deploy a reverter contract in MEVM
	reverterAddr, _, _, err := testcontract.DeployReverter(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Reverter contract deployed at: %s", reverterAddr.String())

	// load deployer private key
	privKey := r.GetSolanaPrivKey()
	r.ResolveSolanaATA(privKey, privKey.PublicKey(), r.SPLAddr)

	data := []byte("hello spl deposit and call")
	connectedPda, err := solanacontracts.ComputeConnectedPdaAddress(runner.ConnectedSPLProgramID)
	require.NoError(r, err)

	connectedPdaAta := r.ResolveSolanaATA(r.GetSolanaPrivKey(), connectedPda, r.SPLAddr)
	connectedPdaBalanceBefore, err := r.SolanaClient.GetTokenAccountBalance(
		r.Ctx,
		connectedPdaAta,
		rpc.CommitmentConfirmed,
	)
	require.NoError(r, err)
	r.Logger.Info("connected pda balance of SPL before revert: %s", connectedPdaBalanceBefore.Value.Amount)

	// create encoded msg
	msg := solanacontracts.ExecuteMsg{
		Accounts: []solanacontracts.AccountMeta{
			{PublicKey: [32]byte(connectedPda.Bytes()), IsWritable: true},
			{PublicKey: [32]byte(connectedPdaAta.Bytes()), IsWritable: true},
			{PublicKey: [32]byte(r.SPLAddr), IsWritable: false},
			{PublicKey: [32]byte(r.ComputePdaAddress().Bytes()), IsWritable: false},
			{PublicKey: [32]byte(solana.TokenProgramID.Bytes()), IsWritable: false},
			{PublicKey: [32]byte(solana.SystemProgramID.Bytes()), IsWritable: false},
		},
		Data: data,
	}

	msgEncoded, err := msg.Encode()
	require.NoError(r, err)

	// #nosec G115 e2eTest - always in range
	sig := r.SPLDepositAndCall(&privKey, uint64(amount), r.SPLAddr, reverterAddr, data, &solanacontracts.RevertOptions{
		RevertAddress:    runner.ConnectedSPLProgramID,
		CallOnRevert:     true,
		RevertMessage:    msgEncoded,
		OnRevertGasLimit: 500000,
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "solana_deposit_spl_and_call_revert_with_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)
	require.Equal(r, cctx.GetCurrentOutboundParam().Receiver, runner.ConnectedSPLProgramID.String())

	require.Contains(r, cctx.CctxStatus.ErrorMessage, utils.ErrHashRevertFoo)

	// verify state and balances are updated
	connectedPdaInfo, err := r.SolanaClient.GetAccountInfo(r.Ctx, connectedPda)
	require.NoError(r, err)
	type ConnectedPdaInfo struct {
		Discriminator     [8]byte
		LastSender        ethcommon.Address
		LastMessage       string
		LastRevertSender  solana.PublicKey
		LastRevertMessage string
	}
	pda := ConnectedPdaInfo{}
	err = borsh.Deserialize(&pda, connectedPdaInfo.Bytes())
	require.NoError(r, err)

	require.Equal(r, "hello spl deposit and call", pda.LastRevertMessage)
	privkey := r.GetSolanaPrivKey()
	require.Equal(r, privkey.PublicKey().String(), pda.LastRevertSender.String())

	connectedPdaBalanceAfter, err := r.SolanaClient.GetTokenAccountBalance(
		r.Ctx,
		connectedPdaAta,
		rpc.CommitmentConfirmed,
	)
	require.NoError(r, err)
	r.Logger.Info("connected pda balance of SPL after revert: %s", connectedPdaBalanceAfter.Value.Amount)

	require.True(
		r,
		utils.ParseUint(r, connectedPdaBalanceAfter.Value.Amount).
			GT(utils.ParseUint(r, connectedPdaBalanceBefore.Value.Amount)),
	)
}
