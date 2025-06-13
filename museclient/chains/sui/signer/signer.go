package signer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	suiptb "github.com/pattonkan/sui-go/sui"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/logs"
	"github.com/RWAs-labs/muse/pkg/bg"
	"github.com/RWAs-labs/muse/pkg/contracts/sui"
	cctypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// Signer Sui outbound transaction signer.
type Signer struct {
	*base.Signer
	client      RPC
	gateway     *sui.Gateway
	withdrawCap *withdrawCap

	musecore interfaces.MusecoreClient
}

// RPC represents Sui rpc.
type RPC interface {
	SuiXGetLatestSuiSystemState(ctx context.Context) (models.SuiSystemStateSummary, error)
	GetOwnedObjectID(ctx context.Context, ownerAddress, structType string) (string, error)
	GetObjectParsedData(ctx context.Context, objectID string) (models.SuiParsedData, error)
	SuiMultiGetObjects(ctx context.Context, req models.SuiMultiGetObjectsRequest) ([]*models.SuiObjectResponse, error)
	GetSuiCoinObjectRefs(ctx context.Context, owner string, minBalanceMist uint64) ([]*suiptb.ObjectRef, error)

	MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error)
	SuiExecuteTransactionBlock(
		ctx context.Context,
		req models.SuiExecuteTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
	InspectTransactionBlock(
		ctx context.Context,
		req models.SuiDevInspectTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
	SuiGetTransactionBlock(
		ctx context.Context,
		req models.SuiGetTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error)
}

// New Signer constructor.
func New(
	baseSigner *base.Signer,
	client RPC,
	gateway *sui.Gateway,
	musecore interfaces.MusecoreClient,
) *Signer {
	return &Signer{
		Signer:      baseSigner,
		client:      client,
		gateway:     gateway,
		musecore:    musecore,
		withdrawCap: &withdrawCap{},
	}
}

// ProcessCCTX schedules outbound cross-chain transaction.
// Build --> Sign --> Broadcast --(async)--> Wait for execution --> PostOutboundTracker
func (s *Signer) ProcessCCTX(ctx context.Context, cctx *cctypes.CrossChainTx, museHeight uint64) error {
	var (
		outboundID = base.OutboundIDFromCCTX(cctx)
		nonce      = cctx.GetCurrentOutboundParam().TssNonce
	)

	s.MarkOutbound(outboundID, true)
	defer func() { s.MarkOutbound(outboundID, false) }()

	// prepare logger
	logger := s.Logger().Std.With().Uint64(logs.FieldNonce, nonce).Logger()
	ctx = logger.WithContext(ctx)

	// skip if gateway nonce does not match CCTX nonce:
	// 1. this will avoid unnecessary gateway nonce mismatch error in the 'withdraw_impl'
	// 2. this also avoid unexpected gateway version bump and cause subsequent txs to fail
	gatewayNonce, err := s.getGatewayNonce(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get gateway nonce")
	}

	if gatewayNonce != nonce {
		logger.Info().Msgf("gateway nonce %d does not match CCTX nonce %d, skip broadcast", gatewayNonce, nonce)
		return nil
	}

	withdrawTxBuilder, err := s.createWithdrawTxBuilder(cctx, museHeight)
	if err != nil {
		return errors.Wrap(err, "unable to create withdrawal tx builder")
	}

	// always need a cancel tx as fallback
	cancelTxBuilder, err := s.createCancelTxBuilder(ctx, cctx, museHeight)
	if err != nil {
		return errors.Wrap(err, "unable to create cancel tx builder")
	}

	var (
		txDigest      string
		validReceiver = true
	)

	// check CCTX receiver address format
	receiver := cctx.GetCurrentOutboundParam().Receiver
	if err := ValidateAddress(receiver); err != nil {
		validReceiver = false
		logger.Error().Err(err).Str("receiver", receiver).Msg("Invalid receiver address")
	}

	// broadcast tx according to compliance check result
	if validReceiver && s.PassesCompliance(cctx) {
		txDigest, err = s.broadcastWithdrawalWithFallback(ctx, withdrawTxBuilder, cancelTxBuilder)
	} else {
		txDigest, err = s.broadcastCancelTx(ctx, cancelTxBuilder)
	}

	if err != nil {
		// todo we might need additional error handling
		// for the case when the tx is already broadcasted by another museclient
		// (e.g. suppress error)
		return errors.Wrap(err, "unable to broadcast tx")
	}

	// report outbound tracker
	task := func(ctx context.Context) error { return s.reportOutboundTracker(ctx, nonce, txDigest) }
	bg.Work(ctx, task, bg.WithName("report_outbound_tracker"))

	return nil
}

// signTx signs a tx with TSS and returns a base64 encoded signature.
func (s *Signer) signTx(ctx context.Context, tx models.TxnMetaData, museHeight, nonce uint64) (string, error) {
	digest, err := sui.Digest(tx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get digest")
	}

	// send TSS signature request.
	sig65B, err := s.TSS().Sign(ctx, wrapDigest(digest), museHeight, nonce, s.Chain().ChainId)
	if err != nil {
		return "", errors.Wrap(err, "unable to sign tx")
	}

	// serialize signature
	sigBase64, err := sui.SerializeSignatureECDSA(sig65B, s.TSS().PubKey().AsECDSA())
	if err != nil {
		return "", errors.Wrap(err, "unable to serialize tx signature")
	}

	return sigBase64, nil
}

// SignTxWithCancel signs original tx and cancel tx in one go to save TSS keysign time.
//
// Note: this function is not used due to tx simulation issue in Sui SDK,
// but we can sign both tx and cancel tx in one go once Sui SDK is updated.
func (s *Signer) SignTxWithCancel(
	ctx context.Context,
	tx models.TxnMetaData,
	txCancel models.TxnMetaData,
	museHeight, nonce uint64,
) (sig string, sigCancel string, err error) {
	digests := make([][]byte, 2)

	// tx digest
	digest, err := sui.Digest(tx)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to get tx digest")
	}
	digests[0] = wrapDigest(digest)

	// cancel tx digest
	digestCancel, err := sui.Digest(txCancel)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to get cancel tx digest")
	}
	digests[1] = wrapDigest(digestCancel)

	// sign both digests with TSS
	sig65Bs, err := s.TSS().SignBatch(ctx, digests, museHeight, nonce, s.Chain().ChainId)
	if err != nil {
		return "", "", errors.Wrapf(err, "unable to sign %d tx(s) with TSS", len(digests))
	}

	// should never mismatch
	if len(sig65Bs) != len(digests) {
		return "", "", fmt.Errorf("expected %d signatures, got %d", len(digests), len(sig65Bs))
	}

	// serialize signatures
	sig, err = sui.SerializeSignatureECDSA(sig65Bs[0], s.TSS().PubKey().AsECDSA())
	if err != nil {
		return "", "", errors.Wrap(err, "unable to serialize tx signature")
	}

	sigCancel, err = sui.SerializeSignatureECDSA(sig65Bs[1], s.TSS().PubKey().AsECDSA())
	if err != nil {
		return "", "", errors.Wrap(err, "unable to serialize tx cancel signature")
	}

	return sig, sigCancel, nil
}

// ValidateAddress checks whether the input string is a valid Sui address
// For WithdrawAndCall, the receiver is the target package ID. It follows same format, so we use same validation for both
func ValidateAddress(addr string) error {
	addr, hasPrefix := strings.CutPrefix(addr, "0x")
	switch {
	case !hasPrefix:
		return errors.New("address must start with 0x")
	case addr != strings.ToLower(addr):
		return errors.New("address must be lowercase")
	case len(addr) != 64:
		// accept full Sui address format only to make the validation easier
		return errors.New("address must be 64 characters")
	}

	if _, err := hex.DecodeString(addr); err != nil {
		return errors.Wrapf(err, "address %s is not valid hex", addr)
	}

	return nil
}

// wrapDigest wraps the digest with sha256.
func wrapDigest(digest [32]byte) []byte {
	// another hashing is required for ECDSA.
	// see: https://docs.sui.io/concepts/cryptography/transaction-auth/signatures#signature-requirements
	digestWrapped := sha256.Sum256(digest[:])
	return digestWrapped[:]
}
