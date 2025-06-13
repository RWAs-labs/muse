package observer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"slices"
	"sort"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/protocol-contracts/pkg/erc20custody.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/museconnector.non-eth.sol"
	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/chains/evm/client"
	"github.com/RWAs-labs/muse/museclient/chains/evm/common"
	"github.com/RWAs-labs/muse/museclient/compliance"
	"github.com/RWAs-labs/muse/museclient/config"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/metrics"
	"github.com/RWAs-labs/muse/museclient/musecore"
	clienttypes "github.com/RWAs-labs/muse/museclient/types"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/pkg/memo"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// ProcessInboundTrackers observes inbound trackers from musecore
func (ob *Observer) ProcessInboundTrackers(ctx context.Context) error {
	trackers, err := ob.MusecoreClient().GetInboundTrackersForChain(ctx, ob.Chain().ChainId)
	if err != nil {
		return err
	}

	for _, tracker := range trackers {
		// query tx and receipt
		tx, _, err := ob.transactionByHash(ctx, tracker.TxHash)
		if err != nil {
			return errors.Wrapf(
				err,
				"error getting transaction for inbound %s chain %d",
				tracker.TxHash,
				ob.Chain().ChainId,
			)
		}

		receipt, err := ob.evmClient.TransactionReceipt(ctx, ethcommon.HexToHash(tracker.TxHash))
		if err != nil {
			return errors.Wrapf(
				err,
				"error getting receipt for inbound %s chain %d",
				tracker.TxHash,
				ob.Chain().ChainId,
			)
		}
		ob.Logger().Inbound.Info().Msgf("checking tracker for inbound %s chain %d", tracker.TxHash, ob.Chain().ChainId)

		// try processing the tracker for v2 inbound
		// filter error if event is not found, in this case we run v1 tracker process
		if err := ob.ProcessInboundTrackerV2(ctx, tx, receipt); err != nil &&
			!errors.Is(err, ErrEventNotFound) && !errors.Is(err, ErrGatewayNotSet) {
			return err
		} else if err == nil {
			// continue with next tracker
			continue
		}

		// try processing the tracker for v1 inbound
		switch tracker.CoinType {
		case coin.CoinType_Muse:
			_, err = ob.checkAndVoteInboundTokenMuse(ctx, tx, receipt, true)
		case coin.CoinType_ERC20:
			_, err = ob.checkAndVoteInboundTokenERC20(ctx, tx, receipt, true)
		case coin.CoinType_Gas:
			_, err = ob.checkAndVoteInboundTokenGas(ctx, tx, receipt, true)
		default:
			return fmt.Errorf(
				"unknown coin type %s for inbound %s chain %d",
				tracker.CoinType,
				tx.Hash,
				ob.Chain().ChainId,
			)
		}
		if err != nil {
			return errors.Wrapf(err, "error checking and voting for inbound %s chain %d", tx.Hash, ob.Chain().ChainId)
		}
	}
	return nil
}

// ObserveInbound observes the evm chain for inbounds and posts votes to musecore
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	logger := ob.Logger().Inbound.With().Str(logs.FieldMethod, "observe_inbound").Logger()

	// keep last block up-to-date
	if err := ob.updateLastBlock(ctx); err != nil {
		return err
	}

	// uncomment this line to stop observing inbound and test observation with inbound trackers
	// https://github.com/RWAs-labs/muse/blob/3879b5ef8b418542c82a4383263604222f0605c6/e2e/e2etests/test_inbound_trackers.go#L19
	// TODO: implement a better way to disable inbound observation
	// https://github.com/RWAs-labs/muse/issues/3186
	//return nil

	// scan SAFE confirmed blocks
	startBlockSafe, endBlockSafe := ob.GetScanRangeInboundSafe(config.MaxBlocksPerScan)
	if startBlockSafe < endBlockSafe {
		// observe inbounds in block range [startBlock, endBlock-1]
		lastScannedNew := ob.observeInboundInBlockRange(ctx, startBlockSafe, endBlockSafe-1)

		// save last scanned block to both memory and db
		if lastScannedNew > ob.LastBlockScanned() {
			logger.Debug().
				Uint64("from", startBlockSafe).
				Uint64("to", lastScannedNew).
				Msg("observed blocks for inbounds")
			if err := ob.SaveLastBlockScanned(lastScannedNew); err != nil {
				return errors.Wrapf(err, "unable to save last scanned block %d", lastScannedNew)
			}
		}
	}

	// scan FAST confirmed blocks if available
	_, endBlockFast := ob.GetScanRangeInboundFast(config.MaxBlocksPerScan)
	if endBlockSafe < endBlockFast {
		ob.observeInboundInBlockRange(ctx, endBlockSafe, endBlockFast-1)
	}

	return nil
}

