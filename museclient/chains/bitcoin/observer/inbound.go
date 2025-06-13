package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	cosmosmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/common"
	"github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/musecore"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/memo"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// ObserveInbound observes the Bitcoin chain for inbounds and post votes to musecore
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	logger := ob.Logger().Inbound.With().Str(logs.FieldMethod, "observe_inbound").Logger()

	// keep last block up-to-date
	if err := ob.updateLastBlock(ctx); err != nil {
		return err
	}

	// scan SAFE confirmed blocks
	startBlockSafe, endBlockSafe := ob.GetScanRangeInboundSafe(config.MaxBlocksPerScan)
	if startBlockSafe < endBlockSafe {
		// observe inbounds for the block range [startBlock, endBlock-1]
		lastScannedNew, err := ob.observeInboundInBlockRange(ctx, startBlockSafe, endBlockSafe-1)
		if err != nil {
			logger.Error().
				Err(err).
				Uint64("from", startBlockSafe).
				Uint64("to", endBlockSafe-1).
				Msg("error observing inbounds in block range")
		}

		// save last scanned block to both memory and db
		if lastScannedNew > ob.LastBlockScanned() {
			logger.Info().
				Uint64("from", startBlockSafe).
				Uint64("to", lastScannedNew).
				Msg("observed blocks for inbounds")
			if err := ob.SaveLastBlockScanned(lastScannedNew); err != nil {
				return errors.Wrapf(err, "unable to save last scanned Bitcoin block %d", lastScannedNew)
			}
		}
	}

	// scan FAST confirmed blocks if available
	_, endBlockFast := ob.GetScanRangeInboundFast(config.MaxBlocksPerScan)
	if endBlockSafe < endBlockFast {
		_, err := ob.observeInboundInBlockRange(ctx, endBlockSafe, endBlockFast-1)
		if err != nil {
			logger.Error().
				Err(err).
				Uint64("from", endBlockSafe).
				Uint64("to", endBlockFast-1).
				Msg("error observing inbounds in block range (fast)")
		}
	}

	return nil
}

// observeInboundInBlockRange observes inbounds for given block range [startBlock, toBlock (inclusive)]
// It returns the last successfully scanned block height, so the caller knows where to resume next time
func (ob *Observer) observeInboundInBlockRange(ctx context.Context, startBlock, toBlock uint64) (uint64, error) {
	for blockNumber := startBlock; blockNumber <= toBlock; blockNumber++ {
		// query incoming gas asset to TSS address
		// #nosec G115 always in range
		res, err := ob.GetBlockByNumberCached(ctx, int64(blockNumber))
		if err != nil {
			// we have to re-scan this block next time
			return blockNumber - 1, errors.Wrapf(err, "error getting bitcoin block %d", blockNumber)
		}

		if len(res.Block.Tx) <= 1 {
			continue
		}

		// filter incoming txs to TSS address
		tssAddress := ob.TSSAddressString()

		// #nosec G115 always positive
		events, err := FilterAndParseIncomingTx(
			ctx,
			ob.rpc,
			res.Block.Tx,
			uint64(res.Block.Height),
			tssAddress,
			ob.logger.Inbound,
			ob.netParams,
		)
		if err != nil {
			// we have to re-scan this block next time
			return blockNumber - 1, errors.Wrapf(err, "error filtering incoming txs for block %d", blockNumber)
		}

		// post inbound vote message to musecore
		for _, event := range events {
			msg := ob.GetInboundVoteFromBtcEvent(event)
			if msg != nil {
				// skip early observed inbound that is not eligible for fast confirmation
				if msg.ConfirmationMode == crosschaintypes.ConfirmationMode_FAST {
					eligible, err := ob.IsInboundEligibleForFastConfirmation(ctx, msg)
					if err != nil {
						return blockNumber - 1, errors.Wrapf(
							err,
							"unable to determine inbound fast confirmation eligibility for tx %s",
							event.TxHash,
						)
					}
					if !eligible {
						continue
					}
				}

				_, err = ob.PostVoteInbound(ctx, msg, musecore.PostVoteInboundExecutionGasLimit)
				if err != nil {
					// we have to re-scan this block next time
					return blockNumber - 1, errors.Wrapf(err, "error posting inbound vote for tx %s", event.TxHash)
				}
			}
		}
	}

	// successful processed all blocks in [startBlock, toBlock]
	return toBlock, nil
}

// FilterAndParseIncomingTx given txs list returned by the "getblock 2" RPC command, return the txs that are relevant to us
// relevant tx must have the following vouts as the first two vouts:
// vout0: p2wpkh to the TSS address (targetAddress)
// vout1: OP_RETURN memo, base64 encoded
func FilterAndParseIncomingTx(
	ctx context.Context,
	rpc RPC,
	txs []btcjson.TxRawResult,
	blockNumber uint64,
	tssAddress string,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
) ([]*BTCInboundEvent, error) {
	events := make([]*BTCInboundEvent, 0)

	for idx, tx := range txs {
		if idx == 0 {
			// the first tx is coinbase; we do not process coinbase tx
			continue
		}

		event, err := GetBtcEventWithWitness(
			ctx,
			rpc,
			tx,
			tssAddress,
			blockNumber,
			logger,
			netParams,
			common.CalcDepositorFee,
		)
		if err != nil {
			// unable to parse the tx, the caller should retry
			return nil, errors.Wrapf(err, "error getting btc event for tx %s in block %d", tx.Txid, blockNumber)
		}

		if event != nil {
			events = append(events, event)
		}
	}

	return events, nil
}

