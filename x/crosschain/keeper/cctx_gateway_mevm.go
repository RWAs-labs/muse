package keeper

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cctxerror "github.com/RWAs-labs/muse/pkg/errors"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// CCTXGatewayMEVM is implementation of CCTXGateway interface for MEVM
type CCTXGatewayMEVM struct {
	crosschainKeeper Keeper
}

// NewCCTXGatewayMEVM returns new instance of CCTXGatewayMEVM
func NewCCTXGatewayMEVM(crosschainKeeper Keeper) CCTXGatewayMEVM {
	return CCTXGatewayMEVM{
		crosschainKeeper: crosschainKeeper,
	}
}

// InitiateOutbound handles evm deposit and immediately validates pending outbound
func (c CCTXGatewayMEVM) InitiateOutbound(
	ctx sdk.Context,
	config InitiateOutboundConfig,
) (newCCTXStatus types.CctxStatus, err error) {
	switch config.CCTX.InboundParams.Status {
	case types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE:
		// abort if CCTX has insufficient depositor fee for Bitcoin, the CCTX can't be reverted in this case
		// because there is no fund to pay for the revert tx
		c.crosschainKeeper.ProcessAbort(ctx, config.CCTX, types.StatusMessages{
			ErrorMessageOutbound: "insufficient depositor fee",
			StatusMessage:        "inbound observation failed",
		})
		return types.CctxStatus_Aborted, nil
	case types.InboundStatus_INVALID_MEMO:
		// when invalid memo is reported, the CCTX is reverted to the sender
		newCCTXStatus = c.crosschainKeeper.ValidateOutboundMEVM(ctx, config.CCTX, errors.New("invalid memo"), true)
		return newCCTXStatus, nil
	case types.InboundStatus_SUCCESS:
		// process the deposit normally
		tmpCtx, commit := ctx.CacheContext()
		isContractReverted, err := c.crosschainKeeper.HandleEVMDeposit(tmpCtx, config.CCTX)

		if err != nil && !isContractReverted {
			// exceptional case; internal error; should abort CCTX
			// use ctx as tmpCtx is dismissed to not save any side effects performed during the evm deposit
			c.crosschainKeeper.ProcessAbort(ctx, config.CCTX, types.StatusMessages{
				StatusMessage:        "outbound failed but the universal contract did not revert",
				ErrorMessageOutbound: cctxerror.NewCCTXErrorJSONMessage("failed to deposit tokens in MEVM", err),
			})
			return types.CctxStatus_Aborted, err
		}

		newCCTXStatus = c.crosschainKeeper.ValidateOutboundMEVM(ctx, config.CCTX, err, isContractReverted)
		if newCCTXStatus == types.CctxStatus_OutboundMined || newCCTXStatus == types.CctxStatus_PendingRevert {
			commit()
		}

		return newCCTXStatus, nil
	default:
		// unknown observation status, abort the CCTX
		c.crosschainKeeper.ProcessAbort(ctx, config.CCTX, types.StatusMessages{
			ErrorMessageOutbound: fmt.Sprintf("invalid observation status %d", config.CCTX.InboundParams.Status),
			StatusMessage:        "inbound observation failed",
		})
		return types.CctxStatus_Aborted, nil
	}
}
