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

// MuseConnectorEthMetaData contains all meta data concerning the MuseConnectorEth contract.
var MuseConnectorEthMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"museToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tssAddress_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tssAddressUpdater_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pauserAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getLockedAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onReceive\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"museValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onRevert\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingMuseValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceTssAddressUpdater\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"send\",\"inputs\":[{\"name\":\"input\",\"type\":\"tuple\",\"internalType\":\"structMuseInterfaces.SendInput\",\"components\":[{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destinationGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"museValueAndGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"museParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tssAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tssAddressUpdater\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauserAddress\",\"inputs\":[{\"name\":\"pauserAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTssAddress\",\"inputs\":[{\"name\":\"tssAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"museToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserAddressUpdated\",\"inputs\":[{\"name\":\"callerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newTssAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TSSAddressUpdated\",\"inputs\":[{\"name\":\"callerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newTssAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TSSAddressUpdaterUpdated\",\"inputs\":[{\"name\":\"callerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newTssUpdaterAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MuseReceived\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"museValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MuseReverted\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"remainingMuseValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MuseSent\",\"inputs\":[{\"name\":\"sourceTxOriginAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"museTxSenderAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"museValueAndGas\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationGasLimit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"museParams\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CallerIsNotPauser\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CallerIsNotTss\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CallerIsNotTssOrUpdater\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CallerIsNotTssUpdater\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExceedsMaxSupply\",\"inputs\":[{\"name\":\"maxSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MuseTransferError\",\"inputs\":[]}]",
	Bin: "0x60a060405234801561001057600080fd5b5060405161175238038061175283398101604081905261002f91610114565b6000805460ff19169055838383836001600160a01b038416158061005a57506001600160a01b038316155b8061006c57506001600160a01b038216155b8061007e57506001600160a01b038116155b1561009c5760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b03938416608052600180549385166001600160a01b0319948516179055600280549285169290931691909117909155600080549190921661010002610100600160a81b03199091161790555061016892505050565b80516001600160a01b038116811461010f57600080fd5b919050565b6000806000806080858703121561012a57600080fd5b610133856100f8565b9350610141602086016100f8565b925061014f604086016100f8565b915061015d606086016100f8565b905092959194509250565b6080516115b461019e6000396000818160f40152818161027b015281816103a701528181610ae90152610d6701526115b46000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80636128480f1161008c5780639122c344116100665780639122c344146101ec578063942a5e16146101ff578063ec02690114610212578063f7fb869b1461022557600080fd5b80636128480f146101c9578063779e3b63146101dc5780638456cb59146101e457600080fd5b8063328a01d0116100c8578063328a01d01461016b5780633f4ba83a1461018b5780635b112591146101935780635c975abb146101b357600080fd5b806321e093b1146100ef578063252bc8861461014057806329dd214d14610156575b600080fd5b6101167f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b61014861024a565b604051908152602001610137565b610169610164366004611057565b610300565b005b6002546101169073ffffffffffffffffffffffffffffffffffffffff1681565b6101696105eb565b6001546101169073ffffffffffffffffffffffffffffffffffffffff1681565b60005460ff166040519015158152602001610137565b6101696101d73660046110f7565b61064d565b61016961077a565b61016961089f565b6101696101fa3660046110f7565b6108ff565b61016961020d366004611119565b610a3f565b6101696102203660046111c3565b610d1e565b60005461011690610100900473ffffffffffffffffffffffffffffffffffffffff1681565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156102d7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102fb91906111fe565b905090565b60015473ffffffffffffffffffffffffffffffffffffffff163314610358576040517fff70ace20000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8681166004830152602482018690526000917f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af11580156103f2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104169190611217565b90508061044f576040517f20878f6200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8215610588578573ffffffffffffffffffffffffffffffffffffffff16633749c51a6040518060a001604052808c8c8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060208082018c905273ffffffffffffffffffffffffffffffffffffffff8b16604080840191909152606083018b90528051601f8a0183900483028101830190915288815260809092019190899089908190840183828082843760009201919091525050509152506040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b168152610555919060040161129d565b600060405180830381600087803b15801561056f57600080fd5b505af1158015610583573d6000803e3d6000fd5b505050505b818673ffffffffffffffffffffffffffffffffffffffff16887ff1302855733b40d8acb467ee990b6d56c05c80e28ebcabfa6e6f3f57cb50d6988c8c8a8a8a6040516105d895949392919061137a565b60405180910390a4505050505050505050565b600054610100900473ffffffffffffffffffffffffffffffffffffffff163314610643576040517f4677a0d300000000000000000000000000000000000000000000000000000000815233600482015260240161034f565b61064b610e99565b565b600054610100900473ffffffffffffffffffffffffffffffffffffffff1633146106a5576040517f4677a0d300000000000000000000000000000000000000000000000000000000815233600482015260240161034f565b73ffffffffffffffffffffffffffffffffffffffff81166106f2576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffff0000000000000000000000000000000000000000ff1661010073ffffffffffffffffffffffffffffffffffffffff8416908102919091179091556040805133815260208101929092527fd41d83655d484bdf299598751c371b2d92088667266fe3774b25a97bdd5d039791015b60405180910390a150565b60025473ffffffffffffffffffffffffffffffffffffffff1633146107cd576040517fe700765e00000000000000000000000000000000000000000000000000000000815233600482015260240161034f565b60015473ffffffffffffffffffffffffffffffffffffffff1661081c576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600154600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691821790556040805133815260208101929092527f5104c9abdc7d111c2aeb4ce890ac70274a4be2ee83f46a62551be5e6ebc82dd091015b60405180910390a1565b600054610100900473ffffffffffffffffffffffffffffffffffffffff1633146108f7576040517f4677a0d300000000000000000000000000000000000000000000000000000000815233600482015260240161034f565b61064b610f11565b60015473ffffffffffffffffffffffffffffffffffffffff16331480159061093f575060025473ffffffffffffffffffffffffffffffffffffffff163314155b15610978576040517fcdfcef9700000000000000000000000000000000000000000000000000000000815233600482015260240161034f565b73ffffffffffffffffffffffffffffffffffffffff81166109c5576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040805133815260208101929092527fe79965b5c67dcfb2cf5fe152715e4a7256cee62a3d5dd8484fd8a8539eb8beff910161076f565b610a47610f6c565b60015473ffffffffffffffffffffffffffffffffffffffff163314610a9a576040517fff70ace200000000000000000000000000000000000000000000000000000000815233600482015260240161034f565b6040517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8a81166004830152602482018690526000917f00000000000000000000000000000000000000000000000000000000000000009091169063a9059cbb906044016020604051808303816000875af1158015610b34573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b589190611217565b905080610b91576040517f20878f6200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8215610ccd578973ffffffffffffffffffffffffffffffffffffffff16633ff0693c6040518060c001604052808d73ffffffffffffffffffffffffffffffffffffffff1681526020018c81526020018b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060208082018b905260408083018b90528051601f8a0183900483028101830190915288815260609092019190899089908190840183828082843760009201919091525050509152506040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b168152610c9a91906004016113b3565b600060405180830381600087803b158015610cb457600080fd5b505af1158015610cc8573d6000803e3d6000fd5b505050505b81867f521fb0b407c2eb9b1375530e9b9a569889992140a688bc076aa72c1712012c888c8c8c8c8b8b8b604051610d0a9796959493929190611448565b60405180910390a350505050505050505050565b610d26610f6c565b6040517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152608082013560448201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906323b872dd906064016020604051808303816000875af1158015610dc5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610de99190611217565b905080610e22576040517f20878f6200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8135337f7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e432610e5460208701876114a5565b60808801356040890135610e6b60608b018b6114a5565b610e7860a08d018d6114a5565b604051610e8d9998979695949392919061150a565b60405180910390a35050565b610ea1610fa9565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610895565b610f19610f6c565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258610eec3390565b60005460ff161561064b576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005460ff1661064b576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008083601f840112610ff757600080fd5b50813567ffffffffffffffff81111561100f57600080fd5b60208301915083602082850101111561102757600080fd5b9250929050565b803573ffffffffffffffffffffffffffffffffffffffff8116811461105257600080fd5b919050565b60008060008060008060008060c0898b03121561107357600080fd5b883567ffffffffffffffff81111561108a57600080fd5b6110968b828c01610fe5565b909950975050602089013595506110af60408a0161102e565b945060608901359350608089013567ffffffffffffffff8111156110d257600080fd5b6110de8b828c01610fe5565b999c989b50969995989497949560a00135949350505050565b60006020828403121561110957600080fd5b6111128261102e565b9392505050565b600080600080600080600080600060e08a8c03121561113757600080fd5b6111408a61102e565b985060208a0135975060408a013567ffffffffffffffff81111561116357600080fd5b61116f8c828d01610fe5565b90985096505060608a0135945060808a0135935060a08a013567ffffffffffffffff81111561119d57600080fd5b6111a98c828d01610fe5565b9a9d999c50979a9699959894979660c00135949350505050565b6000602082840312156111d557600080fd5b813567ffffffffffffffff8111156111ec57600080fd5b820160c0818503121561111257600080fd5b60006020828403121561121057600080fd5b5051919050565b60006020828403121561122957600080fd5b8151801515811461111257600080fd5b6000815180845260005b8181101561125f57602081850181015186830182015201611243565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000825160a060208401526112b960c0840182611239565b90506020840151604084015273ffffffffffffffffffffffffffffffffffffffff60408501511660608401526060840151608084015260808401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160a08501526113288282611239565b95945050505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b60608152600061138e606083018789611331565b85602084015282810360408401526113a7818587611331565b98975050505050505050565b6020815273ffffffffffffffffffffffffffffffffffffffff8251166020820152602082015160408201526000604083015160c060608401526113f960e0840182611239565b905060608401516080840152608084015160a084015260a08401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160c08501526113288282611239565b73ffffffffffffffffffffffffffffffffffffffff8816815286602082015260a06040820152600061147e60a083018789611331565b8560608401528281036080840152611497818587611331565b9a9950505050505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126114da57600080fd5b83018035915067ffffffffffffffff8211156114f557600080fd5b60200191503681900382131561102757600080fd5b73ffffffffffffffffffffffffffffffffffffffff8a16815260c06020820152600061153a60c083018a8c611331565b8860408401528760608401528281036080840152611559818789611331565b905082810360a084015261156e818587611331565b9c9b50505050505050505050505056fea26469706673582212202898e62086c182d504ccd73f12057f5453b7103ddf6ba6c1ed17b3978fdd0d4f64736f6c634300081a0033",
}

