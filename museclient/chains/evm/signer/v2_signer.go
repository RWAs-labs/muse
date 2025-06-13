package signer

import (
	"context"
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/RWAs-labs/muse/museclient/chains/evm/common"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// SignOutboundFromCCTXV2 signs an outbound transaction from a CCTX with protocol contract v2
func (signer *Signer) SignOutboundFromCCTXV2(
	ctx context.Context,
	cctx *types.CrossChainTx,
	outboundData *OutboundData,
) (*ethtypes.Transaction, error) {
	outboundType := common.ParseOutboundTypeFromCCTX(*cctx)
	switch outboundType {
	case common.OutboundTypeGasWithdraw, common.OutboundTypeGasWithdrawRevert:
		return signer.SignGasWithdraw(ctx, outboundData)
	case common.OutboundTypeERC20Withdraw, common.OutboundTypeERC20WithdrawRevert:
		return signer.signERC20CustodyWithdraw(ctx, outboundData)
	case common.OutboundTypeERC20WithdrawAndCall:
		return signer.signERC20CustodyWithdrawAndCall(ctx, outboundData)
	case common.OutboundTypeGasWithdrawAndCall, common.OutboundTypeCall:
		// both gas withdraw and call and no-asset call uses gateway execute
		// no-asset call simply hash msg.value == 0
		return signer.signGatewayExecute(ctx, outboundData)
	case common.OutboundTypeGasWithdrawRevertAndCallOnRevert:
		return signer.signGatewayExecuteRevert(ctx, cctx.InboundParams.Sender, outboundData)
	case common.OutboundTypeERC20WithdrawRevertAndCallOnRevert:
		return signer.signERC20CustodyWithdrawRevert(ctx, cctx.InboundParams.Sender, outboundData)
	}
	return nil, fmt.Errorf("unsupported outbound type %d", outboundType)
}
