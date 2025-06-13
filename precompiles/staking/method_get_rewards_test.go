package staking

import (
	"math/big"
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	precompiletypes "github.com/RWAs-labs/muse/precompiles/types"
	"github.com/RWAs-labs/muse/testutil/sample"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"
)

func Test_GetRewards(t *testing.T) {
	t.Run("should return empty rewards list to a non staker", func(t *testing.T) {
		/* ARRANGE */
		s := newTestSuite(t)

		// Create validator.
		validator := sample.Validator(t, rand.New(rand.NewSource(42)))
		s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)

		// Create staker.
		stakerEVMAddr := sample.EthAddress()

		/* ACT */
		// Call getRewards.
		getRewardsMethod := s.stkContractABI.Methods[GetRewardsMethodName]

		s.mockVMContract.Input = packInputArgs(
			t,
			getRewardsMethod,
			[]interface{}{stakerEVMAddr, validator.GetOperator()}...,
		)

		/* ASSERT */
		bytes, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)
		res, err := getRewardsMethod.Outputs.Unpack(bytes)
		require.NoError(t, err)
		require.Empty(t, res[0])
	})

	t.Run("should return the mrc20 rewards list for a staker", func(t *testing.T) {
		/* ARRANGE */
		s := newTestSuite(t)
		s.sdkKeepers.DistributionKeeper.FeePool.Set(s.ctx, distrtypes.InitialFeePool())

		// Create validator.
		validator := sample.Validator(t, rand.New(rand.NewSource(42)))
		s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)

		// Create staker.
		stakerEVMAddr := sample.EthAddress()
		stakerCosmosAddr, err := precompiletypes.GetCosmosAddress(s.sdkKeepers.BankKeeper, stakerEVMAddr)
		require.NoError(t, err)

		// Become a staker.
		stakeThroughCosmosAPI(
			t,
			s.ctx,
			s.sdkKeepers.BankKeeper,
			s.sdkKeepers.StakingKeeper,
			validator,
			stakerCosmosAddr,
			math.NewInt(100),
		)

		valAddr, err := sdk.ValAddressFromBech32(validator.GetOperator())
		require.NoError(t, err)
		err = s.sdkKeepers.DistributionKeeper.Hooks().
			AfterDelegationModified(s.ctx, stakerCosmosAddr, valAddr)
		require.NoError(t, err)

		/* Distribute 1000 MRC20 tokens to the staking contract */
		distributeMRC20(t, s, big.NewInt(1000))

		// TODO: Simulate a distribution of rewards.
		// emissions.BeginBlocker(s.ctx, *s.sdkKeepers.EmissionsKeeper)
		// staking.BeginBlocker(s.ctx, &s.sdkKeepers.StakingKeeper)
		// distribution.BeginBlocker(s.ctx, abci.RequestBeginBlock{}, s.sdkKeepers.DistributionKeeper)
		// s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 1)

		/* ACT */
		// Call getRewards.
		getRewardsMethod := s.stkContractABI.Methods[GetRewardsMethodName]

		s.mockVMContract.Input = packInputArgs(
			t,
			getRewardsMethod,
			[]interface{}{stakerEVMAddr, validator.GetOperator()}...,
		)

		bytes, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		/* ASSERT */
		_, err = getRewardsMethod.Outputs.Unpack(bytes)
		require.NoError(t, err)
	})
}
