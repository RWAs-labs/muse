package signer

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/sui/client"
	"github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/museclient/keys"
	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	"github.com/RWAs-labs/muse/museclient/testutils/testlog"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/contracts/sui"
	"github.com/RWAs-labs/muse/testutil/sample"
	cc "github.com/RWAs-labs/muse/x/crosschain/types"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSigner(t *testing.T) {
	t.Run("ProcessCCTX", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		const museHeight = 1000

		// Given cctx
		nonce := uint64(123)
		amount := math.NewUint(100_000)
		receiver := "0xdecb47015beebed053c19ef48fe4d722fa3870f567133d235ebe3a70da7b0000"

		cctx := sample.CrossChainTxV2(t, "0xABC123")
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams = []*cc.OutboundParams{{
			Receiver:        receiver,
			ReceiverChainId: ts.Chain.ChainId,
			CoinType:        coin.CoinType_Gas,
			Amount:          amount,
			TssNonce:        nonce,
			GasPrice:        "1000",
			CallOptions: &cc.CallOptions{
				GasLimit: 42,
			},
		}}

		// Given mocked gateway nonce
		ts.MockGatewayNonce(nonce)

		// Given mocked WithdrawCapID
		const withdrawCapID = "0xWithdrawCapID"
		ts.MockWithdrawCapID(withdrawCapID)

		// Given expected MoveCall
		txBytes := base64.StdEncoding.EncodeToString([]byte("raw_tx_bytes"))

		ts.MockMoveCall(func(req models.MoveCallRequest) {
			require.Equal(t, ts.TSS.PubKey().AddressSui(), req.Signer)
			require.Equal(t, ts.Gateway.PackageID(), req.PackageObjectId)
			require.Equal(t, "withdraw", req.Function)

			expectedArgs := []any{
				ts.Gateway.ObjectID(),
				amount.String(),
				fmt.Sprintf("%d", nonce),
				receiver,
				"42000",
				withdrawCapID,
			}
			require.Equal(t, expectedArgs, req.Arguments)
		}, txBytes)

		// Given expected SuiExecuteTransactionBlock
		const digest = "0xTransactionBlockDigest"
		ts.MockExec(func(req models.SuiExecuteTransactionBlockRequest) {
			require.Equal(t, txBytes, req.TxBytes)
			require.NotEmpty(t, req.Signature)
		}, digest)

		// Given included tx from Sui RPC
		ts.SuiMock.
			On("SuiGetTransactionBlock", mock.Anything, mock.Anything).
			Return(models.SuiTransactionBlockResponse{
				Digest: digest,
				Effects: models.SuiEffects{
					Status: models.ExecutionStatus{
						Status: client.TxStatusSuccess,
					},
				},
				Checkpoint: "1000000",
			}, nil)

		// ACT
		err := ts.Signer.ProcessCCTX(ts.Ctx, cctx, museHeight)

		// ASSERT
		require.NoError(t, err)

		// Wait for vote posting
		wait := func() bool {
			if len(ts.TrackerBag) == 0 {
				return false
			}

			vote := ts.TrackerBag[0]
			return vote.hash == digest && vote.nonce == nonce
		}

		require.Eventually(t, wait, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("ProcessCCTX restricted address", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		const museHeight = 1000

		// Given cctx
		nonce := uint64(123)
		amount := math.NewUint(100_000)
		receiver := "0xdecb47015beebed053c19ef48fe4d722fa3870f567133d235ebe3a70da7b0000"

		cctx := sample.CrossChainTxV2(t, "0xABC123")
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams = []*cc.OutboundParams{{
			Receiver:        receiver,
			ReceiverChainId: ts.Chain.ChainId,
			CoinType:        coin.CoinType_Gas,
			Amount:          amount,
			TssNonce:        nonce,
			GasPrice:        "1000",
			CallOptions: &cc.CallOptions{
				GasLimit: 42,
			},
		}}

		// Given compliance config
		cfg := config.Config{
			ComplianceConfig: config.ComplianceConfig{
				RestrictedAddresses: []string{receiver},
			},
		}
		config.SetRestrictedAddressesFromConfig(cfg)

		// Given mocked gateway nonce
		ts.MockGatewayNonce(nonce)

		// Given mocked WithdrawCapID
		const withdrawCapID = "0xWithdrawCapID"
		ts.MockWithdrawCapID(withdrawCapID)

		// Given expected MoveCall
		txBytes := base64.StdEncoding.EncodeToString([]byte("raw_tx_bytes"))

		ts.MockMoveCall(func(req models.MoveCallRequest) {
			require.Equal(t, ts.TSS.PubKey().AddressSui(), req.Signer)
			require.Equal(t, ts.Gateway.PackageID(), req.PackageObjectId)
			require.Equal(t, "increase_nonce", req.Function)

			expectedArgs := []any{
				ts.Gateway.ObjectID(),
				fmt.Sprintf("%d", nonce),
				withdrawCapID,
			}
			require.Equal(t, expectedArgs, req.Arguments)
		}, txBytes)

		// Given expected SuiExecuteTransactionBlock
		const digest = "0xTransactionBlockDigest"
		ts.MockExec(func(req models.SuiExecuteTransactionBlockRequest) {
			require.Equal(t, txBytes, req.TxBytes)
			require.NotEmpty(t, req.Signature)
		}, digest)

		// Given included tx from Sui RPC
		ts.SuiMock.
			On("SuiGetTransactionBlock", mock.Anything, mock.Anything).
			Return(models.SuiTransactionBlockResponse{
				Digest: digest,
				Effects: models.SuiEffects{
					Status: models.ExecutionStatus{
						Status: client.TxStatusSuccess,
					},
				},
				Checkpoint: "1000000",
			}, nil)

		// ACT
		err := ts.Signer.ProcessCCTX(ts.Ctx, cctx, museHeight)

		// ASSERT
		require.NoError(t, err)

		// Wait for vote posting
		wait := func() bool {
			if len(ts.TrackerBag) == 0 {
				return false
			}

			vote := ts.TrackerBag[0]
			return vote.hash == digest && vote.nonce == nonce
		}

		require.Eventually(t, wait, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("ProcessCCTX invalid receiver address", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		const museHeight = 1000

		// Given cctx
		nonce := uint64(123)
		amount := math.NewUint(100_000)

		// Given invalid receiver address, it's a EVM address
		receiver := "0x547a07f0564e0c8d48c4ae53305eabdef87e9610"

		cctx := sample.CrossChainTxV2(t, "0xABC123")
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams = []*cc.OutboundParams{{
			Receiver:        receiver,
			ReceiverChainId: ts.Chain.ChainId,
			CoinType:        coin.CoinType_Gas,
			Amount:          amount,
			TssNonce:        nonce,
			GasPrice:        "1000",
			CallOptions: &cc.CallOptions{
				GasLimit: 42,
			},
		}}

		// Given mocked gateway nonce
		ts.MockGatewayNonce(nonce)

		// Given mocked WithdrawCapID
		const withdrawCapID = "0xWithdrawCapID"
		ts.MockWithdrawCapID(withdrawCapID)

		// Given expected MoveCall
		txBytes := base64.StdEncoding.EncodeToString([]byte("raw_tx_bytes"))

		ts.MockMoveCall(func(req models.MoveCallRequest) {
			require.Equal(t, ts.TSS.PubKey().AddressSui(), req.Signer)
			require.Equal(t, ts.Gateway.PackageID(), req.PackageObjectId)
			require.Equal(t, "increase_nonce", req.Function)

			expectedArgs := []any{
				ts.Gateway.ObjectID(),
				fmt.Sprintf("%d", nonce),
				withdrawCapID,
			}
			require.Equal(t, expectedArgs, req.Arguments)
		}, txBytes)

		// Given expected SuiExecuteTransactionBlock
		const digest = "0xTransactionBlockDigest"
		ts.MockExec(func(req models.SuiExecuteTransactionBlockRequest) {
			require.Equal(t, txBytes, req.TxBytes)
			require.NotEmpty(t, req.Signature)
		}, digest)

		// Given included tx from Sui RPC
		ts.SuiMock.
			On("SuiGetTransactionBlock", mock.Anything, mock.Anything).
			Return(models.SuiTransactionBlockResponse{
				Digest: digest,
				Effects: models.SuiEffects{
					Status: models.ExecutionStatus{
						Status: client.TxStatusSuccess,
					},
				},
				Checkpoint: "1000000",
			}, nil)

		// ACT
		err := ts.Signer.ProcessCCTX(ts.Ctx, cctx, museHeight)

		// ASSERT
		require.NoError(t, err)

		// Wait for vote posting
		wait := func() bool {
			if len(ts.TrackerBag) == 0 {
				return false
			}

			vote := ts.TrackerBag[0]
			return vote.hash == digest && vote.nonce == nonce
		}

		require.Eventually(t, wait, 5*time.Second, 100*time.Millisecond)
	})
}

func Test_ValidSuiAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{
			name:    "Valid full-length address",
			address: "0x2a4c5a97b561ac5b38edc4b4e9b2c183c57b56df5b1ea2f1c6f2e4a44b92d59f",
			wantErr: false,
		},
		{
			name:    "Uppercase addresses are explicitly rejected",
			address: "0X2A4C5A97B561AC5B38EDC4B4E9B2C183C57B56DF5B1EA2F1C6F2E4A44B92D59F",
			wantErr: true,
		},
		{
			name:    "Short addresses are explicitly rejected",
			address: "0x1a",
			wantErr: true,
		},
		{
			name:    "Missing 0x prefix",
			address: "2a4c5a97b561ac5b38edc4b4e9b2c183c57b56df5b1ea2f1c6f2e4a44b92d59f",
			wantErr: true,
		},
		{
			name:    "Too long address",
			address: "0x" + "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899a1",
			wantErr: true,
		},
		{
			name:    "Invalid hex characters",
			address: "0xZZZZZZ0000000000000000000000000000000000000000000000000000000000",
			wantErr: true,
		},
		{
			name:    "Empty string",
			address: "",
			wantErr: true,
		},
		{
			name:    "Only 0x",
			address: "0x",
			wantErr: true,
		},
		{
			name:    "Minimal valid single-byte address",
			address: "0x0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAddress(tt.address)
			if tt.wantErr {
				require.Error(t, err, tt.address)
				return
			}

			require.NoError(t, err, tt.address)
		})
	}
}

