package keeper

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	connectormevm "github.com/RWAs-labs/protocol-contracts/pkg/museconnectormevm.sol"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/pkg/contracts/sui"
	"github.com/RWAs-labs/muse/pkg/crypto"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

var _ evmtypes.EvmHooks = Hooks{}

type Hooks struct {
	k Keeper
}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing is a wrapper for calling the EVM PostTxProcessing hook on
// the module keeper
func (h Hooks) PostTxProcessing(ctx sdk.Context, msg *core.Message, receipt *ethtypes.Receipt) error {
	return h.k.PostTxProcessing(ctx, msg, receipt)
}

// PostTxProcessing implements EvmHooks.PostTxProcessing.
func (k Keeper) PostTxProcessing(
	ctx sdk.Context,
	msg *core.Message,
	receipt *ethtypes.Receipt,
) error {
	var emittingContract ethcommon.Address
	if msg.To != nil {
		emittingContract = *msg.To
	}
	return k.ProcessLogs(ctx, receipt.Logs, emittingContract, msg.From.Hex())
}

// ProcessLogs post-processes logs emitted by a mEVM contract; if the log contains Withdrawal event
// from registered MRC20 contract, new CCTX will be created to trigger and track outbound
// transaction.
// Returning error from process logs does the following:
// - revert the whole tx.
// - clear the logs
// TODO: implement unit tests
// https://github.com/RWAs-labs/muse/issues/1759
// TODO: refactor and simplify
// https://github.com/RWAs-labs/muse/issues/2627
func (k Keeper) ProcessLogs(
	ctx sdk.Context,
	logs []*ethtypes.Log,
	emittingAddress ethcommon.Address,
	txOrigin string,
) error {
	system, found := k.fungibleKeeper.GetSystemContract(ctx)
	if !found {
		return fmt.Errorf("cannot find system contract")
	}
	connectorMEVMAddr := ethcommon.HexToAddress(system.ConnectorMevm)
	if connectorMEVMAddr == (ethcommon.Address{}) {
		return fmt.Errorf("connectorMEVM address is empty")
	}
	gatewayAddr := ethcommon.HexToAddress(system.Gateway)

	// read the logs and process inbounds from emitted events
	// run the processing for the v1 and the v2 protocol contracts
	for _, log := range logs {
		if !crypto.IsEmptyAddress(gatewayAddr) {
			if err := k.ProcessMEVMInboundV2(ctx, log, gatewayAddr, txOrigin); err != nil {
				return errors.Wrap(err, "failed to process MEVM inbound V2")
			}
		}
		if err := k.ProcessMEVMInboundV1(ctx, log, connectorMEVMAddr, emittingAddress, txOrigin); err != nil {
			return errors.Wrap(err, "failed to process MEVM inbound V1")
		}
	}

	return nil
}

