package legacy

import (
	"math/big"

	museconnectoreth "github.com/RWAs-labs/protocol-contracts/pkg/museconnector.eth.sol"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	cctxtypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestMessagePassingExternalChains tests message passing between external EVM chains
// TODO: Use two external EVM chains for these tests
// https://github.com/RWAs-labs/muse/issues/2185
func TestMessagePassingExternalChains(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := utils.ParseBigInt(r, args[0])

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	r.Logger.Info("Approving ConnectorEth to spend deployer's MuseEth")
	auth := r.EVMAuth

	tx, err := r.MuseEth.Approve(auth, r.ConnectorEthAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("Approve tx receipt: %d", receipt.Status)
	r.Logger.Info("Calling ConnectorEth.Send")
	tx, err = r.ConnectorEth.Send(auth, museconnectoreth.MuseInterfacesSendInput{
		DestinationChainId:  chainID,
		DestinationAddress:  r.EVMAddress().Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		MuseValueAndGas:     amount,
		MuseParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("ConnectorEth.Send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("ConnectorEth.Send tx receipt: status %d", receipt.Status)
	r.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := r.ConnectorEth.ParseMuseSent(*log)
		if err == nil {
			r.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			r.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			r.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			r.Logger.Info("    Muse Value: %d", sentLog.MuseValueAndGas)
		}
	}

	r.Logger.Info("Waiting for ConnectorEth.Send CCTX to be mined...")
	r.Logger.Info("  INTX hash: %s", receipt.TxHash.String())

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_OutboundMined)

	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	require.NoError(r, err)
	utils.RequireTxSuccessful(r, receipt)

	for _, log := range receipt.Logs {
		event, err := r.ConnectorEth.ParseMuseReceived(*log)
		if err == nil {
			r.Logger.Info("Received MuseSent event:")
			r.Logger.Info("  Dest Addr: %s", event.DestinationAddress)
			r.Logger.Info("  Muse Value: %d", event.MuseValue)
			r.Logger.Info("  src chainid: %d", event.SourceChainId)

			comp := event.MuseValue.Cmp(cctx.GetCurrentOutboundParam().Amount.BigInt())
			require.Equal(r, 0, comp, "Muse value mismatch")
		}
	}
}
