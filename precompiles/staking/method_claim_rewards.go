package staking

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/precompiles/bank"
	precompiletypes "github.com/RWAs-labs/muse/precompiles/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

// claimRewards claims all the rewards for a delegator from a validator.
// As F1 Cosmos distribution scheme implements an all or nothing withdrawal, the precompile will
// withdraw all the rewards for the delegator, filter MRC20 and unlock them to the delegator EVM address.
func (c *Contract) claimRewards(
	ctx sdk.Context,
	evm *vm.EVM,
	_ *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 2 {
		return nil, &precompiletypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		}
	}

	delegatorAddr, validatorAddr, err := unpackClaimRewardsArgs(args)
	if err != nil {
		return nil, err
	}

	// Get delegator Cosmos address.
	delegatorCosmosAddr, err := precompiletypes.GetCosmosAddress(c.bankKeeper, delegatorAddr)
	if err != nil {
		return nil, err
	}

	// Get validator Cosmos address.
	validatorCosmosAddr, err := sdk.ValAddressFromBech32(validatorAddr)
	if err != nil {
		return nil, err
	}

	// Withdraw all the delegation rewards.
	// The F1 Cosmos distribution scheme implements an all or nothing withdrawal.
	// The coins could be of multiple denomination, and a mix of MRC20 and Cosmos coins.
	coins, err := c.distributionKeeper.WithdrawDelegationRewards(ctx, delegatorCosmosAddr, validatorCosmosAddr)
	if err != nil {
		return nil, precompiletypes.ErrUnexpected{
			When: "WithdrawDelegationRewards",
			Got:  err.Error(),
		}
	}

	// For all the MRC20 coins withdrawed:
	// - Check the amount to unlock is valid.
	// - Burn the Cosmos coins.
	// - Unlock the MRC20 coins.
	for _, coin := range coins {
		// Filter out invalid coins.
		if !coin.IsValid() || !coin.Amount.IsPositive() || !precompiletypes.CoinIsMRC20(coin.Denom) {
			continue
		}

		// Notice that instead of returning errors we just skip the coin. This is because there might be
		// more than one MRC20 coin in the delegation rewards, and we want to unlock as many as possible.
		// Coins are locked in the bank precompile, so it should be possible to unlock them afterwards.
		var (
			mrc20Addr   = common.HexToAddress(strings.TrimPrefix(coin.Denom, config.MRC20DenomPrefix))
			mrc20Amount = coin.Amount.BigInt()
		)

		// Check if bank address has enough MRC20 balance.
		// This check is also made inside UnlockMRC20, but repeat it here to avoid burning the coins.
		if err := c.fungibleKeeper.CheckMRC20Balance(ctx, mrc20Addr, bank.ContractAddress, mrc20Amount); err != nil {
			ctx.Logger().Error(
				"Claimed invalid amount of MRC20 Validator Rewards",
				"Total", mrc20Amount,
				"Denom", precompiletypes.MRC20ToCosmosDenom(mrc20Addr),
			)

			continue
		}

		coinSet := sdk.NewCoins(coin)

		// Send the coins to the fungible module to burn them.
		if err := c.bankKeeper.SendCoinsFromAccountToModule(ctx, delegatorCosmosAddr, fungibletypes.ModuleName, coinSet); err != nil {
			continue
		}

		if err := c.bankKeeper.BurnCoins(ctx, fungibletypes.ModuleName, coinSet); err != nil {
			return nil, &precompiletypes.ErrUnexpected{
				When: "BurnCoins",
				Got:  err.Error(),
			}
		}

		// Finally, unlock the MRC20 coins.
		if err := c.fungibleKeeper.UnlockMRC20(ctx, mrc20Addr, delegatorAddr, bank.ContractAddress, mrc20Amount); err != nil {
			return nil, &precompiletypes.ErrUnexpected{
				When: "UnlockMRC20",
				Got:  err.Error(),
			}
		}

		// Emit an event per MRC20 coin unlocked.
		// This keeps events as granular and deterministic as possible.
		if err := c.addClaimRewardsLog(ctx, evm.StateDB, delegatorAddr, mrc20Addr, validatorCosmosAddr, mrc20Amount); err != nil {
			return nil, &precompiletypes.ErrUnexpected{
				When: "AddClaimRewardLog",
				Got:  err.Error(),
			}
		}

		ctx.Logger().Debug(
			"Claimed MRC20 rewards",
			"Delegator", delegatorCosmosAddr,
			"Denom", precompiletypes.MRC20ToCosmosDenom(mrc20Addr),
			"Amount", coin.Amount,
		)
	}

	return method.Outputs.Pack(true)
}

func unpackClaimRewardsArgs(args []interface{}) (delegator common.Address, validator string, err error) {
	delegator, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, "", &precompiletypes.ErrInvalidAddr{
			Got: delegator.String(),
		}
	}

	validator, ok = args[1].(string)
	if !ok {
		return common.Address{}, "", &precompiletypes.ErrInvalidAddr{
			Got: validator,
		}
	}

	return delegator, validator, nil
}
