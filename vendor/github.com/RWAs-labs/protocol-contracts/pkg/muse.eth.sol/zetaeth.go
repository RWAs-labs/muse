// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package muse

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

// MuseEthMetaData contains all meta data concerning the MuseEth contract.
var MuseEthMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610e47380380610e4783398101604081905261002f9161022c565b604051806040016040528060048152602001635a65746160e01b815250604051806040016040528060048152602001635a45544160e01b81525081600390816100789190610305565b5060046100858282610305565b5050506100b78261009a6100be60201b60201c565b6100a89060ff16600a6104c2565b6100b290846104d5565b6100c3565b50506104ff565b601290565b6001600160a01b0382166100f25760405163ec442f0560e01b8152600060048201526024015b60405180910390fd5b6100fe60008383610102565b5050565b6001600160a01b03831661012d57806002600082825461012291906104ec565b9091555061019f9050565b6001600160a01b038316600090815260208190526040902054818110156101805760405163391434e360e21b81526001600160a01b038516600482015260248101829052604481018390526064016100e9565b6001600160a01b03841660009081526020819052604090209082900390555b6001600160a01b0382166101bb576002805482900390556101da565b6001600160a01b03821660009081526020819052604090208054820190555b816001600160a01b0316836001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161021f91815260200190565b60405180910390a3505050565b6000806040838503121561023f57600080fd5b82516001600160a01b038116811461025657600080fd5b6020939093015192949293505050565b634e487b7160e01b600052604160045260246000fd5b600181811c9082168061029057607f821691505b6020821081036102b057634e487b7160e01b600052602260045260246000fd5b50919050565b601f82111561030057806000526020600020601f840160051c810160208510156102dd5750805b601f840160051c820191505b818110156102fd57600081556001016102e9565b50505b505050565b81516001600160401b0381111561031e5761031e610266565b6103328161032c845461027c565b846102b6565b6020601f821160018114610366576000831561034e5750848201515b600019600385901b1c1916600184901b1784556102fd565b600084815260208120601f198516915b828110156103965787850151825560209485019460019092019101610376565b50848210156103b45786840151600019600387901b60f8161c191681555b50505050600190811b01905550565b634e487b7160e01b600052601160045260246000fd5b6001815b6001841115610414578085048111156103f8576103f86103c3565b600184161561040657908102905b60019390931c9280026103dd565b935093915050565b60008261042b575060016104bc565b81610438575060006104bc565b816001811461044e576002811461045857610474565b60019150506104bc565b60ff841115610469576104696103c3565b50506001821b6104bc565b5060208310610133831016604e8410600b8410161715610497575081810a6104bc565b6104a460001984846103d9565b80600019048211156104b8576104b86103c3565b0290505b92915050565b60006104ce838361041c565b9392505050565b80820281158282048414176104bc576104bc6103c3565b808201808211156104bc576104bc6103c3565b6109398061050e6000396000f3fe608060405234801561001057600080fd5b50600436106100a35760003560e01c8063313ce5671161007657806395d89b411161005b57806395d89b4114610153578063a9059cbb1461015b578063dd62ed3e1461016e57600080fd5b8063313ce5671461010e57806370a082311461011d57600080fd5b806306fdde03146100a8578063095ea7b3146100c657806318160ddd146100e957806323b872dd146100fb575b600080fd5b6100b06101b4565b6040516100bd9190610725565b60405180910390f35b6100d96100d43660046107ba565b610246565b60405190151581526020016100bd565b6002545b6040519081526020016100bd565b6100d96101093660046107e4565b610260565b604051601281526020016100bd565b6100ed61012b366004610821565b73ffffffffffffffffffffffffffffffffffffffff1660009081526020819052604090205490565b6100b0610284565b6100d96101693660046107ba565b610293565b6100ed61017c366004610843565b73ffffffffffffffffffffffffffffffffffffffff918216600090815260016020908152604080832093909416825291909152205490565b6060600380546101c390610876565b80601f01602080910402602001604051908101604052809291908181526020018280546101ef90610876565b801561023c5780601f106102115761010080835404028352916020019161023c565b820191906000526020600020905b81548152906001019060200180831161021f57829003601f168201915b5050505050905090565b6000336102548185856102a1565b60019150505b92915050565b60003361026e8582856102b3565b610279858585610387565b506001949350505050565b6060600480546101c390610876565b600033610254818585610387565b6102ae8383836001610432565b505050565b73ffffffffffffffffffffffffffffffffffffffff8381166000908152600160209081526040808320938616835292905220547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81146103815781811015610372576040517ffb8f41b200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8416600482015260248101829052604481018390526064015b60405180910390fd5b61038184848484036000610432565b50505050565b73ffffffffffffffffffffffffffffffffffffffff83166103d7576040517f96c6fd1e00000000000000000000000000000000000000000000000000000000815260006004820152602401610369565b73ffffffffffffffffffffffffffffffffffffffff8216610427576040517fec442f0500000000000000000000000000000000000000000000000000000000815260006004820152602401610369565b6102ae83838361057a565b73ffffffffffffffffffffffffffffffffffffffff8416610482576040517fe602df0500000000000000000000000000000000000000000000000000000000815260006004820152602401610369565b73ffffffffffffffffffffffffffffffffffffffff83166104d2576040517f94280d6200000000000000000000000000000000000000000000000000000000815260006004820152602401610369565b73ffffffffffffffffffffffffffffffffffffffff80851660009081526001602090815260408083209387168352929052208290558015610381578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258460405161056c91815260200190565b60405180910390a350505050565b73ffffffffffffffffffffffffffffffffffffffff83166105b25780600260008282546105a791906108c9565b909155506106649050565b73ffffffffffffffffffffffffffffffffffffffff831660009081526020819052604090205481811015610638576040517fe450d38c00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff851660048201526024810182905260448101839052606401610369565b73ffffffffffffffffffffffffffffffffffffffff841660009081526020819052604090209082900390555b73ffffffffffffffffffffffffffffffffffffffff821661068d576002805482900390556106b9565b73ffffffffffffffffffffffffffffffffffffffff821660009081526020819052604090208054820190555b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161071891815260200190565b60405180910390a3505050565b602081526000825180602084015260005b818110156107535760208186018101516040868401015201610736565b5060006040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b803573ffffffffffffffffffffffffffffffffffffffff811681146107b557600080fd5b919050565b600080604083850312156107cd57600080fd5b6107d683610791565b946020939093013593505050565b6000806000606084860312156107f957600080fd5b61080284610791565b925061081060208501610791565b929592945050506040919091013590565b60006020828403121561083357600080fd5b61083c82610791565b9392505050565b6000806040838503121561085657600080fd5b61085f83610791565b915061086d60208401610791565b90509250929050565b600181811c9082168061088a57607f821691505b6020821081036108c3577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b8082018082111561025a577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fdfea2646970667358221220855466f44bede3a825ee0445c5e828ad627f363663f502c82442076d31cd264d64736f6c634300081a0033",
}

