package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const TypeMsgPauseMrc20 = "pause_mrc20"

var _ sdk.Msg = &MsgPauseMRC20{}

func NewMsgPauseMRC20(creator string, mrc20 []string) *MsgPauseMRC20 {
	return &MsgPauseMRC20{
		Creator:        creator,
		Mrc20Addresses: mrc20,
	}
}

func (msg *MsgPauseMRC20) Route() string {
	return RouterKey
}

func (msg *MsgPauseMRC20) Type() string {
	return TypeMsgPauseMrc20
}

func (msg *MsgPauseMRC20) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPauseMRC20) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPauseMRC20) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Mrc20Addresses) == 0 {
		return cosmoserrors.Wrap(sdkerrors.ErrInvalidRequest, "no mrc20 to update")
	}

	// check if all mrc20 addresses are valid
	for _, mrc20 := range msg.Mrc20Addresses {
		if !ethcommon.IsHexAddress(mrc20) {
			return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid mrc20 contract address (%s)", mrc20)
		}
	}
	return nil
}