// MuseConnectorEthABI is the input ABI used to generate the binding from.
// Deprecated: Use MuseConnectorEthMetaData.ABI instead.
var MuseConnectorEthABI = MuseConnectorEthMetaData.ABI

// MuseConnectorEthBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MuseConnectorEthMetaData.Bin instead.
var MuseConnectorEthBin = MuseConnectorEthMetaData.Bin

// DeployMuseConnectorEth deploys a new Ethereum contract, binding an instance of MuseConnectorEth to it.
func DeployMuseConnectorEth(auth *bind.TransactOpts, backend bind.ContractBackend, museToken_ common.Address, tssAddress_ common.Address, tssAddressUpdater_ common.Address, pauserAddress_ common.Address) (common.Address, *types.Transaction, *MuseConnectorEth, error) {
	parsed, err := MuseConnectorEthMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MuseConnectorEthBin), backend, museToken_, tssAddress_, tssAddressUpdater_, pauserAddress_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MuseConnectorEth{MuseConnectorEthCaller: MuseConnectorEthCaller{contract: contract}, MuseConnectorEthTransactor: MuseConnectorEthTransactor{contract: contract}, MuseConnectorEthFilterer: MuseConnectorEthFilterer{contract: contract}}, nil
}

// MuseConnectorEth is an auto generated Go binding around an Ethereum contract.
type MuseConnectorEth struct {
	MuseConnectorEthCaller     // Read-only binding to the contract
	MuseConnectorEthTransactor // Write-only binding to the contract
	MuseConnectorEthFilterer   // Log filterer for contract events
}

