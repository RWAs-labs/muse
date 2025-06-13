package runner

import (
	"fmt"
	"math/big"
	"time"

	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/cenkalti/backoff/v4"
	query "github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/contracts/gatewaymevmcaller"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/retry"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

var gasLimit = big.NewInt(250000)

// ApproveETHMRC20 approves ETH MRC20 on EVM to a specific address
func (r *E2ERunner) ApproveETHMRC20(allowed ethcommon.Address) {
	r.approveMRC20(allowed, r.ETHMRC20)
}

// ApproveERC20MRC20 approves ERC20 MRC20 on EVM to a specific address
func (r *E2ERunner) ApproveERC20MRC20(allowed ethcommon.Address) {
	r.approveMRC20(allowed, r.ERC20MRC20)
}

// ApproveBTCMRC20 approves BTC MRC20 on EVM to a specific address
func (r *E2ERunner) ApproveBTCMRC20(allowed ethcommon.Address) {
	r.approveMRC20(allowed, r.BTCMRC20)
}

// ApproveSOLMRC20 approves SOL MRC20 on EVM to a specific address
func (r *E2ERunner) ApproveSOLMRC20(allowed ethcommon.Address) {
	r.approveMRC20(allowed, r.SOLMRC20)
}

// ApproveSPLMRC20 approves SPL MRC20 on EVM to a specific address
func (r *E2ERunner) ApproveSPLMRC20(allowed ethcommon.Address) {
	r.approveMRC20(allowed, r.SPLMRC20)
}

// ApproveSUIMRC20 approves SUI MRC20 on EVM to a specific address
func (r *E2ERunner) ApproveSUIMRC20(allowed ethcommon.Address) {
	r.approveMRC20(allowed, r.SUIMRC20)
}

// ApproveFungibleTokenMRC20 approves Sui fungible token MRC20 on EVM to a specific address
func (r *E2ERunner) ApproveFungibleTokenMRC20(allowed ethcommon.Address) {
	r.approveMRC20(allowed, r.SuiTokenMRC20)
}

// ApproveTONMRC20 approves TON MRC20 on EVM to a specific address
func (r *E2ERunner) ApproveTONMRC20(allowed ethcommon.Address) {
	r.approveMRC20(allowed, r.TONMRC20)
}

// approveMRC20 approves MRC20 on EVM to a specific address
// check if allowance is zero before calling this method
// allow a high amount to avoid multiple approvals
func (r *E2ERunner) approveMRC20(allowed ethcommon.Address, mrc20 *mrc20.MRC20) {
	allowance, err := mrc20.Allowance(&bind.CallOpts{}, r.Account.EVMAddress(), allowed)
	require.NoError(r, err)

	// approve 1M*1e18 if allowance is below 1k
	thousand := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000))
	if allowance.Cmp(thousand) < 0 {
		r.Logger.Info("Approving %s to %s", r.Account.EVMAddress().String(), allowed.String())
		tx, err := mrc20.Approve(r.MEVMAuth, allowed, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000000)))
		require.NoError(r, err)
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
		require.True(r, receipt.Status == 1, "approval failed")
	}
}

