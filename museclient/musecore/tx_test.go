package musecore

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	"github.com/RWAs-labs/go-tss/blame"
	"github.com/RWAs-labs/muse/museclient/keys"
	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSigner   = "jack"
	sampleHash   = "FA51DB4412144F1130669F2BAE8CB44AADBD8D85958DBFFCB0FE236878097E1A"
	ethBlockHash = "1a17bcc359e84ba8ae03b17ec425f97022cd11c3e279f6bdf7a96fcffa12b366"
)

func Test_GasPriceMultiplier(t *testing.T) {
	tt := []struct {
		name       string
		chain      chains.Chain
		multiplier float64
	}{
		{
			name:       "get Ethereum multiplier",
			chain:      chains.Ethereum,
			multiplier: 1.2,
		},
		{
			name:       "get Goerli multiplier",
			chain:      chains.Goerli,
			multiplier: 1.2,
		},
		{
			name:       "get BSC multiplier",
			chain:      chains.BscMainnet,
			multiplier: 1.2,
		},
		{
			name:       "get BSC Testnet multiplier",
			chain:      chains.BscTestnet,
			multiplier: 1.2,
		},
		{
			name:       "get Polygon multiplier",
			chain:      chains.Polygon,
			multiplier: 1.2,
		},
		{
			name:       "get Mumbai Testnet multiplier",
			chain:      chains.Mumbai,
			multiplier: 1.2,
		},
		{
			name:       "get Bitcoin multiplier",
			chain:      chains.BitcoinMainnet,
			multiplier: 2.0,
		},
		{
			name:       "get Bitcoin Testnet multiplier",
			chain:      chains.BitcoinTestnet,
			multiplier: 2.0,
		},
		{
			name:       "get Solana multiplier",
			chain:      chains.SolanaMainnet,
			multiplier: 1.0,
		},
		{
			name:       "get Solana devnet multiplier",
			chain:      chains.SolanaDevnet,
			multiplier: 1.0,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			multiplier := GasPriceMultiplier(tc.chain)
			require.Equal(t, tc.multiplier, multiplier)
		})
	}
}

func TestMusecore_PostGasPrice(t *testing.T) {
	ctx := context.Background()

	extraGRPC := withDummyServer(100)
	setupMockServer(t, observertypes.RegisterQueryServer, skipMethod, nil, nil, extraGRPC...)

	client := setupMusecoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
	)

	t.Run("post gas price success", func(t *testing.T) {
		hash, err := client.PostVoteGasPrice(ctx, chains.BscMainnet, 1000000, 0, 1234)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})

	// Test for failed broadcast, it will take several seconds to complete. Excluding to reduce runtime.
	//
	//t.Run("post gas price fail", func(t *testing.T) {
	//	musecoreBroadcast = MockBroadcastError
	//	hash, err := client.PostGasPrice(chains.BscMainnet, 1000000, "100", 1234)
	//	require.ErrorContains(t, err, "post gasprice failed")
	//	require.Equal(t, "", hash)
	//})
}

func TestMusecore_AddOutboundTracker(t *testing.T) {
	ctx := context.Background()

	const nonce = 123
	chainID := chains.BscMainnet.ChainId

	method := "/musechain.musecore.crosschain.Query/OutboundTracker"
	input := &crosschaintypes.QueryGetOutboundTrackerRequest{
		ChainID: chains.BscMainnet.ChainId,
		Nonce:   nonce,
	}
	output := &crosschaintypes.QueryGetOutboundTrackerResponse{
		OutboundTracker: crosschaintypes.OutboundTracker{
			Index:    "456",
			ChainId:  chainID,
			Nonce:    nonce,
			HashList: nil,
		},
	}

	extraGRPC := withDummyServer(100)
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, output, extraGRPC...)

	tendermintMock := mocks.NewSDKClientWithErr(t, nil, 0)

	client := setupMusecoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withCometBFT(tendermintMock),
	)

	t.Run("add tx hash success", func(t *testing.T) {
		tendermintMock.SetBroadcastTxHash(sampleHash)
		hash, err := client.PostOutboundTracker(ctx, chainID, nonce, "")
		assert.NoError(t, err)
		assert.Equal(t, sampleHash, hash)
	})

	t.Run("add tx hash fail", func(t *testing.T) {
		tendermintMock.SetError(errors.New("broadcast error"))
		hash, err := client.PostOutboundTracker(ctx, chainID, nonce, "")
		assert.Error(t, err)
		assert.Empty(t, hash)
	})
}

