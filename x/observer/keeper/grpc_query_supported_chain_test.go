package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/pkg/chains"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
)

func TestKeeper_SupportedChains(t *testing.T) {
	t.Run("should return supported chains", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.SupportedChains(wctx, nil)
		require.NoError(t, err)
		require.Equal(t, []chains.Chain{}, res.Chains)
	})
}
