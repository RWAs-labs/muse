package runner

import (
	"log"
	"math/big"
	"time"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/constant"
)

// ETHDeposit calls Deposit of Gateway with gas token on EVM
func (r *E2ERunner) ETHDeposit(
	receiver ethcommon.Address,
	amount *big.Int,
	revertOptions gatewayevm.RevertOptions,
	wait bool,
) *ethtypes.Transaction {
	r.Lock()
	defer r.Unlock()

	// set the value of the transaction
	previousValue := r.EVMAuth.Value
	defer func() {
		r.EVMAuth.Value = previousValue
	}()
	r.EVMAuth.Value = amount

	tx, err := r.GatewayEVM.Deposit0(r.EVMAuth, receiver, revertOptions)
	require.NoError(r, err)

	if wait {
		logDepositInfoAndWaitForTxReceipt(r, tx, "eth_deposit")
	}

	return tx
}

// DepositEtherDeployer sends Ethers into MEVM using V2 protocol contracts
func (r *E2ERunner) DepositEtherDeployer() ethcommon.Hash {
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(100)) // 100 eth
	tx := r.ETHDeposit(r.EVMAddress(), amount, gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)}, true)
	return tx.Hash()
}

// DepositERC20Deployer sends ERC20 into MEVM using v2 protocol contracts
func (r *E2ERunner) DepositERC20Deployer() ethcommon.Hash {
	r.Logger.Print("⏳ depositing ERC20 into MEVM")

	r.ApproveERC20OnEVM(r.GatewayEVMAddr)

	oneHundred := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(100))
	tx := r.ERC20Deposit(r.EVMAddress(), oneHundred, gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})
	return tx.Hash()
}

// ETHDepositAndCall calls DepositAndCall of Gateway with gas token on EVM
func (r *E2ERunner) ETHDepositAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewayevm.RevertOptions,
) *ethtypes.Transaction {
	// set the value of the transaction
	previousValue := r.EVMAuth.Value
	defer func() {
		r.EVMAuth.Value = previousValue
	}()
	r.EVMAuth.Value = amount

	tx, err := r.GatewayEVM.DepositAndCall(r.EVMAuth, receiver, payload, revertOptions)
	require.NoError(r, err)

	logDepositInfoAndWaitForTxReceipt(r, tx, "eth_deposit_and_call")

	return tx
}

// ERC20Deposit calls Deposit of Gateway with erc20 token on EVM
func (r *E2ERunner) ERC20Deposit(
	receiver ethcommon.Address,
	amount *big.Int,
	revertOptions gatewayevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayEVM.Deposit(r.EVMAuth, receiver, amount, r.ERC20Addr, revertOptions)
	require.NoError(r, err)

	logDepositInfoAndWaitForTxReceipt(r, tx, "erc20_deposit")

	return tx
}

// ERC20DepositAndCall calls DepositAndCall of Gateway with erc20 token on EVM
func (r *E2ERunner) ERC20DepositAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewayevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayEVM.DepositAndCall0(
		r.EVMAuth,
		receiver,
		amount,
		r.ERC20Addr,
		payload,
		revertOptions,
	)
	require.NoError(r, err)

	logDepositInfoAndWaitForTxReceipt(r, tx, "erc20_deposit_and_call")

	return tx
}

// EVMToZEMVCall calls Call of Gateway on EVM
func (r *E2ERunner) EVMToZEMVCall(
	receiver ethcommon.Address,
	payload []byte,
	revertOptions gatewayevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayEVM.Call(r.EVMAuth, receiver, payload, revertOptions)
	require.NoError(r, err)

	return tx
}

// WaitForTxReceiptOnEVM waits for a tx receipt on EVM
func (r *E2ERunner) WaitForTxReceiptOnEVM(tx *ethtypes.Transaction) {
	r.Lock()
	defer r.Unlock()

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)
}

// MintERC20OnEVM mints ERC20 on EVM
// amount is a multiple of 1e18
func (r *E2ERunner) MintERC20OnEVM(amountERC20 int64) {
	r.Lock()
	defer r.Unlock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountERC20))

	tx, err := r.ERC20.Mint(r.EVMAuth, amount)
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)

	r.Logger.Info("Mint receipt tx hash: %s", tx.Hash().Hex())
}

// SendERC20OnEVM sends ERC20 to an address on EVM
// this allows the ERC20 contract deployer to funds other accounts on EVM
// amountERC20 is a multiple of 1e18
func (r *E2ERunner) SendERC20OnEVM(address ethcommon.Address, amountERC20 int64) *ethtypes.Transaction {
	// the deployer might be sending ERC20 in different goroutines
	r.Lock()
	defer r.Unlock()

	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(amountERC20))

	// transfer
	tx, err := r.ERC20.Transfer(r.EVMAuth, address, amount)
	require.NoError(r, err)

	return tx
}

// ApproveERC20OnEVM approves ERC20 on EVM to a specific address
// check if allowance is zero before calling this method
// allow a high amount to avoid multiple approvals
func (r *E2ERunner) ApproveERC20OnEVM(allowed ethcommon.Address) {
	allowance, err := r.ERC20.Allowance(&bind.CallOpts{}, r.Account.EVMAddress(), r.GatewayEVMAddr)
	require.NoError(r, err)

	// approve 1M*1e18 if allowance is zero
	if allowance.Cmp(big.NewInt(0)) == 0 {
		tx, err := r.ERC20.Approve(r.EVMAuth, allowed, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000000)))
		require.NoError(r, err)
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		require.True(r, receipt.Status == 1, "approval failed")
	}
}

// DonateEtherToTSS donates ether to TSS
func (r *E2ERunner) DonateEtherToTSS(amount *big.Int) (*ethtypes.Transaction, error) {
	return r.LegacySendEther(r.TSSAddress, amount, []byte(constant.DonationMessage))
}

// AnvilMineBlocks mines blocks on Anvil localnet
// the block time is provided in seconds
// the method returns a function to stop the mining
func (r *E2ERunner) AnvilMineBlocks(url string, blockTime int) (func(), error) {
	stop := make(chan struct{})

	client, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				time.Sleep(time.Duration(blockTime) * time.Second)

				var result interface{}
				err = client.CallContext(r.Ctx, &result, "evm_mine")
				if err != nil {
					log.Fatalf("Failed to mine a new block: %v", err)
				}
			}
		}
	}()
	return func() {
		close(stop)
	}, nil
}

func logDepositInfoAndWaitForTxReceipt(
	r *E2ERunner,
	tx *ethtypes.Transaction,
	name string,
) {
	r.Logger.EVMTransaction(*tx, name)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt, name+" failed")

	r.Logger.EVMReceipt(*receipt, name)
	r.Logger.GatewayDeposit(r.GatewayEVM, *receipt, name)
}
