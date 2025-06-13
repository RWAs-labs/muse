package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/RWAs-labs/muse/pkg/crypto"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

// LockMRC20 locks MRC20 tokens in the specified address
// The caller must have approved the locker contract to spend the amount of MRC20 tokens.
// Warning: This function does not mint cosmos coins, if the depositor needs to be rewarded
// it has to be implemented by the caller of this function.
func (k Keeper) LockMRC20(
	ctx sdk.Context,
	mrc20Address, spender, owner, locker common.Address,
	amount *big.Int,
) error {
	// owner is the EOA owner of the MRC20 tokens.
	// spender is the EOA allowed to spend MRC20 on owner's behalf.
	// locker is the address that will lock the MRC20 tokens, i.e: bank precompile.
	if err := k.CheckMRC20Allowance(ctx, owner, spender, mrc20Address, amount); err != nil {
		return errors.Wrap(err, "failed allowance check")
	}

	// Check amount_to_be_locked <= total_erc20_balance - already_locked
	// Max amount of MRC20 tokens that exists in mEVM are the total supply.
	totalSupply, err := k.MRC20TotalSupply(ctx, mrc20Address)
	if err != nil {
		return errors.Wrap(err, "failed totalSupply check")
	}

	// The alreadyLocked amount is the amount of MRC20 tokens that have been locked by the locker.
	// TODO: Implement list of whitelisted locker addresses (https://github.com/RWAs-labs/muse/issues/2991)
	alreadyLocked, err := k.MRC20BalanceOf(ctx, mrc20Address, locker)
	if err != nil {
		return errors.Wrap(err, "failed getting the MRC20 already locked amount")
	}

	if !k.IsValidDepositAmount(totalSupply, alreadyLocked, amount) {
		return errors.Wrap(fungibletypes.ErrInvalidAmount, "amount to be locked is not valid")
	}

	// Initiate a transferFrom the owner to the locker. This will lock the MRC20 tokens.
	// locker has to initiate the transaction and have enough allowance from owner.
	transferred, err := k.MRC20TransferFrom(ctx, mrc20Address, spender, owner, locker, amount)
	if err != nil {
		return errors.Wrap(err, "failed executing transferFrom")
	}

	if !transferred {
		return fmt.Errorf("transferFrom returned false (no success)")
	}

	return nil
}

// UnlockMRC20 unlocks MRC20 tokens and sends them to the owner.
// Warning: Before unlocking MRC20 tokens, the caller must check if
// the owner has enough collateral (cosmos coins) to be exchanged (burnt) for the MRC20 tokens.
func (k Keeper) UnlockMRC20(
	ctx sdk.Context,
	mrc20Address, owner, locker common.Address,
	amount *big.Int,
) error {
	// Check if the account locking the MRC20 tokens has enough balance.
	if err := k.CheckMRC20Balance(ctx, mrc20Address, locker, amount); err != nil {
		return errors.Wrap(err, "failed balance check")
	}

	// transfer from the EOA locking the assets to the owner.
	transferred, err := k.MRC20Transfer(ctx, mrc20Address, locker, owner, amount)
	if err != nil {
		return errors.Wrap(err, "failed executing transfer")
	}

	if !transferred {
		return fmt.Errorf("transfer returned false (no success)")
	}

	return nil
}

// CheckMRC20Allowance checks if the allowance of MRC20 tokens,
// is equal or greater than the provided amount.
func (k Keeper) CheckMRC20Allowance(
	ctx sdk.Context,
	owner, spender, mrc20Address common.Address,
	amount *big.Int,
) error {
	if amount.Sign() <= 0 || amount == nil {
		return fungibletypes.ErrInvalidAmount
	}

	if crypto.IsEmptyAddress(owner) || crypto.IsEmptyAddress(spender) {
		return fungibletypes.ErrZeroAddress
	}

	if err := k.IsValidMRC20(ctx, mrc20Address); err != nil {
		return errors.Wrap(err, "MRC20 is not valid")
	}

	allowanceValue, err := k.MRC20Allowance(ctx, mrc20Address, owner, spender)
	if err != nil {
		return errors.Wrap(err, "failed while checking spender's allowance")
	}

	if allowanceValue.Cmp(amount) < 0 || allowanceValue.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("invalid allowance, got %s, wanted %s", allowanceValue.String(), amount.String())
	}

	return nil
}

// CheckMRC20Balance checks if the balance of MRC20 tokens,
// is equal or greater than the provided amount.
func (k Keeper) CheckMRC20Balance(
	ctx sdk.Context,
	mrc20Address, owner common.Address,
	amount *big.Int,
) error {
	if amount.Sign() <= 0 || amount == nil {
		return fungibletypes.ErrInvalidAmount
	}

	if err := k.IsValidMRC20(ctx, mrc20Address); err != nil {
		return errors.Wrap(err, "MRC20 is not valid")
	}

	if crypto.IsEmptyAddress(owner) {
		return fungibletypes.ErrZeroAddress
	}

	// Check the MRC20 balance of a given account.
	// function balanceOf(address account)
	balance, err := k.MRC20BalanceOf(ctx, mrc20Address, owner)
	if err != nil {
		return errors.Wrap(err, "failed getting owner's MRC20 balance")
	}

	if balance.Cmp(amount) < 0 {
		return fmt.Errorf("invalid balance, got %s, wanted %s", balance.String(), amount.String())
	}

	return nil
}

// IsValidMRC20 returns an error whenever a MRC20 is not whitelisted or paused.
func (k Keeper) IsValidMRC20(ctx sdk.Context, mrc20Address common.Address) error {
	if crypto.IsEmptyAddress(mrc20Address) {
		return fungibletypes.ErrMRC20ZeroAddress
	}

	t, found := k.GetForeignCoins(ctx, mrc20Address.String())
	if !found {
		return fungibletypes.ErrMRC20NotWhiteListed
	}

	if t.Paused {
		return fungibletypes.ErrPausedMRC20
	}

	return nil
}

// IsValidDepositAmount checks "totalSupply >= amount_to_be_locked + amount_already_locked".
// A failure here means the user is trying to lock more than the available MRC20 supply.
// This suggests that an actor is minting MRC20 tokens out of thin air.
func (k Keeper) IsValidDepositAmount(totalSupply, alreadyLocked, amountToDeposit *big.Int) bool {
	if totalSupply == nil || alreadyLocked == nil || amountToDeposit == nil {
		return false
	}

	return totalSupply.Cmp(alreadyLocked.Add(alreadyLocked, amountToDeposit)) >= 0
}
