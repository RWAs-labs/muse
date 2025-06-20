package backend

import (
	"fmt"
	"math/big"

	tmlog "cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/ethermint/indexer"
	ethermint "github.com/RWAs-labs/ethermint/types"
	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/RWAs-labs/muse/rpc/backend/mocks"
	rpctypes "github.com/RWAs-labs/muse/rpc/types"
	"github.com/RWAs-labs/muse/testutil/sample"
)

func (suite *BackendTestSuite) TestGetSyntheticTransactionByHash() {
	hash := sample.Hash().Hex()
	_, txRes := suite.buildSyntheticTxResult(hash)

	suite.backend.indexer = nil
	client := suite.backend.clientCtx.Client.(*mocks.Client)
	query := fmt.Sprintf(
		"%s.%s='%s'",
		evmtypes.TypeMsgEthereumTx,
		evmtypes.AttributeKeyEthereumTxHash,
		common.HexToHash(hash).Hex(),
	)
	queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
	RegisterBaseFee(queryClient, sdkmath.NewInt(1))
	RegisterTxSearchWithTxResult(client, query, []byte{}, txRes)
	RegisterBlock(client, 1, nil)
	RegisterBlockResultsWithTxResults(client, 1, []*abci.ExecTxResult{&txRes})

	res, err := suite.backend.GetTransactionByHash(common.HexToHash(hash))
	suite.Require().NoError(err)

	// assert fields
	suite.Require().Equal(hash, res.Hash.Hex())
	nonce, _ := hexutil.DecodeUint64(res.Nonce.String())
	suite.Require().Equal(uint64(1), nonce)
	suite.Require().Equal(int64(1), res.BlockNumber.ToInt().Int64())
	suite.Require().Equal("0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7", res.To.Hex())
	suite.Require().Equal("0x735b14BB79463307AAcBED86DAf3322B1e6226aB", res.From.Hex())
	txIndex, _ := hexutil.DecodeUint64(res.TransactionIndex.String())
	suite.Require().Equal(uint64(8888), txIndex)
	txType, _ := hexutil.DecodeUint64(res.Type.String())
	suite.Require().Equal(uint64(88), txType)
	suite.Require().Equal(int64(7001), res.ChainID.ToInt().Int64())
	suite.Require().Equal(int64(1000), res.Value.ToInt().Int64())
	gas, _ := hexutil.DecodeUint64(res.Gas.String())
	suite.Require().Equal(uint64(21000), gas)
	suite.Require().Equal("0x1234", res.Input.String())
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), res.V)
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), res.R)
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), res.S)
}

func (suite *BackendTestSuite) TestGetSyntheticTransactionReceiptByHash() {
	hash := sample.Hash().Hex()
	_, txRes := suite.buildSyntheticTxResult(hash)

	suite.backend.indexer = nil
	client := suite.backend.clientCtx.Client.(*mocks.Client)
	query := fmt.Sprintf(
		"%s.%s='%s'",
		evmtypes.TypeMsgEthereumTx,
		evmtypes.AttributeKeyEthereumTxHash,
		common.HexToHash(hash).Hex(),
	)
	RegisterTxSearchWithTxResult(client, query, []byte{}, txRes)
	RegisterBlock(client, 1, nil)
	RegisterBlockResultsWithTxResults(client, 1, []*abci.ExecTxResult{&txRes})

	res, err := suite.backend.GetTransactionReceipt(common.HexToHash(hash))
	suite.Require().NoError(err)

	// assert fields
	suite.Require().Equal(common.HexToHash(hash), res["transactionHash"])
	blockNumber, _ := hexutil.DecodeUint64(res["blockNumber"].(hexutil.Uint64).String())
	suite.Require().Equal(uint64(1), blockNumber)
	toAddress := common.HexToAddress("0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7")
	fromAddress := common.HexToAddress("0x735b14BB79463307AAcBED86DAf3322B1e6226aB")
	suite.Require().Equal(&toAddress, res["to"])
	suite.Require().Equal(fromAddress, res["from"])
	status, _ := hexutil.DecodeUint64(res["status"].(hexutil.Uint).String())
	suite.Require().Equal(uint64(1), status)
	txType, _ := hexutil.DecodeUint64(res["type"].(hexutil.Uint).String())
	suite.Require().Equal(uint64(88), txType)
	txIndex, _ := hexutil.DecodeUint64(res["transactionIndex"].(hexutil.Uint64).String())
	suite.Require().Equal(uint64(8888), txIndex)
}

