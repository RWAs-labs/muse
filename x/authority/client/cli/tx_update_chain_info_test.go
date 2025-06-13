package cli_test

import (
	"testing"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/testutil/testdata"
	"github.com/RWAs-labs/muse/x/authority/client/cli"
	"github.com/stretchr/testify/require"
)

func TestReadChainInfoFromFile(t *testing.T) {
	t.Run("successfully read file", func(t *testing.T) {
		fs := testdata.TypesFiles

		chain, err := cli.ReadChainFromFile(fs, "types/chain.json")
		require.NoError(t, err)

		require.EqualValues(t, chains.Chain{
			ChainId:     1,
			Network:     chains.Network_muse,
			NetworkType: chains.NetworkType_devnet,
			Vm:          chains.Vm_svm,
			Consensus:   chains.Consensus_solana_consensus,
			IsExternal:  true,
			CctxGateway: chains.CCTXGateway_mevm,
			Name:        "testchain",
		}, chain)
	})

	t.Run("file not found", func(t *testing.T) {
		fs := testdata.TypesFiles

		_, err := cli.ReadChainFromFile(fs, "types/chain_not_found.json")
		require.Error(t, err)
	})
}
