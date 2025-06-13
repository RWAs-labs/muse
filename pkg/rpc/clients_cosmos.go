package rpc

import (
	"context"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
)

// GetUpgradePlan returns the current upgrade plan or nil if there is no plan.
func (c *Clients) GetUpgradePlan(ctx context.Context) (*upgradetypes.Plan, error) {
	in := &upgradetypes.QueryCurrentPlanRequest{}

	resp, err := c.Upgrade.CurrentPlan(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get current upgrade plan")
	}

	return resp.Plan, nil
}

// GetMuseTokenSupplyOnNode returns the muse token supply on the node
func (c *Clients) GetMuseTokenSupplyOnNode(ctx context.Context) (sdkmath.Int, error) {
	in := &banktypes.QuerySupplyOfRequest{Denom: config.BaseDenom}

	resp, err := c.Bank.SupplyOf(ctx, in)
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "failed to get muse token supply")
	}

	return resp.GetAmount().Amount, nil
}
