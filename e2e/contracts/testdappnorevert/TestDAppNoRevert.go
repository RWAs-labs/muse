// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testdappnorevert

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

// TestDAppNoRevertMetaData contains all meta data concerning the TestDAppNoRevert contract.
var TestDAppNoRevertMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_connector\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_museToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ErrorTransferringMuse\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMessageType\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"HelloWorldEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"RevertedHelloWorldEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"HELLO_WORLD_MESSAGE_TYPE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"connector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"museTxSenderAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"destinationAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"museValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"internalType\":\"structMuseInterfaces.MuseMessage\",\"name\":\"museMessage\",\"type\":\"tuple\"}],\"name\":\"onMuseMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"destinationAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"doRevert\",\"type\":\"bool\"}],\"name\":\"sendHelloWorld\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"muse\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610cf4380380610cf48339818101604052810190610032919061011d565b816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505061015d565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100ea826100bf565b9050919050565b6100fa816100df565b811461010557600080fd5b50565b600081519050610117816100f1565b92915050565b60008060408385031215610134576101336100ba565b5b600061014285828601610108565b925050602061015385828601610108565b9150509250929050565b610b888061016c6000396000f3fe60806040526004361061004a5760003560e01c80633749c51a1461004f5780637caca3041461007857806383f3084f146100945780638ac44a3f146100bf578063e8f9cb3a146100ea575b600080fd5b34801561005b57600080fd5b5061007660048036038101906100719190610533565b610115565b005b610092600480360381019061008d9190610648565b6101ae565b005b3480156100a057600080fd5b506100a9610497565b6040516100b691906106be565b60405180910390f35b3480156100cb57600080fd5b506100d46104bb565b6040516100e191906106f2565b60405180910390f35b3480156100f657600080fd5b506100ff6104df565b60405161010c91906106be565b60405180910390f35b6000818060800190610127919061071c565b81019061013491906107ab565b915050600015158115151461017e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161017590610848565b60405180910390fd5b7f3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d160405160405180910390a15050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663095ea7b360008054906101000a900473ffffffffffffffffffffffffffffffffffffffff16856040518363ffffffff1660e01b815260040161022d929190610877565b6020604051808303816000875af115801561024c573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061027091906108b5565b90506000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330876040518463ffffffff1660e01b81526004016102d3939291906108e2565b6020604051808303816000875af11580156102f2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061031691906108b5565b90508180156103225750805b610358576040517f2bd0ba5000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ec0269016040518060c00160405280888152602001896040516020016103b69190610961565b60405160208183030381529060405281526020016203d09081526020017f6e0182194bb1deba01849afd3e035a0b70ce7cb069e482ee663519c76cf569b48760405160200161040692919061098b565b604051602081830303815290604052815260200187815260200160405160200161042f906109da565b6040516020818303038152906040528152506040518263ffffffff1660e01b815260040161045d9190610b30565b600060405180830381600087803b15801561047757600080fd5b505af115801561048b573d6000803e3d6000fd5b50505050505050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b7f6e0182194bb1deba01849afd3e035a0b70ce7cb069e482ee663519c76cf569b481565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600080fd5b600080fd5b600080fd5b600060a0828403121561052a5761052961050f565b5b81905092915050565b60006020828403121561054957610548610505565b5b600082013567ffffffffffffffff8111156105675761056661050a565b5b61057384828501610514565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006105a78261057c565b9050919050565b6105b78161059c565b81146105c257600080fd5b50565b6000813590506105d4816105ae565b92915050565b6000819050919050565b6105ed816105da565b81146105f857600080fd5b50565b60008135905061060a816105e4565b92915050565b60008115159050919050565b61062581610610565b811461063057600080fd5b50565b6000813590506106428161061c565b92915050565b6000806000806080858703121561066257610661610505565b5b6000610670878288016105c5565b9450506020610681878288016105fb565b9350506040610692878288016105fb565b92505060606106a387828801610633565b91505092959194509250565b6106b88161059c565b82525050565b60006020820190506106d360008301846106af565b92915050565b6000819050919050565b6106ec816106d9565b82525050565b600060208201905061070760008301846106e3565b92915050565b600080fd5b600080fd5b600080fd5b600080833560016020038436030381126107395761073861070d565b5b80840192508235915067ffffffffffffffff82111561075b5761075a610712565b5b60208301925060018202360383131561077757610776610717565b5b509250929050565b610788816106d9565b811461079357600080fd5b50565b6000813590506107a58161077f565b92915050565b600080604083850312156107c2576107c1610505565b5b60006107d085828601610796565b92505060206107e185828601610633565b9150509250929050565b600082825260208201905092915050565b7f6d65737361676520736179732072657665727400000000000000000000000000600082015250565b60006108326013836107eb565b915061083d826107fc565b602082019050919050565b6000602082019050818103600083015261086181610825565b9050919050565b610871816105da565b82525050565b600060408201905061088c60008301856106af565b6108996020830184610868565b9392505050565b6000815190506108af8161061c565b92915050565b6000602082840312156108cb576108ca610505565b5b60006108d9848285016108a0565b91505092915050565b60006060820190506108f760008301866106af565b61090460208301856106af565b6109116040830184610868565b949350505050565b60008160601b9050919050565b600061093182610919565b9050919050565b600061094382610926565b9050919050565b61095b6109568261059c565b610938565b82525050565b600061096d828461094a565b60148201915081905092915050565b61098581610610565b82525050565b60006040820190506109a060008301856106e3565b6109ad602083018461097c565b9392505050565b50565b60006109c46000836107eb565b91506109cf826109b4565b600082019050919050565b600060208201905081810360008301526109f3816109b7565b9050919050565b610a03816105da565b82525050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610a43578082015181840152602081019050610a28565b60008484015250505050565b6000601f19601f8301169050919050565b6000610a6b82610a09565b610a758185610a14565b9350610a85818560208601610a25565b610a8e81610a4f565b840191505092915050565b600060c083016000830151610ab160008601826109fa565b5060208301518482036020860152610ac98282610a60565b9150506040830151610ade60408601826109fa565b5060608301518482036060860152610af68282610a60565b9150506080830151610b0b60808601826109fa565b5060a083015184820360a0860152610b238282610a60565b9150508091505092915050565b60006020820190508181036000830152610b4a8184610a99565b90509291505056fea264697066735822122058690b26f0645b6f20806288a08764ea93439a663b2bff50b96d229cba0652c964736f6c63430008170033",
}

// TestDAppNoRevertABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDAppNoRevertMetaData.ABI instead.
var TestDAppNoRevertABI = TestDAppNoRevertMetaData.ABI

// TestDAppNoRevertBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestDAppNoRevertMetaData.Bin instead.
var TestDAppNoRevertBin = TestDAppNoRevertMetaData.Bin

// DeployTestDAppNoRevert deploys a new Ethereum contract, binding an instance of TestDAppNoRevert to it.
func DeployTestDAppNoRevert(auth *bind.TransactOpts, backend bind.ContractBackend, _connector common.Address, _museToken common.Address) (common.Address, *types.Transaction, *TestDAppNoRevert, error) {
	parsed, err := TestDAppNoRevertMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDAppNoRevertBin), backend, _connector, _museToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestDAppNoRevert{TestDAppNoRevertCaller: TestDAppNoRevertCaller{contract: contract}, TestDAppNoRevertTransactor: TestDAppNoRevertTransactor{contract: contract}, TestDAppNoRevertFilterer: TestDAppNoRevertFilterer{contract: contract}}, nil
}

// TestDAppNoRevert is an auto generated Go binding around an Ethereum contract.
type TestDAppNoRevert struct {
	TestDAppNoRevertCaller     // Read-only binding to the contract
	TestDAppNoRevertTransactor // Write-only binding to the contract
	TestDAppNoRevertFilterer   // Log filterer for contract events
}

// TestDAppNoRevertCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestDAppNoRevertCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppNoRevertTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestDAppNoRevertTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppNoRevertFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestDAppNoRevertFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppNoRevertSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestDAppNoRevertSession struct {
	Contract     *TestDAppNoRevert // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestDAppNoRevertCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestDAppNoRevertCallerSession struct {
	Contract *TestDAppNoRevertCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// TestDAppNoRevertTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestDAppNoRevertTransactorSession struct {
	Contract     *TestDAppNoRevertTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// TestDAppNoRevertRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestDAppNoRevertRaw struct {
	Contract *TestDAppNoRevert // Generic contract binding to access the raw methods on
}

// TestDAppNoRevertCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestDAppNoRevertCallerRaw struct {
	Contract *TestDAppNoRevertCaller // Generic read-only contract binding to access the raw methods on
}

// TestDAppNoRevertTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestDAppNoRevertTransactorRaw struct {
	Contract *TestDAppNoRevertTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestDAppNoRevert creates a new instance of TestDAppNoRevert, bound to a specific deployed contract.
func NewTestDAppNoRevert(address common.Address, backend bind.ContractBackend) (*TestDAppNoRevert, error) {
	contract, err := bindTestDAppNoRevert(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevert{TestDAppNoRevertCaller: TestDAppNoRevertCaller{contract: contract}, TestDAppNoRevertTransactor: TestDAppNoRevertTransactor{contract: contract}, TestDAppNoRevertFilterer: TestDAppNoRevertFilterer{contract: contract}}, nil
}

// NewTestDAppNoRevertCaller creates a new read-only instance of TestDAppNoRevert, bound to a specific deployed contract.
func NewTestDAppNoRevertCaller(address common.Address, caller bind.ContractCaller) (*TestDAppNoRevertCaller, error) {
	contract, err := bindTestDAppNoRevert(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevertCaller{contract: contract}, nil
}

// NewTestDAppNoRevertTransactor creates a new write-only instance of TestDAppNoRevert, bound to a specific deployed contract.
func NewTestDAppNoRevertTransactor(address common.Address, transactor bind.ContractTransactor) (*TestDAppNoRevertTransactor, error) {
	contract, err := bindTestDAppNoRevert(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevertTransactor{contract: contract}, nil
}

// NewTestDAppNoRevertFilterer creates a new log filterer instance of TestDAppNoRevert, bound to a specific deployed contract.
func NewTestDAppNoRevertFilterer(address common.Address, filterer bind.ContractFilterer) (*TestDAppNoRevertFilterer, error) {
	contract, err := bindTestDAppNoRevert(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevertFilterer{contract: contract}, nil
}

// bindTestDAppNoRevert binds a generic wrapper to an already deployed contract.
func bindTestDAppNoRevert(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestDAppNoRevertMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDAppNoRevert *TestDAppNoRevertRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDAppNoRevert.Contract.TestDAppNoRevertCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDAppNoRevert *TestDAppNoRevertRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.TestDAppNoRevertTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDAppNoRevert *TestDAppNoRevertRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.TestDAppNoRevertTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDAppNoRevert *TestDAppNoRevertCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDAppNoRevert.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDAppNoRevert *TestDAppNoRevertTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDAppNoRevert *TestDAppNoRevertTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.contract.Transact(opts, method, params...)
}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDAppNoRevert *TestDAppNoRevertCaller) HELLOWORLDMESSAGETYPE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TestDAppNoRevert.contract.Call(opts, &out, "HELLO_WORLD_MESSAGE_TYPE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDAppNoRevert *TestDAppNoRevertSession) HELLOWORLDMESSAGETYPE() ([32]byte, error) {
	return _TestDAppNoRevert.Contract.HELLOWORLDMESSAGETYPE(&_TestDAppNoRevert.CallOpts)
}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDAppNoRevert *TestDAppNoRevertCallerSession) HELLOWORLDMESSAGETYPE() ([32]byte, error) {
	return _TestDAppNoRevert.Contract.HELLOWORLDMESSAGETYPE(&_TestDAppNoRevert.CallOpts)
}

// Connector is a free data retrieval call binding the contract method 0x83f3084f.
//
// Solidity: function connector() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertCaller) Connector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestDAppNoRevert.contract.Call(opts, &out, "connector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Connector is a free data retrieval call binding the contract method 0x83f3084f.
//
// Solidity: function connector() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertSession) Connector() (common.Address, error) {
	return _TestDAppNoRevert.Contract.Connector(&_TestDAppNoRevert.CallOpts)
}

// Connector is a free data retrieval call binding the contract method 0x83f3084f.
//
// Solidity: function connector() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertCallerSession) Connector() (common.Address, error) {
	return _TestDAppNoRevert.Contract.Connector(&_TestDAppNoRevert.CallOpts)
}

