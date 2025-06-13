package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// PauseMRC20 pauses a list of MRC20 tokens
// Authorized: admin policy group groupEmergency.
func (k msgServer) PauseMRC20(
	goCtx context.Context,
	msg *types.MsgPauseMRC20,
) (*types.MsgPauseMRC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// iterate all foreign coins and set paused status
	for _, mrc20 := range msg.Mrc20Addresses {
		fc, found := k.GetForeignCoins(ctx, mrc20)
		if !found {
			return nil, cosmoserrors.Wrapf(types.ErrForeignCoinNotFound, "foreign coin not found %s", mrc20)
		}
		// Set status to paused
		fc.Paused = true
		k.SetForeignCoins(ctx, fc)
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventMRC20Paused{
			MsgTypeUrl:     sdk.MsgTypeURL(&types.MsgPauseMRC20{}),
			Mrc20Addresses: msg.Mrc20Addresses,
			Signer:         msg.Creator,
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event",
			"event", "EventMRC20Paused",
			"error", err.Error(),
		)
		return nil, cosmoserrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}

	return &types.MsgPauseMRC20Response{}, nil
}
