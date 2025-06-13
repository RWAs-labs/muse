// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package museconnector

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

// MuseInterfacesSendInput is an auto generated low-level Go binding around an user-defined struct.
type MuseInterfacesSendInput struct {
	DestinationChainId  *big.Int
	DestinationAddress  []byte
	DestinationGasLimit *big.Int
	Message             []byte
	MuseValueAndGas     *big.Int
	MuseParams          []byte
}

// MuseConnectorBaseMetaData contains all meta data concerning the MuseConnectorBase contract.
var MuseConnectorBaseMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"museToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tssAddress_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tssAddressUpdater_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pauserAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onReceive\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"museValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onRevert\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingMuseValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceTssAddressUpdater\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"send\",\"inputs\":[{\"name\":\"input\",\"type\":\"tuple\",\"internalType\":\"structMuseInterfaces.SendInput\",\"components\":[{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destinationGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"museValueAndGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"museParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tssAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tssAddressUpdater\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauserAddress\",\"inputs\":[{\"name\":\"pauserAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTssAddress\",\"inputs\":[{\"name\":\"tssAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"museToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserAddressUpdated\",\"inputs\":[{\"name\":\"callerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newTssAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TSSAddressUpdated\",\"inputs\":[{\"name\":\"callerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newTssAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TSSAddressUpdaterUpdated\",\"inputs\":[{\"name\":\"callerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newTssUpdaterAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MuseReceived\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"museValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MuseReverted\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"remainingMuseValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MuseSent\",\"inputs\":[{\"name\":\"sourceTxOriginAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"museTxSenderAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"museValueAndGas\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationGasLimit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"museParams\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CallerIsNotPauser\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CallerIsNotTss\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CallerIsNotTssOrUpdater\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CallerIsNotTssUpdater\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExceedsMaxSupply\",\"inputs\":[{\"name\":\"maxSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MuseTransferError\",\"inputs\":[]}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610b9f380380610b9f83398101604081905261002f9161010b565b6000805460ff191690556001600160a01b038416158061005657506001600160a01b038316155b8061006857506001600160a01b038216155b8061007a57506001600160a01b038116155b156100985760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b03938416608052600180549385166001600160a01b0319948516179055600280549285169290931691909117909155600080549190921661010002610100600160a81b031990911617905561015f565b80516001600160a01b038116811461010657600080fd5b919050565b6000806000806080858703121561012157600080fd5b61012a856100ef565b9350610138602086016100ef565b9250610146604086016100ef565b9150610154606086016100ef565b905092959194509250565b608051610a26610179600039600060e90152610a266000f3fe608060405234801561001057600080fd5b50600436106100df5760003560e01c80636128480f1161008c5780639122c344116100665780639122c344146101d0578063942a5e16146101e3578063ec026901146101fc578063f7fb869b1461020d57600080fd5b80636128480f146101ad578063779e3b63146101c05780638456cb59146101c857600080fd5b80633f4ba83a116100bd5780633f4ba83a1461016f5780635b112591146101775780635c975abb1461019757600080fd5b806321e093b1146100e457806329dd214d14610135578063328a01d01461014f575b600080fd5b61010b7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b61014d610143366004610849565b5050505050505050565b005b60025461010b9073ffffffffffffffffffffffffffffffffffffffff1681565b61014d610232565b60015461010b9073ffffffffffffffffffffffffffffffffffffffff1681565b60005460ff16604051901515815260200161012c565b61014d6101bb3660046108e9565b610299565b61014d6103c6565b61014d6104eb565b61014d6101de3660046108e9565b61054b565b61014d6101f136600461090b565b505050505050505050565b61014d61020a3660046109b5565b50565b60005461010b90610100900473ffffffffffffffffffffffffffffffffffffffff1681565b600054610100900473ffffffffffffffffffffffffffffffffffffffff16331461028f576040517f4677a0d30000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b61029761068b565b565b600054610100900473ffffffffffffffffffffffffffffffffffffffff1633146102f1576040517f4677a0d3000000000000000000000000000000000000000000000000000000008152336004820152602401610286565b73ffffffffffffffffffffffffffffffffffffffff811661033e576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffff0000000000000000000000000000000000000000ff1661010073ffffffffffffffffffffffffffffffffffffffff8416908102919091179091556040805133815260208101929092527fd41d83655d484bdf299598751c371b2d92088667266fe3774b25a97bdd5d039791015b60405180910390a150565b60025473ffffffffffffffffffffffffffffffffffffffff163314610419576040517fe700765e000000000000000000000000000000000000000000000000000000008152336004820152602401610286565b60015473ffffffffffffffffffffffffffffffffffffffff16610468576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600154600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691821790556040805133815260208101929092527f5104c9abdc7d111c2aeb4ce890ac70274a4be2ee83f46a62551be5e6ebc82dd091015b60405180910390a1565b600054610100900473ffffffffffffffffffffffffffffffffffffffff163314610543576040517f4677a0d3000000000000000000000000000000000000000000000000000000008152336004820152602401610286565b610297610703565b60015473ffffffffffffffffffffffffffffffffffffffff16331480159061058b575060025473ffffffffffffffffffffffffffffffffffffffff163314155b156105c4576040517fcdfcef97000000000000000000000000000000000000000000000000000000008152336004820152602401610286565b73ffffffffffffffffffffffffffffffffffffffff8116610611576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040805133815260208101929092527fe79965b5c67dcfb2cf5fe152715e4a7256cee62a3d5dd8484fd8a8539eb8beff91016103bb565b61069361075e565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016104e1565b61070b61079a565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586106de3390565b60005460ff16610297576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005460ff1615610297576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008083601f8401126107e957600080fd5b50813567ffffffffffffffff81111561080157600080fd5b60208301915083602082850101111561081957600080fd5b9250929050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461084457600080fd5b919050565b60008060008060008060008060c0898b03121561086557600080fd5b883567ffffffffffffffff81111561087c57600080fd5b6108888b828c016107d7565b909950975050602089013595506108a160408a01610820565b945060608901359350608089013567ffffffffffffffff8111156108c457600080fd5b6108d08b828c016107d7565b999c989b50969995989497949560a00135949350505050565b6000602082840312156108fb57600080fd5b61090482610820565b9392505050565b600080600080600080600080600060e08a8c03121561092957600080fd5b6109328a610820565b985060208a0135975060408a013567ffffffffffffffff81111561095557600080fd5b6109618c828d016107d7565b90985096505060608a0135945060808a0135935060a08a013567ffffffffffffffff81111561098f57600080fd5b61099b8c828d016107d7565b9a9d999c50979a9699959894979660c00135949350505050565b6000602082840312156109c757600080fd5b813567ffffffffffffffff8111156109de57600080fd5b820160c0818503121561090457600080fdfea2646970667358221220f1af43c04120e57dee0ccdf659d588e4828c4abe2c22be29c11469808eac0a9d64736f6c634300081a0033",
}

// MuseConnectorBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use MuseConnectorBaseMetaData.ABI instead.
var MuseConnectorBaseABI = MuseConnectorBaseMetaData.ABI

// MuseConnectorBaseBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MuseConnectorBaseMetaData.Bin instead.
var MuseConnectorBaseBin = MuseConnectorBaseMetaData.Bin

// DeployMuseConnectorBase deploys a new Ethereum contract, binding an instance of MuseConnectorBase to it.
func DeployMuseConnectorBase(auth *bind.TransactOpts, backend bind.ContractBackend, museToken_ common.Address, tssAddress_ common.Address, tssAddressUpdater_ common.Address, pauserAddress_ common.Address) (common.Address, *types.Transaction, *MuseConnectorBase, error) {
	parsed, err := MuseConnectorBaseMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MuseConnectorBaseBin), backend, museToken_, tssAddress_, tssAddressUpdater_, pauserAddress_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MuseConnectorBase{MuseConnectorBaseCaller: MuseConnectorBaseCaller{contract: contract}, MuseConnectorBaseTransactor: MuseConnectorBaseTransactor{contract: contract}, MuseConnectorBaseFilterer: MuseConnectorBaseFilterer{contract: contract}}, nil
}

// MuseConnectorBase is an auto generated Go binding around an Ethereum contract.
type MuseConnectorBase struct {
	MuseConnectorBaseCaller     // Read-only binding to the contract
	MuseConnectorBaseTransactor // Write-only binding to the contract
	MuseConnectorBaseFilterer   // Log filterer for contract events
}

// MuseConnectorBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type MuseConnectorBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseConnectorBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MuseConnectorBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseConnectorBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MuseConnectorBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseConnectorBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MuseConnectorBaseSession struct {
	Contract     *MuseConnectorBase // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MuseConnectorBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MuseConnectorBaseCallerSession struct {
	Contract *MuseConnectorBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// MuseConnectorBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MuseConnectorBaseTransactorSession struct {
	Contract     *MuseConnectorBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// MuseConnectorBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type MuseConnectorBaseRaw struct {
	Contract *MuseConnectorBase // Generic contract binding to access the raw methods on
}

// MuseConnectorBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MuseConnectorBaseCallerRaw struct {
	Contract *MuseConnectorBaseCaller // Generic read-only contract binding to access the raw methods on
}

// MuseConnectorBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MuseConnectorBaseTransactorRaw struct {
	Contract *MuseConnectorBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMuseConnectorBase creates a new instance of MuseConnectorBase, bound to a specific deployed contract.
func NewMuseConnectorBase(address common.Address, backend bind.ContractBackend) (*MuseConnectorBase, error) {
	contract, err := bindMuseConnectorBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBase{MuseConnectorBaseCaller: MuseConnectorBaseCaller{contract: contract}, MuseConnectorBaseTransactor: MuseConnectorBaseTransactor{contract: contract}, MuseConnectorBaseFilterer: MuseConnectorBaseFilterer{contract: contract}}, nil
}

// NewMuseConnectorBaseCaller creates a new read-only instance of MuseConnectorBase, bound to a specific deployed contract.
func NewMuseConnectorBaseCaller(address common.Address, caller bind.ContractCaller) (*MuseConnectorBaseCaller, error) {
	contract, err := bindMuseConnectorBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBaseCaller{contract: contract}, nil
}

// NewMuseConnectorBaseTransactor creates a new write-only instance of MuseConnectorBase, bound to a specific deployed contract.
func NewMuseConnectorBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*MuseConnectorBaseTransactor, error) {
	contract, err := bindMuseConnectorBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBaseTransactor{contract: contract}, nil
}

// NewMuseConnectorBaseFilterer creates a new log filterer instance of MuseConnectorBase, bound to a specific deployed contract.
func NewMuseConnectorBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*MuseConnectorBaseFilterer, error) {
	contract, err := bindMuseConnectorBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBaseFilterer{contract: contract}, nil
}

