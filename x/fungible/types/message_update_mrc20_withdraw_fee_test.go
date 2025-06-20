package types_test

import (
	"testing"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestMsgUpdateMRC20WithdrawFee_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdateMRC20WithdrawFee
		err  error
	}{
		{
			name: "invalid address",
			msg: types.NewMsgUpdateMRC20WithdrawFee(
				"invalid_address",
				sample.EthAddress().String(),
				math.NewUint(1),
				math.Uint{},
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new system contract address",
			msg: types.NewMsgUpdateMRC20WithdrawFee(
				sample.AccAddress(),
				"invalid_address",
				math.NewUint(1),
				math.Uint{},
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "both withdraw fee and gas limit nil",
			msg: types.NewMsgUpdateMRC20WithdrawFee(
				sample.AccAddress(),
				sample.EthAddress().String(),
				math.Uint{},
				math.Uint{},
			),
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid message",
			msg: types.NewMsgUpdateMRC20WithdrawFee(
				sample.AccAddress(),
				sample.EthAddress().String(),
				math.NewUint(42),
				math.NewUint(42),
			),
		},
		{
			name: "withdraw fee can be zero",
			msg: types.NewMsgUpdateMRC20WithdrawFee(
				sample.AccAddress(),
				sample.EthAddress().String(),
				math.ZeroUint(),
				math.NewUint(42),
			),
		},
		{
			name: "gas limit can be zero",
			msg: types.NewMsgUpdateMRC20WithdrawFee(
				sample.AccAddress(),
				sample.EthAddress().String(),
				math.ZeroUint(),
				math.NewUint(42),
			),
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

func TestMsgUpdateMRC20WithdrawFee_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateMRC20WithdrawFee
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateMRC20WithdrawFee{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateMRC20WithdrawFee{
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

func TestMsgUpdateMRC20WithdrawFee_Type(t *testing.T) {
	msg := types.MsgUpdateMRC20WithdrawFee{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateMRC20WithdrawFee, msg.Type())
}

func TestMsgUpdateMRC20WithdrawFee_Route(t *testing.T) {
	msg := types.MsgUpdateMRC20WithdrawFee{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateMRC20WithdrawFee_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateMRC20WithdrawFee{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
