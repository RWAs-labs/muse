package simulation_test

import (
	"fmt"
	"testing"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/simulation"
	"github.com/RWAs-labs/muse/x/fungible/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/stretchr/testify/require"
)

func TestDecodeStore(t *testing.T) {
	k, _, _, _ := keepertest.FungibleKeeper(t)
	cdc := k.GetCodec()
	dec := simulation.NewDecodeStore(cdc)
	systemContract := sample.SystemContract()
	foreignCoins := sample.ForeignCoins(t, sample.EthAddress().String())

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: []byte(types.SystemContractKey), Value: cdc.MustMarshal(systemContract)},
			{Key: []byte(types.ForeignCoinsKeyPrefix), Value: cdc.MustMarshal(&foreignCoins)},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"SystemContract", fmt.Sprintf("%v\n%v", *systemContract, *systemContract)},
		{"ForeignCoins", fmt.Sprintf("%v\n%v", foreignCoins, foreignCoins)},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]))
		})
	}
}
