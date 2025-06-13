package observer

import (
	"bytes"
	"context"
	"encoding/hex"
	"sort"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewayevm.sol"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/RWAs-labs/muse/museclient/chains/evm/common"
	"github.com/RWAs-labs/muse/museclient/compliance"
	"github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/musecore"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/pkg/crypto"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// isEventProcessable checks if the event is processable
func (ob *Observer) isEventProcessable(
	sender, receiver ethcommon.Address,
	txHash ethcommon.Hash,
	payload []byte,
) bool {
	// compliance check
	if config.ContainRestrictedAddress(sender.Hex(), receiver.Hex()) {
		compliance.PrintComplianceLog(
			ob.Logger().Inbound,
			ob.Logger().Compliance,
			false,
			ob.Chain().ChainId,
			txHash.Hex(),
			sender.Hex(),
			receiver.Hex(),
			"Deposit",
		)
		return false
	}

	// donation check
	if bytes.Equal(payload, []byte(constant.DonationMessage)) {
		logFields := map[string]any{
			"chain": ob.Chain().ChainId,
			"tx":    txHash.Hex(),
		}
		ob.Logger().Inbound.Info().Fields(logFields).
			Msgf("thank you rich folk for your donation!")
		return false
	}

	return true
}

// observeGatewayDeposit queries the gateway contract for deposit events
// returns the last block successfully scanned
func (ob *Observer) observeGatewayDeposit(
	ctx context.Context,
	startBlock, toBlock uint64,
	rawLogs []ethtypes.Log,
) (uint64, error) {
	// filter ERC20CustodyDeposited logs
	gatewayAddr, gatewayContract, err := ob.getGatewayContract()
	if err != nil {
		// lastScanned is startBlock - 1
		return startBlock - 1, errors.Wrap(err, "can't get gateway contract")
	}

	// parse and validate events
	events := ob.parseAndValidateDepositEvents(rawLogs, gatewayAddr, gatewayContract)

	// post to musecore
	lastScanned := uint64(0)
	for _, event := range events {
		// remember which block we are scanning (there could be multiple events in the same block)
		if event.Raw.BlockNumber > lastScanned {
			lastScanned = event.Raw.BlockNumber
		}

		// check if the event is processable
		if !ob.isEventProcessable(event.Sender, event.Receiver, event.Raw.TxHash, event.Payload) {
			continue
		}

		msg := ob.newDepositInboundVote(event)

		// skip early observed inbound that is not eligible for fast confirmation
		if msg.ConfirmationMode == types.ConfirmationMode_FAST {
			eligible, err := ob.IsInboundEligibleForFastConfirmation(ctx, &msg)
			if err != nil {
				return lastScanned - 1, errors.Wrapf(
					err,
					"unable to determine inbound fast confirmation eligibility for tx %s",
					event.Raw.TxHash,
				)
			}
			if !eligible {
				continue
			}
		}

		_, err = ob.PostVoteInbound(ctx, &msg, musecore.PostVoteInboundExecutionGasLimit)
		if err != nil {
			// decrement the last scanned block so we have to re-scan from this block next time
			return lastScanned - 1, errors.Wrap(err, "error posting vote inbound")
		}
	}

	// successfully processed all events in [startBlock, toBlock]
	return toBlock, nil
}

// parseAndValidateDepositEvents collects and sorts events by block number, tx index, and log index
func (ob *Observer) parseAndValidateDepositEvents(
	rawLogs []ethtypes.Log,
	gatewayAddr ethcommon.Address,
	gatewayContract *gatewayevm.GatewayEVM,
) []*gatewayevm.GatewayEVMDeposited {
	validEvents := make([]*gatewayevm.GatewayEVMDeposited, 0)
	for _, log := range rawLogs {
		err := common.ValidateEvmTxLog(&log, gatewayAddr, "", common.TopicsGatewayDeposit)
		if err != nil {
			continue
		}
		depositedEvent, err := gatewayContract.ParseDeposited(log)
		if err != nil {
			ob.Logger().
				Inbound.Warn().
				Stringer(logs.FieldTx, log.TxHash).
				Uint64(logs.FieldBlock, log.BlockNumber).
				Msg("invalid Deposited event")
			continue
		}
		validEvents = append(validEvents, depositedEvent)
	}

	// order events by height, tx index and event index (ascending)
	// this ensures the first event is observed if there are multiple in the same tx
	sort.SliceStable(validEvents, func(i, j int) bool {
		if validEvents[i].Raw.BlockNumber == validEvents[j].Raw.BlockNumber {
			if validEvents[i].Raw.TxIndex == validEvents[j].Raw.TxIndex {
				return validEvents[i].Raw.Index < validEvents[j].Raw.Index
			}
			return validEvents[i].Raw.TxIndex < validEvents[j].Raw.TxIndex
		}
		return validEvents[i].Raw.BlockNumber < validEvents[j].Raw.BlockNumber
	})

	// filter events from same tx
	filtered := make([]*gatewayevm.GatewayEVMDeposited, 0)
	guard := make(map[string]bool)
	for _, event := range validEvents {
		// guard against multiple events in the same tx
		if guard[event.Raw.TxHash.Hex()] {
			ob.Logger().
				Inbound.Warn().
				Stringer(logs.FieldTx, event.Raw.TxHash).
				Msg("multiple Deposited events in same tx")
			continue
		}
		guard[event.Raw.TxHash.Hex()] = true
		filtered = append(filtered, event)
	}

	return filtered
}

// newDepositInboundVote creates a MsgVoteInbound message for a Deposit event
func (ob *Observer) newDepositInboundVote(event *gatewayevm.GatewayEVMDeposited) types.MsgVoteInbound {
	// if event.Asset is zero, it's a native token
	coinType := coin.CoinType_ERC20
	if crypto.IsEmptyAddress(event.Asset) {
		coinType = coin.CoinType_Gas
	}

	// to maintain compatibility with previous gateway version, deposit event with a non-empty payload is considered as a call
	isCrossChainCall := false
	if len(event.Payload) > 0 {
		isCrossChainCall = true
	}

	// determine confirmation mode
	confirmationMode := types.ConfirmationMode_FAST
	if ob.IsBlockConfirmedForInboundSafe(event.Raw.BlockNumber) {
		confirmationMode = types.ConfirmationMode_SAFE
	}

	return *types.NewMsgVoteInbound(
		ob.MusecoreClient().GetKeys().GetOperatorAddress().String(),
		event.Sender.Hex(),
		ob.Chain().ChainId,
		"",
		event.Receiver.Hex(),
		ob.MusecoreClient().Chain().ChainId,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Payload),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		musecore.PostVoteInboundCallOptionsGasLimit,
		coinType,
		event.Asset.Hex(),
		uint64(event.Raw.Index),
		types.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		types.InboundStatus_SUCCESS,
		confirmationMode,
		types.WithEVMRevertOptions(event.RevertOptions),
		types.WithCrossChainCall(isCrossChainCall),
	)
}

