// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testmrc20

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

// TestMRC20MetaData contains all meta data concerning the TestMRC20 contract.
var TestMRC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainid_\",\"type\":\"uint256\"},{\"internalType\":\"enumCoinType\",\"name\":\"coinType_\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CallerIsNotFungibleModule\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasFeeTransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LowAllowance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"LowBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroGasCoin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroGasPrice\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"from\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"UpdatedGasLimit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"protocolFlatFee\",\"type\":\"uint256\"}],\"name\":\"UpdatedProtocolFlatFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"systemContract\",\"type\":\"address\"}],\"name\":\"UpdatedSystemContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasfee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"protocolFlatFee\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CHAIN_ID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"COIN_TYPE\",\"outputs\":[{\"internalType\":\"enumCoinType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FUNGIBLE_MODULE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GAS_LIMIT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PROTOCOL_FLAT_FEE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SYSTEM_CONTRACT_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gatewayAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newField\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"newPublicField\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"}],\"name\":\"updateGasLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newField_\",\"type\":\"uint256\"}],\"name\":\"updateNewField\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"protocolFlatFee\",\"type\":\"uint256\"}],\"name\":\"updateProtocolFlatFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"updateSystemContractAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"to\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawGasFee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c06040523480156200001157600080fd5b5060405162002644380380620026448339818101604052810190620000379190620000aa565b8160808181525050806002811115620000555762000054620000fb565b5b60a08160028111156200006d576200006c620000fb565b5b60f81b8152505050506200015a565b6000815190506200008d816200012f565b92915050565b600081519050620000a48162000140565b92915050565b60008060408385031215620000c457620000c36200012a565b5b6000620000d48582860162000093565b9250506020620000e7858286016200007c565b9150509250929050565b6000819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b600080fd5b600381106200013d57600080fd5b50565b6200014b81620000f1565b81146200015757600080fd5b50565b60805160a05160f81c6124b3620001916000396000610b03015260008181610a2701528181610fc701526110fc01526124b36000f3fe608060405234801561001057600080fd5b50600436106101c45760003560e01c80638b851b95116100f9578063c701262611610097578063dd62ed3e11610071578063dd62ed3e14610538578063eddeb12314610568578063f2441b3214610584578063f687d12a146105a2576101c4565b8063c7012626146104cd578063c835d7cc146104fd578063d9eeebed14610519576101c4565b8063a457c2d7116100d3578063a457c2d714610431578063a7605f4514610461578063a9059cbb1461047f578063b92894ba146104af576101c4565b80638b851b95146103d757806395d89b41146103f5578063a3413d0314610413576101c4565b80633ce4a5bc116101665780634d8943bb116101405780634d8943bb1461034f57806370a082311461036d578063732bb0e41461039d57806385e1f4d0146103b9576101c4565b80633ce4a5bc146102d157806342966c68146102ef57806347e7ef241461031f576101c4565b806318160ddd116101a257806318160ddd1461023557806323b872dd14610253578063313ce5671461028357806339509351146102a1576101c4565b806306fdde03146101c9578063091d2788146101e7578063095ea7b314610205575b600080fd5b6101d16105be565b6040516101de9190612029565b60405180910390f35b6101ef610650565b6040516101fc919061204b565b60405180910390f35b61021f600480360381019061021a9190611cea565b610656565b60405161022c9190611f77565b60405180910390f35b61023d610674565b60405161024a919061204b565b60405180910390f35b61026d60048036038101906102689190611c97565b61067e565b60405161027a9190611f77565b60405180910390f35b61028b610776565b6040516102989190612066565b60405180910390f35b6102bb60048036038101906102b69190611cea565b61078d565b6040516102c89190611f77565b60405180910390f35b6102d9610833565b6040516102e69190611efc565b60405180910390f35b61030960048036038101906103049190611db3565b61084b565b6040516103169190611f77565b60405180910390f35b61033960048036038101906103349190611cea565b610860565b6040516103469190611f77565b60405180910390f35b6103576109cc565b604051610364919061204b565b60405180910390f35b61038760048036038101906103829190611bfd565b6109d2565b604051610394919061204b565b60405180910390f35b6103b760048036038101906103b29190611db3565b610a1b565b005b6103c1610a25565b6040516103ce919061204b565b60405180910390f35b6103df610a49565b6040516103ec9190611efc565b60405180910390f35b6103fd610a6f565b60405161040a9190612029565b60405180910390f35b61041b610b01565b604051610428919061200e565b60405180910390f35b61044b60048036038101906104469190611cea565b610b25565b6040516104589190611f77565b60405180910390f35b610469610c88565b604051610476919061204b565b60405180910390f35b61049960048036038101906104949190611cea565b610c8e565b6040516104a69190611f77565b60405180910390f35b6104b7610cac565b6040516104c49190612029565b60405180910390f35b6104e760048036038101906104e29190611d57565b610d3a565b6040516104f49190611f77565b60405180910390f35b61051760048036038101906105129190611bfd565b610e90565b005b610521610f83565b60405161052f929190611f4e565b60405180910390f35b610552600480360381019061054d9190611c57565b6111f0565b60405161055f919061204b565b60405180910390f35b610582600480360381019061057d9190611db3565b611277565b005b61058c611331565b6040516105999190611efc565b60405180910390f35b6105bc60048036038101906105b79190611db3565b611355565b005b6060600680546105cd906122af565b80601f01602080910402602001604051908101604052809291908181526020018280546105f9906122af565b80156106465780601f1061061b57610100808354040283529160200191610646565b820191906000526020600020905b81548152906001019060200180831161062957829003601f168201915b5050505050905090565b60015481565b600061066a61066361140f565b8484611417565b6001905092915050565b6000600554905090565b600061068b8484846115d0565b6000600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006106d661140f565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508281101561074d576040517f10bad14700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61076a8561075961140f565b858461076591906121bf565b611417565b60019150509392505050565b6000600860009054906101000a900460ff16905090565b600081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006107d961140f565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610822919061210f565b925050819055506001905092915050565b73735b14bb79463307aacbed86daf3322b1e6226ab81565b6000610857338361182c565b60019050919050565b600073735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141580156108fe575060008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614155b15610935576040517fddb5de5e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61093f83836119e4565b8273ffffffffffffffffffffffffffffffffffffffff167f67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab373735b14bb79463307aacbed86daf3322b1e6226ab60405160200161099c9190611ee1565b604051602081830303815290604052846040516109ba929190611f92565b60405180910390a26001905092915050565b60025481565b6000600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b8060098190555050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600860019054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b606060078054610a7e906122af565b80601f0160208091040260200160405190810160405280929190818152602001828054610aaa906122af565b8015610af75780601f10610acc57610100808354040283529160200191610af7565b820191906000526020600020905b815481529060010190602001808311610ada57829003601f168201915b5050505050905090565b7f000000000000000000000000000000000000000000000000000000000000000081565b600081600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000610b7161140f565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015610be4576040517f10bad14700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000610c2e61140f565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610c7791906121bf565b925050819055506001905092915050565b60095481565b6000610ca2610c9b61140f565b84846115d0565b6001905092915050565b600a8054610cb9906122af565b80601f0160208091040260200160405190810160405280929190818152602001828054610ce5906122af565b8015610d325780601f10610d0757610100808354040283529160200191610d32565b820191906000526020600020905b815481529060010190602001808311610d1557829003601f168201915b505050505081565b6000806000610d47610f83565b915091508173ffffffffffffffffffffffffffffffffffffffff166323b872dd3373735b14bb79463307aacbed86daf3322b1e6226ab846040518463ffffffff1660e01b8152600401610d9c93929190611f17565b602060405180830381600087803b158015610db657600080fd5b505af1158015610dca573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610dee9190611d2a565b610e24576040517f0a7cd6d600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e2e338561182c565b3373ffffffffffffffffffffffffffffffffffffffff167f9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955868684600254604051610e7c9493929190611fc2565b60405180910390a260019250505092915050565b73735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610f09576040517f2b2add3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507fd55614e962c5fd6ece71614f6348d702468a997a394dd5e5c1677950226d97ae81604051610f789190611efc565b60405180910390a150565b60008060008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16630be155477f00000000000000000000000000000000000000000000000000000000000000006040518263ffffffff1660e01b8152600401611002919061204b565b60206040518083038186803b15801561101a57600080fd5b505afa15801561102e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110529190611c2a565b9050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156110bb576040517f78fff39600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663d7fd7afb7f00000000000000000000000000000000000000000000000000000000000000006040518263ffffffff1660e01b8152600401611137919061204b565b60206040518083038186803b15801561114f57600080fd5b505afa158015611163573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111879190611de0565b905060008114156111c4576040517fe661aed000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600254600154836111d79190612165565b6111e1919061210f565b90508281945094505050509091565b6000600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b73735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146112f0576040517f2b2add3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806002819055507fef13af88e424b5d15f49c77758542c1938b08b8b95b91ed0751f98ba99000d8f81604051611326919061204b565b60405180910390a150565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b73735b14bb79463307aacbed86daf3322b1e6226ab73ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146113ce576040517f2b2add3d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b806001819055507fff5788270f43bfc1ca41c503606d2594aa3023a1a7547de403a3e2f146a4a80a81604051611404919061204b565b60405180910390a150565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141561147e576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156114e5576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040516115c3919061204b565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415611637576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561169e576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490508181101561171c576040517ffe382aa700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818161172891906121bf565b600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555081600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546117ba919061210f565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161181e919061204b565b60405180910390a350505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611893576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6000600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015611911576040517ffe382aa700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818161191d91906121bf565b600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816005600082825461197291906121bf565b92505081905550600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516119d7919061204b565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415611a4b576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060056000828254611a5d919061210f565b9250508190555080600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254611ab3919061210f565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051611b18919061204b565b60405180910390a35050565b6000611b37611b32846120a6565b612081565b905082815260208101848484011115611b5357611b526123f7565b5b611b5e84828561226d565b509392505050565b600081359050611b7581612438565b92915050565b600081519050611b8a81612438565b92915050565b600081519050611b9f8161244f565b92915050565b600082601f830112611bba57611bb96123f2565b5b8135611bca848260208601611b24565b91505092915050565b600081359050611be281612466565b92915050565b600081519050611bf781612466565b92915050565b600060208284031215611c1357611c12612401565b5b6000611c2184828501611b66565b91505092915050565b600060208284031215611c4057611c3f612401565b5b6000611c4e84828501611b7b565b91505092915050565b60008060408385031215611c6e57611c6d612401565b5b6000611c7c85828601611b66565b9250506020611c8d85828601611b66565b9150509250929050565b600080600060608486031215611cb057611caf612401565b5b6000611cbe86828701611b66565b9350506020611ccf86828701611b66565b9250506040611ce086828701611bd3565b9150509250925092565b60008060408385031215611d0157611d00612401565b5b6000611d0f85828601611b66565b9250506020611d2085828601611bd3565b9150509250929050565b600060208284031215611d4057611d3f612401565b5b6000611d4e84828501611b90565b91505092915050565b60008060408385031215611d6e57611d6d612401565b5b600083013567ffffffffffffffff811115611d8c57611d8b6123fc565b5b611d9885828601611ba5565b9250506020611da985828601611bd3565b9150509250929050565b600060208284031215611dc957611dc8612401565b5b6000611dd784828501611bd3565b91505092915050565b600060208284031215611df657611df5612401565b5b6000611e0484828501611be8565b91505092915050565b611e16816121f3565b82525050565b611e2d611e28826121f3565b612312565b82525050565b611e3c81612205565b82525050565b6000611e4d826120d7565b611e5781856120ed565b9350611e6781856020860161227c565b611e7081612406565b840191505092915050565b611e848161225b565b82525050565b6000611e95826120e2565b611e9f81856120fe565b9350611eaf81856020860161227c565b611eb881612406565b840191505092915050565b611ecc81612244565b82525050565b611edb8161224e565b82525050565b6000611eed8284611e1c565b60148201915081905092915050565b6000602082019050611f116000830184611e0d565b92915050565b6000606082019050611f2c6000830186611e0d565b611f396020830185611e0d565b611f466040830184611ec3565b949350505050565b6000604082019050611f636000830185611e0d565b611f706020830184611ec3565b9392505050565b6000602082019050611f8c6000830184611e33565b92915050565b60006040820190508181036000830152611fac8185611e42565b9050611fbb6020830184611ec3565b9392505050565b60006080820190508181036000830152611fdc8187611e42565b9050611feb6020830186611ec3565b611ff86040830185611ec3565b6120056060830184611ec3565b95945050505050565b60006020820190506120236000830184611e7b565b92915050565b600060208201905081810360008301526120438184611e8a565b905092915050565b60006020820190506120606000830184611ec3565b92915050565b600060208201905061207b6000830184611ed2565b92915050565b600061208b61209c565b905061209782826122e1565b919050565b6000604051905090565b600067ffffffffffffffff8211156120c1576120c06123c3565b5b6120ca82612406565b9050602081019050919050565b600081519050919050565b600081519050919050565b600082825260208201905092915050565b600082825260208201905092915050565b600061211a82612244565b915061212583612244565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561215a57612159612336565b5b828201905092915050565b600061217082612244565b915061217b83612244565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04831182151516156121b4576121b3612336565b5b828202905092915050565b60006121ca82612244565b91506121d583612244565b9250828210156121e8576121e7612336565b5b828203905092915050565b60006121fe82612224565b9050919050565b60008115159050919050565b600081905061221f82612424565b919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b600060ff82169050919050565b600061226682612211565b9050919050565b82818337600083830152505050565b60005b8381101561229a57808201518184015260208101905061227f565b838111156122a9576000848401525b50505050565b600060028204905060018216806122c757607f821691505b602082108114156122db576122da612394565b5b50919050565b6122ea82612406565b810181811067ffffffffffffffff82111715612309576123086123c3565b5b80604052505050565b600061231d82612324565b9050919050565b600061232f82612417565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b60008160601b9050919050565b6003811061243557612434612365565b5b50565b612441816121f3565b811461244c57600080fd5b50565b61245881612205565b811461246357600080fd5b50565b61246f81612244565b811461247a57600080fd5b5056fea2646970667358221220b8278837e775bf149b356e2bff0f3e4931577c5c8ff848545ce78afb4e0ce8cb64736f6c63430008070033",
}