// observeInboundInBlockRange observes inbounds for given block range [startBlock, toBlock (inclusive)]
// It returns the last successfully scanned block height, so the caller knows where to resume next time
func (ob *Observer) observeInboundInBlockRange(ctx context.Context, startBlock, toBlock uint64) uint64 {
	logger := ob.Logger().Inbound.With().
		Str(logs.FieldMethod, "observeInboundInBlockRange").
		Uint64("start_block", startBlock).Uint64("to_block", toBlock).Logger()

	var (
		lastScannedTssRecvd              = toBlock
		lastScannedMuseSent              = startBlock - 1
		lastScannedDeposited             = startBlock - 1
		lastScannedGatewayDeposit        = startBlock - 1
		lastScannedGatewayCall           = startBlock - 1
		lastScannedGatewayDepositAndCall = startBlock - 1
		err                              error
	)

	// we now only take these actions on specific configurable chains
	if !ob.ChainParams().DisableTssBlockScan {
		// query the incoming tx to TSS address (read at most 100 blocks in one go)
		lastScannedTssRecvd, err = ob.observeTSSReceive(ctx, startBlock, toBlock)
		if err != nil {
			logger.Error().Err(err).Msg("error observing TSS received gas asset")
		}

		// filter the outbounds from TSS address to supplement outbound trackers
		// TODO: make this a separate go routine in outbound.go after switching to smart contract V2
		ob.filterTSSOutbound(ctx, startBlock, toBlock)
	}

	logs, err := ob.fetchLogs(ctx, startBlock, toBlock)
	if err != nil {
		ob.Logger().Inbound.Error().Err(err).Msg("get gateway logs")
	} else {
		// handle connector contract deposit
		lastScannedMuseSent, err = ob.observeMuseSent(ctx, startBlock, toBlock, logs)
		if err != nil {
			logger.Error().Err(err).Msg("error observing muse sent events from MuseConnector contract")
		}

		// handle legacy erc20 direct deposit logs
		lastScannedDeposited, err = ob.observeERC20Deposited(ctx, startBlock, toBlock, logs)
		if err != nil {
			logger.Error().Err(err).Msg("error observing deposited events from ERC20Custody contract")
		}

		lastScannedGatewayDeposit, err = ob.observeGatewayDeposit(ctx, startBlock, toBlock, logs)
		if err != nil {
			ob.Logger().Inbound.Error().Err(err).Msg("error observing deposit events from Gateway contract")
		}
		lastScannedGatewayCall, err = ob.observeGatewayCall(ctx, startBlock, toBlock, logs)
		if err != nil {
			ob.Logger().Inbound.Error().Err(err).Msg("error observing call events from Gateway contract")
		}
		lastScannedGatewayDepositAndCall, err = ob.observeGatewayDepositAndCall(ctx, startBlock, toBlock, logs)
		if err != nil {
			ob.Logger().Inbound.Error().Err(err).Msg("error observing depositAndCall events from Gateway contract")
		}
	}

	// note: using the lowest height for all events is not perfect,
	// but it's simple and good enough
	lowestLastScannedBlock := slices.Min([]uint64{
		lastScannedMuseSent,
		lastScannedDeposited,
		lastScannedTssRecvd,
		lastScannedGatewayDeposit,
		lastScannedGatewayCall,
		lastScannedGatewayDepositAndCall,
	})

	return lowestLastScannedBlock
}

