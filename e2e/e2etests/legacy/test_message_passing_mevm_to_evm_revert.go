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

func TestMessagePassingMEVMtoEVMRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := utils.ParseBigInt(r, args[0])

	// Set destination details
	EVMChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	destinationAddress := r.EvmTestDAppAddr

	// Contract call originates from MEVM chain
	r.MEVMAuth.Value = amount
	tx, err := r.WMuse.Deposit(r.MEVMAuth)
	require.NoError(r, err)

	r.MEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wmuse deposit tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wmuse deposit")
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.WMuse.Approve(r.MEVMAuth, r.MevmTestDAppAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("wmuse approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wmuse approve")
	utils.RequireTxSuccessful(r, receipt)

	testDAppMEVM, err := testdapp.NewTestDApp(r.MevmTestDAppAddr, r.MEVMClient)
	require.NoError(r, err)

	// Get MUSE balance before test
	previousBalanceMEVM, err := r.WMuse.BalanceOf(&bind.CallOpts{}, r.MevmTestDAppAddr)
	require.NoError(r, err)

	previousBalanceEVM, err := r.MuseEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	require.NoError(r, err)

	// Call the SendHelloWorld function on the MEVM dapp Contract which would in turn create a new send, to be picked up by the musenode evm hooks
	// set Do revert to true which adds a message to signal the EVM museReceiver to revert the transaction
	tx, err = testDAppMEVM.SendHelloWorld(r.MEVMAuth, destinationAddress, EVMChainID, amount, true)
	require.NoError(r, err)

	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// New inbound message picked up by musenode evm hooks and processed directly to initiate a contract call on EVM which would revert the transaction
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_Reverted)

	// On finalization the Fungible module calls the onRevert function which in turn calls the onMuseRevert function on the sender contract
	receipt, err = r.MEVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	require.NoError(r, err)
	utils.RequireTxSuccessful(r, receipt)

	receivedHelloWorldEvent := false
	for _, log := range receipt.Logs {
		_, err := testDAppMEVM.ParseRevertedHelloWorldEvent(*log)
		if err == nil {
			r.Logger.Info("Received HelloWorld event:")
			receivedHelloWorldEvent = true
		}
	}
	require.True(r, receivedHelloWorldEvent, "expected Reverted HelloWorld event")

	// Check MUSE balance on MEVM TestDApp and check new balance is between previous balance and previous balance + amount
	// New balance is increased because MUSE are sent from the sender but sent back to the contract
	// Contract receive less than the amount because of the gas fee to pay
	newBalanceMEVM, err := r.WMuse.BalanceOf(&bind.CallOpts{}, r.MevmTestDAppAddr)
	require.NoError(r, err)

	previousBalanceAndAmountMEVM := big.NewInt(0).Add(previousBalanceMEVM, amount)

	// check higher than previous balance and lower than previous balance + amount
	invariant := newBalanceMEVM.Cmp(previousBalanceMEVM) <= 0 || newBalanceMEVM.Cmp(previousBalanceAndAmountMEVM) > 0
	require.False(r,
		invariant,
		"expected new balance to be between %s and %s, got %s",
		previousBalanceMEVM.String(),
		previousBalanceAndAmountMEVM.String(),
		newBalanceMEVM.String(),
	)

	// Check MUSE balance on EVM TestDApp and check new balance is previous balance
	newBalanceEVM, err := r.MuseEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	require.NoError(r, err)
	require.Equal(r, 0, newBalanceEVM.Cmp(previousBalanceEVM), "expected new balance to be equal to previous balance")
}
