// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	btcjson "github.com/btcsuite/btcd/btcjson"
	btcutil "github.com/btcsuite/btcd/btcutil"

	chainhash "github.com/btcsuite/btcd/chaincfg/chainhash"

	client "github.com/RWAs-labs/muse/museclient/chains/bitcoin/client"

	context "context"

	json "encoding/json"

	mock "github.com/stretchr/testify/mock"

	rpcclient "github.com/btcsuite/btcd/rpcclient"

	time "time"

	wire "github.com/btcsuite/btcd/wire"
)

// BitcoinClient is an autogenerated mock type for the client type
type BitcoinClient struct {
	mock.Mock
}

// CreateWallet provides a mock function with given fields: ctx, name, opts
func (_m *BitcoinClient) CreateWallet(ctx context.Context, name string, opts ...rpcclient.CreateWalletOpt) (*btcjson.CreateWalletResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CreateWallet")
	}

	var r0 *btcjson.CreateWalletResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...rpcclient.CreateWalletOpt) (*btcjson.CreateWalletResult, error)); ok {
		return rf(ctx, name, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...rpcclient.CreateWalletOpt) *btcjson.CreateWalletResult); ok {
		r0 = rf(ctx, name, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcjson.CreateWalletResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...rpcclient.CreateWalletOpt) error); ok {
		r1 = rf(ctx, name, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EstimateSmartFee provides a mock function with given fields: ctx, confTarget, mode
func (_m *BitcoinClient) EstimateSmartFee(ctx context.Context, confTarget int64, mode *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error) {
	ret := _m.Called(ctx, confTarget, mode)

	if len(ret) == 0 {
		panic("no return value specified for EstimateSmartFee")
	}

	var r0 *btcjson.EstimateSmartFeeResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, *btcjson.EstimateSmartFeeMode) (*btcjson.EstimateSmartFeeResult, error)); ok {
		return rf(ctx, confTarget, mode)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, *btcjson.EstimateSmartFeeMode) *btcjson.EstimateSmartFeeResult); ok {
		r0 = rf(ctx, confTarget, mode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcjson.EstimateSmartFeeResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, *btcjson.EstimateSmartFeeMode) error); ok {
		r1 = rf(ctx, confTarget, mode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateToAddress provides a mock function with given fields: ctx, numBlocks, address, maxTries
func (_m *BitcoinClient) GenerateToAddress(ctx context.Context, numBlocks int64, address btcutil.Address, maxTries *int64) ([]*chainhash.Hash, error) {
	ret := _m.Called(ctx, numBlocks, address, maxTries)

	if len(ret) == 0 {
		panic("no return value specified for GenerateToAddress")
	}

	var r0 []*chainhash.Hash
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, btcutil.Address, *int64) ([]*chainhash.Hash, error)); ok {
		return rf(ctx, numBlocks, address, maxTries)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, btcutil.Address, *int64) []*chainhash.Hash); ok {
		r0 = rf(ctx, numBlocks, address, maxTries)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*chainhash.Hash)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, btcutil.Address, *int64) error); ok {
		r1 = rf(ctx, numBlocks, address, maxTries)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBalance provides a mock function with given fields: ctx, account
func (_m *BitcoinClient) GetBalance(ctx context.Context, account string) (btcutil.Amount, error) {
	ret := _m.Called(ctx, account)

	if len(ret) == 0 {
		panic("no return value specified for GetBalance")
	}

	var r0 btcutil.Amount
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (btcutil.Amount, error)); ok {
		return rf(ctx, account)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) btcutil.Amount); ok {
		r0 = rf(ctx, account)
	} else {
		r0 = ret.Get(0).(btcutil.Amount)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockCount provides a mock function with given fields: ctx
func (_m *BitcoinClient) GetBlockCount(ctx context.Context) (int64, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockCount")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHash provides a mock function with given fields: ctx, blockHeight
func (_m *BitcoinClient) GetBlockHash(ctx context.Context, blockHeight int64) (*chainhash.Hash, error) {
	ret := _m.Called(ctx, blockHeight)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockHash")
	}

	var r0 *chainhash.Hash
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*chainhash.Hash, error)); ok {
		return rf(ctx, blockHeight)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *chainhash.Hash); ok {
		r0 = rf(ctx, blockHeight)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chainhash.Hash)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, blockHeight)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHeader provides a mock function with given fields: ctx, hash