// MuseConnectorEthCaller is an auto generated read-only Go binding around an Ethereum contract.
type MuseConnectorEthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseConnectorEthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MuseConnectorEthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseConnectorEthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MuseConnectorEthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseConnectorEthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MuseConnectorEthSession struct {
	Contract     *MuseConnectorEth // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MuseConnectorEthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MuseConnectorEthCallerSession struct {
	Contract *MuseConnectorEthCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// MuseConnectorEthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MuseConnectorEthTransactorSession struct {
	Contract     *MuseConnectorEthTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// MuseConnectorEthRaw is an auto generated low-level Go binding around an Ethereum contract.
type MuseConnectorEthRaw struct {
	Contract *MuseConnectorEth // Generic contract binding to access the raw methods on
}

// MuseConnectorEthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MuseConnectorEthCallerRaw struct {
	Contract *MuseConnectorEthCaller // Generic read-only contract binding to access the raw methods on
}

// MuseConnectorEthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MuseConnectorEthTransactorRaw struct {
	Contract *MuseConnectorEthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMuseConnectorEth creates a new instance of MuseConnectorEth, bound to a specific deployed contract.
func NewMuseConnectorEth(address common.Address, backend bind.ContractBackend) (*MuseConnectorEth, error) {
	contract, err := bindMuseConnectorEth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEth{MuseConnectorEthCaller: MuseConnectorEthCaller{contract: contract}, MuseConnectorEthTransactor: MuseConnectorEthTransactor{contract: contract}, MuseConnectorEthFilterer: MuseConnectorEthFilterer{contract: contract}}, nil
}

// NewMuseConnectorEthCaller creates a new read-only instance of MuseConnectorEth, bound to a specific deployed contract.
func NewMuseConnectorEthCaller(address common.Address, caller bind.ContractCaller) (*MuseConnectorEthCaller, error) {
	contract, err := bindMuseConnectorEth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthCaller{contract: contract}, nil
}

// NewMuseConnectorEthTransactor creates a new write-only instance of MuseConnectorEth, bound to a specific deployed contract.
func NewMuseConnectorEthTransactor(address common.Address, transactor bind.ContractTransactor) (*MuseConnectorEthTransactor, error) {
	contract, err := bindMuseConnectorEth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthTransactor{contract: contract}, nil
}

// NewMuseConnectorEthFilterer creates a new log filterer instance of MuseConnectorEth, bound to a specific deployed contract.
func NewMuseConnectorEthFilterer(address common.Address, filterer bind.ContractFilterer) (*MuseConnectorEthFilterer, error) {
	contract, err := bindMuseConnectorEth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthFilterer{contract: contract}, nil
}

// bindMuseConnectorEth binds a generic wrapper to an already deployed contract.
func bindMuseConnectorEth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MuseConnectorEthMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MuseConnectorEth *MuseConnectorEthRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MuseConnectorEth.Contract.MuseConnectorEthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MuseConnectorEth *MuseConnectorEthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.MuseConnectorEthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MuseConnectorEth *MuseConnectorEthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.MuseConnectorEthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MuseConnectorEth *MuseConnectorEthCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MuseConnectorEth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MuseConnectorEth *MuseConnectorEthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MuseConnectorEth *MuseConnectorEthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.contract.Transact(opts, method, params...)
}

