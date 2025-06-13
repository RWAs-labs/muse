package staking

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/RWAs-labs/muse/precompiles/logs"
)

func (c *Contract) addStakeLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	staker common.Address,
	validator string,
	amount *big.Int,
) error {
	event := c.Abi().Events[StakeEventName]

	valAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return err
	}

	// staker and validator are indexed event params
	topics, err := logs.MakeTopics(event, []interface{}{staker}, []interface{}{common.BytesToAddress(valAddr.Bytes())})
	if err != nil {
		return err
	}

	// amount is part of event data
	data, err := logs.PackArguments([]logs.Argument{
		{Type: "uint256", Value: amount},
	})
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}

func (c *Contract) addUnstakeLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	staker common.Address,
	validator string,
	amount *big.Int,
) error {
	event := c.Abi().Events[UnstakeEventName]
	valAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return err
	}

	// staker and validator are indexed event params
	topics, err := logs.MakeTopics(event, []interface{}{staker}, []interface{}{common.BytesToAddress(valAddr.Bytes())})
	if err != nil {
		return err
	}

	// amount is part of event data
	data, err := logs.PackArguments([]logs.Argument{
		{Type: "uint256", Value: amount},
	})
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}

func (c *Contract) addMoveStakeLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	staker common.Address,
	validatorSrc string,
	validatorDst string,
	amount *big.Int,
) error {
	event := c.Abi().Events[MoveStakeEventName]
	validatorSrcAddr, err := sdk.ValAddressFromBech32(validatorSrc)
	if err != nil {
		return err
	}

	validatorDstAddr, err := sdk.ValAddressFromBech32(validatorDst)
	if err != nil {
		return err
	}

	// staker and validators are indexed event params
	topics, err := logs.MakeTopics(
		event,
		[]interface{}{staker},
		[]interface{}{common.BytesToAddress(validatorSrcAddr.Bytes())},
		[]interface{}{common.BytesToAddress(validatorDstAddr.Bytes())},
	)
	if err != nil {
		return err
	}

	// amount is part of event data
	data, err := logs.PackArguments([]logs.Argument{
		{Type: "uint256", Value: amount},
	})
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}

func (c *Contract) addDistributeLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	distributor common.Address,
	mrc20Token common.Address,
	amount *big.Int,
) error {
	event := c.Abi().Events[DistributeEventName]

	topics, err := logs.MakeTopics(
		event,
		[]interface{}{distributor},
		[]interface{}{mrc20Token},
	)
	if err != nil {
		return err
	}

	data, err := logs.PackArguments([]logs.Argument{
		{Type: "uint256", Value: amount},
	})
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}

func (c *Contract) addClaimRewardsLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	delegator common.Address,
	mrc20Token common.Address,
	validator sdk.ValAddress,
	amount *big.Int,
) error {
	event := c.Abi().Events[ClaimRewardsEventName]

	topics, err := logs.MakeTopics(
		event,
		[]interface{}{delegator},
		[]interface{}{mrc20Token},
		[]interface{}{common.BytesToAddress(validator.Bytes())},
	)
	if err != nil {
		return err
	}

	data, err := logs.PackArguments([]logs.Argument{
		{Type: "uint256", Value: amount},
	})
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}
