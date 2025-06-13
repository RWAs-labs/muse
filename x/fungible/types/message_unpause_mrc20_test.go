package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestMsgUnpauseMRC20_ValidateBasic(t *testing.T) {
	tt := []struct {
		name    string
		msg     *types.MsgUnpauseMRC20
		wantErr bool
	}{
		{
			name: "valid unpause message",
			msg: types.NewMsgUnpauseMRC20(
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
			msg: types.NewMsgUnpauseMRC20(
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
			msg: types.NewMsgUnpauseMRC20(
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
			msg: types.NewMsgUnpauseMRC20(
				sample.AccAddress(),
				[]string{},
			),
			wantErr: true,
		},
		{
			name: "invalid mrc20 address",
			msg: types.NewMsgUnpauseMRC20(
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

func TestMsgUnpauseMRC20_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUnpauseMRC20
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUnpauseMRC20{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUnpauseMRC20{
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

func TestMsgUnpauseMRC20_Type(t *testing.T) {
	msg := types.MsgUnpauseMRC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUnpauseMRC20, msg.Type())
}

func TestMsgUnpauseMRC20_Route(t *testing.T) {
	msg := types.MsgUnpauseMRC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUnpauseMRC20_GetSignBytes(t *testing.T) {
	msg := types.MsgUnpauseMRC20{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
