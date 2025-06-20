package backend

import (
	"context"
	"math/big"
	"testing"

	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	tmrpcclient "github.com/cometbft/cometbft/rpc/client"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/rpc/backend/mocks"
	rpc "github.com/RWAs-labs/muse/rpc/types"
)

// Client defines a mocked object that implements the Tendermint JSON-RPC Client
// interface. It allows for performing Client queries without having to run a
// Tendermint RPC Client server.
//
// To use a mock method it has to be registered in a given test.
var _ tmrpcclient.Client = &mocks.Client{}

// Tx Search
func RegisterTxSearch(client *mocks.Client, query string, txBz []byte) {
	resulTxs := []*tmrpctypes.ResultTx{{Tx: txBz}}
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(&tmrpctypes.ResultTxSearch{Txs: resulTxs, TotalCount: 1}, nil)
}

func RegisterTxSearchWithTxResult(client *mocks.Client, query string, txBz []byte, res abci.ExecTxResult) {
	resulTxs := []*tmrpctypes.ResultTx{{Tx: txBz, Height: 1, TxResult: res}}
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(&tmrpctypes.ResultTxSearch{Txs: resulTxs, TotalCount: 1}, nil)
}

func RegisterTxSearchEmpty(client *mocks.Client, query string) {
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(&tmrpctypes.ResultTxSearch{}, nil)
}

func RegisterTxSearchError(client *mocks.Client, query string) {
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(nil, errortypes.ErrInvalidRequest)
}

// Broadcast Tx
func RegisterBroadcastTx(client *mocks.Client, tx types.Tx) {
	client.On("BroadcastTxSync", context.Background(), tx).
		Return(&tmrpctypes.ResultBroadcastTx{}, nil)
}

func RegisterBroadcastTxAny(client *mocks.Client) {
	client.On("BroadcastTxSync", context.Background(), mock.Anything).
		Return(&tmrpctypes.ResultBroadcastTx{}, nil)
}

func RegisterBroadcastTxError(client *mocks.Client, tx types.Tx) {
	client.On("BroadcastTxSync", context.Background(), tx).
		Return(nil, errortypes.ErrInvalidRequest)
}

// Unconfirmed Transactions
func RegisterUnconfirmedTxs(client *mocks.Client, limit *int, txs []types.Tx) {
	client.On("UnconfirmedTxs", rpc.ContextWithHeight(1), limit).
		Return(&tmrpctypes.ResultUnconfirmedTxs{Txs: txs}, nil)
}

func RegisterUnconfirmedTxsEmpty(client *mocks.Client, limit *int) {
	client.On("UnconfirmedTxs", rpc.ContextWithHeight(1), limit).
		Return(&tmrpctypes.ResultUnconfirmedTxs{
			Txs: make([]types.Tx, 2),
		}, nil)
}

func RegisterUnconfirmedTxsError(client *mocks.Client, limit *int) {
	client.On("UnconfirmedTxs", rpc.ContextWithHeight(1), limit).
		Return(nil, errortypes.ErrInvalidRequest)
}

// Status
func RegisterStatus(client *mocks.Client) {
	client.On("Status", rpc.ContextWithHeight(1)).
		Return(&tmrpctypes.ResultStatus{}, nil)
}

func RegisterStatusError(client *mocks.Client) {
	client.On("Status", rpc.ContextWithHeight(1)).
		Return(nil, errortypes.ErrInvalidRequest)
}