// GetInboundVoteFromBtcEvent converts a BTCInboundEvent to a MsgVoteInbound to enable voting on the inbound on musecore
//
// Returns:
//   - a valid MsgVoteInbound message, or
//   - nil if no valid message can be created for whatever reasons:
//     invalid data, not processable, invalid amount, etc.
func (ob *Observer) GetInboundVoteFromBtcEvent(event *BTCInboundEvent) *crosschaintypes.MsgVoteInbound {
	// prepare logger fields
	lf := map[string]any{
		logs.FieldMethod: "GetInboundVoteFromBtcEvent",
		logs.FieldTx:     event.TxHash,
	}

	// decode event memo bytes
	// if the memo is invalid, we set the status in the event, the inbound will be observed as invalid
	err := event.DecodeMemoBytes(ob.Chain().ChainId)
	if err != nil {
		ob.Logger().Inbound.Info().Fields(lf).Msgf("invalid memo bytes: %s", hex.EncodeToString(event.MemoBytes))
		event.Status = crosschaintypes.InboundStatus_INVALID_MEMO
	}

	// check if the event is processable
	if !ob.IsEventProcessable(*event) {
		return nil
	}

	// convert the amount to integer (satoshis)
	amountSats, err := common.GetSatoshis(event.Value)
	if err != nil {
		ob.Logger().Inbound.Error().Err(err).Fields(lf).Msgf("can't convert value %f to satoshis", event.Value)
		return nil
	}
	amountInt := big.NewInt(amountSats)

	// create inbound vote message contract V1 for legacy memo
	if event.MemoStd == nil {
		return ob.NewInboundVoteFromLegacyMemo(event, amountInt)
	}

	// create inbound vote message for standard memo
	return ob.NewInboundVoteFromStdMemo(event, amountInt)
}

// GetSenderAddressByVin get the sender address from the transaction input (vin)
func GetSenderAddressByVin(
	ctx context.Context,
	rpc RPC,
	vin btcjson.Vin,
	net *chaincfg.Params,
) (string, error) {
	// query previous raw transaction by txid
	hash, err := chainhash.NewHashFromStr(vin.Txid)
	if err != nil {
		return "", err
	}

	// this requires running bitcoin node with 'txindex=1'
	tx, err := rpc.GetRawTransaction(ctx, hash)
	if err != nil {
		return "", errors.Wrapf(err, "error getting raw transaction %s", vin.Txid)
	}

	// #nosec G115 - always in range
	if len(tx.MsgTx().TxOut) <= int(vin.Vout) {
		return "", fmt.Errorf("vout index %d out of range for tx %s", vin.Vout, vin.Txid)
	}

	// decode sender address from previous pkScript
	pkScript := tx.MsgTx().TxOut[vin.Vout].PkScript

	return common.DecodeSenderFromScript(pkScript, net)
}

// NewInboundVoteFromLegacyMemo creates a MsgVoteInbound message for inbound that uses legacy memo
func (ob *Observer) NewInboundVoteFromLegacyMemo(
	event *BTCInboundEvent,
	amountSats *big.Int,
) *crosschaintypes.MsgVoteInbound {
	// determine confirmation mode
	confirmationMode := crosschaintypes.ConfirmationMode_FAST
	if ob.IsBlockConfirmedForInboundSafe(event.BlockNumber) {
		confirmationMode = crosschaintypes.ConfirmationMode_SAFE
	}

	return crosschaintypes.NewMsgVoteInbound(
		ob.MusecoreClient().GetKeys().GetOperatorAddress().String(),
		event.FromAddress,
		ob.Chain().ChainId,
		event.FromAddress,
		event.ToAddress,
		ob.MusecoreClient().Chain().ChainId,
		cosmosmath.NewUintFromBigInt(amountSats),
		hex.EncodeToString(event.MemoBytes),
		event.TxHash,
		event.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // no arbitrary call for deposit to MuseChain
		event.Status,
		confirmationMode,
		crosschaintypes.WithCrossChainCall(len(event.MemoBytes) > 0),
	)
}

// NewInboundVoteFromStdMemo creates a MsgVoteInbound message for inbound that uses standard memo
func (ob *Observer) NewInboundVoteFromStdMemo(
	event *BTCInboundEvent,
	amountSats *big.Int,
) *crosschaintypes.MsgVoteInbound {
	// inject the 'revertAddress' specified in the memo, so that
	// musecore will create a revert outbound that points to the custom revert address.
	revertOptions := crosschaintypes.RevertOptions{
		RevertAddress: event.MemoStd.RevertOptions.RevertAddress,
		AbortAddress:  event.MemoStd.RevertOptions.AbortAddress,
	}

	// check if the memo is a cross-chain call, or simple token deposit
	isCrosschainCall := event.MemoStd.OpCode == memo.OpCodeCall || event.MemoStd.OpCode == memo.OpCodeDepositAndCall

	// determine confirmation mode
	confirmationMode := crosschaintypes.ConfirmationMode_FAST
	if ob.IsBlockConfirmedForInboundSafe(event.BlockNumber) {
		confirmationMode = crosschaintypes.ConfirmationMode_SAFE
	}

	return crosschaintypes.NewMsgVoteInbound(
		ob.MusecoreClient().GetKeys().GetOperatorAddress().String(),
		event.FromAddress,
		ob.Chain().ChainId,
		event.FromAddress,
		event.MemoStd.Receiver.Hex(),
		ob.MusecoreClient().Chain().ChainId,
		cosmosmath.NewUintFromBigInt(amountSats),
		hex.EncodeToString(event.MemoStd.Payload),
		event.TxHash,
		event.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // no arbitrary call for deposit to MuseChain
		event.Status,
		confirmationMode,
		crosschaintypes.WithRevertOptions(revertOptions),
		crosschaintypes.WithCrossChainCall(isCrosschainCall),
	)
}
