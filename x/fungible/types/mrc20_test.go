package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestNewMRC20Data(t *testing.T) {
	mrc20 := types.NewMRC20Data("name", "symbol", 8)
	require.Equal(t, "name", mrc20.Name)
	require.Equal(t, "symbol", mrc20.Symbol)
	require.Equal(t, uint8(8), mrc20.Decimals)
}
