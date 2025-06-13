package types

import (
	cosmoserror "cosmossdk.io/errors"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgUpdateMRC20WithdrawFee = "update_mrc20_withdraw_fee"

var _ sdk.Msg = &MsgUpdateMRC20WithdrawFee{}

func NewMsgUpdateMRC20WithdrawFee(
	creator string,
	mrc20 string,
	newFee math.Uint,
	newGasLimit math.Uint,
) *MsgUpdateMRC20WithdrawFee {
	return &MsgUpdateMRC20WithdrawFee{
		Creator:        creator,
		Mrc20Address:   mrc20,
		NewWithdrawFee: newFee,
		NewGasLimit:    newGasLimit,
	}
}

func (msg *MsgUpdateMRC20WithdrawFee) Route() string {
	return RouterKey
}

func (msg *MsgUpdateMRC20WithdrawFee) Type() string {
	return TypeMsgUpdateMRC20WithdrawFee
}

func (msg *MsgUpdateMRC20WithdrawFee) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateMRC20WithdrawFee) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateMRC20WithdrawFee) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	// check if the system contract address is valid
	if !ethcommon.IsHexAddress(msg.Mrc20Address) {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidAddress, "invalid system contract address (%s)", msg.Mrc20Address)
	}
	if msg.NewWithdrawFee.IsNil() && msg.NewGasLimit.IsNil() {
		return cosmoserror.Wrapf(sdkerrors.ErrInvalidRequest, "nothing to update")
	}

	return nil
}
