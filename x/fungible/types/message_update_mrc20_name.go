package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateMRC20Name = "update_mrc20_name"

var _ sdk.Msg = &MsgUpdateMRC20Name{}

func NewMsgUpdateMRC20Name(creator, mrc20, name, symbol string) *MsgUpdateMRC20Name {
	return &MsgUpdateMRC20Name{
		Creator:      creator,
		Mrc20Address: mrc20,
		Name:         name,
		Symbol:       symbol,
	}
}

func (msg *MsgUpdateMRC20Name) Route() string {
	return RouterKey
}

func (msg *MsgUpdateMRC20Name) Type() string {
	return TypeMsgUpdateMRC20Name
}

func (msg *MsgUpdateMRC20Name) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateMRC20Name) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateMRC20Name) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !ethcommon.IsHexAddress(msg.Mrc20Address) {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", msg.Mrc20Address)
	}

	if msg.Name == "" && msg.Symbol == "" {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidRequest, "nothing to update")
	}

	return nil
}
