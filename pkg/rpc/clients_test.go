package rpc

import (
	"context"
	"net"
	"testing"

	sdkmath "cosmossdk.io/math"

	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	feemarkettypes "github.com/RWAs-labs/ethermint/x/feemarket/types"
	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/testutil/sample"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	lightclienttypes "github.com/RWAs-labs/muse/x/lightclient/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

const skipMethod = "skip"
const gRPCListenPath = "127.0.0.1:47392"

// setupMockServer setup mock musecore GRPC server
func setupMockServer(
	t *testing.T,
	serviceFunc any, method string, input any, expectedOutput any,
	extra ...grpcmock.ServerOption,
) *grpcmock.Server {
	listener, err := net.Listen("tcp", gRPCListenPath)
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

func setupMusecoreClients(t *testing.T) Clients {
	c, err := NewGRPCClients(gRPCListenPath, grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)

	return c
}

func TestMusecore_GetBallot(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryBallotByIdentifierResponse{
		BallotIdentifier: "123",
		Voters:           nil,
		ObservationType:  0,
		BallotStatus:     0,
	}
	input := observertypes.QueryBallotByIdentifierRequest{BallotIdentifier: "123"}
	method := "/musechain.musecore.observer.Query/BallotByIdentifier"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetBallotByID(ctx, "123")
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestMusecore_GetCrosschainFlags(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetCrosschainFlagsResponse{CrosschainFlags: observertypes.CrosschainFlags{
		IsInboundEnabled:      true,
		IsOutboundEnabled:     false,
		GasPriceIncreaseFlags: nil,
	}}
	input := observertypes.QueryGetCrosschainFlagsRequest{}
	method := "/musechain.musecore.observer.Query/CrosschainFlags"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetCrosschainFlags(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrosschainFlags, resp)
}

func TestMusecore_GetRateLimiterFlags(t *testing.T) {
	ctx := context.Background()

	// create sample flags
	rateLimiterFlags := sample.RateLimiterFlags()
	expectedOutput := crosschaintypes.QueryRateLimiterFlagsResponse{
		RateLimiterFlags: rateLimiterFlags,
	}

	// setup mock server
	input := crosschaintypes.QueryRateLimiterFlagsRequest{}
	method := "/musechain.musecore.crosschain.Query/RateLimiterFlags"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	// query
	resp, err := client.GetRateLimiterFlags(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.RateLimiterFlags, resp)
}

func TestMusecore_HeaderEnabledChains(t *testing.T) {
	ctx := context.Background()

	expectedOutput := lightclienttypes.QueryHeaderEnabledChainsResponse{
		HeaderEnabledChains: []lightclienttypes.HeaderSupportedChain{
			{
				ChainId: chains.Ethereum.ChainId,
				Enabled: true,
			},
			{
				ChainId: chains.BitcoinMainnet.ChainId,
				Enabled: true,
			},
		},
	}
	input := lightclienttypes.QueryHeaderEnabledChainsRequest{}
	method := "/musechain.musecore.lightclient.Query/HeaderEnabledChains"
	setupMockServer(t, lightclienttypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetBlockHeaderEnabledChains(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.HeaderEnabledChains, resp)
}

func TestMusecore_GetChainParamsForChainID(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetChainParamsForChainResponse{ChainParams: &observertypes.ChainParams{
		ChainId:               123,
		BallotThreshold:       sdkmath.LegacyZeroDec(),
		MinObserverDelegation: sdkmath.LegacyZeroDec(),
	}}
	input := observertypes.QueryGetChainParamsForChainRequest{ChainId: 123}
	method := "/musechain.musecore.observer.Query/GetChainParamsForChain"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetChainParamsForChainID(ctx, 123)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainParams, resp)
}

func TestMusecore_GetChainParams(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetChainParamsResponse{ChainParams: &observertypes.ChainParamsList{
		ChainParams: []*observertypes.ChainParams{
			{
				ChainId:               123,
				MinObserverDelegation: sdkmath.LegacyZeroDec(),
				BallotThreshold:       sdkmath.LegacyZeroDec(),
			},
		},
	}}
	input := observertypes.QueryGetChainParamsRequest{}
	method := "/musechain.musecore.observer.Query/GetChainParams"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetChainParams(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainParams.ChainParams, resp)
}

func TestMusecore_GetUpgradePlan(t *testing.T) {
	ctx := context.Background()

	expectedOutput := upgradetypes.QueryCurrentPlanResponse{
		Plan: &upgradetypes.Plan{
			Name:   "big upgrade",
			Height: 100,
		},
	}
	input := upgradetypes.QueryCurrentPlanRequest{}
	method := "/cosmos.upgrade.v1beta1.Query/CurrentPlan"
	setupMockServer(t, upgradetypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetUpgradePlan(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Plan, resp)
}

func TestMusecore_GetAllCctx(t *testing.T) {
	ctx := context.Background()

	expectedOutput := crosschaintypes.QueryAllCctxResponse{
		CrossChainTx: []*crosschaintypes.CrossChainTx{
			{
				Index: "cross-chain4456",
			},
		},
		Pagination: nil,
	}
	input := crosschaintypes.QueryAllCctxRequest{}
	method := "/musechain.musecore.crosschain.Query/CctxAll"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetAllCctx(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestMusecore_GetCctxByHash(t *testing.T) {
	ctx := context.Background()

	expectedOutput := crosschaintypes.QueryGetCctxResponse{CrossChainTx: &crosschaintypes.CrossChainTx{
		Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3",
	}}
	input := crosschaintypes.QueryGetCctxRequest{
		Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3",
	}
	method := "/musechain.musecore.crosschain.Query/Cctx"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetCctxByHash(ctx, "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3")
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestMusecore_GetCctxByNonce(t *testing.T) {
	ctx := context.Background()

	expectedOutput := crosschaintypes.QueryGetCctxResponse{CrossChainTx: &crosschaintypes.CrossChainTx{
		Index: "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3",
	}}
	input := crosschaintypes.QueryGetCctxByNonceRequest{
		ChainID: 7000,
		Nonce:   55,
	}
	method := "/musechain.musecore.crosschain.Query/CctxByNonce"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetCctxByNonce(ctx, 7000, 55)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
}

func TestMusecore_GetObserverList(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryObserverSetResponse{
		Observers: []string{
			"muse19jr7nl82lrktge35f52x9g5y5prmvchmk40zhg",
			"muse1cxj07f3ju484ry2cnnhxl5tryyex7gev0yzxtj",
			"muse1hjct6q7npsspsg3dgvzk3sdf89spmlpf7rqmnw",
		},
	}
	input := observertypes.QueryObserverSet{}
	method := "/musechain.musecore.observer.Query/ObserverSet"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetObserverList(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Observers, resp)
}

func TestMusecore_GetRateLimiterInput(t *testing.T) {
	ctx := context.Background()

	expectedOutput := &crosschaintypes.QueryRateLimiterInputResponse{
		Height:                  10,
		CctxsMissed:             []*crosschaintypes.CrossChainTx{sample.CrossChainTx(t, "1-1")},
		CctxsPending:            []*crosschaintypes.CrossChainTx{sample.CrossChainTx(t, "1-2")},
		TotalPending:            1,
		PastCctxsValue:          "123456",
		PendingCctxsValue:       "1234",
		LowestPendingCctxHeight: 2,
	}
	input := crosschaintypes.QueryRateLimiterInputRequest{Window: 10}
	method := "/musechain.musecore.crosschain.Query/RateLimiterInput"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetRateLimiterInput(ctx, 10)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, resp)
}

func TestMusecore_ListPendingCctx(t *testing.T) {
	ctx := context.Background()

	expectedOutput := crosschaintypes.QueryListPendingCctxResponse{
		CrossChainTx: []*crosschaintypes.CrossChainTx{
			{
				Index: "cross-chain4456",
			},
		},
		TotalPending: 1,
	}
	input := crosschaintypes.QueryListPendingCctxRequest{ChainId: 7000}
	method := "/musechain.musecore.crosschain.Query/ListPendingCctx"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, totalPending, err := client.ListPendingCCTX(ctx, 7000)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.CrossChainTx, resp)
	require.Equal(t, expectedOutput.TotalPending, totalPending)
}

func TestMusecore_GetAbortedMuseAmount(t *testing.T) {
	ctx := context.Background()

	expectedOutput := crosschaintypes.QueryMuseAccountingResponse{AbortedMuseAmount: "1080999"}
	input := crosschaintypes.QueryMuseAccountingRequest{}
	method := "/musechain.musecore.crosschain.Query/MuseAccounting"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetAbortedMuseAmount(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.AbortedMuseAmount, resp)
}

// Need to test after refactor
func TestMusecore_GetGenesisSupply(t *testing.T) {
}

func TestMusecore_GetMuseTokenSupplyOnNode(t *testing.T) {
	ctx := context.Background()

	expectedOutput := banktypes.QuerySupplyOfResponse{
		Amount: types.Coin{
			Denom:  config.BaseDenom,
			Amount: sdkmath.NewInt(329438),
		}}
	input := banktypes.QuerySupplyOfRequest{Denom: config.BaseDenom}
	method := "/cosmos.bank.v1beta1.Query/SupplyOf"
	setupMockServer(t, banktypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetMuseTokenSupplyOnNode(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.GetAmount().Amount, resp)
}

func TestMusecore_GetBlockHeight(t *testing.T) {
	ctx := context.Background()

	method := "/musechain.musecore.crosschain.Query/LastMuseHeight"
	input := &crosschaintypes.QueryLastMuseHeightRequest{}
	output := &crosschaintypes.QueryLastMuseHeightResponse{Height: 12345}

	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, output)

	client := setupMusecoreClients(t)

	t.Run("last block height", func(t *testing.T) {
		height, err := client.GetBlockHeight(ctx)
		require.NoError(t, err)
		require.Equal(t, int64(12345), height)
	})
}

func TestMusecore_GetLatestMuseBlock(t *testing.T) {
	ctx := context.Background()

	expectedOutput := cmtservice.GetLatestBlockResponse{
		SdkBlock: &cmtservice.Block{
			Header:     cmtservice.Header{},
			Data:       tmtypes.Data{},
			Evidence:   tmtypes.EvidenceList{},
			LastCommit: nil,
		},
	}
	input := cmtservice.GetLatestBlockRequest{}
	method := "/cosmos.base.tendermint.v1beta1.Service/GetLatestBlock"
	setupMockServer(t, cmtservice.RegisterServiceServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetLatestMuseBlock(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.SdkBlock, resp)
}

func TestMusecore_GetNodeInfo(t *testing.T) {
	ctx := context.Background()

	expectedOutput := cmtservice.GetNodeInfoResponse{
		DefaultNodeInfo:    nil,
		ApplicationVersion: &cmtservice.VersionInfo{},
	}
	input := cmtservice.GetNodeInfoRequest{}
	method := "/cosmos.base.tendermint.v1beta1.Service/GetNodeInfo"
	setupMockServer(t, cmtservice.RegisterServiceServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetNodeInfo(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestMusecore_GetBaseGasPrice(t *testing.T) {
	ctx := context.Background()

	expectedOutput := feemarkettypes.QueryParamsResponse{
		Params: feemarkettypes.Params{
			BaseFee: sdkmath.NewInt(23455),
		},
	}
	input := feemarkettypes.QueryParamsRequest{}
	method := "/ethermint.feemarket.v1.Query/Params"
	setupMockServer(t, feemarkettypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetBaseGasPrice(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Params.BaseFee.Int64(), resp)
}

func TestMusecore_GetNonceByChain(t *testing.T) {
	ctx := context.Background()

	chain := chains.BscMainnet
	expectedOutput := observertypes.QueryGetChainNoncesResponse{
		ChainNonces: observertypes.ChainNonces{
			Creator:         "",
			ChainId:         chain.ChainId,
			Nonce:           8446,
			Signers:         nil,
			FinalizedHeight: 0,
		},
	}
	input := observertypes.QueryGetChainNoncesRequest{ChainId: chain.ChainId}
	method := "/musechain.musecore.observer.Query/ChainNonces"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetNonceByChain(ctx, chain)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainNonces, resp)
}

func TestMusecore_GetAllNodeAccounts(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryAllNodeAccountResponse{
		NodeAccount: []*observertypes.NodeAccount{
			{
				Operator:       "muse19jr7nl82lrktge35f52x9g5y5prmvchmk40zhg",
				GranteeAddress: "muse1kxhesgcvl6j5upupd9m3d3g3gfz4l3pcpqfnw6",
				GranteePubkey:  nil,
				NodeStatus:     0,
			},
		},
	}
	input := observertypes.QueryAllNodeAccountRequest{}
	method := "/musechain.musecore.observer.Query/NodeAccountAll"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetAllNodeAccounts(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.NodeAccount, resp)
}

func TestMusecore_GetKeyGen(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetKeygenResponse{
		Keygen: &observertypes.Keygen{
			Status:         observertypes.KeygenStatus_KeyGenSuccess,
			GranteePubkeys: nil,
			BlockNumber:    5646,
		}}
	input := observertypes.QueryGetKeygenRequest{}
	method := "/musechain.musecore.observer.Query/Keygen"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetKeyGen(ctx)
	require.NoError(t, err)
	require.Equal(t, *expectedOutput.Keygen, resp)
}

func TestMusecore_GetBallotByID(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryBallotByIdentifierResponse{
		BallotIdentifier: "ballot1235",
	}
	input := observertypes.QueryBallotByIdentifierRequest{BallotIdentifier: "ballot1235"}
	method := "/musechain.musecore.observer.Query/BallotByIdentifier"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetBallot(ctx, "ballot1235")
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestMusecore_GetInboundTrackersForChain(t *testing.T) {
	ctx := context.Background()

	chainID := chains.BscMainnet.ChainId
	expectedOutput := crosschaintypes.QueryAllInboundTrackerByChainResponse{
		InboundTracker: []crosschaintypes.InboundTracker{
			{
				ChainId:  chainID,
				TxHash:   "DC76A6DCCC3AA62E89E69042ADC44557C50D59E4D3210C37D78DC8AE49B3B27F",
				CoinType: coin.CoinType_Gas,
			},
		},
	}
	input := crosschaintypes.QueryAllInboundTrackerByChainRequest{ChainId: chainID}
	method := "/musechain.musecore.crosschain.Query/InboundTrackerAllByChain"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetInboundTrackersForChain(ctx, chainID)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.InboundTracker, resp)
}

func TestMusecore_GetTss(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetTSSResponse{
		TSS: observertypes.TSS{
			TssPubkey:           "musepub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TssParticipantList:  nil,
			OperatorAddressList: nil,
			FinalizedMuseHeight: 1000,
			KeyGenMuseHeight:    900,
		},
	}
	input := observertypes.QueryGetTSSRequest{}
	method := "/musechain.musecore.observer.Query/TSS"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetTSS(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.TSS, resp)
}

func TestMusecore_GetEthTssAddress(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetTssAddressResponse{
		Eth: "0x70e967acfcc17c3941e87562161406d41676fd83",
		Btc: "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
	}
	input := observertypes.QueryGetTssAddressRequest{}
	method := "/musechain.musecore.observer.Query/GetTssAddress"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetEVMTSSAddress(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Eth, resp)
}

func TestMusecore_GetBtcTssAddress(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryGetTssAddressResponse{
		Eth: "0x70e967acfcc17c3941e87562161406d41676fd83",
		Btc: "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
	}
	input := observertypes.QueryGetTssAddressRequest{BitcoinChainId: 8332}
	method := "/musechain.musecore.observer.Query/GetTssAddress"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetBTCTSSAddress(ctx, 8332)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Btc, resp)
}

func TestMusecore_GetTssHistory(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryTssHistoryResponse{
		TssList: []observertypes.TSS{
			{
				TssPubkey:           "musepub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
				TssParticipantList:  nil,
				OperatorAddressList: nil,
				FinalizedMuseHeight: 46546,
				KeyGenMuseHeight:    6897,
			},
		},
	}
	input := observertypes.QueryTssHistoryRequest{}
	method := "/musechain.musecore.observer.Query/TssHistory"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetTSSHistory(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.TssList, resp)
}

func TestMusecore_GetOutboundTracker(t *testing.T) {
	chain := chains.BscMainnet
	expectedOutput := crosschaintypes.QueryGetOutboundTrackerResponse{
		OutboundTracker: crosschaintypes.OutboundTracker{
			Index:    "tracker12345",
			ChainId:  chain.ChainId,
			Nonce:    456,
			HashList: nil,
		},
	}
	input := crosschaintypes.QueryGetOutboundTrackerRequest{
		ChainID: chain.ChainId,
		Nonce:   456,
	}
	method := "/musechain.musecore.crosschain.Query/OutboundTracker"
	setupMockServer(t, crosschaintypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	ctx := context.Background()
	resp, err := client.GetOutboundTracker(ctx, chain, 456)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.OutboundTracker, *resp)
}

func TestMusecore_GetPendingNoncesByChain(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryPendingNoncesByChainResponse{
		PendingNonces: observertypes.PendingNonces{
			NonceLow:  0,
			NonceHigh: 0,
			ChainId:   chains.Ethereum.ChainId,
			Tss:       "",
		},
	}
	input := observertypes.QueryPendingNoncesByChainRequest{ChainId: chains.Ethereum.ChainId}
	method := "/musechain.musecore.observer.Query/PendingNoncesByChain"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetPendingNoncesByChain(ctx, chains.Ethereum.ChainId)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.PendingNonces, resp)
}

func TestMusecore_GetBlockHeaderChainState(t *testing.T) {
	ctx := context.Background()

	chainID := chains.BscMainnet.ChainId
	expectedOutput := lightclienttypes.QueryGetChainStateResponse{ChainState: &lightclienttypes.ChainState{
		ChainId:         chainID,
		LatestHeight:    5566654,
		EarliestHeight:  4454445,
		LatestBlockHash: nil,
	}}
	input := lightclienttypes.QueryGetChainStateRequest{ChainId: chainID}
	method := "/musechain.musecore.lightclient.Query/ChainState"
	setupMockServer(t, lightclienttypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetBlockHeaderChainState(ctx, chainID)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainState, resp)
}

func TestMusecore_GetSupportedChains(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QuerySupportedChainsResponse{
		Chains: []chains.Chain{
			{
				ChainId:     chains.BitcoinMainnet.ChainId,
				Network:     chains.BscMainnet.Network,
				NetworkType: chains.BscMainnet.NetworkType,
				Vm:          chains.BscMainnet.Vm,
				Consensus:   chains.BscMainnet.Consensus,
				IsExternal:  chains.BscMainnet.IsExternal,
				Name:        chains.BscMainnet.Name,
			},
			{
				ChainId:     chains.Ethereum.ChainId,
				Network:     chains.Ethereum.Network,
				NetworkType: chains.Ethereum.NetworkType,
				Vm:          chains.Ethereum.Vm,
				Consensus:   chains.Ethereum.Consensus,
				IsExternal:  chains.Ethereum.IsExternal,
				Name:        chains.Ethereum.Name,
			},
		},
	}
	input := observertypes.QuerySupportedChains{}
	method := "/musechain.musecore.observer.Query/SupportedChains"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetSupportedChains(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Chains, resp)
}

func TestMusecore_GetAdditionalChains(t *testing.T) {
	ctx := context.Background()

	expectedOutput := authoritytypes.QueryGetChainInfoResponse{
		ChainInfo: authoritytypes.ChainInfo{
			Chains: []chains.Chain{
				chains.BitcoinMainnet,
				chains.Ethereum,
			},
		},
	}
	input := observertypes.QuerySupportedChains{}
	method := "/musechain.musecore.authority.Query/ChainInfo"

	setupMockServer(t, authoritytypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetAdditionalChains(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.ChainInfo.Chains, resp)
}

func TestMusecore_GetPendingNonces(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryAllPendingNoncesResponse{
		PendingNonces: []observertypes.PendingNonces{
			{
				NonceLow:  225,
				NonceHigh: 226,
				ChainId:   8332,
				Tss:       "bc1qm24wp577nk8aacckv8np465z3dvmu7ry45el6y",
			},
		},
	}
	input := observertypes.QueryAllPendingNoncesRequest{}
	method := "/musechain.musecore.observer.Query/PendingNoncesAll"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.GetPendingNonces(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}

func TestMusecore_Prove(t *testing.T) {
	ctx := context.Background()

	chainId := chains.BscMainnet.ChainId
	txHash := "9c8d02b6956b9c78ecb6090a8160faaa48e7aecfd0026fcdf533721d861436a3"
	blockHash := "0000000000000000000172c9a64f86f208b867a84dc7a0b7c75be51e750ed8eb"
	txIndex := 555
	expectedOutput := lightclienttypes.QueryProveResponse{
		Valid: true,
	}
	input := lightclienttypes.QueryProveRequest{
		ChainId:   chainId,
		TxHash:    txHash,
		Proof:     nil,
		BlockHash: blockHash,
		TxIndex:   int64(txIndex),
	}
	method := "/musechain.musecore.lightclient.Query/Prove"
	setupMockServer(t, lightclienttypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.Prove(ctx, blockHash, txHash, int64(txIndex), nil, chainId)
	require.NoError(t, err)
	require.Equal(t, expectedOutput.Valid, resp)
}

func TestMusecore_HasVoted(t *testing.T) {
	ctx := context.Background()

	expectedOutput := observertypes.QueryHasVotedResponse{HasVoted: true}
	input := observertypes.QueryHasVotedRequest{
		BallotIdentifier: "123456asdf",
		VoterAddress:     "muse1l40mm7meacx03r4lp87s9gkxfan32xnznp42u6",
	}
	method := "/musechain.musecore.observer.Query/HasVoted"
	setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)

	client := setupMusecoreClients(t)

	resp, err := client.HasVoted(ctx, "123456asdf", "muse1l40mm7meacx03r4lp87s9gkxfan32xnznp42u6")
	require.NoError(t, err)
	require.Equal(t, expectedOutput.HasVoted, resp)
}
