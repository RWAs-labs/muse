package testutils

import (
	"github.com/RWAs-labs/protocol-contracts/pkg/erc20custody.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/museconnector.non-eth.sol"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// ParseReceiptMuseSent parses a MuseSent event from a receipt
func ParseReceiptMuseSent(
	receipt *ethtypes.Receipt,
	connector *museconnector.MuseConnectorNonEth,
) *museconnector.MuseConnectorNonEthMuseSent {
	for _, log := range receipt.Logs {
		event, err := connector.ParseMuseSent(*log)
		if err == nil && event != nil {
			return event // found
		}
	}
	return nil
}

// ParseReceiptERC20Deposited parses an Deposited event from a receipt
func ParseReceiptERC20Deposited(
	receipt *ethtypes.Receipt,
	custody *erc20custody.ERC20Custody,
) *erc20custody.ERC20CustodyDeposited {
	for _, log := range receipt.Logs {
		event, err := custody.ParseDeposited(*log)
		if err == nil && event != nil {
			return event // found
		}
	}
	return nil
}