func (suite *BackendTestSuite) TestGetSyntheticTransactionByBlockNumberAndIndex() {
	hash := sample.Hash().Hex()
	tx, txRes := suite.buildSyntheticTxResult(hash)

	suite.backend.indexer = nil
	client := suite.backend.clientCtx.Client.(*mocks.Client)
	queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
	RegisterBlock(client, 1, []types.Tx{tx})
	RegisterBlockResultsWithTxResults(client, 1, []*abci.ExecTxResult{&txRes})
	RegisterBaseFee(queryClient, sdkmath.NewInt(1))

	res, err := suite.backend.GetTransactionByBlockNumberAndIndex(rpctypes.BlockNumber(1), 0)
	suite.Require().NoError(err)

	// assert fields
	suite.Require().Equal(hash, res.Hash.Hex())
	nonce, _ := hexutil.DecodeUint64(res.Nonce.String())
	suite.Require().Equal(uint64(1), nonce)
	suite.Require().Equal("0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7", res.To.Hex())
	suite.Require().Equal("0x735b14BB79463307AAcBED86DAf3322B1e6226aB", res.From.Hex())
	txType, _ := hexutil.DecodeUint64(res.Type.String())
	suite.Require().Equal(uint64(88), txType)
	suite.Require().Equal(int64(7001), res.ChainID.ToInt().Int64())
	suite.Require().Equal(int64(1000), res.Value.ToInt().Int64())
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), res.V)
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), res.R)
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), res.S)
}

func (suite *BackendTestSuite) TestGetSyntheticTransactionByBlockNumberAndIndexWithRealTransaction() {
	hash := sample.Hash().Hex()
	tx, txRes := suite.buildSyntheticTxResult(hash)
	msgEthereumTx, _ := suite.buildEthereumTx()

	realTx := suite.signAndEncodeEthTx(msgEthereumTx)

	suite.backend.indexer = nil
	client := suite.backend.clientCtx.Client.(*mocks.Client)
	queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
	// synthetic tx with real MsgEthereumTx
	RegisterBlock(client, 1, []types.Tx{realTx, tx})
	RegisterBlockResultsWithTxResults(client, 1, []*abci.ExecTxResult{{}, &txRes})
	RegisterBaseFee(queryClient, sdkmath.NewInt(1))

	res, err := suite.backend.GetTransactionByBlockNumberAndIndex(rpctypes.BlockNumber(1), 1)
	suite.Require().NoError(err)

	// assert fields
	suite.Require().Equal(hash, res.Hash.Hex())
	nonce, _ := hexutil.DecodeUint64(res.Nonce.String())
	suite.Require().Equal(uint64(1), nonce)
	suite.Require().Equal("0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7", res.To.Hex())
	suite.Require().Equal("0x735b14BB79463307AAcBED86DAf3322B1e6226aB", res.From.Hex())
	txType, _ := hexutil.DecodeUint64(res.Type.String())
	suite.Require().Equal(uint64(88), txType)
	suite.Require().Equal(int64(7001), res.ChainID.ToInt().Int64())
	suite.Require().Equal(int64(1000), res.Value.ToInt().Int64())
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), res.V)
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), res.R)
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), res.S)
}

