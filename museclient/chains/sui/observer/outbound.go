package observer

import (
	"context"
	"strconv"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/chains/sui/client"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/museclient/musecore"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/contracts/sui"
	cctypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// 50 SUI
// https://docs.sui.io/concepts/tokenomics/gas-in-sui#gas-budgets
const maxGasLimit = 50_000_000_000

// OutboundCreated checks if the outbound tx exists in the memory
// and has valid nonce & signature
func (ob *Observer) OutboundCreated(nonce uint64) bool {
	_, ok := ob.getTx(nonce)
	return ok
}

// ProcessOutboundTrackers loads all freshly-included Sui transactions in-memory
// for further voting by Observer-Signer.
func (ob *Observer) ProcessOutboundTrackers(ctx context.Context) error {
	chainID := ob.Chain().ChainId

	trackers, err := ob.MusecoreClient().GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
	if err != nil {
		return errors.Wrap(err, "unable to get outbound trackers")
	}

	for _, tracker := range trackers {
		nonce := tracker.Nonce

		// already loaded
		if _, ok := ob.getTx(nonce); ok {
			continue
		}

		// should not happen
		if len(tracker.HashList) == 0 {
			// we don't want to block other cctxs, so let's error and continue
			ob.Logger().Outbound.Error().
				Str(logs.FieldMethod, "ProcessOutboundTrackers").
				Uint64(logs.FieldNonce, nonce).
				Str(logs.FieldTracker, tracker.Index).
				Msg("Tracker hash list is empty!")
			continue
		}

		digest := tracker.HashList[0].TxHash

		cctx, err := ob.MusecoreClient().GetCctxByNonce(ctx, chainID, tracker.Nonce)
		if err != nil {
			return errors.Wrapf(err, "unable to get cctx by nonce %d (sui digest %q)", tracker.Nonce, digest)
		}

		if err := ob.loadOutboundTx(ctx, cctx, digest); err != nil {
			// we don't want to block other cctxs, so let's error and continue
			ob.Logger().Outbound.
				Error().Err(err).
				Str(logs.FieldMethod, "ProcessOutboundTrackers").
				Uint64(logs.FieldNonce, nonce).
				Str(logs.FieldTx, digest).
				Msg("Unable to load outbound transaction")
		}
	}

	return nil
}

// VoteOutbound calculates outbound result based on cctx and in-mem Sui tx
// and votes the ballot to musecore.
func (ob *Observer) VoteOutbound(ctx context.Context, cctx *cctypes.CrossChainTx) error {
	chainID := ob.Chain().ChainId
	params := cctx.GetCurrentOutboundParam()
	nonce := params.TssNonce

	// should be fetched by ProcessOutboundTrackers routine
	// if exists, we can safely assume it's authentic and nonce is valid
	tx, ok := ob.getTx(nonce)
	if !ok {
		return errors.Errorf("missing tx for nonce %d", nonce)
	}

	// used checkpoint instead of block height
	checkpoint, err := strconv.ParseUint(tx.Checkpoint, 10, 64)
	if err != nil {
		return errors.Wrap(err, "unable to parse checkpoint")
	}

	// parse outbound event
	event, content, err := ob.gateway.ParseOutboundEvent(tx)
	if err != nil {
		return errors.Wrap(err, "unable to parse outbound event")
	}

	// determine amount, status and coinType
	var (
		amount   = content.TokenAmount()
		status   = chains.ReceiveStatus_success
		coinType = cctx.InboundParams.CoinType
	)

	// cancelled transaction means the outbound is failed
	// - set amount to CCTX's amount to bypass amount check in musecore
	// - set status to failed to revert the CCTX in musecore
	if event.IsCancelTx() {
		amount = params.Amount
		status = chains.ReceiveStatus_failed
	}

	// Gas parameters
	// Gas price *might* change once per epoch (~24h), so using the latest value is fine.
	// #nosec G115 - always in range
	outboundGasPrice := math.NewInt(int64(ob.getLatestGasPrice()))

	// This might happen after musecore restart when PostGasPrice has not been called yet. retry later.
	if outboundGasPrice.IsZero() {
		return errors.New("latest gas price is zero")
	}

	outboundGasUsed, err := parseGasUsed(tx)
	if err != nil {
		return errors.Wrap(err, "unable to parse gas used")
	}

	// Create message
	msg := cctypes.NewMsgVoteOutbound(
		ob.MusecoreClient().GetKeys().GetOperatorAddress().String(),
		cctx.Index,
		tx.Digest,
		checkpoint,
		outboundGasUsed,
		outboundGasPrice,
		maxGasLimit,
		amount,
		status,
		chainID,
		nonce,
		coinType,
		cctypes.ConfirmationMode_SAFE,
	)

	if err := ob.postVoteOutbound(ctx, msg); err != nil {
		return errors.Wrap(err, "unable to post vote outbound")
	}

	ob.unsetTx(nonce)

	return nil
}

// loadOutboundTx loads cross-chain outbound tx by digest and ensures its authenticity.
func (ob *Observer) loadOutboundTx(ctx context.Context, cctx *cctypes.CrossChainTx, digest string) error {
	res, err := ob.client.SuiGetTransactionBlock(ctx, models.SuiGetTransactionBlockRequest{
		Digest: digest,
		Options: models.SuiTransactionBlockOptions{
			ShowEvents:  true,
			ShowInput:   true,
			ShowEffects: true,
		},
	})

	if err != nil {
		return errors.Wrap(err, "unable to get tx")
	}

	if err := ob.validateOutbound(cctx, res); err != nil {
		return errors.Wrap(err, "tx validation failed")
	}

	ob.setTx(res, cctx.GetCurrentOutboundParam().TssNonce)

	return nil
}

// validateOutbound validates the authenticity of the outbound transaction.
func (ob *Observer) validateOutbound(cctx *cctypes.CrossChainTx, tx models.SuiTransactionBlockResponse) error {
	// a valid outbound should be successful
	if tx.Effects.Status.Status != client.TxStatusSuccess {
		return errors.Errorf("tx failed with error: %s", tx.Effects.Status.Error)
	}

	// parse outbound event
	_, content, err := ob.gateway.ParseOutboundEvent(tx)
	if err != nil {
		return errors.Wrap(err, "unable to parse outbound event")
	}

	// tx nonce should match CCTX nonce
	txNonce := content.TxNonce()
	nonce := cctx.GetCurrentOutboundParam().TssNonce
	if txNonce != nonce {
		return errors.Errorf("nonce mismatch (tx nonce %d, cctx nonce %d)", txNonce, nonce)
	}

	// check tx signature
	if len(tx.Transaction.TxSignatures) == 0 {
		return errors.New("missing tx signature")
	}

	pubKey, _, err := sui.DeserializeSignatureECDSA(tx.Transaction.TxSignatures[0])
	if err != nil {
		return errors.Wrap(err, "unable to deserialize tx signature")
	}

	if !ob.TSS().PubKey().AsECDSA().Equal(pubKey) {
		return errors.New("pubKey mismatch")
	}

	return nil
}

func (ob *Observer) postVoteOutbound(ctx context.Context, msg *cctypes.MsgVoteOutbound) error {
	const gasLimit = musecore.PostVoteOutboundGasLimit

	retryGasLimit := musecore.PostVoteOutboundRetryGasLimit
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = musecore.PostVoteOutboundRevertGasLimit
	}

	museTxHash, ballot, err := ob.MusecoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	switch {
	case err != nil:
		return errors.Wrap(err, "unable to post vote outbound")
	case museTxHash != "":
		ob.Logger().Outbound.Info().
			Str(logs.FieldTx, msg.ObservedOutboundHash).
			Str(logs.FieldMuseTx, museTxHash).
			Str(logs.FieldBallot, ballot).
			Msg("PostVoteOutbound: posted outbound vote successfully")
	}

	return nil
}

