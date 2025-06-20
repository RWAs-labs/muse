package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestGetGasStabilityPoolAddress(t *testing.T) {
	address := types.GasStabilityPoolAddress()
	require.False(t, address.Empty())
}

func TestGetGasStabilityPoolAddressEVM(t *testing.T) {
	address := types.GasStabilityPoolAddressEVM()
	require.NotEmpty(t, address)
}
