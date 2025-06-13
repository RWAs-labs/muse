package e2etests

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/e2e/contracts/testdistribute"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/precompiles/bank"
	"github.com/RWAs-labs/muse/precompiles/staking"
	precompiletypes "github.com/RWAs-labs/muse/precompiles/types"
)

func TestPrecompilesDistributeAndClaimThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	var (
		// Addresses.
		staker               = r.EVMAddress()
		distrContractAddress = staking.ContractAddress
		lockerAddress        = bank.ContractAddress

		// Stake amount.
		stakeAmt = new(big.Int)

		// MRC20 distribution.
		mrc20Address  = r.ERC20MRC20Addr
		mrc20Denom    = precompiletypes.MRC20ToCosmosDenom(mrc20Address)
		mrc20DistrAmt = big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1e6))

		// carry is carried from the TestPrecompilesDistributeName test. It's applicable only to locker address.
		// This is needed because there's no easy way to retrieve that balance from the locker.
		carry              = big.NewInt(6210810988040846448)
		mrc20DistrAmtCarry = new(big.Int).Add(mrc20DistrAmt, carry)
		oneThousand        = big.NewInt(1e3)
		oneThousandOne     = big.NewInt(1001)
		fiveHundred        = big.NewInt(500)
		fiveHundredOne     = big.NewInt(501)
		zero               = big.NewInt(0)
		stake              = "1000000000000000000000"

		previousGasLimit = r.MEVMAuth.GasLimit
	)

	// stakeAmt has to be as big as the validator self delegation.
	// This way the rewards will be distributed 50%.
	_, ok := stakeAmt.SetString(stake, 10)
	require.True(r, ok)

	// Set new gas limit to avoid out of gas errors.
	r.MEVMAuth.GasLimit = 10_000_000

	distrContract, err := staking.NewIStaking(distrContractAddress, r.MEVMClient)
	require.NoError(r, err, "failed to create distribute contract caller")

	// testDstrContract  is the dApp contract that uses the staking precompile under the hood.
	_, tx, testDstrContract, err := testdistribute.DeployTestDistribute(r.MEVMAuth, r.MEVMClient)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "deployment of disitributor caller contract failed")

	// Save first validator bech32 address and ValAddress as it will be used through the test.
	validatorAddr, validatorValAddr := getValidatorAddresses(r, distrContract)

	// Reset the test after it finishes.
	defer resetDistributionTest(r, distrContract, lockerAddress, previousGasLimit, staker, validatorValAddr)

	// Get ERC20MRC20.
	txHash := r.LegacyDepositERC20WithAmountAndMessage(staker, mrc20DistrAmt, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	// There is no delegation, so the response should be empty.
	dv, err := testDstrContract.GetDelegatorValidatorsThroughContract(
		&bind.CallOpts{},
		staker,
	)
	require.NoError(r, err)
	require.Empty(r, dv, "DelegatorValidators response should be empty")

	// There should be no rewards.
	rewards, err := testDstrContract.GetRewardsThroughContract(&bind.CallOpts{}, staker, validatorAddr)
	require.NoError(r, err)
	require.Empty(r, rewards, "rewards should be empty when there are no delegations")

	// Stake with spender so it's registered as a delegator.
	err = stakeThroughCosmosAPI(r, validatorValAddr, staker, stakeAmt)
	require.NoError(r, err)

	// Check initial balances.
	balanceShouldBe(r, mrc20DistrAmt, checkMRC20Balance(r, staker))
	balanceShouldBe(r, carry, checkMRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, mrc20Denom))

	tx, err = testDstrContract.DistributeThroughContract(r.MEVMAuth, mrc20Address, oneThousand)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail when there's no allowance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, mrc20DistrAmt, checkMRC20Balance(r, staker))
	balanceShouldBe(r, carry, checkMRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, mrc20Denom))

	// Allow 500.
	approveAllowance(r, distrContractAddress, fiveHundred)

	tx, err = testDstrContract.DistributeThroughContract(r.MEVMAuth, mrc20Address, fiveHundredOne)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than allowed")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, mrc20DistrAmt, checkMRC20Balance(r, staker))
	balanceShouldBe(r, carry, checkMRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, mrc20Denom))

	// Raise the allowance to 1000.
	approveAllowance(r, distrContractAddress, oneThousand)

	// Shouldn't be able to distribute more than owned balance.
	tx, err = testDstrContract.DistributeThroughContract(r.MEVMAuth, mrc20Address, oneThousandOne)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than owned balance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, mrc20DistrAmt, checkMRC20Balance(r, staker))
	balanceShouldBe(r, carry, checkMRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, mrc20Denom))

	// Raise the allowance to max tokens.
	approveAllowance(r, distrContractAddress, mrc20DistrAmt)

	// Should be able to distribute an amount which is within balance and allowance.
	tx, err = testDstrContract.DistributeThroughContract(r.MEVMAuth, mrc20Address, mrc20DistrAmt)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "distribute should succeed when distributing within balance and allowance")

	balanceShouldBe(r, zero, checkMRC20Balance(r, staker))
	balanceShouldBe(r, mrc20DistrAmtCarry, checkMRC20Balance(r, lockerAddress))
	balanceShouldBe(r, mrc20DistrAmt, checkCosmosBalance(r, r.FeeCollectorAddress, mrc20Denom))

	eventDitributed, err := distrContract.ParseDistributed(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, mrc20Address, eventDitributed.Mrc20Token)
	require.Equal(r, staker, eventDitributed.Mrc20Distributor)
	require.Equal(r, mrc20DistrAmt.Uint64(), eventDitributed.Amount.Uint64())

	// After one block the rewards should have been distributed and fee collector should have 0 MRC20 balance.
	r.WaitForBlocks(1)
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, mrc20Denom))

	// DelegatorValidators returns the list of validator this delegator has delegated to.
	// The result should include the validator address.
	dv, err = testDstrContract.GetDelegatorValidatorsThroughContract(&bind.CallOpts{}, staker)
	require.NoError(r, err)
	require.Contains(r, dv, validatorAddr, "DelegatorValidators response should include validator address")

	// Get rewards and check it contains mrc20 tokens.
	rewards, err = testDstrContract.GetRewardsThroughContract(&bind.CallOpts{}, staker, validatorAddr)
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(rewards), 2)
	found := false
	for _, coin := range rewards {
		if strings.Contains(coin.Denom, config.MRC20DenomPrefix) {
			found = true
			break
		}
	}
	require.True(r, found, "rewards should include the MRC20 token")

	tx, err = testDstrContract.ClaimRewardsThroughContract(r.MEVMAuth, staker, validatorAddr)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "claim rewards should succeed")

	// Before claiming rewards the MRC20 balance is 0. After claiming rewards the MRC20 balance should be 14239697290875601808.
	// Which is the amount of MRC20 distributed, divided by two validators, and subtracted the commissions.
	mrc20RewardsAmt, ok := big.NewInt(0).SetString("14239697290875601808", 10)
	require.True(r, ok)
	balanceShouldBe(r, mrc20RewardsAmt, checkMRC20Balance(r, staker))

	eventClaimed, err := distrContract.ParseClaimedRewards(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, mrc20Address, eventClaimed.Mrc20Token)
	require.Equal(r, staker, eventClaimed.ClaimAddress)
	require.Equal(r, common.BytesToAddress(validatorValAddr.Bytes()), eventClaimed.Validator)
	require.Equal(r, mrc20RewardsAmt.Uint64(), eventClaimed.Amount.Uint64())

	// Locker final balance should be mrc20Distributed with carry - mrc20RewardsAmt.
	lockerFinalBalance := big.NewInt(0).Sub(mrc20DistrAmtCarry, mrc20RewardsAmt)
	balanceShouldBe(r, lockerFinalBalance, checkMRC20Balance(r, lockerAddress))

	// Staker final cosmos balance should be 0.
	balanceShouldBe(r, zero, checkCosmosBalance(r, sdk.AccAddress(staker.Bytes()), mrc20Denom))
}