func (ob *Observer) fetchLogs(ctx context.Context, startBlock, toBlock uint64) ([]ethtypes.Log, error) {
	gatewayAddr, _, err := ob.getGatewayContract()
	if err != nil {
		return nil, errors.Wrap(err, "can't get gateway contract")
	}

	erc20Addr, _, err := ob.getERC20CustodyContract()
	if err != nil {
		return nil, errors.Wrap(err, "can't get erc20 custody contract")
	}

	connectorAddr, _, err := ob.getConnectorContract()
	if err != nil {
		return nil, errors.Wrap(err, "can't get connector contract")
	}

	addresses := []ethcommon.Address{gatewayAddr, erc20Addr, connectorAddr}

	logs, err := ob.evmClient.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(startBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: addresses,
	})
	if err != nil {
		return nil, errors.Wrap(err, "filter logs")
	}

	// increment prom counter
	metrics.GetFilterLogsPerChain.WithLabelValues(ob.Chain().Name).Inc()

	return logs, nil
}

// observeMuseSent queries the MuseSent event from the connector contract and posts to musecore
// returns the last block successfully scanned
func (ob *Observer) observeMuseSent(
	ctx context.Context,
	startBlock, toBlock uint64,
	logs []ethtypes.Log,
) (uint64, error) {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return 0, err
	}

	// filter MuseSent logs
	addrConnector, connector, err := ob.getConnectorContract()
	if err != nil {
		// we have to re-scan from this block next time
		return startBlock - 1, errors.Wrap(err, "error getting connector contract")
	}

	// collect and sort events by block number, then tx index, then log index (ascending)
	events := make([]*museconnector.MuseConnectorNonEthMuseSent, 0)
	for _, log := range logs {
		// sanity check tx event
		err := common.ValidateEvmTxLog(&log, addrConnector, "", common.TopicsMuseSent)
		if err != nil {
			continue
		}
		event, err := connector.ParseMuseSent(log)
		if err == nil {
			events = append(events, event)
			continue
		}
		ob.Logger().Inbound.Warn().
			Err(err).
			Msgf("ObserveMuseSent: invalid MuseSent event in tx %s on chain %d at height %d",
				log.TxHash.Hex(), ob.Chain().ChainId, log.BlockNumber)
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Raw.BlockNumber == events[j].Raw.BlockNumber {
			if events[i].Raw.TxIndex == events[j].Raw.TxIndex {
				return events[i].Raw.Index < events[j].Raw.Index
			}
			return events[i].Raw.TxIndex < events[j].Raw.TxIndex
		}
		return events[i].Raw.BlockNumber < events[j].Raw.BlockNumber
	})

	// increment prom counter
	metrics.GetFilterLogsPerChain.WithLabelValues(ob.Chain().Name).Inc()

	// post to musecore
	beingScanned := uint64(0)
	guard := make(map[string]bool)
	for _, event := range events {
		// remember which block we are scanning (there could be multiple events in the same block)
		if event.Raw.BlockNumber > beingScanned {
			beingScanned = event.Raw.BlockNumber
		}
		// guard against multiple events in the same tx
		if guard[event.Raw.TxHash.Hex()] {
			ob.Logger().Inbound.Warn().
				Msgf("ObserveMuseSent: multiple remote call events detected in tx %s", event.Raw.TxHash)
			continue
		}
		guard[event.Raw.TxHash.Hex()] = true

		msg := ob.buildInboundVoteMsgForMuseSentEvent(app, event)
		if msg == nil {
			continue
		}

		const gasLimit = musecore.PostVoteInboundMessagePassingExecutionGasLimit
		if _, err = ob.PostVoteInbound(ctx, msg, gasLimit); err != nil {
			// we have to re-scan from this block next time
			return beingScanned - 1, errors.Wrap(err, "error posting inbound vote")
		}
	}

	// successful processed all events in [startBlock, toBlock]
	return toBlock, nil
}