// observeGatewayCall queries the gateway contract for call events
// returns the last block successfully scanned
// TODO: there are lot of similarities between this function and ObserveGatewayDeposit
// logic should be factorized using interfaces and generics
// https://github.com/RWAs-labs/muse/issues/2493
func (ob *Observer) observeGatewayCall(
	ctx context.Context,
	startBlock, toBlock uint64,
	rawLogs []ethtypes.Log,
) (uint64, error) {
	gatewayAddr, gatewayContract, err := ob.getGatewayContract()
	if err != nil {
		// lastScanned is startBlock - 1
		return startBlock - 1, errors.Wrap(err, "can't get gateway contract")
	}

	events := ob.parseAndValidateCallEvents(rawLogs, gatewayAddr, gatewayContract)
	lastScanned := uint64(0)
	for _, event := range events {
		if event.Raw.BlockNumber > lastScanned {
			lastScanned = event.Raw.BlockNumber
		}

		if !ob.isEventProcessable(event.Sender, event.Receiver, event.Raw.TxHash, event.Payload) {
			continue
		}

		msg := ob.newCallInboundVote(event)

		ob.Logger().Inbound.Info().
			Msgf("ObserveGateway: Call inbound detected on chain %d tx %s block %d from %s value message %s",
				ob.Chain().
					ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, event.Sender.Hex(), hex.EncodeToString(event.Payload))

		_, err = ob.PostVoteInbound(ctx, &msg, musecore.PostVoteInboundExecutionGasLimit)
		if err != nil {
			return lastScanned - 1, errors.Wrap(err, "error posting vote inbound")
		}
	}

	return toBlock, nil
}