// Muse is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function muse() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertCaller) Muse(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestDAppNoRevert.contract.Call(opts, &out, "muse")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Muse is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function muse() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertSession) Muse() (common.Address, error) {
	return _TestDAppNoRevert.Contract.Muse(&_TestDAppNoRevert.CallOpts)
}

// Muse is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function muse() view returns(address)
func (_TestDAppNoRevert *TestDAppNoRevertCallerSession) Muse() (common.Address, error) {
	return _TestDAppNoRevert.Contract.Muse(&_TestDAppNoRevert.CallOpts)
}

// OnMuseMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onMuseMessage((bytes,uint256,address,uint256,bytes) museMessage) returns()
func (_TestDAppNoRevert *TestDAppNoRevertTransactor) OnMuseMessage(opts *bind.TransactOpts, museMessage MuseInterfacesMuseMessage) (*types.Transaction, error) {
	return _TestDAppNoRevert.contract.Transact(opts, "onMuseMessage", museMessage)
}

// OnMuseMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onMuseMessage((bytes,uint256,address,uint256,bytes) museMessage) returns()
func (_TestDAppNoRevert *TestDAppNoRevertSession) OnMuseMessage(museMessage MuseInterfacesMuseMessage) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.OnMuseMessage(&_TestDAppNoRevert.TransactOpts, museMessage)
}

// OnMuseMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onMuseMessage((bytes,uint256,address,uint256,bytes) museMessage) returns()
func (_TestDAppNoRevert *TestDAppNoRevertTransactorSession) OnMuseMessage(museMessage MuseInterfacesMuseMessage) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.OnMuseMessage(&_TestDAppNoRevert.TransactOpts, museMessage)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDAppNoRevert *TestDAppNoRevertTransactor) SendHelloWorld(opts *bind.TransactOpts, destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDAppNoRevert.contract.Transact(opts, "sendHelloWorld", destinationAddress, destinationChainId, value, doRevert)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDAppNoRevert *TestDAppNoRevertSession) SendHelloWorld(destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.SendHelloWorld(&_TestDAppNoRevert.TransactOpts, destinationAddress, destinationChainId, value, doRevert)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDAppNoRevert *TestDAppNoRevertTransactorSession) SendHelloWorld(destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDAppNoRevert.Contract.SendHelloWorld(&_TestDAppNoRevert.TransactOpts, destinationAddress, destinationChainId, value, doRevert)
}

