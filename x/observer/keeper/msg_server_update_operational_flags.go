package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/observer/types"
)

func (k msgServer) UpdateOperationalFlags(
	goCtx context.Context,
	msg *types.MsgUpdateOperationalFlags,
) (*types.MsgUpdateOperationalFlagsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	k.Keeper.SetOperationalFlags(ctx, msg.OperationalFlags)

	return &types.MsgUpdateOperationalFlagsResponse{}, nil
}