// bindMuseConnectorBase binds a generic wrapper to an already deployed contract.
func bindMuseConnectorBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MuseConnectorBaseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MuseConnectorBase *MuseConnectorBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MuseConnectorBase.Contract.MuseConnectorBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MuseConnectorBase *MuseConnectorBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.MuseConnectorBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MuseConnectorBase *MuseConnectorBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.MuseConnectorBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MuseConnectorBase *MuseConnectorBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MuseConnectorBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MuseConnectorBase *MuseConnectorBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MuseConnectorBase *MuseConnectorBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.contract.Transact(opts, method, params...)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MuseConnectorBase *MuseConnectorBaseCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MuseConnectorBase.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MuseConnectorBase *MuseConnectorBaseSession) Paused() (bool, error) {
	return _MuseConnectorBase.Contract.Paused(&_MuseConnectorBase.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MuseConnectorBase *MuseConnectorBaseCallerSession) Paused() (bool, error) {
	return _MuseConnectorBase.Contract.Paused(&_MuseConnectorBase.CallOpts)
}

// PauserAddress is a free data retrieval call binding the contract method 0xf7fb869b.
//
// Solidity: function pauserAddress() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseCaller) PauserAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorBase.contract.Call(opts, &out, "pauserAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PauserAddress is a free data retrieval call binding the contract method 0xf7fb869b.
//
// Solidity: function pauserAddress() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseSession) PauserAddress() (common.Address, error) {
	return _MuseConnectorBase.Contract.PauserAddress(&_MuseConnectorBase.CallOpts)
}

// PauserAddress is a free data retrieval call binding the contract method 0xf7fb869b.
//
// Solidity: function pauserAddress() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseCallerSession) PauserAddress() (common.Address, error) {
	return _MuseConnectorBase.Contract.PauserAddress(&_MuseConnectorBase.CallOpts)
}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseCaller) TssAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorBase.contract.Call(opts, &out, "tssAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseSession) TssAddress() (common.Address, error) {
	return _MuseConnectorBase.Contract.TssAddress(&_MuseConnectorBase.CallOpts)
}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseCallerSession) TssAddress() (common.Address, error) {
	return _MuseConnectorBase.Contract.TssAddress(&_MuseConnectorBase.CallOpts)
}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseCaller) TssAddressUpdater(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorBase.contract.Call(opts, &out, "tssAddressUpdater")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseSession) TssAddressUpdater() (common.Address, error) {
	return _MuseConnectorBase.Contract.TssAddressUpdater(&_MuseConnectorBase.CallOpts)
}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseCallerSession) TssAddressUpdater() (common.Address, error) {
	return _MuseConnectorBase.Contract.TssAddressUpdater(&_MuseConnectorBase.CallOpts)
}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseCaller) MuseToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorBase.contract.Call(opts, &out, "museToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseSession) MuseToken() (common.Address, error) {
	return _MuseConnectorBase.Contract.MuseToken(&_MuseConnectorBase.CallOpts)
}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_MuseConnectorBase *MuseConnectorBaseCallerSession) MuseToken() (common.Address, error) {
	return _MuseConnectorBase.Contract.MuseToken(&_MuseConnectorBase.CallOpts)
}

// OnReceive is a paid mutator transaction binding the contract method 0x29dd214d.
//
// Solidity: function onReceive(bytes museTxSenderAddress, uint256 sourceChainId, address destinationAddress, uint256 museValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactor) OnReceive(opts *bind.TransactOpts, museTxSenderAddress []byte, sourceChainId *big.Int, destinationAddress common.Address, museValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorBase.contract.Transact(opts, "onReceive", museTxSenderAddress, sourceChainId, destinationAddress, museValue, message, internalSendHash)
}

// OnReceive is a paid mutator transaction binding the contract method 0x29dd214d.
//
// Solidity: function onReceive(bytes museTxSenderAddress, uint256 sourceChainId, address destinationAddress, uint256 museValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorBase *MuseConnectorBaseSession) OnReceive(museTxSenderAddress []byte, sourceChainId *big.Int, destinationAddress common.Address, museValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.OnReceive(&_MuseConnectorBase.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, museValue, message, internalSendHash)
}

// OnReceive is a paid mutator transaction binding the contract method 0x29dd214d.
//
// Solidity: function onReceive(bytes museTxSenderAddress, uint256 sourceChainId, address destinationAddress, uint256 museValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactorSession) OnReceive(museTxSenderAddress []byte, sourceChainId *big.Int, destinationAddress common.Address, museValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.OnReceive(&_MuseConnectorBase.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, museValue, message, internalSendHash)
}

// OnRevert is a paid mutator transaction binding the contract method 0x942a5e16.
//
// Solidity: function onRevert(address museTxSenderAddress, uint256 sourceChainId, bytes destinationAddress, uint256 destinationChainId, uint256 remainingMuseValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactor) OnRevert(opts *bind.TransactOpts, museTxSenderAddress common.Address, sourceChainId *big.Int, destinationAddress []byte, destinationChainId *big.Int, remainingMuseValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorBase.contract.Transact(opts, "onRevert", museTxSenderAddress, sourceChainId, destinationAddress, destinationChainId, remainingMuseValue, message, internalSendHash)
}

