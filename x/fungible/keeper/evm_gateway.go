package keeper

import (
	"math/big"

	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	"github.com/RWAs-labs/protocol-contracts/pkg/gatewaymevm.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/RWAs-labs/protocol-contracts/pkg/revert.sol"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/RWAs-labs/muse/pkg/crypto"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// gatewayGasLimit is the gas limit for the gateway functions
var gatewayGasLimit = big.NewInt(1_500_000)

// CallUpdateGatewayAddress calls the updateGatewayAddress function on the MRC20 contract
// function updateGatewayAddress(address addr)
func (k Keeper) CallUpdateGatewayAddress(
	ctx sdk.Context,
	mrc20Address common.Address,
	newGatewayAddress common.Address,
) (*evmtypes.MsgEthereumTxResponse, error) {
	mrc20ABI, err := mrc20.MRC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	return k.CallEVM(
		ctx,
		*mrc20ABI,
		types.ModuleAddressEVM,
		mrc20Address,
		BigIntZero,
		gatewayGasLimit,
		true,
		false,
		"updateGatewayAddress",
		newGatewayAddress,
	)
}

// CallDepositAndCallMRC20 calls the depositAndCall (MRC20 version) function on the gateway contract
// Callable only by the fungible module account
// returns directly CallEVM()
// function depositAndCall(
//
//	    mContext calldata context,
//	    address mrc20,
//	    uint256 amount,
//	    address target,
//	    bytes calldata message
//	)
func (k Keeper) CallDepositAndCallMRC20(
	ctx sdk.Context,
	context gatewaymevm.MessageContext,
	mrc20 common.Address,
	amount *big.Int,
	target common.Address,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	gatewayABI, err := gatewaymevm.GatewayMEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	systemContract, found := k.GetSystemContract(ctx)
	if !found {
		return nil, types.ErrSystemContractNotFound
	}
	gatewayAddr := common.HexToAddress(systemContract.Gateway)
	if crypto.IsEmptyAddress(gatewayAddr) {
		return nil, types.ErrGatewayContractNotSet
	}

	// NOTE:
	// depositAndCall: MUSE version for depositAndCall method
	// depositAndCall0: MRC20 version for depositAndCall method
	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		gatewayGasLimit,
		true,
		false,
		"depositAndCall0",
		context,
		mrc20,
		amount,
		target,
		message,
	)
}

// CallExecute calls the execute function on the gateway contract
// function execute(
//
//	mContext calldata context,
//	address mrc20,
//	uint256 amount,
//	address target,
//	bytes calldata message
//
// )
func (k Keeper) CallExecute(
	ctx sdk.Context,
	context gatewaymevm.MessageContext,
	mrc20 common.Address,
	amount *big.Int,
	target common.Address,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	gatewayABI, err := gatewaymevm.GatewayMEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	systemContract, found := k.GetSystemContract(ctx)
	if !found {
		return nil, types.ErrSystemContractNotFound
	}
	gatewayAddr := common.HexToAddress(systemContract.Gateway)
	if crypto.IsEmptyAddress(gatewayAddr) {
		return nil, types.ErrGatewayContractNotSet
	}

	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		gatewayGasLimit,
		true,
		false,
		"execute",
		context,
		mrc20,
		amount,
		target,
		message,
	)
}

// CallExecuteRevert calls the executeRevert function on the gateway contract
//
//	function executeRevert(
//	address target,
//	RevertContext revertContext,
//	)
func (k Keeper) CallExecuteRevert(
	ctx sdk.Context,
	inboundSender string,
	mrc20 common.Address,
	amount *big.Int,
	target common.Address,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	gatewayABI, err := gatewaymevm.GatewayMEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	systemContract, found := k.GetSystemContract(ctx)
	if !found {
		return nil, types.ErrSystemContractNotFound
	}
	gatewayAddr := common.HexToAddress(systemContract.Gateway)
	if crypto.IsEmptyAddress(gatewayAddr) {
		return nil, types.ErrGatewayContractNotSet
	}

	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		gatewayGasLimit,
		true,
		false,
		"executeRevert",
		target,
		revert.RevertContext{
			Sender:        common.HexToAddress(inboundSender),
			Asset:         mrc20,
			Amount:        amount,
			RevertMessage: message,
		},
	)
}

// CallDepositAndRevert calls the depositAndRevert function on the gateway contract
//
// function depositAndRevert(
//
//	address mrc20,
//	uint256 amount,
//	address target,
//	RevertContext revertContext
//
// )
func (k Keeper) CallDepositAndRevert(
	ctx sdk.Context,
	inboundSender string,
	mrc20 common.Address,
	amount *big.Int,
	target common.Address,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	gatewayABI, err := gatewaymevm.GatewayMEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	systemContract, found := k.GetSystemContract(ctx)
	if !found {
		return nil, types.ErrSystemContractNotFound
	}
	gatewayAddr := common.HexToAddress(systemContract.Gateway)
	if crypto.IsEmptyAddress(gatewayAddr) {
		return nil, types.ErrGatewayContractNotSet
	}

	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		gatewayGasLimit,
		true,
		false,
		"depositAndRevert",
		mrc20,
		amount,
		target,
		revert.RevertContext{
			Sender:        common.HexToAddress(inboundSender),
			Asset:         mrc20,
			Amount:        amount,
			RevertMessage: message,
		},
	)
}

// CallExecuteAbort calls the executeAbort function on the gateway contract
//
//	function executeAbort(
//	address target,
//	AbortContext abortContext,
//	)
func (k Keeper) CallExecuteAbort(
	ctx sdk.Context,
	inboundSender string,
	mrc20 common.Address,
	amount *big.Int,
	outgoing bool,
	chainID *big.Int,
	target common.Address,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, error) {
	gatewayABI, err := gatewaymevm.GatewayMEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	systemContract, found := k.GetSystemContract(ctx)
	if !found {
		return nil, types.ErrSystemContractNotFound
	}
	gatewayAddr := common.HexToAddress(systemContract.Gateway)
	if crypto.IsEmptyAddress(gatewayAddr) {
		return nil, types.ErrGatewayContractNotSet
	}

	// TODO: set correct sender for non EVM chains
	// https://github.com/RWAs-labs/muse/issues/3532
	return k.CallEVM(
		ctx,
		*gatewayABI,
		types.ModuleAddressEVM,
		gatewayAddr,
		BigIntZero,
		gatewayGasLimit,
		true,
		false,
		"executeAbort",
		target,
		revert.AbortContext{
			Sender:        common.HexToAddress(inboundSender).Bytes(),
			Asset:         mrc20,
			Amount:        amount,
			Outgoing:      outgoing,
			ChainID:       chainID,
			RevertMessage: message,
		},
	)
}
