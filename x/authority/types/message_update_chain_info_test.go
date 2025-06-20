package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/authority/types"
)

func TestMsgUpdateChainInfo_ValidateBasic(t *testing.T) {
	tests := []struct {
		name        string
		msg         *types.MsgUpdateChainInfo
		errContains string
	}{
		{
			name: "valid message",
			msg:  types.NewMsgUpdateChainInfo(sample.AccAddress(), sample.Chain(42)),
		},
		{
			name:        "invalid creator address",
			msg:         types.NewMsgUpdateChainInfo("invalid", sample.Chain(42)),
			errContains: "invalid creator address",
		},
		{
			name:        "invalid chain",
			msg:         types.NewMsgUpdateChainInfo(sample.AccAddress(), chains.Chain{ChainId: -1}),
			errContains: "invalid chain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.errContains != "" {
				require.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUpdateChainInfo_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgUpdateChainInfo
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgUpdateChainInfo(signer, sample.Chain(42)),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgUpdateChainInfo("invalid", sample.Chain(42)),
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

func TestMsgUpdateChainInfo_Type(t *testing.T) {
	msg := types.NewMsgUpdateChainInfo(sample.AccAddress(), sample.Chain(42))
	require.Equal(t, types.TypeMsgUpdateChainInfo, msg.Type())
}

func TestMsgUpdateChainInfo_Route(t *testing.T) {
	msg := types.NewMsgUpdateChainInfo(sample.AccAddress(), sample.Chain(42))
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateChainInfo_GetSignBytes(t *testing.T) {
	msg := types.NewMsgUpdateChainInfo(sample.AccAddress(), sample.Chain(42))
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