// Block
func RegisterBlock(
	client *mocks.Client,
	height int64,
	txs []types.Tx,
) (*tmrpctypes.ResultBlock, error) {
	// without tx
	if len(txs) == 0 {
		emptyBlock := types.MakeBlock(height, []types.Tx{}, nil, nil)
		emptyBlock.ChainID = ChainID
		blockHash := common.BigToHash(big.NewInt(height)).Bytes()
		resBlock := &tmrpctypes.ResultBlock{Block: emptyBlock, BlockID: types.BlockID{Hash: bytes.HexBytes(blockHash)}}
		client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).Return(resBlock, nil)
		return resBlock, nil
	}

	// with tx
	block := types.MakeBlock(height, txs, nil, nil)
	block.ChainID = ChainID
	blockHash := common.BigToHash(big.NewInt(height)).Bytes()
	resBlock := &tmrpctypes.ResultBlock{Block: block, BlockID: types.BlockID{Hash: bytes.HexBytes(blockHash)}}
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).Return(resBlock, nil)
	return resBlock, nil
}

// Block returns error
func RegisterBlockError(client *mocks.Client, height int64) {
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(nil, errortypes.ErrInvalidRequest)
}

// Block not found
func RegisterBlockNotFound(
	client *mocks.Client,
	height int64,
) (*tmrpctypes.ResultBlock, error) {
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(&tmrpctypes.ResultBlock{Block: nil}, nil)

	return &tmrpctypes.ResultBlock{Block: nil}, nil
}

func TestRegisterBlock(t *testing.T) {
	client := mocks.NewClient(t)
	height := rpc.BlockNumber(1).Int64()
	RegisterBlock(client, height, nil)

	res, err := client.Block(rpc.ContextWithHeight(height), &height)

	emptyBlock := types.MakeBlock(height, []types.Tx{}, nil, nil)
	emptyBlock.ChainID = ChainID
	blockHash := common.BigToHash(big.NewInt(height)).Bytes()
	resBlock := &tmrpctypes.ResultBlock{Block: emptyBlock, BlockID: types.BlockID{Hash: blockHash}}
	require.Equal(t, resBlock, res)
	require.NoError(t, err)
}