func (_m *BitcoinClient) GetBlockHeader(ctx context.Context, hash *chainhash.Hash) (*wire.BlockHeader, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockHeader")
	}

	var r0 *wire.BlockHeader
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash) (*wire.BlockHeader, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash) *wire.BlockHeader); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*wire.BlockHeader)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *chainhash.Hash) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHeightByStr provides a mock function with given fields: ctx, blockHash
func (_m *BitcoinClient) GetBlockHeightByStr(ctx context.Context, blockHash string) (int64, error) {
	ret := _m.Called(ctx, blockHash)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockHeightByStr")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int64, error)); ok {
		return rf(ctx, blockHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int64); ok {
		r0 = rf(ctx, blockHash)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, blockHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockVerbose provides a mock function with given fields: ctx, hash
func (_m *BitcoinClient) GetBlockVerbose(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockVerbose")
	}

	var r0 *btcjson.GetBlockVerboseTxResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash) *btcjson.GetBlockVerboseTxResult); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcjson.GetBlockVerboseTxResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *chainhash.Hash) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockVerboseByStr provides a mock function with given fields: ctx, blockHash
func (_m *BitcoinClient) GetBlockVerboseByStr(ctx context.Context, blockHash string) (*btcjson.GetBlockVerboseTxResult, error) {
	ret := _m.Called(ctx, blockHash)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockVerboseByStr")
	}

	var r0 *btcjson.GetBlockVerboseTxResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*btcjson.GetBlockVerboseTxResult, error)); ok {
		return rf(ctx, blockHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *btcjson.GetBlockVerboseTxResult); ok {
		r0 = rf(ctx, blockHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcjson.GetBlockVerboseTxResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, blockHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEstimatedFeeRate provides a mock function with given fields: ctx, confTarget
func (_m *BitcoinClient) GetEstimatedFeeRate(ctx context.Context, confTarget int64) (uint64, error) {
	ret := _m.Called(ctx, confTarget)

	if len(ret) == 0 {
		panic("no return value specified for GetEstimatedFeeRate")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (uint64, error)); ok {
		return rf(ctx, confTarget)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) uint64); ok {
		r0 = rf(ctx, confTarget)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, confTarget)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMempoolEntry provides a mock function with given fields: ctx, txHash
func (_m *BitcoinClient) GetMempoolEntry(ctx context.Context, txHash string) (*btcjson.GetMempoolEntryResult, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for GetMempoolEntry")
	}

	var r0 *btcjson.GetMempoolEntryResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*btcjson.GetMempoolEntryResult, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *btcjson.GetMempoolEntryResult); ok {
		r0 = rf(ctx, txHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcjson.GetMempoolEntryResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMempoolTxsAndFees provides a mock function with given fields: ctx, childHash
func (_m *BitcoinClient) GetMempoolTxsAndFees(ctx context.Context, childHash string) (client.MempoolTxsAndFees, error) {
	ret := _m.Called(ctx, childHash)

	if len(ret) == 0 {
		panic("no return value specified for GetMempoolTxsAndFees")
	}

	var r0 client.MempoolTxsAndFees
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (client.MempoolTxsAndFees, error)); ok {
		return rf(ctx, childHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) client.MempoolTxsAndFees); ok {
		r0 = rf(ctx, childHash)
	} else {
		r0 = ret.Get(0).(client.MempoolTxsAndFees)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, childHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNetworkInfo provides a mock function with given fields: ctx
func (_m *BitcoinClient) GetNetworkInfo(ctx context.Context) (*btcjson.GetNetworkInfoResult, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetNetworkInfo")
	}

	var r0 *btcjson.GetNetworkInfoResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*btcjson.GetNetworkInfoResult, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *btcjson.GetNetworkInfoResult); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcjson.GetNetworkInfoResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNewAddress provides a mock function with given fields: ctx, account
func (_m *BitcoinClient) GetNewAddress(ctx context.Context, account string) (btcutil.Address, error) {
	ret := _m.Called(ctx, account)

	if len(ret) == 0 {
		panic("no return value specified for GetNewAddress")
	}

	var r0 btcutil.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (btcutil.Address, error)); ok {
		return rf(ctx, account)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) btcutil.Address); ok {
		r0 = rf(ctx, account)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(btcutil.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRawMempool provides a mock function with given fields: ctx
func (_m *BitcoinClient) GetRawMempool(ctx context.Context) ([]*chainhash.Hash, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetRawMempool")
	}

	var r0 []*chainhash.Hash
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*chainhash.Hash, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*chainhash.Hash); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*chainhash.Hash)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRawTransaction provides a mock function with given fields: ctx, hash
func (_m *BitcoinClient) GetRawTransaction(ctx context.Context, hash *chainhash.Hash) (*btcutil.Tx, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetRawTransaction")
	}

	var r0 *btcutil.Tx
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash) (*btcutil.Tx, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash) *btcutil.Tx); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcutil.Tx)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *chainhash.Hash) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRawTransactionByStr provides a mock function with given fields: ctx, hash
func (_m *BitcoinClient) GetRawTransactionByStr(ctx context.Context, hash string) (*btcutil.Tx, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetRawTransactionByStr")
	}

	var r0 *btcutil.Tx
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*btcutil.Tx, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *btcutil.Tx); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcutil.Tx)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRawTransactionResult provides a mock function with given fields: ctx, hash, res
func (_m *BitcoinClient) GetRawTransactionResult(ctx context.Context, hash *chainhash.Hash, res *btcjson.GetTransactionResult) (btcjson.TxRawResult, error) {
	ret := _m.Called(ctx, hash, res)

	if len(ret) == 0 {
		panic("no return value specified for GetRawTransactionResult")
	}

	var r0 btcjson.TxRawResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash, *btcjson.GetTransactionResult) (btcjson.TxRawResult, error)); ok {
		return rf(ctx, hash, res)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash, *btcjson.GetTransactionResult) btcjson.TxRawResult); ok {
		r0 = rf(ctx, hash, res)
	} else {
		r0 = ret.Get(0).(btcjson.TxRawResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *chainhash.Hash, *btcjson.GetTransactionResult) error); ok {
		r1 = rf(ctx, hash, res)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRawTransactionVerbose provides a mock function with given fields: ctx, hash
func (_m *BitcoinClient) GetRawTransactionVerbose(ctx context.Context, hash *chainhash.Hash) (*btcjson.TxRawResult, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetRawTransactionVerbose")
	}

	var r0 *btcjson.TxRawResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash) (*btcjson.TxRawResult, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash) *btcjson.TxRawResult); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcjson.TxRawResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *chainhash.Hash) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransaction provides a mock function with given fields: ctx, hash
