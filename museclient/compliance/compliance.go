// Package compliance provides functions to check for compliance of cross-chain transactions
package compliance

import (
	"github.com/rs/zerolog"

	"github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/museclient/logs"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// IsCCTXRestricted returns true if the cctx involves restricted addresses
func IsCCTXRestricted(cctx *crosschaintypes.CrossChainTx) bool {
	sender := cctx.InboundParams.Sender
	receiver := cctx.GetCurrentOutboundParam().Receiver

	return config.ContainRestrictedAddress(sender, receiver)
}

// PrintComplianceLog prints compliance log with fields [chain, cctx/inbound, chain, sender, receiver, token]
func PrintComplianceLog(
	logger, complianceLogger zerolog.Logger,
	outbound bool,
	chainID int64,
	identifier, sender, receiver, token string,
) {
	var (
		message string
		fields  map[string]any
	)

	if outbound {
		message = "Restricted address detected in cctx"
		fields = map[string]any{
			logs.FieldChain:    chainID,
			logs.FieldCctx:     identifier,
			logs.FieldCoinType: token,
			"sender":           sender,
			"receiver":         receiver,
		}
	} else {
		message = "Restricted address detected in inbound"
		fields = map[string]any{
			logs.FieldChain:    chainID,
			logs.FieldCoinType: token,
			logs.FieldTx:       identifier,
			"sender":           sender,
			"receiver":         receiver,
		}
	}

	logger.Warn().Fields(fields).Msg(message)
	complianceLogger.Warn().Fields(fields).Msg(message)
}
