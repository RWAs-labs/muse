// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gatewaymevm

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// AbortContext is an auto generated low-level Go binding around an user-defined struct.
type AbortContext struct {
	Sender        []byte
	Asset         common.Address
	Amount        *big.Int
	Outgoing      bool
	ChainID       *big.Int
	RevertMessage []byte
}

// CallOptions is an auto generated low-level Go binding around an user-defined struct.
type CallOptions struct {
	GasLimit        *big.Int
	IsArbitraryCall bool
}

// MessageContext is an auto generated low-level Go binding around an user-defined struct.
type MessageContext struct {
	Sender    []byte
	SenderEVM common.Address
	ChainID   *big.Int
}

// RevertContext is an auto generated low-level Go binding around an user-defined struct.
type RevertContext struct {
	Sender        common.Address
	Asset         common.Address
	Amount        *big.Int
	RevertMessage []byte
}

// RevertOptions is an auto generated low-level Go binding around an user-defined struct.
type RevertOptions struct {
	RevertAddress    common.Address
	CallOnRevert     bool
	AbortAddress     common.Address
	RevertMessage    []byte
	OnRevertGasLimit *big.Int
}

// GatewayMEVMMetaData contains all meta data concerning the GatewayMEVM contract.
var GatewayMEVMMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_MESSAGE_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_GAS_LIMIT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAUSER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PROTOCOL_ADDRESS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"call\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"mrc20\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callOptions\",\"type\":\"tuple\",\"internalType\":\"structCallOptions\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isArbitraryCall\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"revertOptions\",\"type\":\"tuple\",\"internalType\":\"structRevertOptions\",\"components\":[{\"name\":\"revertAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callOnRevert\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"abortAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"onRevertGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"mrc20\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositAndCall\",\"inputs\":[{\"name\":\"context\",\"type\":\"tuple\",\"internalType\":\"structMessageContext\",\"components\":[{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"senderEVM\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"chainID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositAndCall\",\"inputs\":[{\"name\":\"context\",\"type\":\"tuple\",\"internalType\":\"structMessageContext\",\"components\":[{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"senderEVM\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"chainID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"mrc20\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositAndRevert\",\"inputs\":[{\"name\":\"mrc20\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertContext\",\"type\":\"tuple\",\"internalType\":\"structRevertContext\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"execute\",\"inputs\":[{\"name\":\"context\",\"type\":\"tuple\",\"internalType\":\"structMessageContext\",\"components\":[{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"senderEVM\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"chainID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"mrc20\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeAbort\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"abortContext\",\"type\":\"tuple\",\"internalType\":\"structAbortContext\",\"components\":[{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"outgoing\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"chainID\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executeRevert\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertContext\",\"type\":\"tuple\",\"internalType\":\"structRevertContext\",\"components\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"museToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"admin_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"mrc20\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertOptions\",\"type\":\"tuple\",\"internalType\":\"structRevertOptions\",\"components\":[{\"name\":\"revertAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callOnRevert\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"abortAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"onRevertGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRevertOptions\",\"components\":[{\"name\":\"revertAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callOnRevert\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"abortAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"onRevertGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawAndCall\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structCallOptions\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isArbitraryCall\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRevertOptions\",\"components\":[{\"name\":\"revertAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callOnRevert\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"abortAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"onRevertGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawAndCall\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"mrc20\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"callOptions\",\"type\":\"tuple\",\"internalType\":\"structCallOptions\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isArbitraryCall\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"revertOptions\",\"type\":\"tuple\",\"internalType\":\"structRevertOptions\",\"components\":[{\"name\":\"revertAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callOnRevert\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"abortAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"onRevertGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"museToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Called\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"mrc20\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"callOptions\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structCallOptions\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isArbitraryCall\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"revertOptions\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRevertOptions\",\"components\":[{\"name\":\"revertAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callOnRevert\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"abortAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"onRevertGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"mrc20\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasfee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"protocolFlatFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"callOptions\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structCallOptions\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isArbitraryCall\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"revertOptions\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRevertOptions\",\"components\":[{\"name\":\"revertAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callOnRevert\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"abortAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"onRevertGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawnAndCalled\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"receiver\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"mrc20\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"gasfee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"protocolFlatFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"callOptions\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structCallOptions\",\"components\":[{\"name\":\"gasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"isArbitraryCall\",\"type\":\"bool\",\"internalType\":\"bool\"}]},{\"name\":\"revertOptions\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structRevertOptions\",\"components\":[{\"name\":\"revertAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"callOnRevert\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"abortAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"revertMessage\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"onRevertGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CallOnRevertNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CallerIsNotProtocol\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedMuseSent\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"GasFeeTransferFailed\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientGasLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMRC20Amount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMuseAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTarget\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MessageSizeExceeded\",\"inputs\":[{\"name\":\"provided\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maximum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyWMUSEOrProtocol\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"WithdrawalFailed\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MUSENotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MRC20BurnFailed\",\"inputs\":[{\"name\":\"mrc20\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MRC20DepositFailed\",\"inputs\":[{\"name\":\"mrc20\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MRC20TransferFailed\",\"inputs\":[{\"name\":\"mrc20\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
	Bin: "0x60a06040523060805234801561001457600080fd5b5061001d610022565b6100d4565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100725760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100d15780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b608051613bd26100fd6000396000818161226001528181612289015261243c0152613bd26000f3fe6080604052600436106101dc5760003560e01c80635c975abb116101025780639d4ba46511610095578063c39aca3711610064578063c39aca3714610669578063d547741f14610689578063e63ab1e9146106a9578063f45346dc146106dd57600080fd5b80639d4ba465146105be578063a217fddf146105de578063ad3cb1cc146105f3578063bcf7f32b1461064957600080fd5b80638456cb59116100d15780638456cb591461051357806391d148541461052857806397a1cef11461058d57806397d340f5146105a857600080fd5b80635c975abb146104855780637b15118b146104bc5780637c0dcb5f146104dc5780637ce1ffeb146104fc57600080fd5b80632722feee1161017a5780633f4ba83a116101495780633f4ba83a14610428578063485cc9551461043d5780634f1ef2861461045d57806352d1902d1461047057600080fd5b80632722feee146103a05780632810ae63146103c85780632f2ff15d146103e857806336568abe1461040857600080fd5b80632095dedb116101b65780632095dedb146102cb57806321501a95146102eb57806321e093b11461030b578063248a9ca31461034357600080fd5b806301ffc9a71461025657806306cb89831461028b578063184b0793146102ab57600080fd5b36610251576101e96106fd565b6000546001600160a01b0316331480159061021857503373735b14bb79463307aacbed86daf3322b1e6226ab14155b1561024f576040517fb3af013700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b005b600080fd5b34801561026257600080fd5b50610276610271366004612d0f565b61075b565b60405190151581526020015b60405180910390f35b34801561029757600080fd5b5061024f6102a6366004612e97565b6107f4565b3480156102b757600080fd5b5061024f6102c6366004612f67565b6108e2565b3480156102d757600080fd5b5061024f6102e6366004612fb7565b610a23565b3480156102f757600080fd5b5061024f610306366004613020565b610b05565b34801561031757600080fd5b5060005461032b906001600160a01b031681565b6040516001600160a01b039091168152602001610282565b34801561034f57600080fd5b5061039261035e3660046130ac565b60009081527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602052604090206001015490565b604051908152602001610282565b3480156103ac57600080fd5b5061032b73735b14bb79463307aacbed86daf3322b1e6226ab81565b3480156103d457600080fd5b5061024f6103e33660046130c5565b610d03565b3480156103f457600080fd5b5061024f610403366004613185565b610d3d565b34801561041457600080fd5b5061024f610423366004613185565b610d87565b34801561043457600080fd5b5061024f610dd8565b34801561044957600080fd5b5061024f6104583660046131aa565b610e0d565b61024f61046b3660046131d8565b611063565b34801561047c57600080fd5b5061039261107e565b34801561049157600080fd5b507fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16610276565b3480156104c857600080fd5b5061024f6104d736600461321e565b6110ad565b3480156104e857600080fd5b5061024f6104f7366004613290565b611264565b34801561050857600080fd5b50610392620186a081565b34801561051f57600080fd5b5061024f611481565b34801561053457600080fd5b50610276610543366004613185565b60009182527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408084206001600160a01b0393909316845291905290205460ff1690565b34801561059957600080fd5b5061024f6103e3366004613315565b3480156105b457600080fd5b5061039261080081565b3480156105ca57600080fd5b5061024f6105d9366004613379565b6114b3565b3480156105ea57600080fd5b50610392600081565b3480156105ff57600080fd5b5061063c6040518060400160405280600581526020017f352e302e3000000000000000000000000000000000000000000000000000000081525081565b6040516102829190613429565b34801561065557600080fd5b5061024f61066436600461343c565b611703565b34801561067557600080fd5b5061024f61068436600461343c565b61185c565b34801561069557600080fd5b5061024f6106a4366004613185565b611a0c565b3480156106b557600080fd5b506103927f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a81565b3480156106e957600080fd5b5061024f6106f83660046134da565b611a50565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff1615610759576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b0000000000000000000000000000000000000000000000000000000014806107ee57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6107fc6106fd565b620186a08235101561083a576040517f60ee124700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61080061084a606083018361351c565b610855915085613581565b11156108b657610868606082018261351c565b610873915084613581565b6040517fcd6f4e6d000000000000000000000000000000000000000000000000000000008152600481019190915261080060248201526044015b60405180910390fd5b6108da868686866108cc368890038801886135c9565b6108d587613621565b611bf8565b505050505050565b6108ea611d8d565b3373735b14bb79463307aacbed86daf3322b1e6226ab14610937576040517f42c0407e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61093f6106fd565b6001600160a01b03821661097f576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517fc9028a360000000000000000000000000000000000000000000000000000000081526001600160a01b0383169063c9028a36906109c4908490600401613758565b600060405180830381600087803b1580156109de57600080fd5b505af11580156109f2573d6000803e3d6000fd5b50505050610a1f60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5050565b610a2b611d8d565b3373735b14bb79463307aacbed86daf3322b1e6226ab14610a78576040517f42c0407e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610a806106fd565b6001600160a01b038216610ac0576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f2d4cfb7e0000000000000000000000000000000000000000000000000000000081526001600160a01b03831690632d4cfb7e906109c49084906004016137c8565b610b0d611d8d565b3373735b14bb79463307aacbed86daf3322b1e6226ab14610b5a576040517f42c0407e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610b626106fd565b6001600160a01b038316610ba2576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83600003610bdc576040517f19c08f4900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b03831673735b14bb79463307aacbed86daf3322b1e6226ab1480610c0f57506001600160a01b03831630145b15610c46576040517f82d5d76a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610c508484611e34565b6000546040517f5bcfd6160000000000000000000000000000000000000000000000000000000081526001600160a01b0380861692635bcfd61692610ca1928a92169089908890889060040161386f565b600060405180830381600087803b158015610cbb57600080fd5b505af1158015610ccf573d6000803e3d6000fd5b50505050610cfc60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5050505050565b610d0b6106fd565b6040517fe4dd681d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154610d7781611fe2565b610d818383611fec565b50505050565b6001600160a01b0381163314610dc9576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610dd382826120d9565b505050565b7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a610e0281611fe2565b610e0a61219d565b50565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000810460ff16159067ffffffffffffffff16600081158015610e585750825b905060008267ffffffffffffffff166001148015610e755750303b155b905081158015610e83575080155b15610eba576040517ff92ee8a900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b84547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001660011785558315610f1b5784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680100000000000000001785555b6001600160a01b0387161580610f3857506001600160a01b038616155b15610f6f576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610f7761222d565b610f7f61222d565b610f87612235565b610f8f612245565b610f9a600087611fec565b50610fc57f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a87611fec565b50600080547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b038916179055831561105a5784547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50505050505050565b61106b612255565b61107482612325565b610a1f8282612330565b6000611088612431565b507f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc90565b6110b56106fd565b86516000036110f0576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8560000361112a576040517f5d67094f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620186a082351015611168576040517f60ee124700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610800611178606083018361351c565b611183915085613581565b111561119657610868606082018261351c565b60006111a487878535612493565b90506000336001600160a01b03167fd90f94752d2b12f364f4a2237ebe1aff24ba6127585376bf4935f6a7be17dd2a8a898b868c6001600160a01b0316634d8943bb6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611215573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061123991906138ec565b8c8c8c8c60405161125299989796959493929190613986565b60405180910390a35050505050505050565b61126c6106fd565b83516000036112a7576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b826000036112e1576040517f5d67094f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6108006112f1606083018361351c565b9050111561134657611306606082018261351c565b6040517fcd6f4e6d0000000000000000000000000000000000000000000000000000000081526108ad925061080090600401918252602082015260400190565b60006113528484612649565b90506000336001600160a01b03167f07bf64173efd8f3dfb9e4eb3834bab9d5b85a3d89a1c6425797329de0668502c87868886896001600160a01b0316634d8943bb6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156113c3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906113e791906138ec565b60405180604001604052808c6001600160a01b031663091d27886040518163ffffffff1660e01b8152600401602060405180830381865afa158015611430573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061145491906138ec565b81526001602090910152604051611472969594939291908c90613a10565b60405180910390a35050505050565b7f65d7a28e3265b37a6474929f336521b332c1681b933f6cb9f3376673440d862a6114ab81611fe2565b610e0a6126b7565b6114bb611d8d565b3373735b14bb79463307aacbed86daf3322b1e6226ab14611508576040517f42c0407e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6115106106fd565b6001600160a01b038416158061152d57506001600160a01b038216155b15611564576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8260000361159e576040517f5d67094f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b03821673735b14bb79463307aacbed86daf3322b1e6226ab14806115d157506001600160a01b03821630145b15611608576040517f82d5d76a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611613848385612730565b611663576040517f40a143ba0000000000000000000000000000000000000000000000000000000081526001600160a01b03808616600483015283166024820152604481018490526064016108ad565b6040517fc9028a360000000000000000000000000000000000000000000000000000000081526001600160a01b0383169063c9028a36906116a8908490600401613758565b600060405180830381600087803b1580156116c257600080fd5b505af11580156116d6573d6000803e3d6000fd5b50505050610d8160017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b61170b611d8d565b3373735b14bb79463307aacbed86daf3322b1e6226ab14611758576040517f42c0407e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6117606106fd565b6001600160a01b038516158061177d57506001600160a01b038316155b156117b4576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040517f5bcfd6160000000000000000000000000000000000000000000000000000000081526001600160a01b03841690635bcfd61690611801908990899089908890889060040161386f565b600060405180830381600087803b15801561181b57600080fd5b505af115801561182f573d6000803e3d6000fd5b505050506108da60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b611864611d8d565b3373735b14bb79463307aacbed86daf3322b1e6226ab146118b1576040517f42c0407e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6118b96106fd565b6001600160a01b03851615806118d657506001600160a01b038316155b1561190d576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b83600003611947576040517f5d67094f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b03831673735b14bb79463307aacbed86daf3322b1e6226ab148061197a57506001600160a01b03831630145b156119b1576040517f82d5d76a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6119bc858486612730565b6117b4576040517f40a143ba0000000000000000000000000000000000000000000000000000000081526001600160a01b03808716600483015284166024820152604481018590526064016108ad565b60008281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b6268006020526040902060010154611a4681611fe2565b610d8183836120d9565b3373735b14bb79463307aacbed86daf3322b1e6226ab14611a9d576040517f42c0407e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611aa56106fd565b6001600160a01b0383161580611ac257506001600160a01b038116155b15611af9576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81600003611b33576040517f5d67094f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600160a01b03811673735b14bb79463307aacbed86daf3322b1e6226ab1480611b6657506001600160a01b03811630145b15611b9d576040517f82d5d76a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611ba8838284612730565b610dd3576040517f40a143ba0000000000000000000000000000000000000000000000000000000081526001600160a01b03808516600483015282166024820152604481018390526064016108ad565b8551600003611c33576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81516040517ffc5fecd5000000000000000000000000000000000000000000000000000000008152600481019190915260009081906001600160a01b0388169063fc5fecd5906024016040805180830381865afa158015611c98573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611cbc9190613a92565b91509150611ce0823373735b14bb79463307aacbed86daf3322b1e6226ab846127cb565b611d42576040517f667011120000000000000000000000000000000000000000000000000000000081526001600160a01b038316600482015273735b14bb79463307aacbed86daf3322b1e6226ab6024820152604481018290526064016108ad565b866001600160a01b0316336001600160a01b03167f306ee13f48319a123b222c69908e44dcf91abffc20cacc502e3cf5a4ff23e0e48a89898989604051611252959493929190613ac0565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0080547ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01611e08576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60029055565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b600054611e4c906001600160a01b03163330856127cb565b611e8b576040517ff279af7e000000000000000000000000000000000000000000000000000000008152306004820152602481018390526044016108ad565b6000546040517f2e1a7d4d000000000000000000000000000000000000000000000000000000008152600481018490526001600160a01b0390911690632e1a7d4d90602401600060405180830381600087803b158015611eea57600080fd5b505af1925050508015611efb575060015b611f43576040517ff279af7e0000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152602481018390526044016108ad565b6000816001600160a01b03168360405160006040518083038185875af1925050503d8060008114611f90576040519150601f19603f3d011682016040523d82523d6000602084013e611f95565b606091505b5050905080610dd3576040517ff279af7e0000000000000000000000000000000000000000000000000000000081526001600160a01b0383166004820152602481018490526044016108ad565b610e0a8133612872565b60008281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff166120cf576000848152602082815260408083206001600160a01b0387168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556120853390565b6001600160a01b0316836001600160a01b0316857f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a460019150506107ee565b60009150506107ee565b60008281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602081815260408084206001600160a01b038616855290915282205460ff16156120cf576000848152602082815260408083206001600160a01b038716808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339287917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a460019150506107ee565b6121a56128ff565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001681557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b6040516001600160a01b03909116815260200160405180910390a150565b61075961295a565b61223d61295a565b6107596129c1565b61224d61295a565b610759612a12565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614806122ee57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166122e27f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b6001600160a01b031614155b15610759576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000610a1f81611fe2565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561238a575060408051601f3d908101601f19168201909252612387918101906138ec565b60015b6123cb576040517f4c9c8ce30000000000000000000000000000000000000000000000000000000081526001600160a01b03831660048201526024016108ad565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc8114612427576040517faa1d49a4000000000000000000000000000000000000000000000000000000008152600481018290526024016108ad565b610dd38383612a1a565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614610759576040517fe07c8dba00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000806000846001600160a01b031663fc5fecd5856040518263ffffffff1660e01b81526004016124c691815260200190565b6040805180830381865afa1580156124e2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906125069190613a92565b9150915061252a823373735b14bb79463307aacbed86daf3322b1e6226ab846127cb565b61258c576040517f667011120000000000000000000000000000000000000000000000000000000081526001600160a01b038316600482015273735b14bb79463307aacbed86daf3322b1e6226ab6024820152604481018290526064016108ad565b612598853330896127cb565b6125ec576040517f489ca9b70000000000000000000000000000000000000000000000000000000081526001600160a01b0386166004820152336024820152306044820152606481018790526084016108ad565b6125f68587612a70565b61263e576040517f7112ae770000000000000000000000000000000000000000000000000000000081526001600160a01b0386166004820152602481018790526044016108ad565b9150505b9392505050565b60006126428383846001600160a01b031663091d27886040518163ffffffff1660e01b8152600401602060405180830381865afa15801561268e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126b291906138ec565b612493565b6126bf6106fd565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011781557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2583361220f565b6040517f47e7ef240000000000000000000000000000000000000000000000000000000081526001600160a01b03838116600483015260248201839052600091908516906347e7ef24906044016020604051808303816000875af19250505080156127b8575060408051601f3d908101601f191682019092526127b591810190613b63565b60015b6127c457506000612642565b9050612642565b6040517f23b872dd0000000000000000000000000000000000000000000000000000000081526001600160a01b038481166004830152838116602483015260448201839052600091908616906323b872dd906064016020604051808303816000875af192505050801561285b575060408051601f3d908101601f1916820190925261285891810190613b63565b60015b6128675750600061286a565b90505b949350505050565b60008281527f02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800602090815260408083206001600160a01b038516845290915290205460ff16610a1f576040517fe2517d3f0000000000000000000000000000000000000000000000000000000081526001600160a01b0382166004820152602481018390526044016108ad565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f033005460ff16610759576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a005468010000000000000000900460ff16610759576040517fd7e6bcf800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6129c961295a565b7fcd5ed15c6e187e77e9aee88184c21f4f2182ab5827cb3b7e07fbedcd63f0330080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169055565b611e0e61295a565b612a2382612b02565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a2805115612a6857610dd38282612baa565b610a1f612c20565b6040517f42966c68000000000000000000000000000000000000000000000000000000008152600481018290526000906001600160a01b038416906342966c68906024016020604051808303816000875af1925050508015612aef575060408051601f3d908101601f19168201909252612aec91810190613b63565b60015b612afb575060006107ee565b90506107ee565b806001600160a01b03163b600003612b51576040517f4c9c8ce30000000000000000000000000000000000000000000000000000000081526001600160a01b03821660048201526024016108ad565b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff0000000000000000000000000000000000000000166001600160a01b0392909216919091179055565b6060600080846001600160a01b031684604051612bc79190613b80565b600060405180830381855af49150503d8060008114612c02576040519150601f19603f3d011682016040523d82523d6000602084013e612c07565b606091505b5091509150612c17858383612c58565b95945050505050565b3415610759576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606082612c6d57612c6882612ccd565b612642565b8151158015612c8457506001600160a01b0384163b155b15612cc6576040517f9996b3150000000000000000000000000000000000000000000000000000000081526001600160a01b03851660048201526024016108ad565b5080612642565b805115612cdd5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600060208284031215612d2157600080fd5b81357fffffffff000000000000000000000000000000000000000000000000000000008116811461264257600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600082601f830112612d9157600080fd5b813567ffffffffffffffff811115612dab57612dab612d51565b604051601f8201601f19908116603f0116810167ffffffffffffffff81118282101715612dda57612dda612d51565b604052818152838201602001851015612df257600080fd5b816020850160208301376000918101602001919091529392505050565b6001600160a01b0381168114610e0a57600080fd5b60008083601f840112612e3657600080fd5b50813567ffffffffffffffff811115612e4e57600080fd5b602083019150836020828501011115612e6657600080fd5b9250929050565b600060408284031215612e7f57600080fd5b50919050565b600060a08284031215612e7f57600080fd5b60008060008060008060c08789031215612eb057600080fd5b863567ffffffffffffffff811115612ec757600080fd5b612ed389828a01612d80565b9650506020870135612ee481612e0f565b9450604087013567ffffffffffffffff811115612f0057600080fd5b612f0c89828a01612e24565b9095509350612f2090508860608901612e6d565b915060a087013567ffffffffffffffff811115612f3c57600080fd5b612f4889828a01612e85565b9150509295509295509295565b600060808284031215612e7f57600080fd5b60008060408385031215612f7a57600080fd5b8235612f8581612e0f565b9150602083013567ffffffffffffffff811115612fa157600080fd5b612fad85828601612f55565b9150509250929050565b60008060408385031215612fca57600080fd5b8235612fd581612e0f565b9150602083013567ffffffffffffffff811115612ff157600080fd5b830160c0818603121561300357600080fd5b809150509250929050565b600060608284031215612e7f57600080fd5b60008060008060006080868803121561303857600080fd5b853567ffffffffffffffff81111561304f57600080fd5b61305b8882890161300e565b95505060208601359350604086013561307381612e0f565b9250606086013567ffffffffffffffff81111561308f57600080fd5b61309b88828901612e24565b969995985093965092949392505050565b6000602082840312156130be57600080fd5b5035919050565b600080600080600080600060e0888a0312156130e057600080fd5b873567ffffffffffffffff8111156130f757600080fd5b6131038a828b01612d80565b9750506020880135955060408801359450606088013567ffffffffffffffff81111561312e57600080fd5b61313a8a828b01612e24565b909550935061314e90508960808a01612e6d565b915060c088013567ffffffffffffffff81111561316a57600080fd5b6131768a828b01612e85565b91505092959891949750929550565b6000806040838503121561319857600080fd5b82359150602083013561300381612e0f565b600080604083850312156131bd57600080fd5b82356131c881612e0f565b9150602083013561300381612e0f565b600080604083850312156131eb57600080fd5b82356131f681612e0f565b9150602083013567ffffffffffffffff81111561321257600080fd5b612fad85828601612d80565b600080600080600080600060e0888a03121561323957600080fd5b873567ffffffffffffffff81111561325057600080fd5b61325c8a828b01612d80565b97505060208801359550604088013561327481612e0f565b9450606088013567ffffffffffffffff81111561312e57600080fd5b600080600080608085870312156132a657600080fd5b843567ffffffffffffffff8111156132bd57600080fd5b6132c987828801612d80565b9450506020850135925060408501356132e181612e0f565b9150606085013567ffffffffffffffff8111156132fd57600080fd5b61330987828801612e85565b91505092959194509250565b6000806000806080858703121561332b57600080fd5b843567ffffffffffffffff81111561334257600080fd5b61334e87828801612d80565b9450506020850135925060408501359150606085013567ffffffffffffffff8111156132fd57600080fd5b6000806000806080858703121561338f57600080fd5b843561339a81612e0f565b93506020850135925060408501356133b181612e0f565b9150606085013567ffffffffffffffff8111156133cd57600080fd5b61330987828801612f55565b60005b838110156133f45781810151838201526020016133dc565b50506000910152565b600081518084526134158160208601602086016133d9565b601f01601f19169290920160200192915050565b60208152600061264260208301846133fd565b60008060008060008060a0878903121561345557600080fd5b863567ffffffffffffffff81111561346c57600080fd5b61347889828a0161300e565b965050602087013561348981612e0f565b94506040870135935060608701356134a081612e0f565b9250608087013567ffffffffffffffff8111156134bc57600080fd5b6134c889828a01612e24565b979a9699509497509295939492505050565b6000806000606084860312156134ef57600080fd5b83356134fa81612e0f565b925060208401359150604084013561351181612e0f565b809150509250925092565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261355157600080fd5b83018035915067ffffffffffffffff82111561356c57600080fd5b602001915036819003821315612e6657600080fd5b808201808211156107ee577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8015158114610e0a57600080fd5b600060408284031280156135dc57600080fd5b506040805190810167ffffffffffffffff8111828210171561360057613600612d51565b604052823581526020830135613615816135bb565b60208201529392505050565b600060a0823603121561363357600080fd5b60405160a0810167ffffffffffffffff8111828210171561365657613656612d51565b604052823561366481612e0f565b81526020830135613674816135bb565b6020820152604083013561368781612e0f565b6040820152606083013567ffffffffffffffff8111156136a657600080fd5b6136b236828601612d80565b606083015250608092830135928101929092525090565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126136fe57600080fd5b830160208101925035905067ffffffffffffffff81111561371e57600080fd5b803603821315612e6657600080fd5b818352818160208501375060006020828401015260006020601f19601f840116840101905092915050565b602081526000823561376981612e0f565b6001600160a01b038116602084015250602083013561378781612e0f565b6001600160a01b038116604084015250600060408401359050806060840152506137b460608401846136c9565b608080850152612c1760a08501828461372d565b6020815260006137d883846136c9565b60c060208501526137ed60e08501828461372d565b91505060208401356137fe81612e0f565b6001600160a01b0316604084810191909152840135606080850191909152840135613828816135bb565b8015156080850152506000608085013590508060a08501525061384e60a08501856136c9565b601f198584030160c086015261386583828461372d565b9695505050505050565b60808152600061387f87886136c9565b6060608085015261389460e08501828461372d565b91505060208801356138a581612e0f565b6001600160a01b0390811660a085015260408981013560c08601529088166020850152830186905282810360608401526138e081858761372d565b98975050505050505050565b6000602082840312156138fe57600080fd5b5051919050565b6000813561391281612e0f565b6001600160a01b03168352602082013561392b816135bb565b15156020840152604082013561394081612e0f565b6001600160a01b0316604084015261395b60608301836136c9565b60a0606086015261397060a08601828461372d565b6080948501359590940194909452509092915050565b6101208152600061399b61012083018c6133fd565b6001600160a01b038b16602084015289604084015288606084015287608084015282810360a08401526139cf81878961372d565b853560c0850152905060208501356139e6816135bb565b151560e0840152828103610100840152613a008185613905565b9c9b505050505050505050505050565b61012081526000613a2561012083018a6133fd565b6001600160a01b03891660208401528760408401528660608401528560808401528281038060a085015260008252613a6c60c0850187805182526020908101511515910152565b6020810161010085015250613a846020820185613905565b9a9950505050505050505050565b60008060408385031215613aa557600080fd5b8251613ab081612e0f565b6020939093015192949293505050565b60a081526000613ad360a08301886133fd565b8281036020840152613ae681878961372d565b85516040850152602086015115156060850152905082810360808401526001600160a01b0384511681526020840151151560208201526001600160a01b036040850151166040820152606084015160a06060830152613b4860a08301826133fd565b90506080850151608083015280925050509695505050505050565b600060208284031215613b7557600080fd5b8151612642816135bb565b60008251613b928184602087016133d9565b919091019291505056fea2646970667358221220ad7acef9e4b39e104bb8f9b45093f60727c371457d5b96b12237a69ea277706364736f6c634300081a0033",
}

