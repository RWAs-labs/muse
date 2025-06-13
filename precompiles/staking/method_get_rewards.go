package staking

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	precompiletypes "github.com/RWAs-labs/muse/precompiles/types"
)

// getRewards returns the list of MRC20 cosmos coins, available for withdrawal by the delegator.
func (c *Contract) getRewards(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 2 {
		return nil, &precompiletypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		}
	}

	delegatorAddr, validatorAddr, err := unpackGetRewardsArgs(args)
	if err != nil {
		return nil, err
	}

	// Get delegator Cosmos address.
	delegatorCosmosAddr, err := precompiletypes.GetCosmosAddress(c.bankKeeper, delegatorAddr)
	if err != nil {
		return nil, err
	}

	// Query the delegation rewards through the distribution keeper querier.
	dstrQuerier := distrkeeper.NewQuerier(c.distributionKeeper)

	res, err := dstrQuerier.DelegationRewards(ctx, &distrtypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delegatorCosmosAddr.String(),
		ValidatorAddress: validatorAddr,
	})

	// DelegationRewards returns an error if the delegation does not exist.
	// In this case, simply return an empty list of rewards, so external contracts
	// can process this case without failing.
	if err != nil {
		if errors.Is(err, stakingtypes.ErrNoDelegation) {
			rewards := make([]DecCoin, 0)
			return method.Outputs.Pack(rewards)
		}

		return nil, &precompiletypes.ErrUnexpected{
			When: "DelegationRewards",
			Got:  err.Error(),
		}
	}

	coins := res.GetRewards()
	if !coins.IsValid() {
		return nil, precompiletypes.ErrUnexpected{
			When: "GetRewards",
			Got:  "invalid coins",
		}
	}

	rewards := make([]DecCoin, 0)
	for _, coin := range coins {
		rewards = append(rewards, DecCoin{
			Denom:  coin.Denom,
			Amount: coin.Amount.BigInt(),
		})
	}

	return method.Outputs.Pack(rewards)
}

func unpackGetRewardsArgs(args []interface{}) (delegator common.Address, validator string, err error) {
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
