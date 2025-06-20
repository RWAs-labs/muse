package sample

import (
	"testing"

	"cosmossdk.io/math"

	"github.com/RWAs-labs/muse/x/emissions/types"
)

func WithdrawableEmissions(t *testing.T) types.WithdrawableEmissions {
	addr := AccAddress()
	r := newRandFromStringSeed(t, addr)

	return types.WithdrawableEmissions{
		Address: AccAddress(),
		Amount:  math.NewInt(r.Int63()),
	}
}