// parseAndValidateCallEvents collects and sorts events by block number, tx index, and log index
func (ob *Observer) parseAndValidateCallEvents(
	rawLogs []ethtypes.Log,
	gatewayAddr ethcommon.Address,
	gatewayContract *gatewayevm.GatewayEVM,
) []*gatewayevm.GatewayEVMCalled {
	validEvents := make([]*gatewayevm.GatewayEVMCalled, 0)
	for _, log := range rawLogs {
		err := common.ValidateEvmTxLog(&log, gatewayAddr, "", common.TopicsGatewayCall)
		if err != nil {
			continue
		}
		calledEvent, err := gatewayContract.ParseCalled(log)
		if err != nil {
			ob.Logger().
				Inbound.Warn().
				Stringer(logs.FieldTx, log.TxHash).
				Uint64(logs.FieldBlock, log.BlockNumber).
				Msg("invalid Called event")
			continue
		}
		validEvents = append(validEvents, calledEvent)
	}

	// order events by height, tx index and event index (ascending)
	// this ensures the first event is observed if there are multiple in the same tx
	sort.SliceStable(validEvents, func(i, j int) bool {
		if validEvents[i].Raw.BlockNumber == validEvents[j].Raw.BlockNumber {
			if validEvents[i].Raw.TxIndex == validEvents[j].Raw.TxIndex {
				return validEvents[i].Raw.Index < validEvents[j].Raw.Index
			}
			return validEvents[i].Raw.TxIndex < validEvents[j].Raw.TxIndex
		}
		return validEvents[i].Raw.BlockNumber < validEvents[j].Raw.BlockNumber
	})

	// filter events from same tx
	filtered := make([]*gatewayevm.GatewayEVMCalled, 0)
	guard := make(map[string]bool)
	for _, event := range validEvents {
		// guard against multiple events in the same tx
		if guard[event.Raw.TxHash.Hex()] {
			ob.Logger().Inbound.Warn().Stringer(logs.FieldTx, event.Raw.TxHash).Msg("multiple Called events in same tx")
			continue
		}
		guard[event.Raw.TxHash.Hex()] = true
		filtered = append(filtered, event)
	}

	return filtered
}

// newCallInboundVote creates a MsgVoteInbound message for a Call event
func (ob *Observer) newCallInboundVote(event *gatewayevm.GatewayEVMCalled) types.MsgVoteInbound {
	return *types.NewMsgVoteInbound(
		ob.MusecoreClient().GetKeys().GetOperatorAddress().String(),
		event.Sender.Hex(),
		ob.Chain().ChainId,
		"",
		event.Receiver.Hex(),
		ob.MusecoreClient().Chain().ChainId,
		sdkmath.ZeroUint(),
		hex.EncodeToString(event.Payload),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		musecore.PostVoteInboundCallOptionsGasLimit,
		coin.CoinType_NoAssetCall,
		"",
		uint64(event.Raw.Index),
		types.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		types.InboundStatus_SUCCESS,
		types.ConfirmationMode_SAFE,
		types.WithEVMRevertOptions(event.RevertOptions),
	)
}

// observeGatewayDepositAndCall queries the gateway contract for deposit and call events
// returns the last block successfully scanned
func (ob *Observer) observeGatewayDepositAndCall(
	ctx context.Context,
	startBlock, toBlock uint64,
	rawLogs []ethtypes.Log,
) (uint64, error) {
	gatewayAddr, gatewayContract, err := ob.getGatewayContract()
	if err != nil {
		// lastScanned is startBlock - 1
		return startBlock - 1, errors.Wrap(err, "can't get gateway contract")
	}

	events := ob.parseAndValidateDepositAndCallEvents(rawLogs, gatewayAddr, gatewayContract)

	lastScanned := uint64(0)
	for _, event := range events {
		// remember which block we are scanning (there could be multiple events in the same block)
		if event.Raw.BlockNumber > lastScanned {
			lastScanned = event.Raw.BlockNumber
		}

		// check if the event is processable
		if !ob.isEventProcessable(event.Sender, event.Receiver, event.Raw.TxHash, event.Payload) {
			continue
		}

		msg := ob.newDepositAndCallInboundVote(event)

		ob.Logger().Inbound.Info().
			Msgf("ObserveGateway: DepositAndCall inbound detected on chain %d tx %s block %d from %s value %s message %s",
				ob.Chain().
					ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, event.Sender.Hex(), event.Amount.String(), hex.EncodeToString(event.Payload))

		_, err = ob.PostVoteInbound(ctx, &msg, musecore.PostVoteInboundExecutionGasLimit)
		if err != nil {
			// decrement the last scanned block so we have to re-scan from this block next time
			return lastScanned - 1, errors.Wrap(err, "error posting vote inbound")
		}
	}

	// successfully processed all events in [startBlock, toBlock]
	return toBlock, nil
}