// MuseEthABI is the input ABI used to generate the binding from.
// Deprecated: Use MuseEthMetaData.ABI instead.
var MuseEthABI = MuseEthMetaData.ABI

// MuseEthBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MuseEthMetaData.Bin instead.
var MuseEthBin = MuseEthMetaData.Bin

// DeployMuseEth deploys a new Ethereum contract, binding an instance of MuseEth to it.
func DeployMuseEth(auth *bind.TransactOpts, backend bind.ContractBackend, creator common.Address, initialSupply *big.Int) (common.Address, *types.Transaction, *MuseEth, error) {
	parsed, err := MuseEthMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MuseEthBin), backend, creator, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MuseEth{MuseEthCaller: MuseEthCaller{contract: contract}, MuseEthTransactor: MuseEthTransactor{contract: contract}, MuseEthFilterer: MuseEthFilterer{contract: contract}}, nil
}

// MuseEth is an auto generated Go binding around an Ethereum contract.
type MuseEth struct {
	MuseEthCaller     // Read-only binding to the contract
	MuseEthTransactor // Write-only binding to the contract
	MuseEthFilterer   // Log filterer for contract events
}

// MuseEthCaller is an auto generated read-only Go binding around an Ethereum contract.
type MuseEthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseEthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MuseEthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseEthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MuseEthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseEthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MuseEthSession struct {
	Contract     *MuseEth          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MuseEthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MuseEthCallerSession struct {
	Contract *MuseEthCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// MuseEthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MuseEthTransactorSession struct {
	Contract     *MuseEthTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MuseEthRaw is an auto generated low-level Go binding around an Ethereum contract.
type MuseEthRaw struct {
	Contract *MuseEth // Generic contract binding to access the raw methods on
}

// MuseEthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MuseEthCallerRaw struct {
	Contract *MuseEthCaller // Generic read-only contract binding to access the raw methods on
}

// MuseEthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MuseEthTransactorRaw struct {
	Contract *MuseEthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMuseEth creates a new instance of MuseEth, bound to a specific deployed contract.
func NewMuseEth(address common.Address, backend bind.ContractBackend) (*MuseEth, error) {
	contract, err := bindMuseEth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MuseEth{MuseEthCaller: MuseEthCaller{contract: contract}, MuseEthTransactor: MuseEthTransactor{contract: contract}, MuseEthFilterer: MuseEthFilterer{contract: contract}}, nil
}

// NewMuseEthCaller creates a new read-only instance of MuseEth, bound to a specific deployed contract.
func NewMuseEthCaller(address common.Address, caller bind.ContractCaller) (*MuseEthCaller, error) {
	contract, err := bindMuseEth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MuseEthCaller{contract: contract}, nil
}

// NewMuseEthTransactor creates a new write-only instance of MuseEth, bound to a specific deployed contract.
func NewMuseEthTransactor(address common.Address, transactor bind.ContractTransactor) (*MuseEthTransactor, error) {
	contract, err := bindMuseEth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MuseEthTransactor{contract: contract}, nil
}

// NewMuseEthFilterer creates a new log filterer instance of MuseEth, bound to a specific deployed contract.
func NewMuseEthFilterer(address common.Address, filterer bind.ContractFilterer) (*MuseEthFilterer, error) {
	contract, err := bindMuseEth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MuseEthFilterer{contract: contract}, nil
}

// bindMuseEth binds a generic wrapper to an already deployed contract.
func bindMuseEth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MuseEthMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MuseEth *MuseEthRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MuseEth.Contract.MuseEthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MuseEth *MuseEthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseEth.Contract.MuseEthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MuseEth *MuseEthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MuseEth.Contract.MuseEthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MuseEth *MuseEthCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MuseEth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MuseEth *MuseEthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseEth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MuseEth *MuseEthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MuseEth.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MuseEth *MuseEthCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MuseEth.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MuseEth *MuseEthSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MuseEth.Contract.Allowance(&_MuseEth.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_MuseEth *MuseEthCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _MuseEth.Contract.Allowance(&_MuseEth.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MuseEth *MuseEthCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MuseEth.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MuseEth *MuseEthSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _MuseEth.Contract.BalanceOf(&_MuseEth.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_MuseEth *MuseEthCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _MuseEth.Contract.BalanceOf(&_MuseEth.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MuseEth *MuseEthCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _MuseEth.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MuseEth *MuseEthSession) Decimals() (uint8, error) {
	return _MuseEth.Contract.Decimals(&_MuseEth.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_MuseEth *MuseEthCallerSession) Decimals() (uint8, error) {
	return _MuseEth.Contract.Decimals(&_MuseEth.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MuseEth *MuseEthCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MuseEth.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MuseEth *MuseEthSession) Name() (string, error) {
	return _MuseEth.Contract.Name(&_MuseEth.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_MuseEth *MuseEthCallerSession) Name() (string, error) {
	return _MuseEth.Contract.Name(&_MuseEth.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MuseEth *MuseEthCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MuseEth.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MuseEth *MuseEthSession) Symbol() (string, error) {
	return _MuseEth.Contract.Symbol(&_MuseEth.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_MuseEth *MuseEthCallerSession) Symbol() (string, error) {
	return _MuseEth.Contract.Symbol(&_MuseEth.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MuseEth *MuseEthCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MuseEth.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MuseEth *MuseEthSession) TotalSupply() (*big.Int, error) {
	return _MuseEth.Contract.TotalSupply(&_MuseEth.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_MuseEth *MuseEthCallerSession) TotalSupply() (*big.Int, error) {
	return _MuseEth.Contract.TotalSupply(&_MuseEth.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MuseEth *MuseEthTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MuseEth.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MuseEth *MuseEthSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MuseEth.Contract.Approve(&_MuseEth.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_MuseEth *MuseEthTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _MuseEth.Contract.Approve(&_MuseEth.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MuseEth *MuseEthTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MuseEth.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MuseEth *MuseEthSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MuseEth.Contract.Transfer(&_MuseEth.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_MuseEth *MuseEthTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MuseEth.Contract.Transfer(&_MuseEth.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MuseEth *MuseEthTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MuseEth.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MuseEth *MuseEthSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MuseEth.Contract.TransferFrom(&_MuseEth.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_MuseEth *MuseEthTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _MuseEth.Contract.TransferFrom(&_MuseEth.TransactOpts, from, to, value)
}

// MuseEthApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the MuseEth contract.
type MuseEthApprovalIterator struct {
	Event *MuseEthApproval // Event containing the contract specifics and raw log

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
func (it *MuseEthApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseEthApproval)
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
		it.Event = new(MuseEthApproval)
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
func (it *MuseEthApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseEthApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseEthApproval represents a Approval event raised by the MuseEth contract.
type MuseEthApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MuseEth *MuseEthFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*MuseEthApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MuseEth.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &MuseEthApprovalIterator{contract: _MuseEth.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MuseEth *MuseEthFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *MuseEthApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _MuseEth.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseEthApproval)
				if err := _MuseEth.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_MuseEth *MuseEthFilterer) ParseApproval(log types.Log) (*MuseEthApproval, error) {
	event := new(MuseEthApproval)
	if err := _MuseEth.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseEthTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the MuseEth contract.
type MuseEthTransferIterator struct {
	Event *MuseEthTransfer // Event containing the contract specifics and raw log

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
func (it *MuseEthTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseEthTransfer)
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
		it.Event = new(MuseEthTransfer)
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
func (it *MuseEthTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseEthTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseEthTransfer represents a Transfer event raised by the MuseEth contract.
type MuseEthTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MuseEth *MuseEthFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*MuseEthTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MuseEth.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MuseEthTransferIterator{contract: _MuseEth.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MuseEth *MuseEthFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *MuseEthTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MuseEth.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseEthTransfer)
				if err := _MuseEth.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_MuseEth *MuseEthFilterer) ParseTransfer(log types.Log) (*MuseEthTransfer, error) {
	event := new(MuseEthTransfer)
	if err := _MuseEth.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
