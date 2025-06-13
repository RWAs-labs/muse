package e2etests

import (
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestERC20DepositAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveERC20OnEVM(r.GatewayEVMAddr)

	payload := randomPayload(r)

	r.AssertTestDAppEVMCalled(false, payload, amount)

	// perform the deposit
	tx := r.ERC20DepositAndCall(r.TestDAppV2MEVMAddr, amount, []byte("revert"), gatewayevm.RevertOptions{
		RevertAddress:    r.TestDAppV2EVMAddr,
		CallOnRevert:     true,
		RevertMessage:    []byte(payload),
		OnRevertGasLimit: big.NewInt(200000),
	})

	// wait for the cctx to be reverted
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit_and_call")
	require.Equal(r, crosschaintypes.CctxStatus_Reverted, cctx.CctxStatus.Status)

	// check the payload was received on the contract
	r.AssertTestDAppEVMCalled(true, payload, big.NewInt(0))

	// check expected sender was used
	senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(payload),
	)
	require.NoError(r, err)
	require.Equal(r, r.EVMAuth.From, senderForMsg)
}
