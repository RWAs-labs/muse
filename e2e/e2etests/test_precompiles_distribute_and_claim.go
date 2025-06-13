package e2etests

import (
	"math/big"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/precompiles/bank"
	"github.com/RWAs-labs/muse/precompiles/staking"
	precompiletypes "github.com/RWAs-labs/muse/precompiles/types"
)

func TestPrecompilesDistributeAndClaim(r *runner.E2ERunner, args []string) {
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

		// Amounts to test with.
		higherThanBalance = big.NewInt(0).Add(mrc20DistrAmt, big.NewInt(1))
		fiveHundred       = big.NewInt(500)
		fiveHundredOne    = big.NewInt(501)
		zero              = big.NewInt(0)
		stake             = "1000000000000000000000"

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

	// Retrieve the list of validators.
	validators, err := distrContract.GetAllValidators(&bind.CallOpts{})
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(validators), 2)

	// Save first validator bech32 address and as it will be used through the test.
	validatorAddr, validatorValAddr := getValidatorAddresses(r, distrContract)

	// Reset the test after it finishes.
	defer resetDistributionTest(r, distrContract, lockerAddress, previousGasLimit, staker, validatorValAddr)

	// Get ERC20MRC20.
	txHash := r.LegacyDepositERC20WithAmountAndMessage(staker, mrc20DistrAmt, []byte{})
	utils.WaitCctxMinedByInboundHash(r.Ctx, txHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

	// There is no delegation, so the response should be empty.
	dv, err := distrContract.GetDelegatorValidators(&bind.CallOpts{}, staker)
	require.NoError(r, err)
	require.Empty(r, dv, "DelegatorValidators response should be empty")

	// Shares at this point should be 0.
	sharesBefore, err := distrContract.GetShares(&bind.CallOpts{}, r.MEVMAuth.From, validatorAddr)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesBefore.Int64(), "shares should be 0 when there are no delegations")

	// There should be no rewards.
	rewards, err := distrContract.GetRewards(&bind.CallOpts{}, staker, validatorAddr)
	require.NoError(r, err)
	require.Empty(r, rewards, "rewards should be empty when there are no delegations")

	// Stake with spender so it's registered as a delegator.
	err = stakeThroughCosmosAPI(r, validatorValAddr, staker, stakeAmt)
	require.NoError(r, err)

	// Check initial balances.
	balanceShouldBe(r, mrc20DistrAmt, checkMRC20Balance(r, staker))
	balanceShouldBe(r, zero, checkMRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, mrc20Denom))

	// Failed attempt!
	tx, err := distrContract.Distribute(r.MEVMAuth, mrc20Address, mrc20DistrAmt)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail when there's no allowance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, mrc20DistrAmt, checkMRC20Balance(r, staker))
	balanceShouldBe(r, zero, checkMRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, mrc20Denom))

	// Allow 500.
	approveAllowance(r, distrContractAddress, fiveHundred)

	// Failed attempt! Shouldn't be able to distribute more than allowed.
	tx, err = distrContract.Distribute(r.MEVMAuth, mrc20Address, fiveHundredOne)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than allowed")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, mrc20DistrAmt, checkMRC20Balance(r, staker))
	balanceShouldBe(r, zero, checkMRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, mrc20Denom))

	// Raise the allowance to the maximum MRC20 amount.
	approveAllowance(r, distrContractAddress, mrc20DistrAmt)

	// Failed attempt! Shouldn't be able to distribute more than owned balance.
	tx, err = distrContract.Distribute(r.MEVMAuth, mrc20Address, higherThanBalance)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt, "distribute should fail trying to distribute more than owned balance")

	// Balances shouldn't change after a failed attempt.
	balanceShouldBe(r, mrc20DistrAmt, checkMRC20Balance(r, staker))
	balanceShouldBe(r, zero, checkMRC20Balance(r, lockerAddress))
	balanceShouldBe(r, zero, checkCosmosBalance(r, r.FeeCollectorAddress, mrc20Denom))

	// Should be able to distribute an amount which is within balance and allowance.
	tx, err = distrContract.Distribute(r.MEVMAuth, mrc20Address, mrc20DistrAmt)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "distribute should succeed when distributing within balance and allowance")

	balanceShouldBe(r, zero, checkMRC20Balance(r, staker))
	balanceShouldBe(r, mrc20DistrAmt, checkMRC20Balance(r, lockerAddress))
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
	dv, err = distrContract.GetDelegatorValidators(&bind.CallOpts{}, staker)
	require.NoError(r, err)
	require.Contains(r, dv, validatorAddr, "DelegatorValidators response should include validator address")

	// Get rewards and check it contains mrc20 tokens.
	rewards, err = distrContract.GetRewards(&bind.CallOpts{}, staker, validatorAddr)
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

	// Claim the rewards, they'll be unlocked as MRC20 tokens.
	tx, err = distrContract.ClaimRewards(r.MEVMAuth, staker, validatorAddr)
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

	// Locker final balance should be mrc20Disitributed - mrc20RewardsAmt.
	lockerFinalBalance := big.NewInt(0).Sub(mrc20DistrAmt, mrc20RewardsAmt)
	balanceShouldBe(r, lockerFinalBalance, checkMRC20Balance(r, lockerAddress))

	// Staker final cosmos balance should be 0.
	balanceShouldBe(r, zero, checkCosmosBalance(r, sdk.AccAddress(staker.Bytes()), mrc20Denom))
}

