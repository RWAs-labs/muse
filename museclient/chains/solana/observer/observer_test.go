package observer_test

import (
	"testing"

	"github.com/RWAs-labs/muse/museclient/db"
	"github.com/RWAs-labs/muse/museclient/keys"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/chains/solana/observer"
	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/testutil/sample"
)

// MockSolanaObserver creates a mock Solana observer with custom chain, TSS, params etc
func MockSolanaObserver(
	t *testing.T,
	chain chains.Chain,
	solClient interfaces.SolanaRPCClient,
	chainParams observertypes.ChainParams,
	musecoreClient interfaces.MusecoreClient,
	tss interfaces.TSSSigner,
) *observer.Observer {
	// use mock musecore client if not provided
	if musecoreClient == nil {
		musecoreClient = mocks.NewMusecoreClient(t).WithKeys(&keys.Keys{})
	}

	// use mock tss if not provided
	if tss == nil {
		tss = mocks.NewTSS(t)
	}

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	baseObserver, err := base.NewObserver(
		chain,
		chainParams,
		musecoreClient,
		tss,
		1000,
		nil,
		database,
		base.DefaultLogger(),
	)
	require.NoError(t, err)

	ob, err := observer.New(baseObserver, solClient, chainParams.GatewayAddress)
	require.NoError(t, err)

	return ob
}

func Test_LoadLastTxScanned(t *testing.T) {
	// prepare params
	chain := chains.SolanaDevnet
	params := sample.ChainParams(chain.ChainId)
	params.GatewayAddress = sample.SolanaAddress(t)

	// create observer
	ob := MockSolanaObserver(t, chain, nil, *params, nil, nil)

	t.Run("should load last block scanned", func(t *testing.T) {
		// write sample last tx to db
		lastTx := sample.SolanaSignature(t).String()
		ob.WriteLastTxScannedToDB(lastTx)

		// load last tx scanned
		err := ob.LoadLastTxScanned()
		require.NoError(t, err)
		require.Equal(t, lastTx, ob.LastTxScanned())
	})
}
