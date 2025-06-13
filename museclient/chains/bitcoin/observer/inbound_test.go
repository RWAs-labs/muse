package observer

import (
	"bytes"
	"context"
	"encoding/hex"
	"math"
	"math/big"
	"path"
	"testing"

	cosmosmath "cosmossdk.io/math"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/memo"

	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/common"
	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/museclient/keys"
	"github.com/RWAs-labs/muse/museclient/testutils"
	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	"github.com/RWAs-labs/muse/museclient/testutils/testrpc"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/testutil"
	"github.com/RWAs-labs/muse/testutil/sample"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// mockDepositFeeCalculator returns a mock depositor fee calculator that returns the given fee and error.
func mockDepositFeeCalculator(fee float64, err error) common.DepositorFeeCalculator {
	return func(_ context.Context, _ common.RPC, _ *btcjson.TxRawResult, _ *chaincfg.Params) (float64, error) {
		return fee, err
	}
}

func TestAvgFeeRateBlock828440(t *testing.T) {
	// load archived block 828440
	var blockVb btcjson.GetBlockVerboseTxResult
	testutils.LoadObjectFromJSONFile(
		t,
		&blockVb,
		path.Join(TestDataDir, testutils.TestDataPathBTC, "block_trimmed_8332_828440.json"),
	)

	// https://mempool.space/block/000000000000000000025ca01d2c1094b8fd3bacc5468cc3193ced6a14618c27
	var blockMb testutils.MempoolBlock
	testutils.LoadObjectFromJSONFile(
		t,
		&blockMb,
		path.Join(TestDataDir, testutils.TestDataPathBTC, "block_mempool.space_8332_828440.json"),
	)

	gasRate, err := common.CalcBlockAvgFeeRate(&blockVb, &chaincfg.MainNetParams)
	require.NoError(t, err)
	require.Equal(t, int64(blockMb.Extras.AvgFeeRate), gasRate)
}

func TestAvgFeeRateBlock828440Errors(t *testing.T) {
	// load archived block 828440
	var blockVb btcjson.GetBlockVerboseTxResult
	testutils.LoadObjectFromJSONFile(
		t,
		&blockVb,
		path.Join(TestDataDir, testutils.TestDataPathBTC, "block_trimmed_8332_828440.json"),
	)

	t.Run("block has no transactions", func(t *testing.T) {
		emptyVb := btcjson.GetBlockVerboseTxResult{Tx: []btcjson.TxRawResult{}}
		_, err := common.CalcBlockAvgFeeRate(&emptyVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "block has no transactions")
	})
	t.Run("it's okay if block has only coinbase tx", func(t *testing.T) {
		coinbaseVb := btcjson.GetBlockVerboseTxResult{Tx: []btcjson.TxRawResult{
			blockVb.Tx[0],
		}}
		_, err := common.CalcBlockAvgFeeRate(&coinbaseVb, &chaincfg.MainNetParams)
		require.NoError(t, err)
	})
	t.Run("tiny block weight should fail", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Weight = 3
		_, err := common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "block weight 3 too small")
	})
	t.Run("block weight should not be less than coinbase tx weight", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Weight = blockVb.Tx[0].Weight - 1
		_, err := common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "less than coinbase tx weight")
	})
	t.Run("invalid block height should fail", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Height = 0
		_, err := common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid block height")

		invalidVb.Height = math.MaxInt32 + 1
		_, err = common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid block height")
	})
	t.Run("failed to decode coinbase tx", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Tx = []btcjson.TxRawResult{blockVb.Tx[0], blockVb.Tx[1]}
		invalidVb.Tx[0].Hex = "invalid hex"
		_, err := common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "failed to decode coinbase tx")
	})
	t.Run("1st tx is not coinbase", func(t *testing.T) {
		invalidVb := blockVb
		invalidVb.Tx = []btcjson.TxRawResult{blockVb.Tx[1], blockVb.Tx[0]}
		_, err := common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "not coinbase tx")
	})
	t.Run("miner earned less than subsidy", func(t *testing.T) {
		invalidVb := blockVb
		coinbaseTxBytes, err := hex.DecodeString(blockVb.Tx[0].Hex)
		require.NoError(t, err)
		coinbaseTx, err := btcutil.NewTxFromBytes(coinbaseTxBytes)
		require.NoError(t, err)
		msgTx := coinbaseTx.MsgTx()

		// reduce subsidy by 1 satoshi
		for i := range msgTx.TxOut {
			if i == 0 {
				msgTx.TxOut[i].Value = blockchain.CalcBlockSubsidy(int32(blockVb.Height), &chaincfg.MainNetParams) - 1
			} else {
				msgTx.TxOut[i].Value = 0
			}
		}
		// calculate fee rate on modified coinbase tx
		var buf bytes.Buffer
		err = msgTx.Serialize(&buf)
		require.NoError(t, err)
		invalidVb.Tx[0].Hex = hex.EncodeToString(buf.Bytes())
		_, err = common.CalcBlockAvgFeeRate(&invalidVb, &chaincfg.MainNetParams)
		require.Error(t, err)
		require.ErrorContains(t, err, "less than subsidy")
	})
}

