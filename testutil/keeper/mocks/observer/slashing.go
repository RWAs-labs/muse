// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// ObserverSlashingKeeper is an autogenerated mock type for the ObserverSlashingKeeper type
type ObserverSlashingKeeper struct {
	mock.Mock
}

// IsTombstoned provides a mock function with given fields: ctx, addr
func (_m *ObserverSlashingKeeper) IsTombstoned(ctx context.Context, addr types.ConsAddress) bool {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for IsTombstoned")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, types.ConsAddress) bool); ok {
		r0 = rf(ctx, addr)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SetValidatorSigningInfo provides a mock function with given fields: ctx, address, info
func (_m *ObserverSlashingKeeper) SetValidatorSigningInfo(ctx context.Context, address types.ConsAddress, info slashingtypes.ValidatorSigningInfo) error {
	ret := _m.Called(ctx, address, info)

	if len(ret) == 0 {
		panic("no return value specified for SetValidatorSigningInfo")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.ConsAddress, slashingtypes.ValidatorSigningInfo) error); ok {
		r0 = rf(ctx, address, info)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewObserverSlashingKeeper creates a new instance of ObserverSlashingKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewObserverSlashingKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *ObserverSlashingKeeper {
	mock := &ObserverSlashingKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
