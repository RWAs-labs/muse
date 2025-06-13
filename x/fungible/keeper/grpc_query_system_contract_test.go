package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestKeeper_SystemContract(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		res, err := k.SystemContract(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if system contract not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		res, err := k.SystemContract(ctx, &types.QueryGetSystemContractRequest{})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return system contract if found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		sc := types.SystemContract{
			SystemContract: sample.EthAddress().Hex(),
			ConnectorMevm:  sample.EthAddress().Hex(),
		}
		k.SetSystemContract(ctx, sc)
		res, err := k.SystemContract(ctx, &types.QueryGetSystemContractRequest{})
		require.NoError(t, err)
		require.Equal(t, &types.QueryGetSystemContractResponse{
			SystemContract: sc,
		}, res)
	})
}
