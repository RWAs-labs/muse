package cli

import (
	"github.com/RWAs-labs/go-tss/blame"

	"github.com/RWAs-labs/muse/x/observer/types"
)

func ConvertNodes(n []blame.Node) (nodes []*types.Node) {
	for _, node := range n {
		var entry types.Node
		entry.PubKey = node.Pubkey
		entry.BlameSignature = node.BlameSignature
		entry.BlameData = node.BlameData

		nodes = append(nodes, &entry)
	}
	return
}
