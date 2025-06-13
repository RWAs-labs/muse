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

// UpdateMRC20WithdrawFee updates the withdraw fee and gas limit of a mrc20 token
func (k msgServer) UpdateMRC20WithdrawFee(
	goCtx context.Context,
	msg *types.MsgUpdateMRC20WithdrawFee,
) (*types.MsgUpdateMRC20WithdrawFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check signer permission
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// check the mrc20 exists
	mrc20Addr := ethcommon.HexToAddress(msg.Mrc20Address)
	if mrc20Addr == (ethcommon.Address{}) {
		return nil, cosmoserrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"invalid mrc20 contract address (%s)",
			msg.Mrc20Address,
		)
	}
	coin, found := k.GetForeignCoins(ctx, msg.Mrc20Address)
	if !found {
		return nil, cosmoserrors.Wrapf(
			types.ErrForeignCoinNotFound,
			"no foreign coin match requested mrc20 address (%s)",
			msg.Mrc20Address,
		)
	}

	// get the previous fee
	oldWithdrawFee, err := k.QueryProtocolFlatFee(ctx, mrc20Addr)
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to query protocol flat fee (%s)", err.Error())
	}
	oldGasLimit, err := k.QueryGasLimit(ctx, mrc20Addr)
	if err != nil {
		return nil, cosmoserrors.Wrapf(types.ErrContractCall, "failed to query gas limit (%s)", err.Error())
	}

	// call the contract methods
	tmpCtx, commit := ctx.CacheContext()
	if !msg.NewWithdrawFee.IsNil() {
		_, err = k.UpdateMRC20ProtocolFlatFee(tmpCtx, mrc20Addr, msg.NewWithdrawFee.BigInt())
		if err != nil {
			return nil, cosmoserrors.Wrapf(
				types.ErrContractCall,
				"failed to call mrc20 contract updateProtocolFlatFee method (%s)",
				err.Error(),
			)
		}
	}
	if !msg.NewGasLimit.IsNil() {
		_, err = k.UpdateMRC20GasLimit(tmpCtx, mrc20Addr, msg.NewGasLimit.BigInt())
		if err != nil {
			return nil, cosmoserrors.Wrapf(
				types.ErrContractCall,
				"failed to call mrc20 contract updateGasLimit method (%s)",
				err.Error(),
			)
		}
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventMRC20WithdrawFeeUpdated{
			MsgTypeUrl:     sdk.MsgTypeURL(&types.MsgUpdateMRC20WithdrawFee{}),
			ChainId:        coin.ForeignChainId,
			CoinType:       coin.CoinType,
			Mrc20Address:   mrc20Addr.Hex(),
			OldWithdrawFee: oldWithdrawFee.String(),
			NewWithdrawFee: msg.NewWithdrawFee.String(),
			Signer:         msg.Creator,
			OldGasLimit:    oldGasLimit.String(),
			NewGasLimit:    msg.NewGasLimit.String(),
		},
	)
	if err != nil {
		k.Logger(ctx).Error("failed to emit event", "error", err.Error())
		return nil, cosmoserrors.Wrapf(types.ErrEmitEvent, "failed to emit event (%s)", err.Error())
	}
	commit()

	return &types.MsgUpdateMRC20WithdrawFeeResponse{}, nil
}