// GetLockedAmount is a free data retrieval call binding the contract method 0x252bc886.
//
// Solidity: function getLockedAmount() view returns(uint256)
func (_MuseConnectorEth *MuseConnectorEthCaller) GetLockedAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MuseConnectorEth.contract.Call(opts, &out, "getLockedAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockedAmount is a free data retrieval call binding the contract method 0x252bc886.
//
// Solidity: function getLockedAmount() view returns(uint256)
func (_MuseConnectorEth *MuseConnectorEthSession) GetLockedAmount() (*big.Int, error) {
	return _MuseConnectorEth.Contract.GetLockedAmount(&_MuseConnectorEth.CallOpts)
}

// GetLockedAmount is a free data retrieval call binding the contract method 0x252bc886.
//
// Solidity: function getLockedAmount() view returns(uint256)
func (_MuseConnectorEth *MuseConnectorEthCallerSession) GetLockedAmount() (*big.Int, error) {
	return _MuseConnectorEth.Contract.GetLockedAmount(&_MuseConnectorEth.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MuseConnectorEth *MuseConnectorEthCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MuseConnectorEth.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MuseConnectorEth *MuseConnectorEthSession) Paused() (bool, error) {
	return _MuseConnectorEth.Contract.Paused(&_MuseConnectorEth.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MuseConnectorEth *MuseConnectorEthCallerSession) Paused() (bool, error) {
	return _MuseConnectorEth.Contract.Paused(&_MuseConnectorEth.CallOpts)
}

// PauserAddress is a free data retrieval call binding the contract method 0xf7fb869b.
//
// Solidity: function pauserAddress() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthCaller) PauserAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorEth.contract.Call(opts, &out, "pauserAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PauserAddress is a free data retrieval call binding the contract method 0xf7fb869b.
//
// Solidity: function pauserAddress() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthSession) PauserAddress() (common.Address, error) {
	return _MuseConnectorEth.Contract.PauserAddress(&_MuseConnectorEth.CallOpts)
}

// PauserAddress is a free data retrieval call binding the contract method 0xf7fb869b.
//
// Solidity: function pauserAddress() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthCallerSession) PauserAddress() (common.Address, error) {
	return _MuseConnectorEth.Contract.PauserAddress(&_MuseConnectorEth.CallOpts)
}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthCaller) TssAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorEth.contract.Call(opts, &out, "tssAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthSession) TssAddress() (common.Address, error) {
	return _MuseConnectorEth.Contract.TssAddress(&_MuseConnectorEth.CallOpts)
}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthCallerSession) TssAddress() (common.Address, error) {
	return _MuseConnectorEth.Contract.TssAddress(&_MuseConnectorEth.CallOpts)
}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthCaller) TssAddressUpdater(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorEth.contract.Call(opts, &out, "tssAddressUpdater")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthSession) TssAddressUpdater() (common.Address, error) {
	return _MuseConnectorEth.Contract.TssAddressUpdater(&_MuseConnectorEth.CallOpts)
}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthCallerSession) TssAddressUpdater() (common.Address, error) {
	return _MuseConnectorEth.Contract.TssAddressUpdater(&_MuseConnectorEth.CallOpts)
}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthCaller) MuseToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorEth.contract.Call(opts, &out, "museToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthSession) MuseToken() (common.Address, error) {
	return _MuseConnectorEth.Contract.MuseToken(&_MuseConnectorEth.CallOpts)
}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_MuseConnectorEth *MuseConnectorEthCallerSession) MuseToken() (common.Address, error) {
	return _MuseConnectorEth.Contract.MuseToken(&_MuseConnectorEth.CallOpts)
}

// OnReceive is a paid mutator transaction binding the contract method 0x29dd214d.
//
// Solidity: function onReceive(bytes museTxSenderAddress, uint256 sourceChainId, address destinationAddress, uint256 museValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorEth *MuseConnectorEthTransactor) OnReceive(opts *bind.TransactOpts, museTxSenderAddress []byte, sourceChainId *big.Int, destinationAddress common.Address, museValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorEth.contract.Transact(opts, "onReceive", museTxSenderAddress, sourceChainId, destinationAddress, museValue, message, internalSendHash)
}

// OnReceive is a paid mutator transaction binding the contract method 0x29dd214d.
//
// Solidity: function onReceive(bytes museTxSenderAddress, uint256 sourceChainId, address destinationAddress, uint256 museValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorEth *MuseConnectorEthSession) OnReceive(museTxSenderAddress []byte, sourceChainId *big.Int, destinationAddress common.Address, museValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.OnReceive(&_MuseConnectorEth.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, museValue, message, internalSendHash)
}

// OnReceive is a paid mutator transaction binding the contract method 0x29dd214d.
//
// Solidity: function onReceive(bytes museTxSenderAddress, uint256 sourceChainId, address destinationAddress, uint256 museValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorEth *MuseConnectorEthTransactorSession) OnReceive(museTxSenderAddress []byte, sourceChainId *big.Int, destinationAddress common.Address, museValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.OnReceive(&_MuseConnectorEth.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, museValue, message, internalSendHash)
}

// OnRevert is a paid mutator transaction binding the contract method 0x942a5e16.
//
// Solidity: function onRevert(address museTxSenderAddress, uint256 sourceChainId, bytes destinationAddress, uint256 destinationChainId, uint256 remainingMuseValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorEth *MuseConnectorEthTransactor) OnRevert(opts *bind.TransactOpts, museTxSenderAddress common.Address, sourceChainId *big.Int, destinationAddress []byte, destinationChainId *big.Int, remainingMuseValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorEth.contract.Transact(opts, "onRevert", museTxSenderAddress, sourceChainId, destinationAddress, destinationChainId, remainingMuseValue, message, internalSendHash)
}

// OnRevert is a paid mutator transaction binding the contract method 0x942a5e16.
//
// Solidity: function onRevert(address museTxSenderAddress, uint256 sourceChainId, bytes destinationAddress, uint256 destinationChainId, uint256 remainingMuseValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorEth *MuseConnectorEthSession) OnRevert(museTxSenderAddress common.Address, sourceChainId *big.Int, destinationAddress []byte, destinationChainId *big.Int, remainingMuseValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.OnRevert(&_MuseConnectorEth.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, destinationChainId, remainingMuseValue, message, internalSendHash)
}

