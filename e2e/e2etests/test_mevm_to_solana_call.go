package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	solanacontract "github.com/RWAs-labs/muse/pkg/contracts/solana"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestMEVMToSolanaCall executes simple call from MEVM to Solana, calling connected program
func TestMEVMToSolanaCall(r *runner.E2ERunner, _ []string) {
	// approve amount is 1 SOL
	approvedAmount := new(big.Int).SetUint64(solana.LAMPORTS_PER_SOL)

	r.MEVMAuth.GasLimit = 10000000
	// withdraw and call
	tx := r.CallSOLMRC20(
		runner.ConnectedProgramID,
		approvedAmount,
		[]byte("simple call"),
		gatewaymevm.RevertOptions{
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check pda account info of connected program
	connectedPda, err := solanacontract.ComputeConnectedPdaAddress(runner.ConnectedProgramID)
	require.NoError(r, err)
	connectedPdaInfo, err := r.SolanaClient.GetAccountInfo(r.Ctx, connectedPda)
	require.NoError(r, err)

	type ConnectedPdaInfo struct {
		Discriminator [8]byte
		LastSender    common.Address
		LastMessage   string
	}
	pda := ConnectedPdaInfo{}
	err = borsh.Deserialize(&pda, connectedPdaInfo.Bytes())
	require.NoError(r, err)

	require.Equal(r, "simple call", pda.LastMessage)
	require.Equal(r, r.MEVMAuth.From.String(), common.BytesToAddress(pda.LastSender[:]).String())
}
