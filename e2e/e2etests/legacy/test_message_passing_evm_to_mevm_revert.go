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

// fungibleModuleAddress is a constant representing the EVM address of the Fungible module account
const fungibleModuleAddress = "0x735b14BB79463307AAcBED86DAf3322B1e6226aB"

func TestMessagePassingEVMtoMEVMRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	fungibleEthAddress := ethcommon.HexToAddress(fungibleModuleAddress)
	require.True(r, fungibleEthAddress != ethcommon.Address{}, "invalid fungible module address")

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

	// Get MUSE balance before test
	previousBalanceMEVM, err := r.WMuse.BalanceOf(&bind.CallOpts{}, r.MevmTestDAppAddr)
	require.NoError(r, err)

	previousBalanceEVM, err := r.MuseEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	require.NoError(r, err)

	previousFungibleBalance, err := r.WMuse.BalanceOf(&bind.CallOpts{}, fungibleEthAddress)
	require.NoError(r, err)

	// Call the SendHelloWorld function on the EVM dapp Contract which would in turn create a new send, to be picked up by the muse-clients
	// set Do revert to true which adds a message to signal the MEVM museReceiver to revert the transaction
	tx, err = testDAppEVM.SendHelloWorld(r.EVMAuth, destinationAddress, mEVMChainID, amount, true)
	require.NoError(r, err)

	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// New inbound message picked up by muse-clients and voted on by observers to initiate a contract call on mEVM which would revert the transaction
	// A revert transaction is created and gets fialized on the original sender chain.
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_Reverted)

	// On finalization the Tss address calls the onRevert function which in turn calls the onMuseRevert function on the sender contract
	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	require.NoError(r, err)
	utils.RequireTxSuccessful(r, receipt)

	receivedHelloWorldEvent := false
	for _, log := range receipt.Logs {
		_, err := testDAppEVM.ParseRevertedHelloWorldEvent(*log)
		if err == nil {
			r.Logger.Info("Received RevertHelloWorld event:")
			receivedHelloWorldEvent = true
		}
	}
	require.True(r, receivedHelloWorldEvent, "expected Reverted HelloWorld event")

	// Check MUSE balance on MEVM TestDApp and check new balance is previous balance
	newBalanceMEVM, err := r.WMuse.BalanceOf(&bind.CallOpts{}, r.MevmTestDAppAddr)
	require.NoError(r, err)
	require.Equal(
		r,
		0,
		newBalanceMEVM.Cmp(previousBalanceMEVM),
		"expected new balance to be %s, got %s",
		previousBalanceMEVM.String(),
		newBalanceMEVM.String(),
	)

	// Check MUSE balance on EVM TestDApp and check new balance is between previous balance and previous balance + amount
	// New balance is increased because MUSE are sent from the sender but sent back to the contract
	// New balance is less than previous balance + amount because of the gas fee to pay
	newBalanceEVM, err := r.MuseEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	require.NoError(r, err)

	previousBalanceAndAmountEVM := big.NewInt(0).Add(previousBalanceEVM, amount)

	// check higher than previous balance and lower than previous balance + amount
	invariant := newBalanceEVM.Cmp(previousBalanceEVM) <= 0 || newBalanceEVM.Cmp(previousBalanceAndAmountEVM) > 0
	require.False(
		r,
		invariant,
		"expected new balance to be between %s and %s, got %s",
		previousBalanceEVM.String(),
		previousBalanceAndAmountEVM.String(),
		newBalanceEVM.String(),
	)

	// Check MUSE balance on Fungible Module and check new balance is previous balance
	newFungibleBalance, err := r.WMuse.BalanceOf(&bind.CallOpts{}, fungibleEthAddress)
	require.NoError(r, err)

	require.Equal(
		r,
		0,
		newFungibleBalance.Cmp(previousFungibleBalance),
		"expected new balance to be %s, got %s",
		previousFungibleBalance.String(),
		newFungibleBalance.String(),
	)
}
