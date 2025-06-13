package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// UpdateMRC20LiquidityCap updates the liquidity cap for a MRC20 token.
//
// Authorized: admin policy group 2.
func (k msgServer) UpdateMRC20LiquidityCap(
	goCtx context.Context,
	msg *types.MsgUpdateMRC20LiquidityCap,
) (*types.MsgUpdateMRC20LiquidityCapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check authorization
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// fetch the foreign coin
	coin, found := k.GetForeignCoins(ctx, msg.Mrc20Address)
	if !found {
		return nil, types.ErrForeignCoinNotFound
	}

	// update the liquidity cap
	coin.LiquidityCap = msg.LiquidityCap
	k.SetForeignCoins(ctx, coin)

	return &types.MsgUpdateMRC20LiquidityCapResponse{}, nil
}