// TestDAppNoRevertHelloWorldEventIterator is returned from FilterHelloWorldEvent and is used to iterate over the raw logs and unpacked data for HelloWorldEvent events raised by the TestDAppNoRevert contract.
type TestDAppNoRevertHelloWorldEventIterator struct {
	Event *TestDAppNoRevertHelloWorldEvent // Event containing the contract specifics and raw log

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
func (it *TestDAppNoRevertHelloWorldEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestDAppNoRevertHelloWorldEvent)
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
		it.Event = new(TestDAppNoRevertHelloWorldEvent)
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
func (it *TestDAppNoRevertHelloWorldEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestDAppNoRevertHelloWorldEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestDAppNoRevertHelloWorldEvent represents a HelloWorldEvent event raised by the TestDAppNoRevert contract.
type TestDAppNoRevertHelloWorldEvent struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterHelloWorldEvent is a free log retrieval operation binding the contract event 0x3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d1.
//
// Solidity: event HelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) FilterHelloWorldEvent(opts *bind.FilterOpts) (*TestDAppNoRevertHelloWorldEventIterator, error) {

	logs, sub, err := _TestDAppNoRevert.contract.FilterLogs(opts, "HelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevertHelloWorldEventIterator{contract: _TestDAppNoRevert.contract, event: "HelloWorldEvent", logs: logs, sub: sub}, nil
}

// WatchHelloWorldEvent is a free log subscription operation binding the contract event 0x3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d1.
//
// Solidity: event HelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) WatchHelloWorldEvent(opts *bind.WatchOpts, sink chan<- *TestDAppNoRevertHelloWorldEvent) (event.Subscription, error) {

	logs, sub, err := _TestDAppNoRevert.contract.WatchLogs(opts, "HelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestDAppNoRevertHelloWorldEvent)
				if err := _TestDAppNoRevert.contract.UnpackLog(event, "HelloWorldEvent", log); err != nil {
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

// ParseHelloWorldEvent is a log parse operation binding the contract event 0x3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d1.
//
// Solidity: event HelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) ParseHelloWorldEvent(log types.Log) (*TestDAppNoRevertHelloWorldEvent, error) {
	event := new(TestDAppNoRevertHelloWorldEvent)
	if err := _TestDAppNoRevert.contract.UnpackLog(event, "HelloWorldEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestDAppNoRevertRevertedHelloWorldEventIterator is returned from FilterRevertedHelloWorldEvent and is used to iterate over the raw logs and unpacked data for RevertedHelloWorldEvent events raised by the TestDAppNoRevert contract.
type TestDAppNoRevertRevertedHelloWorldEventIterator struct {
	Event *TestDAppNoRevertRevertedHelloWorldEvent // Event containing the contract specifics and raw log

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
func (it *TestDAppNoRevertRevertedHelloWorldEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestDAppNoRevertRevertedHelloWorldEvent)
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
		it.Event = new(TestDAppNoRevertRevertedHelloWorldEvent)
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
func (it *TestDAppNoRevertRevertedHelloWorldEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestDAppNoRevertRevertedHelloWorldEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestDAppNoRevertRevertedHelloWorldEvent represents a RevertedHelloWorldEvent event raised by the TestDAppNoRevert contract.
type TestDAppNoRevertRevertedHelloWorldEvent struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRevertedHelloWorldEvent is a free log retrieval operation binding the contract event 0x4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a521.
//
// Solidity: event RevertedHelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) FilterRevertedHelloWorldEvent(opts *bind.FilterOpts) (*TestDAppNoRevertRevertedHelloWorldEventIterator, error) {

	logs, sub, err := _TestDAppNoRevert.contract.FilterLogs(opts, "RevertedHelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return &TestDAppNoRevertRevertedHelloWorldEventIterator{contract: _TestDAppNoRevert.contract, event: "RevertedHelloWorldEvent", logs: logs, sub: sub}, nil
}

// WatchRevertedHelloWorldEvent is a free log subscription operation binding the contract event 0x4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a521.
//
// Solidity: event RevertedHelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) WatchRevertedHelloWorldEvent(opts *bind.WatchOpts, sink chan<- *TestDAppNoRevertRevertedHelloWorldEvent) (event.Subscription, error) {

	logs, sub, err := _TestDAppNoRevert.contract.WatchLogs(opts, "RevertedHelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestDAppNoRevertRevertedHelloWorldEvent)
				if err := _TestDAppNoRevert.contract.UnpackLog(event, "RevertedHelloWorldEvent", log); err != nil {
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

// ParseRevertedHelloWorldEvent is a log parse operation binding the contract event 0x4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a521.
//
// Solidity: event RevertedHelloWorldEvent()
func (_TestDAppNoRevert *TestDAppNoRevertFilterer) ParseRevertedHelloWorldEvent(log types.Log) (*TestDAppNoRevertRevertedHelloWorldEvent, error) {
	event := new(TestDAppNoRevertRevertedHelloWorldEvent)
	if err := _TestDAppNoRevert.contract.UnpackLog(event, "RevertedHelloWorldEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
