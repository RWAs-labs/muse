package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/RWAs-labs/muse/pkg/chains"
	musesimulation "github.com/RWAs-labs/muse/testutil/simulation"
	"github.com/RWAs-labs/muse/x/crosschain/keeper"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// SimulateMsgAbortStuckCCTX generates a MsgAbortStuckCCTX with random values
func SimulateMsgAbortStuckCCTX(k keeper.Keeper) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accounts []simtypes.Account, _ string,
	) (OperationMsg simtypes.OperationMsg, futureOps []simtypes.FutureOperation, err error) {
		// Pick a ethereum chain to abort a stuck cctx
		chainID := chains.GoerliLocalnet.ChainId
		supportedChains := k.GetObserverKeeper().GetSupportedChains(ctx)
		if len(supportedChains) == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAbortStuckCCTX,
				"no supported chains found",
			), nil, nil
		}

		for _, chain := range supportedChains {
			if chains.IsEthereumChain(chain.ChainId, []chains.Chain{}) {
				chainID = chain.ChainId
			}
		}

		policyAccount, err := musesimulation.GetPolicyAccount(ctx, k.GetAuthorityKeeper(), accounts)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgAbortStuckCCTX, err.Error()), nil, nil
		}

		authAccount := k.GetAuthKeeper().GetAccount(ctx, policyAccount.Address)
		spendable := k.GetBankKeeper().SpendableCoins(ctx, authAccount.GetAddress())

		tss, found := k.GetObserverKeeper().GetTSS(ctx)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAbortStuckCCTX,
				"no TSS found",
			), nil, nil
		}

		pendingNonces, found := k.GetObserverKeeper().GetPendingNonces(ctx, tss.TssPubkey, chainID)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAbortStuckCCTX,
				"no pending nonces found",
			), nil, nil
		}

		// If nonce low is the same as nonce high, it means that there are no pending nonces to add trackers for
		if pendingNonces.NonceLow == pendingNonces.NonceHigh {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAbortStuckCCTX,
				"no pending nonces found",
			), nil, nil
		}
		// Pick a random pending nonce
		nonce := 0
		switch {
		case pendingNonces.NonceHigh <= 1:
			nonce = int(pendingNonces.NonceLow)
		case pendingNonces.NonceLow == 0:
			nonce = r.Intn(int(pendingNonces.NonceHigh))
		default:
			nonce = r.Intn(int(pendingNonces.NonceHigh)-int(pendingNonces.NonceLow)) + int(pendingNonces.NonceLow)
		}

		nonceToCctx, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tss.TssPubkey, chainID, int64(nonce))
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAbortStuckCCTX,
				"no cctx found",
			), nil, nil
		}

		cctx, found := k.GetCrossChainTx(ctx, nonceToCctx.CctxIndex)
		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAbortStuckCCTX,
				"no cctx found",
			), nil, nil
		}

		msg := types.MsgAbortStuckCCTX{
			Creator:   policyAccount.Address.String(),
			CctxIndex: cctx.Index,
		}

		err = msg.ValidateBasic()
		if err != nil {
			return simtypes.NoOpMsg(
				types.ModuleName,
				TypeMsgAbortStuckCCTX,
				"unable to validate MsgAbortStuckCCTX msg",
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
