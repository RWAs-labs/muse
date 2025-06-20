package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/RWAs-labs/muse/x/observer/types"
)

func (k Keeper) SupportedChains(
	goCtx context.Context,
	_ *types.QuerySupportedChains,
) (*types.QuerySupportedChainsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chains := k.GetSupportedChains(ctx)
	return &types.QuerySupportedChainsResponse{Chains: chains}, nil
}