// TestMRC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use TestMRC20MetaData.ABI instead.
var TestMRC20ABI = TestMRC20MetaData.ABI

// TestMRC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestMRC20MetaData.Bin instead.
var TestMRC20Bin = TestMRC20MetaData.Bin

// DeployTestMRC20 deploys a new Ethereum contract, binding an instance of TestMRC20 to it.
func DeployTestMRC20(auth *bind.TransactOpts, backend bind.ContractBackend, chainid_ *big.Int, coinType_ uint8) (common.Address, *types.Transaction, *TestMRC20, error) {
	parsed, err := TestMRC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestMRC20Bin), backend, chainid_, coinType_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestMRC20{TestMRC20Caller: TestMRC20Caller{contract: contract}, TestMRC20Transactor: TestMRC20Transactor{contract: contract}, TestMRC20Filterer: TestMRC20Filterer{contract: contract}}, nil
}

// TestMRC20 is an auto generated Go binding around an Ethereum contract.
type TestMRC20 struct {
	TestMRC20Caller     // Read-only binding to the contract
	TestMRC20Transactor // Write-only binding to the contract
	TestMRC20Filterer   // Log filterer for contract events
}

// TestMRC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type TestMRC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestMRC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type TestMRC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestMRC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestMRC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestMRC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestMRC20Session struct {
	Contract     *TestMRC20        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestMRC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestMRC20CallerSession struct {
	Contract *TestMRC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TestMRC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestMRC20TransactorSession struct {
	Contract     *TestMRC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TestMRC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type TestMRC20Raw struct {
	Contract *TestMRC20 // Generic contract binding to access the raw methods on
}

// TestMRC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestMRC20CallerRaw struct {
	Contract *TestMRC20Caller // Generic read-only contract binding to access the raw methods on
}

// TestMRC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestMRC20TransactorRaw struct {
	Contract *TestMRC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTestMRC20 creates a new instance of TestMRC20, bound to a specific deployed contract.
func NewTestMRC20(address common.Address, backend bind.ContractBackend) (*TestMRC20, error) {
	contract, err := bindTestMRC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestMRC20{TestMRC20Caller: TestMRC20Caller{contract: contract}, TestMRC20Transactor: TestMRC20Transactor{contract: contract}, TestMRC20Filterer: TestMRC20Filterer{contract: contract}}, nil
}

// NewTestMRC20Caller creates a new read-only instance of TestMRC20, bound to a specific deployed contract.
func NewTestMRC20Caller(address common.Address, caller bind.ContractCaller) (*TestMRC20Caller, error) {
	contract, err := bindTestMRC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestMRC20Caller{contract: contract}, nil
}

// NewTestMRC20Transactor creates a new write-only instance of TestMRC20, bound to a specific deployed contract.
func NewTestMRC20Transactor(address common.Address, transactor bind.ContractTransactor) (*TestMRC20Transactor, error) {
	contract, err := bindTestMRC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestMRC20Transactor{contract: contract}, nil
}

// NewTestMRC20Filterer creates a new log filterer instance of TestMRC20, bound to a specific deployed contract.
func NewTestMRC20Filterer(address common.Address, filterer bind.ContractFilterer) (*TestMRC20Filterer, error) {
	contract, err := bindTestMRC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestMRC20Filterer{contract: contract}, nil
}

// bindTestMRC20 binds a generic wrapper to an already deployed contract.
func bindTestMRC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TestMRC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestMRC20 *TestMRC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestMRC20.Contract.TestMRC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestMRC20 *TestMRC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestMRC20.Contract.TestMRC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestMRC20 *TestMRC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestMRC20.Contract.TestMRC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestMRC20 *TestMRC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestMRC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestMRC20 *TestMRC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestMRC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestMRC20 *TestMRC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestMRC20.Contract.contract.Transact(opts, method, params...)
}