// observeERC20Deposited queries the ERC20CustodyDeposited event from the ERC20Custody contract and posts to musecore
// returns the last block successfully scanned
func (ob *Observer) observeERC20Deposited(
	ctx context.Context,
	startBlock, toBlock uint64,
	logs []ethtypes.Log,
) (uint64, error) {
	// filter ERC20CustodyDeposited logs
	addrCustody, erc20custodyContract, err := ob.getERC20CustodyContract()
	if err != nil {
		// we have to re-scan from this block next time
		return startBlock - 1, errors.Wrap(err, "error getting ERC20Custody contract")
	}

	// collect and sort events by block number, then tx index, then log index (ascending)
	events := make([]*erc20custody.ERC20CustodyDeposited, 0)
	for _, log := range logs {
		// sanity check tx event
		err := common.ValidateEvmTxLog(&log, addrCustody, "", common.TopicsDeposited)
		if err != nil {
			continue
		}
		event, err := erc20custodyContract.ParseDeposited(log)
		if err == nil {
			events = append(events, event)
			continue
		}
		ob.Logger().Inbound.Warn().
			Err(err).
			Msgf("ObserveERC20Deposited: invalid Deposited event in tx %s on chain %d at height %d",
				log.TxHash.Hex(), ob.Chain().ChainId, log.BlockNumber)
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Raw.BlockNumber == events[j].Raw.BlockNumber {
			if events[i].Raw.TxIndex == events[j].Raw.TxIndex {
				return events[i].Raw.Index < events[j].Raw.Index
			}
			return events[i].Raw.TxIndex < events[j].Raw.TxIndex
		}
		return events[i].Raw.BlockNumber < events[j].Raw.BlockNumber
	})

	// increment prom counter
	metrics.GetFilterLogsPerChain.WithLabelValues(ob.Chain().Name).Inc()

	// post to zeatcore
	guard := make(map[string]bool)
	beingScanned := uint64(0)
	for _, event := range events {
		// remember which block we are scanning (there could be multiple events in the same block)
		if event.Raw.BlockNumber > beingScanned {
			beingScanned = event.Raw.BlockNumber
		}
		tx, _, err := ob.transactionByHash(ctx, event.Raw.TxHash.Hex())
		if err != nil {
			// we have to re-scan from this block next time
			return beingScanned - 1, errors.Wrapf(err, "error getting transaction %s", event.Raw.TxHash.Hex())
		}
		sender := ethcommon.HexToAddress(tx.From)

		// guard against multiple events in the same tx
		if guard[event.Raw.TxHash.Hex()] {
			ob.Logger().Inbound.Warn().
				Msgf("ObserveERC20Deposited: multiple remote call events detected in tx %s", event.Raw.TxHash)
			continue
		}
		guard[event.Raw.TxHash.Hex()] = true

		msg := ob.buildInboundVoteMsgForDepositedEvent(event, sender)
		if msg != nil {
			_, err = ob.PostVoteInbound(ctx, msg, musecore.PostVoteInboundExecutionGasLimit)
			if err != nil {
				// we have to re-scan from this block next time
				return beingScanned - 1, errors.Wrap(err, "error posting inbound vote")
			}
		}
	}
	// successful processed all events in [startBlock, toBlock]
	return toBlock, nil
}

// observeTSSReceive queries the incoming gas asset to TSS address and posts to musecore
// returns the last block successfully scanned
func (ob *Observer) observeTSSReceive(ctx context.Context, startBlock, toBlock uint64) (uint64, error) {
	// query incoming gas asset
	for bn := startBlock; bn <= toBlock; bn++ {
		// observe TSS received gas token in block 'bn'
		err := ob.observeTSSReceiveInBlock(ctx, bn)
		if err != nil {
			// we have to re-scan from this block next time
			return bn - 1, errors.Wrapf(err, "error observing TSS received gas asset in block %d", bn)
		}
	}

	// successful processed all gas asset deposits in [startBlock, toBlock]
	return toBlock, nil
}

