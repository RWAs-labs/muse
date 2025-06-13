package legacy

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/testdapp"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	cctxtypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestMessagePassingEVMtoMEVM(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := utils.ParseBigInt(r, args[0])

	// Set destination details
	mEVMChainID, err := r.MEVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	destinationAddress := r.MevmTestDAppAddr

	// Contract call originates from EVM chain
	tx, err := r.MuseEth.Approve(r.EVMAuth, r.EvmTestDAppAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("Approve tx receipt: %d", receipt.Status)
	testDAppEVM, err := testdapp.NewTestDApp(r.EvmTestDAppAddr, r.EVMClient)
	require.NoError(r, err)

	// Get MUSE balance on MEVM TestDApp
	previousBalanceMEVM, err := r.WMuse.BalanceOf(&bind.CallOpts{}, r.MevmTestDAppAddr)
	require.NoError(r, err)
	previousBalanceEVM, err := r.MuseEth.BalanceOf(&bind.CallOpts{}, r.EVMAuth.From)
	require.NoError(r, err)

	// Call the SendHelloWorld function on the EVM dapp Contract which would in turn create a new send, to be picked up by the muse-clients
	// set Do revert to false which adds a message to signal the MEVM museReceiver to not revert the transaction
	tx, err = testDAppEVM.SendHelloWorld(r.EVMAuth, destinationAddress, mEVMChainID, amount, false)
	require.NoError(r, err)
	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// New inbound message picked up by muse-clients and voted on by observers to initiate a contract call on mEVM
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_OutboundMined)

	r.Logger.Info("🔄 Cctx mined for contract call chain mevm %s", cctx.Index)

	// On finalization the Fungible module calls the onReceive function which in turn calls the onMuseMessage function on the destination contract
	receipt, err = r.MEVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	require.NoError(r, err)
	utils.RequireTxSuccessful(r, receipt)

	testDAppMEVM, err := testdapp.NewTestDApp(r.MevmTestDAppAddr, r.MEVMClient)
	require.NoError(r, err)

	// Check event emitted
	receivedHelloWorldEvent := false
	for _, log := range receipt.Logs {
		_, err := testDAppMEVM.ParseHelloWorldEvent(*log)
		if err == nil {
			r.Logger.Info("Received HelloWorld event")
			receivedHelloWorldEvent = true
		}
	}
	require.True(r, receivedHelloWorldEvent, "expected HelloWorld event")

	// Check MUSE balance on MEVM TestDApp and check new balance is previous balance + amount
	newBalanceMEVM, err := r.WMuse.BalanceOf(&bind.CallOpts{}, r.MevmTestDAppAddr)
	require.NoError(r, err)
	require.Equal(r, 0, newBalanceMEVM.Cmp(big.NewInt(0).Add(previousBalanceMEVM, amount)))

	// Check MUSE balance on EVM TestDApp and check new balance is previous balance - amount
	newBalanceEVM, err := r.MuseEth.BalanceOf(&bind.CallOpts{}, r.EVMAuth.From)
	require.NoError(r, err)
	require.Equal(r, 0, newBalanceEVM.Cmp(big.NewInt(0).Sub(previousBalanceEVM, amount)))
}