// CHAINID is a free data retrieval call binding the contract method 0x85e1f4d0.
//
// Solidity: function CHAIN_ID() view returns(uint256)
func (_TestMRC20 *TestMRC20Caller) CHAINID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "CHAIN_ID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CHAINID is a free data retrieval call binding the contract method 0x85e1f4d0.
//
// Solidity: function CHAIN_ID() view returns(uint256)
func (_TestMRC20 *TestMRC20Session) CHAINID() (*big.Int, error) {
	return _TestMRC20.Contract.CHAINID(&_TestMRC20.CallOpts)
}

// CHAINID is a free data retrieval call binding the contract method 0x85e1f4d0.
//
// Solidity: function CHAIN_ID() view returns(uint256)
func (_TestMRC20 *TestMRC20CallerSession) CHAINID() (*big.Int, error) {
	return _TestMRC20.Contract.CHAINID(&_TestMRC20.CallOpts)
}

// COINTYPE is a free data retrieval call binding the contract method 0xa3413d03.
//
// Solidity: function COIN_TYPE() view returns(uint8)
func (_TestMRC20 *TestMRC20Caller) COINTYPE(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "COIN_TYPE")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// COINTYPE is a free data retrieval call binding the contract method 0xa3413d03.
//
// Solidity: function COIN_TYPE() view returns(uint8)
func (_TestMRC20 *TestMRC20Session) COINTYPE() (uint8, error) {
	return _TestMRC20.Contract.COINTYPE(&_TestMRC20.CallOpts)
}