func Test_GetInboundVoteFromBtcEvent(t *testing.T) {
	r := sample.Rand()

	// can use any bitcoin chain for testing
	chain := chains.BitcoinMainnet

	// create test observer
	ob := newTestSuite(t, chain)
	ob.musecore.WithKeys(&keys.Keys{}).WithMuseChain()

	// test cases
	tests := []struct {
		name              string
		event             *BTCInboundEvent
		observationStatus crosschaintypes.InboundStatus
		nilVote           bool
	}{
		{
			name: "should return vote for standard memo",
			event: &BTCInboundEvent{
				FromAddress: sample.BTCAddressP2WPKH(t, r, &chaincfg.MainNetParams).String(),
				// a deposit and call
				MemoBytes: testutil.HexToBytes(
					t,
					"5a0110032d07a9cbd57dcca3e2cf966c88bc874445b6e3b60d68656c6c6f207361746f736869",
				),
			},
			observationStatus: crosschaintypes.InboundStatus_SUCCESS,
		},
		{
			name: "should return vote for legacy memo",
			event: &BTCInboundEvent{
				// raw address + payload
				MemoBytes: testutil.HexToBytes(t, "2d07a9cbd57dcca3e2cf966c88bc874445b6e3b668656c6c6f207361746f736869"),
			},
			observationStatus: crosschaintypes.InboundStatus_SUCCESS,
		},
		{
			name: "should return vote for invalid memo",
			event: &BTCInboundEvent{
				// standard memo that carries payload only, receiver address is empty
				MemoBytes: testutil.HexToBytes(t, "5a0110020d68656c6c6f207361746f736869"),
			},
			observationStatus: crosschaintypes.InboundStatus_INVALID_MEMO,
		},
		{
			name: "should return nil on donation message",
			event: &BTCInboundEvent{
				MemoBytes: []byte(constant.DonationMessage),
			},
			nilVote: true,
		},
		{
			name: "should return nil on invalid deposit value",
			event: &BTCInboundEvent{
				Value:     -1, // invalid value
				MemoBytes: testutil.HexToBytes(t, "2d07a9cbd57dcca3e2cf966c88bc874445b6e3b668656c6c6f207361746f736869"),
			},
			nilVote: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ob.GetInboundVoteFromBtcEvent(tt.event)
			if tt.nilVote {
				require.Nil(t, msg)
			} else {
				require.NotNil(t, msg)
				require.EqualValues(t, tt.observationStatus, msg.Status)
			}
		})
	}
}

func TestGetSenderAddressByVin(t *testing.T) {
	ctx := context.Background()

	// https://mempool.space/tx/3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867
	txHash := "3618e869f9e87863c0f1cc46dbbaa8b767b4a5d6d60b143c2c50af52b257e867"
	chain := chains.BitcoinMainnet
	net := &chaincfg.MainNetParams

	t.Run("should get sender address from tx", func(t *testing.T) {
		// vin from the archived P2WPKH tx
		// https://mempool.space/tx/c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697
		txHash := "c5d224963832fc0b9a597251c2342a17b25e481a88cc9119008e8f8296652697"
		rpcClient := testrpc.CreateBTCRPCAndLoadTx(t, TestDataDir, chain.ChainId, txHash)

		// get sender address
		txVin := btcjson.Vin{Txid: txHash, Vout: 2}
		sender, err := GetSenderAddressByVin(ctx, rpcClient, txVin, net)
		require.NoError(t, err)
		require.Equal(t, "bc1q68kxnq52ahz5vd6c8czevsawu0ux9nfrzzrh6e", sender)
	})

	t.Run("should return error on invalid txHash", func(t *testing.T) {
		rpcClient := mocks.NewBitcoinClient(t)
		// use invalid tx hash
		txVin := btcjson.Vin{Txid: "invalid tx hash", Vout: 2}
		sender, err := GetSenderAddressByVin(ctx, rpcClient, txVin, net)
		require.Error(t, err)
		require.Empty(t, sender)
	})

	t.Run("should return error when RPC client fails to get raw tx", func(t *testing.T) {
		// create mock rpc client that returns rpc error
		rpcClient := mocks.NewBitcoinClient(t)
		rpcClient.On("GetRawTransaction", mock.Anything, mock.Anything).Return(nil, errors.New("rpc error"))

		// get sender address
		txVin := btcjson.Vin{Txid: txHash, Vout: 2}
		sender, err := GetSenderAddressByVin(ctx, rpcClient, txVin, net)
		require.ErrorContains(t, err, "error getting raw transaction")
		require.Empty(t, sender)
	})

	t.Run("should return error on invalid output index", func(t *testing.T) {
		// create mock rpc client with preloaded tx
		rpcClient := testrpc.CreateBTCRPCAndLoadTx(t, TestDataDir, chain.ChainId, txHash)

		// invalid output index
		txVin := btcjson.Vin{Txid: txHash, Vout: 3}
		sender, err := GetSenderAddressByVin(ctx, rpcClient, txVin, net)
		require.ErrorContains(t, err, "out of range")
		require.Empty(t, sender)
	})
}

