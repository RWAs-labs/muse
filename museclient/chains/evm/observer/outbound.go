package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/RWAs-labs/protocol-contracts/pkg/erc20custody.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/museconnector.non-eth.sol"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/chains/evm/common"
	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/compliance"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/musecore"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// ProcessOutboundTrackers processes outbound trackers
func (ob *Observer) ProcessOutboundTrackers(ctx context.Context) error {
	chainID := ob.Chain().ChainId
	trackers, err := ob.MusecoreClient().GetAllOutboundTrackerByChain(ctx, ob.Chain().ChainId, interfaces.Ascending)
	if err != nil {
		return errors.Wrap(err, "GetAllOutboundTrackerByChain error")
	}

	// keep last block up-to-date
	if err := ob.updateLastBlock(ctx); err != nil {
		return err
	}

	// prepare logger fields
	logger := ob.Logger().Outbound.With().
		Str(logs.FieldMethod, "ProcessOutboundTrackers").
		Int64(logs.FieldChain, chainID).
		Logger()

	// process outbound trackers
	for _, tracker := range trackers {
		// go to next tracker if this one already has a confirmed tx
		nonce := tracker.Nonce
		if ob.isTxConfirmed(nonce) {
			continue
		}

		// check each txHash and save tx and receipt if it's legit and confirmed
		txCount := 0
		var outboundReceipt *ethtypes.Receipt
		var outbound *ethtypes.Transaction
		for _, txHash := range tracker.HashList {
			if receipt, tx, ok := ob.checkConfirmedTx(ctx, txHash.TxHash, nonce); ok {
				txCount++
				outboundReceipt = receipt
				outbound = tx
				logger.Info().Msgf("confirmed outbound %s for chain %d nonce %d", txHash.TxHash, chainID, nonce)
				if txCount > 1 {
					logger.Error().
						Msgf("checkConfirmedTx passed, txCount %d chain %d nonce %d receipt %v tx %v", txCount, chainID, nonce, receipt, tx)
				}
			}
		}

		// should be only one txHash confirmed for each nonce.
		if txCount == 1 {
			ob.setTxNReceipt(nonce, outboundReceipt, outbound)
		} else if txCount > 1 {
			// should not happen. We can't tell which txHash is true. It might happen (e.g. bug, glitchy/hacked endpoint)
			ob.Logger().Outbound.Error().Msgf("WatchOutbound: confirmed multiple (%d) outbound for chain %d nonce %d", txCount, chainID, nonce)
		} else {
			if tracker.MaxReached() {
				ob.Logger().Outbound.Error().Msgf("WatchOutbound: outbound tracker is full of hashes for chain %d nonce %d", chainID, nonce)
			}
		}
	}

	return nil
}

// postVoteOutbound posts vote to musecore for the confirmed outbound
func (ob *Observer) postVoteOutbound(
	ctx context.Context,
	cctxIndex string,
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
	receiveValue *big.Int,
	receiveStatus chains.ReceiveStatus,
	nonce uint64,
	coinType coin.CoinType,
	logger zerolog.Logger,
) {
	chainID := ob.Chain().ChainId

	signerAddress := ob.MusecoreClient().GetKeys().GetOperatorAddress()

	msg := crosschaintypes.NewMsgVoteOutbound(
		signerAddress.String(),
		cctxIndex,
		receipt.TxHash.Hex(),
		receipt.BlockNumber.Uint64(),
		receipt.GasUsed,
		math.NewIntFromBigInt(transaction.GasPrice()),
		transaction.Gas(),
		math.NewUintFromBigInt(receiveValue),
		receiveStatus,
		chainID,
		nonce,
		coinType,
		crosschaintypes.ConfirmationMode_SAFE,
	)

	const gasLimit = musecore.PostVoteOutboundGasLimit

	retryGasLimit := musecore.PostVoteOutboundRetryGasLimit
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = musecore.PostVoteOutboundRevertGasLimit
	}

	// post vote to musecore
	logFields := map[string]any{
		"chain":    chainID,
		"nonce":    nonce,
		"outbound": receipt.TxHash.String(),
	}
	museTxHash, ballot, err := ob.MusecoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	if err != nil {
		logger.Error().
			Err(err).
			Fields(logFields).
			Msgf("PostVoteOutbound: error posting vote for chain %d", chainID)
		return
	}

	// print vote tx hash and ballot
	if museTxHash != "" {
		logFields["vote"] = museTxHash
		logFields["ballot"] = ballot
		logger.Info().Fields(logFields).Msgf("PostVoteOutbound: posted vote for chain %d", chainID)
	}
}