// parseAndValidateDepositAndCallEvents collects and sorts events by block number, tx index, and log index
func (ob *Observer) parseAndValidateDepositAndCallEvents(
	rawLogs []ethtypes.Log,
	gatewayAddr ethcommon.Address,
	gatewayContract *gatewayevm.GatewayEVM,
) []*gatewayevm.GatewayEVMDepositedAndCalled {
	// collect and sort validEvents by block number, then tx index, then log index (ascending)
	validEvents := make([]*gatewayevm.GatewayEVMDepositedAndCalled, 0)
	for _, log := range rawLogs {
		err := common.ValidateEvmTxLog(&log, gatewayAddr, "", common.TopicsGatewayDepositAndCall)
		if err != nil {
			continue
		}
		depositAndCallEvent, err := gatewayContract.ParseDepositedAndCalled(log)
		if err != nil {
			ob.Logger().
				Inbound.Warn().
				Stringer(logs.FieldTx, log.TxHash).
				Uint64(logs.FieldBlock, log.BlockNumber).
				Msg("invalid DepositedAndCalled event")
			continue
		}
		validEvents = append(validEvents, depositAndCallEvent)
	}

	// order events by height, tx index and event index (ascending)
	// this ensures the first event is observed if there are multiple in the same tx
	sort.SliceStable(validEvents, func(i, j int) bool {
		if validEvents[i].Raw.BlockNumber == validEvents[j].Raw.BlockNumber {
			if validEvents[i].Raw.TxIndex == validEvents[j].Raw.TxIndex {
				return validEvents[i].Raw.Index < validEvents[j].Raw.Index
			}
			return validEvents[i].Raw.TxIndex < validEvents[j].Raw.TxIndex
		}
		return validEvents[i].Raw.BlockNumber < validEvents[j].Raw.BlockNumber
	})

	// filter events from same tx
	filtered := make([]*gatewayevm.GatewayEVMDepositedAndCalled, 0)
	guard := make(map[string]bool)
	for _, event := range validEvents {
		// guard against multiple events in the same tx
		if guard[event.Raw.TxHash.Hex()] {
			ob.Logger().
				Inbound.Warn().
				Stringer(logs.FieldTx, event.Raw.TxHash).
				Msg("multiple DepositedAndCalled events in same tx")
			continue
		}
		guard[event.Raw.TxHash.Hex()] = true
		filtered = append(filtered, event)
	}

	return filtered
}

// newDepositAndCallInboundVote creates a MsgVoteInbound message for a Deposit event
func (ob *Observer) newDepositAndCallInboundVote(event *gatewayevm.GatewayEVMDepositedAndCalled) types.MsgVoteInbound {
	// if event.Asset is zero, it's a native token
	coinType := coin.CoinType_ERC20
	if crypto.IsEmptyAddress(event.Asset) {
		coinType = coin.CoinType_Gas
	}

	return *types.NewMsgVoteInbound(
		ob.MusecoreClient().GetKeys().GetOperatorAddress().String(),
		event.Sender.Hex(),
		ob.Chain().ChainId,
		"",
		event.Receiver.Hex(),
		ob.MusecoreClient().Chain().ChainId,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Payload),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		coinType,
		event.Asset.Hex(),
		uint64(event.Raw.Index),
		types.ProtocolContractVersion_V2,
		false, // currently not relevant since calls are not arbitrary
		types.InboundStatus_SUCCESS,
		types.ConfirmationMode_SAFE,
		types.WithEVMRevertOptions(event.RevertOptions),
		types.WithCrossChainCall(true),
	)
}