func Test_NewInboundVoteFromLegacyMemo(t *testing.T) {
	// can use any bitcoin chain for testing
	chain := chains.BitcoinMainnet

	// create test observer
	ob := newTestSuite(t, chain)
	ob.musecore.WithKeys(&keys.Keys{}).WithMuseChain()

	t.Run("should create new inbound vote msg V2", func(t *testing.T) {
		// create test event
		event := createTestBtcEvent(t, &chaincfg.MainNetParams, []byte("dummy memo"), nil)

		// test amount
		amountSats := big.NewInt(1000)

		// mock SAFE confirmed block
		ob.WithLastBlock(event.BlockNumber + ob.ChainParams().InboundConfirmationSafe())

		// expected vote
		expectedVote := crosschaintypes.MsgVoteInbound{
			Sender:             event.FromAddress,
			SenderChainId:      chain.ChainId,
			TxOrigin:           event.FromAddress,
			Receiver:           event.ToAddress,
			ReceiverChain:      ob.MusecoreClient().Chain().ChainId,
			Amount:             cosmosmath.NewUint(amountSats.Uint64()),
			Message:            hex.EncodeToString(event.MemoBytes),
			InboundHash:        event.TxHash,
			InboundBlockHeight: event.BlockNumber,
			CallOptions: &crosschaintypes.CallOptions{
				GasLimit: 0,
			},
			CoinType:                coin.CoinType_Gas,
			ProtocolContractVersion: crosschaintypes.ProtocolContractVersion_V2,
			RevertOptions:           crosschaintypes.NewEmptyRevertOptions(), // always empty with legacy memo
			IsCrossChainCall:        true,
			Status:                  crosschaintypes.InboundStatus_SUCCESS,
			ConfirmationMode:        crosschaintypes.ConfirmationMode_SAFE,
		}

		// create new inbound vote V1
		vote := ob.NewInboundVoteFromLegacyMemo(&event, amountSats)
		require.Equal(t, expectedVote, *vote)
	})
}

func Test_NewInboundVoteFromStdMemo(t *testing.T) {
	// can use any bitcoin chain for testing
	chain := chains.BitcoinMainnet

	// create test observer
	ob := newTestSuite(t, chain)
	ob.musecore.WithKeys(&keys.Keys{}).WithMuseChain()

	t.Run("should create new inbound vote msg with standard memo", func(t *testing.T) {
		// create revert options
		r := sample.Rand()
		revertOptions := crosschaintypes.NewEmptyRevertOptions()
		revertOptions.RevertAddress = sample.BTCAddressP2WPKH(t, r, &chaincfg.MainNetParams).String()

		// create test event
		receiver := sample.EthAddress()
		event := createTestBtcEvent(t, &chaincfg.MainNetParams, []byte("dymmy"), &memo.InboundMemo{
			FieldsV0: memo.FieldsV0{
				Receiver:      receiver,
				Payload:       []byte("some payload"),
				RevertOptions: revertOptions,
			},
		})

		// test amount
		amountSats := big.NewInt(1000)

		// mock SAFE confirmed block
		ob.WithLastBlock(event.BlockNumber + ob.ChainParams().InboundConfirmationSafe())

		// expected vote
		memoBytesExpected := event.MemoStd.Payload
		expectedVote := crosschaintypes.MsgVoteInbound{
			Sender:             event.FromAddress,
			SenderChainId:      chain.ChainId,
			TxOrigin:           event.FromAddress,
			Receiver:           event.MemoStd.Receiver.Hex(),
			ReceiverChain:      ob.MusecoreClient().Chain().ChainId,
			Amount:             cosmosmath.NewUint(amountSats.Uint64()),
			Message:            hex.EncodeToString(memoBytesExpected), // a simulated legacy memo
			InboundHash:        event.TxHash,
			InboundBlockHeight: event.BlockNumber,
			CallOptions: &crosschaintypes.CallOptions{
				GasLimit: 0,
			},
			CoinType:                coin.CoinType_Gas,
			ProtocolContractVersion: crosschaintypes.ProtocolContractVersion_V2,
			RevertOptions: crosschaintypes.RevertOptions{
				RevertAddress: revertOptions.RevertAddress, // should be overridden by revert address
			},
			Status:           crosschaintypes.InboundStatus_SUCCESS,
			ConfirmationMode: crosschaintypes.ConfirmationMode_SAFE,
		}

		// create new inbound vote V2 with standard memo
		vote := ob.NewInboundVoteFromStdMemo(&event, amountSats)
		require.Equal(t, expectedVote, *vote)
	})
}
