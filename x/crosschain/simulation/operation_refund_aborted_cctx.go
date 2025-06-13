package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/RWAs-labs/muse/testutil/sample"
	musesimulation "github.com/RWAs-labs/muse/testutil/simulation"
	"github.com/RWAs-labs/muse/x/crosschain/keeper"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// SimulateMsgRefundAbortedCCTX generates a MsgRefundAbortedCCTX with random values
func SimulateMsgRefundAbortedCCTX(k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		// Fetch the account from the auth keeper which can then be used to fetch spendable coins}
		policyAccount, err := musesimulation.GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgRefundAbortedCCTX, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		cctxList := k.GetAllCrossChainTx(ctx)
		abortedCctx := types.CrossChainTx{}
		abortedCctxFound := false

		for _, cctx := range cctxList {
			if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
				if !cctx.InboundParams.CoinType.SupportsRefund() {
					continue
				}
				if cctx.CctxStatus.IsAbortRefunded {
					continue
				}
				abortedCctx = cctx
				abortedCctxFound = true
				break
			}
		}
		if !abortedCctxFound {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgRefundAbortedCCTX, "no aborted cctx found"), nil, nil
		}

		msg := types.MsgRefundAbortedCCTX{
			Creator:       policyAccount.Address.String(),
			CctxIndex:     abortedCctx.Index,
			RefundAddress: sample.EthAddressFromRand(r).String(),
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgRefundAbortedCCTX,
				"unable to validate MsgRefundAbortedCCTX msg",
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
