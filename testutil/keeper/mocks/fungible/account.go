// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// FungibleAccountKeeper is an autogenerated mock type for the FungibleAccountKeeper type
type FungibleAccountKeeper struct {
	mock.Mock
}

// GetAccount provides a mock function with given fields: ctx, addr
func (_m *FungibleAccountKeeper) GetAccount(ctx context.Context, addr types.AccAddress) types.AccountI {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for GetAccount")
	}

	var r0 types.AccountI
	if rf, ok := ret.Get(0).(func(context.Context, types.AccAddress) types.AccountI); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.AccountI)
		}
	}

	return r0
}

// GetModuleAccount provides a mock function with given fields: ctx, name
func (_m *FungibleAccountKeeper) GetModuleAccount(ctx context.Context, name string) types.ModuleAccountI {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetModuleAccount")
	}

	var r0 types.ModuleAccountI
	if rf, ok := ret.Get(0).(func(context.Context, string) types.ModuleAccountI); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.ModuleAccountI)
		}
	}

	return r0
}

// GetSequence provides a mock function with given fields: ctx, addr
func (_m *FungibleAccountKeeper) GetSequence(ctx context.Context, addr types.AccAddress) (uint64, error) {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for GetSequence")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.AccAddress) (uint64, error)); ok {
		return rf(ctx, addr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.AccAddress) uint64); ok {
		r0 = rf(ctx, addr)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.AccAddress) error); ok {
		r1 = rf(ctx, addr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HasAccount provides a mock function with given fields: ctx, addr
func (_m *FungibleAccountKeeper) HasAccount(ctx context.Context, addr types.AccAddress) bool {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for HasAccount")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, types.AccAddress) bool); ok {
		r0 = rf(ctx, addr)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewAccountWithAddress provides a mock function with given fields: ctx, addr
func (_m *FungibleAccountKeeper) NewAccountWithAddress(ctx context.Context, addr types.AccAddress) types.AccountI {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for NewAccountWithAddress")
	}

	var r0 types.AccountI
	if rf, ok := ret.Get(0).(func(context.Context, types.AccAddress) types.AccountI); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.AccountI)
		}
	}

	return r0
}

// SetAccount provides a mock function with given fields: ctx, acc
func (_m *FungibleAccountKeeper) SetAccount(ctx context.Context, acc types.AccountI) {
	_m.Called(ctx, acc)
}

// NewFungibleAccountKeeper creates a new instance of FungibleAccountKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFungibleAccountKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *FungibleAccountKeeper {
	mock := &FungibleAccountKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
