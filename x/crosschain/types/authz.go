package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// GetAllAuthzMuseclientTxTypes returns all the authz types for required for museclient
func GetAllAuthzMuseclientTxTypes() []string {
	return []string{
		sdk.MsgTypeURL(&MsgVoteGasPrice{}),
		sdk.MsgTypeURL(&MsgVoteInbound{}),
		sdk.MsgTypeURL(&MsgVoteOutbound{}),
		sdk.MsgTypeURL(&MsgAddOutboundTracker{}),
		sdk.MsgTypeURL(&observertypes.MsgVoteTSS{}),
		sdk.MsgTypeURL(&observertypes.MsgVoteBlame{}),
		sdk.MsgTypeURL(&observertypes.MsgVoteBlockHeader{}),
	}
}