func TestPrecompilesDistributeNonMRC20(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	// Increase the gasLimit. It's required because of the gas consumed by precompiled functions.
	previousGasLimit := r.MEVMAuth.GasLimit
	r.MEVMAuth.GasLimit = 10_000_000
	defer func() {
		r.MEVMAuth.GasLimit = previousGasLimit
	}()

	spender, dstrAddress := r.EVMAddress(), staking.ContractAddress

	// Create a staking contract caller.
	dstrContract, err := staking.NewIStaking(dstrAddress, r.MEVMClient)
	require.NoError(r, err, "Failed to create staking contract caller")

	// Deposit and approve 50 WMUSE for the test.
	approveAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(50))
	r.LegacyDepositAndApproveWMuse(approveAmount)

	// Allow the staking contract to spend 25 WMuse tokens.
	tx, err := r.WMuse.Approve(r.MEVMAuth, dstrAddress, big.NewInt(25))
	require.NoError(r, err, "Error approving allowance for staking contract")
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.EqualValues(r, uint64(1), receipt.Status, "approve allowance tx failed")

	// Check the allowance of the staking in WMuse tokens. Should be 25.
	allowance, err := r.WMuse.Allowance(&bind.CallOpts{Context: r.Ctx}, spender, dstrAddress)
	require.NoError(r, err, "Error retrieving staking allowance")
	require.EqualValues(r, uint64(25), allowance.Uint64(), "Error allowance for staking contract")

	// Call Distribute with 25 Non MRC20 tokens. Should fail.
	tx, err = dstrContract.Distribute(r.MEVMAuth, r.WMuseAddr, big.NewInt(25))
	require.NoError(r, err, "Error calling staking.distribute()")
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	require.Equal(r, uint64(0), receipt.Status, "Non MRC20 deposit should fail")
}

// checkCosmosBalance checks the cosmos coin balance for an address. The coin is specified by its denom.
func checkCosmosBalance(r *runner.E2ERunner, address sdk.AccAddress, denom string) *big.Int {
	bal, err := r.BankClient.Balance(
		r.Ctx,
		&banktypes.QueryBalanceRequest{Address: address.String(), Denom: denom},
	)
	require.NoError(r, err)

	return bal.Balance.Amount.BigInt()
}

func stakeThroughCosmosAPI(
	r *runner.E2ERunner,
	validator sdk.ValAddress,
	staker common.Address,
	amount *big.Int,
) error {
	msg := stakingtypes.NewMsgDelegate(
		sdk.AccAddress(staker.Bytes()).String(),
		validator.String(),
		sdk.Coin{
			Denom:  config.BaseDenom,
			Amount: math.NewIntFromBigInt(amount),
		},
	)

	_, err := r.MuseTxServer.BroadcastTx(sdk.AccAddress(staker.Bytes()).String(), msg)
	if err != nil {
		return err
	}

	return nil
}

func resetDistributionTest(
	r *runner.E2ERunner,
	distrContract *staking.IStaking,
	lockerAddress common.Address,
	previousGasLimit uint64,
	staker common.Address,
	validator sdk.ValAddress,
) {
	validatorAddr, _ := getValidatorAddresses(r, distrContract)

	amount, err := distrContract.GetShares(&bind.CallOpts{}, r.MEVMAuth.From, validatorAddr)
	require.NoError(r, err)

	// Restore the gas limit.
	r.MEVMAuth.GasLimit = previousGasLimit

	// Reset the allowance to 0; this is needed when running upgrade tests where this test runs twice.
	tx, err := r.ERC20MRC20.Approve(r.MEVMAuth, lockerAddress, big.NewInt(0))
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "resetting allowance failed")

	// Reset balance to 0 for spender; this is needed when running upgrade tests where this test runs twice.
	balance, err := r.ERC20MRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// Burn all ERC20 balance.
	tx, err = r.ERC20MRC20.Transfer(
		r.MEVMAuth,
		common.HexToAddress("0x000000000000000000000000000000000000dEaD"),
		balance,
	)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "Resetting balance failed")

	// Clean the delegation.
	// Delegator will always delegate on the first validator.
	msg := stakingtypes.NewMsgUndelegate(
		sdk.AccAddress(staker.Bytes()).String(),
		validator.String(),
		sdk.Coin{
			Denom:  config.BaseDenom,
			Amount: math.NewIntFromBigInt(amount.Div(amount, big.NewInt(1e18))),
		},
	)

	_, err = r.MuseTxServer.BroadcastTx(sdk.AccAddress(staker.Bytes()).String(), msg)
	require.NoError(r, err)
}

func getValidatorAddresses(r *runner.E2ERunner, distrContract *staking.IStaking) (string, sdk.ValAddress) {
	// distrContract, err := staking.NewIStaking(staking.ContractAddress, r.MEVMClient)
	// require.NoError(r, err, "failed to create distribute contract caller")

	// Retrieve the list of validators.
	validators, err := distrContract.GetAllValidators(&bind.CallOpts{})
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(validators), 2)

	// Save first validators as it will be used through the test.
	validatorAddr, err := sdk.ValAddressFromBech32(validators[0].OperatorAddress)
	require.NoError(r, err)

	return validators[0].OperatorAddress, validatorAddr
}