func TestMusecore_SetTSS(t *testing.T) {
	ctx := context.Background()

	extraGRPC := withDummyServer(100)
	setupMockServer(t, crosschaintypes.RegisterMsgServer, skipMethod, nil, nil, extraGRPC...)

	client := setupMusecoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
	)

	t.Run("set tss success", func(t *testing.T) {
		hash, err := client.PostVoteTSS(
			ctx,
			"musepub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			9987,
			chains.ReceiveStatus_success,
		)
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestMusecore_PostBlameData(t *testing.T) {
	ctx := context.Background()

	extraGRPC := withDummyServer(100)
	setupMockServer(t, observertypes.RegisterQueryServer, skipMethod, nil, nil, extraGRPC...)

	client := setupMusecoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
	)

	t.Run("post blame data success", func(t *testing.T) {
		hash, err := client.PostVoteBlameData(
			ctx,
			&blame.Blame{
				FailReason: "",
				IsUnicast:  false,
				BlameNodes: nil,
			},
			chains.BscMainnet.ChainId,
			"102394876-bsc",
		)
		assert.NoError(t, err)
		assert.Equal(t, sampleHash, hash)
	})
}

func TestMusecore_PostVoteInbound(t *testing.T) {
	ctx := context.Background()

	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())

	expectedOutput := observertypes.QueryHasVotedResponse{HasVoted: false}
	input := observertypes.QueryHasVotedRequest{
		BallotIdentifier: "0xd204175fc8500bcea563049cce918fa55134bd2d415d3fe137144f55e572b5ff",
		VoterAddress:     address.String(),
	}
	method := "/musechain.musecore.observer.Query/HasVoted"

	extraGRPC := withDummyServer(100)
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput, extraGRPC...)

	client := setupMusecoreClient(t,
		withDefaultObserverKeys(),
		withAccountRetriever(t, 100, 100),
		withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
	)

	t.Run("post inbound vote already voted", func(t *testing.T) {
		hash, _, err := client.PostVoteInbound(ctx, 100, 200, &crosschaintypes.MsgVoteInbound{
			Creator: address.String(),
		})
		require.NoError(t, err)
		require.Equal(t, sampleHash, hash)
	})
}

func TestMusecore_GetInboundVoteMessage(t *testing.T) {
	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())
	t.Run("get inbound vote message", func(t *testing.T) {
		msg := GetInboundVoteMessage(
			address.String(),
			chains.Ethereum.ChainId,
			"",
			address.String(),
			chains.MuseChainMainnet.ChainId,
			math.NewUint(500),
			"",
			"", 12345,
			1000,
			coin.CoinType_Gas,
			"amuse",
			address.String(),
			0,
			types.InboundStatus_SUCCESS,
		)
		require.Equal(t, address.String(), msg.Creator)
	})
}

func TestMusecore_MonitorVoteInboundResult(t *testing.T) {
	ctx := context.Background()

	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())
	client := setupMusecoreClient(t,
		withObserverKeys(keys.NewKeysWithKeybase(mocks.NewKeyring(), address, testSigner, "")),
		withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0)),
	)

	t.Run("monitor inbound vote", func(t *testing.T) {
		err := client.MonitorVoteInboundResult(ctx, sampleHash, 1000, &crosschaintypes.MsgVoteInbound{
			Creator: address.String(),
		})

		require.NoError(t, err)
	})
}

func TestMusecore_PostVoteOutbound(t *testing.T) {
	const (
		blockHeight = 1234
		accountNum  = 10
		accountSeq  = 10
	)

	ctx := context.Background()

	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())

	expectedOutput := observertypes.QueryHasVotedResponse{HasVoted: false}
	input := observertypes.QueryHasVotedRequest{
		BallotIdentifier: "0xf52f379287561dd07869de72b09fb56b7f6dfdda65b01c25882722e315f333f1",
		VoterAddress:     address.String(),
	}
	method := "/musechain.musecore.observer.Query/HasVoted"

	extraGRPC := withDummyServer(blockHeight)

	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput, extraGRPC...)
	require.NotNil(t, server)

	client := setupMusecoreClient(t,
		withDefaultObserverKeys(),
		withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0).SetBroadcastTxHash(sampleHash)),
		withAccountRetriever(t, accountNum, accountSeq),
	)

	msg := crosschaintypes.NewMsgVoteOutbound(
		address.String(),
		sampleHash,
		sampleHash,
		blockHeight,
		1000,
		math.NewInt(100),
		1200,
		math.NewUint(500),
		chains.ReceiveStatus_success,
		chains.Ethereum.ChainId,
		10001,
		coin.CoinType_Gas,
		crosschaintypes.ConfirmationMode_SAFE,
	)

	hash, ballot, err := client.PostVoteOutbound(ctx, 100_000, 200_000, msg)

	assert.NoError(t, err)
	assert.Equal(t, sampleHash, hash)
	assert.Equal(t, "0xf52f379287561dd07869de72b09fb56b7f6dfdda65b01c25882722e315f333f1", ballot)
}

func TestMusecore_MonitorVoteOutboundResult(t *testing.T) {
	ctx := context.Background()

	address := sdktypes.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())
	client := setupMusecoreClient(t,
		withObserverKeys(keys.NewKeysWithKeybase(mocks.NewKeyring(), address, testSigner, "")),
		withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0)),
	)

	t.Run("monitor outbound vote", func(t *testing.T) {
		msg := &crosschaintypes.MsgVoteOutbound{Creator: address.String()}

		err := client.MonitorVoteOutboundResult(ctx, sampleHash, 1000, msg)
		assert.NoError(t, err)
	})
}