type testSuite struct {
	t   *testing.T
	Ctx context.Context

	Chain chains.Chain

	TSS      *mocks.TSS
	Musecore *mocks.MusecoreClient
	SuiMock  *mocks.SuiClient
	Gateway  *sui.Gateway

	*Signer

	TrackerBag []testTracker
}

func newTestSuite(t *testing.T) *testSuite {
	var (
		ctx = context.Background()

		chain       = chains.SuiMainnet
		chainParams = mocks.MockChainParams(chain.ChainId, 10)

		tss      = mocks.NewTSS(t)
		musecore = mocks.NewMusecoreClient(t).WithKeys(&keys.Keys{})

		testLogger = testlog.New(t)
		logger     = base.Logger{Std: testLogger.Logger, Compliance: testLogger.Logger}
	)

	suiMock := mocks.NewSuiClient(t)

	gw, err := sui.NewGatewayFromPairID(chainParams.GatewayAddress)
	require.NoError(t, err)

	baseSigner := base.NewSigner(chain, tss, logger)
	signer := New(baseSigner, suiMock, gw, musecore)

	ts := &testSuite{
		t:        t,
		Ctx:      ctx,
		Chain:    chain,
		TSS:      tss,
		Musecore: musecore,
		SuiMock:  suiMock,
		Gateway:  gw,
		Signer:   signer,
	}

	// Setup mocks
	ts.Musecore.On("Chain").Return(chain).Maybe()

	ts.setupTrackersBag()

	return ts
}