// checkAndVoteInboundTokenMuse checks and votes on the given inbound Muse token
func (ob *Observer) checkAndVoteInboundTokenMuse(
	ctx context.Context,
	tx *client.Transaction,
	receipt *ethtypes.Receipt,
	vote bool,
) (string, error) {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return "", err
	}

	// check confirmations
	if !ob.IsBlockConfirmedForInboundSafe(receipt.BlockNumber.Uint64()) {
		return "", fmt.Errorf(
			"inbound %s has not been confirmed yet: receipt block %d",
			tx.Hash,
			receipt.BlockNumber.Uint64(),
		)
	}

	// get muse connector contract
	addrConnector, connector, err := ob.getConnectorContract()
	if err != nil {
		return "", err
	}

	// build inbound vote message and post vote
	var msg *types.MsgVoteInbound
	for _, log := range receipt.Logs {
		event, err := connector.ParseMuseSent(*log)
		if err == nil && event != nil {
			// sanity check tx event
			err = common.ValidateEvmTxLog(&event.Raw, addrConnector, tx.Hash, common.TopicsMuseSent)
			if err == nil {
				msg = ob.buildInboundVoteMsgForMuseSentEvent(app, event)
			} else {
				ob.Logger().Inbound.Error().Err(err).Msgf("CheckEvmTxLog error on inbound %s chain %d", tx.Hash, ob.Chain().ChainId)
				return "", err
			}
			break // only one event is allowed per tx
		}
	}
	if msg == nil {
		// no event, restricted tx, etc.
		ob.Logger().Inbound.Info().Msgf("no MuseSent event found for inbound %s chain %d", tx.Hash, ob.Chain().ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(ctx, msg, musecore.PostVoteInboundMessagePassingExecutionGasLimit)
	}

	return msg.Digest(), nil
}

