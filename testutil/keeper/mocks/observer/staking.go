// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	context "context"

	cosmos_sdktypes "github.com/cosmos/cosmos-sdk/types"
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ObserverStakingKeeper is an autogenerated mock type for the ObserverStakingKeeper type
type ObserverStakingKeeper struct {
	mock.Mock
}

// GetAllValidators provides a mock function with given fields: ctx
func (_m *ObserverStakingKeeper) GetAllValidators(ctx context.Context) ([]types.Validator, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllValidators")
	}

	var r0 []types.Validator
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]types.Validator, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []types.Validator); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Validator)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDelegation provides a mock function with given fields: ctx, delAddr, valAddr
func (_m *ObserverStakingKeeper) GetDelegation(ctx context.Context, delAddr cosmos_sdktypes.AccAddress, valAddr cosmos_sdktypes.ValAddress) (types.Delegation, error) {
	ret := _m.Called(ctx, delAddr, valAddr)

	if len(ret) == 0 {
		panic("no return value specified for GetDelegation")
	}

	var r0 types.Delegation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, cosmos_sdktypes.ValAddress) (types.Delegation, error)); ok {
		return rf(ctx, delAddr, valAddr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.AccAddress, cosmos_sdktypes.ValAddress) types.Delegation); ok {
		r0 = rf(ctx, delAddr, valAddr)
	} else {
		r0 = ret.Get(0).(types.Delegation)
	}

	if rf, ok := ret.Get(1).(func(context.Context, cosmos_sdktypes.AccAddress, cosmos_sdktypes.ValAddress) error); ok {
		r1 = rf(ctx, delAddr, valAddr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValidator provides a mock function with given fields: ctx, addr
func (_m *ObserverStakingKeeper) GetValidator(ctx context.Context, addr cosmos_sdktypes.ValAddress) (types.Validator, error) {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for GetValidator")
	}

	var r0 types.Validator
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.ValAddress) (types.Validator, error)); ok {
		return rf(ctx, addr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, cosmos_sdktypes.ValAddress) types.Validator); ok {
		r0 = rf(ctx, addr)
	} else {
		r0 = ret.Get(0).(types.Validator)
	}

	if rf, ok := ret.Get(1).(func(context.Context, cosmos_sdktypes.ValAddress) error); ok {
		r1 = rf(ctx, addr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetValidator provides a mock function with given fields: ctx, validator
func (_m *ObserverStakingKeeper) SetValidator(ctx context.Context, validator types.Validator) error {
	ret := _m.Called(ctx, validator)

	if len(ret) == 0 {
		panic("no return value specified for SetValidator")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.Validator) error); ok {
		r0 = rf(ctx, validator)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewObserverStakingKeeper creates a new instance of ObserverStakingKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewObserverStakingKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *ObserverStakingKeeper {
	mock := &ObserverStakingKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
