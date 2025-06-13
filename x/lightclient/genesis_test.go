package lightclient_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/proofs"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/nullify"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/lightclient"
	"github.com/RWAs-labs/muse/x/lightclient/types"
)

func TestGenesis(t *testing.T) {
	t.Run("can import and export genesis", func(t *testing.T) {
		genesisState := types.GenesisState{
			BlockHeaderVerification: sample.BlockHeaderVerification(),
			BlockHeaders: []proofs.BlockHeader{
				sample.BlockHeader(sample.Hash().Bytes()),
				sample.BlockHeader(sample.Hash().Bytes()),
				sample.BlockHeader(sample.Hash().Bytes()),
			},
			ChainStates: []types.ChainState{
				sample.ChainState(chains.Ethereum.ChainId),
				sample.ChainState(chains.BitcoinMainnet.ChainId),
				sample.ChainState(chains.BscMainnet.ChainId),
			},
		}

		// Init and export
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		lightclient.InitGenesis(ctx, *k, genesisState)
		got := lightclient.ExportGenesis(ctx, *k)
		require.NotNil(t, got)

		// Compare genesis after init and export
		nullify.Fill(&genesisState)
		nullify.Fill(got)
		require.Equal(t, genesisState, *got)
	})

	t.Run("can export genesis with empty state", func(t *testing.T) {
		// Export genesis with empty state
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		got := lightclient.ExportGenesis(ctx, *k)
		require.NotNil(t, got)

		// Compare genesis after export
		expected := types.GenesisState{
			BlockHeaderVerification: types.DefaultBlockHeaderVerification(),
			BlockHeaders:            []proofs.BlockHeader(nil),
			ChainStates:             []types.ChainState(nil),
		}
		require.Equal(t, expected, *got)
		require.Equal(
			t,
			expected.BlockHeaderVerification.HeaderSupportedChains,
			got.BlockHeaderVerification.HeaderSupportedChains,
		)
	})
}