// GatewayMEVMABI is the input ABI used to generate the binding from.
// Deprecated: Use GatewayMEVMMetaData.ABI instead.
var GatewayMEVMABI = GatewayMEVMMetaData.ABI

// GatewayMEVMBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GatewayMEVMMetaData.Bin instead.
var GatewayMEVMBin = GatewayMEVMMetaData.Bin

// DeployGatewayMEVM deploys a new Ethereum contract, binding an instance of GatewayMEVM to it.
func DeployGatewayMEVM(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *GatewayMEVM, error) {
	parsed, err := GatewayMEVMMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GatewayMEVMBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GatewayMEVM{GatewayMEVMCaller: GatewayMEVMCaller{contract: contract}, GatewayMEVMTransactor: GatewayMEVMTransactor{contract: contract}, GatewayMEVMFilterer: GatewayMEVMFilterer{contract: contract}}, nil
}

// GatewayMEVM is an auto generated Go binding around an Ethereum contract.
type GatewayMEVM struct {
	GatewayMEVMCaller     // Read-only binding to the contract
	GatewayMEVMTransactor // Write-only binding to the contract
	GatewayMEVMFilterer   // Log filterer for contract events
}

// GatewayMEVMCaller is an auto generated read-only Go binding around an Ethereum contract.
type GatewayMEVMCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GatewayMEVMTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GatewayMEVMTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GatewayMEVMFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GatewayMEVMFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GatewayMEVMSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GatewayMEVMSession struct {
	Contract     *GatewayMEVM      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GatewayMEVMCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GatewayMEVMCallerSession struct {
	Contract *GatewayMEVMCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// GatewayMEVMTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GatewayMEVMTransactorSession struct {
	Contract     *GatewayMEVMTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// GatewayMEVMRaw is an auto generated low-level Go binding around an Ethereum contract.
type GatewayMEVMRaw struct {
	Contract *GatewayMEVM // Generic contract binding to access the raw methods on
}

// GatewayMEVMCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GatewayMEVMCallerRaw struct {
	Contract *GatewayMEVMCaller // Generic read-only contract binding to access the raw methods on
}

// GatewayMEVMTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GatewayMEVMTransactorRaw struct {
	Contract *GatewayMEVMTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGatewayMEVM creates a new instance of GatewayMEVM, bound to a specific deployed contract.
func NewGatewayMEVM(address common.Address, backend bind.ContractBackend) (*GatewayMEVM, error) {
	contract, err := bindGatewayMEVM(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVM{GatewayMEVMCaller: GatewayMEVMCaller{contract: contract}, GatewayMEVMTransactor: GatewayMEVMTransactor{contract: contract}, GatewayMEVMFilterer: GatewayMEVMFilterer{contract: contract}}, nil
}

// NewGatewayMEVMCaller creates a new read-only instance of GatewayMEVM, bound to a specific deployed contract.
func NewGatewayMEVMCaller(address common.Address, caller bind.ContractCaller) (*GatewayMEVMCaller, error) {
	contract, err := bindGatewayMEVM(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMCaller{contract: contract}, nil
}

// NewGatewayMEVMTransactor creates a new write-only instance of GatewayMEVM, bound to a specific deployed contract.
func NewGatewayMEVMTransactor(address common.Address, transactor bind.ContractTransactor) (*GatewayMEVMTransactor, error) {
	contract, err := bindGatewayMEVM(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMTransactor{contract: contract}, nil
}

// NewGatewayMEVMFilterer creates a new log filterer instance of GatewayMEVM, bound to a specific deployed contract.
func NewGatewayMEVMFilterer(address common.Address, filterer bind.ContractFilterer) (*GatewayMEVMFilterer, error) {
	contract, err := bindGatewayMEVM(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMFilterer{contract: contract}, nil
}

// bindGatewayMEVM binds a generic wrapper to an already deployed contract.
func bindGatewayMEVM(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GatewayMEVMMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GatewayMEVM *GatewayMEVMRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GatewayMEVM.Contract.GatewayMEVMCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GatewayMEVM *GatewayMEVMRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.GatewayMEVMTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GatewayMEVM *GatewayMEVMRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.GatewayMEVMTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GatewayMEVM *GatewayMEVMCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GatewayMEVM.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GatewayMEVM *GatewayMEVMTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GatewayMEVM *GatewayMEVMTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GatewayMEVM.Contract.DEFAULTADMINROLE(&_GatewayMEVM.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GatewayMEVM.Contract.DEFAULTADMINROLE(&_GatewayMEVM.CallOpts)
}

// MAXMESSAGESIZE is a free data retrieval call binding the contract method 0x97d340f5.
//
// Solidity: function MAX_MESSAGE_SIZE() view returns(uint256)
func (_GatewayMEVM *GatewayMEVMCaller) MAXMESSAGESIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "MAX_MESSAGE_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXMESSAGESIZE is a free data retrieval call binding the contract method 0x97d340f5.
//
// Solidity: function MAX_MESSAGE_SIZE() view returns(uint256)
func (_GatewayMEVM *GatewayMEVMSession) MAXMESSAGESIZE() (*big.Int, error) {
	return _GatewayMEVM.Contract.MAXMESSAGESIZE(&_GatewayMEVM.CallOpts)
}

// MAXMESSAGESIZE is a free data retrieval call binding the contract method 0x97d340f5.
//
// Solidity: function MAX_MESSAGE_SIZE() view returns(uint256)
func (_GatewayMEVM *GatewayMEVMCallerSession) MAXMESSAGESIZE() (*big.Int, error) {
	return _GatewayMEVM.Contract.MAXMESSAGESIZE(&_GatewayMEVM.CallOpts)
}

// MINGASLIMIT is a free data retrieval call binding the contract method 0x7ce1ffeb.
//
// Solidity: function MIN_GAS_LIMIT() view returns(uint256)
func (_GatewayMEVM *GatewayMEVMCaller) MINGASLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "MIN_GAS_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINGASLIMIT is a free data retrieval call binding the contract method 0x7ce1ffeb.
//
// Solidity: function MIN_GAS_LIMIT() view returns(uint256)
func (_GatewayMEVM *GatewayMEVMSession) MINGASLIMIT() (*big.Int, error) {
	return _GatewayMEVM.Contract.MINGASLIMIT(&_GatewayMEVM.CallOpts)
}

// MINGASLIMIT is a free data retrieval call binding the contract method 0x7ce1ffeb.
//
// Solidity: function MIN_GAS_LIMIT() view returns(uint256)
func (_GatewayMEVM *GatewayMEVMCallerSession) MINGASLIMIT() (*big.Int, error) {
	return _GatewayMEVM.Contract.MINGASLIMIT(&_GatewayMEVM.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMSession) PAUSERROLE() ([32]byte, error) {
	return _GatewayMEVM.Contract.PAUSERROLE(&_GatewayMEVM.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMCallerSession) PAUSERROLE() ([32]byte, error) {
	return _GatewayMEVM.Contract.PAUSERROLE(&_GatewayMEVM.CallOpts)
}

// PROTOCOLADDRESS is a free data retrieval call binding the contract method 0x2722feee.
//
// Solidity: function PROTOCOL_ADDRESS() view returns(address)
func (_GatewayMEVM *GatewayMEVMCaller) PROTOCOLADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "PROTOCOL_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PROTOCOLADDRESS is a free data retrieval call binding the contract method 0x2722feee.
//
// Solidity: function PROTOCOL_ADDRESS() view returns(address)
func (_GatewayMEVM *GatewayMEVMSession) PROTOCOLADDRESS() (common.Address, error) {
	return _GatewayMEVM.Contract.PROTOCOLADDRESS(&_GatewayMEVM.CallOpts)
}

// PROTOCOLADDRESS is a free data retrieval call binding the contract method 0x2722feee.
//
// Solidity: function PROTOCOL_ADDRESS() view returns(address)
func (_GatewayMEVM *GatewayMEVMCallerSession) PROTOCOLADDRESS() (common.Address, error) {
	return _GatewayMEVM.Contract.PROTOCOLADDRESS(&_GatewayMEVM.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_GatewayMEVM *GatewayMEVMCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_GatewayMEVM *GatewayMEVMSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _GatewayMEVM.Contract.UPGRADEINTERFACEVERSION(&_GatewayMEVM.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_GatewayMEVM *GatewayMEVMCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _GatewayMEVM.Contract.UPGRADEINTERFACEVERSION(&_GatewayMEVM.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GatewayMEVM.Contract.GetRoleAdmin(&_GatewayMEVM.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GatewayMEVM.Contract.GetRoleAdmin(&_GatewayMEVM.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GatewayMEVM *GatewayMEVMCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GatewayMEVM *GatewayMEVMSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GatewayMEVM.Contract.HasRole(&_GatewayMEVM.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GatewayMEVM *GatewayMEVMCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GatewayMEVM.Contract.HasRole(&_GatewayMEVM.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GatewayMEVM *GatewayMEVMCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GatewayMEVM *GatewayMEVMSession) Paused() (bool, error) {
	return _GatewayMEVM.Contract.Paused(&_GatewayMEVM.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GatewayMEVM *GatewayMEVMCallerSession) Paused() (bool, error) {
	return _GatewayMEVM.Contract.Paused(&_GatewayMEVM.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMSession) ProxiableUUID() ([32]byte, error) {
	return _GatewayMEVM.Contract.ProxiableUUID(&_GatewayMEVM.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_GatewayMEVM *GatewayMEVMCallerSession) ProxiableUUID() ([32]byte, error) {
	return _GatewayMEVM.Contract.ProxiableUUID(&_GatewayMEVM.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GatewayMEVM *GatewayMEVMCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GatewayMEVM *GatewayMEVMSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GatewayMEVM.Contract.SupportsInterface(&_GatewayMEVM.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GatewayMEVM *GatewayMEVMCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GatewayMEVM.Contract.SupportsInterface(&_GatewayMEVM.CallOpts, interfaceId)
}

// Withdraw0 is a free data retrieval call binding the contract method 0x97a1cef1.
//
// Solidity: function withdraw(bytes , uint256 , uint256 , (address,bool,address,bytes,uint256) ) view returns()
func (_GatewayMEVM *GatewayMEVMCaller) Withdraw0(opts *bind.CallOpts, arg0 []byte, arg1 *big.Int, arg2 *big.Int, arg3 RevertOptions) error {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "withdraw0", arg0, arg1, arg2, arg3)

	if err != nil {
		return err
	}

	return err

}

// Withdraw0 is a free data retrieval call binding the contract method 0x97a1cef1.
//
// Solidity: function withdraw(bytes , uint256 , uint256 , (address,bool,address,bytes,uint256) ) view returns()
func (_GatewayMEVM *GatewayMEVMSession) Withdraw0(arg0 []byte, arg1 *big.Int, arg2 *big.Int, arg3 RevertOptions) error {
	return _GatewayMEVM.Contract.Withdraw0(&_GatewayMEVM.CallOpts, arg0, arg1, arg2, arg3)
}

// Withdraw0 is a free data retrieval call binding the contract method 0x97a1cef1.
//
// Solidity: function withdraw(bytes , uint256 , uint256 , (address,bool,address,bytes,uint256) ) view returns()
func (_GatewayMEVM *GatewayMEVMCallerSession) Withdraw0(arg0 []byte, arg1 *big.Int, arg2 *big.Int, arg3 RevertOptions) error {
	return _GatewayMEVM.Contract.Withdraw0(&_GatewayMEVM.CallOpts, arg0, arg1, arg2, arg3)
}

// WithdrawAndCall is a free data retrieval call binding the contract method 0x2810ae63.
//
// Solidity: function withdrawAndCall(bytes , uint256 , uint256 , bytes , (uint256,bool) , (address,bool,address,bytes,uint256) ) view returns()
func (_GatewayMEVM *GatewayMEVMCaller) WithdrawAndCall(opts *bind.CallOpts, arg0 []byte, arg1 *big.Int, arg2 *big.Int, arg3 []byte, arg4 CallOptions, arg5 RevertOptions) error {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "withdrawAndCall", arg0, arg1, arg2, arg3, arg4, arg5)

	if err != nil {
		return err
	}

	return err

}

// WithdrawAndCall is a free data retrieval call binding the contract method 0x2810ae63.
//
// Solidity: function withdrawAndCall(bytes , uint256 , uint256 , bytes , (uint256,bool) , (address,bool,address,bytes,uint256) ) view returns()
func (_GatewayMEVM *GatewayMEVMSession) WithdrawAndCall(arg0 []byte, arg1 *big.Int, arg2 *big.Int, arg3 []byte, arg4 CallOptions, arg5 RevertOptions) error {
	return _GatewayMEVM.Contract.WithdrawAndCall(&_GatewayMEVM.CallOpts, arg0, arg1, arg2, arg3, arg4, arg5)
}

// WithdrawAndCall is a free data retrieval call binding the contract method 0x2810ae63.
//
// Solidity: function withdrawAndCall(bytes , uint256 , uint256 , bytes , (uint256,bool) , (address,bool,address,bytes,uint256) ) view returns()
func (_GatewayMEVM *GatewayMEVMCallerSession) WithdrawAndCall(arg0 []byte, arg1 *big.Int, arg2 *big.Int, arg3 []byte, arg4 CallOptions, arg5 RevertOptions) error {
	return _GatewayMEVM.Contract.WithdrawAndCall(&_GatewayMEVM.CallOpts, arg0, arg1, arg2, arg3, arg4, arg5)
}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_GatewayMEVM *GatewayMEVMCaller) MuseToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GatewayMEVM.contract.Call(opts, &out, "museToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_GatewayMEVM *GatewayMEVMSession) MuseToken() (common.Address, error) {
	return _GatewayMEVM.Contract.MuseToken(&_GatewayMEVM.CallOpts)
}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_GatewayMEVM *GatewayMEVMCallerSession) MuseToken() (common.Address, error) {
	return _GatewayMEVM.Contract.MuseToken(&_GatewayMEVM.CallOpts)
}

// Call is a paid mutator transaction binding the contract method 0x06cb8983.
//
// Solidity: function call(bytes receiver, address mrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) Call(opts *bind.TransactOpts, receiver []byte, mrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "call", receiver, mrc20, message, callOptions, revertOptions)
}

// Call is a paid mutator transaction binding the contract method 0x06cb8983.
//
// Solidity: function call(bytes receiver, address mrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayMEVM *GatewayMEVMSession) Call(receiver []byte, mrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Call(&_GatewayMEVM.TransactOpts, receiver, mrc20, message, callOptions, revertOptions)
}

// Call is a paid mutator transaction binding the contract method 0x06cb8983.
//
// Solidity: function call(bytes receiver, address mrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) Call(receiver []byte, mrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Call(&_GatewayMEVM.TransactOpts, receiver, mrc20, message, callOptions, revertOptions)
}

// Deposit is a paid mutator transaction binding the contract method 0xf45346dc.
//
// Solidity: function deposit(address mrc20, uint256 amount, address target) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) Deposit(opts *bind.TransactOpts, mrc20 common.Address, amount *big.Int, target common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "deposit", mrc20, amount, target)
}

// Deposit is a paid mutator transaction binding the contract method 0xf45346dc.
//
// Solidity: function deposit(address mrc20, uint256 amount, address target) returns()
func (_GatewayMEVM *GatewayMEVMSession) Deposit(mrc20 common.Address, amount *big.Int, target common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Deposit(&_GatewayMEVM.TransactOpts, mrc20, amount, target)
}

// Deposit is a paid mutator transaction binding the contract method 0xf45346dc.
//
// Solidity: function deposit(address mrc20, uint256 amount, address target) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) Deposit(mrc20 common.Address, amount *big.Int, target common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Deposit(&_GatewayMEVM.TransactOpts, mrc20, amount, target)
}

// DepositAndCall is a paid mutator transaction binding the contract method 0x21501a95.
//
// Solidity: function depositAndCall((bytes,address,uint256) context, uint256 amount, address target, bytes message) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) DepositAndCall(opts *bind.TransactOpts, context MessageContext, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "depositAndCall", context, amount, target, message)
}

// DepositAndCall is a paid mutator transaction binding the contract method 0x21501a95.
//
// Solidity: function depositAndCall((bytes,address,uint256) context, uint256 amount, address target, bytes message) returns()
func (_GatewayMEVM *GatewayMEVMSession) DepositAndCall(context MessageContext, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.DepositAndCall(&_GatewayMEVM.TransactOpts, context, amount, target, message)
}

// DepositAndCall is a paid mutator transaction binding the contract method 0x21501a95.
//
// Solidity: function depositAndCall((bytes,address,uint256) context, uint256 amount, address target, bytes message) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) DepositAndCall(context MessageContext, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.DepositAndCall(&_GatewayMEVM.TransactOpts, context, amount, target, message)
}

// DepositAndCall0 is a paid mutator transaction binding the contract method 0xc39aca37.
//
// Solidity: function depositAndCall((bytes,address,uint256) context, address mrc20, uint256 amount, address target, bytes message) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) DepositAndCall0(opts *bind.TransactOpts, context MessageContext, mrc20 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "depositAndCall0", context, mrc20, amount, target, message)
}

// DepositAndCall0 is a paid mutator transaction binding the contract method 0xc39aca37.
//
// Solidity: function depositAndCall((bytes,address,uint256) context, address mrc20, uint256 amount, address target, bytes message) returns()
func (_GatewayMEVM *GatewayMEVMSession) DepositAndCall0(context MessageContext, mrc20 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.DepositAndCall0(&_GatewayMEVM.TransactOpts, context, mrc20, amount, target, message)
}

// DepositAndCall0 is a paid mutator transaction binding the contract method 0xc39aca37.
//
// Solidity: function depositAndCall((bytes,address,uint256) context, address mrc20, uint256 amount, address target, bytes message) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) DepositAndCall0(context MessageContext, mrc20 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.DepositAndCall0(&_GatewayMEVM.TransactOpts, context, mrc20, amount, target, message)
}

// DepositAndRevert is a paid mutator transaction binding the contract method 0x9d4ba465.
//
// Solidity: function depositAndRevert(address mrc20, uint256 amount, address target, (address,address,uint256,bytes) revertContext) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) DepositAndRevert(opts *bind.TransactOpts, mrc20 common.Address, amount *big.Int, target common.Address, revertContext RevertContext) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "depositAndRevert", mrc20, amount, target, revertContext)
}

// DepositAndRevert is a paid mutator transaction binding the contract method 0x9d4ba465.
//
// Solidity: function depositAndRevert(address mrc20, uint256 amount, address target, (address,address,uint256,bytes) revertContext) returns()
func (_GatewayMEVM *GatewayMEVMSession) DepositAndRevert(mrc20 common.Address, amount *big.Int, target common.Address, revertContext RevertContext) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.DepositAndRevert(&_GatewayMEVM.TransactOpts, mrc20, amount, target, revertContext)
}

// DepositAndRevert is a paid mutator transaction binding the contract method 0x9d4ba465.
//
// Solidity: function depositAndRevert(address mrc20, uint256 amount, address target, (address,address,uint256,bytes) revertContext) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) DepositAndRevert(mrc20 common.Address, amount *big.Int, target common.Address, revertContext RevertContext) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.DepositAndRevert(&_GatewayMEVM.TransactOpts, mrc20, amount, target, revertContext)
}

// Execute is a paid mutator transaction binding the contract method 0xbcf7f32b.
//
// Solidity: function execute((bytes,address,uint256) context, address mrc20, uint256 amount, address target, bytes message) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) Execute(opts *bind.TransactOpts, context MessageContext, mrc20 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "execute", context, mrc20, amount, target, message)
}

// Execute is a paid mutator transaction binding the contract method 0xbcf7f32b.
//
// Solidity: function execute((bytes,address,uint256) context, address mrc20, uint256 amount, address target, bytes message) returns()
func (_GatewayMEVM *GatewayMEVMSession) Execute(context MessageContext, mrc20 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Execute(&_GatewayMEVM.TransactOpts, context, mrc20, amount, target, message)
}

// Execute is a paid mutator transaction binding the contract method 0xbcf7f32b.
//
// Solidity: function execute((bytes,address,uint256) context, address mrc20, uint256 amount, address target, bytes message) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) Execute(context MessageContext, mrc20 common.Address, amount *big.Int, target common.Address, message []byte) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Execute(&_GatewayMEVM.TransactOpts, context, mrc20, amount, target, message)
}

// ExecuteAbort is a paid mutator transaction binding the contract method 0x2095dedb.
//
// Solidity: function executeAbort(address target, (bytes,address,uint256,bool,uint256,bytes) abortContext) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) ExecuteAbort(opts *bind.TransactOpts, target common.Address, abortContext AbortContext) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "executeAbort", target, abortContext)
}

// ExecuteAbort is a paid mutator transaction binding the contract method 0x2095dedb.
//
// Solidity: function executeAbort(address target, (bytes,address,uint256,bool,uint256,bytes) abortContext) returns()
func (_GatewayMEVM *GatewayMEVMSession) ExecuteAbort(target common.Address, abortContext AbortContext) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.ExecuteAbort(&_GatewayMEVM.TransactOpts, target, abortContext)
}

