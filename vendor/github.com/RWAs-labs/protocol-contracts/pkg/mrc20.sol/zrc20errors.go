// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mrc20

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

// MRC20ErrorsMetaData contains all meta data concerning the MRC20Errors contract.
var MRC20ErrorsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"error\",\"name\":\"CallerIsNotFungibleModule\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasFeeTransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSender\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LowAllowance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LowBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroGasCoin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroGasPrice\",\"inputs\":[]}]",
}

// MRC20ErrorsABI is the input ABI used to generate the binding from.
// Deprecated: Use MRC20ErrorsMetaData.ABI instead.
var MRC20ErrorsABI = MRC20ErrorsMetaData.ABI

// MRC20Errors is an auto generated Go binding around an Ethereum contract.
type MRC20Errors struct {
	MRC20ErrorsCaller     // Read-only binding to the contract
	MRC20ErrorsTransactor // Write-only binding to the contract
	MRC20ErrorsFilterer   // Log filterer for contract events
}

// MRC20ErrorsCaller is an auto generated read-only Go binding around an Ethereum contract.
type MRC20ErrorsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MRC20ErrorsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MRC20ErrorsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MRC20ErrorsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MRC20ErrorsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MRC20ErrorsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MRC20ErrorsSession struct {
	Contract     *MRC20Errors      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MRC20ErrorsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MRC20ErrorsCallerSession struct {
	Contract *MRC20ErrorsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// MRC20ErrorsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MRC20ErrorsTransactorSession struct {
	Contract     *MRC20ErrorsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// MRC20ErrorsRaw is an auto generated low-level Go binding around an Ethereum contract.
type MRC20ErrorsRaw struct {
	Contract *MRC20Errors // Generic contract binding to access the raw methods on
}

// MRC20ErrorsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MRC20ErrorsCallerRaw struct {
	Contract *MRC20ErrorsCaller // Generic read-only contract binding to access the raw methods on
}

// MRC20ErrorsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MRC20ErrorsTransactorRaw struct {
	Contract *MRC20ErrorsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMRC20Errors creates a new instance of MRC20Errors, bound to a specific deployed contract.
func NewMRC20Errors(address common.Address, backend bind.ContractBackend) (*MRC20Errors, error) {
	contract, err := bindMRC20Errors(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MRC20Errors{MRC20ErrorsCaller: MRC20ErrorsCaller{contract: contract}, MRC20ErrorsTransactor: MRC20ErrorsTransactor{contract: contract}, MRC20ErrorsFilterer: MRC20ErrorsFilterer{contract: contract}}, nil
}

// NewMRC20ErrorsCaller creates a new read-only instance of MRC20Errors, bound to a specific deployed contract.
func NewMRC20ErrorsCaller(address common.Address, caller bind.ContractCaller) (*MRC20ErrorsCaller, error) {
	contract, err := bindMRC20Errors(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MRC20ErrorsCaller{contract: contract}, nil
}

// NewMRC20ErrorsTransactor creates a new write-only instance of MRC20Errors, bound to a specific deployed contract.
func NewMRC20ErrorsTransactor(address common.Address, transactor bind.ContractTransactor) (*MRC20ErrorsTransactor, error) {
	contract, err := bindMRC20Errors(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MRC20ErrorsTransactor{contract: contract}, nil
}

// NewMRC20ErrorsFilterer creates a new log filterer instance of MRC20Errors, bound to a specific deployed contract.
func NewMRC20ErrorsFilterer(address common.Address, filterer bind.ContractFilterer) (*MRC20ErrorsFilterer, error) {
	contract, err := bindMRC20Errors(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MRC20ErrorsFilterer{contract: contract}, nil
}

// bindMRC20Errors binds a generic wrapper to an already deployed contract.
func bindMRC20Errors(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MRC20ErrorsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MRC20Errors *MRC20ErrorsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MRC20Errors.Contract.MRC20ErrorsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MRC20Errors *MRC20ErrorsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MRC20Errors.Contract.MRC20ErrorsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MRC20Errors *MRC20ErrorsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MRC20Errors.Contract.MRC20ErrorsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MRC20Errors *MRC20ErrorsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MRC20Errors.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MRC20Errors *MRC20ErrorsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MRC20Errors.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MRC20Errors *MRC20ErrorsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MRC20Errors.Contract.contract.Transact(opts, method, params...)
}
