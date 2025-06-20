// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testdapp

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

// TestDAppMetaData contains all meta data concerning the TestDApp contract.
var TestDAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_connector\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_museToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ErrorTransferringMuse\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidMessageType\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"HelloWorldEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"RevertedHelloWorldEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"HELLO_WORLD_MESSAGE_TYPE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"connector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"museTxSenderAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"destinationAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"museValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"internalType\":\"structMuseInterfaces.MuseMessage\",\"name\":\"museMessage\",\"type\":\"tuple\"}],\"name\":\"onMuseMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"museTxSenderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sourceChainId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"destinationAddress\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"remainingMuseValue\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"internalType\":\"structMuseInterfaces.MuseRevert\",\"name\":\"museRevert\",\"type\":\"tuple\"}],\"name\":\"onMuseRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"destinationAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"destinationChainId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"doRevert\",\"type\":\"bool\"}],\"name\":\"sendHelloWorld\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"muse\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405162000ebe38038062000ebe8339818101604052810190610034919061011f565b816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505061015f565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006100ec826100c1565b9050919050565b6100fc816100e1565b811461010757600080fd5b50565b600081519050610119816100f3565b92915050565b60008060408385031215610136576101356100bc565b5b60006101448582860161010a565b92505060206101558582860161010a565b9150509250929050565b610d4f806200016f6000396000f3fe6080604052600436106100555760003560e01c80633749c51a1461005a5780633ff0693c146100835780637caca304146100ac57806383f3084f146100c85780638ac44a3f146100f3578063e8f9cb3a1461011e575b600080fd5b34801561006657600080fd5b50610081600480360381019061007c9190610600565b610149565b005b34801561008f57600080fd5b506100aa60048036038101906100a59190610668565b6101e2565b005b6100c660048036038101906100c1919061077d565b61027b565b005b3480156100d457600080fd5b506100dd610564565b6040516100ea91906107f3565b60405180910390f35b3480156100ff57600080fd5b50610108610588565b6040516101159190610827565b60405180910390f35b34801561012a57600080fd5b506101336105ac565b60405161014091906107f3565b60405180910390f35b600081806080019061015b9190610851565b81019061016891906108e0565b91505060001515811515146101b2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101a99061097d565b60405180910390fd5b7f3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d160405160405180910390a15050565b6000818060a001906101f49190610851565b81019061020191906108e0565b915050600115158115151461024b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161024290610a0f565b60405180910390fd5b7f4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a52160405160405180910390a15050565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663095ea7b360008054906101000a900473ffffffffffffffffffffffffffffffffffffffff16856040518363ffffffff1660e01b81526004016102fa929190610a3e565b6020604051808303816000875af1158015610319573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061033d9190610a7c565b90506000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166323b872dd3330876040518463ffffffff1660e01b81526004016103a093929190610aa9565b6020604051808303816000875af11580156103bf573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103e39190610a7c565b90508180156103ef5750805b610425576040517f2bd0ba5000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663ec0269016040518060c00160405280888152602001896040516020016104839190610b28565b60405160208183030381529060405281526020016203d09081526020017f6e0182194bb1deba01849afd3e035a0b70ce7cb069e482ee663519c76cf569b4876040516020016104d3929190610b52565b60405160208183030381529060405281526020018781526020016040516020016104fc90610ba1565b6040516020818303038152906040528152506040518263ffffffff1660e01b815260040161052a9190610cf7565b600060405180830381600087803b15801561054457600080fd5b505af1158015610558573d6000803e3d6000fd5b50505050505050505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b7f6e0182194bb1deba01849afd3e035a0b70ce7cb069e482ee663519c76cf569b481565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600080fd5b600080fd5b600080fd5b600060a082840312156105f7576105f66105dc565b5b81905092915050565b600060208284031215610616576106156105d2565b5b600082013567ffffffffffffffff811115610634576106336105d7565b5b610640848285016105e1565b91505092915050565b600060c0828403121561065f5761065e6105dc565b5b81905092915050565b60006020828403121561067e5761067d6105d2565b5b600082013567ffffffffffffffff81111561069c5761069b6105d7565b5b6106a884828501610649565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006106dc826106b1565b9050919050565b6106ec816106d1565b81146106f757600080fd5b50565b600081359050610709816106e3565b92915050565b6000819050919050565b6107228161070f565b811461072d57600080fd5b50565b60008135905061073f81610719565b92915050565b60008115159050919050565b61075a81610745565b811461076557600080fd5b50565b60008135905061077781610751565b92915050565b60008060008060808587031215610797576107966105d2565b5b60006107a5878288016106fa565b94505060206107b687828801610730565b93505060406107c787828801610730565b92505060606107d887828801610768565b91505092959194509250565b6107ed816106d1565b82525050565b600060208201905061080860008301846107e4565b92915050565b6000819050919050565b6108218161080e565b82525050565b600060208201905061083c6000830184610818565b92915050565b600080fd5b600080fd5b600080fd5b6000808335600160200384360303811261086e5761086d610842565b5b80840192508235915067ffffffffffffffff8211156108905761088f610847565b5b6020830192506001820236038313156108ac576108ab61084c565b5b509250929050565b6108bd8161080e565b81146108c857600080fd5b50565b6000813590506108da816108b4565b92915050565b600080604083850312156108f7576108f66105d2565b5b6000610905858286016108cb565b925050602061091685828601610768565b9150509250929050565b600082825260208201905092915050565b7f6d65737361676520736179732072657665727400000000000000000000000000600082015250565b6000610967601383610920565b915061097282610931565b602082019050919050565b600060208201905081810360008301526109968161095a565b9050919050565b7f74686520317374206f7574626f756e6420776173206e6f74206361757365642060008201527f62792072657665727420666c616720696e206d65737361676500000000000000602082015250565b60006109f9603983610920565b9150610a048261099d565b604082019050919050565b60006020820190508181036000830152610a28816109ec565b9050919050565b610a388161070f565b82525050565b6000604082019050610a5360008301856107e4565b610a606020830184610a2f565b9392505050565b600081519050610a7681610751565b92915050565b600060208284031215610a9257610a916105d2565b5b6000610aa084828501610a67565b91505092915050565b6000606082019050610abe60008301866107e4565b610acb60208301856107e4565b610ad86040830184610a2f565b949350505050565b60008160601b9050919050565b6000610af882610ae0565b9050919050565b6000610b0a82610aed565b9050919050565b610b22610b1d826106d1565b610aff565b82525050565b6000610b348284610b11565b60148201915081905092915050565b610b4c81610745565b82525050565b6000604082019050610b676000830185610818565b610b746020830184610b43565b9392505050565b50565b6000610b8b600083610920565b9150610b9682610b7b565b600082019050919050565b60006020820190508181036000830152610bba81610b7e565b9050919050565b610bca8161070f565b82525050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610c0a578082015181840152602081019050610bef565b60008484015250505050565b6000601f19601f8301169050919050565b6000610c3282610bd0565b610c3c8185610bdb565b9350610c4c818560208601610bec565b610c5581610c16565b840191505092915050565b600060c083016000830151610c786000860182610bc1565b5060208301518482036020860152610c908282610c27565b9150506040830151610ca56040860182610bc1565b5060608301518482036060860152610cbd8282610c27565b9150506080830151610cd26080860182610bc1565b5060a083015184820360a0860152610cea8282610c27565b9150508091505092915050565b60006020820190508181036000830152610d118184610c60565b90509291505056fea2646970667358221220c212f42910156db82872f9a01c6651ca87289ee6dade044f8864279558273d1164736f6c63430008120033",
}

