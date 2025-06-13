// Package signer implements the ChainSigner interface for BTC
package signer

import (
	"bytes"
	"context"
	"encoding/hex"
	"runtime/debug"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/client"
	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/observer"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/retry"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

const (
	// broadcastBackoff is the backoff duration for retrying broadcast
	broadcastBackoff = time.Second * 6

	// broadcastRetries is the maximum number of retries for broadcasting a transaction
	broadcastRetries = 10
)

type RPC interface {
	GetNetworkInfo(ctx context.Context) (*btcjson.GetNetworkInfoResult, error)
	GetRawTransaction(ctx context.Context, hash *chainhash.Hash) (*btcutil.Tx, error)
	GetEstimatedFeeRate(ctx context.Context, confTarget int64) (uint64, error)
	SendRawTransaction(ctx context.Context, tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)
	GetMempoolTxsAndFees(ctx context.Context, childHash string) (client.MempoolTxsAndFees, error)
}

// Signer deals with signing & broadcasting BTC transactions.
type Signer struct {
	*base.Signer
	rpc      RPC
	isRegnet bool
}

// New creates a new Bitcoin signer
func New(baseSigner *base.Signer, rpc RPC) *Signer {
	return &Signer{
		Signer:   baseSigner,
		rpc:      rpc,
		isRegnet: chains.IsBitcoinRegnet(baseSigner.Chain().ChainId),
	}
}

// Broadcast sends the signed transaction to the network
func (signer *Signer) Broadcast(ctx context.Context, signedTx *wire.MsgTx) error {
	var outBuff bytes.Buffer
	if err := signedTx.Serialize(&outBuff); err != nil {
		return errors.Wrap(err, "unable to serialize tx")
	}

	signer.Logger().Std.Info().
		Str(logs.FieldTx, signedTx.TxHash().String()).
		Str("signer.tx_payload", hex.EncodeToString(outBuff.Bytes())).
		Msg("Broadcasting transaction")

	_, err := signer.rpc.SendRawTransaction(ctx, signedTx, true)
	if err != nil {
		return errors.Wrap(err, "unable to broadcast raw tx")
	}

	return nil
}

// TryProcessOutbound signs and broadcasts a BTC transaction from a new outbound
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *types.CrossChainTx,
	observer *observer.Observer,
	height uint64,
) {
	outboundID := base.OutboundIDFromCCTX(cctx)
	signer.MarkOutbound(outboundID, true)

	// end outbound process on panic
	defer func() {
		signer.MarkOutbound(outboundID, false)
		if err := recover(); err != nil {
			signer.Logger().
				Std.Error().
				Str(logs.FieldMethod, "TryProcessOutbound").
				Str(logs.FieldCctx, cctx.Index).
				Interface("panic", err).
				Str("stack_trace", string(debug.Stack())).
				Msg("caught panic error")
		}
	}()

	// prepare logger
	params := cctx.GetCurrentOutboundParam()
	lf := map[string]any{
		logs.FieldMethod: "TryProcessOutbound",
		logs.FieldCctx:   cctx.Index,
		logs.FieldNonce:  params.TssNonce,
	}
	signerAddress, err := observer.MusecoreClient().GetKeys().GetAddress()
	if err != nil {
		return
	}
	lf["signer"] = signerAddress.String()
	logger := signer.Logger().Std.With().Fields(lf).Logger()

	// query network info to get minRelayFee (typically 1000 satoshis)
	networkInfo, err := signer.rpc.GetNetworkInfo(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed get bitcoin network info")
		return
	}
	minRelayFee := networkInfo.RelayFee
	if minRelayFee <= 0 {
		logger.Error().Msgf("invalid minimum relay fee: %f", minRelayFee)
		return
	}

	// compliance check
	isRestricted := !signer.PassesCompliance(cctx)

	// setup outbound data
	txData, err := NewOutboundData(cctx, height, minRelayFee, isRestricted, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to setup Bitcoin outbound data")
		return
	}

	var (
		signedTx       *wire.MsgTx
		stuckTx, found = observer.LastStuckOutbound()
		rbfTx          = found && stuckTx.Nonce == txData.nonce
	)

	// sign outbound
	if !rbfTx {
		// sign withdraw tx
		signedTx, err = signer.SignWithdrawTx(ctx, txData, observer)
		if err != nil {
			logger.Error().Err(err).Msg("SignWithdrawTx failed")
			return
		}
		logger.Info().Str(logs.FieldTx, signedTx.TxID()).Msg("SignWithdrawTx succeed")
	} else {
		// sign RBF tx
		signedTx, err = signer.SignRBFTx(ctx, txData, stuckTx.Tx)
		if err != nil {
			logger.Error().Err(err).Msg("SignRBFTx failed")
			return
		}
		logger.Info().Str(logs.FieldTx, signedTx.TxID()).Msg("SignRBFTx succeed")
	}

	// broadcast signed outbound
	signer.BroadcastOutbound(ctx, signedTx, params.TssNonce, rbfTx, cctx, observer)
}

// BroadcastOutbound sends the signed transaction to the Bitcoin network
func (signer *Signer) BroadcastOutbound(
	ctx context.Context,
	tx *wire.MsgTx,
	nonce uint64,
	rbfTx bool,
	cctx *types.CrossChainTx,
	ob *observer.Observer,
) {
	txHash := tx.TxID()

	// prepare logger fields
	logger := signer.Logger().Std.With().
		Str(logs.FieldMethod, "BroadcastOutbound").
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txHash).
		Str(logs.FieldCctx, cctx.Index).
		Logger()

	// double check to ensure the tx being replaced is still the last outbound.
	// when CCTX gets stuck at nonce 'N', the pending nonce will stop incrementing
	// and stay at 'N' or 'N+1' (at most).
	if rbfTx && ob.GetPendingNonce() > nonce+1 {
		logger.Warn().Msgf("RBF tx nonce is outdated, skipping broadcasting")
		return
	}

	// try broacasting tx with backoff in case of RPC error
	broadcast := func() error {
		return retry.Retry(signer.Broadcast(ctx, tx))
	}

	bo := backoff.NewConstantBackOff(broadcastBackoff)
	boWithMaxRetries := backoff.WithMaxRetries(bo, broadcastRetries)
	if err := retry.DoWithBackoff(broadcast, boWithMaxRetries); err != nil {
		logger.Error().Err(err).Msgf("unable to broadcast Bitcoin outbound")
	}
	logger.Info().Msg("broadcasted Bitcoin outbound successfully")

	// save tx local db and ignore db error.
	// db error is not critical and should not block outbound tracker.
	if err := ob.SaveBroadcastedTx(txHash, nonce); err != nil {
		logger.Error().Err(err).Msg("unable to save broadcasted Bitcoin outbound")
	}

	// add tx to outbound tracker so that all observers know about it
	museHash, err := ob.MusecoreClient().PostOutboundTracker(ctx, ob.Chain().ChainId, nonce, txHash)
	if err != nil {
		logger.Err(err).Msg("unable to add Bitcoin outbound tracker")
	} else {
		logger.Info().Str(logs.FieldMuseTx, museHash).Msg("add Bitcoin outbound tracker successfully")
	}

	// try including this outbound as early as possible, no need to wait for outbound tracker
	_, included := ob.TryIncludeOutbound(ctx, cctx, txHash)
	if included {
		logger.Info().Msg("included newly broadcasted Bitcoin outbound")
	}
}
