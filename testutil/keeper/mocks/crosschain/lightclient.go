// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	proofs "github.com/RWAs-labs/muse/pkg/proofs"

	types "github.com/cosmos/cosmos-sdk/types"
)

// CrosschainLightclientKeeper is an autogenerated mock type for the CrosschainLightclientKeeper type
type CrosschainLightclientKeeper struct {
	mock.Mock
}

// VerifyProof provides a mock function with given fields: ctx, proof, chainID, blockHash, txIndex
func (_m *CrosschainLightclientKeeper) VerifyProof(ctx types.Context, proof *proofs.Proof, chainID int64, blockHash string, txIndex int64) ([]byte, error) {
	ret := _m.Called(ctx, proof, chainID, blockHash, txIndex)

	if len(ret) == 0 {
		panic("no return value specified for VerifyProof")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, *proofs.Proof, int64, string, int64) ([]byte, error)); ok {
		return rf(ctx, proof, chainID, blockHash, txIndex)
	}
	if rf, ok := ret.Get(0).(func(types.Context, *proofs.Proof, int64, string, int64) []byte); ok {
		r0 = rf(ctx, proof, chainID, blockHash, txIndex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, *proofs.Proof, int64, string, int64) error); ok {
		r1 = rf(ctx, proof, chainID, blockHash, txIndex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCrosschainLightclientKeeper creates a new instance of CrosschainLightclientKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCrosschainLightclientKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *CrosschainLightclientKeeper {
	mock := &CrosschainLightclientKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
