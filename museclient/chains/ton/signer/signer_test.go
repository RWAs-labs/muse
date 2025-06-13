package signer

import (
	"context"
	"encoding/hex"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/ton/liteapi"
	"github.com/RWAs-labs/muse/museclient/keys"
	"github.com/RWAs-labs/muse/museclient/testutils"
	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	toncontracts "github.com/RWAs-labs/muse/pkg/contracts/ton"
	"github.com/RWAs-labs/muse/testutil/sample"
	cc "github.com/RWAs-labs/muse/x/crosschain/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"
)

func TestSigner(t *testing.T) {
	// ARRANGE
	ts := newTestSuite(t)

	// Given TON signer
	signer := New(ts.baseSigner, ts.liteClient, ts.gw)

	// Given a sample TON receiver
	receiver := ton.MustParseAccountID("0QAyaVdkvWSuax8luWhDXY_0X9Am1ASWlJz4OI7M-jqcM5wK")

	const (
		museHeight = 123
		outboundID = "abc123"
		nonce      = 2
	)

	amount := tonCoins(t, "5")

	// Given CCTX
	cctx := sample.CrossChainTx(t, "123")
	cctx.InboundParams.CoinType = coin.CoinType_Gas
	cctx.OutboundParams = []*cc.OutboundParams{{
		Receiver:        receiver.ToRaw(),
		ReceiverChainId: ts.chain.ChainId,
		CoinType:        coin.CoinType_Gas,
		Amount:          amount,
		TssNonce:        nonce,
	}}

	// Given expected withdrawal
	withdrawal := toncontracts.Withdrawal{
		Recipient: receiver,
		Amount:    amount,
		Seqno:     nonce,
	}

	ts.Sign(&withdrawal)

	// Given expected liteapi calls
	lt, hash := uint64(400), decodeHash(t, "df8a01053f50a74503dffe6802f357bf0e665bd1f3d082faccfebdea93cddfeb")
	ts.OnGetAccountState(ts.gw.AccountID(), tlb.ShardAccount{LastTransLt: lt, LastTransHash: hash})

	ts.OnSendMessage(0, nil)

	withdrawalTX := sample.TONWithdrawal(t, ts.gw.AccountID(), withdrawal)
	ts.OnGetTransactionsSince(ts.gw.AccountID(), lt, ton.Bits256(hash), []ton.Transaction{withdrawalTX}, nil)

	// ACT
	signer.TryProcessOutbound(ts.ctx, cctx, ts.musecore, museHeight)

	// ASSERT
	// Make sure signer send the tx the chain AND published the outbound tracker
	require.Len(t, ts.trackerBag, 1)

	tracker := ts.trackerBag[0]

	require.Equal(t, uint64(nonce), tracker.nonce)
	require.Equal(t, liteapi.TransactionToHashString(withdrawalTX), tracker.hash)
}

func TestExitCodeRegex(t *testing.T) {
	for _, tt := range []string{
		`unable to send external message: error code: 0 message: 
		cannot apply external message to current state : 
		External message was not accepted\nCannot run message on account: inbound external message rejected by 
		transaction CC8803E21EDA7E6487D191380725A82CD75316E1C131496E1A5636751CE60347:
		\nexitcode=109, steps=108, gas_used=0\nVM Log (truncated):\n...INT 0\nexecute THROWIFNOT 
		105\nexecute MYADDR\nexecute XCHG s1,s4\nexecute SDEQ\nexecute THROWIF 112\nexecute OVER\nexecute 
		EQINT 0\nexecute THROWIF 106\nexecute GETGLOB
		3\nexecute NEQ\nexecute THROWIF 109\ndefault exception handler, terminating vm with exit code 109\n`,

		`unable to send external message: error code: 0 message: cannot apply external message to current state : 
		External message was not accepted\nCannot run message on account: 
		inbound external message rejected by transaction 
		6CCBB83C7D9BFBFDB40541F35AD069714856F18B4850C1273A117DF6BFADE1C6:\nexitcode=109, steps=108, 
		gas_used=0\nVM Log (truncated):\n...INT 0....`,
	} {
		require.True(t, exitCodeErrorRegex.MatchString(tt))

		exitCode, ok := extractExitCode(tt)
		require.True(t, ok)
		require.Equal(t, uint32(109), exitCode)
	}
}

