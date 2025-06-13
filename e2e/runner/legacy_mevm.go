package runner

import (
	"math/big"

	museconnectoreth "github.com/RWAs-labs/protocol-contracts/pkg/museconnector.eth.sol"
	connectormevm "github.com/RWAs-labs/protocol-contracts/pkg/museconnectormevm.sol"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/utils"
)

// LegacySendMuseOnEvm sends MUSE to an address on EVM using legacy protocol contracts
// this allows the MUSE contract deployer to funds other accounts on EVM
func (r *E2ERunner) LegacySendMuseOnEvm(address ethcommon.Address, museAmount int64) *ethtypes.Transaction {
	// the deployer might be sending MUSE in different goroutines
	r.Lock()
	defer r.Unlock()

	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(museAmount))
	tx, err := r.MuseEth.Transfer(r.EVMAuth, address, amount)
	require.NoError(r, err)

	return tx
}

// LegacyDepositMuse deposits MUSE on MuseChain from the MUSE smart contract on EVM using legacy protocol contracts
func (r *E2ERunner) LegacyDepositMuse() ethcommon.Hash {
	amount := big.NewInt(1e18)
	amount = amount.Mul(amount, big.NewInt(100)) // 100 Muse

	return r.LegacyDepositMuseWithAmount(r.EVMAddress(), amount)
}

// LegacyDepositMuseWithAmount deposits MUSE on MuseChain from the MUSE smart contract on EVM with the specified amount using legacy protocol contracts
func (r *E2ERunner) LegacyDepositMuseWithAmount(to ethcommon.Address, amount *big.Int) ethcommon.Hash {
	tx, err := r.MuseEth.Approve(r.EVMAuth, r.ConnectorEthAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "approve")
	r.requireTxSuccessful(receipt, "approve tx failed")

	// query the chain ID using mevm client
	museChainID, err := r.MEVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	tx, err = r.ConnectorEth.Send(r.EVMAuth, museconnectoreth.MuseInterfacesSendInput{
		// TODO: allow user to specify destination chain id
		// https://github.com/RWAs-labs/muse-private/issues/41
		DestinationChainId:  museChainID,
		DestinationAddress:  to.Bytes(),
		DestinationGasLimit: big.NewInt(250_000),
		Message:             nil,
		MuseValueAndGas:     amount,
		MuseParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("Send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "send")
	r.requireTxSuccessful(receipt, "send tx failed")

	r.Logger.Info("  Logs:")
	for _, log := range receipt.Logs {
		sentLog, err := r.ConnectorEth.ParseMuseSent(*log)
		if err == nil {
			r.Logger.Info("    Connector: %s", r.ConnectorEthAddr.String())
			r.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
			r.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
			r.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
			r.Logger.Info("    Muse Value: %d", sentLog.MuseValueAndGas)
			r.Logger.Info("    Block Num: %d", log.BlockNumber)
		}
	}

	return tx.Hash()
}

// LegacyDepositAndApproveWMuse deposits and approves WMUSE on MuseChain from the MUSE smart contract on MEVM using legacy protocol contracts
func (r *E2ERunner) LegacyDepositAndApproveWMuse(amount *big.Int) {
	r.MEVMAuth.Value = amount
	tx, err := r.WMuse.Deposit(r.MEVMAuth)
	require.NoError(r, err)

	r.MEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wmuse deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wmuse deposit")
	r.requireTxSuccessful(receipt, "deposit failed")

	tx, err = r.WMuse.Approve(r.MEVMAuth, r.ConnectorMEVMAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("wmuse approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wmuse approve")
	r.requireTxSuccessful(receipt, "approve failed, logs: %+v", receipt.Logs)
}

// LegacyWithdrawMuse withdraws MUSE from MuseChain to the MUSE smart contract on EVM using legacy protocol contracts
// waitReceipt specifies whether to wait for the tx receipt and check if the tx was successful
func (r *E2ERunner) LegacyWithdrawMuse(amount *big.Int, waitReceipt bool) *ethtypes.Transaction {
	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	tx, err := r.ConnectorMEVM.Send(r.MEVMAuth, connectormevm.MuseInterfacesSendInput{
		DestinationChainId:  chainID,
		DestinationAddress:  r.EVMAddress().Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		MuseValueAndGas:     amount,
		MuseParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("send tx hash: %s", tx.Hash().Hex())

	if waitReceipt {
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
		r.Logger.EVMReceipt(*receipt, "send")
		r.requireTxSuccessful(receipt, "send failed, logs: %+v", receipt.Logs)

		r.Logger.Info("  Logs:")
		for _, log := range receipt.Logs {
			sentLog, err := r.ConnectorMEVM.ParseMuseSent(*log)
			if err == nil {
				r.Logger.Info("    Dest Addr: %s", ethcommon.BytesToAddress(sentLog.DestinationAddress).Hex())
				r.Logger.Info("    Dest Chain: %d", sentLog.DestinationChainId)
				r.Logger.Info("    Dest Gas: %d", sentLog.DestinationGasLimit)
				r.Logger.Info("    Muse Value: %d", sentLog.MuseValueAndGas)
			}
		}
	}

	return tx
}

// LegacyWithdrawEther withdraws Ether from MuseChain to the MUSE smart contract on EVM using legacy protocol contracts
func (r *E2ERunner) LegacyWithdrawEther(amount *big.Int) *ethtypes.Transaction {
	// withdraw
	tx, err := r.ETHMRC20.Withdraw(r.MEVMAuth, r.EVMAddress().Bytes(), amount)
	require.NoError(r, err)

	r.Logger.EVMTransaction(*tx, "withdraw")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt, "withdraw failed")

	r.Logger.EVMReceipt(*receipt, "withdraw")
	r.Logger.MRC20Withdrawal(r.ETHMRC20, *receipt, "withdraw")

	return tx
}

// LegacyWithdrawERC20 withdraws an ERC20 token from MuseChain to the MUSE smart contract on EVM using legacy protocol contracts
func (r *E2ERunner) LegacyWithdrawERC20(amount *big.Int) *ethtypes.Transaction {
	tx, err := r.ERC20MRC20.Withdraw(r.MEVMAuth, r.EVMAddress().Bytes(), amount)
	require.NoError(r, err)

	r.Logger.EVMTransaction(*tx, "withdraw")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := r.ERC20MRC20.ParseWithdrawal(*log)
		if err != nil {
			continue
		}
		r.Logger.Info(
			"  logs: from %s, to %x, value %d, gasfee %d",
			event.From.Hex(),
			event.To,
			event.Value,
			event.GasFee,
		)
	}

	return tx
}
