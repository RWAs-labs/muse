package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestGetAllAuthzMuseclientTxTypes(t *testing.T) {
	require.Equal(t, []string{"/musechain.musecore.crosschain.MsgVoteGasPrice",
		"/musechain.musecore.crosschain.MsgVoteInbound",
		"/musechain.musecore.crosschain.MsgVoteOutbound",
		"/musechain.musecore.crosschain.MsgAddOutboundTracker",
		"/musechain.musecore.observer.MsgVoteTSS",
		"/musechain.musecore.observer.MsgVoteBlame",
		"/musechain.musecore.observer.MsgVoteBlockHeader"},
		crosschaintypes.GetAllAuthzMuseclientTxTypes())
}