// checkAndVoteInboundTokenERC20 checks and votes on the given inbound ERC20 token
func (ob *Observer) checkAndVoteInboundTokenERC20(
	ctx context.Context,
	tx *client.Transaction,
	receipt *ethtypes.Receipt,
	vote bool,
) (string, error) {
	// check confirmations
	if !ob.IsBlockConfirmedForInboundSafe(receipt.BlockNumber.Uint64()) {
		return "", fmt.Errorf(
			"inbound %s has not been confirmed yet: receipt block %d",
			tx.Hash,
			receipt.BlockNumber.Uint64(),
		)
	}

	// get erc20 custody contract
	addrCustody, custody, err := ob.getERC20CustodyContract()
	if err != nil {
		return "", err
	}
	sender := ethcommon.HexToAddress(tx.From)

	// build inbound vote message and post vote
	var msg *types.MsgVoteInbound
	for _, log := range receipt.Logs {
		museDeposited, err := custody.ParseDeposited(*log)
		if err == nil && museDeposited != nil {
			// sanity check tx event
			err = common.ValidateEvmTxLog(&museDeposited.Raw, addrCustody, tx.Hash, common.TopicsDeposited)
			if err == nil {
				msg = ob.buildInboundVoteMsgForDepositedEvent(museDeposited, sender)
			} else {
				ob.Logger().Inbound.Error().Err(err).Msgf("CheckEvmTxLog error on inbound %s chain %d", tx.Hash, ob.Chain().ChainId)
				return "", err
			}
			break // only one event is allowed per tx
		}
	}
	if msg == nil {
		// no event, donation, restricted tx, etc.
		ob.Logger().Inbound.Info().Msgf("no Deposited event found for inbound %s chain %d", tx.Hash, ob.Chain().ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(ctx, msg, musecore.PostVoteInboundExecutionGasLimit)
	}

	return msg.Digest(), nil
}

// checkAndVoteInboundTokenGas checks and votes on the given inbound gas token
func (ob *Observer) checkAndVoteInboundTokenGas(
	ctx context.Context,
	tx *client.Transaction,
	receipt *ethtypes.Receipt,
	vote bool,
) (string, error) {
	// check confirmations
	if !ob.IsBlockConfirmedForInboundSafe(receipt.BlockNumber.Uint64()) {
		return "", fmt.Errorf(
			"inbound %s has not been confirmed yet: receipt block %d",
			tx.Hash,
			receipt.BlockNumber.Uint64(),
		)
	}

	// checks receiver and tx status
	if ethcommon.HexToAddress(tx.To) != ob.TSS().PubKey().AddressEVM() {
		return "", fmt.Errorf("tx.To %s is not TSS address", tx.To)
	}
	if receipt.Status != ethtypes.ReceiptStatusSuccessful {
		return "", errors.New("not a successful tx")
	}
	sender := ethcommon.HexToAddress(tx.From)

	// build inbound vote message and post vote
	msg := ob.buildInboundVoteMsgForTokenSentToTSS(tx, sender, receipt.BlockNumber.Uint64())
	if msg == nil {
		// donation, restricted tx, etc.
		ob.Logger().Inbound.Info().Msgf("no vote message built for inbound %s chain %d", tx.Hash, ob.Chain().ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(ctx, msg, musecore.PostVoteInboundExecutionGasLimit)
	}

	return msg.Digest(), nil
}

// buildInboundVoteMsgForDepositedEvent builds a inbound vote message for a Deposited event
func (ob *Observer) buildInboundVoteMsgForDepositedEvent(
	event *erc20custody.ERC20CustodyDeposited,
	sender ethcommon.Address,
) *types.MsgVoteInbound {
	// compliance check
	maybeReceiver := ""
	parsedAddress, _, err := memo.DecodeLegacyMemoHex(hex.EncodeToString(event.Message))
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		maybeReceiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(sender.Hex(), clienttypes.BytesToEthHex(event.Recipient), maybeReceiver) {
		compliance.PrintComplianceLog(
			ob.Logger().Inbound,
			ob.Logger().Compliance,
			false,
			ob.Chain().ChainId,
			event.Raw.TxHash.Hex(),
			sender.Hex(),
			clienttypes.BytesToEthHex(event.Recipient),
			"ERC20",
		)
		return nil
	}

	// donation check
	if bytes.Equal(event.Message, []byte(constant.DonationMessage)) {
		ob.Logger().Inbound.Info().
			Msgf("thank you rich folk for your donation! tx %s chain %d", event.Raw.TxHash.Hex(), ob.Chain().ChainId)
		return nil
	}
	message := hex.EncodeToString(event.Message)
	ob.Logger().Inbound.Info().
		Msgf("ERC20CustodyDeposited inbound detected on chain %d tx %s block %d from %s value %s message %s",
			ob.Chain().
				ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, sender.Hex(), event.Amount.String(), message)

	return musecore.GetInboundVoteMessage(
		sender.Hex(),
		ob.Chain().ChainId,
		"",
		clienttypes.BytesToEthHex(event.Recipient),
		ob.MusecoreClient().Chain().ChainId,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		coin.CoinType_ERC20,
		event.Asset.String(),
		ob.MusecoreClient().GetKeys().GetOperatorAddress().String(),
		uint64(event.Raw.Index),
		types.InboundStatus_SUCCESS,
	)
}

// buildInboundVoteMsgForMuseSentEvent builds a inbound vote message for a MuseSent event
func (ob *Observer) buildInboundVoteMsgForMuseSentEvent(
	appContext *zctx.AppContext,
	event *museconnector.MuseConnectorNonEthMuseSent,
) *types.MsgVoteInbound {
	// note that this is most likely muse chain
	destChain, err := appContext.GetChain(event.DestinationChainId.Int64())
	if err != nil {
		ob.Logger().Inbound.Warn().Err(err).Msgf("chain id %d not supported", event.DestinationChainId.Int64())
		return nil
	}

	destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)

	// compliance check
	sender := event.MuseTxSenderAddress.Hex()
	if config.ContainRestrictedAddress(sender, destAddr, event.SourceTxOriginAddress.Hex()) {
		compliance.PrintComplianceLog(ob.Logger().Inbound, ob.Logger().Compliance,
			false, ob.Chain().ChainId, event.Raw.TxHash.Hex(), sender, destAddr, "Muse")
		return nil
	}

	if !destChain.IsMuse() {
		if strings.EqualFold(destAddr, destChain.Params().MuseTokenContractAddress) {
			ob.Logger().Inbound.Warn().
				Msgf("potential attack attempt: %s destination address is MUSE token contract address", destAddr)
			return nil
		}
	}
	message := base64.StdEncoding.EncodeToString(event.Message)
	ob.Logger().Inbound.Info().Msgf("MuseSent inbound detected on chain %d tx %s block %d from %s value %s message %s",
		ob.Chain().
			ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, sender, event.MuseValueAndGas.String(), message)

	return musecore.GetInboundVoteMessage(
		sender,
		ob.Chain().ChainId,
		event.SourceTxOriginAddress.Hex(),
		destAddr,
		destChain.ID(),
		sdkmath.NewUintFromBigInt(event.MuseValueAndGas),
		message,
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		event.DestinationGasLimit.Uint64(),
		coin.CoinType_Muse,
		"",
		ob.MusecoreClient().GetKeys().GetOperatorAddress().String(),
		uint64(event.Raw.Index),
		types.InboundStatus_SUCCESS,
	)
}

// buildInboundVoteMsgForTokenSentToTSS builds a inbound vote message for a token sent to TSS
func (ob *Observer) buildInboundVoteMsgForTokenSentToTSS(
	tx *client.Transaction,
	sender ethcommon.Address,
	blockNumber uint64,
) *types.MsgVoteInbound {
	message := tx.Input

	// compliance check
	maybeReceiver := ""
	parsedAddress, _, err := memo.DecodeLegacyMemoHex(message)
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		maybeReceiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(sender.Hex(), maybeReceiver) {
		compliance.PrintComplianceLog(ob.Logger().Inbound, ob.Logger().Compliance,
			false, ob.Chain().ChainId, tx.Hash, sender.Hex(), sender.Hex(), "Gas")
		return nil
	}

	// donation check
	// #nosec G703 err is already checked
	data, _ := hex.DecodeString(message)
	if bytes.Equal(data, []byte(constant.DonationMessage)) {
		ob.Logger().Inbound.Info().
			Msgf("thank you rich folk for your donation! tx %s chain %d", tx.Hash, ob.Chain().ChainId)
		return nil
	}
	ob.Logger().Inbound.Info().Msgf("TSS inbound detected on chain %d tx %s block %d from %s value %s message %s",
		ob.Chain().ChainId, tx.Hash, blockNumber, sender.Hex(), tx.Value.String(), message)

	return musecore.GetInboundVoteMessage(
		sender.Hex(),
		ob.Chain().ChainId,
		sender.Hex(),
		sender.Hex(),
		ob.MusecoreClient().Chain().ChainId,
		sdkmath.NewUintFromBigInt(tx.Value),
		message,
		tx.Hash,
		blockNumber,
		90_000,
		coin.CoinType_Gas,
		"",
		ob.MusecoreClient().GetKeys().GetOperatorAddress().String(),
		0, // not a smart contract call
		types.InboundStatus_SUCCESS,
	)
}

// observeTSSReceiveInBlock queries the incoming gas asset to TSS address in a single block and posts votes
func (ob *Observer) observeTSSReceiveInBlock(ctx context.Context, blockNumber uint64) error {
	block, err := ob.GetBlockByNumberCached(ctx, blockNumber)
	if err != nil {
		return errors.Wrapf(err, "error getting block %d for chain %d", blockNumber, ob.Chain().ChainId)
	}
	for i := range block.Transactions {
		tx := block.Transactions[i]
		if ethcommon.HexToAddress(tx.To) == ob.TSS().PubKey().AddressEVM() {
			receipt, err := ob.evmClient.TransactionReceipt(ctx, ethcommon.HexToHash(tx.Hash))
			if err != nil {
				return errors.Wrapf(err, "error getting receipt for inbound %s chain %d", tx.Hash, ob.Chain().ChainId)
			}

			_, err = ob.checkAndVoteInboundTokenGas(ctx, &tx, receipt, true)
			if err != nil {
				return errors.Wrapf(
					err,
					"error checking and voting inbound gas asset for inbound %s chain %d",
					tx.Hash,
					ob.Chain().ChainId,
				)
			}
		}
	}
	return nil
}