// VoteOutboundIfConfirmed checks outbound status and returns (continueKeysign, error)
func (ob *Observer) VoteOutboundIfConfirmed(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
) (bool, error) {
	// skip if outbound is not confirmed
	nonce := cctx.GetCurrentOutboundParam().TssNonce
	if !ob.isTxConfirmed(nonce) {
		return true, nil
	}
	receipt, transaction := ob.getTxNReceipt(nonce)
	sendID := fmt.Sprintf("%d-%d", ob.Chain().ChainId, nonce)
	logger := ob.Logger().Outbound.With().Str("sendID", sendID).Logger()

	// get connector and erc20Custody contracts
	connectorAddr, connector, err := ob.getConnectorContract()
	if err != nil {
		return true, errors.Wrapf(err, "error getting muse connector for chain %d", ob.Chain().ChainId)
	}
	custodyAddr, custody, err := ob.getERC20CustodyContract()
	if err != nil {
		return true, errors.Wrapf(err, "error getting erc20 custody for chain %d", ob.Chain().ChainId)
	}
	gatewayAddr, gateway, err := ob.getGatewayContract()
	if err != nil {
		return true, errors.Wrap(err, "error getting gateway for chain")
	}
	_, custodyV2, err := ob.getERC20CustodyV2Contract()
	if err != nil {
		return true, errors.Wrapf(err, "error getting erc20 custody v2 for chain %d", ob.Chain().ChainId)
	}

	// define a few common variables
	var (
		receiveValue  *big.Int
		receiveStatus chains.ReceiveStatus
		cointype      = cctx.InboundParams.CoinType
	)

	// cancelled transaction means the outbound is failed
	// - set amount to CCTX's amount to bypass amount check in musecore
	// - set status to failed to revert the CCTX in musecore
	if compliance.IsCCTXRestricted(cctx) {
		receiveValue = cctx.GetCurrentOutboundParam().Amount.BigInt()
		receiveStatus = chains.ReceiveStatus_failed
		ob.postVoteOutbound(ctx, cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
		return false, nil
	}

	// parse the received value from the outbound receipt
	receiveValue, receiveStatus, err = parseOutboundReceivedValue(
		cctx,
		receipt,
		transaction,
		cointype,
		connectorAddr,
		connector,
		custodyAddr,
		custody,
		custodyV2,
		gatewayAddr,
		gateway,
	)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("VoteOutboundIfConfirmed: error parsing outbound event for chain %d txhash %s", ob.Chain().ChainId, receipt.TxHash)
		return true, err
	}

	// post vote to musecore
	ob.postVoteOutbound(ctx, cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
	return false, nil
}

// parseOutboundReceivedValue parses the received value and status from the outbound receipt
// The received value is the amount of Muse/ERC20/Gas token (released from connector/custody/TSS) sent to the receiver
// TODO: simplify this function and reduce the number of argument
// https://github.com/RWAs-labs/muse/issues/2627
// https://github.com/RWAs-labs/muse/pull/2666#discussion_r1718379784
func parseOutboundReceivedValue(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
	cointype coin.CoinType,
	connectorAddress ethcommon.Address,
	connector *museconnector.MuseConnectorNonEth,
	custodyAddress ethcommon.Address,
	custody *erc20custody.ERC20Custody,
	custodyV2 *erc20custody.ERC20Custody,
	gatewayAddress ethcommon.Address,
	gateway *gatewayevm.GatewayEVM,
) (*big.Int, chains.ReceiveStatus, error) {
	// determine the receive status and value
	// https://docs.nethereum.com/en/latest/nethereum-receipt-status/
	receiveValue := big.NewInt(0)
	receiveStatus := chains.ReceiveStatus_failed
	if receipt.Status == ethtypes.ReceiptStatusSuccessful {
		receiveValue = transaction.Value()
		receiveStatus = chains.ReceiveStatus_success
	}

	// parse outbound event for protocol contract v2
	if cctx.ProtocolContractVersion == crosschaintypes.ProtocolContractVersion_V2 {
		return parseOutboundEventV2(cctx, receipt, transaction, custodyAddress, custodyV2, gatewayAddress, gateway)
	}

	// parse receive value from the outbound receipt for Muse and ERC20
	switch cointype {
	case coin.CoinType_Muse:
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			receivedLog, revertedLog, err := parseAndCheckMuseEvent(cctx, receipt, connectorAddress, connector)
			if err != nil {
				return nil, chains.ReceiveStatus_failed, err
			}
			// use the value in MuseReceived/MuseReverted event for vote message
			if receivedLog != nil {
				receiveValue = receivedLog.MuseValue
			} else if revertedLog != nil {
				receiveValue = revertedLog.RemainingMuseValue
			}
		}
	case coin.CoinType_ERC20:
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			withdrawn, err := parseAndCheckWithdrawnEvent(cctx, receipt, custodyAddress, custody)
			if err != nil {
				return nil, chains.ReceiveStatus_failed, err
			}
			// use the value in Withdrawn event for vote message
			receiveValue = withdrawn.Amount
		}
	case coin.CoinType_Gas, coin.CoinType_Cmd:
		// nothing to do for CoinType_Gas/CoinType_Cmd, no need to parse event
	default:
		return nil, chains.ReceiveStatus_failed, fmt.Errorf("unknown coin type %s", cointype)
	}
	return receiveValue, receiveStatus, nil
}

