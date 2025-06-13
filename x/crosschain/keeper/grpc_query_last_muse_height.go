package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RWAs-labs/muse/x/crosschain/types"
)

func (k Keeper) LastMuseHeight(
	goCtx context.Context,
	req *types.QueryLastMuseHeightRequest,
) (*types.QueryLastMuseHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	height := ctx.BlockHeight()
	if height < 0 {
		return nil, status.Error(codes.OutOfRange, "height out of range")
	}
	return &types.QueryLastMuseHeightResponse{
		Height: height,
	}, nil
}
