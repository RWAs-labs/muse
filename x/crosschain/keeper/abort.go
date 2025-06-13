package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

// ProcessAbort processes the abort of a cctx
// It refunds the amount to the abort address and try calling onAbort
// StatusMessages contains current status messages for cctx (outbound and revert)
func (k Keeper) ProcessAbort(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	messages types.StatusMessages,
) {
	// only support cctx with v2 and with a defined abort address
	// also abort can't be processed if the abort amount is already refunded
	if cctx.ProtocolContractVersion != types.ProtocolContractVersion_V2 ||
		cctx.RevertOptions.AbortAddress == "" ||
		cctx.CctxStatus.IsAbortRefunded {
		messages.ErrorMessageAbort = "abort processing not supported for this cctx"

		cctx.CctxStatus.UpdateStatusAndErrorMessages(types.CctxStatus_Aborted, messages)
		return
	}

	abortedAmount := GetAbortedAmount(*cctx)
	abortAddress := ethcommon.HexToAddress(cctx.RevertOptions.AbortAddress)

	connectedChainID, outgoing, err := cctx.GetConnectedChainID()
	if err != nil {
		messages.ErrorMessageAbort = "failed to get connected chain ID: " + err.Error()
		cctx.CctxStatus.UpdateStatusAndErrorMessages(types.CctxStatus_Aborted, messages)
		return
	}

	// use a temporary context to not commit any state change if processing the abort logs fails
	// this is to avoid an inconsistent state where onAbort is called by created cctx inside are not processed
	tmpCtx, commit := ctx.CacheContext()

	// process the abort on the mevm
	evmTxResponse, err := k.fungibleKeeper.ProcessAbort(
		tmpCtx,
		cctx.InboundParams.Sender,
		abortedAmount.BigInt(),
		outgoing,
		connectedChainID,
		cctx.InboundParams.CoinType,
		cctx.InboundParams.Asset,
		abortAddress,
		cctx.RevertOptions.RevertMessage,
	)

	if evmTxResponse != nil && !fungibletypes.IsContractReverted(evmTxResponse, err) {
		logs := evmtypes.LogsToEthereum(evmTxResponse.Logs)
		if len(logs) > 0 {
			tmpCtx = tmpCtx.WithValue(InCCTXIndexKey, cctx.Index)
			txOrigin := cctx.InboundParams.TxOrigin
			if txOrigin == "" {
				txOrigin = cctx.GetInboundParams().Sender
			}

			// process logs to process cctx events initiated during the contract call
			processLogsErr := k.ProcessLogs(tmpCtx, logs, abortAddress, txOrigin)
			if processLogsErr != nil {
				// this happens if the cctx events are not processed correctly with invalid withdrawals
				// in this situation we want the CCTX to be reverted, we don't commit the state so the contract call is not persisted
				// the contract call is considered as reverted
				messages.ErrorMessageAbort = "failed to process logs for abort: " + err.Error()
				cctx.CctxStatus.UpdateStatusAndErrorMessages(types.CctxStatus_Aborted, messages)
				return
			}
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
					sdk.NewAttribute("action", "ProcessAbort"),
					sdk.NewAttribute("contract", cctx.RevertOptions.AbortAddress),
					sdk.NewAttribute("data", ""),
					sdk.NewAttribute("cctxIndex", cctx.Index),
				),
			)
		}
	}
	if err != nil {
		messages.ErrorMessageAbort = "failed to process abort: " + err.Error()
	}
	// note: we still set this value to true if onAbort reverted because the funds will still be deposited to the abortAddress
	if err == nil || errors.Is(err, fungibletypes.ErrOnAbortFailed) {
		cctx.CctxStatus.IsAbortRefunded = true
	}

	// commit state change from the deposit and eventual cctx events
	commit()

	cctx.CctxStatus.UpdateStatusAndErrorMessages(types.CctxStatus_Aborted, messages)

	return
}

// LegacyRefundAbortedAmountOnMuseChain refunds the amount of the cctx on MuseChain in case of aborted cctx
// For v2 cctx this logic has been replaced by using ProcessAbort of the fungible module
// TODO: Remove once only v2 workflow is supported
// https://github.com/RWAs-labs/muse/issues/2627
func (k Keeper) LegacyRefundAbortedAmountOnMuseChain(
	ctx sdk.Context,
	cctx types.CrossChainTx,
	refundAddress ethcommon.Address,
) error {
	coinType := cctx.InboundParams.CoinType
	switch coinType {
	case coin.CoinType_Gas:
		return k.LegacyRefundAbortedAmountOnMuseChainGas(ctx, cctx, refundAddress)
	case coin.CoinType_Muse:
		return k.LegacyRefundAbortedAmountOnMuseChainMuse(ctx, cctx, refundAddress)
	case coin.CoinType_ERC20:
		return k.LegacyRefundAbortedAmountOnMuseChainERC20(ctx, cctx, refundAddress)
	default:
		return fmt.Errorf("unsupported coin type for refund on MuseChain : %s", coinType)
	}
}