// ETHWithdraw calls Withdraw of Gateway with gas token on MEVM
func (r *E2ERunner) ETHWithdraw(
	receiver ethcommon.Address,
	amount *big.Int,
	revertOptions gatewaymevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayMEVM.Withdraw(
		r.MEVMAuth,
		receiver.Bytes(),
		amount,
		r.ETHMRC20Addr,
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ETHWithdrawAndArbitraryCall calls WithdrawAndCall of Gateway with gas token on MEVM using arbitrary call
func (r *E2ERunner) ETHWithdrawAndArbitraryCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewaymevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayMEVM.WithdrawAndCall0(
		r.MEVMAuth,
		receiver.Bytes(),
		amount,
		r.ETHMRC20Addr,
		payload,
		gatewaymevm.CallOptions{GasLimit: gasLimit, IsArbitraryCall: true},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ETHWithdrawAndCall calls WithdrawAndCall of Gateway with gas token on MEVM using authenticated call
func (r *E2ERunner) ETHWithdrawAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewaymevm.RevertOptions,
	gasLimit *big.Int,
) *ethtypes.Transaction {
	tx, err := r.GatewayMEVM.WithdrawAndCall0(
		r.MEVMAuth,
		receiver.Bytes(),
		amount,
		r.ETHMRC20Addr,
		payload,
		gatewaymevm.CallOptions{
			IsArbitraryCall: false,
			GasLimit:        gasLimit,
		},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ETHWithdrawAndCallThroughContract calls WithdrawAndCall of Gateway with gas token on MEVM using authenticated call
// through contract
func (r *E2ERunner) ETHWithdrawAndCallThroughContract(
	gatewayMEVMCaller *gatewaymevmcaller.GatewayMEVMCaller,
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewaymevmcaller.RevertOptions,
) *ethtypes.Transaction {
	tx, err := gatewayMEVMCaller.WithdrawAndCallGatewayMEVM(
		r.MEVMAuth,
		receiver.Bytes(),
		amount,
		r.ETHMRC20Addr,
		payload,
		gatewaymevmcaller.CallOptions{
			IsArbitraryCall: false,
			GasLimit:        gasLimit,
		},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ERC20Withdraw calls Withdraw of Gateway with erc20 token on MEVM
func (r *E2ERunner) ERC20Withdraw(
	receiver ethcommon.Address,
	amount *big.Int,
	revertOptions gatewaymevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayMEVM.Withdraw(
		r.MEVMAuth,
		receiver.Bytes(),
		amount,
		r.ERC20MRC20Addr,
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ERC20WithdrawAndArbitraryCall calls WithdrawAndCall of Gateway with erc20 token on MEVM using arbitrary call
func (r *E2ERunner) ERC20WithdrawAndArbitraryCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewaymevm.RevertOptions,
) *ethtypes.Transaction {
	// this function take more gas than default 500k
	// so we need to increase the gas limit
	previousGasLimit := r.MEVMAuth.GasLimit
	r.MEVMAuth.GasLimit = 10000000
	defer func() {
		r.MEVMAuth.GasLimit = previousGasLimit
	}()

	tx, err := r.GatewayMEVM.WithdrawAndCall0(
		r.MEVMAuth,
		receiver.Bytes(),
		amount,
		r.ERC20MRC20Addr,
		payload,
		gatewaymevm.CallOptions{GasLimit: gasLimit, IsArbitraryCall: true},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// ERC20WithdrawAndCall calls WithdrawAndCall of Gateway with erc20 token on MEVM using authenticated call
func (r *E2ERunner) ERC20WithdrawAndCall(
	receiver ethcommon.Address,
	amount *big.Int,
	payload []byte,
	revertOptions gatewaymevm.RevertOptions,
) *ethtypes.Transaction {
	// this function take more gas than default 500k
	// so we need to increase the gas limit
	previousGasLimit := r.MEVMAuth.GasLimit
	r.MEVMAuth.GasLimit = 10000000
	defer func() {
		r.MEVMAuth.GasLimit = previousGasLimit
	}()

	tx, err := r.GatewayMEVM.WithdrawAndCall0(
		r.MEVMAuth,
		receiver.Bytes(),
		amount,
		r.ERC20MRC20Addr,
		payload,
		gatewaymevm.CallOptions{GasLimit: gasLimit, IsArbitraryCall: false},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// MEVMToEMVArbitraryCall calls Call of Gateway on MEVM using arbitrary call
func (r *E2ERunner) MEVMToEMVArbitraryCall(
	receiver ethcommon.Address,
	payload []byte,
	revertOptions gatewaymevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayMEVM.Call(
		r.MEVMAuth,
		receiver.Bytes(),
		r.ETHMRC20Addr,
		payload,
		gatewaymevm.CallOptions{GasLimit: gasLimit, IsArbitraryCall: true},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// MEVMToEMVCall calls authenticated Call of Gateway on MEVM using authenticated call
func (r *E2ERunner) MEVMToEMVCall(
	receiver ethcommon.Address,
	payload []byte,
	revertOptions gatewaymevm.RevertOptions,
) *ethtypes.Transaction {
	tx, err := r.GatewayMEVM.Call(
		r.MEVMAuth,
		receiver.Bytes(),
		r.ETHMRC20Addr,
		payload,
		gatewaymevm.CallOptions{
			GasLimit:        gasLimit,
			IsArbitraryCall: false,
		},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// MEVMToEMVCallThroughContract calls authenticated Call of Gateway on MEVM through contract using authenticated call
func (r *E2ERunner) MEVMToEMVCallThroughContract(
	gatewayMEVMCaller *gatewaymevmcaller.GatewayMEVMCaller,
	receiver ethcommon.Address,
	payload []byte,
	revertOptions gatewaymevmcaller.RevertOptions,
) *ethtypes.Transaction {
	tx, err := gatewayMEVMCaller.CallGatewayMEVM(
		r.MEVMAuth,
		receiver.Bytes(),
		r.ETHMRC20Addr,
		payload,
		gatewaymevmcaller.CallOptions{
			GasLimit:        gasLimit,
			IsArbitraryCall: false,
		},
		revertOptions,
	)
	require.NoError(r, err)

	return tx
}

// WaitForBlocks waits for a specific number of blocks to be generated
// The parameter n is the number of blocks to wait for
func (r *E2ERunner) WaitForBlocks(n int64) {
	height, err := r.CctxClient.LastMuseHeight(r.Ctx, &types.QueryLastMuseHeightRequest{})
	if err != nil {
		return
	}
	call := func() error {
		return retry.Retry(r.waitForBlock(height.Height + n))
	}
	retryBuffer := uint64(20)
	bo := backoff.NewConstantBackOff(time.Second * 6)
	// #nosec G115 always in range
	boWithMaxRetries := backoff.WithMaxRetries(bo, uint64(n)+retryBuffer)
	err = retry.DoWithBackoff(call, boWithMaxRetries)
	require.NoError(r, err, "failed to wait for %d blocks", n)
}

// WaitForTSSGeneration waits for a specific number of TSS to be generated
// The parameter n is the number of TSS to wait for
func (r *E2ERunner) WaitForTSSGeneration(tssNumber int64) {
	call := func() error {
		return retry.Retry(r.checkNumberOfTSSGenerated(tssNumber))
	}
	bo := backoff.NewConstantBackOff(time.Second * 5)
	boWithMaxRetries := backoff.WithMaxRetries(bo, 10)
	err := retry.DoWithBackoff(call, boWithMaxRetries)
	require.NoError(r, err, "failed to wait for %d tss generation", tssNumber)
}

// checkNumberOfTSSGenerated checks the number of TSS generated
// if the number of tss is less that the `tssNumber` provided we return an error
func (r *E2ERunner) checkNumberOfTSSGenerated(tssNumber int64) error {
	tssList, err := r.ObserverClient.TssHistory(r.Ctx, &observertypes.QueryTssHistoryRequest{})
	if err != nil {
		return err
	}
	if int64(len(tssList.TssList)) < tssNumber {
		return fmt.Errorf("waiting for %d tss generation, number of TSS :%d", tssNumber, len(tssList.TssList))
	}
	return nil
}

func (r *E2ERunner) waitForBlock(n int64) error {
	height, err := r.CctxClient.LastMuseHeight(r.Ctx, &types.QueryLastMuseHeightRequest{})
	if err != nil {
		return err
	}
	if height.Height < n {
		return fmt.Errorf("waiting for height: %d, current height: %d", n, height.Height)
	}
	return nil
}

// WaitForTxReceiptOnMEVM waits for a tx receipt on MEVM
func (r *E2ERunner) WaitForTxReceiptOnMEVM(tx *ethtypes.Transaction) {
	r.Lock()
	defer r.Unlock()

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.requireTxSuccessful(receipt)
}

// WaitForMinedCCTX waits for a cctx to be mined from a tx
func (r *E2ERunner) WaitForMinedCCTX(txHash ethcommon.Hash) {
	r.Lock()
	defer r.Unlock()

	cctx := utils.WaitCctxMinedByInboundHash(
		r.Ctx,
		txHash.Hex(),
		r.CctxClient,
		r.Logger,
		r.CctxTimeout,
	)
	utils.RequireCCTXStatus(r, cctx, types.CctxStatus_OutboundMined)
}

// WaitForMinedCCTXFromIndex waits for a cctx to be mined from its index
func (r *E2ERunner) WaitForMinedCCTXFromIndex(index string) *types.CrossChainTx {
	return r.waitForMinedCCTXFromIndex(index, types.CctxStatus_OutboundMined)
}

func (r *E2ERunner) waitForMinedCCTXFromIndex(index string, status types.CctxStatus) *types.CrossChainTx {
	r.Lock()
	defer r.Unlock()

	cctx := utils.WaitCCTXMinedByIndex(r.Ctx, index, r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, status)

	return cctx
}

// WaitForSpecificCCTX scans for cctx by filters and ensures it's mined
func (r *E2ERunner) WaitForSpecificCCTX(
	filter func(*types.CrossChainTx) bool,
	status types.CctxStatus,
	timeout time.Duration,
) *types.CrossChainTx {
	var (
		ctx      = r.Ctx
		start    = time.Now()
		reqQuery = &types.QueryAllCctxRequest{
			Pagination: &query.PageRequest{
				Limit:   50,
				Reverse: false,
			},
		}
	)

	for time.Since(start) < timeout {
		res, err := r.CctxClient.CctxAll(ctx, reqQuery)
		require.NoError(r, err)

		for i := range res.CrossChainTx {
			tx := res.CrossChainTx[i]
			if filter(tx) {
				return r.waitForMinedCCTXFromIndex(tx.Index, status)
			}
		}

		time.Sleep(time.Second)
	}

	r.Logger.Error("WaitForSpecificCCTX: No CCTX found. Timed out")
	r.FailNow()

	return nil
}

// skipChainOperations checks if the chain operations should be skipped for E2E
func (r *E2ERunner) skipChainOperations(chainID int64) bool {
	skip := r.IsRunningUpgrade() && chains.IsTONChain(chainID, nil)

	if skip {
		r.Logger.Print("Skipping chain operations for chain %d", chainID)
	}

	return skip
}

// AddInboundTracker adds an inbound tracker from the tx hash
func (r *E2ERunner) AddInboundTracker(coinType coin.CoinType, txHash string) {
	require.NotNil(r, r.MuseTxServer)

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	msg := types.NewMsgAddInboundTracker(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		chainID.Int64(),
		coinType,
		txHash,
	)
	_, err = r.MuseTxServer.BroadcastTx(utils.EmergencyPolicyName, msg)
	require.NoError(r, err)
}