func (_m *BitcoinClient) GetTransaction(ctx context.Context, hash *chainhash.Hash) (*btcjson.GetTransactionResult, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetTransaction")
	}

	var r0 *btcjson.GetTransactionResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash) (*btcjson.GetTransactionResult, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *chainhash.Hash) *btcjson.GetTransactionResult); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*btcjson.GetTransactionResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *chainhash.Hash) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByStr provides a mock function with given fields: ctx, hash
func (_m *BitcoinClient) GetTransactionByStr(ctx context.Context, hash string) (*chainhash.Hash, *btcjson.GetTransactionResult, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionByStr")
	}

	var r0 *chainhash.Hash
	var r1 *btcjson.GetTransactionResult
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*chainhash.Hash, *btcjson.GetTransactionResult, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *chainhash.Hash); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chainhash.Hash)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *btcjson.GetTransactionResult); ok {
		r1 = rf(ctx, hash)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*btcjson.GetTransactionResult)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, hash)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetTransactionFeeAndRate provides a mock function with given fields: ctx, tx
func (_m *BitcoinClient) GetTransactionFeeAndRate(ctx context.Context, tx *btcjson.TxRawResult) (int64, int64, error) {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionFeeAndRate")
	}

	var r0 int64
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *btcjson.TxRawResult) (int64, int64, error)); ok {
		return rf(ctx, tx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *btcjson.TxRawResult) int64); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *btcjson.TxRawResult) int64); ok {
		r1 = rf(ctx, tx)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, *btcjson.TxRawResult) error); ok {
		r2 = rf(ctx, tx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Healthcheck provides a mock function with given fields: ctx
func (_m *BitcoinClient) Healthcheck(ctx context.Context) (time.Time, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Healthcheck")
	}

	var r0 time.Time
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (time.Time, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) time.Time); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ImportAddress provides a mock function with given fields: ctx, address
func (_m *BitcoinClient) ImportAddress(ctx context.Context, address string) error {
	ret := _m.Called(ctx, address)

	if len(ret) == 0 {
		panic("no return value specified for ImportAddress")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, address)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsTxStuckInMempool provides a mock function with given fields: ctx, txHash, maxWaitBlocks
func (_m *BitcoinClient) IsTxStuckInMempool(ctx context.Context, txHash string, maxWaitBlocks int64) (bool, time.Duration, error) {
	ret := _m.Called(ctx, txHash, maxWaitBlocks)

	if len(ret) == 0 {
		panic("no return value specified for IsTxStuckInMempool")
	}

	var r0 bool
	var r1 time.Duration
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) (bool, time.Duration, error)); ok {
		return rf(ctx, txHash, maxWaitBlocks)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) bool); ok {
		r0 = rf(ctx, txHash, maxWaitBlocks)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64) time.Duration); ok {
		r1 = rf(ctx, txHash, maxWaitBlocks)
	} else {
		r1 = ret.Get(1).(time.Duration)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, int64) error); ok {
		r2 = rf(ctx, txHash, maxWaitBlocks)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ListUnspent provides a mock function with given fields: ctx
func (_m *BitcoinClient) ListUnspent(ctx context.Context) ([]btcjson.ListUnspentResult, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListUnspent")
	}

	var r0 []btcjson.ListUnspentResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]btcjson.ListUnspentResult, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []btcjson.ListUnspentResult); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]btcjson.ListUnspentResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUnspentMinMaxAddresses provides a mock function with given fields: ctx, minConf, maxConf, addresses