// LegacyRefundAbortedAmountOnMuseChainGas refunds the amount of the cctx on MuseChain in case of aborted cctx with cointype gas
// TODO: Remove once only v2 workflow is supported
// https://github.com/RWAs-labs/muse/issues/2627
func (k Keeper) LegacyRefundAbortedAmountOnMuseChainGas(
	ctx sdk.Context,
	cctx types.CrossChainTx,
	refundAddress ethcommon.Address,
) error {
	// refund in gas token to refund address
	// Refund the the amount was previously
	refundAmount := GetAbortedAmount(cctx)
	if refundAmount.IsNil() || refundAmount.IsZero() {
		return errors.New("no amount to refund")
	}
	chainID := cctx.InboundParams.SenderChainId
	// get the mrc20 contract address
	fcSenderChain, found := k.fungibleKeeper.GetGasCoinForForeignCoin(ctx, chainID)
	if !found {
		return types.ErrForeignCoinNotFound
	}
	mrc20 := ethcommon.HexToAddress(fcSenderChain.Mrc20ContractAddress)
	if mrc20 == (ethcommon.Address{}) {
		return errorsmod.Wrapf(types.ErrForeignCoinNotFound, "mrc20 contract address not found for chain %d", chainID)
	}
	// deposit the amount to the tx origin instead of receiver as this is a refund
	if _, err := k.fungibleKeeper.DepositMRC20(ctx, mrc20, refundAddress, refundAmount.BigInt()); err != nil {
		return errors.New("failed to refund muse on MuseChain" + err.Error())
	}
	return nil
}

// LegacyRefundAbortedAmountOnMuseChainMuse refunds the amount of the cctx on MuseChain in case of aborted cctx with cointype muse
// TODO: Remove once only v2 workflow is supported
// https://github.com/RWAs-labs/muse/issues/2627
func (k Keeper) LegacyRefundAbortedAmountOnMuseChainMuse(
	ctx sdk.Context,
	cctx types.CrossChainTx,
	refundAddress ethcommon.Address,
) error {
	// if coin type is Muse, handle this as a deposit MUSE to mEVM.
	refundAmount := GetAbortedAmount(cctx)
	chainID := cctx.InboundParams.SenderChainId
	// check if chain is an EVM chain
	if !chains.IsEVMChain(chainID, k.GetAuthorityKeeper().GetAdditionalChainList(ctx)) {
		return errors.New("only EVM chains are supported for refund when coin type is Muse")
	}
	if cctx.InboundParams.Amount.IsNil() || cctx.InboundParams.Amount.IsZero() {
		return errors.New("no amount to refund")
	}
	// deposit the amount to refund address
	if err := k.fungibleKeeper.DepositCoinMuse(ctx, refundAddress, refundAmount.BigInt()); err != nil {
		return fmt.Errorf("failed to refund muse on MuseChain: %w", err)
	}
	return nil
}

// LegacyRefundAbortedAmountOnMuseChainERC20 refunds the amount of the cctx on MuseChain in case of aborted cctx
// NOTE: GetCurrentOutboundParam should contain the last up to date cctx amount
// Refund address should already be validated before calling this function
// TODO: Remove once only v2 workflow is supported
// https://github.com/RWAs-labs/muse/issues/2627
func (k Keeper) LegacyRefundAbortedAmountOnMuseChainERC20(
	ctx sdk.Context,
	cctx types.CrossChainTx,
	refundAddress ethcommon.Address,
) error {
	refundAmount := GetAbortedAmount(cctx)
	// preliminary checks
	if cctx.InboundParams.CoinType != coin.CoinType_ERC20 {
		return errors.New("unsupported coin type for refund on MuseChain")
	}
	if !chains.IsEVMChain(cctx.InboundParams.SenderChainId, k.GetAuthorityKeeper().GetAdditionalChainList(ctx)) {
		return errors.New("only EVM chains are supported for refund on MuseChain")
	}

	if refundAmount.IsNil() || refundAmount.IsZero() {
		return errors.New("no amount to refund")
	}

	chainID, _, err := cctx.GetConnectedChainID()
	if err != nil {
		return errors.Wrap(err, "failed to get connected chain ID")
	}

	// get address of the mrc20
	fc, found := k.fungibleKeeper.GetForeignCoinFromAsset(
		ctx,
		cctx.InboundParams.Asset,
		chainID,
	)
	if !found {
		return fmt.Errorf("asset %s mrc not found", cctx.InboundParams.Asset)
	}
	mrc20 := ethcommon.HexToAddress(fc.Mrc20ContractAddress)
	if mrc20 == (ethcommon.Address{}) {
		return fmt.Errorf("asset %s invalid mrc address", cctx.InboundParams.Asset)
	}

	// deposit the amount to the sender
	if _, err := k.fungibleKeeper.DepositMRC20(ctx, mrc20, refundAddress, refundAmount.BigInt()); err != nil {
		return errors.New("failed to deposit mrc20 on MuseChain" + err.Error())
	}

	return nil
}
