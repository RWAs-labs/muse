package musecore

import (
	"context"
	"errors"
	"net"
	"testing"

	sdkmath "cosmossdk.io/math"
	feemarkettypes "github.com/RWAs-labs/ethermint/x/feemarket/types"
	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"

	"github.com/RWAs-labs/muse/museclient/keys"
	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	"github.com/RWAs-labs/muse/pkg/chains"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestHandleBroadcastError(t *testing.T) {
	type response struct {
		retry  bool
		report bool
	}
	testCases := map[error]response{
		errors.New("nonce too low"):                       {retry: false, report: true},
		errors.New("replacement transaction underpriced"): {retry: false, report: false},
		errors.New("already known"):                       {retry: false, report: true},
		errors.New(""):                                    {retry: true, report: false},
	}
	for input, output := range testCases {
		retry, report := HandleBroadcastError(input, 100, 1, "")
		require.Equal(t, output.report, report)
		require.Equal(t, output.retry, retry)
	}
}

func TestBroadcast(t *testing.T) {
	ctx := context.Background()

	address := types.AccAddress(mocks.TestKeyringPair.PubKey().Address().Bytes())

	//Setup server for multiple grpc calls
	listener, err := net.Listen("tcp", "127.0.0.1:9090")
	require.NoError(t, err)

	server := grpcmock.MockUnstartedServer(
		grpcmock.RegisterService(crosschaintypes.RegisterQueryServer),
		grpcmock.RegisterService(feemarkettypes.RegisterQueryServer),
		grpcmock.RegisterService(authtypes.RegisterQueryServer),
		grpcmock.WithPlanner(planner.FirstMatch()),
		grpcmock.WithListener(listener),
		func(s *grpcmock.Server) {
			method := "/musechain.musecore.crosschain.Query/LastMuseHeight"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(crosschaintypes.QueryLastMuseHeightRequest{}).
				Return(crosschaintypes.QueryLastMuseHeightResponse{Height: 0})

			method = "/ethermint.feemarket.v1.Query/Params"
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(feemarkettypes.QueryParamsRequest{}).
				Return(feemarkettypes.QueryParamsResponse{
					Params: feemarkettypes.Params{
						BaseFee: sdkmath.NewInt(23455),
					},
				})
		},
	)(t)

	server.Serve()
	defer server.Close()

	observerKeys := keys.NewKeysWithKeybase(mocks.NewKeyring(), address, testSigner, "")

	t.Run("broadcast success", func(t *testing.T) {
		client := setupMusecoreClient(t,
			withObserverKeys(observerKeys),
			withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0)),
		)

		msg := crosschaintypes.NewMsgVoteGasPrice(address.String(), chains.Ethereum.ChainId, 10000, 1000, 1)
		authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
		require.NoError(t, err)

		_, err = client.Broadcast(ctx, 10_000, authzMsg, authzSigner)
		require.NoError(t, err)
	})

	t.Run("broadcast failed", func(t *testing.T) {
		client := setupMusecoreClient(t,
			withObserverKeys(observerKeys),
			withCometBFT(
				mocks.NewSDKClientWithErr(t, errors.New("account sequence mismatch, expected 5 got 4"), 32),
			),
		)

		msg := crosschaintypes.NewMsgVoteGasPrice(address.String(), chains.Ethereum.ChainId, 10000, 1000, 1)
		authzMsg, authzSigner, err := WrapMessageWithAuthz(msg)
		require.NoError(t, err)

		_, err = client.Broadcast(ctx, 10_000, authzMsg, authzSigner)
		require.Error(t, err)
	})
}