func (_m *BitcoinClient) ListUnspentMinMaxAddresses(ctx context.Context, minConf int, maxConf int, addresses []btcutil.Address) ([]btcjson.ListUnspentResult, error) {
	ret := _m.Called(ctx, minConf, maxConf, addresses)

	if len(ret) == 0 {
		panic("no return value specified for ListUnspentMinMaxAddresses")
	}

	var r0 []btcjson.ListUnspentResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, []btcutil.Address) ([]btcjson.ListUnspentResult, error)); ok {
		return rf(ctx, minConf, maxConf, addresses)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, []btcutil.Address) []btcjson.ListUnspentResult); ok {
		r0 = rf(ctx, minConf, maxConf, addresses)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]btcjson.ListUnspentResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, []btcutil.Address) error); ok {
		r1 = rf(ctx, minConf, maxConf, addresses)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ping provides a mock function with given fields: ctx
func (_m *BitcoinClient) Ping(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Ping")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RawRequest provides a mock function with given fields: ctx, method, params
func (_m *BitcoinClient) RawRequest(ctx context.Context, method string, params []json.RawMessage) (json.RawMessage, error) {
	ret := _m.Called(ctx, method, params)

	if len(ret) == 0 {
		panic("no return value specified for RawRequest")
	}

	var r0 json.RawMessage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []json.RawMessage) (json.RawMessage, error)); ok {
		return rf(ctx, method, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []json.RawMessage) json.RawMessage); ok {
		r0 = rf(ctx, method, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(json.RawMessage)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []json.RawMessage) error); ok {
		r1 = rf(ctx, method, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendRawTransaction provides a mock function with given fields: ctx, tx, allowHighFees
func (_m *BitcoinClient) SendRawTransaction(ctx context.Context, tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error) {
	ret := _m.Called(ctx, tx, allowHighFees)

	if len(ret) == 0 {
		panic("no return value specified for SendRawTransaction")
	}

	var r0 *chainhash.Hash
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *wire.MsgTx, bool) (*chainhash.Hash, error)); ok {
		return rf(ctx, tx, allowHighFees)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *wire.MsgTx, bool) *chainhash.Hash); ok {
		r0 = rf(ctx, tx, allowHighFees)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*chainhash.Hash)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *wire.MsgTx, bool) error); ok {
		r1 = rf(ctx, tx, allowHighFees)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBitcoinClient creates a new instance of BitcoinClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBitcoinClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *BitcoinClient {
	mock := &BitcoinClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