// TestDAppABI is the input ABI used to generate the binding from.
// Deprecated: Use TestDAppMetaData.ABI instead.
var TestDAppABI = TestDAppMetaData.ABI

// TestDAppBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestDAppMetaData.Bin instead.
var TestDAppBin = TestDAppMetaData.Bin

// DeployTestDApp deploys a new Ethereum contract, binding an instance of TestDApp to it.
func DeployTestDApp(auth *bind.TransactOpts, backend bind.ContractBackend, _connector common.Address, _museToken common.Address) (common.Address, *types.Transaction, *TestDApp, error) {
	parsed, err := TestDAppMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestDAppBin), backend, _connector, _museToken)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestDApp{TestDAppCaller: TestDAppCaller{contract: contract}, TestDAppTransactor: TestDAppTransactor{contract: contract}, TestDAppFilterer: TestDAppFilterer{contract: contract}}, nil
}

// TestDApp is an auto generated Go binding around an Ethereum contract.
type TestDApp struct {
	TestDAppCaller     // Read-only binding to the contract
	TestDAppTransactor // Write-only binding to the contract
	TestDAppFilterer   // Log filterer for contract events
}

// TestDAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestDAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestDAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestDAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestDAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestDAppSession struct {
	Contract     *TestDApp         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestDAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestDAppCallerSession struct {
	Contract *TestDAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TestDAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestDAppTransactorSession struct {
	Contract     *TestDAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TestDAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestDAppRaw struct {
	Contract *TestDApp // Generic contract binding to access the raw methods on
}

// TestDAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestDAppCallerRaw struct {
	Contract *TestDAppCaller // Generic read-only contract binding to access the raw methods on
}

// TestDAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestDAppTransactorRaw struct {
	Contract *TestDAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestDApp creates a new instance of TestDApp, bound to a specific deployed contract.
func NewTestDApp(address common.Address, backend bind.ContractBackend) (*TestDApp, error) {
	contract, err := bindTestDApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestDApp{TestDAppCaller: TestDAppCaller{contract: contract}, TestDAppTransactor: TestDAppTransactor{contract: contract}, TestDAppFilterer: TestDAppFilterer{contract: contract}}, nil
}

// NewTestDAppCaller creates a new read-only instance of TestDApp, bound to a specific deployed contract.
func NewTestDAppCaller(address common.Address, caller bind.ContractCaller) (*TestDAppCaller, error) {
	contract, err := bindTestDApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppCaller{contract: contract}, nil
}

// NewTestDAppTransactor creates a new write-only instance of TestDApp, bound to a specific deployed contract.
func NewTestDAppTransactor(address common.Address, transactor bind.ContractTransactor) (*TestDAppTransactor, error) {
	contract, err := bindTestDApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestDAppTransactor{contract: contract}, nil
}

// NewTestDAppFilterer creates a new log filterer instance of TestDApp, bound to a specific deployed contract.
func NewTestDAppFilterer(address common.Address, filterer bind.ContractFilterer) (*TestDAppFilterer, error) {
	contract, err := bindTestDApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestDAppFilterer{contract: contract}, nil
}

// bindTestDApp binds a generic wrapper to an already deployed contract.
func bindTestDApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestDAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDApp *TestDAppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDApp.Contract.TestDAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDApp *TestDAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDApp.Contract.TestDAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDApp *TestDAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDApp.Contract.TestDAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestDApp *TestDAppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestDApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestDApp *TestDAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestDApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestDApp *TestDAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestDApp.Contract.contract.Transact(opts, method, params...)
}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDApp *TestDAppCaller) HELLOWORLDMESSAGETYPE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TestDApp.contract.Call(opts, &out, "HELLO_WORLD_MESSAGE_TYPE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDApp *TestDAppSession) HELLOWORLDMESSAGETYPE() ([32]byte, error) {
	return _TestDApp.Contract.HELLOWORLDMESSAGETYPE(&_TestDApp.CallOpts)
}

