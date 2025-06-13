package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// ProcessMEVMInboundV2 processes the logs emitted by the mEVM contract for V2 protocol contracts
// it parses logs from GatewayMEVM contract and updates the crosschain module state
func (k Keeper) ProcessMEVMInboundV2(
	ctx sdk.Context,
	log *ethtypes.Log,
	gatewayAddr ethcommon.Address,
	txOrigin string,
) error {
	// try to parse a withdrawal event from the log
	withdrawalEvent, callEvent, withdrawalAndCallEvent, err := types.ParseGatewayEvent(*log, gatewayAddr)
	if err == nil && (withdrawalEvent != nil || callEvent != nil || withdrawalAndCallEvent != nil) {
		var inbound *types.MsgVoteInbound

		// parse data from event and validate
		var mrc20 ethcommon.Address
		var value *big.Int
		var receiver []byte
		var contractAddress ethcommon.Address
		if withdrawalEvent != nil {
			mrc20 = withdrawalEvent.Mrc20
			value = withdrawalEvent.Value
			receiver = withdrawalEvent.Receiver
			contractAddress = withdrawalEvent.Raw.Address
		} else if callEvent != nil {
			mrc20 = callEvent.Mrc20
			value = big.NewInt(0)
			receiver = callEvent.Receiver
			contractAddress = callEvent.Raw.Address
		} else {
			mrc20 = withdrawalAndCallEvent.Mrc20
			value = withdrawalAndCallEvent.Value
			receiver = withdrawalAndCallEvent.Receiver
			contractAddress = withdrawalAndCallEvent.Raw.Address
		}

		// get several information necessary for processing the inbound
		foreignCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, mrc20.Hex())
		if !found {
			ctx.Logger().
				Info(fmt.Sprintf("cannot find foreign coin with contract address %s", contractAddress.Hex()))
			return nil
		}
		receiverChain, found := k.museObserverKeeper.GetSupportedChainFromChainID(ctx, foreignCoin.ForeignChainId)
		if !found {
			return errorsmod.Wrapf(
				observertypes.ErrSupportedChains,
				"chain with chainID %d not supported",
				foreignCoin.ForeignChainId,
			)
		}
		gasLimitQueried, err := k.fungibleKeeper.QueryGasLimit(
			ctx,
			ethcommon.HexToAddress(foreignCoin.Mrc20ContractAddress),
		)
		if err != nil {
			return err
		}

		// validate data of the withdrawal event
		coinType := foreignCoin.CoinType
		if callEvent != nil {
			coinType = coin.CoinType_NoAssetCall
		}
		if err := k.validateOutbound(ctx, foreignCoin.ForeignChainId, coinType, value, receiver); err != nil {
			return err
		}

		// create inbound object depending on the event type
		if withdrawalEvent != nil {
			inbound, err = types.NewWithdrawalInbound(
				ctx,
				txOrigin,
				foreignCoin.CoinType,
				foreignCoin.Asset,
				withdrawalEvent,
				receiverChain,
				gasLimitQueried,
			)
			if err != nil {
				return err
			}
		} else if callEvent != nil {
			inbound, err = types.NewCallInbound(
				ctx,
				txOrigin,
				callEvent,
				receiverChain,
				gasLimitQueried,
			)
			if err != nil {
				return err
			}
		} else {
			inbound, err = types.NewWithdrawAndCallInbound(
				ctx,
				txOrigin,
				foreignCoin.CoinType,
				foreignCoin.Asset,
				withdrawalAndCallEvent,
				receiverChain,
				gasLimitQueried,
			)
			if err != nil {
				return err
			}
		}

		if inbound == nil {
			return errors.New("ParseGatewayEvent: invalid log - no event found")
		}

		// validate inbound for processing
		cctx, err := k.ValidateInbound(ctx, inbound, false)
		if err != nil {
			return err
		}
		if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
			return errors.New("cctx aborted")
		}

		EmitMRCWithdrawCreated(ctx, *cctx)
	}

	return nil
}
