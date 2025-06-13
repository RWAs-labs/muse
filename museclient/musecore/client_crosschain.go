package musecore

import (
	"context"
	"sort"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/metrics"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

func (c *Client) ListPendingCCTX(ctx context.Context, chain chains.Chain) ([]*types.CrossChainTx, uint64, error) {
	list, total, err := c.Clients.ListPendingCCTX(ctx, chain.ChainId)

	if err == nil {
		value := float64(total)

		metrics.PendingTxsPerChain.WithLabelValues(chain.Name).Set(value)
	}

	return list, total, err
}

// GetAllOutboundTrackerByChain returns all outbound trackers for a chain
func (c *Client) GetAllOutboundTrackerByChain(
	ctx context.Context,
	chainID int64,
	order interfaces.Order,
) ([]types.OutboundTracker, error) {
	in := &types.QueryAllOutboundTrackerByChainRequest{
		Chain: chainID,
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      2000,
			CountTotal: false,
			Reverse:    false,
		},
	}

	resp, err := c.Crosschain.OutboundTrackerAllByChain(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all outbound trackers")
	}

	if order == interfaces.Ascending {
		sort.SliceStable(resp.OutboundTracker, func(i, j int) bool {
			return resp.OutboundTracker[i].Nonce < resp.OutboundTracker[j].Nonce
		})
	} else if order == interfaces.Descending {
		sort.SliceStable(resp.OutboundTracker, func(i, j int) bool {
			return resp.OutboundTracker[i].Nonce > resp.OutboundTracker[j].Nonce
		})
	}

	return resp.OutboundTracker, nil
}
