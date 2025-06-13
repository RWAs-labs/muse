package simulation_test

import (
	"fmt"
	"testing"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/emissions/simulation"
	"github.com/RWAs-labs/muse/x/emissions/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"
)

func TestDecodeStore(t *testing.T) {
	k, _, _, _ := keepertest.EmissionsKeeper(t)
	cdc := k.GetCodec()
	dec := simulation.NewDecodeStore(cdc)
	withdrawableEmissions := sample.WithdrawableEmissions(t)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.KeyPrefix(types.WithdrawableEmissionsKey), Value: cdc.MustMarshal(&withdrawableEmissions)},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{
			"withdrawable emissions",
			fmt.Sprintf(
				"key %s value A %v value B %v",
				types.WithdrawableEmissionsKey,
				withdrawableEmissions,
				withdrawableEmissions,
			),
		},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]))
		})
	}
}
