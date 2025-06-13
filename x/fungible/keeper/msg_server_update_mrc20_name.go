package keeper

import (
	"context"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"

	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// UpdateMRC20Name updates the name and/or the symbol of a mrc20 token
func (k msgServer) UpdateMRC20Name(
	goCtx context.Context,
	msg *types.MsgUpdateMRC20Name,
) (*types.MsgUpdateMRC20NameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check signer permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// check the mrc20 is valid
	mrc20Addr := ethcommon.HexToAddress(msg.Mrc20Address)
	if mrc20Addr == (ethcommon.Address{}) {
		return nil, cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid mrc20 contract address (%s)",
			msg.Mrc20Address,
		)
	}

	// check the mrc20 exists
	fc, found := k.GetForeignCoins(ctx, msg.Mrc20Address)
	if !found {
		return nil, cosmoserrors.Wrapf(
			types.ErrForeignCoinNotFound,
			"no foreign coin match requested mrc20 address (%s)",
			msg.Mrc20Address,
		)
	}

	// call the contract methods and update the object
	if msg.Name != "" {
		if err := k.MRC20SetName(ctx, mrc20Addr, msg.Name); err != nil {
			return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to update mrc20 name (%s)", err.Error())
		}
		fc.Name = msg.Name
	}

	if msg.Symbol != "" {
		if err = k.MRC20SetSymbol(ctx, mrc20Addr, msg.Symbol); err != nil {
			return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to update mrc20 symbol (%s)", err.Error())
		}
		fc.Symbol = msg.Symbol
	}

	// save the object
	k.SetForeignCoins(ctx, fc)

	return &types.MsgUpdateMRC20NameResponse{}, nil
}
