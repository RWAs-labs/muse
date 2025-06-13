package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestNewMsgUpdateMRC20Name_ValidateBasics(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdateMRC20Name
		err  error
	}{
		{
			name: "valid message",
			msg: types.NewMsgUpdateMRC20Name(
				sample.AccAddress(),
				sample.EthAddress().String(),
				"foo",
				"bar",
			),
		},
		{
			name: "valid message with empty name",
			msg: types.NewMsgUpdateMRC20Name(
				sample.AccAddress(),
				sample.EthAddress().String(),
				"",
				"bar",
			),
		},
		{
			name: "valid message with empty symbol",
			msg: types.NewMsgUpdateMRC20Name(
				sample.AccAddress(),
				sample.EthAddress().String(),
				"foo",
				"",
			),
		},

		{
			name: "invalid address",
			msg: types.NewMsgUpdateMRC20Name(
				"invalid_address",
				sample.EthAddress().String(),
				"foo",
				"bar",
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid contract address",
			msg: types.NewMsgUpdateMRC20Name(
				sample.AccAddress(),
				"invalid_address",
				"foo",
				"bar",
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "nothing to update",
			msg: types.NewMsgUpdateMRC20Name(
				sample.AccAddress(),
				sample.EthAddress().String(),
				"",
				"",
			),
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewMsgUpdateMRC20Name_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateMRC20Name
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateMRC20Name{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateMRC20Name{
				Creator: "invalid",
			},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				require.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestNewMsgUpdateMRC20Name_Type(t *testing.T) {
	msg := types.MsgUpdateMRC20Name{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateMRC20Name, msg.Type())
}

func TestNewMsgUpdateMRC20Name_Route(t *testing.T) {
	msg := types.MsgUpdateMRC20Name{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestNewMsgUpdateMRC20Name_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateMRC20Name{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
