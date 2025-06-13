package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestMsgPauseMRC20_ValidateBasic(t *testing.T) {
	tt := []struct {
		name    string
		msg     *types.MsgPauseMRC20
		wantErr bool
	}{
		{
			name: "valid pause message",
			msg: types.NewMsgPauseMRC20(
				sample.AccAddress(),
				[]string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
			),
			wantErr: false,
		},
		{
			name: "valid unpause message",
			msg: types.NewMsgPauseMRC20(
				sample.AccAddress(),
				[]string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
			),
			wantErr: false,
		},
		{
			name: "invalid creator address",
			msg: types.NewMsgPauseMRC20(
				"invalid",
				[]string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
			),
			wantErr: true,
		},
		{
			name: "invalid empty mrc20 address",
			msg: types.NewMsgPauseMRC20(
				sample.AccAddress(),
				[]string{},
			),
			wantErr: true,
		},
		{
			name: "invalid mrc20 address",
			msg: types.NewMsgPauseMRC20(
				sample.AccAddress(),
				[]string{
					sample.EthAddress().String(),
					"invalid",
					sample.EthAddress().String(),
				},
			),
			wantErr: true,
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgPauseMRC20_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgPauseMRC20
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgPauseMRC20{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgPauseMRC20{
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

func TestMsgPauseMRC20_Type(t *testing.T) {
	msg := types.MsgPauseMRC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgPauseMrc20, msg.Type())
}

func TestMsgPauseMRC20_Route(t *testing.T) {
	msg := types.MsgPauseMRC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgPauseMRC20_GetSignBytes(t *testing.T) {
	msg := types.MsgPauseMRC20{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