// parseAndCheckMuseEvent parses and checks MuseReceived/MuseReverted event from the outbound receipt
// It either returns an MuseReceived or an MuseReverted event, or an error if no event found
func parseAndCheckMuseEvent(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	connectorAddr ethcommon.Address,
	connector *museconnector.MuseConnectorNonEth,
) (*museconnector.MuseConnectorNonEthMuseReceived, *museconnector.MuseConnectorNonEthMuseReverted, error) {
	params := cctx.GetCurrentOutboundParam()
	for _, vLog := range receipt.Logs {
		// try parsing MuseReceived event
		received, err := connector.MuseConnectorNonEthFilterer.ParseMuseReceived(*vLog)
		if err == nil {
			err = common.ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), common.TopicsMuseReceived)
			if err != nil {
				return nil, nil, errors.Wrap(err, "error validating MuseReceived event")
			}
			if !strings.EqualFold(received.DestinationAddress.Hex(), params.Receiver) {
				return nil, nil, fmt.Errorf("receiver address mismatch in MuseReceived event, want %s got %s",
					params.Receiver, received.DestinationAddress.Hex())
			}
			if received.MuseValue.Cmp(params.Amount.BigInt()) != 0 {
				return nil, nil, fmt.Errorf("amount mismatch in MuseReceived event, want %s got %s",
					params.Amount.String(), received.MuseValue.String())
			}
			if ethcommon.BytesToHash(received.InternalSendHash[:]).Hex() != cctx.Index {
				return nil, nil, fmt.Errorf("cctx index mismatch in MuseReceived event, want %s got %s",
					cctx.Index, hex.EncodeToString(received.InternalSendHash[:]))
			}
			return received, nil, nil
		}
		// try parsing MuseReverted event
		reverted, err := connector.MuseConnectorNonEthFilterer.ParseMuseReverted(*vLog)
		if err == nil {
			err = common.ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), common.TopicsMuseReverted)
			if err != nil {
				return nil, nil, errors.Wrap(err, "error validating MuseReverted event")
			}
			if !strings.EqualFold(
				ethcommon.BytesToAddress(reverted.DestinationAddress[:]).Hex(),
				cctx.InboundParams.Sender,
			) {
				return nil, nil, fmt.Errorf("receiver address mismatch in MuseReverted event, want %s got %s",
					cctx.InboundParams.Sender, ethcommon.BytesToAddress(reverted.DestinationAddress[:]).Hex())
			}
			if reverted.RemainingMuseValue.Cmp(params.Amount.BigInt()) != 0 {
				return nil, nil, fmt.Errorf("amount mismatch in MuseReverted event, want %s got %s",
					params.Amount.String(), reverted.RemainingMuseValue.String())
			}
			if ethcommon.BytesToHash(reverted.InternalSendHash[:]).Hex() != cctx.Index {
				return nil, nil, fmt.Errorf("cctx index mismatch in MuseReverted event, want %s got %s",
					cctx.Index, hex.EncodeToString(reverted.InternalSendHash[:]))
			}
			return nil, reverted, nil
		}
	}
	return nil, nil, errors.New("no MuseReceived/MuseReverted event found")
}