// ProcessMEVMInboundV1 processes the logs emitted by the mEVM contract for V1 protocol contracts
// it parses logs from Connector and MRC20 contracts and processes them accordingly
func (k Keeper) ProcessMEVMInboundV1(
	ctx sdk.Context,
	log *ethtypes.Log,
	connectorMEVMAddr,
	emittingAddress ethcommon.Address,
	txOrigin string,
) error {
	eventMRC20Withdrawal, errMrc20 := ParseMRC20WithdrawalEvent(*log)
	eventMUSESent, errMuseSent := ParseMuseSentEvent(*log, connectorMEVMAddr)
	if errMrc20 != nil && errMuseSent != nil {
		// This log does not contain any of the two events
		return nil
	}
	if eventMRC20Withdrawal != nil && eventMUSESent != nil {
		// This log contains both events, this is not possible
		ctx.Logger().
			Error(fmt.Sprintf("ProcessLogs: log contains both MRC20Withdrawal and MuseSent events, %s , %s", log.Topics, log.Data))
		return nil
	}

	// if eventMrc20Withdrawal is not nil we will try to validate it and see if it can be processed
	if eventMRC20Withdrawal != nil {
		// Check if the contract is a registered MRC20 contract. If its not a registered MRC20 contract, we can discard this event as it is not relevant
		coin, foundCoin := k.fungibleKeeper.GetForeignCoins(ctx, eventMRC20Withdrawal.Raw.Address.Hex())
		if !foundCoin {
			ctx.Logger().
				Info(fmt.Sprintf("cannot find foreign coin with contract address %s", eventMRC20Withdrawal.Raw.Address.Hex()))
			return nil
		}

		// If Validation fails, we will not process the event and return and error. This condition means that the event was correct, and emitted from a registered MRC20 contract
		// But the information entered by the user is incorrect. In this case we can return an error and roll back the transaction
		if err := k.ValidateMRC20WithdrawEvent(ctx, eventMRC20Withdrawal, coin.ForeignChainId, coin.CoinType); err != nil {
			return err
		}
		// If the event is valid, we will process it and create a new CCTX
		// If the process fails, we will return an error and roll back the transaction
		if err := k.ProcessMRC20WithdrawalEvent(ctx, eventMRC20Withdrawal, emittingAddress, txOrigin); err != nil {
			return err
		}
	}
	// if eventMuseSent is not nil we will try to validate it and see if it can be processed
	if eventMUSESent != nil {
		if err := k.ProcessMuseSentEvent(ctx, eventMUSESent, emittingAddress, txOrigin); err != nil {
			return err
		}
	}
	return nil
}

// ProcessMRC20WithdrawalEvent creates a new CCTX to process the withdrawal event
// error indicates system error and non-recoverable; should abort
func (k Keeper) ProcessMRC20WithdrawalEvent(
	ctx sdk.Context,
	event *mrc20.MRC20Withdrawal,
	emittingContract ethcommon.Address,
	txOrigin string,
) error {
	ctx.Logger().Info(fmt.Sprintf("MRC20 withdrawal to %s amount %d", hex.EncodeToString(event.To), event.Value))
	foreignCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, event.Raw.Address.Hex())
	if !found {
		return fmt.Errorf("cannot find foreign coin with emittingContract address %s", event.Raw.Address.Hex())
	}

	receiverChain, found := k.museObserverKeeper.GetSupportedChainFromChainID(ctx, foreignCoin.ForeignChainId)
	if !found {
		return errorsmod.Wrapf(
			observertypes.ErrSupportedChains,
			"chain with chainID %d not supported",
			foreignCoin.ForeignChainId,
		)
	}

	senderChain, err := chains.MuseChainFromCosmosChainID(ctx.ChainID())
	if err != nil {
		return fmt.Errorf("ProcessMRC20WithdrawalEvent: failed to convert chainID: %s", err.Error())
	}

	toAddr, err := receiverChain.EncodeAddress(event.To)
	if err != nil {
		return fmt.Errorf("cannot encode address %s: %s", event.To, err.Error())
	}

	gasLimit, err := k.fungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(foreignCoin.Mrc20ContractAddress))
	if err != nil {
		return fmt.Errorf("cannot query gas limit: %s", err.Error())
	}

	// gasLimit+uint64(event.Raw.Index) to generate different cctx for multiple events in the same tx.
	msg := types.NewMsgVoteInbound(
		"",
		emittingContract.Hex(),
		senderChain.ChainId,
		txOrigin,
		toAddr,
		foreignCoin.ForeignChainId,
		sdkmath.NewUintFromBigInt(event.Value),
		"",
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		gasLimit.Uint64(),
		foreignCoin.CoinType,
		foreignCoin.Asset,
		uint64(event.Raw.Index),
		types.ProtocolContractVersion_V1,
		false, // not relevant for v1
		types.InboundStatus_SUCCESS,
		types.ConfirmationMode_SAFE,
	)

	cctx, err := k.ValidateInbound(ctx, msg, false)
	if err != nil {
		return err
	}

	if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
		return errors.New("cctx aborted")
	}

	EmitMRCWithdrawCreated(ctx, *cctx)

	return nil
}

