package context

import (
	"testing"

	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	"github.com/RWAs-labs/muse/pkg/chains"
	observer "github.com/RWAs-labs/muse/x/observer/types"
	"github.com/stretchr/testify/require"
)

func TestChainRegistry(t *testing.T) {
	// Given chains & chainParams
	var (
		btc       = &chains.BitcoinMainnet
		btcParams = makeParams(btc.ChainId, true)

		eth       = &chains.Ethereum
		ethParams = makeParams(eth.ChainId, true)

		matic       = &chains.Polygon
		maticParams = makeParams(matic.ChainId, true)

		// NOT supported!
		opt       = &chains.OptimismSepolia
		optParams = makeParams(opt.ChainId, false)

		sol       = &chains.SolanaMainnet
		solParams = makeParams(sol.ChainId, true)

		// Musechain itself
		muse       = &chains.MuseChainMainnet
		museParams = makeParams(muse.ChainId, true)
	)

	t.Run("Sample Flow", func(t *testing.T) {
		// Given registry
		r := NewChainRegistry(nil)

		// With some chains added
		require.NoError(t, r.Set(btc.ChainId, btc, btcParams))
		require.NoError(t, r.Set(eth.ChainId, eth, ethParams))
		require.NoError(t, r.Set(matic.ChainId, matic, maticParams))
		require.NoError(t, r.Set(sol.ChainId, sol, solParams))
		require.NoError(t, r.Set(muse.ChainId, muse, museParams))

		// With failures on invalid data
		require.Error(t, r.Set(0, btc, btcParams))
		require.Error(t, r.Set(btc.ChainId, btc, nil))
		require.Error(t, r.Set(btc.ChainId, nil, btcParams))
		require.Error(t, r.Set(123, btc, btcParams))
		require.Error(t, r.Set(btc.ChainId, btc, ethParams))

		// With failure on adding unsupported chains
		require.ErrorIs(t, r.Set(opt.ChainId, opt, optParams), ErrChainNotSupported)

		// Should return a proper chain list
		expectedChains := []int64{
			btc.ChainId,
			eth.ChainId,
			matic.ChainId,
			sol.ChainId,
			muse.ChainId,
		}

		require.ElementsMatch(t, expectedChains, r.ChainIDs())

		// Should return not found error
		_, err := r.Get(123)
		require.ErrorIs(t, err, ErrChainNotFound)

		// Let's check ETH
		ethChain, err := r.Get(eth.ChainId)
		require.NoError(t, err)
		require.True(t, ethChain.IsEVM())
		require.False(t, ethChain.IsBitcoin())
		require.False(t, ethChain.IsSolana())
		require.Equal(t, ethParams, ethChain.Params())
	})
}

func makeParams(id int64, supported bool) *observer.ChainParams {
	cp := mocks.MockChainParams(id, 123)
	cp.IsSupported = supported

	return &cp
}