func (ob *Observer) getTx(nonce uint64) (models.SuiTransactionBlockResponse, bool) {
	ob.txMu.RLock()
	defer ob.txMu.RUnlock()

	tx, ok := ob.txMap[nonce]

	return tx, ok
}

func (ob *Observer) setTx(tx models.SuiTransactionBlockResponse, nonce uint64) {
	ob.txMu.Lock()
	defer ob.txMu.Unlock()

	ob.txMap[nonce] = tx
}

func (ob *Observer) unsetTx(nonce uint64) {
	ob.txMu.Lock()
	defer ob.txMu.Unlock()

	delete(ob.txMap, nonce)
}

func parseGasUsed(tx models.SuiTransactionBlockResponse) (uint64, error) {
	gas := tx.Effects.GasUsed

	compCost, err := parseUint64(gas.ComputationCost)
	if err != nil {
		return 0, errors.Wrap(err, "comp cost")
	}

	storageCost, err := parseUint64(gas.StorageCost)
	if err != nil {
		return 0, errors.Wrap(err, "storage cost")
	}

	storageRebate, err := parseUint64(gas.StorageRebate)
	if err != nil {
		return 0, errors.Wrap(err, "storage rebate")
	}

	return compCost + storageCost - storageRebate, nil
}

func parseUint64(v string) (uint64, error) {
	return strconv.ParseUint(v, 10, 64)
}