// OnRevert is a paid mutator transaction binding the contract method 0x942a5e16.
//
// Solidity: function onRevert(address museTxSenderAddress, uint256 sourceChainId, bytes destinationAddress, uint256 destinationChainId, uint256 remainingMuseValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorEth *MuseConnectorEthTransactorSession) OnRevert(museTxSenderAddress common.Address, sourceChainId *big.Int, destinationAddress []byte, destinationChainId *big.Int, remainingMuseValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.OnRevert(&_MuseConnectorEth.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, destinationChainId, remainingMuseValue, message, internalSendHash)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MuseConnectorEth *MuseConnectorEthTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorEth.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MuseConnectorEth *MuseConnectorEthSession) Pause() (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.Pause(&_MuseConnectorEth.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MuseConnectorEth *MuseConnectorEthTransactorSession) Pause() (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.Pause(&_MuseConnectorEth.TransactOpts)
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_MuseConnectorEth *MuseConnectorEthTransactor) RenounceTssAddressUpdater(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorEth.contract.Transact(opts, "renounceTssAddressUpdater")
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_MuseConnectorEth *MuseConnectorEthSession) RenounceTssAddressUpdater() (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.RenounceTssAddressUpdater(&_MuseConnectorEth.TransactOpts)
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_MuseConnectorEth *MuseConnectorEthTransactorSession) RenounceTssAddressUpdater() (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.RenounceTssAddressUpdater(&_MuseConnectorEth.TransactOpts)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_MuseConnectorEth *MuseConnectorEthTransactor) Send(opts *bind.TransactOpts, input MuseInterfacesSendInput) (*types.Transaction, error) {
	return _MuseConnectorEth.contract.Transact(opts, "send", input)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_MuseConnectorEth *MuseConnectorEthSession) Send(input MuseInterfacesSendInput) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.Send(&_MuseConnectorEth.TransactOpts, input)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_MuseConnectorEth *MuseConnectorEthTransactorSession) Send(input MuseInterfacesSendInput) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.Send(&_MuseConnectorEth.TransactOpts, input)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MuseConnectorEth *MuseConnectorEthTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorEth.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MuseConnectorEth *MuseConnectorEthSession) Unpause() (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.Unpause(&_MuseConnectorEth.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MuseConnectorEth *MuseConnectorEthTransactorSession) Unpause() (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.Unpause(&_MuseConnectorEth.TransactOpts)
}

// UpdatePauserAddress is a paid mutator transaction binding the contract method 0x6128480f.
//
// Solidity: function updatePauserAddress(address pauserAddress_) returns()
func (_MuseConnectorEth *MuseConnectorEthTransactor) UpdatePauserAddress(opts *bind.TransactOpts, pauserAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorEth.contract.Transact(opts, "updatePauserAddress", pauserAddress_)
}

// UpdatePauserAddress is a paid mutator transaction binding the contract method 0x6128480f.
//
// Solidity: function updatePauserAddress(address pauserAddress_) returns()
func (_MuseConnectorEth *MuseConnectorEthSession) UpdatePauserAddress(pauserAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.UpdatePauserAddress(&_MuseConnectorEth.TransactOpts, pauserAddress_)
}

// UpdatePauserAddress is a paid mutator transaction binding the contract method 0x6128480f.
//
// Solidity: function updatePauserAddress(address pauserAddress_) returns()
func (_MuseConnectorEth *MuseConnectorEthTransactorSession) UpdatePauserAddress(pauserAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.UpdatePauserAddress(&_MuseConnectorEth.TransactOpts, pauserAddress_)
}

// UpdateTssAddress is a paid mutator transaction binding the contract method 0x9122c344.
//
// Solidity: function updateTssAddress(address tssAddress_) returns()
func (_MuseConnectorEth *MuseConnectorEthTransactor) UpdateTssAddress(opts *bind.TransactOpts, tssAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorEth.contract.Transact(opts, "updateTssAddress", tssAddress_)
}

// UpdateTssAddress is a paid mutator transaction binding the contract method 0x9122c344.
//
// Solidity: function updateTssAddress(address tssAddress_) returns()
func (_MuseConnectorEth *MuseConnectorEthSession) UpdateTssAddress(tssAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.UpdateTssAddress(&_MuseConnectorEth.TransactOpts, tssAddress_)
}

// UpdateTssAddress is a paid mutator transaction binding the contract method 0x9122c344.
//
// Solidity: function updateTssAddress(address tssAddress_) returns()
func (_MuseConnectorEth *MuseConnectorEthTransactorSession) UpdateTssAddress(tssAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorEth.Contract.UpdateTssAddress(&_MuseConnectorEth.TransactOpts, tssAddress_)
}

// MuseConnectorEthPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the MuseConnectorEth contract.
type MuseConnectorEthPausedIterator struct {
	Event *MuseConnectorEthPaused // Event containing the contract specifics and raw log

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
func (it *MuseConnectorEthPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorEthPaused)
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
		it.Event = new(MuseConnectorEthPaused)
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
func (it *MuseConnectorEthPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorEthPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorEthPaused represents a Paused event raised by the MuseConnectorEth contract.
type MuseConnectorEthPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MuseConnectorEth *MuseConnectorEthFilterer) FilterPaused(opts *bind.FilterOpts) (*MuseConnectorEthPausedIterator, error) {

	logs, sub, err := _MuseConnectorEth.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthPausedIterator{contract: _MuseConnectorEth.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MuseConnectorEth *MuseConnectorEthFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *MuseConnectorEthPaused) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorEth.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorEthPaused)
				if err := _MuseConnectorEth.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) ParsePaused(log types.Log) (*MuseConnectorEthPaused, error) {
	event := new(MuseConnectorEthPaused)
	if err := _MuseConnectorEth.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorEthPauserAddressUpdatedIterator is returned from FilterPauserAddressUpdated and is used to iterate over the raw logs and unpacked data for PauserAddressUpdated events raised by the MuseConnectorEth contract.
type MuseConnectorEthPauserAddressUpdatedIterator struct {
	Event *MuseConnectorEthPauserAddressUpdated // Event containing the contract specifics and raw log

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
func (it *MuseConnectorEthPauserAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorEthPauserAddressUpdated)
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
		it.Event = new(MuseConnectorEthPauserAddressUpdated)
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
func (it *MuseConnectorEthPauserAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorEthPauserAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorEthPauserAddressUpdated represents a PauserAddressUpdated event raised by the MuseConnectorEth contract.
type MuseConnectorEthPauserAddressUpdated struct {
	CallerAddress common.Address
	NewTssAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterPauserAddressUpdated is a free log retrieval operation binding the contract event 0xd41d83655d484bdf299598751c371b2d92088667266fe3774b25a97bdd5d0397.
//
// Solidity: event PauserAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorEth *MuseConnectorEthFilterer) FilterPauserAddressUpdated(opts *bind.FilterOpts) (*MuseConnectorEthPauserAddressUpdatedIterator, error) {

	logs, sub, err := _MuseConnectorEth.contract.FilterLogs(opts, "PauserAddressUpdated")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthPauserAddressUpdatedIterator{contract: _MuseConnectorEth.contract, event: "PauserAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchPauserAddressUpdated is a free log subscription operation binding the contract event 0xd41d83655d484bdf299598751c371b2d92088667266fe3774b25a97bdd5d0397.
//
// Solidity: event PauserAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorEth *MuseConnectorEthFilterer) WatchPauserAddressUpdated(opts *bind.WatchOpts, sink chan<- *MuseConnectorEthPauserAddressUpdated) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorEth.contract.WatchLogs(opts, "PauserAddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorEthPauserAddressUpdated)
				if err := _MuseConnectorEth.contract.UnpackLog(event, "PauserAddressUpdated", log); err != nil {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) ParsePauserAddressUpdated(log types.Log) (*MuseConnectorEthPauserAddressUpdated, error) {
	event := new(MuseConnectorEthPauserAddressUpdated)
	if err := _MuseConnectorEth.contract.UnpackLog(event, "PauserAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorEthTSSAddressUpdatedIterator is returned from FilterTSSAddressUpdated and is used to iterate over the raw logs and unpacked data for TSSAddressUpdated events raised by the MuseConnectorEth contract.
type MuseConnectorEthTSSAddressUpdatedIterator struct {
	Event *MuseConnectorEthTSSAddressUpdated // Event containing the contract specifics and raw log

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
func (it *MuseConnectorEthTSSAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorEthTSSAddressUpdated)
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
		it.Event = new(MuseConnectorEthTSSAddressUpdated)
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
func (it *MuseConnectorEthTSSAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorEthTSSAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorEthTSSAddressUpdated represents a TSSAddressUpdated event raised by the MuseConnectorEth contract.
type MuseConnectorEthTSSAddressUpdated struct {
	CallerAddress common.Address
	NewTssAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTSSAddressUpdated is a free log retrieval operation binding the contract event 0xe79965b5c67dcfb2cf5fe152715e4a7256cee62a3d5dd8484fd8a8539eb8beff.
//
// Solidity: event TSSAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorEth *MuseConnectorEthFilterer) FilterTSSAddressUpdated(opts *bind.FilterOpts) (*MuseConnectorEthTSSAddressUpdatedIterator, error) {

	logs, sub, err := _MuseConnectorEth.contract.FilterLogs(opts, "TSSAddressUpdated")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthTSSAddressUpdatedIterator{contract: _MuseConnectorEth.contract, event: "TSSAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchTSSAddressUpdated is a free log subscription operation binding the contract event 0xe79965b5c67dcfb2cf5fe152715e4a7256cee62a3d5dd8484fd8a8539eb8beff.
//
// Solidity: event TSSAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorEth *MuseConnectorEthFilterer) WatchTSSAddressUpdated(opts *bind.WatchOpts, sink chan<- *MuseConnectorEthTSSAddressUpdated) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorEth.contract.WatchLogs(opts, "TSSAddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorEthTSSAddressUpdated)
				if err := _MuseConnectorEth.contract.UnpackLog(event, "TSSAddressUpdated", log); err != nil {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) ParseTSSAddressUpdated(log types.Log) (*MuseConnectorEthTSSAddressUpdated, error) {
	event := new(MuseConnectorEthTSSAddressUpdated)
	if err := _MuseConnectorEth.contract.UnpackLog(event, "TSSAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorEthTSSAddressUpdaterUpdatedIterator is returned from FilterTSSAddressUpdaterUpdated and is used to iterate over the raw logs and unpacked data for TSSAddressUpdaterUpdated events raised by the MuseConnectorEth contract.
type MuseConnectorEthTSSAddressUpdaterUpdatedIterator struct {
	Event *MuseConnectorEthTSSAddressUpdaterUpdated // Event containing the contract specifics and raw log

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
func (it *MuseConnectorEthTSSAddressUpdaterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorEthTSSAddressUpdaterUpdated)
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
		it.Event = new(MuseConnectorEthTSSAddressUpdaterUpdated)
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
func (it *MuseConnectorEthTSSAddressUpdaterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorEthTSSAddressUpdaterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorEthTSSAddressUpdaterUpdated represents a TSSAddressUpdaterUpdated event raised by the MuseConnectorEth contract.
type MuseConnectorEthTSSAddressUpdaterUpdated struct {
	CallerAddress        common.Address
	NewTssUpdaterAddress common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterTSSAddressUpdaterUpdated is a free log retrieval operation binding the contract event 0x5104c9abdc7d111c2aeb4ce890ac70274a4be2ee83f46a62551be5e6ebc82dd0.
//
// Solidity: event TSSAddressUpdaterUpdated(address callerAddress, address newTssUpdaterAddress)
func (_MuseConnectorEth *MuseConnectorEthFilterer) FilterTSSAddressUpdaterUpdated(opts *bind.FilterOpts) (*MuseConnectorEthTSSAddressUpdaterUpdatedIterator, error) {

	logs, sub, err := _MuseConnectorEth.contract.FilterLogs(opts, "TSSAddressUpdaterUpdated")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthTSSAddressUpdaterUpdatedIterator{contract: _MuseConnectorEth.contract, event: "TSSAddressUpdaterUpdated", logs: logs, sub: sub}, nil
}

// WatchTSSAddressUpdaterUpdated is a free log subscription operation binding the contract event 0x5104c9abdc7d111c2aeb4ce890ac70274a4be2ee83f46a62551be5e6ebc82dd0.
//
// Solidity: event TSSAddressUpdaterUpdated(address callerAddress, address newTssUpdaterAddress)
func (_MuseConnectorEth *MuseConnectorEthFilterer) WatchTSSAddressUpdaterUpdated(opts *bind.WatchOpts, sink chan<- *MuseConnectorEthTSSAddressUpdaterUpdated) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorEth.contract.WatchLogs(opts, "TSSAddressUpdaterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorEthTSSAddressUpdaterUpdated)
				if err := _MuseConnectorEth.contract.UnpackLog(event, "TSSAddressUpdaterUpdated", log); err != nil {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) ParseTSSAddressUpdaterUpdated(log types.Log) (*MuseConnectorEthTSSAddressUpdaterUpdated, error) {
	event := new(MuseConnectorEthTSSAddressUpdaterUpdated)
	if err := _MuseConnectorEth.contract.UnpackLog(event, "TSSAddressUpdaterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorEthUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the MuseConnectorEth contract.
type MuseConnectorEthUnpausedIterator struct {
	Event *MuseConnectorEthUnpaused // Event containing the contract specifics and raw log

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
func (it *MuseConnectorEthUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorEthUnpaused)
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
		it.Event = new(MuseConnectorEthUnpaused)
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
func (it *MuseConnectorEthUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorEthUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorEthUnpaused represents a Unpaused event raised by the MuseConnectorEth contract.
type MuseConnectorEthUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MuseConnectorEth *MuseConnectorEthFilterer) FilterUnpaused(opts *bind.FilterOpts) (*MuseConnectorEthUnpausedIterator, error) {

	logs, sub, err := _MuseConnectorEth.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthUnpausedIterator{contract: _MuseConnectorEth.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MuseConnectorEth *MuseConnectorEthFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *MuseConnectorEthUnpaused) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorEth.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorEthUnpaused)
				if err := _MuseConnectorEth.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) ParseUnpaused(log types.Log) (*MuseConnectorEthUnpaused, error) {
	event := new(MuseConnectorEthUnpaused)
	if err := _MuseConnectorEth.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorEthMuseReceivedIterator is returned from FilterMuseReceived and is used to iterate over the raw logs and unpacked data for MuseReceived events raised by the MuseConnectorEth contract.
type MuseConnectorEthMuseReceivedIterator struct {
	Event *MuseConnectorEthMuseReceived // Event containing the contract specifics and raw log

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
func (it *MuseConnectorEthMuseReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorEthMuseReceived)
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
		it.Event = new(MuseConnectorEthMuseReceived)
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
func (it *MuseConnectorEthMuseReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorEthMuseReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorEthMuseReceived represents a MuseReceived event raised by the MuseConnectorEth contract.
type MuseConnectorEthMuseReceived struct {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) FilterMuseReceived(opts *bind.FilterOpts, sourceChainId []*big.Int, destinationAddress []common.Address, internalSendHash [][32]byte) (*MuseConnectorEthMuseReceivedIterator, error) {

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

	logs, sub, err := _MuseConnectorEth.contract.FilterLogs(opts, "MuseReceived", sourceChainIdRule, destinationAddressRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthMuseReceivedIterator{contract: _MuseConnectorEth.contract, event: "MuseReceived", logs: logs, sub: sub}, nil
}

// WatchMuseReceived is a free log subscription operation binding the contract event 0xf1302855733b40d8acb467ee990b6d56c05c80e28ebcabfa6e6f3f57cb50d698.
//
// Solidity: event MuseReceived(bytes museTxSenderAddress, uint256 indexed sourceChainId, address indexed destinationAddress, uint256 museValue, bytes message, bytes32 indexed internalSendHash)
func (_MuseConnectorEth *MuseConnectorEthFilterer) WatchMuseReceived(opts *bind.WatchOpts, sink chan<- *MuseConnectorEthMuseReceived, sourceChainId []*big.Int, destinationAddress []common.Address, internalSendHash [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _MuseConnectorEth.contract.WatchLogs(opts, "MuseReceived", sourceChainIdRule, destinationAddressRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorEthMuseReceived)
				if err := _MuseConnectorEth.contract.UnpackLog(event, "MuseReceived", log); err != nil {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) ParseMuseReceived(log types.Log) (*MuseConnectorEthMuseReceived, error) {
	event := new(MuseConnectorEthMuseReceived)
	if err := _MuseConnectorEth.contract.UnpackLog(event, "MuseReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorEthMuseRevertedIterator is returned from FilterMuseReverted and is used to iterate over the raw logs and unpacked data for MuseReverted events raised by the MuseConnectorEth contract.
type MuseConnectorEthMuseRevertedIterator struct {
	Event *MuseConnectorEthMuseReverted // Event containing the contract specifics and raw log

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
func (it *MuseConnectorEthMuseRevertedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorEthMuseReverted)
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
		it.Event = new(MuseConnectorEthMuseReverted)
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
func (it *MuseConnectorEthMuseRevertedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorEthMuseRevertedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorEthMuseReverted represents a MuseReverted event raised by the MuseConnectorEth contract.
type MuseConnectorEthMuseReverted struct {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) FilterMuseReverted(opts *bind.FilterOpts, destinationChainId []*big.Int, internalSendHash [][32]byte) (*MuseConnectorEthMuseRevertedIterator, error) {

	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	var internalSendHashRule []interface{}
	for _, internalSendHashItem := range internalSendHash {
		internalSendHashRule = append(internalSendHashRule, internalSendHashItem)
	}

	logs, sub, err := _MuseConnectorEth.contract.FilterLogs(opts, "MuseReverted", destinationChainIdRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthMuseRevertedIterator{contract: _MuseConnectorEth.contract, event: "MuseReverted", logs: logs, sub: sub}, nil
}

// WatchMuseReverted is a free log subscription operation binding the contract event 0x521fb0b407c2eb9b1375530e9b9a569889992140a688bc076aa72c1712012c88.
//
// Solidity: event MuseReverted(address museTxSenderAddress, uint256 sourceChainId, uint256 indexed destinationChainId, bytes destinationAddress, uint256 remainingMuseValue, bytes message, bytes32 indexed internalSendHash)
func (_MuseConnectorEth *MuseConnectorEthFilterer) WatchMuseReverted(opts *bind.WatchOpts, sink chan<- *MuseConnectorEthMuseReverted, destinationChainId []*big.Int, internalSendHash [][32]byte) (event.Subscription, error) {

	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	var internalSendHashRule []interface{}
	for _, internalSendHashItem := range internalSendHash {
		internalSendHashRule = append(internalSendHashRule, internalSendHashItem)
	}

	logs, sub, err := _MuseConnectorEth.contract.WatchLogs(opts, "MuseReverted", destinationChainIdRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorEthMuseReverted)
				if err := _MuseConnectorEth.contract.UnpackLog(event, "MuseReverted", log); err != nil {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) ParseMuseReverted(log types.Log) (*MuseConnectorEthMuseReverted, error) {
	event := new(MuseConnectorEthMuseReverted)
	if err := _MuseConnectorEth.contract.UnpackLog(event, "MuseReverted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorEthMuseSentIterator is returned from FilterMuseSent and is used to iterate over the raw logs and unpacked data for MuseSent events raised by the MuseConnectorEth contract.
type MuseConnectorEthMuseSentIterator struct {
	Event *MuseConnectorEthMuseSent // Event containing the contract specifics and raw log

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
func (it *MuseConnectorEthMuseSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorEthMuseSent)
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
		it.Event = new(MuseConnectorEthMuseSent)
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
func (it *MuseConnectorEthMuseSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorEthMuseSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorEthMuseSent represents a MuseSent event raised by the MuseConnectorEth contract.
type MuseConnectorEthMuseSent struct {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) FilterMuseSent(opts *bind.FilterOpts, museTxSenderAddress []common.Address, destinationChainId []*big.Int) (*MuseConnectorEthMuseSentIterator, error) {

	var museTxSenderAddressRule []interface{}
	for _, museTxSenderAddressItem := range museTxSenderAddress {
		museTxSenderAddressRule = append(museTxSenderAddressRule, museTxSenderAddressItem)
	}
	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	logs, sub, err := _MuseConnectorEth.contract.FilterLogs(opts, "MuseSent", museTxSenderAddressRule, destinationChainIdRule)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorEthMuseSentIterator{contract: _MuseConnectorEth.contract, event: "MuseSent", logs: logs, sub: sub}, nil
}

// WatchMuseSent is a free log subscription operation binding the contract event 0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4.
//
// Solidity: event MuseSent(address sourceTxOriginAddress, address indexed museTxSenderAddress, uint256 indexed destinationChainId, bytes destinationAddress, uint256 museValueAndGas, uint256 destinationGasLimit, bytes message, bytes museParams)
func (_MuseConnectorEth *MuseConnectorEthFilterer) WatchMuseSent(opts *bind.WatchOpts, sink chan<- *MuseConnectorEthMuseSent, museTxSenderAddress []common.Address, destinationChainId []*big.Int) (event.Subscription, error) {

	var museTxSenderAddressRule []interface{}
	for _, museTxSenderAddressItem := range museTxSenderAddress {
		museTxSenderAddressRule = append(museTxSenderAddressRule, museTxSenderAddressItem)
	}
	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	logs, sub, err := _MuseConnectorEth.contract.WatchLogs(opts, "MuseSent", museTxSenderAddressRule, destinationChainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorEthMuseSent)
				if err := _MuseConnectorEth.contract.UnpackLog(event, "MuseSent", log); err != nil {
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
func (_MuseConnectorEth *MuseConnectorEthFilterer) ParseMuseSent(log types.Log) (*MuseConnectorEthMuseSent, error) {
	event := new(MuseConnectorEthMuseSent)
	if err := _MuseConnectorEth.contract.UnpackLog(event, "MuseSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