// parseAndCheckWithdrawnEvent parses and checks erc20 Withdrawn event from the outbound receipt
func parseAndCheckWithdrawnEvent(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	custodyAddr ethcommon.Address,
	custody *erc20custody.ERC20Custody,
) (*erc20custody.ERC20CustodyWithdrawn, error) {
	params := cctx.GetCurrentOutboundParam()
	for _, vLog := range receipt.Logs {
		withdrawn, err := custody.ParseWithdrawn(*vLog)
		if err == nil {
			err = common.ValidateEvmTxLog(vLog, custodyAddr, receipt.TxHash.Hex(), common.TopicsWithdrawn)
			if err != nil {
				return nil, errors.Wrap(err, "error validating Withdrawn event")
			}
			if !strings.EqualFold(withdrawn.To.Hex(), params.Receiver) {
				return nil, fmt.Errorf("receiver address mismatch in Withdrawn event, want %s got %s",
					params.Receiver, withdrawn.To.Hex())
			}
			if !strings.EqualFold(withdrawn.Token.Hex(), cctx.InboundParams.Asset) {
				return nil, fmt.Errorf("asset mismatch in Withdrawn event, want %s got %s",
					cctx.InboundParams.Asset, withdrawn.Token.Hex())
			}
			if withdrawn.Amount.Cmp(params.Amount.BigInt()) != 0 {
				return nil, fmt.Errorf("amount mismatch in Withdrawn event, want %s got %s",
					params.Amount.String(), withdrawn.Amount.String())
			}
			return withdrawn, nil
		}
	}
	return nil, errors.New("no ERC20 Withdrawn event found")
}

// filterTSSOutbound filters the outbounds from TSS address to supplement outbound trackers
func (ob *Observer) filterTSSOutbound(ctx context.Context, startBlock, toBlock uint64) {
	// filters the outbounds from TSS address block by block
	for bn := startBlock; bn <= toBlock; bn++ {
		ob.filterTSSOutboundInBlock(ctx, bn)
	}
}