type testSuite struct {
	ctx context.Context
	t   *testing.T

	chain       chains.Chain
	chainParams *observertypes.ChainParams

	liteClient *mocks.TONLiteClient

	musecore *mocks.MusecoreClient
	tss      *mocks.TSS

	gw         *toncontracts.Gateway
	baseSigner *base.Signer

	trackerBag []testTracker
}

type testTracker struct {
	nonce uint64
	hash  string
}

func newTestSuite(t *testing.T) *testSuite {
	var (
		ctx = context.Background()

		chain       = chains.TONTestnet
		chainParams = sample.ChainParams(chain.ChainId)

		liteClient = mocks.NewTONLiteClient(t)

		tss      = mocks.NewTSS(t)
		musecore = mocks.NewMusecoreClient(t).WithKeys(&keys.Keys{})

		testLogger = zerolog.New(zerolog.NewTestWriter(t))
		logger     = base.Logger{Std: testLogger, Compliance: testLogger}

		gwAccountID = ton.MustParseAccountID(testutils.GatewayAddresses[chain.ChainId])
	)

	ts := &testSuite{
		ctx: ctx,
		t:   t,

		chain:       chain,
		chainParams: chainParams,

		liteClient: liteClient,

		musecore: musecore,
		tss:      tss,

		gw:         toncontracts.NewGateway(gwAccountID),
		baseSigner: base.NewSigner(chain, tss, logger),
	}

	// Setup mocks
	ts.musecore.On("Chain").Return(chain).Maybe()

	setupTrackersBag(ts)

	return ts
}

func (ts *testSuite) OnGetAccountState(acc ton.AccountID, state tlb.ShardAccount) *mock.Call {
	return ts.liteClient.On("GetAccountState", mock.Anything, acc).Return(state, nil)
}

func (ts *testSuite) OnSendMessage(id uint32, err error) *mock.Call {
	return ts.liteClient.On("SendMessage", mock.Anything, mock.Anything).Return(id, err)
}

func (ts *testSuite) OnGetTransactionsSince(
	acc ton.AccountID,
	lt uint64,
	hash ton.Bits256,
	txs []ton.Transaction,
	err error,
) *mock.Call {
	return ts.liteClient.
		On("GetTransactionsSince", mock.Anything, acc, lt, hash).
		Return(txs, err)
}

func (ts *testSuite) Sign(msg Signable) {
	hash, err := msg.Hash()
	require.NoError(ts.t, err)

	sig, err := ts.tss.Sign(ts.ctx, hash[:], 0, 0, 0)
	require.NoError(ts.t, err)

	msg.SetSignature(sig)
}

// parses string to TON
func tonCoins(t *testing.T, raw string) math.Uint {
	t.Helper()

	const oneTON = 1_000_000_000

	f, err := strconv.ParseFloat(raw, 64)
	require.NoError(t, err)

	f *= oneTON

	return math.NewUint(uint64(f))
}

func decodeHash(t *testing.T, raw string) tlb.Bits256 {
	t.Helper()

	h, err := hex.DecodeString(raw)
	require.NoError(t, err)

	var hash tlb.Bits256

	copy(hash[:], h)

	return hash
}

func setupTrackersBag(ts *testSuite) {
	catcher := func(args mock.Arguments) {
		require.Equal(ts.t, ts.chain.ChainId, args.Get(1).(int64))
		nonce := args.Get(2).(uint64)
		txHash := args.Get(3).(string)

		ts.t.Logf("Adding outbound tracker: nonce=%d, hash=%s", nonce, txHash)

		ts.trackerBag = append(ts.trackerBag, testTracker{nonce, txHash})
	}

	ts.musecore.On(
		"PostOutboundTracker",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Maybe().Run(catcher).Return("", nil)
}
