package observer

import (
	"context"

	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/chains/ton/liteapi"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/musecore"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	toncontracts "github.com/RWAs-labs/muse/pkg/contracts/ton"
	cctypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// https://tonscan.com/config-parameters (N21: "Computation costs")
// This might changes in the future by TON's gov proposal (very unlikely though)
const maxGasLimit = 1_000_000

type outbound struct {
	tx            *toncontracts.Transaction
	receiveStatus chains.ReceiveStatus
	nonce         uint64
}

// VoteOutboundIfConfirmed checks outbound status and returns (continueKeysign, error)
func (ob *Observer) VoteOutboundIfConfirmed(ctx context.Context, cctx *cctypes.CrossChainTx) (bool, error) {
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	outboundRes, exists := ob.getOutboundByNonce(nonce)
	if !exists {
		return true, nil
	}

	withdrawal, err := outboundRes.tx.Withdrawal()
	if err != nil {
		return false, errors.Wrap(err, "unable to get withdrawal")
	}

	// TODO: Add compliance check
	// https://github.com/RWAs-labs/muse/issues/2916

	if err = ob.postVoteOutbound(ctx, cctx, outboundRes, withdrawal); err != nil {
		return false, errors.Wrap(err, "unable to post vote")
	}

	return false, nil
}

// ProcessOutboundTrackers pulls outbounds trackers from musecore,
// fetches txs from TON and stores them in memory for further use.
func (ob *Observer) ProcessOutboundTrackers(ctx context.Context) error {
	var (
		chainID  = ob.Chain().ChainId
		musecore = ob.MusecoreClient()
	)

	trackers, err := musecore.GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
	if err != nil {
		return errors.Wrap(err, "unable to get outbound trackers")
	}

	for _, tracker := range trackers {
		nonce := tracker.Nonce

		// If outbound is already in memory, skip.
		if _, ok := ob.getOutboundByNonce(nonce); ok {
			continue
		}

		// Let's not block other cctxs from being processed
		cctx, err := musecore.GetCctxByNonce(ctx, chainID, nonce)
		if err != nil {
			ob.Logger().Outbound.
				Error().Err(err).
				Uint64("outbound.nonce", nonce).
				Msg("Unable to get cctx by nonce")

			continue
		}

		for _, txHash := range tracker.HashList {
			if err := ob.processOutboundTracker(ctx, cctx, txHash.TxHash); err != nil {
				ob.Logger().Outbound.
					Error().Err(err).
					Uint64("outbound.nonce", nonce).
					Str("outbound.hash", txHash.TxHash).
					Msg("Unable to check transaction by nonce")
			}
		}
	}

	return nil
}

// processOutboundTracker checks TON tx and stores it in memory for further processing
// by VoteOutboundIfConfirmed.
func (ob *Observer) processOutboundTracker(ctx context.Context, cctx *cctypes.CrossChainTx, txHash string) error {
	if cctx.InboundParams.CoinType != coin.CoinType_Gas {
		return errors.New("only gas cctxs are supported")
	}

	lt, hash, err := liteapi.TransactionHashFromString(txHash)
	if err != nil {
		return errors.Wrap(err, "unable to parse tx hash")
	}

	rawTX, err := ob.client.GetTransaction(ctx, ob.gateway.AccountID(), lt, hash)
	if err != nil {
		return errors.Wrap(err, "unable to get transaction form liteapi")
	}

	tx, err := ob.gateway.ParseTransaction(rawTX)
	if err != nil {
		return errors.Wrap(err, "unable to parse transaction")
	}

	receiveStatus, err := ob.determineReceiveStatus(tx)
	if err != nil {
		return errors.Wrap(err, "unable to determine outbound outcome")
	}

	// TODO: Add compliance check
	// https://github.com/RWAs-labs/muse/issues/2916

	nonce := cctx.GetCurrentOutboundParam().TssNonce
	ob.setOutboundByNonce(outbound{tx, receiveStatus, nonce})

	return nil
}

func (ob *Observer) determineReceiveStatus(tx *toncontracts.Transaction) (chains.ReceiveStatus, error) {
	_, evmSigner, err := extractWithdrawal(tx)
	switch {
	case err != nil:
		return 0, err
	case evmSigner != ob.TSS().PubKey().AddressEVM():
		return 0, errors.New("withdrawal signer is not TSS")
	case !tx.IsSuccess():
		return chains.ReceiveStatus_failed, nil
	default:
		return chains.ReceiveStatus_success, nil
	}
}

// addOutboundTracker publishes outbound tracker to Musecore.
// In most cases will be a noop because the tracker is already published by the signer.
// See Signer{}.trackOutbound(...) for more details.
func (ob *Observer) addOutboundTracker(ctx context.Context, tx *toncontracts.Transaction) error {
	w, evmSigner, err := extractWithdrawal(tx)
	switch {
	case err != nil:
		return err
	case evmSigner != ob.TSS().PubKey().AddressEVM():
		ob.Logger().Inbound.Warn().
			Fields(txLogFields(tx)).
			Str("transaction.ton.signer", evmSigner.String()).
			Msg("observeGateway: addOutboundTracker: withdrawal signer is not TSS. Skipping")

		return nil
	}

	var (
		chainID = ob.Chain().ChainId
		nonce   = uint64(w.Seqno)
		hash    = liteapi.TransactionToHashString(tx.Transaction)
	)

	// note it has a check for noop
	_, err = ob.MusecoreClient().PostOutboundTracker(ctx, chainID, nonce, hash)

	return err
}

// return withdrawal and tx signer
func extractWithdrawal(tx *toncontracts.Transaction) (toncontracts.Withdrawal, eth.Address, error) {
	w, err := tx.Withdrawal()
	if err != nil {
		return toncontracts.Withdrawal{}, eth.Address{}, errors.Wrap(err, "not a withdrawal")
	}

	s, err := w.Signer()
	if err != nil {
		return toncontracts.Withdrawal{}, eth.Address{}, errors.Wrap(err, "unable to get signer")
	}

	return w, s, nil
}

// getOutboundByNonce returns outbound by nonce
func (ob *Observer) getOutboundByNonce(nonce uint64) (outbound, bool) {
	v, ok := ob.outbounds.Get(nonce)
	if !ok {
		return outbound{}, false
	}

	return v.(outbound), true
}

// setOutboundByNonce stores outbound by nonce
func (ob *Observer) setOutboundByNonce(o outbound) {
	ob.outbounds.Add(o.nonce, o)
}

func (ob *Observer) postVoteOutbound(
	ctx context.Context,
	cctx *cctypes.CrossChainTx,
	outboundRes outbound,
	w toncontracts.Withdrawal,
) error {
	// There's no sequential block height. Also, different txs might end up in different shards.
	// tlb.BlockID is essentially a workchain+shard+seqno tuple. We can't use it as a block height, thus zero.
	const tonBlockHeight = 0

	var (
		chainID       = ob.Chain().ChainId
		txHash        = liteapi.TransactionToHashString(outboundRes.tx.Transaction)
		nonce         = cctx.GetCurrentOutboundParam().TssNonce
		signerAddress = ob.MusecoreClient().GetKeys().GetOperatorAddress()
		coinType      = cctx.InboundParams.CoinType
	)

	gasPrice, ok := ob.getLatestGasPrice()

	// should not happen
	if !ok {
		return errors.New("gas price is not set (call PostGasPrice first)")
	}

	// #nosec G115 len always in range
	gasPriceInt := math.NewInt(int64(gasPrice))

	msg := cctypes.NewMsgVoteOutbound(
		signerAddress.String(),
		cctx.Index,
		txHash,
		tonBlockHeight,
		outboundRes.tx.GasUsed().Uint64(),
		gasPriceInt,
		maxGasLimit,
		w.Amount,
		outboundRes.receiveStatus,
		chainID,
		nonce,
		coinType,
		cctypes.ConfirmationMode_SAFE,
	)

	const gasLimit = musecore.PostVoteOutboundGasLimit

	retryGasLimit := musecore.PostVoteOutboundRetryGasLimit
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = musecore.PostVoteOutboundRevertGasLimit
	}

	log := ob.Logger().Outbound.With().
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txHash).
		Logger()

	museTxHash, ballot, err := ob.MusecoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)

	switch {
	case err != nil:
		log.Error().Err(err).Msg("PostVoteOutbound: failed to post vote")
		return err
	case museTxHash != "":
		log.Info().
			Str("outbound.vote_tx_hash", museTxHash).
			Str("outbound.ballot_id", ballot).
			Msg("PostVoteOutbound: posted vote")
	}

	return nil
}