// COINTYPE is a free data retrieval call binding the contract method 0xa3413d03.
//
// Solidity: function COIN_TYPE() view returns(uint8)
func (_TestMRC20 *TestMRC20CallerSession) COINTYPE() (uint8, error) {
	return _TestMRC20.Contract.COINTYPE(&_TestMRC20.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_TestMRC20 *TestMRC20Caller) FUNGIBLEMODULEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "FUNGIBLE_MODULE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_TestMRC20 *TestMRC20Session) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _TestMRC20.Contract.FUNGIBLEMODULEADDRESS(&_TestMRC20.CallOpts)
}

// FUNGIBLEMODULEADDRESS is a free data retrieval call binding the contract method 0x3ce4a5bc.
//
// Solidity: function FUNGIBLE_MODULE_ADDRESS() view returns(address)
func (_TestMRC20 *TestMRC20CallerSession) FUNGIBLEMODULEADDRESS() (common.Address, error) {
	return _TestMRC20.Contract.FUNGIBLEMODULEADDRESS(&_TestMRC20.CallOpts)
}

// GASLIMIT is a free data retrieval call binding the contract method 0x091d2788.
//
// Solidity: function GAS_LIMIT() view returns(uint256)
func (_TestMRC20 *TestMRC20Caller) GASLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "GAS_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GASLIMIT is a free data retrieval call binding the contract method 0x091d2788.
//
// Solidity: function GAS_LIMIT() view returns(uint256)
func (_TestMRC20 *TestMRC20Session) GASLIMIT() (*big.Int, error) {
	return _TestMRC20.Contract.GASLIMIT(&_TestMRC20.CallOpts)
}

// GASLIMIT is a free data retrieval call binding the contract method 0x091d2788.
//
// Solidity: function GAS_LIMIT() view returns(uint256)
func (_TestMRC20 *TestMRC20CallerSession) GASLIMIT() (*big.Int, error) {
	return _TestMRC20.Contract.GASLIMIT(&_TestMRC20.CallOpts)
}

// PROTOCOLFLATFEE is a free data retrieval call binding the contract method 0x4d8943bb.
//
// Solidity: function PROTOCOL_FLAT_FEE() view returns(uint256)
func (_TestMRC20 *TestMRC20Caller) PROTOCOLFLATFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "PROTOCOL_FLAT_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PROTOCOLFLATFEE is a free data retrieval call binding the contract method 0x4d8943bb.
//
// Solidity: function PROTOCOL_FLAT_FEE() view returns(uint256)
func (_TestMRC20 *TestMRC20Session) PROTOCOLFLATFEE() (*big.Int, error) {
	return _TestMRC20.Contract.PROTOCOLFLATFEE(&_TestMRC20.CallOpts)
}

// PROTOCOLFLATFEE is a free data retrieval call binding the contract method 0x4d8943bb.
//
// Solidity: function PROTOCOL_FLAT_FEE() view returns(uint256)
func (_TestMRC20 *TestMRC20CallerSession) PROTOCOLFLATFEE() (*big.Int, error) {
	return _TestMRC20.Contract.PROTOCOLFLATFEE(&_TestMRC20.CallOpts)
}

// SYSTEMCONTRACTADDRESS is a free data retrieval call binding the contract method 0xf2441b32.
//
// Solidity: function SYSTEM_CONTRACT_ADDRESS() view returns(address)
func (_TestMRC20 *TestMRC20Caller) SYSTEMCONTRACTADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "SYSTEM_CONTRACT_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SYSTEMCONTRACTADDRESS is a free data retrieval call binding the contract method 0xf2441b32.
//
// Solidity: function SYSTEM_CONTRACT_ADDRESS() view returns(address)
func (_TestMRC20 *TestMRC20Session) SYSTEMCONTRACTADDRESS() (common.Address, error) {
	return _TestMRC20.Contract.SYSTEMCONTRACTADDRESS(&_TestMRC20.CallOpts)
}