// ExecuteAbort is a paid mutator transaction binding the contract method 0x2095dedb.
//
// Solidity: function executeAbort(address target, (bytes,address,uint256,bool,uint256,bytes) abortContext) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) ExecuteAbort(target common.Address, abortContext AbortContext) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.ExecuteAbort(&_GatewayMEVM.TransactOpts, target, abortContext)
}

// ExecuteRevert is a paid mutator transaction binding the contract method 0x184b0793.
//
// Solidity: function executeRevert(address target, (address,address,uint256,bytes) revertContext) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) ExecuteRevert(opts *bind.TransactOpts, target common.Address, revertContext RevertContext) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "executeRevert", target, revertContext)
}

// ExecuteRevert is a paid mutator transaction binding the contract method 0x184b0793.
//
// Solidity: function executeRevert(address target, (address,address,uint256,bytes) revertContext) returns()
func (_GatewayMEVM *GatewayMEVMSession) ExecuteRevert(target common.Address, revertContext RevertContext) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.ExecuteRevert(&_GatewayMEVM.TransactOpts, target, revertContext)
}

// ExecuteRevert is a paid mutator transaction binding the contract method 0x184b0793.
//
// Solidity: function executeRevert(address target, (address,address,uint256,bytes) revertContext) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) ExecuteRevert(target common.Address, revertContext RevertContext) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.ExecuteRevert(&_GatewayMEVM.TransactOpts, target, revertContext)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GatewayMEVM *GatewayMEVMSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.GrantRole(&_GatewayMEVM.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.GrantRole(&_GatewayMEVM.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address museToken_, address admin_) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) Initialize(opts *bind.TransactOpts, museToken_ common.Address, admin_ common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "initialize", museToken_, admin_)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address museToken_, address admin_) returns()
func (_GatewayMEVM *GatewayMEVMSession) Initialize(museToken_ common.Address, admin_ common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Initialize(&_GatewayMEVM.TransactOpts, museToken_, admin_)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address museToken_, address admin_) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) Initialize(museToken_ common.Address, admin_ common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Initialize(&_GatewayMEVM.TransactOpts, museToken_, admin_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GatewayMEVM *GatewayMEVMTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GatewayMEVM *GatewayMEVMSession) Pause() (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Pause(&_GatewayMEVM.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) Pause() (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Pause(&_GatewayMEVM.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_GatewayMEVM *GatewayMEVMSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.RenounceRole(&_GatewayMEVM.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.RenounceRole(&_GatewayMEVM.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GatewayMEVM *GatewayMEVMSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.RevokeRole(&_GatewayMEVM.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.RevokeRole(&_GatewayMEVM.TransactOpts, role, account)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GatewayMEVM *GatewayMEVMTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GatewayMEVM *GatewayMEVMSession) Unpause() (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Unpause(&_GatewayMEVM.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) Unpause() (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Unpause(&_GatewayMEVM.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_GatewayMEVM *GatewayMEVMTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_GatewayMEVM *GatewayMEVMSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.UpgradeToAndCall(&_GatewayMEVM.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.UpgradeToAndCall(&_GatewayMEVM.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x7c0dcb5f.
//
// Solidity: function withdraw(bytes receiver, uint256 amount, address mrc20, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) Withdraw(opts *bind.TransactOpts, receiver []byte, amount *big.Int, mrc20 common.Address, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "withdraw", receiver, amount, mrc20, revertOptions)
}

// Withdraw is a paid mutator transaction binding the contract method 0x7c0dcb5f.
//
// Solidity: function withdraw(bytes receiver, uint256 amount, address mrc20, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayMEVM *GatewayMEVMSession) Withdraw(receiver []byte, amount *big.Int, mrc20 common.Address, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Withdraw(&_GatewayMEVM.TransactOpts, receiver, amount, mrc20, revertOptions)
}

// Withdraw is a paid mutator transaction binding the contract method 0x7c0dcb5f.
//
// Solidity: function withdraw(bytes receiver, uint256 amount, address mrc20, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) Withdraw(receiver []byte, amount *big.Int, mrc20 common.Address, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Withdraw(&_GatewayMEVM.TransactOpts, receiver, amount, mrc20, revertOptions)
}

// WithdrawAndCall0 is a paid mutator transaction binding the contract method 0x7b15118b.
//
// Solidity: function withdrawAndCall(bytes receiver, uint256 amount, address mrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayMEVM *GatewayMEVMTransactor) WithdrawAndCall0(opts *bind.TransactOpts, receiver []byte, amount *big.Int, mrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayMEVM.contract.Transact(opts, "withdrawAndCall0", receiver, amount, mrc20, message, callOptions, revertOptions)
}

// WithdrawAndCall0 is a paid mutator transaction binding the contract method 0x7b15118b.
//
// Solidity: function withdrawAndCall(bytes receiver, uint256 amount, address mrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayMEVM *GatewayMEVMSession) WithdrawAndCall0(receiver []byte, amount *big.Int, mrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.WithdrawAndCall0(&_GatewayMEVM.TransactOpts, receiver, amount, mrc20, message, callOptions, revertOptions)
}

// WithdrawAndCall0 is a paid mutator transaction binding the contract method 0x7b15118b.
//
// Solidity: function withdrawAndCall(bytes receiver, uint256 amount, address mrc20, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions) returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) WithdrawAndCall0(receiver []byte, amount *big.Int, mrc20 common.Address, message []byte, callOptions CallOptions, revertOptions RevertOptions) (*types.Transaction, error) {
	return _GatewayMEVM.Contract.WithdrawAndCall0(&_GatewayMEVM.TransactOpts, receiver, amount, mrc20, message, callOptions, revertOptions)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GatewayMEVM *GatewayMEVMTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GatewayMEVM.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GatewayMEVM *GatewayMEVMSession) Receive() (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Receive(&_GatewayMEVM.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GatewayMEVM *GatewayMEVMTransactorSession) Receive() (*types.Transaction, error) {
	return _GatewayMEVM.Contract.Receive(&_GatewayMEVM.TransactOpts)
}

// GatewayMEVMCalledIterator is returned from FilterCalled and is used to iterate over the raw logs and unpacked data for Called events raised by the GatewayMEVM contract.
type GatewayMEVMCalledIterator struct {
	Event *GatewayMEVMCalled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GatewayMEVMCalledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GatewayMEVMCalled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GatewayMEVMCalled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GatewayMEVMCalledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GatewayMEVMCalledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GatewayMEVMCalled represents a Called event raised by the GatewayMEVM contract.
type GatewayMEVMCalled struct {
	Sender        common.Address
	Mrc20         common.Address
	Receiver      []byte
	Message       []byte
	CallOptions   CallOptions
	RevertOptions RevertOptions
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterCalled is a free log retrieval operation binding the contract event 0x306ee13f48319a123b222c69908e44dcf91abffc20cacc502e3cf5a4ff23e0e4.
//
// Solidity: event Called(address indexed sender, address indexed mrc20, bytes receiver, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions)
func (_GatewayMEVM *GatewayMEVMFilterer) FilterCalled(opts *bind.FilterOpts, sender []common.Address, mrc20 []common.Address) (*GatewayMEVMCalledIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var mrc20Rule []interface{}
	for _, mrc20Item := range mrc20 {
		mrc20Rule = append(mrc20Rule, mrc20Item)
	}

	logs, sub, err := _GatewayMEVM.contract.FilterLogs(opts, "Called", senderRule, mrc20Rule)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMCalledIterator{contract: _GatewayMEVM.contract, event: "Called", logs: logs, sub: sub}, nil
}

// WatchCalled is a free log subscription operation binding the contract event 0x306ee13f48319a123b222c69908e44dcf91abffc20cacc502e3cf5a4ff23e0e4.
//
// Solidity: event Called(address indexed sender, address indexed mrc20, bytes receiver, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions)
func (_GatewayMEVM *GatewayMEVMFilterer) WatchCalled(opts *bind.WatchOpts, sink chan<- *GatewayMEVMCalled, sender []common.Address, mrc20 []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var mrc20Rule []interface{}
	for _, mrc20Item := range mrc20 {
		mrc20Rule = append(mrc20Rule, mrc20Item)
	}

	logs, sub, err := _GatewayMEVM.contract.WatchLogs(opts, "Called", senderRule, mrc20Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GatewayMEVMCalled)
				if err := _GatewayMEVM.contract.UnpackLog(event, "Called", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCalled is a log parse operation binding the contract event 0x306ee13f48319a123b222c69908e44dcf91abffc20cacc502e3cf5a4ff23e0e4.
//
// Solidity: event Called(address indexed sender, address indexed mrc20, bytes receiver, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions)
func (_GatewayMEVM *GatewayMEVMFilterer) ParseCalled(log types.Log) (*GatewayMEVMCalled, error) {
	event := new(GatewayMEVMCalled)
	if err := _GatewayMEVM.contract.UnpackLog(event, "Called", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GatewayMEVMInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GatewayMEVM contract.
type GatewayMEVMInitializedIterator struct {
	Event *GatewayMEVMInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GatewayMEVMInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GatewayMEVMInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GatewayMEVMInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GatewayMEVMInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GatewayMEVMInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GatewayMEVMInitialized represents a Initialized event raised by the GatewayMEVM contract.
type GatewayMEVMInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_GatewayMEVM *GatewayMEVMFilterer) FilterInitialized(opts *bind.FilterOpts) (*GatewayMEVMInitializedIterator, error) {

	logs, sub, err := _GatewayMEVM.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMInitializedIterator{contract: _GatewayMEVM.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_GatewayMEVM *GatewayMEVMFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GatewayMEVMInitialized) (event.Subscription, error) {

	logs, sub, err := _GatewayMEVM.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GatewayMEVMInitialized)
				if err := _GatewayMEVM.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_GatewayMEVM *GatewayMEVMFilterer) ParseInitialized(log types.Log) (*GatewayMEVMInitialized, error) {
	event := new(GatewayMEVMInitialized)
	if err := _GatewayMEVM.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GatewayMEVMPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the GatewayMEVM contract.
type GatewayMEVMPausedIterator struct {
	Event *GatewayMEVMPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GatewayMEVMPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GatewayMEVMPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GatewayMEVMPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GatewayMEVMPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GatewayMEVMPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GatewayMEVMPaused represents a Paused event raised by the GatewayMEVM contract.
type GatewayMEVMPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GatewayMEVM *GatewayMEVMFilterer) FilterPaused(opts *bind.FilterOpts) (*GatewayMEVMPausedIterator, error) {

	logs, sub, err := _GatewayMEVM.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMPausedIterator{contract: _GatewayMEVM.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GatewayMEVM *GatewayMEVMFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *GatewayMEVMPaused) (event.Subscription, error) {

	logs, sub, err := _GatewayMEVM.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GatewayMEVMPaused)
				if err := _GatewayMEVM.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GatewayMEVM *GatewayMEVMFilterer) ParsePaused(log types.Log) (*GatewayMEVMPaused, error) {
	event := new(GatewayMEVMPaused)
	if err := _GatewayMEVM.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GatewayMEVMRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the GatewayMEVM contract.
type GatewayMEVMRoleAdminChangedIterator struct {
	Event *GatewayMEVMRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GatewayMEVMRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GatewayMEVMRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GatewayMEVMRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GatewayMEVMRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GatewayMEVMRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GatewayMEVMRoleAdminChanged represents a RoleAdminChanged event raised by the GatewayMEVM contract.
type GatewayMEVMRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GatewayMEVM *GatewayMEVMFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*GatewayMEVMRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _GatewayMEVM.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMRoleAdminChangedIterator{contract: _GatewayMEVM.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GatewayMEVM *GatewayMEVMFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *GatewayMEVMRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _GatewayMEVM.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GatewayMEVMRoleAdminChanged)
				if err := _GatewayMEVM.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GatewayMEVM *GatewayMEVMFilterer) ParseRoleAdminChanged(log types.Log) (*GatewayMEVMRoleAdminChanged, error) {
	event := new(GatewayMEVMRoleAdminChanged)
	if err := _GatewayMEVM.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GatewayMEVMRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the GatewayMEVM contract.
type GatewayMEVMRoleGrantedIterator struct {
	Event *GatewayMEVMRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GatewayMEVMRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GatewayMEVMRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GatewayMEVMRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GatewayMEVMRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GatewayMEVMRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GatewayMEVMRoleGranted represents a RoleGranted event raised by the GatewayMEVM contract.
type GatewayMEVMRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GatewayMEVM *GatewayMEVMFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GatewayMEVMRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GatewayMEVM.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMRoleGrantedIterator{contract: _GatewayMEVM.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GatewayMEVM *GatewayMEVMFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *GatewayMEVMRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GatewayMEVM.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GatewayMEVMRoleGranted)
				if err := _GatewayMEVM.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GatewayMEVM *GatewayMEVMFilterer) ParseRoleGranted(log types.Log) (*GatewayMEVMRoleGranted, error) {
	event := new(GatewayMEVMRoleGranted)
	if err := _GatewayMEVM.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GatewayMEVMRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the GatewayMEVM contract.
type GatewayMEVMRoleRevokedIterator struct {
	Event *GatewayMEVMRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GatewayMEVMRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GatewayMEVMRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GatewayMEVMRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GatewayMEVMRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GatewayMEVMRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GatewayMEVMRoleRevoked represents a RoleRevoked event raised by the GatewayMEVM contract.
type GatewayMEVMRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GatewayMEVM *GatewayMEVMFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GatewayMEVMRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GatewayMEVM.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMRoleRevokedIterator{contract: _GatewayMEVM.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GatewayMEVM *GatewayMEVMFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *GatewayMEVMRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GatewayMEVM.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GatewayMEVMRoleRevoked)
				if err := _GatewayMEVM.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GatewayMEVM *GatewayMEVMFilterer) ParseRoleRevoked(log types.Log) (*GatewayMEVMRoleRevoked, error) {
	event := new(GatewayMEVMRoleRevoked)
	if err := _GatewayMEVM.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GatewayMEVMUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the GatewayMEVM contract.
type GatewayMEVMUnpausedIterator struct {
	Event *GatewayMEVMUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GatewayMEVMUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GatewayMEVMUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GatewayMEVMUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GatewayMEVMUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GatewayMEVMUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GatewayMEVMUnpaused represents a Unpaused event raised by the GatewayMEVM contract.
type GatewayMEVMUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GatewayMEVM *GatewayMEVMFilterer) FilterUnpaused(opts *bind.FilterOpts) (*GatewayMEVMUnpausedIterator, error) {

	logs, sub, err := _GatewayMEVM.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMUnpausedIterator{contract: _GatewayMEVM.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GatewayMEVM *GatewayMEVMFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *GatewayMEVMUnpaused) (event.Subscription, error) {

	logs, sub, err := _GatewayMEVM.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GatewayMEVMUnpaused)
				if err := _GatewayMEVM.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GatewayMEVM *GatewayMEVMFilterer) ParseUnpaused(log types.Log) (*GatewayMEVMUnpaused, error) {
	event := new(GatewayMEVMUnpaused)
	if err := _GatewayMEVM.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GatewayMEVMUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the GatewayMEVM contract.
type GatewayMEVMUpgradedIterator struct {
	Event *GatewayMEVMUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GatewayMEVMUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GatewayMEVMUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GatewayMEVMUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GatewayMEVMUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GatewayMEVMUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GatewayMEVMUpgraded represents a Upgraded event raised by the GatewayMEVM contract.
type GatewayMEVMUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_GatewayMEVM *GatewayMEVMFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*GatewayMEVMUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _GatewayMEVM.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMUpgradedIterator{contract: _GatewayMEVM.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_GatewayMEVM *GatewayMEVMFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *GatewayMEVMUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _GatewayMEVM.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GatewayMEVMUpgraded)
				if err := _GatewayMEVM.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_GatewayMEVM *GatewayMEVMFilterer) ParseUpgraded(log types.Log) (*GatewayMEVMUpgraded, error) {
	event := new(GatewayMEVMUpgraded)
	if err := _GatewayMEVM.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GatewayMEVMWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the GatewayMEVM contract.
type GatewayMEVMWithdrawnIterator struct {
	Event *GatewayMEVMWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GatewayMEVMWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GatewayMEVMWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GatewayMEVMWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GatewayMEVMWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GatewayMEVMWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GatewayMEVMWithdrawn represents a Withdrawn event raised by the GatewayMEVM contract.
type GatewayMEVMWithdrawn struct {
	Sender          common.Address
	ChainId         *big.Int
	Receiver        []byte
	Mrc20           common.Address
	Value           *big.Int
	Gasfee          *big.Int
	ProtocolFlatFee *big.Int
	Message         []byte
	CallOptions     CallOptions
	RevertOptions   RevertOptions
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x07bf64173efd8f3dfb9e4eb3834bab9d5b85a3d89a1c6425797329de0668502c.
//
// Solidity: event Withdrawn(address indexed sender, uint256 indexed chainId, bytes receiver, address mrc20, uint256 value, uint256 gasfee, uint256 protocolFlatFee, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions)
func (_GatewayMEVM *GatewayMEVMFilterer) FilterWithdrawn(opts *bind.FilterOpts, sender []common.Address, chainId []*big.Int) (*GatewayMEVMWithdrawnIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _GatewayMEVM.contract.FilterLogs(opts, "Withdrawn", senderRule, chainIdRule)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMWithdrawnIterator{contract: _GatewayMEVM.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x07bf64173efd8f3dfb9e4eb3834bab9d5b85a3d89a1c6425797329de0668502c.
//
// Solidity: event Withdrawn(address indexed sender, uint256 indexed chainId, bytes receiver, address mrc20, uint256 value, uint256 gasfee, uint256 protocolFlatFee, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions)
func (_GatewayMEVM *GatewayMEVMFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *GatewayMEVMWithdrawn, sender []common.Address, chainId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _GatewayMEVM.contract.WatchLogs(opts, "Withdrawn", senderRule, chainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GatewayMEVMWithdrawn)
				if err := _GatewayMEVM.contract.UnpackLog(event, "Withdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawn is a log parse operation binding the contract event 0x07bf64173efd8f3dfb9e4eb3834bab9d5b85a3d89a1c6425797329de0668502c.
//
// Solidity: event Withdrawn(address indexed sender, uint256 indexed chainId, bytes receiver, address mrc20, uint256 value, uint256 gasfee, uint256 protocolFlatFee, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions)
func (_GatewayMEVM *GatewayMEVMFilterer) ParseWithdrawn(log types.Log) (*GatewayMEVMWithdrawn, error) {
	event := new(GatewayMEVMWithdrawn)
	if err := _GatewayMEVM.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GatewayMEVMWithdrawnAndCalledIterator is returned from FilterWithdrawnAndCalled and is used to iterate over the raw logs and unpacked data for WithdrawnAndCalled events raised by the GatewayMEVM contract.
type GatewayMEVMWithdrawnAndCalledIterator struct {
	Event *GatewayMEVMWithdrawnAndCalled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GatewayMEVMWithdrawnAndCalledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GatewayMEVMWithdrawnAndCalled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GatewayMEVMWithdrawnAndCalled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GatewayMEVMWithdrawnAndCalledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GatewayMEVMWithdrawnAndCalledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GatewayMEVMWithdrawnAndCalled represents a WithdrawnAndCalled event raised by the GatewayMEVM contract.
type GatewayMEVMWithdrawnAndCalled struct {
	Sender          common.Address
	ChainId         *big.Int
	Receiver        []byte
	Mrc20           common.Address
	Value           *big.Int
	Gasfee          *big.Int
	ProtocolFlatFee *big.Int
	Message         []byte
	CallOptions     CallOptions
	RevertOptions   RevertOptions
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterWithdrawnAndCalled is a free log retrieval operation binding the contract event 0xd90f94752d2b12f364f4a2237ebe1aff24ba6127585376bf4935f6a7be17dd2a.
//
// Solidity: event WithdrawnAndCalled(address indexed sender, uint256 indexed chainId, bytes receiver, address mrc20, uint256 value, uint256 gasfee, uint256 protocolFlatFee, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions)
func (_GatewayMEVM *GatewayMEVMFilterer) FilterWithdrawnAndCalled(opts *bind.FilterOpts, sender []common.Address, chainId []*big.Int) (*GatewayMEVMWithdrawnAndCalledIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _GatewayMEVM.contract.FilterLogs(opts, "WithdrawnAndCalled", senderRule, chainIdRule)
	if err != nil {
		return nil, err
	}
	return &GatewayMEVMWithdrawnAndCalledIterator{contract: _GatewayMEVM.contract, event: "WithdrawnAndCalled", logs: logs, sub: sub}, nil
}

// WatchWithdrawnAndCalled is a free log subscription operation binding the contract event 0xd90f94752d2b12f364f4a2237ebe1aff24ba6127585376bf4935f6a7be17dd2a.
//
// Solidity: event WithdrawnAndCalled(address indexed sender, uint256 indexed chainId, bytes receiver, address mrc20, uint256 value, uint256 gasfee, uint256 protocolFlatFee, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions)
func (_GatewayMEVM *GatewayMEVMFilterer) WatchWithdrawnAndCalled(opts *bind.WatchOpts, sink chan<- *GatewayMEVMWithdrawnAndCalled, sender []common.Address, chainId []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var chainIdRule []interface{}
	for _, chainIdItem := range chainId {
		chainIdRule = append(chainIdRule, chainIdItem)
	}

	logs, sub, err := _GatewayMEVM.contract.WatchLogs(opts, "WithdrawnAndCalled", senderRule, chainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GatewayMEVMWithdrawnAndCalled)
				if err := _GatewayMEVM.contract.UnpackLog(event, "WithdrawnAndCalled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawnAndCalled is a log parse operation binding the contract event 0xd90f94752d2b12f364f4a2237ebe1aff24ba6127585376bf4935f6a7be17dd2a.
//
// Solidity: event WithdrawnAndCalled(address indexed sender, uint256 indexed chainId, bytes receiver, address mrc20, uint256 value, uint256 gasfee, uint256 protocolFlatFee, bytes message, (uint256,bool) callOptions, (address,bool,address,bytes,uint256) revertOptions)
func (_GatewayMEVM *GatewayMEVMFilterer) ParseWithdrawnAndCalled(log types.Log) (*GatewayMEVMWithdrawnAndCalled, error) {
	event := new(GatewayMEVMWithdrawnAndCalled)
	if err := _GatewayMEVM.contract.UnpackLog(event, "WithdrawnAndCalled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
