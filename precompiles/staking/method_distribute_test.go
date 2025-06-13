package staking

import (
	"math/big"
	"testing"

	precompiletypes "github.com/RWAs-labs/muse/precompiles/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func Test_Distribute(t *testing.T) {
	feeCollectorAddress := authtypes.NewModuleAddress(authtypes.FeeCollectorName).String()

	t.Run("should fail to run distribute as read only method", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)
		distributeMethod := s.stkContractABI.Methods[DistributeMethodName]

		mrc20Denom := precompiletypes.MRC20ToCosmosDenom(s.mrc20Address)

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			distributeMethod,
			[]interface{}{s.mrc20Address, big.NewInt(0)}...,
		)

		// Call method as read only.
		result, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, true)

		// Check error and result.
		require.ErrorIs(t, err, precompiletypes.ErrWriteMethod{
			Method: DistributeMethodName,
		})

		// Result is empty as the write check is done before executing distribute() function.
		// On-chain this would look like reverting, so staticcall is properly reverted.
		require.Empty(t, result)

		// End fee collector balance should be 0.
		balance, err := s.sdkKeepers.BankKeeper.Balance(s.ctx, &banktypes.QueryBalanceRequest{
			Address: feeCollectorAddress,
			Denom:   mrc20Denom,
		})
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance.Balance.Amount.Uint64())
	})

	t.Run("should fail to distribute with 0 token balance", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)
		distributeMethod := s.stkContractABI.Methods[DistributeMethodName]
		mrc20Denom := precompiletypes.MRC20ToCosmosDenom(s.mrc20Address)

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			distributeMethod,
			[]interface{}{s.mrc20Address, big.NewInt(0)}...,
		)

		// Call method.
		success, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.ErrorAs(
			t,
			precompiletypes.ErrInvalidAmount{
				Got: "0",
			},
			err,
		)

		// Unpack and check result boolean.
		res, err := distributeMethod.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)

		// End fee collector balance should be 0.
		balance, err := s.sdkKeepers.BankKeeper.Balance(s.ctx, &banktypes.QueryBalanceRequest{
			Address: feeCollectorAddress,
			Denom:   mrc20Denom,
		})
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance.Balance.Amount.Uint64())
	})

	t.Run("should fail to distribute with 0 allowance", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)
		distributeMethod := s.stkContractABI.Methods[DistributeMethodName]
		mrc20Denom := precompiletypes.MRC20ToCosmosDenom(s.mrc20Address)

		// Set caller balance.
		_, err := s.fungibleKeeper.DepositMRC20(s.ctx, s.mrc20Address, s.defaultCaller, big.NewInt(1000))
		require.NoError(t, err)

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			distributeMethod,
			[]interface{}{s.mrc20Address, big.NewInt(1000)}...,
		)

		// Call method.
		success, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid allowance, got 0")

		// Unpack and check result boolean.
		res, err := distributeMethod.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)

		// End fee collector balance should be 0.
		balance, err := s.sdkKeepers.BankKeeper.Balance(s.ctx, &banktypes.QueryBalanceRequest{
			Address: feeCollectorAddress,
			Denom:   mrc20Denom,
		})
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance.Balance.Amount.Uint64())
	})

	t.Run("should fail to distribute 0 token", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)
		distributeMethod := s.stkContractABI.Methods[DistributeMethodName]
		mrc20Denom := precompiletypes.MRC20ToCosmosDenom(s.mrc20Address)

		// Set caller balance.
		_, err := s.fungibleKeeper.DepositMRC20(s.ctx, s.mrc20Address, s.defaultCaller, big.NewInt(1000))
		require.NoError(t, err)

		// Allow staking to spend MRC20 tokens.
		allowStaking(t, s, big.NewInt(1000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			distributeMethod,
			[]interface{}{s.mrc20Address, big.NewInt(0)}...,
		)

		// Call method.
		success, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid token amount: 0")

		// Unpack and check result boolean.
		res, err := distributeMethod.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)

		// End fee collector balance should be 0.
		balance, err := s.sdkKeepers.BankKeeper.Balance(s.ctx, &banktypes.QueryBalanceRequest{
			Address: feeCollectorAddress,
			Denom:   mrc20Denom,
		})
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance.Balance.Amount.Uint64())
	})

	t.Run("should fail to distribute more than allowed to staking", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)
		distributeMethod := s.stkContractABI.Methods[DistributeMethodName]
		mrc20Denom := precompiletypes.MRC20ToCosmosDenom(s.mrc20Address)

		// Set caller balance.
		_, err := s.fungibleKeeper.DepositMRC20(s.ctx, s.mrc20Address, s.defaultCaller, big.NewInt(1000))
		require.NoError(t, err)

		// Allow staking to spend MRC20 tokens.
		allowStaking(t, s, big.NewInt(999))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			distributeMethod,
			[]interface{}{s.mrc20Address, big.NewInt(1000)}...,
		)

		// Call method.
		success, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid allowance, got 999, wanted 1000")

		// Unpack and check result boolean.
		res, err := distributeMethod.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)

		// End fee collector balance should be 0.
		balance, err := s.sdkKeepers.BankKeeper.Balance(s.ctx, &banktypes.QueryBalanceRequest{
			Address: feeCollectorAddress,
			Denom:   mrc20Denom,
		})
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance.Balance.Amount.Uint64())
	})

	t.Run("should fail to distribute more than user balance", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)
		distributeMethod := s.stkContractABI.Methods[DistributeMethodName]
		mrc20Denom := precompiletypes.MRC20ToCosmosDenom(s.mrc20Address)

		// Set caller balance.
		_, err := s.fungibleKeeper.DepositMRC20(s.ctx, s.mrc20Address, s.defaultCaller, big.NewInt(1000))
		require.NoError(t, err)

		// Allow staking to spend MRC20 tokens.
		allowStaking(t, s, big.NewInt(100000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			distributeMethod,
			[]interface{}{s.mrc20Address, big.NewInt(1001)}...,
		)

		success, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.Error(t, err)
		require.Contains(t, err.Error(), "execution reverted")

		// Unpack and check result boolean.
		res, err := distributeMethod.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.False(t, ok)

		// End fee collector balance should be 0.
		balance, err := s.sdkKeepers.BankKeeper.Balance(s.ctx, &banktypes.QueryBalanceRequest{
			Address: feeCollectorAddress,
			Denom:   mrc20Denom,
		})
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance.Balance.Amount.Uint64())
	})

	t.Run("should distribute and lock MRC20", func(t *testing.T) {
		// Setup test.
		s := newTestSuite(t)
		distributeMethod := s.stkContractABI.Methods[DistributeMethodName]

		// Set caller balance.
		_, err := s.fungibleKeeper.DepositMRC20(s.ctx, s.mrc20Address, s.defaultCaller, big.NewInt(1000))
		require.NoError(t, err)

		// Allow staking to spend MRC20 tokens.
		allowStaking(t, s, big.NewInt(1000))

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			distributeMethod,
			[]interface{}{s.mrc20Address, big.NewInt(1000)}...,
		)

		success, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)

		// Check error.
		require.NoError(t, err)

		// Unpack and check result boolean.
		res, err := distributeMethod.Outputs.Unpack(success)
		require.NoError(t, err)

		ok := res[0].(bool)
		require.True(t, ok)
	})
}
