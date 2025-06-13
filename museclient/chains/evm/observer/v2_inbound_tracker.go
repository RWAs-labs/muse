package observer

import (
	"context"
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/museclient/chains/evm/client"
	"github.com/RWAs-labs/muse/museclient/musecore"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrGatewayNotSet = errors.New("gateway contract not set")
)

// ProcessInboundTrackerV2 processes inbound tracker events from the gateway
// TODO: add test coverage
// https://github.com/RWAs-labs/muse/issues/2669
func (ob *Observer) ProcessInboundTrackerV2(
	ctx context.Context,
	tx *client.Transaction,
	receipt *ethtypes.Receipt,
) error {
	gatewayAddr, gateway, err := ob.getGatewayContract()
	if err != nil {
		ob.Logger().Inbound.Debug().Err(err).Msg("error getting gateway contract for processing inbound tracker")
		return ErrGatewayNotSet
	}

	// check confirmations
	if !ob.IsBlockConfirmedForInboundSafe(receipt.BlockNumber.Uint64()) {
		return fmt.Errorf(
			"inbound %s has not been confirmed yet: receipt block %d",
			tx.Hash,
			receipt.BlockNumber.Uint64(),
		)
	}

	for _, log := range receipt.Logs {
		if log == nil || log.Address != gatewayAddr {
			continue
		}

		// try parsing deposit
		eventDeposit, err := gateway.ParseDeposited(*log)
		if err == nil {
			// check if the event is processable
			if !ob.isEventProcessable(
				eventDeposit.Sender,
				eventDeposit.Receiver,
				eventDeposit.Raw.TxHash,
				eventDeposit.Payload,
			) {
				return fmt.Errorf("event from inbound tracker %s is not processable", tx.Hash)
			}
			msg := ob.newDepositInboundVote(eventDeposit)
			_, err = ob.PostVoteInbound(ctx, &msg, musecore.PostVoteInboundExecutionGasLimit)
			return err
		}

		// try parsing deposit and call
		eventDepositAndCall, err := gateway.ParseDepositedAndCalled(*log)
		if err == nil {
			// check if the event is processable
			if !ob.isEventProcessable(
				eventDepositAndCall.Sender,
				eventDepositAndCall.Receiver,
				eventDepositAndCall.Raw.TxHash,
				eventDepositAndCall.Payload,
			) {
				return fmt.Errorf("event from inbound tracker %s is not processable", tx.Hash)
			}
			msg := ob.newDepositAndCallInboundVote(eventDepositAndCall)
			_, err = ob.PostVoteInbound(ctx, &msg, musecore.PostVoteInboundExecutionGasLimit)
			return err
		}

		// try parsing call
		eventCall, err := gateway.ParseCalled(*log)
		if err == nil {
			// check if the event is processable
			if !ob.isEventProcessable(
				eventCall.Sender,
				eventCall.Receiver,
				eventCall.Raw.TxHash,
				eventCall.Payload,
			) {
				return fmt.Errorf("event from inbound tracker %s is not processable", tx.Hash)
			}
			msg := ob.newCallInboundVote(eventCall)
			_, err = ob.PostVoteInbound(ctx, &msg, musecore.PostVoteInboundExecutionGasLimit)
			return err
		}
	}

	return errors.Wrapf(ErrEventNotFound, "inbound tracker %s", tx.Hash)
}
