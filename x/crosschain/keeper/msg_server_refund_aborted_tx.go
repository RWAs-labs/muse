package keeper

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"golang.org/x/net/context"

	"github.com/RWAs-labs/muse/pkg/coin"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// RefundAbortedCCTX refunds the aborted CCTX.
// It verifies if the CCTX is aborted and not refunded, and if the refund address is valid.
// It refunds the amount to the refund address and sets the CCTX as refunded.
// Refer to documentation for GetRefundAddress for the refund address logic.
// Refer to documentation for GetAbortedAmount for the aborted amount logic.
func (k msgServer) RefundAbortedCCTX(
	goCtx context.Context,
	msg *types.MsgRefundAbortedCCTX,
) (*types.MsgRefundAbortedCCTXResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if authorized
	err := k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, errorsmod.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	// check if the cctx exists
	cctx, found := k.GetCrossChainTx(ctx, msg.CctxIndex)
	if !found {
		return nil, types.ErrCannotFindCctx
	}

	// check if the cctx is aborted
	if cctx.CctxStatus.Status != types.CctxStatus_Aborted {
		return nil, errorsmod.Wrap(types.ErrInvalidStatus, "CCTX is not aborted")
	}
	// check if the cctx is not refunded
	if cctx.CctxStatus.IsAbortRefunded {
		return nil, errorsmod.Wrap(types.ErrUnableProcessRefund, "CCTX is already refunded")
	}

	// Check if aborted amount is available to maintain muse accounting
	if cctx.InboundParams.CoinType == coin.CoinType_Muse {
		err := k.RemoveMuseAbortedAmount(ctx, GetAbortedAmount(cctx))
		// if the muse accounting is not found, it means the muse accounting is not set yet and the refund should not be processed
		if errors.Is(err, types.ErrUnableToFindMuseAccounting) {
			return nil, errorsmod.Wrap(types.ErrUnableProcessRefund, err.Error())
		}
		// if the muse accounting is found but the amount is insufficient, it means the refund can be processed but the muse accounting is not maintained properly
		if errors.Is(err, types.ErrInsufficientMuseAmount) {
			ctx.Logger().Error("Muse Accounting Error: ", err)
		}
	}

	refundAddress, err := GetRefundAddress(msg.RefundAddress)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidAddress, err.Error())
	}
	// refund the amount
	// use temporary context to avoid gas refunding issues and side effects
	tmpCtx, commit := ctx.CacheContext()
	err = k.LegacyRefundAbortedAmountOnMuseChain(tmpCtx, cctx, refundAddress)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUnableProcessRefund, err.Error())
	}
	commit()

	// set the cctx as refunded
	cctx.CctxStatus.SetAbortRefunded()

	k.SetCrossChainTx(ctx, cctx)
	return &types.MsgRefundAbortedCCTXResponse{}, nil
}

// GetRefundAddress gets the proper refund address.
// For BTC sender chain the refund address is the one provided in the message in the RefundAddress field.
// For EVM chain with coin type ERC20 the refund address is the sender , but can be overridden by the RefundAddress field in the message.
// For EVM chain with coin type Muse the refund address is the tx origin, but can be overridden by the RefundAddress field in the message.
// For EVM chain with coin type Gas the refund address is the tx origin, but can be overridden by the RefundAddress field in the message.
func GetRefundAddress(refundAddress string) (ethcommon.Address, error) {
	// make sure a separate refund address is provided for a bitcoin chain as we cannot refund to tx origin or sender in this case
	if refundAddress == "" {
		return ethcommon.Address{}, errorsmod.Wrap(types.ErrInvalidAddress, "refund address is required")
	}
	if !ethcommon.IsHexAddress(refundAddress) {
		return ethcommon.Address{}, errorsmod.Wrap(types.ErrInvalidAddress, "invalid refund address provided")
	}
	ethRefundAddress := ethcommon.HexToAddress(refundAddress)
	// Double check to make sure the refund address is valid
	if ethRefundAddress == (ethcommon.Address{}) {
		return ethcommon.Address{}, errorsmod.Wrap(types.ErrInvalidAddress, "invalid refund address")
	}
	return ethRefundAddress, nil
}
