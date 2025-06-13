package cctx

import (
	"fmt"

	"github.com/RWAs-labs/muse/cmd/musetool/context"
	"github.com/RWAs-labs/muse/pkg/chains"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TrackingDetails tracks the status of a CCTX transaction
type TrackingDetails struct {
	CCTXIdentifier          string       `json:"cctx_identifier"`
	Status                  Status       `json:"status"`
	OutboundChain           chains.Chain `json:"outbound_chain_id"`
	OutboundTssNonce        uint64       `json:"outbound_tss_nonce"`
	OutboundTrackerHashList []string     `json:"outbound_tracker_hash_list"`
	Message                 string       `json:"message"`
}

func NewTrackingDetails() *TrackingDetails {
	return &TrackingDetails{
		CCTXIdentifier: "",
		Status:         Unknown,
	}
}

// UpdateStatusFromMusecoreCCTX updates the status of the TrackingDetails from the musecore CCTX status
func (c *TrackingDetails) UpdateStatusFromMusecoreCCTX(status crosschaintypes.CctxStatus) {
	switch status {
	case crosschaintypes.CctxStatus_PendingOutbound:
		c.Status = PendingOutbound
	case crosschaintypes.CctxStatus_OutboundMined:
		c.Status = OutboundMined
	case crosschaintypes.CctxStatus_Reverted:
		c.Status = Reverted
	case crosschaintypes.CctxStatus_PendingRevert:
		c.Status = PendingRevert
	case crosschaintypes.CctxStatus_Aborted:
		c.Status = Aborted
	default:
		c.Status = Unknown
	}
}

func (c *TrackingDetails) Print() string {
	return fmt.Sprintf("CCTX Identifier: %s Status: %s", c.CCTXIdentifier, c.Status.String())
}

func (c *TrackingDetails) DebugPrint() string {
	return fmt.Sprintf("CCTX Identifier: %s Status: %s Message: %s", c.CCTXIdentifier, c.Status.String(), c.Message)
}

// UpdateCCTXStatus updates the TrackingDetails with status from musecore
func (c *TrackingDetails) UpdateCCTXStatus(ctx *context.Context) {
	var (
		musecoreClient = ctx.GetMuseCoreClient()
		goCtx          = ctx.GetContext()
	)

	CCTX, err := musecoreClient.GetCctxByHash(goCtx, c.CCTXIdentifier)
	if err != nil {
		c.Message = fmt.Sprintf("failed to get cctx: %v", err)
		return
	}

	c.UpdateStatusFromMusecoreCCTX(CCTX.CctxStatus.Status)

	return
}

// UpdateCCTXOutboundDetails updates the TrackingDetails with the outbound chain and nonce
func (c *TrackingDetails) UpdateCCTXOutboundDetails(ctx *context.Context) {
	var (
		musecoreClient = ctx.GetMuseCoreClient()
		goCtx          = ctx.GetContext()
	)
	CCTX, err := musecoreClient.GetCctxByHash(goCtx, c.CCTXIdentifier)
	if err != nil {
		c.Message = fmt.Sprintf("failed to get cctx: %v", err)
	}
	outboundParams := CCTX.GetCurrentOutboundParam()
	if outboundParams == nil {
		c.Message = "outbound params not found"
		return
	}
	chainID := CCTX.GetCurrentOutboundParam().ReceiverChainId

	// This is almost impossible to happen as the cctx would not have been created if the chain was not supported
	chain, found := chains.GetChainFromChainID(chainID, []chains.Chain{})
	if !found {
		c.Message = fmt.Sprintf("receiver chain not supported,chain id: %d", chainID)
	}
	c.OutboundChain = chain
	c.OutboundTssNonce = CCTX.GetCurrentOutboundParam().TssNonce
	return
}

// UpdateHashListAndPendingStatus updates the TrackingDetails with the hash list and updates pending status
// If the tracker is found, it means the outbound is broadcast, but we are waiting for the confirmations
// If the tracker is not found, it means the outbound is not broadcast yet; we are waiting for the tss to sign the outbound
func (c *TrackingDetails) UpdateHashListAndPendingStatus(ctx *context.Context) {
	var (
		musecoreClient = ctx.GetMuseCoreClient()
		goCtx          = ctx.GetContext()
		outboundChain  = c.OutboundChain
		outboundNonce  = c.OutboundTssNonce
	)

	tracker, err := musecoreClient.GetOutboundTracker(goCtx, outboundChain, outboundNonce)
	// the tracker is found that means the outbound has been broadcast, but we are waiting for confirmations
	if err == nil && tracker != nil {
		c.updateOutboundConfirmation()
		var hashList []string
		for _, hash := range tracker.HashList {
			hashList = append(hashList, hash.TxHash)
		}
		c.OutboundTrackerHashList = hashList
		return
	}
	// the cctx is in pending state, but the outbound signing has not been done
	c.updateOutboundSigning()
	return
}

// IsInboundFinalized checks if the inbound voting has been finalized
func (c *TrackingDetails) IsInboundFinalized() bool {
	return !(c.Status == PendingInboundConfirmation || c.Status == PendingInboundVoting)
}

// IsPendingOutbound checks if the cctx is pending processing the outbound transaction (outbound or revert)
func (c *TrackingDetails) IsPendingOutbound() bool {
	return c.Status == PendingOutbound || c.Status == PendingRevert
}

// IsPendingConfirmation checks if the cctx is pending outbound confirmation (outbound or revert
func (c *TrackingDetails) IsPendingConfirmation() bool {
	return c.Status == PendingOutboundConfirmation || c.Status == PendingRevertConfirmation
}

// State transitions for TrackingDetails
// 0 - Inbound Confirmation
func (c *TrackingDetails) updateInboundConfirmation(isConfirmed bool) {
	c.Status = PendingInboundConfirmation
	if isConfirmed {
		c.Status = PendingInboundVoting
	}
}

// 1 - Outbound Signing
func (c *TrackingDetails) updateOutboundSigning() {
	switch {
	case c.Status == PendingOutbound:
		c.Status = PendingOutboundSigning
	case c.Status == PendingRevert:
		c.Status = PendingRevertSigning
	}
}

// 2 - Outbound Confirmation
func (c *TrackingDetails) updateOutboundConfirmation() {
	switch {
	case c.Status == PendingOutbound:
		c.Status = PendingOutboundConfirmation
	case c.Status == PendingRevert:
		c.Status = PendingRevertConfirmation
	}
}

// 3 - Outbound Voting
func (c *TrackingDetails) updateOutboundVoting() {
	switch {
	case c.Status == PendingOutboundConfirmation:
		c.Status = PendingOutboundVoting
	case c.Status == PendingRevertConfirmation:
		c.Status = PendingRevertVoting
	}
}