// filterTSSOutboundInBlock filters the outbounds in a single block to supplement outbound trackers
func (ob *Observer) filterTSSOutboundInBlock(ctx context.Context, blockNumber uint64) {
	// query block and ignore error (we don't rescan as we are only supplementing outbound trackers)
	block, err := ob.GetBlockByNumberCached(ctx, blockNumber)
	if err != nil {
		ob.Logger().
			Outbound.Error().
			Err(err).
			Msgf("error getting block %d for chain %d", blockNumber, ob.Chain().ChainId)
		return
	}

	for i := range block.Transactions {
		tx := block.Transactions[i]
		if ethcommon.HexToAddress(tx.From) == ob.TSS().PubKey().AddressEVM() {
			// #nosec G115 nonce always positive
			nonce := uint64(tx.Nonce)
			if !ob.isTxConfirmed(nonce) {
				if receipt, txx, ok := ob.checkConfirmedTx(ctx, tx.Hash, nonce); ok {
					ob.setTxNReceipt(nonce, receipt, txx)
					ob.Logger().
						Outbound.Info().
						Msgf("TSS outbound detected on chain %d nonce %d tx %s", ob.Chain().ChainId, nonce, tx.Hash)
				}
			}
		}
	}
}

// checkConfirmedTx checks if a txHash is confirmed
// returns (receipt, transaction, true) if confirmed or (nil, nil, false) otherwise
func (ob *Observer) checkConfirmedTx(
	ctx context.Context,
	txHash string,
	nonce uint64,
) (*ethtypes.Receipt, *ethtypes.Transaction, bool) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// prepare logger
	logger := ob.Logger().Outbound.With().
		Str(logs.FieldMethod, "checkConfirmedTx").
		Int64(logs.FieldChain, ob.Chain().ChainId).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txHash).
		Logger()

	// query transaction
	transaction, isPending, err := ob.evmClient.TransactionByHash(ctx, ethcommon.HexToHash(txHash))
	if err != nil {
		logger.Error().Err(err).Msg("TransactionByHash error")
		return nil, nil, false
	}
	if transaction == nil { // should not happen
		logger.Error().Msg("transaction is nil")
		return nil, nil, false
	}
	if isPending {
		// should not happen when we are here. The outbound tracker reporter won't report a pending tx.
		logger.Error().Msg("transaction is pending")
		return nil, nil, false
	}

	// check tx sender and nonce
	signer := ethtypes.NewLondonSigner(big.NewInt(ob.Chain().ChainId))
	from, err := signer.Sender(transaction)
	switch {
	case err != nil:
		logger.Error().Err(err).Msg("local recovery of sender address failed")
		return nil, nil, false
	case from != ob.TSS().PubKey().AddressEVM():
		// might be false positive during TSS upgrade for unconfirmed txs
		// Make sure all deposits/withdrawals are paused during TSS upgrade
		logger.Error().Str("tx.sender", from.String()).Msgf("tx sender is not TSS addresses")
		return nil, nil, false
	case transaction.Nonce() != nonce:
		logger.Error().
			Uint64("tx.nonce", transaction.Nonce()).
			Uint64("tracker.nonce", nonce).
			Msg("tx nonce is not matching tracker nonce")
		return nil, nil, false
	}

	// query receipt
	receipt, err := ob.evmClient.TransactionReceipt(ctx, ethcommon.HexToHash(txHash))
	if err != nil {
		logger.Error().Err(err).Msg("TransactionReceipt error")
		return nil, nil, false
	}
	if receipt == nil { // should not happen
		logger.Error().Msg("receipt is nil")
		return nil, nil, false
	}

	// check confirmations
	txBlock := receipt.BlockNumber.Uint64()
	if !ob.IsBlockConfirmedForOutboundSafe(txBlock) {
		logger.Debug().Uint64("tx_block", txBlock).Uint64("last_block", ob.LastBlock()).Msg("tx not confirmed yet")
		return nil, nil, false
	}

	// cross-check tx inclusion against the block
	// Note: a guard for false BlockNumber in receipt. The blob-carrying tx won't come here
	err = ob.checkTxInclusion(ctx, transaction, receipt)
	if err != nil {
		logger.Error().Err(err).Msg("CheckTxInclusion error")
		return nil, nil, false
	}

	return receipt, transaction, true
}
