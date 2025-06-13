package types_test

import (
	"errors"
	"fmt"
	"testing"

	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestIsContractReverted(t *testing.T) {
	require.True(t, types.IsContractReverted(nil, vm.ErrExecutionReverted))
	require.True(t, types.IsContractReverted(nil, fmt.Errorf("foo : %s", vm.ErrExecutionReverted.Error())))
	require.True(t, types.IsContractReverted(&evmtypes.MsgEthereumTxResponse{VmError: "foo"}, nil))

	require.False(t, types.IsContractReverted(nil, nil))
	require.False(t, types.IsContractReverted(nil, errors.New("foo")))
	require.False(t, types.IsContractReverted(&evmtypes.MsgEthereumTxResponse{VmError: ""}, nil))
}
