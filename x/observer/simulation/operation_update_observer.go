package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	musesimulation "github.com/RWAs-labs/muse/testutil/simulation"
	"github.com/RWAs-labs/muse/x/observer/keeper"
	"github.com/RWAs-labs/muse/x/observer/types"
)

// SimulateUpdateObserver generates a TypeMsgUpdateObserver and delivers it.
func SimulateUpdateObserver(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := musesimulation.GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgUpdateObserver, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		_, randomObserver, observerList, err := musesimulation.GetRandomAccountAndObserver(r, ctx, k, accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgUpdateObserver, err.Error()), nil, nil
		}

		observerMap := make(map[string]bool)
		for _, observer := range observerList {
			observerMap[observer] = true
		}

		validators, err := k.GetStakingKeeper().GetAllValidators(ctx)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgUpdateObserver, err.Error()), nil, nil
		}

		if len(validators) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgUpdateObserver,
				"no validators found",
			), nil, nil
		}

		newObserver := ""
		foundNewObserver := musesimulation.RepeatCheck(func() bool {
			randomValidator := validators[r.Intn(len(validators))]
			randomValidatorAddress, err := types.GetAccAddressFromOperatorAddress(randomValidator.OperatorAddress)
			if err != nil {
				return false
			}
			newObserver = randomValidatorAddress.String()
			err = k.IsValidator(ctx, newObserver)
			if err != nil {
				return false
			}
			if _, ok := observerMap[newObserver]; !ok {
				return true
			}
			return false
		})

		if !foundNewObserver {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgUpdateObserver,
				"no new observer found",
			), nil, nil
		}

		lastBlockCount, found := k.GetLastObserverCount(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgUpdateObserver,
				"no last block count found",
			), nil, nil
		}
		// #nosec G115 - overflow is not a concern here
		if int(lastBlockCount.Count) != len(observerList) {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgUpdateObserver,
				"observer count mismatch",
			), nil, nil
		}

		msg := types.MsgUpdateObserver{
			Creator:            policyAccount.Address.String(),
			OldObserverAddress: randomObserver,
			NewObserverAddress: newObserver,
			UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgUpdateObserver, err.Error()), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           moduletestutil.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			Context:         ctx,
			SimAccount:      policyAccount,
			AccountKeeper:   k.GetAuthKeeper(),
			Bankkeeper:      k.GetBankKeeper(),
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return musesimulation.GenAndDeliverTxWithRandFees(txCtx, true)
	}
}