func (suite *BackendTestSuite) TestGetTransactionByHash() {
	msgEthereumTx, _ := suite.buildEthereumTx()
	txHash := msgEthereumTx.AsTransaction().Hash()

	txBz := suite.signAndEncodeEthTx(msgEthereumTx)
	block := &types.Block{Header: types.Header{Height: 1, ChainID: "test"}, Data: types.Data{Txs: []types.Tx{txBz}}}
	responseDeliver := []*abci.ExecTxResult{
		{
			Code: 0,
			Events: []abci.Event{
				{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
					{Key: "ethereumTxHash", Value: txHash.Hex()},
					{Key: "txIndex", Value: "0"},
					{Key: "amount", Value: "1000"},
					{Key: "txGasUsed", Value: "21000"},
					{Key: "txHash", Value: ""},
					{Key: "recipient", Value: ""},
				}},
			},
		},
	}

	rpcTransaction, err := rpctypes.NewRPCTransaction(
		msgEthereumTx.AsTransaction(),
		common.HexToHash("0x1"),
		1,
		0,
		big.NewInt(1),
		suite.backend.chainID,
	)
	suite.Require().NoError(err)

	testCases := []struct {
		name         string
		registerMock func()
		tx           *evmtypes.MsgEthereumTx
		expRPCTx     *rpctypes.RPCTransaction
		expPass      bool
	}{
		{
			"fail - Block error",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterBlockError(client, 1)
			},
			msgEthereumTx,
			rpcTransaction,
			false,
		},
		{
			"fail - Block Result error",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterBlock(client, 1, []types.Tx{txBz})
				RegisterBlockResultsError(client, 1)
			},
			msgEthereumTx,
			nil,
			true,
		},
		{
			"pass - Base fee error",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBlock(client, 1, []types.Tx{txBz})
				RegisterBlockResults(client, 1)
				RegisterBaseFeeError(queryClient)
			},
			msgEthereumTx,
			rpcTransaction,
			true,
		},
		{
			"pass - Transaction found and returned",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBlock(client, 1, []types.Tx{txBz})
				RegisterBlockResults(client, 1)
				RegisterBaseFee(queryClient, sdkmath.NewInt(1))
			},
			msgEthereumTx,
			rpcTransaction,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.registerMock()

			db := dbm.NewMemDB()
			suite.backend.indexer = indexer.NewKVIndexer(db, tmlog.NewNopLogger(), suite.backend.clientCtx)
			err := suite.backend.indexer.IndexBlock(block, responseDeliver)
			suite.Require().NoError(err)

			rpcTx, err := suite.backend.GetTransactionByHash(common.HexToHash(tc.tx.Hash))

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(rpcTx, tc.expRPCTx)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestGetTransactionsByHashPending() {
	msgEthereumTx, bz := suite.buildEthereumTx()
	rpcTransaction, err := rpctypes.NewRPCTransaction(
		msgEthereumTx.AsTransaction(),
		common.Hash{},
		0,
		0,
		big.NewInt(1),
		suite.backend.chainID,
	)
	suite.Require().NoError(err)

	testCases := []struct {
		name         string
		registerMock func()
		tx           *evmtypes.MsgEthereumTx
		expRPCTx     *rpctypes.RPCTransaction
		expPass      bool
	}{
		{
			"fail - Pending transactions returns error",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterUnconfirmedTxsError(client, nil)
			},
			msgEthereumTx,
			nil,
			true,
		},
		{
			"fail - Tx not found return nil",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterUnconfirmedTxs(client, nil, nil)
			},
			msgEthereumTx,
			nil,
			true,
		},
		{
			"pass - Tx found and returned",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterUnconfirmedTxs(client, nil, types.Txs{bz})
			},
			msgEthereumTx,
			rpcTransaction,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.registerMock()

			rpcTx, err := suite.backend.getTransactionByHashPending(common.HexToHash(tc.tx.Hash))

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(rpcTx, tc.expRPCTx)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestGetTxByEthHash() {
	msgEthereumTx, bz := suite.buildEthereumTx()
	rpcTransaction, err := rpctypes.NewRPCTransaction(
		msgEthereumTx.AsTransaction(),
		common.Hash{},
		0,
		0,
		big.NewInt(1),
		suite.backend.chainID,
	)
	suite.Require().NoError(err)

	testCases := []struct {
		name         string
		registerMock func()
		tx           *evmtypes.MsgEthereumTx
		expRPCTx     *rpctypes.RPCTransaction
		expPass      bool
	}{
		{
			"fail - Indexer disabled can't find transaction",
			func() {
				suite.backend.indexer = nil
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				query := fmt.Sprintf(
					"%s.%s='%s'",
					evmtypes.TypeMsgEthereumTx,
					evmtypes.AttributeKeyEthereumTxHash,
					common.HexToHash(msgEthereumTx.Hash).Hex(),
				)
				RegisterTxSearch(client, query, bz)
			},
			msgEthereumTx,
			rpcTransaction,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.registerMock()

			rpcTx, _, err := suite.backend.GetTxByEthHash(common.HexToHash(tc.tx.Hash))

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(rpcTx, tc.expRPCTx)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestGetTransactionByBlockHashAndIndex() {
	_, bz := suite.buildEthereumTx()

	testCases := []struct {
		name         string
		registerMock func()
		blockHash    common.Hash
		expRPCTx     *rpctypes.RPCTransaction
		expPass      bool
	}{
		{
			"pass - block not found",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterBlockByHashError(client, common.Hash{}, bz)
			},
			common.Hash{},
			nil,
			true,
		},
		{
			"pass - Block results error",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterBlockByHash(client, common.Hash{}, bz)
				RegisterBlockResultsError(client, 1)
			},
			common.Hash{},
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.registerMock()

			rpcTx, err := suite.backend.GetTransactionByBlockHashAndIndex(tc.blockHash, 1)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(rpcTx, tc.expRPCTx)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestGetTransactionByBlockAndIndex() {
	msgEthTx, bz := suite.buildEthereumTx()

	defaultBlock := types.MakeBlock(1, []types.Tx{bz}, nil, nil)
	defaultExecTxResult := []*abci.ExecTxResult{
		{
			Code: 0,
			Events: []abci.Event{
				{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
					{Key: "ethereumTxHash", Value: common.HexToHash(msgEthTx.Hash).Hex()},
					{Key: "txIndex", Value: "0"},
					{Key: "amount", Value: "1000"},
					{Key: "txGasUsed", Value: "21000"},
					{Key: "txHash", Value: ""},
					{Key: "recipient", Value: ""},
				}},
			},
		},
	}

	txFromMsg, err := rpctypes.NewTransactionFromMsg(
		msgEthTx,
		common.BytesToHash(defaultBlock.Hash().Bytes()),
		1,
		0,
		big.NewInt(1),
		suite.backend.chainID,
		nil,
	)
	suite.Require().NoError(err)

	testCases := []struct {
		name         string
		registerMock func()
		block        *tmrpctypes.ResultBlock
		idx          hexutil.Uint
		expRPCTx     *rpctypes.RPCTransaction
		expPass      bool
	}{
		{
			"pass - block txs index out of bound ",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterBlockResults(client, 1)
			},
			&tmrpctypes.ResultBlock{Block: types.MakeBlock(1, []types.Tx{bz}, nil, nil)},
			1,
			nil,
			true,
		},
		{
			"pass - Can't fetch base fee",
			func() {
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterBlockResults(client, 1)
				RegisterBaseFeeError(queryClient)
			},
			&tmrpctypes.ResultBlock{Block: defaultBlock},
			0,
			txFromMsg,
			true,
		},
		{
			"pass - Gets Tx by transaction index",
			func() {
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				db := dbm.NewMemDB()
				suite.backend.indexer = indexer.NewKVIndexer(db, tmlog.NewNopLogger(), suite.backend.clientCtx)
				txBz := suite.signAndEncodeEthTx(msgEthTx)
				block := &types.Block{
					Header: types.Header{Height: 1, ChainID: "test"},
					Data:   types.Data{Txs: []types.Tx{txBz}},
				}
				err := suite.backend.indexer.IndexBlock(block, defaultExecTxResult)
				suite.Require().NoError(err)
				RegisterBlockResults(client, 1)
				RegisterBaseFee(queryClient, sdkmath.NewInt(1))
			},
			&tmrpctypes.ResultBlock{Block: defaultBlock},
			0,
			txFromMsg,
			true,
		},
		{
			"pass - returns the Ethereum format transaction by the Ethereum hash",
			func() {
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterBlockResults(client, 1)
				RegisterBaseFee(queryClient, sdkmath.NewInt(1))
			},
			&tmrpctypes.ResultBlock{Block: defaultBlock},
			0,
			txFromMsg,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.registerMock()

			rpcTx, err := suite.backend.GetTransactionByBlockAndIndex(tc.block, tc.idx)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(rpcTx, tc.expRPCTx)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestGetTransactionByBlockNumberAndIndex() {
	msgEthTx, bz := suite.buildEthereumTx()
	defaultBlock := types.MakeBlock(1, []types.Tx{bz}, nil, nil)
	txFromMsg, err := rpctypes.NewTransactionFromMsg(
		msgEthTx,
		common.BytesToHash(defaultBlock.Hash().Bytes()),
		1,
		0,
		big.NewInt(1),
		suite.backend.chainID,
		nil,
	)
	suite.Require().NoError(err)

	testCases := []struct {
		name         string
		registerMock func()
		blockNum     rpctypes.BlockNumber
		idx          hexutil.Uint
		expRPCTx     *rpctypes.RPCTransaction
		expPass      bool
	}{
		{
			"fail -  block not found return nil",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterBlockError(client, 1)
			},
			0,
			0,
			nil,
			true,
		},
		{
			"pass - returns the transaction identified by block number and index",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBlock(client, 1, []types.Tx{bz})
				RegisterBlockResults(client, 1)
				RegisterBaseFee(queryClient, sdkmath.NewInt(1))
			},
			0,
			0,
			txFromMsg,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.registerMock()

			rpcTx, err := suite.backend.GetTransactionByBlockNumberAndIndex(tc.blockNum, tc.idx)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(rpcTx, tc.expRPCTx)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestGetTransactionByTxIndex() {
	_, bz := suite.buildEthereumTx()

	testCases := []struct {
		name         string
		registerMock func()
		height       int64
		index        uint
		expTxResult  *ethermint.TxResult
		expPass      bool
	}{
		{
			"fail - Ethereum tx with query not found",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				suite.backend.indexer = nil
				RegisterTxSearch(client, "tx.height=0 AND ethereum_tx.txIndex=0", bz)
			},
			0,
			0,
			&ethermint.TxResult{},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.registerMock()

			txResults, _, err := suite.backend.GetTxByTxIndex(tc.height, tc.index)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(txResults, tc.expTxResult)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestQueryTendermintTxIndexer() {
	testCases := []struct {
		name         string
		registerMock func()
		txGetter     func(*rpctypes.ParsedTxs) *rpctypes.ParsedTx
		query        string
		expTxResult  *ethermint.TxResult
		expPass      bool
	}{
		{
			"fail - Ethereum tx with query not found",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterTxSearchEmpty(client, "")
			},
			func(txs *rpctypes.ParsedTxs) *rpctypes.ParsedTx {
				return &rpctypes.ParsedTx{}
			},
			"",
			&ethermint.TxResult{},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.registerMock()

			txResults, _, err := suite.backend.queryTendermintTxIndexer(tc.query, tc.txGetter)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(txResults, tc.expTxResult)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestGetTransactionReceipt() {
	msgEthereumTx, _ := suite.buildEthereumTx()
	txHash := msgEthereumTx.AsTransaction().Hash()

	txBz := suite.signAndEncodeEthTx(msgEthereumTx)

	testCases := []struct {
		name         string
		registerMock func()
		tx           *evmtypes.MsgEthereumTx
		block        *types.Block
		blockResult  []*abci.ExecTxResult
		expTxReceipt map[string]interface{}
		expPass      bool
	}{
		{
			"fail - Receipts do not match ",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterBlock(client, 1, []types.Tx{txBz})
				RegisterBlockResults(client, 1)
			},
			msgEthereumTx,
			&types.Block{Header: types.Header{Height: 1}, Data: types.Data{Txs: []types.Tx{txBz}}},
			[]*abci.ExecTxResult{
				{
					Code: 0,
					Events: []abci.Event{
						{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
							{Key: "ethereumTxHash", Value: txHash.Hex()},
							{Key: "txIndex", Value: "0"},
							{Key: "amount", Value: "1000"},
							{Key: "txGasUsed", Value: "21000"},
							{Key: "txHash", Value: ""},
							{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
						}},
					},
				},
			},
			map[string]interface{}(nil),
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
			tc.registerMock()

			db := dbm.NewMemDB()
			suite.backend.indexer = indexer.NewKVIndexer(db, tmlog.NewNopLogger(), suite.backend.clientCtx)
			err := suite.backend.indexer.IndexBlock(tc.block, tc.blockResult)
			suite.Require().NoError(err)

			txReceipt, err := suite.backend.GetTransactionReceipt(common.HexToHash(tc.tx.Hash))
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(txReceipt, tc.expTxReceipt)
			} else {
				suite.Require().NotEqual(txReceipt, tc.expTxReceipt)
			}
		})
	}
}
