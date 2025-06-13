package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	musesimulation "github.com/RWAs-labs/muse/testutil/simulation"
	"github.com/RWAs-labs/muse/x/crosschain/keeper"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// SimulateMsgRemoveOutboundTracker generates a MsgRemoveOutboundTracker with random values
func SimulateMsgRemoveOutboundTracker(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		policyAccount, err := musesimulation.GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgRemoveOutboundTracker, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		trackers := k.GetAllOutboundTracker(ctx)

		if len(trackers) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgRemoveOutboundTracker,
				"no outbound trackers found",
			), nil, nil
		}

		randomTracker := trackers[r.Intn(len(trackers))]

		msg := types.MsgRemoveOutboundTracker{
			ChainId: randomTracker.ChainId,
			Nonce:   randomTracker.Nonce,
			Creator: policyAccount.Address.String(),
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgRemoveOutboundTracker,
				"unable to validate MsgRemoveOutboundTracker",
			), nil, err
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