// SYSTEMCONTRACTADDRESS is a free data retrieval call binding the contract method 0xf2441b32.
//
// Solidity: function SYSTEM_CONTRACT_ADDRESS() view returns(address)
func (_TestMRC20 *TestMRC20CallerSession) SYSTEMCONTRACTADDRESS() (common.Address, error) {
	return _TestMRC20.Contract.SYSTEMCONTRACTADDRESS(&_TestMRC20.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestMRC20 *TestMRC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestMRC20 *TestMRC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TestMRC20.Contract.Allowance(&_TestMRC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestMRC20 *TestMRC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TestMRC20.Contract.Allowance(&_TestMRC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestMRC20 *TestMRC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestMRC20 *TestMRC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _TestMRC20.Contract.BalanceOf(&_TestMRC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestMRC20 *TestMRC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TestMRC20.Contract.BalanceOf(&_TestMRC20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestMRC20 *TestMRC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestMRC20 *TestMRC20Session) Decimals() (uint8, error) {
	return _TestMRC20.Contract.Decimals(&_TestMRC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestMRC20 *TestMRC20CallerSession) Decimals() (uint8, error) {
	return _TestMRC20.Contract.Decimals(&_TestMRC20.CallOpts)
}

// GatewayAddress is a free data retrieval call binding the contract method 0x8b851b95.
//
// Solidity: function gatewayAddress() view returns(address)
func (_TestMRC20 *TestMRC20Caller) GatewayAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "gatewayAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GatewayAddress is a free data retrieval call binding the contract method 0x8b851b95.
//
// Solidity: function gatewayAddress() view returns(address)
func (_TestMRC20 *TestMRC20Session) GatewayAddress() (common.Address, error) {
	return _TestMRC20.Contract.GatewayAddress(&_TestMRC20.CallOpts)
}

// GatewayAddress is a free data retrieval call binding the contract method 0x8b851b95.
//
// Solidity: function gatewayAddress() view returns(address)
func (_TestMRC20 *TestMRC20CallerSession) GatewayAddress() (common.Address, error) {
	return _TestMRC20.Contract.GatewayAddress(&_TestMRC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestMRC20 *TestMRC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestMRC20 *TestMRC20Session) Name() (string, error) {
	return _TestMRC20.Contract.Name(&_TestMRC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestMRC20 *TestMRC20CallerSession) Name() (string, error) {
	return _TestMRC20.Contract.Name(&_TestMRC20.CallOpts)
}

// NewField is a free data retrieval call binding the contract method 0xa7605f45.
//
// Solidity: function newField() view returns(uint256)
func (_TestMRC20 *TestMRC20Caller) NewField(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "newField")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NewField is a free data retrieval call binding the contract method 0xa7605f45.
//
// Solidity: function newField() view returns(uint256)
func (_TestMRC20 *TestMRC20Session) NewField() (*big.Int, error) {
	return _TestMRC20.Contract.NewField(&_TestMRC20.CallOpts)
}

// NewField is a free data retrieval call binding the contract method 0xa7605f45.
//
// Solidity: function newField() view returns(uint256)
func (_TestMRC20 *TestMRC20CallerSession) NewField() (*big.Int, error) {
	return _TestMRC20.Contract.NewField(&_TestMRC20.CallOpts)
}

// NewPublicField is a free data retrieval call binding the contract method 0xb92894ba.
//
// Solidity: function newPublicField() view returns(string)
func (_TestMRC20 *TestMRC20Caller) NewPublicField(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "newPublicField")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NewPublicField is a free data retrieval call binding the contract method 0xb92894ba.
//
// Solidity: function newPublicField() view returns(string)
func (_TestMRC20 *TestMRC20Session) NewPublicField() (string, error) {
	return _TestMRC20.Contract.NewPublicField(&_TestMRC20.CallOpts)
}

// NewPublicField is a free data retrieval call binding the contract method 0xb92894ba.
//
// Solidity: function newPublicField() view returns(string)
func (_TestMRC20 *TestMRC20CallerSession) NewPublicField() (string, error) {
	return _TestMRC20.Contract.NewPublicField(&_TestMRC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestMRC20 *TestMRC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestMRC20 *TestMRC20Session) Symbol() (string, error) {
	return _TestMRC20.Contract.Symbol(&_TestMRC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestMRC20 *TestMRC20CallerSession) Symbol() (string, error) {
	return _TestMRC20.Contract.Symbol(&_TestMRC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestMRC20 *TestMRC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestMRC20 *TestMRC20Session) TotalSupply() (*big.Int, error) {
	return _TestMRC20.Contract.TotalSupply(&_TestMRC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestMRC20 *TestMRC20CallerSession) TotalSupply() (*big.Int, error) {
	return _TestMRC20.Contract.TotalSupply(&_TestMRC20.CallOpts)
}

// WithdrawGasFee is a free data retrieval call binding the contract method 0xd9eeebed.
//
// Solidity: function withdrawGasFee() view returns(address, uint256)
func (_TestMRC20 *TestMRC20Caller) WithdrawGasFee(opts *bind.CallOpts) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _TestMRC20.contract.Call(opts, &out, "withdrawGasFee")

	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// WithdrawGasFee is a free data retrieval call binding the contract method 0xd9eeebed.
//
// Solidity: function withdrawGasFee() view returns(address, uint256)
func (_TestMRC20 *TestMRC20Session) WithdrawGasFee() (common.Address, *big.Int, error) {
	return _TestMRC20.Contract.WithdrawGasFee(&_TestMRC20.CallOpts)
}

// WithdrawGasFee is a free data retrieval call binding the contract method 0xd9eeebed.
//
// Solidity: function withdrawGasFee() view returns(address, uint256)
func (_TestMRC20 *TestMRC20CallerSession) WithdrawGasFee() (common.Address, *big.Int, error) {
	return _TestMRC20.Contract.WithdrawGasFee(&_TestMRC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.Approve(&_TestMRC20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.Approve(&_TestMRC20.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Transactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Session) Burn(amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.Burn(&_TestMRC20.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20TransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.Burn(&_TestMRC20.TransactOpts, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "decreaseAllowance", spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Session) DecreaseAllowance(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.DecreaseAllowance(&_TestMRC20.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20TransactorSession) DecreaseAllowance(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.DecreaseAllowance(&_TestMRC20.TransactOpts, spender, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address to, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Transactor) Deposit(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "deposit", to, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address to, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Session) Deposit(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.Deposit(&_TestMRC20.TransactOpts, to, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address to, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20TransactorSession) Deposit(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.Deposit(&_TestMRC20.TransactOpts, to, amount)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "increaseAllowance", spender, amount)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Session) IncreaseAllowance(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.IncreaseAllowance(&_TestMRC20.TransactOpts, spender, amount)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20TransactorSession) IncreaseAllowance(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.IncreaseAllowance(&_TestMRC20.TransactOpts, spender, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.Transfer(&_TestMRC20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.Transfer(&_TestMRC20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.TransferFrom(&_TestMRC20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.TransferFrom(&_TestMRC20.TransactOpts, sender, recipient, amount)
}

// UpdateGasLimit is a paid mutator transaction binding the contract method 0xf687d12a.
//
// Solidity: function updateGasLimit(uint256 gasLimit) returns()
func (_TestMRC20 *TestMRC20Transactor) UpdateGasLimit(opts *bind.TransactOpts, gasLimit *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "updateGasLimit", gasLimit)
}

// UpdateGasLimit is a paid mutator transaction binding the contract method 0xf687d12a.
//
// Solidity: function updateGasLimit(uint256 gasLimit) returns()
func (_TestMRC20 *TestMRC20Session) UpdateGasLimit(gasLimit *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.UpdateGasLimit(&_TestMRC20.TransactOpts, gasLimit)
}

// UpdateGasLimit is a paid mutator transaction binding the contract method 0xf687d12a.
//
// Solidity: function updateGasLimit(uint256 gasLimit) returns()
func (_TestMRC20 *TestMRC20TransactorSession) UpdateGasLimit(gasLimit *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.UpdateGasLimit(&_TestMRC20.TransactOpts, gasLimit)
}

// UpdateNewField is a paid mutator transaction binding the contract method 0x732bb0e4.
//
// Solidity: function updateNewField(uint256 newField_) returns()
func (_TestMRC20 *TestMRC20Transactor) UpdateNewField(opts *bind.TransactOpts, newField_ *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "updateNewField", newField_)
}

// UpdateNewField is a paid mutator transaction binding the contract method 0x732bb0e4.
//
// Solidity: function updateNewField(uint256 newField_) returns()
func (_TestMRC20 *TestMRC20Session) UpdateNewField(newField_ *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.UpdateNewField(&_TestMRC20.TransactOpts, newField_)
}

// UpdateNewField is a paid mutator transaction binding the contract method 0x732bb0e4.
//
// Solidity: function updateNewField(uint256 newField_) returns()
func (_TestMRC20 *TestMRC20TransactorSession) UpdateNewField(newField_ *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.UpdateNewField(&_TestMRC20.TransactOpts, newField_)
}

// UpdateProtocolFlatFee is a paid mutator transaction binding the contract method 0xeddeb123.
//
// Solidity: function updateProtocolFlatFee(uint256 protocolFlatFee) returns()
func (_TestMRC20 *TestMRC20Transactor) UpdateProtocolFlatFee(opts *bind.TransactOpts, protocolFlatFee *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "updateProtocolFlatFee", protocolFlatFee)
}

// UpdateProtocolFlatFee is a paid mutator transaction binding the contract method 0xeddeb123.
//
// Solidity: function updateProtocolFlatFee(uint256 protocolFlatFee) returns()
func (_TestMRC20 *TestMRC20Session) UpdateProtocolFlatFee(protocolFlatFee *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.UpdateProtocolFlatFee(&_TestMRC20.TransactOpts, protocolFlatFee)
}

// UpdateProtocolFlatFee is a paid mutator transaction binding the contract method 0xeddeb123.
//
// Solidity: function updateProtocolFlatFee(uint256 protocolFlatFee) returns()
func (_TestMRC20 *TestMRC20TransactorSession) UpdateProtocolFlatFee(protocolFlatFee *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.UpdateProtocolFlatFee(&_TestMRC20.TransactOpts, protocolFlatFee)
}

// UpdateSystemContractAddress is a paid mutator transaction binding the contract method 0xc835d7cc.
//
// Solidity: function updateSystemContractAddress(address addr) returns()
func (_TestMRC20 *TestMRC20Transactor) UpdateSystemContractAddress(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "updateSystemContractAddress", addr)
}

// UpdateSystemContractAddress is a paid mutator transaction binding the contract method 0xc835d7cc.
//
// Solidity: function updateSystemContractAddress(address addr) returns()
func (_TestMRC20 *TestMRC20Session) UpdateSystemContractAddress(addr common.Address) (*types.Transaction, error) {
	return _TestMRC20.Contract.UpdateSystemContractAddress(&_TestMRC20.TransactOpts, addr)
}

// UpdateSystemContractAddress is a paid mutator transaction binding the contract method 0xc835d7cc.
//
// Solidity: function updateSystemContractAddress(address addr) returns()
func (_TestMRC20 *TestMRC20TransactorSession) UpdateSystemContractAddress(addr common.Address) (*types.Transaction, error) {
	return _TestMRC20.Contract.UpdateSystemContractAddress(&_TestMRC20.TransactOpts, addr)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc7012626.
//
// Solidity: function withdraw(bytes to, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Transactor) Withdraw(opts *bind.TransactOpts, to []byte, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.contract.Transact(opts, "withdraw", to, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc7012626.
//
// Solidity: function withdraw(bytes to, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20Session) Withdraw(to []byte, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.Withdraw(&_TestMRC20.TransactOpts, to, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xc7012626.
//
// Solidity: function withdraw(bytes to, uint256 amount) returns(bool)
func (_TestMRC20 *TestMRC20TransactorSession) Withdraw(to []byte, amount *big.Int) (*types.Transaction, error) {
	return _TestMRC20.Contract.Withdraw(&_TestMRC20.TransactOpts, to, amount)
}

// TestMRC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TestMRC20 contract.
type TestMRC20ApprovalIterator struct {
	Event *TestMRC20Approval // Event containing the contract specifics and raw log

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
func (it *TestMRC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestMRC20Approval)
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
		it.Event = new(TestMRC20Approval)
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
func (it *TestMRC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestMRC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestMRC20Approval represents a Approval event raised by the TestMRC20 contract.
type TestMRC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestMRC20 *TestMRC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TestMRC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TestMRC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TestMRC20ApprovalIterator{contract: _TestMRC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestMRC20 *TestMRC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TestMRC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TestMRC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestMRC20Approval)
				if err := _TestMRC20.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_TestMRC20 *TestMRC20Filterer) ParseApproval(log types.Log) (*TestMRC20Approval, error) {
	event := new(TestMRC20Approval)
	if err := _TestMRC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestMRC20DepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the TestMRC20 contract.
type TestMRC20DepositIterator struct {
	Event *TestMRC20Deposit // Event containing the contract specifics and raw log

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
func (it *TestMRC20DepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestMRC20Deposit)
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
		it.Event = new(TestMRC20Deposit)
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
func (it *TestMRC20DepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestMRC20DepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestMRC20Deposit represents a Deposit event raised by the TestMRC20 contract.
type TestMRC20Deposit struct {
	From  []byte
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab3.
//
// Solidity: event Deposit(bytes from, address indexed to, uint256 value)
func (_TestMRC20 *TestMRC20Filterer) FilterDeposit(opts *bind.FilterOpts, to []common.Address) (*TestMRC20DepositIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestMRC20.contract.FilterLogs(opts, "Deposit", toRule)
	if err != nil {
		return nil, err
	}
	return &TestMRC20DepositIterator{contract: _TestMRC20.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab3.
//
// Solidity: event Deposit(bytes from, address indexed to, uint256 value)
func (_TestMRC20 *TestMRC20Filterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *TestMRC20Deposit, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestMRC20.contract.WatchLogs(opts, "Deposit", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestMRC20Deposit)
				if err := _TestMRC20.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0x67fc7bdaed5b0ec550d8706b87d60568ab70c6b781263c70101d54cd1564aab3.
//
// Solidity: event Deposit(bytes from, address indexed to, uint256 value)
func (_TestMRC20 *TestMRC20Filterer) ParseDeposit(log types.Log) (*TestMRC20Deposit, error) {
	event := new(TestMRC20Deposit)
	if err := _TestMRC20.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestMRC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TestMRC20 contract.
type TestMRC20TransferIterator struct {
	Event *TestMRC20Transfer // Event containing the contract specifics and raw log

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
func (it *TestMRC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestMRC20Transfer)
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
		it.Event = new(TestMRC20Transfer)
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
func (it *TestMRC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestMRC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestMRC20Transfer represents a Transfer event raised by the TestMRC20 contract.
type TestMRC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestMRC20 *TestMRC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestMRC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestMRC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestMRC20TransferIterator{contract: _TestMRC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestMRC20 *TestMRC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TestMRC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestMRC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestMRC20Transfer)
				if err := _TestMRC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_TestMRC20 *TestMRC20Filterer) ParseTransfer(log types.Log) (*TestMRC20Transfer, error) {
	event := new(TestMRC20Transfer)
	if err := _TestMRC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestMRC20UpdatedGasLimitIterator is returned from FilterUpdatedGasLimit and is used to iterate over the raw logs and unpacked data for UpdatedGasLimit events raised by the TestMRC20 contract.
type TestMRC20UpdatedGasLimitIterator struct {
	Event *TestMRC20UpdatedGasLimit // Event containing the contract specifics and raw log

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
func (it *TestMRC20UpdatedGasLimitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestMRC20UpdatedGasLimit)
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
		it.Event = new(TestMRC20UpdatedGasLimit)
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
func (it *TestMRC20UpdatedGasLimitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestMRC20UpdatedGasLimitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestMRC20UpdatedGasLimit represents a UpdatedGasLimit event raised by the TestMRC20 contract.
type TestMRC20UpdatedGasLimit struct {
	GasLimit *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUpdatedGasLimit is a free log retrieval operation binding the contract event 0xff5788270f43bfc1ca41c503606d2594aa3023a1a7547de403a3e2f146a4a80a.
//
// Solidity: event UpdatedGasLimit(uint256 gasLimit)
func (_TestMRC20 *TestMRC20Filterer) FilterUpdatedGasLimit(opts *bind.FilterOpts) (*TestMRC20UpdatedGasLimitIterator, error) {

	logs, sub, err := _TestMRC20.contract.FilterLogs(opts, "UpdatedGasLimit")
	if err != nil {
		return nil, err
	}
	return &TestMRC20UpdatedGasLimitIterator{contract: _TestMRC20.contract, event: "UpdatedGasLimit", logs: logs, sub: sub}, nil
}

// WatchUpdatedGasLimit is a free log subscription operation binding the contract event 0xff5788270f43bfc1ca41c503606d2594aa3023a1a7547de403a3e2f146a4a80a.
//
// Solidity: event UpdatedGasLimit(uint256 gasLimit)
func (_TestMRC20 *TestMRC20Filterer) WatchUpdatedGasLimit(opts *bind.WatchOpts, sink chan<- *TestMRC20UpdatedGasLimit) (event.Subscription, error) {

	logs, sub, err := _TestMRC20.contract.WatchLogs(opts, "UpdatedGasLimit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestMRC20UpdatedGasLimit)
				if err := _TestMRC20.contract.UnpackLog(event, "UpdatedGasLimit", log); err != nil {
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

// ParseUpdatedGasLimit is a log parse operation binding the contract event 0xff5788270f43bfc1ca41c503606d2594aa3023a1a7547de403a3e2f146a4a80a.
//
// Solidity: event UpdatedGasLimit(uint256 gasLimit)
func (_TestMRC20 *TestMRC20Filterer) ParseUpdatedGasLimit(log types.Log) (*TestMRC20UpdatedGasLimit, error) {
	event := new(TestMRC20UpdatedGasLimit)
	if err := _TestMRC20.contract.UnpackLog(event, "UpdatedGasLimit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestMRC20UpdatedProtocolFlatFeeIterator is returned from FilterUpdatedProtocolFlatFee and is used to iterate over the raw logs and unpacked data for UpdatedProtocolFlatFee events raised by the TestMRC20 contract.
type TestMRC20UpdatedProtocolFlatFeeIterator struct {
	Event *TestMRC20UpdatedProtocolFlatFee // Event containing the contract specifics and raw log

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
func (it *TestMRC20UpdatedProtocolFlatFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestMRC20UpdatedProtocolFlatFee)
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
		it.Event = new(TestMRC20UpdatedProtocolFlatFee)
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
func (it *TestMRC20UpdatedProtocolFlatFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestMRC20UpdatedProtocolFlatFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestMRC20UpdatedProtocolFlatFee represents a UpdatedProtocolFlatFee event raised by the TestMRC20 contract.
type TestMRC20UpdatedProtocolFlatFee struct {
	ProtocolFlatFee *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpdatedProtocolFlatFee is a free log retrieval operation binding the contract event 0xef13af88e424b5d15f49c77758542c1938b08b8b95b91ed0751f98ba99000d8f.
//
// Solidity: event UpdatedProtocolFlatFee(uint256 protocolFlatFee)
func (_TestMRC20 *TestMRC20Filterer) FilterUpdatedProtocolFlatFee(opts *bind.FilterOpts) (*TestMRC20UpdatedProtocolFlatFeeIterator, error) {

	logs, sub, err := _TestMRC20.contract.FilterLogs(opts, "UpdatedProtocolFlatFee")
	if err != nil {
		return nil, err
	}
	return &TestMRC20UpdatedProtocolFlatFeeIterator{contract: _TestMRC20.contract, event: "UpdatedProtocolFlatFee", logs: logs, sub: sub}, nil
}

// WatchUpdatedProtocolFlatFee is a free log subscription operation binding the contract event 0xef13af88e424b5d15f49c77758542c1938b08b8b95b91ed0751f98ba99000d8f.
//
// Solidity: event UpdatedProtocolFlatFee(uint256 protocolFlatFee)
func (_TestMRC20 *TestMRC20Filterer) WatchUpdatedProtocolFlatFee(opts *bind.WatchOpts, sink chan<- *TestMRC20UpdatedProtocolFlatFee) (event.Subscription, error) {

	logs, sub, err := _TestMRC20.contract.WatchLogs(opts, "UpdatedProtocolFlatFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestMRC20UpdatedProtocolFlatFee)
				if err := _TestMRC20.contract.UnpackLog(event, "UpdatedProtocolFlatFee", log); err != nil {
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

// ParseUpdatedProtocolFlatFee is a log parse operation binding the contract event 0xef13af88e424b5d15f49c77758542c1938b08b8b95b91ed0751f98ba99000d8f.
//
// Solidity: event UpdatedProtocolFlatFee(uint256 protocolFlatFee)
func (_TestMRC20 *TestMRC20Filterer) ParseUpdatedProtocolFlatFee(log types.Log) (*TestMRC20UpdatedProtocolFlatFee, error) {
	event := new(TestMRC20UpdatedProtocolFlatFee)
	if err := _TestMRC20.contract.UnpackLog(event, "UpdatedProtocolFlatFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestMRC20UpdatedSystemContractIterator is returned from FilterUpdatedSystemContract and is used to iterate over the raw logs and unpacked data for UpdatedSystemContract events raised by the TestMRC20 contract.
type TestMRC20UpdatedSystemContractIterator struct {
	Event *TestMRC20UpdatedSystemContract // Event containing the contract specifics and raw log

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
func (it *TestMRC20UpdatedSystemContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestMRC20UpdatedSystemContract)
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
		it.Event = new(TestMRC20UpdatedSystemContract)
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
func (it *TestMRC20UpdatedSystemContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestMRC20UpdatedSystemContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestMRC20UpdatedSystemContract represents a UpdatedSystemContract event raised by the TestMRC20 contract.
type TestMRC20UpdatedSystemContract struct {
	SystemContract common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpdatedSystemContract is a free log retrieval operation binding the contract event 0xd55614e962c5fd6ece71614f6348d702468a997a394dd5e5c1677950226d97ae.
//
// Solidity: event UpdatedSystemContract(address systemContract)
func (_TestMRC20 *TestMRC20Filterer) FilterUpdatedSystemContract(opts *bind.FilterOpts) (*TestMRC20UpdatedSystemContractIterator, error) {

	logs, sub, err := _TestMRC20.contract.FilterLogs(opts, "UpdatedSystemContract")
	if err != nil {
		return nil, err
	}
	return &TestMRC20UpdatedSystemContractIterator{contract: _TestMRC20.contract, event: "UpdatedSystemContract", logs: logs, sub: sub}, nil
}

// WatchUpdatedSystemContract is a free log subscription operation binding the contract event 0xd55614e962c5fd6ece71614f6348d702468a997a394dd5e5c1677950226d97ae.
//
// Solidity: event UpdatedSystemContract(address systemContract)
func (_TestMRC20 *TestMRC20Filterer) WatchUpdatedSystemContract(opts *bind.WatchOpts, sink chan<- *TestMRC20UpdatedSystemContract) (event.Subscription, error) {

	logs, sub, err := _TestMRC20.contract.WatchLogs(opts, "UpdatedSystemContract")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestMRC20UpdatedSystemContract)
				if err := _TestMRC20.contract.UnpackLog(event, "UpdatedSystemContract", log); err != nil {
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

// ParseUpdatedSystemContract is a log parse operation binding the contract event 0xd55614e962c5fd6ece71614f6348d702468a997a394dd5e5c1677950226d97ae.
//
// Solidity: event UpdatedSystemContract(address systemContract)
func (_TestMRC20 *TestMRC20Filterer) ParseUpdatedSystemContract(log types.Log) (*TestMRC20UpdatedSystemContract, error) {
	event := new(TestMRC20UpdatedSystemContract)
	if err := _TestMRC20.contract.UnpackLog(event, "UpdatedSystemContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestMRC20WithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the TestMRC20 contract.
type TestMRC20WithdrawalIterator struct {
	Event *TestMRC20Withdrawal // Event containing the contract specifics and raw log

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
func (it *TestMRC20WithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestMRC20Withdrawal)
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
		it.Event = new(TestMRC20Withdrawal)
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
func (it *TestMRC20WithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestMRC20WithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestMRC20Withdrawal represents a Withdrawal event raised by the TestMRC20 contract.
type TestMRC20Withdrawal struct {
	From            common.Address
	To              []byte
	Value           *big.Int
	Gasfee          *big.Int
	ProtocolFlatFee *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955.
//
// Solidity: event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee, uint256 protocolFlatFee)
func (_TestMRC20 *TestMRC20Filterer) FilterWithdrawal(opts *bind.FilterOpts, from []common.Address) (*TestMRC20WithdrawalIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _TestMRC20.contract.FilterLogs(opts, "Withdrawal", fromRule)
	if err != nil {
		return nil, err
	}
	return &TestMRC20WithdrawalIterator{contract: _TestMRC20.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955.
//
// Solidity: event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee, uint256 protocolFlatFee)
func (_TestMRC20 *TestMRC20Filterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *TestMRC20Withdrawal, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _TestMRC20.contract.WatchLogs(opts, "Withdrawal", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestMRC20Withdrawal)
				if err := _TestMRC20.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955.
//
// Solidity: event Withdrawal(address indexed from, bytes to, uint256 value, uint256 gasfee, uint256 protocolFlatFee)
func (_TestMRC20 *TestMRC20Filterer) ParseWithdrawal(log types.Log) (*TestMRC20Withdrawal, error) {
	event := new(TestMRC20Withdrawal)
	if err := _TestMRC20.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