// OnRevert is a paid mutator transaction binding the contract method 0x942a5e16.
//
// Solidity: function onRevert(address museTxSenderAddress, uint256 sourceChainId, bytes destinationAddress, uint256 destinationChainId, uint256 remainingMuseValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorBase *MuseConnectorBaseSession) OnRevert(museTxSenderAddress common.Address, sourceChainId *big.Int, destinationAddress []byte, destinationChainId *big.Int, remainingMuseValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.OnRevert(&_MuseConnectorBase.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, destinationChainId, remainingMuseValue, message, internalSendHash)
}

// OnRevert is a paid mutator transaction binding the contract method 0x942a5e16.
//
// Solidity: function onRevert(address museTxSenderAddress, uint256 sourceChainId, bytes destinationAddress, uint256 destinationChainId, uint256 remainingMuseValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactorSession) OnRevert(museTxSenderAddress common.Address, sourceChainId *big.Int, destinationAddress []byte, destinationChainId *big.Int, remainingMuseValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.OnRevert(&_MuseConnectorBase.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, destinationChainId, remainingMuseValue, message, internalSendHash)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorBase.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MuseConnectorBase *MuseConnectorBaseSession) Pause() (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.Pause(&_MuseConnectorBase.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactorSession) Pause() (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.Pause(&_MuseConnectorBase.TransactOpts)
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactor) RenounceTssAddressUpdater(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorBase.contract.Transact(opts, "renounceTssAddressUpdater")
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_MuseConnectorBase *MuseConnectorBaseSession) RenounceTssAddressUpdater() (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.RenounceTssAddressUpdater(&_MuseConnectorBase.TransactOpts)
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactorSession) RenounceTssAddressUpdater() (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.RenounceTssAddressUpdater(&_MuseConnectorBase.TransactOpts)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactor) Send(opts *bind.TransactOpts, input MuseInterfacesSendInput) (*types.Transaction, error) {
	return _MuseConnectorBase.contract.Transact(opts, "send", input)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_MuseConnectorBase *MuseConnectorBaseSession) Send(input MuseInterfacesSendInput) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.Send(&_MuseConnectorBase.TransactOpts, input)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactorSession) Send(input MuseInterfacesSendInput) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.Send(&_MuseConnectorBase.TransactOpts, input)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorBase.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MuseConnectorBase *MuseConnectorBaseSession) Unpause() (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.Unpause(&_MuseConnectorBase.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactorSession) Unpause() (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.Unpause(&_MuseConnectorBase.TransactOpts)
}

// UpdatePauserAddress is a paid mutator transaction binding the contract method 0x6128480f.
//
// Solidity: function updatePauserAddress(address pauserAddress_) returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactor) UpdatePauserAddress(opts *bind.TransactOpts, pauserAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorBase.contract.Transact(opts, "updatePauserAddress", pauserAddress_)
}

// UpdatePauserAddress is a paid mutator transaction binding the contract method 0x6128480f.
//
// Solidity: function updatePauserAddress(address pauserAddress_) returns()
func (_MuseConnectorBase *MuseConnectorBaseSession) UpdatePauserAddress(pauserAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.UpdatePauserAddress(&_MuseConnectorBase.TransactOpts, pauserAddress_)
}

// UpdatePauserAddress is a paid mutator transaction binding the contract method 0x6128480f.
//
// Solidity: function updatePauserAddress(address pauserAddress_) returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactorSession) UpdatePauserAddress(pauserAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.UpdatePauserAddress(&_MuseConnectorBase.TransactOpts, pauserAddress_)
}

// UpdateTssAddress is a paid mutator transaction binding the contract method 0x9122c344.
//
// Solidity: function updateTssAddress(address tssAddress_) returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactor) UpdateTssAddress(opts *bind.TransactOpts, tssAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorBase.contract.Transact(opts, "updateTssAddress", tssAddress_)
}

// UpdateTssAddress is a paid mutator transaction binding the contract method 0x9122c344.
//
// Solidity: function updateTssAddress(address tssAddress_) returns()
func (_MuseConnectorBase *MuseConnectorBaseSession) UpdateTssAddress(tssAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.UpdateTssAddress(&_MuseConnectorBase.TransactOpts, tssAddress_)
}

// UpdateTssAddress is a paid mutator transaction binding the contract method 0x9122c344.
//
// Solidity: function updateTssAddress(address tssAddress_) returns()
func (_MuseConnectorBase *MuseConnectorBaseTransactorSession) UpdateTssAddress(tssAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorBase.Contract.UpdateTssAddress(&_MuseConnectorBase.TransactOpts, tssAddress_)
}

// MuseConnectorBasePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the MuseConnectorBase contract.
type MuseConnectorBasePausedIterator struct {
	Event *MuseConnectorBasePaused // Event containing the contract specifics and raw log

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
func (it *MuseConnectorBasePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorBasePaused)
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
		it.Event = new(MuseConnectorBasePaused)
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
func (it *MuseConnectorBasePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorBasePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorBasePaused represents a Paused event raised by the MuseConnectorBase contract.
type MuseConnectorBasePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) FilterPaused(opts *bind.FilterOpts) (*MuseConnectorBasePausedIterator, error) {

	logs, sub, err := _MuseConnectorBase.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBasePausedIterator{contract: _MuseConnectorBase.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *MuseConnectorBasePaused) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorBase.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorBasePaused)
				if err := _MuseConnectorBase.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_MuseConnectorBase *MuseConnectorBaseFilterer) ParsePaused(log types.Log) (*MuseConnectorBasePaused, error) {
	event := new(MuseConnectorBasePaused)
	if err := _MuseConnectorBase.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorBasePauserAddressUpdatedIterator is returned from FilterPauserAddressUpdated and is used to iterate over the raw logs and unpacked data for PauserAddressUpdated events raised by the MuseConnectorBase contract.
type MuseConnectorBasePauserAddressUpdatedIterator struct {
	Event *MuseConnectorBasePauserAddressUpdated // Event containing the contract specifics and raw log

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
func (it *MuseConnectorBasePauserAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorBasePauserAddressUpdated)
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
		it.Event = new(MuseConnectorBasePauserAddressUpdated)
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
func (it *MuseConnectorBasePauserAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorBasePauserAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorBasePauserAddressUpdated represents a PauserAddressUpdated event raised by the MuseConnectorBase contract.
type MuseConnectorBasePauserAddressUpdated struct {
	CallerAddress common.Address
	NewTssAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterPauserAddressUpdated is a free log retrieval operation binding the contract event 0xd41d83655d484bdf299598751c371b2d92088667266fe3774b25a97bdd5d0397.
//
// Solidity: event PauserAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) FilterPauserAddressUpdated(opts *bind.FilterOpts) (*MuseConnectorBasePauserAddressUpdatedIterator, error) {

	logs, sub, err := _MuseConnectorBase.contract.FilterLogs(opts, "PauserAddressUpdated")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBasePauserAddressUpdatedIterator{contract: _MuseConnectorBase.contract, event: "PauserAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchPauserAddressUpdated is a free log subscription operation binding the contract event 0xd41d83655d484bdf299598751c371b2d92088667266fe3774b25a97bdd5d0397.
//
// Solidity: event PauserAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) WatchPauserAddressUpdated(opts *bind.WatchOpts, sink chan<- *MuseConnectorBasePauserAddressUpdated) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorBase.contract.WatchLogs(opts, "PauserAddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorBasePauserAddressUpdated)
				if err := _MuseConnectorBase.contract.UnpackLog(event, "PauserAddressUpdated", log); err != nil {
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

// ParsePauserAddressUpdated is a log parse operation binding the contract event 0xd41d83655d484bdf299598751c371b2d92088667266fe3774b25a97bdd5d0397.
//
// Solidity: event PauserAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) ParsePauserAddressUpdated(log types.Log) (*MuseConnectorBasePauserAddressUpdated, error) {
	event := new(MuseConnectorBasePauserAddressUpdated)
	if err := _MuseConnectorBase.contract.UnpackLog(event, "PauserAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorBaseTSSAddressUpdatedIterator is returned from FilterTSSAddressUpdated and is used to iterate over the raw logs and unpacked data for TSSAddressUpdated events raised by the MuseConnectorBase contract.
type MuseConnectorBaseTSSAddressUpdatedIterator struct {
	Event *MuseConnectorBaseTSSAddressUpdated // Event containing the contract specifics and raw log

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
func (it *MuseConnectorBaseTSSAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorBaseTSSAddressUpdated)
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
		it.Event = new(MuseConnectorBaseTSSAddressUpdated)
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
func (it *MuseConnectorBaseTSSAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorBaseTSSAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorBaseTSSAddressUpdated represents a TSSAddressUpdated event raised by the MuseConnectorBase contract.
type MuseConnectorBaseTSSAddressUpdated struct {
	CallerAddress common.Address
	NewTssAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTSSAddressUpdated is a free log retrieval operation binding the contract event 0xe79965b5c67dcfb2cf5fe152715e4a7256cee62a3d5dd8484fd8a8539eb8beff.
//
// Solidity: event TSSAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) FilterTSSAddressUpdated(opts *bind.FilterOpts) (*MuseConnectorBaseTSSAddressUpdatedIterator, error) {

	logs, sub, err := _MuseConnectorBase.contract.FilterLogs(opts, "TSSAddressUpdated")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBaseTSSAddressUpdatedIterator{contract: _MuseConnectorBase.contract, event: "TSSAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchTSSAddressUpdated is a free log subscription operation binding the contract event 0xe79965b5c67dcfb2cf5fe152715e4a7256cee62a3d5dd8484fd8a8539eb8beff.
//
// Solidity: event TSSAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) WatchTSSAddressUpdated(opts *bind.WatchOpts, sink chan<- *MuseConnectorBaseTSSAddressUpdated) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorBase.contract.WatchLogs(opts, "TSSAddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorBaseTSSAddressUpdated)
				if err := _MuseConnectorBase.contract.UnpackLog(event, "TSSAddressUpdated", log); err != nil {
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

// ParseTSSAddressUpdated is a log parse operation binding the contract event 0xe79965b5c67dcfb2cf5fe152715e4a7256cee62a3d5dd8484fd8a8539eb8beff.
//
// Solidity: event TSSAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) ParseTSSAddressUpdated(log types.Log) (*MuseConnectorBaseTSSAddressUpdated, error) {
	event := new(MuseConnectorBaseTSSAddressUpdated)
	if err := _MuseConnectorBase.contract.UnpackLog(event, "TSSAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorBaseTSSAddressUpdaterUpdatedIterator is returned from FilterTSSAddressUpdaterUpdated and is used to iterate over the raw logs and unpacked data for TSSAddressUpdaterUpdated events raised by the MuseConnectorBase contract.
type MuseConnectorBaseTSSAddressUpdaterUpdatedIterator struct {
	Event *MuseConnectorBaseTSSAddressUpdaterUpdated // Event containing the contract specifics and raw log

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
func (it *MuseConnectorBaseTSSAddressUpdaterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorBaseTSSAddressUpdaterUpdated)
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
		it.Event = new(MuseConnectorBaseTSSAddressUpdaterUpdated)
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
func (it *MuseConnectorBaseTSSAddressUpdaterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorBaseTSSAddressUpdaterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorBaseTSSAddressUpdaterUpdated represents a TSSAddressUpdaterUpdated event raised by the MuseConnectorBase contract.
type MuseConnectorBaseTSSAddressUpdaterUpdated struct {
	CallerAddress        common.Address
	NewTssUpdaterAddress common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterTSSAddressUpdaterUpdated is a free log retrieval operation binding the contract event 0x5104c9abdc7d111c2aeb4ce890ac70274a4be2ee83f46a62551be5e6ebc82dd0.
//
// Solidity: event TSSAddressUpdaterUpdated(address callerAddress, address newTssUpdaterAddress)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) FilterTSSAddressUpdaterUpdated(opts *bind.FilterOpts) (*MuseConnectorBaseTSSAddressUpdaterUpdatedIterator, error) {

	logs, sub, err := _MuseConnectorBase.contract.FilterLogs(opts, "TSSAddressUpdaterUpdated")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBaseTSSAddressUpdaterUpdatedIterator{contract: _MuseConnectorBase.contract, event: "TSSAddressUpdaterUpdated", logs: logs, sub: sub}, nil
}

// WatchTSSAddressUpdaterUpdated is a free log subscription operation binding the contract event 0x5104c9abdc7d111c2aeb4ce890ac70274a4be2ee83f46a62551be5e6ebc82dd0.
//
// Solidity: event TSSAddressUpdaterUpdated(address callerAddress, address newTssUpdaterAddress)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) WatchTSSAddressUpdaterUpdated(opts *bind.WatchOpts, sink chan<- *MuseConnectorBaseTSSAddressUpdaterUpdated) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorBase.contract.WatchLogs(opts, "TSSAddressUpdaterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorBaseTSSAddressUpdaterUpdated)
				if err := _MuseConnectorBase.contract.UnpackLog(event, "TSSAddressUpdaterUpdated", log); err != nil {
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

// ParseTSSAddressUpdaterUpdated is a log parse operation binding the contract event 0x5104c9abdc7d111c2aeb4ce890ac70274a4be2ee83f46a62551be5e6ebc82dd0.
//
// Solidity: event TSSAddressUpdaterUpdated(address callerAddress, address newTssUpdaterAddress)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) ParseTSSAddressUpdaterUpdated(log types.Log) (*MuseConnectorBaseTSSAddressUpdaterUpdated, error) {
	event := new(MuseConnectorBaseTSSAddressUpdaterUpdated)
	if err := _MuseConnectorBase.contract.UnpackLog(event, "TSSAddressUpdaterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorBaseUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the MuseConnectorBase contract.
type MuseConnectorBaseUnpausedIterator struct {
	Event *MuseConnectorBaseUnpaused // Event containing the contract specifics and raw log

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
func (it *MuseConnectorBaseUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorBaseUnpaused)
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
		it.Event = new(MuseConnectorBaseUnpaused)
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
func (it *MuseConnectorBaseUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorBaseUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorBaseUnpaused represents a Unpaused event raised by the MuseConnectorBase contract.
type MuseConnectorBaseUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) FilterUnpaused(opts *bind.FilterOpts) (*MuseConnectorBaseUnpausedIterator, error) {

	logs, sub, err := _MuseConnectorBase.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBaseUnpausedIterator{contract: _MuseConnectorBase.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *MuseConnectorBaseUnpaused) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorBase.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorBaseUnpaused)
				if err := _MuseConnectorBase.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_MuseConnectorBase *MuseConnectorBaseFilterer) ParseUnpaused(log types.Log) (*MuseConnectorBaseUnpaused, error) {
	event := new(MuseConnectorBaseUnpaused)
	if err := _MuseConnectorBase.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorBaseMuseReceivedIterator is returned from FilterMuseReceived and is used to iterate over the raw logs and unpacked data for MuseReceived events raised by the MuseConnectorBase contract.
type MuseConnectorBaseMuseReceivedIterator struct {
	Event *MuseConnectorBaseMuseReceived // Event containing the contract specifics and raw log

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
func (it *MuseConnectorBaseMuseReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorBaseMuseReceived)
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
		it.Event = new(MuseConnectorBaseMuseReceived)
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
func (it *MuseConnectorBaseMuseReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorBaseMuseReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorBaseMuseReceived represents a MuseReceived event raised by the MuseConnectorBase contract.
type MuseConnectorBaseMuseReceived struct {
	MuseTxSenderAddress []byte
	SourceChainId       *big.Int
	DestinationAddress  common.Address
	MuseValue           *big.Int
	Message             []byte
	InternalSendHash    [32]byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterMuseReceived is a free log retrieval operation binding the contract event 0xf1302855733b40d8acb467ee990b6d56c05c80e28ebcabfa6e6f3f57cb50d698.
//
// Solidity: event MuseReceived(bytes museTxSenderAddress, uint256 indexed sourceChainId, address indexed destinationAddress, uint256 museValue, bytes message, bytes32 indexed internalSendHash)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) FilterMuseReceived(opts *bind.FilterOpts, sourceChainId []*big.Int, destinationAddress []common.Address, internalSendHash [][32]byte) (*MuseConnectorBaseMuseReceivedIterator, error) {

	var sourceChainIdRule []interface{}
	for _, sourceChainIdItem := range sourceChainId {
		sourceChainIdRule = append(sourceChainIdRule, sourceChainIdItem)
	}
	var destinationAddressRule []interface{}
	for _, destinationAddressItem := range destinationAddress {
		destinationAddressRule = append(destinationAddressRule, destinationAddressItem)
	}

	var internalSendHashRule []interface{}
	for _, internalSendHashItem := range internalSendHash {
		internalSendHashRule = append(internalSendHashRule, internalSendHashItem)
	}

	logs, sub, err := _MuseConnectorBase.contract.FilterLogs(opts, "MuseReceived", sourceChainIdRule, destinationAddressRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBaseMuseReceivedIterator{contract: _MuseConnectorBase.contract, event: "MuseReceived", logs: logs, sub: sub}, nil
}

// WatchMuseReceived is a free log subscription operation binding the contract event 0xf1302855733b40d8acb467ee990b6d56c05c80e28ebcabfa6e6f3f57cb50d698.
//
// Solidity: event MuseReceived(bytes museTxSenderAddress, uint256 indexed sourceChainId, address indexed destinationAddress, uint256 museValue, bytes message, bytes32 indexed internalSendHash)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) WatchMuseReceived(opts *bind.WatchOpts, sink chan<- *MuseConnectorBaseMuseReceived, sourceChainId []*big.Int, destinationAddress []common.Address, internalSendHash [][32]byte) (event.Subscription, error) {

	var sourceChainIdRule []interface{}
	for _, sourceChainIdItem := range sourceChainId {
		sourceChainIdRule = append(sourceChainIdRule, sourceChainIdItem)
	}
	var destinationAddressRule []interface{}
	for _, destinationAddressItem := range destinationAddress {
		destinationAddressRule = append(destinationAddressRule, destinationAddressItem)
	}

	var internalSendHashRule []interface{}
	for _, internalSendHashItem := range internalSendHash {
		internalSendHashRule = append(internalSendHashRule, internalSendHashItem)
	}

	logs, sub, err := _MuseConnectorBase.contract.WatchLogs(opts, "MuseReceived", sourceChainIdRule, destinationAddressRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorBaseMuseReceived)
				if err := _MuseConnectorBase.contract.UnpackLog(event, "MuseReceived", log); err != nil {
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

// ParseMuseReceived is a log parse operation binding the contract event 0xf1302855733b40d8acb467ee990b6d56c05c80e28ebcabfa6e6f3f57cb50d698.
//
// Solidity: event MuseReceived(bytes museTxSenderAddress, uint256 indexed sourceChainId, address indexed destinationAddress, uint256 museValue, bytes message, bytes32 indexed internalSendHash)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) ParseMuseReceived(log types.Log) (*MuseConnectorBaseMuseReceived, error) {
	event := new(MuseConnectorBaseMuseReceived)
	if err := _MuseConnectorBase.contract.UnpackLog(event, "MuseReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorBaseMuseRevertedIterator is returned from FilterMuseReverted and is used to iterate over the raw logs and unpacked data for MuseReverted events raised by the MuseConnectorBase contract.
type MuseConnectorBaseMuseRevertedIterator struct {
	Event *MuseConnectorBaseMuseReverted // Event containing the contract specifics and raw log

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
func (it *MuseConnectorBaseMuseRevertedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorBaseMuseReverted)
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
		it.Event = new(MuseConnectorBaseMuseReverted)
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
func (it *MuseConnectorBaseMuseRevertedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorBaseMuseRevertedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorBaseMuseReverted represents a MuseReverted event raised by the MuseConnectorBase contract.
type MuseConnectorBaseMuseReverted struct {
	MuseTxSenderAddress common.Address
	SourceChainId       *big.Int
	DestinationChainId  *big.Int
	DestinationAddress  []byte
	RemainingMuseValue  *big.Int
	Message             []byte
	InternalSendHash    [32]byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterMuseReverted is a free log retrieval operation binding the contract event 0x521fb0b407c2eb9b1375530e9b9a569889992140a688bc076aa72c1712012c88.
//
// Solidity: event MuseReverted(address museTxSenderAddress, uint256 sourceChainId, uint256 indexed destinationChainId, bytes destinationAddress, uint256 remainingMuseValue, bytes message, bytes32 indexed internalSendHash)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) FilterMuseReverted(opts *bind.FilterOpts, destinationChainId []*big.Int, internalSendHash [][32]byte) (*MuseConnectorBaseMuseRevertedIterator, error) {

	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	var internalSendHashRule []interface{}
	for _, internalSendHashItem := range internalSendHash {
		internalSendHashRule = append(internalSendHashRule, internalSendHashItem)
	}

	logs, sub, err := _MuseConnectorBase.contract.FilterLogs(opts, "MuseReverted", destinationChainIdRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBaseMuseRevertedIterator{contract: _MuseConnectorBase.contract, event: "MuseReverted", logs: logs, sub: sub}, nil
}

// WatchMuseReverted is a free log subscription operation binding the contract event 0x521fb0b407c2eb9b1375530e9b9a569889992140a688bc076aa72c1712012c88.
//
// Solidity: event MuseReverted(address museTxSenderAddress, uint256 sourceChainId, uint256 indexed destinationChainId, bytes destinationAddress, uint256 remainingMuseValue, bytes message, bytes32 indexed internalSendHash)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) WatchMuseReverted(opts *bind.WatchOpts, sink chan<- *MuseConnectorBaseMuseReverted, destinationChainId []*big.Int, internalSendHash [][32]byte) (event.Subscription, error) {

	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	var internalSendHashRule []interface{}
	for _, internalSendHashItem := range internalSendHash {
		internalSendHashRule = append(internalSendHashRule, internalSendHashItem)
	}

	logs, sub, err := _MuseConnectorBase.contract.WatchLogs(opts, "MuseReverted", destinationChainIdRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorBaseMuseReverted)
				if err := _MuseConnectorBase.contract.UnpackLog(event, "MuseReverted", log); err != nil {
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

// ParseMuseReverted is a log parse operation binding the contract event 0x521fb0b407c2eb9b1375530e9b9a569889992140a688bc076aa72c1712012c88.
//
// Solidity: event MuseReverted(address museTxSenderAddress, uint256 sourceChainId, uint256 indexed destinationChainId, bytes destinationAddress, uint256 remainingMuseValue, bytes message, bytes32 indexed internalSendHash)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) ParseMuseReverted(log types.Log) (*MuseConnectorBaseMuseReverted, error) {
	event := new(MuseConnectorBaseMuseReverted)
	if err := _MuseConnectorBase.contract.UnpackLog(event, "MuseReverted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorBaseMuseSentIterator is returned from FilterMuseSent and is used to iterate over the raw logs and unpacked data for MuseSent events raised by the MuseConnectorBase contract.
type MuseConnectorBaseMuseSentIterator struct {
	Event *MuseConnectorBaseMuseSent // Event containing the contract specifics and raw log

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
func (it *MuseConnectorBaseMuseSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorBaseMuseSent)
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
		it.Event = new(MuseConnectorBaseMuseSent)
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
func (it *MuseConnectorBaseMuseSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorBaseMuseSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorBaseMuseSent represents a MuseSent event raised by the MuseConnectorBase contract.
type MuseConnectorBaseMuseSent struct {
	SourceTxOriginAddress common.Address
	MuseTxSenderAddress   common.Address
	DestinationChainId    *big.Int
	DestinationAddress    []byte
	MuseValueAndGas       *big.Int
	DestinationGasLimit   *big.Int
	Message               []byte
	MuseParams            []byte
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterMuseSent is a free log retrieval operation binding the contract event 0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4.
//
// Solidity: event MuseSent(address sourceTxOriginAddress, address indexed museTxSenderAddress, uint256 indexed destinationChainId, bytes destinationAddress, uint256 museValueAndGas, uint256 destinationGasLimit, bytes message, bytes museParams)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) FilterMuseSent(opts *bind.FilterOpts, museTxSenderAddress []common.Address, destinationChainId []*big.Int) (*MuseConnectorBaseMuseSentIterator, error) {

	var museTxSenderAddressRule []interface{}
	for _, museTxSenderAddressItem := range museTxSenderAddress {
		museTxSenderAddressRule = append(museTxSenderAddressRule, museTxSenderAddressItem)
	}
	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	logs, sub, err := _MuseConnectorBase.contract.FilterLogs(opts, "MuseSent", museTxSenderAddressRule, destinationChainIdRule)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorBaseMuseSentIterator{contract: _MuseConnectorBase.contract, event: "MuseSent", logs: logs, sub: sub}, nil
}

// WatchMuseSent is a free log subscription operation binding the contract event 0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4.
//
// Solidity: event MuseSent(address sourceTxOriginAddress, address indexed museTxSenderAddress, uint256 indexed destinationChainId, bytes destinationAddress, uint256 museValueAndGas, uint256 destinationGasLimit, bytes message, bytes museParams)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) WatchMuseSent(opts *bind.WatchOpts, sink chan<- *MuseConnectorBaseMuseSent, museTxSenderAddress []common.Address, destinationChainId []*big.Int) (event.Subscription, error) {

	var museTxSenderAddressRule []interface{}
	for _, museTxSenderAddressItem := range museTxSenderAddress {
		museTxSenderAddressRule = append(museTxSenderAddressRule, museTxSenderAddressItem)
	}
	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	logs, sub, err := _MuseConnectorBase.contract.WatchLogs(opts, "MuseSent", museTxSenderAddressRule, destinationChainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorBaseMuseSent)
				if err := _MuseConnectorBase.contract.UnpackLog(event, "MuseSent", log); err != nil {
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

// ParseMuseSent is a log parse operation binding the contract event 0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4.
//
// Solidity: event MuseSent(address sourceTxOriginAddress, address indexed museTxSenderAddress, uint256 indexed destinationChainId, bytes destinationAddress, uint256 museValueAndGas, uint256 destinationGasLimit, bytes message, bytes museParams)
func (_MuseConnectorBase *MuseConnectorBaseFilterer) ParseMuseSent(log types.Log) (*MuseConnectorBaseMuseSent, error) {
	event := new(MuseConnectorBaseMuseSent)
	if err := _MuseConnectorBase.contract.UnpackLog(event, "MuseSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
