package precompiles

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	ethermint "github.com/RWAs-labs/ethermint/types"
	"github.com/RWAs-labs/muse/testutil/keeper"
	ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func Test_StatefulContracts(t *testing.T) {
	k, ctx, sdkk, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec

	var expectedContracts int
	for _, enabled := range EnabledStatefulContracts {
		if enabled {
			expectedContracts++
		}
	}

	// StatefulContracts() should return all the enabled contracts.
	contracts := StatefulContracts(
		k,
		&sdkk.StakingKeeper,
		sdkk.BankKeeper,
		sdkk.DistributionKeeper,
		appCodec,
		gasConfig,
	)
	require.NotNil(t, contracts, "StatefulContracts() should not return a nil slice")
	require.Len(t, contracts, expectedContracts, "StatefulContracts() should return all the enabled contracts")

	for _, customContractFn := range contracts {
		// Extract the contract function.
		contract := customContractFn(ctx, ethparams.Rules{})

		// Check the contract function returns a valid address.
		contractAddr := contract.Address()
		require.NotNil(t, contractAddr, "The called contract should have a valid address")
	}
}
