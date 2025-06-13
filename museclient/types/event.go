package types

import (
	"bytes"
	"encoding/hex"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/pkg/contracts/solana"
	"github.com/RWAs-labs/muse/pkg/crypto"
	"github.com/RWAs-labs/muse/pkg/memo"
)

// InboundCategory is an enum representing the category of an inbound event
type InboundCategory int

const (
	// InboundCategoryUnknown represents an unknown inbound
	InboundCategoryUnknown InboundCategory = iota

	// InboundCategoryProcessable represents a processable inbound
	InboundCategoryProcessable

	// InboundCategoryDonation represents a donation inbound
	InboundCategoryDonation

	// InboundCategoryRestricted represents a restricted inbound
	InboundCategoryRestricted
)

// InboundEvent represents an inbound event
// TODO: we should consider using this generic struct when it applies (e.g. for Bitcoin, Solana, etc.)
// https://github.com/RWAs-labs/muse/issues/2495
type InboundEvent struct {
	// SenderChainID is the chain ID of the sender
	SenderChainID int64

	// Sender is the sender address
	Sender string

	// Receiver is the receiver address
	Receiver string

	// TxOrigin is the origin of the transaction
	TxOrigin string

	// Value is the amount of token
	Amount uint64

	// Memo is the memo attached to the inbound
	Memo []byte

	// BlockNumber is the block number of the inbound
	BlockNumber uint64

	// TxHash is the hash of the inbound
	TxHash string

	// Index is the index of the event
	Index uint32

	// CoinType is the coin type of the inbound
	CoinType coin.CoinType

	// Asset is the asset of the inbound
	Asset string

	// IsCrossChainCall is true if the inbound is a cross-chain call
	IsCrossChainCall bool

	// RevertOptions are optional revert options
	RevertOptions *solana.RevertOptions
}

// DecodeMemo decodes the receiver from the memo bytes
func (event *InboundEvent) DecodeMemo() error {
	// skip decoding donation tx as it won't go through musecore
	if bytes.Equal(event.Memo, []byte(constant.DonationMessage)) {
		return nil
	}

	// decode receiver address from memo
	parsedAddress, _, err := memo.DecodeLegacyMemoHex(hex.EncodeToString(event.Memo))
	if err != nil { // unreachable code
		return errors.Wrap(err, "invalid memo hex")
	}

	// ensure the receiver is valid
	if crypto.IsEmptyAddress(parsedAddress) {
		return errors.New("got empty receiver address from memo")
	}
	event.Receiver = parsedAddress.Hex()

	return nil
}

// Category returns the category of the inbound event
func (event *InboundEvent) Category() InboundCategory {
	// parse memo-specified receiver
	receiver := ""
	parsedAddress, _, err := memo.DecodeLegacyMemoHex(hex.EncodeToString(event.Memo))
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		receiver = parsedAddress.Hex()
	}

	// check restricted addresses
	if config.ContainRestrictedAddress(event.Sender, event.Receiver, event.TxOrigin, receiver) {
		return InboundCategoryRestricted
	}

	// donation check
	if bytes.Equal(event.Memo, []byte(constant.DonationMessage)) {
		return InboundCategoryDonation
	}

	return InboundCategoryProcessable
}
