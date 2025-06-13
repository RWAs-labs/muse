package musecore

import (
	"context"
	"net"
	"testing"

	sdkmath "cosmossdk.io/math"

	feemarkettypes "github.com/RWAs-labs/ethermint/x/feemarket/types"
	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	keyinterfaces "github.com/RWAs-labs/muse/museclient/keys/interfaces"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/museclient/keys"
	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	cometbftrpc "github.com/cometbft/cometbft/rpc/client"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
)

const skipMethod = "skip"

// setupMockServer setup mock musecore GRPC server
func setupMockServer(
	t *testing.T,
	serviceFunc any, method string, input any, expectedOutput any,
	extra ...grpcmock.ServerOption,
) *grpcmock.Server {
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
	require.NoError(t, err)

	opts := []grpcmock.ServerOption{
		grpcmock.RegisterService(serviceFunc),
		grpcmock.WithPlanner(planner.FirstMatch()),
		grpcmock.WithListener(listener),
	}

	opts = append(opts, extra...)

	if method != skipMethod {
		opts = append(opts, func(s *grpcmock.Server) {
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(input).
				Return(expectedOutput)
		})
	}

	server := grpcmock.MockUnstartedServer(opts...)(t)

	server.Serve()

	t.Cleanup(func() {
		require.NoError(t, server.Close())
	})

	return server
}

func withDummyServer(museBlockHeight int64) []grpcmock.ServerOption {
	return []grpcmock.ServerOption{
		grpcmock.RegisterService(crosschaintypes.RegisterQueryServer),
		grpcmock.RegisterService(crosschaintypes.RegisterMsgServer),
		grpcmock.RegisterService(feemarkettypes.RegisterQueryServer),
		grpcmock.RegisterService(authtypes.RegisterQueryServer),
		func(s *grpcmock.Server) {
			// Block Height
			s.ExpectUnary("/musechain.musecore.crosschain.Query/LastMuseHeight").
				UnlimitedTimes().
				Return(crosschaintypes.QueryLastMuseHeightResponse{Height: museBlockHeight})

			// London Base Fee
			s.ExpectUnary("/ethermint.feemarket.v1.Query/Params").
				UnlimitedTimes().
				Return(feemarkettypes.QueryParamsResponse{
					Params: feemarkettypes.Params{BaseFee: sdkmath.NewInt(100)},
				})
		},
	}
}

type clientTestConfig struct {
	keys keyinterfaces.ObserverKeys
	opts []Opt
}

type clientTestOpt func(*clientTestConfig)

func withObserverKeys(keys keyinterfaces.ObserverKeys) clientTestOpt {
	return func(cfg *clientTestConfig) { cfg.keys = keys }
}

func withDefaultObserverKeys() clientTestOpt {
	var (
		key     = mocks.TestKeyringPair
		address = types.AccAddress(key.PubKey().Address().Bytes())
		keyRing = mocks.NewKeyring()
	)

	return withObserverKeys(keys.NewKeysWithKeybase(keyRing, address, testSigner, ""))
}

func withCometBFT(client cometbftrpc.Client) clientTestOpt {
	return func(cfg *clientTestConfig) { cfg.opts = append(cfg.opts, WithCometBFTClient(client)) }
}

func withAccountRetriever(t *testing.T, accNum uint64, accSeq uint64) clientTestOpt {
	ctrl := gomock.NewController(t)
	ac := mock.NewMockAccountRetriever(ctrl)
	ac.EXPECT().
		GetAccountNumberSequence(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(accNum, accSeq, nil)

	return func(cfg *clientTestConfig) {
		cfg.opts = append(cfg.opts, WithCustomAccountRetriever(ac))
	}
}

func setupMusecoreClient(t *testing.T, opts ...clientTestOpt) *Client {
	const (
		chainIP = "127.0.0.1"
		signer  = testSigner
		chainID = "musechain_7000-1"
	)

	var cfg clientTestConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.keys == nil {
		cfg.keys = &keys.Keys{}
	}

	c, err := NewClient(
		cfg.keys,
		chainIP, signer,
		chainID,
		zerolog.Nop(),
		cfg.opts...,
	)

	require.NoError(t, err)

	return c
}

// Need to test after refactor
func TestMusecore_GetGenesisSupply(t *testing.T) {
}

func TestMusecore_GetMuseHotKeyBalance(t *testing.T) {
	ctx := context.Background()

	expectedOutput := banktypes.QueryBalanceResponse{
		Balance: &types.Coin{
			Denom:  config.BaseDenom,
			Amount: sdkmath.NewInt(55646484),
		},
	}
	input := banktypes.QueryBalanceRequest{
		Address: types.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes()).String(),
		Denom:   config.BaseDenom,
	}
	method := "/cosmos.bank.v1beta1.Query/Balance"
	setupMockServer(t, banktypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClient(
		t,
		withDefaultObserverKeys(),
		withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0)),
	)

	// should be able to get balance of signer
	client.keys = keys.NewKeysWithKeybase(mocks.NewKeyring(), types.AccAddress{}, "bob", "")
	resp, err := client.GetMuseHotKeyBalance(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Balance.Amount, resp)

	// should return error on empty signer
	client.keys = keys.NewKeysWithKeybase(mocks.NewKeyring(), types.AccAddress{}, "", "")
	resp, err = client.GetMuseHotKeyBalance(ctx)
	require.Error(t, err)
	require.Equal(t, sdkmath.ZeroInt(), resp)
}

func TestMusecore_GetAllOutboundTrackerByChain(t *testing.T) {
	ctx := context.Background()

	chain := chains.BscMainnet
	expectedOutput := crosschaintypes.QueryAllOutboundTrackerByChainResponse{
		OutboundTracker: []crosschaintypes.OutboundTracker{
			{
				Index:    "tracker23456",
				ChainId:  chain.ChainId,
				Nonce:    123456,
				HashList: nil,
			},
		},
	}
	input := crosschaintypes.QueryAllOutboundTrackerByChainRequest{
		Chain: chain.ChainId,
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      2000,
			CountTotal: false,
			Reverse:    false,
		},
	}
	method := "/musechain.musecore.crosschain.Query/OutboundTrackerAllByChain"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClient(
		t,
		withDefaultObserverKeys(),
		withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0)),
	)

	resp, err := client.GetAllOutboundTrackerByChain(ctx, chain.ChainId, interfaces.Ascending)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.OutboundTracker, resp)

	resp, err = client.GetAllOutboundTrackerByChain(ctx, chain.ChainId, interfaces.Descending)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.OutboundTracker, resp)
}

func TestMusecore_SubscribeNewBlocks(t *testing.T) {
	ctx := context.Background()
	cometBFTClient := mocks.NewSDKClientWithErr(t, nil, 0)
	client := setupMusecoreClient(
		t,
		withDefaultObserverKeys(),
		withCometBFT(cometBFTClient),
	)

	newBlockChan, err := client.NewBlockSubscriber(ctx)
	require.NoError(t, err)

	height := int64(10)

	cometBFTClient.PublishToSubscribers(coretypes.ResultEvent{
		Data: cometbfttypes.EventDataNewBlock{
			Block: &cometbfttypes.Block{
				Header: cometbfttypes.Header{
					Height: height,
				},
			},
		},
	})

	newBlockEvent := <-newBlockChan
	require.Equal(t, height, newBlockEvent.Block.Header.Height)
}
