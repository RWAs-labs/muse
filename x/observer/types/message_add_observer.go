package types

import (
	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/RWAs-labs/muse/pkg/crypto"
)

const TypeMsgAddObserver = "add_observer"

var _ sdk.Msg = &MsgAddObserver{}

func NewMsgAddObserver(
	creator string,
	observerAdresss string,
	museclientGranteePubKey string,
	addNodeAccountOnly bool,
) *MsgAddObserver {
	return &MsgAddObserver{
		Creator:                 creator,
		ObserverAddress:         observerAdresss,
		MuseclientGranteePubkey: museclientGranteePubKey,
		AddNodeAccountOnly:      addNodeAccountOnly,
	}
}

func (msg *MsgAddObserver) Route() string {
	return RouterKey
}

func (msg *MsgAddObserver) Type() string {
	return TypeMsgAddObserver
}

func (msg *MsgAddObserver) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddObserver) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddObserver) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.ObserverAddress)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid observer address (%s)", err)
	}
	_, err = crypto.NewPubKey(msg.MuseclientGranteePubkey)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid museclient grantee pubkey (%s)", err)
	}
	_, err = crypto.GetAddressFromPubkeyString(msg.MuseclientGranteePubkey)
	if err != nil {
		return cosmoserrors.Wrapf(sdkerrors.ErrInvalidPubKey, "invalid museclient grantee pubkey (%s)", err)
	}
	return nil
}
