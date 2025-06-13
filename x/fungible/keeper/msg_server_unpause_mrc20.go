package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// UnpauseMRC20 unpauses the MRC20 token
// Authorized: admin policy group groupOperational.
func (k msgServer) UnpauseMRC20(
	goCtx context.Context,
	msg *types.MsgUnpauseMRC20,
) (*types.MsgUnpauseMRC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// iterate all foreign coins and set unpaused status
	for _, mrc20 := range msg.Mrc20Addresses {
		fc, found := k.GetForeignCoins(ctx, mrc20)
		if !found {
			return nil, cosmoserrors.Wrapf(types.ErrForeignCoinNotFound, "foreign coin not found %s", mrc20)
		}
		// Set status to unpaused
		fc.Paused = false
		k.SetForeignCoins(ctx, fc)
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventMRC20Unpaused{
			MsgTypeUrl:     sdk.MsgTypeURL(&types.MsgUnpauseMRC20{}),
			Mrc20Addresses: msg.Mrc20Addresses,
			Signer:         msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event",
			"event", "EventMRC20Unpaused",
			"error", err.Error(),
		)
		return nil, cosmoserrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}

	return &types.MsgUnpauseMRC20Response{}, nil
}
