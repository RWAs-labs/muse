package legacy

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/testdapp"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	cctxtypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestMessagePassingRevertSuccessExternalChains tests message passing with successful revert between external EVM chains
// TODO: Use two external EVM chains for these tests
// https://github.com/RWAs-labs/muse/issues/2185
func TestMessagePassingRevertSuccessExternalChains(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := utils.ParseBigInt(r, args[0])

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	auth := r.EVMAuth

	tx, err := r.MuseEth.Approve(auth, r.EvmTestDAppAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)
	r.Logger.Info("Approve tx receipt: %d", receipt.Status)

	r.Logger.Info("Calling TestDApp.SendHello on contract address %s", r.EvmTestDAppAddr.Hex())
	testDApp, err := testdapp.NewTestDApp(r.EvmTestDAppAddr, r.EVMClient)
	require.NoError(r, err)

	res2, err := r.BankClient.SupplyOf(r.Ctx, &banktypes.QuerySupplyOfRequest{Denom: "amuse"})
	require.NoError(r, err)

	r.Logger.Info("$$$ Before: SUPPLY OF AMUSE: %d", res2.Amount.Amount)

	tx, err = testDApp.SendHelloWorld(auth, r.EvmTestDAppAddr, chainID, amount, true)
	require.NoError(r, err)

	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("TestDApp.SendHello tx receipt: status %d", receipt.Status)

	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_Reverted)

	outTxHash := cctx.GetCurrentOutboundParam().Hash
	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(outTxHash))
	require.NoError(r, err)

	for _, log := range receipt.Logs {
		event, err := r.ConnectorEth.ParseMuseReverted(*log)
		if err == nil {
			r.Logger.Info("MuseReverted event: ")
			r.Logger.Info("  Dest Addr: %s", ethcommon.BytesToAddress(event.DestinationAddress).Hex())
			r.Logger.Info("  Dest Chain: %d", event.DestinationChainId)
			r.Logger.Info("  RemainingMuseValue: %d", event.RemainingMuseValue)
			r.Logger.Info("  Message: %x", event.Message)
		}
	}

	res3, err := r.BankClient.SupplyOf(r.Ctx, &banktypes.QuerySupplyOfRequest{Denom: "amuse"})
	require.NoError(r, err)

	r.Logger.Info("$$$ After: SUPPLY OF AMUSE: %d", res3.Amount.Amount.BigInt())
	r.Logger.Info("$$$ Diff: SUPPLY OF AMUSE: %d", res3.Amount.Amount.Sub(res2.Amount.Amount).BigInt())
}
