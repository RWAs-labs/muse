package emissions_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/nullify"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/emissions"
	"github.com/RWAs-labs/muse/x/emissions/types"
)

func TestGenesis(t *testing.T) {
	t.Run("should init and export for valid state", func(t *testing.T) {
		params := types.DefaultParams()

		genesisState := types.GenesisState{
			Params: params,
			WithdrawableEmissions: []types.WithdrawableEmissions{
				sample.WithdrawableEmissions(t),
				sample.WithdrawableEmissions(t),
				sample.WithdrawableEmissions(t),
			},
		}

		// Init and export
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		emissions.InitGenesis(ctx, *k, genesisState)
		got := emissions.ExportGenesis(ctx, *k)
		require.NotNil(t, got)

		// Compare genesis after init and export
		nullify.Fill(&genesisState)
		nullify.Fill(got)
		require.Equal(t, genesisState, *got)
	})

	t.Run("should error for invalid params", func(t *testing.T) {
		params := types.DefaultParams()
		params.ObserverSlashAmount = sdkmath.NewInt(-1)

		genesisState := types.GenesisState{
			Params: params,
			WithdrawableEmissions: []types.WithdrawableEmissions{
				sample.WithdrawableEmissions(t),
				sample.WithdrawableEmissions(t),
				sample.WithdrawableEmissions(t),
			},
		}

		// Init and export
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		require.Panics(t, func() {
			emissions.InitGenesis(ctx, *k, genesisState)
		})
	})
}