func (ts *testSuite) MockGatewayNonce(nonce uint64) {
	ts.SuiMock.On("GetObjectParsedData", mock.Anything, mock.Anything).Return(models.SuiParsedData{
		SuiMoveObject: models.SuiMoveObject{
			Fields: map[string]any{"nonce": fmt.Sprintf("%d", nonce)},
		},
	}, nil)
}

func (ts *testSuite) MockWithdrawCapID(id string) {
	tss, structType := ts.TSS.PubKey().AddressSui(), ts.Gateway.WithdrawCapType()
	ts.SuiMock.On("GetOwnedObjectID", mock.Anything, tss, structType).Return(id, nil)
}

func (ts *testSuite) MockMoveCall(assert func(req models.MoveCallRequest), txBytesBase64 string) {
	call := func(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error) {
		assert(req)
		return models.TxnMetaData{TxBytes: txBytesBase64}, nil
	}

	ts.SuiMock.On("MoveCall", mock.Anything, mock.Anything).Return(call)
}

func (ts *testSuite) MockExec(assert func(req models.SuiExecuteTransactionBlockRequest), digest string) {
	call := func(
		ctx context.Context,
		req models.SuiExecuteTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error) {
		assert(req)
		return models.SuiTransactionBlockResponse{
			Effects: models.SuiEffects{
				Status: models.ExecutionStatus{
					Status: client.TxStatusSuccess,
				},
			},
			Digest: digest,
		}, nil
	}

	ts.SuiMock.On("SuiExecuteTransactionBlock", mock.Anything, mock.Anything).Return(call)
}

type testTracker struct {
	nonce uint64
	hash  string
}

func (ts *testSuite) setupTrackersBag() {
	catcher := func(args mock.Arguments) {
		require.Equal(ts.t, ts.Chain.ChainId, args.Get(1).(int64))
		nonce := args.Get(2).(uint64)
		txHash := args.Get(3).(string)

		ts.t.Logf("Adding outbound tracker: nonce=%d, hash=%s", nonce, txHash)

		ts.TrackerBag = append(ts.TrackerBag, testTracker{nonce, txHash})
	}

	ts.Musecore.On(
		"PostOutboundTracker",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Maybe().Run(catcher).Return("", nil)
}
