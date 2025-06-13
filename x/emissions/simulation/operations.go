package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/RWAs-labs/muse/x/emissions/keeper"
	"github.com/RWAs-labs/muse/x/emissions/types"
)

var TypeMsgWithdrawEmission = sdk.MsgTypeURL(&types.MsgWithdrawEmission{})

// Simulation operation weights constants
// Operation weights are used by the simulation program to simulate the weight of different operations.
// This decides what percentage of a certain type of operation is part of a block.
// Based on the weights assigned in the cosmos sdk modules , 100 seems to the max weight used , and therefore guarantees that at least one operation of that type is present in a block.
// Operation weights are used by the `SimulateFromSeed`
// function to pick a random operation based on the weights.The functions with higher weights are more likely to be picked.

// Therefore, this decides the percentage of a certain operation that is part of a block.

// Based on the weights assigned in the cosmos sdk modules,
// 100 seems to the max weight used,and we should use relative weights
// to signify the number of each operation in a block.

const (
	DefaultWeightMsgWithdrawEmissionType = 100

	OpWeightMsgWithdrawEmissionType = "op_weight_msg_withdraw_emission_type" // #nosec G101 not a hardcoded credential 	// #nosec G101 not a hardcoded credential
)

func WeightedOperations(
	appParams simtypes.AppParams, k keeper.Keeper) simulation.WeightedOperations {
	var (
		weightMsgWithdrawEmissionType int
	)

	appParams.GetOrGenerate(OpWeightMsgWithdrawEmissionType, &weightMsgWithdrawEmissionType, nil,
		func(_ *rand.Rand) {
			weightMsgWithdrawEmissionType = DefaultWeightMsgWithdrawEmissionType
		})

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgWithdrawEmissionType,
			SimulateMsgWithdrawEmissions(k),
		),
	}
}
