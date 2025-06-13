package types

import (
	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateMRC20LiquidityCap = "update_mrc20_liquidity_cap"

var _ sdk.Msg = &MsgUpdateMRC20LiquidityCap{}

func NewMsgUpdateMRC20LiquidityCap(creator string, mrc20 string, liquidityCap math.Uint) *MsgUpdateMRC20LiquidityCap {
	return &MsgUpdateMRC20LiquidityCap{
		Creator:      creator,
		Mrc20Address: mrc20,
		LiquidityCap: liquidityCap,
	}
}

func (msg *MsgUpdateMRC20LiquidityCap) Route() string {
	return RouterKey
}

func (msg *MsgUpdateMRC20LiquidityCap) Type() string {
	return TypeMsgUpdateMRC20LiquidityCap
}

func (msg *MsgUpdateMRC20LiquidityCap) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateMRC20LiquidityCap) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateMRC20LiquidityCap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !ethcommon.IsHexAddress(msg.Mrc20Address) {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", msg.Mrc20Address)
	}

	return nil
}
