package common

import "time"

const (
	// OutboundInclusionTimeout is the timeout for waiting for an outbound to be included in a block
	OutboundInclusionTimeout = 20 * time.Minute

	// ReorgProtectBlockCount is confirmations count to protect against reorg
	// Short 1~2 block reorgs could happen often on Ethereum due to network congestion or block production race conditions
	ReorgProtectBlockCount = 2

	// OutboundTrackerReportTimeout is the timeout for waiting for an outbound tracker report
	OutboundTrackerReportTimeout = 10 * time.Minute

	// TopicsMuseSent is the number of topics for a Muse sent event
	// [signature, museTxSenderAddress, destinationChainId]
	// https://github.com/RWAs-labs/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/MuseConnector.base.sol#L34
	TopicsMuseSent = 3

	// TopicsMuseReceived is the number of topics for a Muse received event
	// [signature, sourceChainId, destinationAddress, internalSendHash]
	// https://github.com/RWAs-labs/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/MuseConnector.base.sol#L45
	TopicsMuseReceived = 4

	// TopicsMuseReverted is the number of topics for a Muse reverted event
	// [signature, destinationChainId, internalSendHash]
	// https://github.com/RWAs-labs/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/MuseConnector.base.sol#L54
	TopicsMuseReverted = 3

	// TopicsWithdrawn is the number of topics for a withdrawn event
	// [signature, recipient, asset]
	// https://github.com/RWAs-labs/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ERC20Custody.sol#L43
	TopicsWithdrawn = 3

	// TopicsDeposited is the number of topics for a deposited event
	// [signature, asset]
	// https://github.com/RWAs-labs/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ERC20Custody.sol#L42
	TopicsDeposited = 2

	// V2 contracts

	// TopicsGatewayDeposit is the number of topics for a gateway deposit event
	// [signature, sender, receiver]
	TopicsGatewayDeposit = 3

	// TopicsGatewayDepositAndCall is the number of topics for a gateway deposit and call event
	// [signature, sender, receiver]
	TopicsGatewayDepositAndCall = 3

	// TopicsGatewayCall is the number of topics for a gateway call event
	// [signature, sender, receiver]
	TopicsGatewayCall = 3

	// TopicsGatewayExecuted is the number of topics for a gateway executed event
	// [signature, destination]
	TopicsGatewayExecuted = 2

	// TopicsGatewayExecutedWithERC20 is the number of topics for a gateway executed with ERC20 event
	// [signature, token, destination]
	TopicsGatewayExecutedWithERC20 = 3

	// TopicsGatewayReverted is the number of topics for a reverted event
	// [signature, destination]
	TopicsGatewayReverted = 3

	// TopicsERC20CustodyWithdraw is the number of topics for an ERC20 custody withdraw event
	// [signature, recipient, asset]
	TopicsERC20CustodyWithdraw = 3

	// TopicsERC20CustodyWithdrawAndCall is the number of topics for an ERC20 custody withdraw and call event
	// [signature, recipient, asset]
	TopicsERC20CustodyWithdrawAndCall = 3
)
