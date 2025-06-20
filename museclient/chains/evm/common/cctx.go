package common

import (
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// OutboundType enumerate the different types of outbound transactions
// NOTE: only used for v2 protocol contracts and currently excludes MUSE withdraws
type OutboundType int

const (
	// OutboundTypeUnknown is an unknown outbound transaction
	OutboundTypeUnknown OutboundType = iota

	// OutboundTypeGasWithdraw is a gas withdraw transaction
	OutboundTypeGasWithdraw

	// OutboundTypeERC20Withdraw is an ERC20 withdraw transaction
	OutboundTypeERC20Withdraw

	// OutboundTypeGasWithdrawAndCall is a gas withdraw and call transaction
	OutboundTypeGasWithdrawAndCall

	// OutboundTypeERC20WithdrawAndCall is an ERC20 withdraw and call transaction
	OutboundTypeERC20WithdrawAndCall

	// OutboundTypeCall is a no-asset call transaction
	OutboundTypeCall

	// OutboundTypeGasWithdrawRevert is a gas withdraw revert
	OutboundTypeGasWithdrawRevert

	// OutboundTypeGasWithdrawRevertAndCallOnRevert is a gas withdraw revert and call on revert
	OutboundTypeGasWithdrawRevertAndCallOnRevert

	// OutboundTypeERC20WithdrawRevert is an ERC20 withdraw revert
	OutboundTypeERC20WithdrawRevert

	// OutboundTypeERC20WithdrawRevertAndCallOnRevert is an ERC20 withdraw revert and call on revert
	OutboundTypeERC20WithdrawRevertAndCallOnRevert
)

// ParseOutboundTypeFromCCTX returns the outbound type from the CCTX
func ParseOutboundTypeFromCCTX(cctx types.CrossChainTx) OutboundType {
	switch cctx.InboundParams.CoinType {
	case coin.CoinType_Gas:
		switch cctx.CctxStatus.Status {
		case types.CctxStatus_PendingOutbound:
			if cctx.InboundParams.IsCrossChainCall {
				return OutboundTypeGasWithdrawAndCall
			} else {
				return OutboundTypeGasWithdraw
			}
		case types.CctxStatus_PendingRevert:
			if cctx.RevertOptions.CallOnRevert {
				return OutboundTypeGasWithdrawRevertAndCallOnRevert
			} else {
				return OutboundTypeGasWithdrawRevert
			}
		}
	case coin.CoinType_ERC20:
		switch cctx.CctxStatus.Status {
		case types.CctxStatus_PendingOutbound:
			if cctx.InboundParams.IsCrossChainCall {
				return OutboundTypeERC20WithdrawAndCall
			} else {
				return OutboundTypeERC20Withdraw
			}
		case types.CctxStatus_PendingRevert:
			if cctx.RevertOptions.CallOnRevert {
				return OutboundTypeERC20WithdrawRevertAndCallOnRevert
			} else {
				return OutboundTypeERC20WithdrawRevert
			}
		}
	case coin.CoinType_NoAssetCall:
		if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
			return OutboundTypeCall
		}
	}

	return OutboundTypeUnknown
}
