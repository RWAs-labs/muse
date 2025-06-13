package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/emissions/keeper"
	"github.com/RWAs-labs/muse/x/emissions/types"
)

func TestMsgServer_UpdateParams(t *testing.T) {
	t.Run("successfully update params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		res, err := msgServer.UpdateParams(ctx, &types.MsgUpdateParams{
			Authority: k.GetAuthority(),
			Params:    types.DefaultParams(),
		})

		require.NoError(t, err)
		require.Empty(t, res)
		params, found := k.GetParams(ctx)
		require.True(t, found)
		require.Equal(t, types.DefaultParams(), params)
	})

	t.Run("fail for wrong authority", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		_, err := msgServer.UpdateParams(ctx, &types.MsgUpdateParams{
			Authority: sample.AccAddress(),
			Params:    types.DefaultParams(),
		})

		require.Error(t, err)
	})

	t.Run("fail for invalid params ,validatorEmissionPercentage is invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		params := types.DefaultParams()
		params.ValidatorEmissionPercentage = "-1.5"
		_, err := msgServer.UpdateParams(ctx, &types.MsgUpdateParams{
			Authority: k.GetAuthority(),
			Params:    params,
		})

		require.ErrorIs(t, err, types.ErrUnableToSetParams)
	})

	t.Run("fail for invalid params ,pending buffer blocks is invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		params := types.DefaultParams()
		params.PendingBallotsDeletionBufferBlocks = -1
		_, err := msgServer.UpdateParams(ctx, &types.MsgUpdateParams{
			Authority: k.GetAuthority(),
			Params:    params,
		})

		require.ErrorIs(t, err, types.ErrUnableToSetParams)
	})
}