// ConsensusParams
func RegisterConsensusParams(client *mocks.Client, height int64) {
	consensusParams := types.DefaultConsensusParams()
	client.On("ConsensusParams", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(&tmrpctypes.ResultConsensusParams{ConsensusParams: *consensusParams}, nil)
}

func RegisterConsensusParamsError(client *mocks.Client, height int64) {
	client.On("ConsensusParams", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(nil, errortypes.ErrInvalidRequest)
}

func TestRegisterConsensusParams(t *testing.T) {
	client := mocks.NewClient(t)
	height := int64(1)
	RegisterConsensusParams(client, height)

	res, err := client.ConsensusParams(rpc.ContextWithHeight(height), &height)
	consensusParams := types.DefaultConsensusParams()
	require.Equal(t, &tmrpctypes.ResultConsensusParams{ConsensusParams: *consensusParams}, res)
	require.NoError(t, err)
}

// BlockResults

func RegisterBlockResultsWithEventLog(client *mocks.Client, height int64) (*tmrpctypes.ResultBlockResults, error) {
	res := &tmrpctypes.ResultBlockResults{
		Height: height,
		TxsResults: []*abci.ExecTxResult{
			{Code: 0, GasUsed: 0, Events: []abci.Event{{
				Type: evmtypes.EventTypeTxLog,
				Attributes: []abci.EventAttribute{{
					Key:   evmtypes.AttributeKeyTxLog,
					Value: "{\"test\": \"hello\"}",
					Index: true,
				}},
			}}},
		},
	}
	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(res, nil)
	return res, nil
}

func RegisterBlockResults(
	client *mocks.Client,
	height int64,
) (*tmrpctypes.ResultBlockResults, error) {
	res := &tmrpctypes.ResultBlockResults{
		Height:     height,
		TxsResults: []*abci.ExecTxResult{{Code: 0, GasUsed: 0}},
	}

	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(res, nil)
	return res, nil
}

func RegisterBlockResultsWithTxResults(
	client *mocks.Client,
	height int64,
	txResults []*abci.ExecTxResult,
) (*tmrpctypes.ResultBlockResults, error) {
	res := &tmrpctypes.ResultBlockResults{
		Height:     height,
		TxsResults: txResults,
	}

	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(res, nil)
	return res, nil
}

func RegisterBlockResultsError(client *mocks.Client, height int64) {
	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(nil, errortypes.ErrInvalidRequest)
}

func TestRegisterBlockResults(t *testing.T) {
	client := mocks.NewClient(t)
	height := int64(1)
	RegisterBlockResults(client, height)

	res, err := client.BlockResults(rpc.ContextWithHeight(height), &height)
	expRes := &tmrpctypes.ResultBlockResults{
		Height:     height,
		TxsResults: []*abci.ExecTxResult{{Code: 0, GasUsed: 0}},
	}
	require.Equal(t, expRes, res)
	require.NoError(t, err)
}

// BlockByHash
func RegisterBlockByHash(
	client *mocks.Client,
	hash common.Hash,
	tx []byte,
) (*tmrpctypes.ResultBlock, error) {
	block := types.MakeBlock(1, []types.Tx{tx}, nil, nil)
	resBlock := &tmrpctypes.ResultBlock{Block: block}

	client.On("BlockByHash", rpc.ContextWithHeight(1), []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(resBlock, nil)
	return resBlock, nil
}

func RegisterBlockByHashError(client *mocks.Client, hash common.Hash, tx []byte) {
	client.On("BlockByHash", rpc.ContextWithHeight(1), []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(nil, errortypes.ErrInvalidRequest)
}

func RegisterBlockByHashNotFound(client *mocks.Client, hash common.Hash, tx []byte) {
	client.On("BlockByHash", rpc.ContextWithHeight(1), []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(nil, nil)
}

func RegisterABCIQueryWithOptions(
	client *mocks.Client,
	height int64,
	path string,
	data bytes.HexBytes,
	opts tmrpcclient.ABCIQueryOptions,
	respValue []byte,
) {
	client.On("ABCIQueryWithOptions", context.Background(), path, data, opts).
		Return(&tmrpctypes.ResultABCIQuery{
			Response: abci.ResponseQuery{
				Value:  respValue,
				Height: height,
			},
		}, nil)
}

func RegisterABCIQueryWithOptionsError(
	clients *mocks.Client,
	path string,
	data bytes.HexBytes,
	opts tmrpcclient.ABCIQueryOptions,
) {
	clients.On("ABCIQueryWithOptions", context.Background(), path, data, opts).
		Return(nil, errortypes.ErrInvalidRequest)
}

func RegisterABCIQueryAccount(
	clients *mocks.Client,
	data bytes.HexBytes,
	opts tmrpcclient.ABCIQueryOptions,
	acc client.Account,
) {
	baseAccount := authtypes.NewBaseAccount(
		acc.GetAddress(),
		acc.GetPubKey(),
		acc.GetAccountNumber(),
		acc.GetSequence(),
	)
	accAny, _ := codectypes.NewAnyWithValue(baseAccount)
	accResponse := authtypes.QueryAccountResponse{Account: accAny}
	respBz, _ := accResponse.Marshal()
	clients.On("ABCIQueryWithOptions", context.Background(), "/cosmos.auth.v1beta1.Query/Account", data, opts).
		Return(&tmrpctypes.ResultABCIQuery{
			Response: abci.ResponseQuery{
				Value:  respBz,
				Height: 1,
			},
		}, nil)
}

func RegisterABCIQuerySimulate(
	clients *mocks.Client,
	opts tmrpcclient.ABCIQueryOptions,
) {
	simResp := &tx.SimulateResponse{
		GasInfo: &sdk.GasInfo{GasWanted: uint64(21000), GasUsed: uint64(21000)},
	}
	respBz, _ := simResp.Marshal()
	clients.On("ABCIQueryWithOptions", context.Background(), "/cosmos.tx.v1beta1.Service/Simulate", mock.Anything, opts).
		Return(&tmrpctypes.ResultABCIQuery{
			Response: abci.ResponseQuery{
				Value:  respBz,
				Height: 1,
			},
		}, nil)
}
