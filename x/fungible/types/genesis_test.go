package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				ForeignCoinsList: []types.ForeignCoins{
					{
						Mrc20ContractAddress: "0",
					},
					{
						Mrc20ContractAddress: "1",
					},
				},
			},
			valid: true,
		},
		{
			desc: "duplicated foreignCoins",
			genState: &types.GenesisState{
				ForeignCoinsList: []types.ForeignCoins{
					{
						Mrc20ContractAddress: "0",
					},
					{
						Mrc20ContractAddress: "0",
					},
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
