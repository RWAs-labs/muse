package observer

import (
	"bytes"
	"encoding/hex"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/compliance"
	"github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/museclient/logs"
	clienttypes "github.com/RWAs-labs/muse/museclient/types"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/pkg/crypto"
	"github.com/RWAs-labs/muse/pkg/memo"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// BTCInboundEvent represents an incoming transaction event
type BTCInboundEvent struct {
	// FromAddress is the first input address
	FromAddress string

	// ToAddress is the MEVM receiver address
	ToAddress string

	// Value is the amount of BTC
	Value float64

	// DepositorFee is the deposit fee
	DepositorFee float64

	// MemoBytes is the memo of inbound
	MemoBytes []byte

	// MemoStd is the standard inbound memo if it can be decoded
	MemoStd *memo.InboundMemo

	// BlockNumber is the block number of the inbound
	BlockNumber uint64

	// TxHash is the hash of the inbound
	TxHash string

	// Status is the status of the inbound event
	Status crosschaintypes.InboundStatus
}

// Category returns the category of the inbound event
func (event *BTCInboundEvent) Category() clienttypes.InboundCategory {
	// compliance check on sender and receiver addresses
	if config.ContainRestrictedAddress(event.FromAddress, event.ToAddress) {
		return clienttypes.InboundCategoryRestricted
	}

	// compliance check on receiver, revert/abort addresses in standard memo
	if event.MemoStd != nil {
		if config.ContainRestrictedAddress(
			event.MemoStd.Receiver.Hex(),
			event.MemoStd.RevertOptions.RevertAddress,
			event.MemoStd.RevertOptions.AbortAddress,
		) {
			return clienttypes.InboundCategoryRestricted
		}
	}

	// donation check
	if bytes.Equal(event.MemoBytes, []byte(constant.DonationMessage)) {
		return clienttypes.InboundCategoryDonation
	}

	return clienttypes.InboundCategoryProcessable
}

// DecodeMemoBytes decodes the contained memo bytes as either standard or legacy memo
// It updates the event object with the decoded data
func (event *BTCInboundEvent) DecodeMemoBytes(chainID int64) error {
	var (
		err            error
		isStandardMemo bool
		memoStd        *memo.InboundMemo
		receiver       ethcommon.Address
	)

	// skip decoding if no memo is found, returning error to revert the inbound
	if bytes.Equal(event.MemoBytes, []byte(noMemoFound)) {
		event.MemoBytes = []byte{}
		return errors.New("no memo found in inbound")
	}

	// skip decoding donation tx as it won't go through musecore
	if bytes.Equal(event.MemoBytes, []byte(constant.DonationMessage)) {
		return nil
	}

	// try to decode the standard memo as the preferred format
	// the standard memo is NOT enabled for Bitcoin mainnet

	if chainID != chains.BitcoinMainnet.ChainId {
		memoStd, isStandardMemo, err = memo.DecodeFromBytes(event.MemoBytes)
	}

	// process standard memo or fallback to legacy memo
	if isStandardMemo {
		// skip standard memo that carries improper data
		if err != nil {
			return errors.Wrap(err, "standard memo contains improper data")
		}

		// validate the content of the standard memo
		err = ValidateStandardMemo(*memoStd, chainID)
		if err != nil {
			return errors.Wrap(err, "invalid standard memo for bitcoin")
		}

		event.MemoStd = memoStd
		receiver = memoStd.Receiver
	} else {
		parsedAddress, payload, err := memo.DecodeLegacyMemoHex(hex.EncodeToString(event.MemoBytes))
		if err != nil { // unreachable code
			return errors.Wrap(err, "invalid legacy memo")
		}
		receiver = parsedAddress

		// update the memo bytes to only contain the data
		event.MemoBytes = payload
	}

	// ensure the receiver is valid
	if crypto.IsEmptyAddress(receiver) {
		return errors.New("got empty receiver address from memo")
	}
	event.ToAddress = receiver.Hex()

	return nil
}

// ValidateStandardMemo validates the standard memo in Bitcoin context
func ValidateStandardMemo(memoStd memo.InboundMemo, chainID int64) error {
	// NoAssetCall will be disabled for Bitcoin until full V2 support
	// https://github.com/RWAs-labs/muse/issues/2711
	if memoStd.OpCode == memo.OpCodeCall {
		return errors.New("NoAssetCall is disabled for Bitcoin")
	}

	// ensure the revert address is a valid and supported BTC address
	revertAddress := memoStd.RevertOptions.RevertAddress
	if revertAddress != "" {
		btcAddress, err := chains.DecodeBtcAddress(revertAddress, chainID)
		if err != nil {
			return errors.Wrapf(err, "invalid revert address in memo: %s", revertAddress)
		}
		if !chains.IsBtcAddressSupported(btcAddress) {
			return fmt.Errorf("unsupported revert address in memo: %s", revertAddress)
		}
	}

	return nil
}

// IsEventProcessable checks if the inbound event is processable
func (ob *Observer) IsEventProcessable(event BTCInboundEvent) bool {
	logFields := map[string]any{logs.FieldTx: event.TxHash}

	switch category := event.Category(); category {
	case clienttypes.InboundCategoryProcessable:
		return true
	case clienttypes.InboundCategoryDonation:
		ob.Logger().Inbound.Info().Fields(logFields).Msgf("thank you rich folk for your donation!")
		return false
	case clienttypes.InboundCategoryRestricted:
		compliance.PrintComplianceLog(ob.logger.Inbound, ob.logger.Compliance,
			false, ob.Chain().ChainId, event.TxHash, event.FromAddress, event.ToAddress, "BTC")
		return false
	default:
		ob.Logger().Inbound.Error().Fields(logFields).Msgf("unreachable code got InboundCategory: %v", category)
		return false
	}
}
