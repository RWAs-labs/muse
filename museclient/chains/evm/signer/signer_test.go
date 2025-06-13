package signer

import (
	"context"
	"math/big"
	"testing"

	"github.com/RWAs-labs/muse/museclient/chains/evm/client"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/museclient/keys"
	"github.com/RWAs-labs/muse/museclient/testutils/testlog"
	"github.com/RWAs-labs/muse/museclient/testutils/testrpc"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/museclient/testutils"
	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/testutil/sample"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

var (
	// Dummy addresses as they are just used as transaction data to be signed
	ConnectorAddress    = sample.EthAddress()
	ERC20CustodyAddress = sample.EthAddress()
)

type testSuite struct {
	*Signer
	tss       *mocks.TSS
	evmServer *testrpc.EVMServer
	client    *client.Client
}

func newTestSuite(t *testing.T) *testSuite {
	ctx := context.Background()

	chain := chains.BscMainnet

	evmServer := testrpc.NewEVMServer(t)

	evmServer.SetChainID(int(chain.ChainId))
	evmServer.MockSendTransaction()

	evmClient, err := client.NewFromEndpoint(ctx, evmServer.Endpoint)
	require.NoError(t, err)

	tss := mocks.NewTSS(t)

	logger := testlog.New(t)

	baseSigner := base.NewSigner(chain, tss, base.Logger{Std: logger.Logger, Compliance: logger.Logger})

	s, err := New(
		baseSigner,
		evmClient,
		ConnectorAddress,
		ERC20CustodyAddress,
		sample.EthAddress(),
	)
	require.NoError(t, err)

	return &testSuite{
		Signer:    s,
		tss:       tss,
		evmServer: evmServer,
		client:    evmClient,
	}
}

func (ts *testSuite) EvmSigner() ethtypes.Signer {
	return ts.client.Signer
}

func getCCTX(t *testing.T) *crosschaintypes.CrossChainTx {
	return testutils.LoadCctxByNonce(t, 56, 68270)
}

func getInvalidCCTX(t *testing.T) *crosschaintypes.CrossChainTx {
	cctx := getCCTX(t)
	// modify receiver chain id to make it invalid
	cctx.GetCurrentOutboundParam().ReceiverChainId = 13378337
	return cctx
}

// verifyTxSender is a helper function to verify the signature of a transaction
//
// signer.Sender() will ecrecover the public key of the transaction internally
// and will fail if the transaction is not valid or has been tampered with
func verifyTxSender(t *testing.T, tx *ethtypes.Transaction, expectedSender ethcommon.Address, signer ethtypes.Signer) {
	senderAddr, err := signer.Sender(tx)
	require.NoError(t, err)
	require.Equal(t, expectedSender.String(), senderAddr.String())
}

// verifyTxBodyBasics is a helper function to verify 'to', 'nonce' and 'amount' of a transaction
func verifyTxBodyBasics(
	t *testing.T,
	tx *ethtypes.Transaction,
	to ethcommon.Address,
	nonce uint64,
	amount *big.Int,
) {
	require.Equal(t, to, *tx.To())
	require.Equal(t, nonce, tx.Nonce())
	require.True(t, amount.Cmp(tx.Value()) == 0)
}

func TestSigner_SetGetConnectorAddress(t *testing.T) {
	evmSigner := newTestSuite(t)

	// Get and compare
	require.Equal(t, ConnectorAddress, evmSigner.GetMuseConnectorAddress())

	// Update and get again
	newConnector := sample.EthAddress()
	evmSigner.SetMuseConnectorAddress(newConnector)
	require.Equal(t, newConnector, evmSigner.GetMuseConnectorAddress())
}

func TestSigner_SetGetERC20CustodyAddress(t *testing.T) {
	evmSigner := newTestSuite(t)
	// Get and compare
	require.Equal(t, ERC20CustodyAddress, evmSigner.GetERC20CustodyAddress())

	// Update and get again
	newCustody := sample.EthAddress()
	evmSigner.SetERC20CustodyAddress(newCustody)
	require.Equal(t, newCustody, evmSigner.GetERC20CustodyAddress())
}

func TestSigner_TryProcessOutbound(t *testing.T) {
	ctx := makeCtx(t)

	// ARRANGE
	// Setup evm signer
	evmSigner := newTestSuite(t)
	cctx := getCCTX(t)

	// Test with mock client that has keys
	client := mocks.NewMusecoreClient(t).
		WithKeys(&keys.Keys{}).
		WithMuseChain().
		WithPostVoteOutbound("", "")

	// mock evm client "NonceAt"
	nonce := uint64(123)
	evmSigner.evmServer.MockNonceAt(nonce)

	// ACT
	evmSigner.TryProcessOutbound(ctx, cctx, client, nonce)

	// ASSERT
	// Check if cctx was signed and broadcasted
	list := evmSigner.GetReportedTxList()
	require.Len(t, *list, 1)
}

func TestSigner_BroadcastOutbound(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner := newTestSuite(t)

	// Setup txData struct
	cctx := getCCTX(t)
	nonce := uint64(123)
	txData, skip, err := NewOutboundData(ctx, cctx, nonce, zerolog.Logger{})
	require.NoError(t, err)
	require.False(t, skip)

	// Mock evm client "NonceAt"
	evmSigner.evmServer.MockNonceAt(nonce)

	t.Run("BroadcastOutbound - should successfully broadcast", func(t *testing.T) {
		// Call SignERC20Withdraw
		tx, err := evmSigner.SignERC20Withdraw(ctx, txData)
		require.NoError(t, err)

		evmSigner.BroadcastOutbound(
			ctx,
			tx,
			cctx,
			zerolog.Logger{},
			mocks.NewMusecoreClient(t),
			txData,
		)

		//Check if cctx was signed and broadcasted
		list := evmSigner.GetReportedTxList()
		require.Len(t, *list, 1)
	})
}

func TestSigner_SignerErrorMsg(t *testing.T) {
	cctx := getCCTX(t)

	msg := ErrorMsg(cctx)
	require.Contains(t, msg, "nonce 68270 chain 56")
}

func makeCtx(t *testing.T) context.Context {
	app := zctx.New(config.New(false), nil, zerolog.Nop())

	bscParams := mocks.MockChainParams(chains.BscMainnet.ChainId, 10)

	err := app.Update(
		[]chains.Chain{chains.BscMainnet, chains.MuseChainMainnet},
		nil,
		map[int64]*observertypes.ChainParams{
			chains.BscMainnet.ChainId: &bscParams,
		},
		observertypes.CrosschainFlags{},
		observertypes.OperationalFlags{},
	)
	require.NoError(t, err, "unable to update app context")

	return zctx.WithAppContext(context.Background(), app)
}

func makeLogger(t *testing.T) zerolog.Logger {
	return zerolog.New(zerolog.NewTestWriter(t))
}