func (k Keeper) ProcessMuseSentEvent(
	ctx sdk.Context,
	event *connectormevm.MuseConnectorMEVMMuseSent,
	emittingContract ethcommon.Address,
	txOrigin string,
) error {
	ctx.Logger().Info(fmt.Sprintf(
		"Muse withdrawal to %s amount %d to chain with chainId %d",
		hex.EncodeToString(event.DestinationAddress),
		event.MuseValueAndGas,
		event.DestinationChainId,
	))

	if err := k.bankKeeper.BurnCoins(
		ctx,
		fungibletypes.ModuleName,
		sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewIntFromBigInt(event.MuseValueAndGas))),
	); err != nil {
		ctx.Logger().Error(fmt.Sprintf("ProcessMuseSentEvent: failed to burn coins from fungible: %s", err.Error()))
		return fmt.Errorf("ProcessMuseSentEvent: failed to burn coins from fungible: %s", err.Error())
	}

	receiverChainID := event.DestinationChainId

	receiverChain, found := k.museObserverKeeper.GetSupportedChainFromChainID(ctx, receiverChainID.Int64())
	if !found {
		return observertypes.ErrSupportedChains
	}

	// Validation if we want to send MUSE to an external chain, but there is no MUSE token.
	chainParams, found := k.museObserverKeeper.GetChainParamsByChainID(ctx, receiverChain.ChainId)
	if !found {
		return observertypes.ErrChainParamsNotFound
	}

	if receiverChain.IsExternalChain() &&
		(chainParams.MuseTokenContractAddress == "" || chainParams.MuseTokenContractAddress == constant.EVMZeroAddress) {
		return types.ErrUnableToSendCoinType
	}

	toAddr := "0x" + hex.EncodeToString(event.DestinationAddress)
	senderChain, err := chains.MuseChainFromCosmosChainID(ctx.ChainID())
	if err != nil {
		return fmt.Errorf("ProcessMuseSentEvent: failed to convert chainID: %s", err.Error())
	}

	amount := sdkmath.NewUintFromBigInt(event.MuseValueAndGas)
	messageString := base64.StdEncoding.EncodeToString(event.Message)

	// Bump gasLimit by event index (which is very unlikely to be larger than 1000) to always have different MuseSent events msgs.
	msg := types.NewMsgVoteInbound(
		"",
		emittingContract.Hex(),
		senderChain.ChainId,
		txOrigin, toAddr,
		receiverChain.ChainId,
		amount,
		messageString,
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		90000,
		coin.CoinType_Muse,
		"",
		uint64(event.Raw.Index),
		types.ProtocolContractVersion_V1,
		false, // not relevant for v1
		types.InboundStatus_SUCCESS,
		types.ConfirmationMode_SAFE,
	)

	cctx, err := k.ValidateInbound(ctx, msg, true)
	if err != nil {
		return err
	}

	if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
		return errors.New("cctx aborted")
	}

	EmitMuseWithdrawCreated(ctx, *cctx)
	return nil
}

// ValidateMRC20WithdrawEvent checks if the MRC20Withdrawal event is valid
// It verifies event information for BTC chains and returns an error if the event is invalid
func (k Keeper) ValidateMRC20WithdrawEvent(
	ctx sdk.Context,
	event *mrc20.MRC20Withdrawal,
	chainID int64,
	coinType coin.CoinType,
) error {
	// The event was parsed; that means the user has deposited tokens to the contract.
	return k.validateOutbound(ctx, chainID, coinType, event.Value, event.To)
}

