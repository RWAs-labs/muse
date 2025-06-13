package coin_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/stretchr/testify/require"
)

func Test_AmusePerMuse(t *testing.T) {
	require.Equal(t, sdkmath.LegacyNewDec(1e18), coin.AmusePerMuse())
}

func Test_GetAmuseDecFromAmountInMuse(t *testing.T) {
	tt := []struct {
		name        string
		museAmount  string
		err         require.ErrorAssertionFunc
		amuseAmount sdkmath.LegacyDec
	}{
		{
			name:        "valid muse amount",
			museAmount:  "210000000",
			err:         require.NoError,
			amuseAmount: sdkmath.LegacyMustNewDecFromStr("210000000000000000000000000"),
		},
		{
			name:        "very high muse amount",
			museAmount:  "21000000000000000000",
			err:         require.NoError,
			amuseAmount: sdkmath.LegacyMustNewDecFromStr("21000000000000000000000000000000000000"),
		},
		{
			name:        "very low muse amount",
			museAmount:  "1",
			err:         require.NoError,
			amuseAmount: sdkmath.LegacyMustNewDecFromStr("1000000000000000000"),
		},
		{
			name:        "zero muse amount",
			museAmount:  "0",
			err:         require.NoError,
			amuseAmount: sdkmath.LegacyMustNewDecFromStr("0"),
		},
		{
			name:        "decimal muse amount",
			museAmount:  "0.1",
			err:         require.NoError,
			amuseAmount: sdkmath.LegacyMustNewDecFromStr("100000000000000000"),
		},
		{
			name:        "invalid muse amount",
			museAmount:  "%%%%%$#",
			err:         require.Error,
			amuseAmount: sdkmath.LegacyMustNewDecFromStr("0"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			amuse, err := coin.GetAmuseDecFromAmountInMuse(tc.museAmount)
			tc.err(t, err)
			if err == nil {
				require.Equal(t, tc.amuseAmount, amuse)
			}
		})
	}

}

func TestGetCoinType(t *testing.T) {
	tests := []struct {
		name    string
		coin    string
		want    coin.CoinType
		wantErr bool
	}{
		{
			name:    "valid coin type 0",
			coin:    "0",
			want:    coin.CoinType(0),
			wantErr: false,
		},
		{
			name:    "valid coin type 1",
			coin:    "1",
			want:    coin.CoinType(1),
			wantErr: false,
		},
		{
			name:    "valid coin type 2",
			coin:    "2",
			want:    coin.CoinType(2),
			wantErr: false,
		},
		{
			name:    "valid coin type 3",
			coin:    "3",
			want:    coin.CoinType(3),
			wantErr: false,
		},
		{
			name:    "invalid coin type negative",
			coin:    "-1",
			wantErr: true,
		},
		{
			name: "invalid coin type large number",
			coin: "4",
			want: coin.CoinType(4),
		},
		{
			name:    "invalid coin type non-integer",
			coin:    "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := coin.GetCoinType(tt.coin)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCoinType_SupportsRefund(t *testing.T) {
	tests := []struct {
		name string
		c    coin.CoinType
		want bool
	}{
		{"should support refund for ERC20", coin.CoinType_ERC20, true},
		{"should support refund forGas", coin.CoinType_Gas, true},
		{"should support refund forMuse", coin.CoinType_Muse, true},
		{"should not support refund forCmd", coin.CoinType_Cmd, false},
		{"should not support refund forUnknown", coin.CoinType(100), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.SupportsRefund(); got != tt.want {
				t.Errorf("FungibleTokenCoinType.SupportsRefund() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoinType_IsAsset(t *testing.T) {
	tests := []struct {
		name string
		c    coin.CoinType
		want bool
	}{
		{"Gas is asset", coin.CoinType_Gas, true},
		{"ERC20 is asset", coin.CoinType_ERC20, true},
		{"Muse is asset", coin.CoinType_Muse, true},
		{"Cmd is not asset", coin.CoinType_Cmd, false},
		{"CoinType_NoAssetCall is irrelevant and not asset", coin.CoinType_NoAssetCall, false},
		{"Unknown coin type is not asset", coin.CoinType(100), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsAsset(); got != tt.want {
				t.Errorf("CoinType.IsAsset() = %v, want %v", got, tt.want)
			}
		})
	}
}