// HELLOWORLDMESSAGETYPE is a free data retrieval call binding the contract method 0x8ac44a3f.
//
// Solidity: function HELLO_WORLD_MESSAGE_TYPE() view returns(bytes32)
func (_TestDApp *TestDAppCallerSession) HELLOWORLDMESSAGETYPE() ([32]byte, error) {
	return _TestDApp.Contract.HELLOWORLDMESSAGETYPE(&_TestDApp.CallOpts)
}

// Connector is a free data retrieval call binding the contract method 0x83f3084f.
//
// Solidity: function connector() view returns(address)
func (_TestDApp *TestDAppCaller) Connector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestDApp.contract.Call(opts, &out, "connector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Connector is a free data retrieval call binding the contract method 0x83f3084f.
//
// Solidity: function connector() view returns(address)
func (_TestDApp *TestDAppSession) Connector() (common.Address, error) {
	return _TestDApp.Contract.Connector(&_TestDApp.CallOpts)
}

// Connector is a free data retrieval call binding the contract method 0x83f3084f.
//
// Solidity: function connector() view returns(address)
func (_TestDApp *TestDAppCallerSession) Connector() (common.Address, error) {
	return _TestDApp.Contract.Connector(&_TestDApp.CallOpts)
}

// Muse is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function muse() view returns(address)
func (_TestDApp *TestDAppCaller) Muse(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestDApp.contract.Call(opts, &out, "muse")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Muse is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function muse() view returns(address)
func (_TestDApp *TestDAppSession) Muse() (common.Address, error) {
	return _TestDApp.Contract.Muse(&_TestDApp.CallOpts)
}

// Muse is a free data retrieval call binding the contract method 0xe8f9cb3a.
//
// Solidity: function muse() view returns(address)
func (_TestDApp *TestDAppCallerSession) Muse() (common.Address, error) {
	return _TestDApp.Contract.Muse(&_TestDApp.CallOpts)
}

// OnMuseMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onMuseMessage((bytes,uint256,address,uint256,bytes) museMessage) returns()
func (_TestDApp *TestDAppTransactor) OnMuseMessage(opts *bind.TransactOpts, museMessage MuseInterfacesMuseMessage) (*types.Transaction, error) {
	return _TestDApp.contract.Transact(opts, "onMuseMessage", museMessage)
}

// OnMuseMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onMuseMessage((bytes,uint256,address,uint256,bytes) museMessage) returns()
func (_TestDApp *TestDAppSession) OnMuseMessage(museMessage MuseInterfacesMuseMessage) (*types.Transaction, error) {
	return _TestDApp.Contract.OnMuseMessage(&_TestDApp.TransactOpts, museMessage)
}

// OnMuseMessage is a paid mutator transaction binding the contract method 0x3749c51a.
//
// Solidity: function onMuseMessage((bytes,uint256,address,uint256,bytes) museMessage) returns()
func (_TestDApp *TestDAppTransactorSession) OnMuseMessage(museMessage MuseInterfacesMuseMessage) (*types.Transaction, error) {
	return _TestDApp.Contract.OnMuseMessage(&_TestDApp.TransactOpts, museMessage)
}

// OnMuseRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onMuseRevert((address,uint256,bytes,uint256,uint256,bytes) museRevert) returns()
func (_TestDApp *TestDAppTransactor) OnMuseRevert(opts *bind.TransactOpts, museRevert MuseInterfacesMuseRevert) (*types.Transaction, error) {
	return _TestDApp.contract.Transact(opts, "onMuseRevert", museRevert)
}

// OnMuseRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onMuseRevert((address,uint256,bytes,uint256,uint256,bytes) museRevert) returns()
func (_TestDApp *TestDAppSession) OnMuseRevert(museRevert MuseInterfacesMuseRevert) (*types.Transaction, error) {
	return _TestDApp.Contract.OnMuseRevert(&_TestDApp.TransactOpts, museRevert)
}

// OnMuseRevert is a paid mutator transaction binding the contract method 0x3ff0693c.
//
// Solidity: function onMuseRevert((address,uint256,bytes,uint256,uint256,bytes) museRevert) returns()
func (_TestDApp *TestDAppTransactorSession) OnMuseRevert(museRevert MuseInterfacesMuseRevert) (*types.Transaction, error) {
	return _TestDApp.Contract.OnMuseRevert(&_TestDApp.TransactOpts, museRevert)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDApp *TestDAppTransactor) SendHelloWorld(opts *bind.TransactOpts, destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDApp.contract.Transact(opts, "sendHelloWorld", destinationAddress, destinationChainId, value, doRevert)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDApp *TestDAppSession) SendHelloWorld(destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDApp.Contract.SendHelloWorld(&_TestDApp.TransactOpts, destinationAddress, destinationChainId, value, doRevert)
}

// SendHelloWorld is a paid mutator transaction binding the contract method 0x7caca304.
//
// Solidity: function sendHelloWorld(address destinationAddress, uint256 destinationChainId, uint256 value, bool doRevert) payable returns()
func (_TestDApp *TestDAppTransactorSession) SendHelloWorld(destinationAddress common.Address, destinationChainId *big.Int, value *big.Int, doRevert bool) (*types.Transaction, error) {
	return _TestDApp.Contract.SendHelloWorld(&_TestDApp.TransactOpts, destinationAddress, destinationChainId, value, doRevert)
}

// TestDAppHelloWorldEventIterator is returned from FilterHelloWorldEvent and is used to iterate over the raw logs and unpacked data for HelloWorldEvent events raised by the TestDApp contract.
type TestDAppHelloWorldEventIterator struct {
	Event *TestDAppHelloWorldEvent // Event containing the contract specifics and raw log

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
func (it *TestDAppHelloWorldEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestDAppHelloWorldEvent)
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
		it.Event = new(TestDAppHelloWorldEvent)
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
func (it *TestDAppHelloWorldEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestDAppHelloWorldEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestDAppHelloWorldEvent represents a HelloWorldEvent event raised by the TestDApp contract.
type TestDAppHelloWorldEvent struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterHelloWorldEvent is a free log retrieval operation binding the contract event 0x3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d1.
//
// Solidity: event HelloWorldEvent()
func (_TestDApp *TestDAppFilterer) FilterHelloWorldEvent(opts *bind.FilterOpts) (*TestDAppHelloWorldEventIterator, error) {

	logs, sub, err := _TestDApp.contract.FilterLogs(opts, "HelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return &TestDAppHelloWorldEventIterator{contract: _TestDApp.contract, event: "HelloWorldEvent", logs: logs, sub: sub}, nil
}

// WatchHelloWorldEvent is a free log subscription operation binding the contract event 0x3399097dded3a4667baa7375fe02dfaec8fb76c75ba8da569c40bd175686b0d1.
//
// Solidity: event HelloWorldEvent()
func (_TestDApp *TestDAppFilterer) WatchHelloWorldEvent(opts *bind.WatchOpts, sink chan<- *TestDAppHelloWorldEvent) (event.Subscription, error) {

	logs, sub, err := _TestDApp.contract.WatchLogs(opts, "HelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestDAppHelloWorldEvent)
				if err := _TestDApp.contract.UnpackLog(event, "HelloWorldEvent", log); err != nil {
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
func (_TestDApp *TestDAppFilterer) ParseHelloWorldEvent(log types.Log) (*TestDAppHelloWorldEvent, error) {
	event := new(TestDAppHelloWorldEvent)
	if err := _TestDApp.contract.UnpackLog(event, "HelloWorldEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestDAppRevertedHelloWorldEventIterator is returned from FilterRevertedHelloWorldEvent and is used to iterate over the raw logs and unpacked data for RevertedHelloWorldEvent events raised by the TestDApp contract.
type TestDAppRevertedHelloWorldEventIterator struct {
	Event *TestDAppRevertedHelloWorldEvent // Event containing the contract specifics and raw log

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
func (it *TestDAppRevertedHelloWorldEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestDAppRevertedHelloWorldEvent)
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
		it.Event = new(TestDAppRevertedHelloWorldEvent)
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
func (it *TestDAppRevertedHelloWorldEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestDAppRevertedHelloWorldEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestDAppRevertedHelloWorldEvent represents a RevertedHelloWorldEvent event raised by the TestDApp contract.
type TestDAppRevertedHelloWorldEvent struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRevertedHelloWorldEvent is a free log retrieval operation binding the contract event 0x4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a521.
//
// Solidity: event RevertedHelloWorldEvent()
func (_TestDApp *TestDAppFilterer) FilterRevertedHelloWorldEvent(opts *bind.FilterOpts) (*TestDAppRevertedHelloWorldEventIterator, error) {

	logs, sub, err := _TestDApp.contract.FilterLogs(opts, "RevertedHelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return &TestDAppRevertedHelloWorldEventIterator{contract: _TestDApp.contract, event: "RevertedHelloWorldEvent", logs: logs, sub: sub}, nil
}

// WatchRevertedHelloWorldEvent is a free log subscription operation binding the contract event 0x4f30bf4846ce4cde02361b3232cd2287313384a7b8e60161a1b2818b6905a521.
//
// Solidity: event RevertedHelloWorldEvent()
func (_TestDApp *TestDAppFilterer) WatchRevertedHelloWorldEvent(opts *bind.WatchOpts, sink chan<- *TestDAppRevertedHelloWorldEvent) (event.Subscription, error) {

	logs, sub, err := _TestDApp.contract.WatchLogs(opts, "RevertedHelloWorldEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestDAppRevertedHelloWorldEvent)
				if err := _TestDApp.contract.UnpackLog(event, "RevertedHelloWorldEvent", log); err != nil {
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
func (_TestDApp *TestDAppFilterer) ParseRevertedHelloWorldEvent(log types.Log) (*TestDAppRevertedHelloWorldEvent, error) {
	event := new(TestDAppRevertedHelloWorldEvent)
	if err := _TestDApp.contract.UnpackLog(event, "RevertedHelloWorldEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
