// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package museconnectormevm

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

// MuseInterfacesMuseMessage is an auto generated low-level Go binding around an user-defined struct.
type MuseInterfacesMuseMessage struct {
	MuseTxSenderAddress []byte
	SourceChainId       *big.Int
	DestinationAddress  common.Address
	MuseValue           *big.Int
	Message             []byte
}

// MuseInterfacesMuseRevert is an auto generated low-level Go binding around an user-defined struct.
type MuseInterfacesMuseRevert struct {
	MuseTxSenderAddress common.Address
	SourceChainId       *big.Int
	DestinationAddress  []byte
	DestinationChainId  *big.Int
	RemainingMuseValue  *big.Int
	Message             []byte
}

// MuseReceiverMetaData contains all meta data concerning the MuseReceiver contract.
var MuseReceiverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"onMuseMessage\",\"inputs\":[{\"name\":\"museMessage\",\"type\":\"tuple\",\"internalType\":\"structMuseInterfaces.MuseMessage\",\"components\":[{\"name\":\"museTxSenderAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"museValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onMuseRevert\",\"inputs\":[{\"name\":\"museRevert\",\"type\":\"tuple\",\"internalType\":\"structMuseInterfaces.MuseRevert\",\"components\":[{\"name\":\"museTxSenderAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingMuseValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"}]",
}

// MuseReceiverABI is the input ABI used to generate the binding from.
// Deprecated: Use MuseReceiverMetaData.ABI instead.
var MuseReceiverABI = MuseReceiverMetaData.ABI

// MuseReceiver is an auto generated Go binding around an Ethereum contract.
type MuseReceiver struct {
	MuseReceiverCaller     // Read-only binding to the contract
	MuseReceiverTransactor // Write-only binding to the contract
	MuseReceiverFilterer   // Log filterer for contract events
}

// MuseReceiverCaller is an auto generated read-only Go binding around an Ethereum contract.
type MuseReceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseReceiverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MuseReceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseReceiverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MuseReceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseReceiverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MuseReceiverSession struct {
	Contract     *MuseReceiver     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MuseReceiverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MuseReceiverCallerSession struct {
	Contract *MuseReceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// MuseReceiverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MuseReceiverTransactorSession struct {
	Contract     *MuseReceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// MuseReceiverRaw is an auto generated low-level Go binding around an Ethereum contract.
type MuseReceiverRaw struct {
	Contract *MuseReceiver // Generic contract binding to access the raw methods on
}

// MuseReceiverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MuseReceiverCallerRaw struct {
	Contract *MuseReceiverCaller // Generic read-only contract binding to access the raw methods on
}

// MuseReceiverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MuseReceiverTransactorRaw struct {
	Contract *MuseReceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMuseReceiver creates a new instance of MuseReceiver, bound to a specific deployed contract.
func NewMuseReceiver(address common.Address, backend bind.ContractBackend) (*MuseReceiver, error) {
	contract, err := bindMuseReceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MuseReceiver{MuseReceiverCaller: MuseReceiverCaller{contract: contract}, MuseReceiverTransactor: MuseReceiverTransactor{contract: contract}, MuseReceiverFilterer: MuseReceiverFilterer{contract: contract}}, nil
}

// NewMuseReceiverCaller creates a new read-only instance of MuseReceiver, bound to a specific deployed contract.
func NewMuseReceiverCaller(address common.Address, caller bind.ContractCaller) (*MuseReceiverCaller, error) {
	contract, err := bindMuseReceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MuseReceiverCaller{contract: contract}, nil
}

// NewMuseReceiverTransactor creates a new write-only instance of MuseReceiver, bound to a specific deployed contract.
func NewMuseReceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*MuseReceiverTransactor, error) {
	contract, err := bindMuseReceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MuseReceiverTransactor{contract: contract}, nil
}

// NewMuseReceiverFilterer creates a new log filterer instance of MuseReceiver, bound to a specific deployed contract.
func NewMuseReceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*MuseReceiverFilterer, error) {
	contract, err := bindMuseReceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MuseReceiverFilterer{contract: contract}, nil
}

// bindMuseReceiver binds a generic wrapper to an already deployed contract.
func bindMuseReceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MuseReceiverMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MuseReceiver *MuseReceiverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MuseReceiver.Contract.MuseReceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MuseReceiver *MuseReceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseReceiver.Contract.MuseReceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MuseReceiver *MuseReceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MuseReceiver.Contract.MuseReceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MuseReceiver *MuseReceiverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MuseReceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MuseReceiver *MuseReceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseReceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MuseReceiver *MuseReceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MuseReceiver.Contract.contract.Transact(opts, method, params...)
}

// OnMuseMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onMuseMessage((bytes,uint256,address,uint256,bytes) museMessage) returns()
func (_MuseReceiver *MuseReceiverTransactor) OnMuseMessage(opts *bind.TransactOpts, museMessage MuseInterfacesMuseMessage) (*types.Transaction, error) {
	return _MuseReceiver.contract.Transact(opts, "onMuseMessage", museMessage)
}

// OnMuseMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onMuseMessage((bytes,uint256,address,uint256,bytes) museMessage) returns()
func (_MuseReceiver *MuseReceiverSession) OnMuseMessage(museMessage MuseInterfacesMuseMessage) (*types.Transaction, error) {
	return _MuseReceiver.Contract.OnMuseMessage(&_MuseReceiver.TransactOpts, museMessage)
}

// OnMuseMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onMuseMessage((bytes,uint256,address,uint256,bytes) museMessage) returns()
func (_MuseReceiver *MuseReceiverTransactorSession) OnMuseMessage(museMessage MuseInterfacesMuseMessage) (*types.Transaction, error) {
	return _MuseReceiver.Contract.OnMuseMessage(&_MuseReceiver.TransactOpts, museMessage)
}

// OnMuseRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onMuseRevert((address,uint256,bytes,uint256,uint256,bytes) museRevert) returns()
func (_MuseReceiver *MuseReceiverTransactor) OnMuseRevert(opts *bind.TransactOpts, museRevert MuseInterfacesMuseRevert) (*types.Transaction, error) {
	return _MuseReceiver.contract.Transact(opts, "onMuseRevert", museRevert)
}

// OnMuseRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onMuseRevert((address,uint256,bytes,uint256,uint256,bytes) museRevert) returns()
func (_MuseReceiver *MuseReceiverSession) OnMuseRevert(museRevert MuseInterfacesMuseRevert) (*types.Transaction, error) {
	return _MuseReceiver.Contract.OnMuseRevert(&_MuseReceiver.TransactOpts, museRevert)
}

// OnMuseRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onMuseRevert((address,uint256,bytes,uint256,uint256,bytes) museRevert) returns()
func (_MuseReceiver *MuseReceiverTransactorSession) OnMuseRevert(museRevert MuseInterfacesMuseRevert) (*types.Transaction, error) {
	return _MuseReceiver.Contract.OnMuseRevert(&_MuseReceiver.TransactOpts, museRevert)
}