// validateOutbound validates the data of a MRC20 Withdrawals and Call event (version 1 or 2)
// it checks if the withdrawal amount is valid and the destination address is supported depending on the chain
func (k Keeper) validateOutbound(
	ctx sdk.Context,
	chainID int64,
	coinType coin.CoinType,
	value *big.Int,
	to []byte,
) error {
	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)
	if chains.IsBitcoinChain(chainID, additionalChains) {
		if value.Cmp(big.NewInt(constant.BTCWithdrawalDustAmount)) < 0 {
			return errorsmod.Wrapf(
				types.ErrInvalidWithdrawalAmount,
				"withdraw amount %s is less than dust amount %d",
				value.String(),
				constant.BTCWithdrawalDustAmount,
			)
		}
		addr, err := chains.DecodeBtcAddress(string(to), chainID)
		if err != nil {
			return errorsmod.Wrapf(types.ErrInvalidAddress, "invalid Bitcoin address %s", string(to))
		}
		if !chains.IsBtcAddressSupported(addr) {
			return errorsmod.Wrapf(types.ErrInvalidAddress, "unsupported Bitcoin address %s", string(to))
		}
	} else if chains.IsSolanaChain(chainID, additionalChains) {
		// The rent exempt check is not needed for MRC20 (SPL) tokens because withdrawing SPL token
		// already needs a non-trivial amount of SOL for potential ATA creation so we can skip the check,
		// and also not needed for simple no asset call.
		if coinType == coin.CoinType_Gas && value.Cmp(big.NewInt(constant.SolanaWalletRentExempt)) < 0 {
			return errorsmod.Wrapf(
				types.ErrInvalidWithdrawalAmount,
				"withdraw amount %s is less than rent exempt %d",
				value.String(),
				constant.SolanaWalletRentExempt,
			)
		}
		_, err := chains.DecodeSolanaWalletAddress(string(to))
		if err != nil {
			return errorsmod.Wrapf(types.ErrInvalidAddress, "invalid Solana address %s", string(to))
		}
	} else if chains.IsSuiChain(chainID, additionalChains) {
		// check the string format of the address is valid

		addr, err := sui.DecodeAddress(to)
		if err != nil {
			return errorsmod.Wrapf(types.ErrInvalidAddress, "invalid Sui address %s", string(to))
		}
		if err := sui.ValidAddress(addr); err != nil {
			return errorsmod.Wrapf(types.ErrInvalidAddress, "invalid Sui address %s", string(to))
		}
	}

	return nil
}

// ParseMRC20WithdrawalEvent tries extracting MRC20Withdrawal event from the input logs using the mrc20 contract;
// It only returns a not-nil event if the event has been correctly validated as a valid withdrawal event
func ParseMRC20WithdrawalEvent(log ethtypes.Log) (*mrc20.MRC20Withdrawal, error) {
	mrc20MEVM, err := mrc20.NewMRC20Filterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, err
	}
	if len(log.Topics) == 0 {
		return nil, fmt.Errorf("ParseMRC20WithdrawalEvent: invalid log - no topics")
	}
	event, err := mrc20MEVM.ParseWithdrawal(log)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// ParseMuseSentEvent tries extracting MuseSent event from connectorMEVM contract;
// returns error if the log entry is not a MuseSent event, or is not emitted from connectorMEVM
// It only returns a not-nil event if all the error checks pass
func ParseMuseSentEvent(
	log ethtypes.Log,
	connectorMEVM ethcommon.Address,
) (*connectormevm.MuseConnectorMEVMMuseSent, error) {
	museConnectorMEVM, err := connectormevm.NewMuseConnectorMEVMFilterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, err
	}
	if len(log.Topics) == 0 {
		return nil, fmt.Errorf("ParseMuseSentEvent: invalid log - no topics")
	}
	event, err := museConnectorMEVM.ParseMuseSent(log)
	if err != nil {
		return nil, err
	}

	if event.Raw.Address != connectorMEVM {
		return nil, fmt.Errorf(
			"ParseMuseSentEvent: event address %s does not match connectorMEVM %s",
			event.Raw.Address.Hex(),
			connectorMEVM.Hex(),
		)
	}
	return event, nil
}
