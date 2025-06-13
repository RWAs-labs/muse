package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/x/emissions/types"
)

func (k Keeper) ShowAvailableEmissions(
	goCtx context.Context,
	req *types.QueryShowAvailableEmissionsRequest,
) (*types.QueryShowAvailableEmissionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	emissions, found := k.GetWithdrawableEmission(ctx, req.Address)
	if !found {
		return &types.QueryShowAvailableEmissionsResponse{
			Amount: sdk.NewCoin(config.BaseDenom, sdkmath.ZeroInt()).String(),
		}, nil
	}
	return &types.QueryShowAvailableEmissionsResponse{
		Amount: sdk.NewCoin(config.BaseDenom, emissions.Amount).String(),
	}, nil
}
