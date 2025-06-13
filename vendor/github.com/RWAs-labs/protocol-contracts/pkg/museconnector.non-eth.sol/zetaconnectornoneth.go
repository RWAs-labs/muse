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

// MuseConnectorNonEthMetaData contains all meta data concerning the MuseConnectorNonEth contract.
var MuseConnectorNonEthMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"museTokenAddress_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tssAddress_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tssAddressUpdater_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pauserAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getLockedAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onReceive\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"museValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onRevert\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"remainingMuseValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pauserAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceTssAddressUpdater\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"send\",\"inputs\":[{\"name\":\"input\",\"type\":\"tuple\",\"internalType\":\"structMuseInterfaces.SendInput\",\"components\":[{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destinationGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"museValueAndGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"museParams\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxSupply\",\"inputs\":[{\"name\":\"maxSupply_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tssAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tssAddressUpdater\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauserAddress\",\"inputs\":[{\"name\":\"pauserAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateTssAddress\",\"inputs\":[{\"name\":\"tssAddress_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"museToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"MaxSupplyUpdated\",\"inputs\":[{\"name\":\"callerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newMaxSupply\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PauserAddressUpdated\",\"inputs\":[{\"name\":\"callerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newTssAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TSSAddressUpdated\",\"inputs\":[{\"name\":\"callerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newTssAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TSSAddressUpdaterUpdated\",\"inputs\":[{\"name\":\"callerAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newTssUpdaterAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MuseReceived\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"museValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MuseReverted\",\"inputs\":[{\"name\":\"museTxSenderAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"sourceChainId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"remainingMuseValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"internalSendHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MuseSent\",\"inputs\":[{\"name\":\"sourceTxOriginAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"museTxSenderAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"destinationAddress\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"museValueAndGas\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationGasLimit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"message\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"museParams\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CallerIsNotPauser\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CallerIsNotTss\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CallerIsNotTssOrUpdater\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CallerIsNotTssUpdater\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExceedsMaxSupply\",\"inputs\":[{\"name\":\"maxSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MuseTransferError\",\"inputs\":[]}]",
	Bin: "0x60a060405260001960035534801561001657600080fd5b506040516119323803806119328339810160408190526100359161011a565b6000805460ff19169055838383836001600160a01b038416158061006057506001600160a01b038316155b8061007257506001600160a01b038216155b8061008457506001600160a01b038116155b156100a25760405163e6c4247b60e01b815260040160405180910390fd5b6001600160a01b03938416608052600180549385166001600160a01b0319948516179055600280549285169290931691909117909155600080549190921661010002610100600160a81b03199091161790555061016e92505050565b80516001600160a01b038116811461011557600080fd5b919050565b6000806000806080858703121561013057600080fd5b610139856100fe565b9350610147602086016100fe565b9250610155604086016100fe565b9150610163606086016100fe565b905092959194509250565b60805161177f6101b36000396000818161010a015281816102ad0152818161038f015281816104b701528181610bf201528181610d1a0152610f44015261177f6000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c80636f8b44b011610097578063942a5e1611610066578063942a5e1614610228578063d5abeb011461023b578063ec02690114610244578063f7fb869b1461025757600080fd5b80636f8b44b0146101f2578063779e3b63146102055780638456cb591461020d5780639122c3441461021557600080fd5b80633f4ba83a116100d35780633f4ba83a146101a15780635b112591146101a95780635c975abb146101c95780636128480f146101df57600080fd5b806321e093b114610105578063252bc8861461015657806329dd214d1461016c578063328a01d014610181575b600080fd5b61012c7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020015b60405180910390f35b61015e61027c565b60405190815260200161014d565b61017f61017a3660046111eb565b610332565b005b60025461012c9073ffffffffffffffffffffffffffffffffffffffff1681565b61017f6106b0565b60015461012c9073ffffffffffffffffffffffffffffffffffffffff1681565b60005460ff16604051901515815260200161014d565b61017f6101ed36600461128b565b610712565b61017f6102003660046112ad565b61083f565b61017f6108cd565b61017f6109f2565b61017f61022336600461128b565b610a52565b61017f6102363660046112c6565b610b92565b61015e60035481565b61017f610252366004611370565b610f04565b60005461012c90610100900473ffffffffffffffffffffffffffffffffffffffff1681565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015610309573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061032d91906113ab565b905090565b60015473ffffffffffffffffffffffffffffffffffffffff16331461038a576040517fff70ace20000000000000000000000000000000000000000000000000000000081523360048201526024015b60405180910390fd5b6003547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166318160ddd6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156103f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061041c91906113ab565b61042690866113c4565b1115610464576003546040517f3d3dbc8300000000000000000000000000000000000000000000000000000000815260040161038191815260200190565b6040517f1e458bee00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff868116600483015260248201869052604482018390527f00000000000000000000000000000000000000000000000000000000000000001690631e458bee90606401600060405180830381600087803b1580156104fb57600080fd5b505af115801561050f573d6000803e3d6000fd5b50508315915061064e9050578473ffffffffffffffffffffffffffffffffffffffff16633749c51a6040518060a001604052808b8b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060208082018b905273ffffffffffffffffffffffffffffffffffffffff8a16604080840191909152606083018a90528051601f890183900483028101830190915287815260809092019190889088908190840183828082843760009201919091525050509152506040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b16815261061b9190600401611468565b600060405180830381600087803b15801561063557600080fd5b505af1158015610649573d6000803e3d6000fd5b505050505b808573ffffffffffffffffffffffffffffffffffffffff16877ff1302855733b40d8acb467ee990b6d56c05c80e28ebcabfa6e6f3f57cb50d6988b8b89898960405161069e959493929190611545565b60405180910390a45050505050505050565b600054610100900473ffffffffffffffffffffffffffffffffffffffff163314610708576040517f4677a0d3000000000000000000000000000000000000000000000000000000008152336004820152602401610381565b61071061102d565b565b600054610100900473ffffffffffffffffffffffffffffffffffffffff16331461076a576040517f4677a0d3000000000000000000000000000000000000000000000000000000008152336004820152602401610381565b73ffffffffffffffffffffffffffffffffffffffff81166107b7576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600080547fffffffffffffffffffffff0000000000000000000000000000000000000000ff1661010073ffffffffffffffffffffffffffffffffffffffff8416908102919091179091556040805133815260208101929092527fd41d83655d484bdf299598751c371b2d92088667266fe3774b25a97bdd5d039791015b60405180910390a150565b60015473ffffffffffffffffffffffffffffffffffffffff163314610892576040517fff70ace2000000000000000000000000000000000000000000000000000000008152336004820152602401610381565b600381905560408051338152602081018390527f26843c619c87f9021bcc4ec5143177198076d9da3c13ce1bb2e941644c69d42e9101610834565b60025473ffffffffffffffffffffffffffffffffffffffff163314610920576040517fe700765e000000000000000000000000000000000000000000000000000000008152336004820152602401610381565b60015473ffffffffffffffffffffffffffffffffffffffff1661096f576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600154600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691821790556040805133815260208101929092527f5104c9abdc7d111c2aeb4ce890ac70274a4be2ee83f46a62551be5e6ebc82dd091015b60405180910390a1565b600054610100900473ffffffffffffffffffffffffffffffffffffffff163314610a4a576040517f4677a0d3000000000000000000000000000000000000000000000000000000008152336004820152602401610381565b6107106110a5565b60015473ffffffffffffffffffffffffffffffffffffffff163314801590610a92575060025473ffffffffffffffffffffffffffffffffffffffff163314155b15610acb576040517fcdfcef97000000000000000000000000000000000000000000000000000000008152336004820152602401610381565b73ffffffffffffffffffffffffffffffffffffffff8116610b18576040517fe6c4247b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040805133815260208101929092527fe79965b5c67dcfb2cf5fe152715e4a7256cee62a3d5dd8484fd8a8539eb8beff9101610834565b610b9a611100565b60015473ffffffffffffffffffffffffffffffffffffffff163314610bed576040517fff70ace2000000000000000000000000000000000000000000000000000000008152336004820152602401610381565b6003547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166318160ddd6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610c5b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c7f91906113ab565b610c8990866113c4565b1115610cc7576003546040517f3d3dbc8300000000000000000000000000000000000000000000000000000000815260040161038191815260200190565b6040517f1e458bee00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8a8116600483015260248201869052604482018390527f00000000000000000000000000000000000000000000000000000000000000001690631e458bee90606401600060405180830381600087803b158015610d5e57600080fd5b505af1158015610d72573d6000803e3d6000fd5b505083159150610eb49050578873ffffffffffffffffffffffffffffffffffffffff16633ff0693c6040518060c001604052808c73ffffffffffffffffffffffffffffffffffffffff1681526020018b81526020018a8a8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050509082525060208082018a905260408083018a90528051601f890183900483028101830190915287815260609092019190889088908190840183828082843760009201919091525050509152506040517fffffffff0000000000000000000000000000000000000000000000000000000060e084901b168152610e81919060040161157e565b600060405180830381600087803b158015610e9b57600080fd5b505af1158015610eaf573d6000803e3d6000fd5b505050505b80857f521fb0b407c2eb9b1375530e9b9a569889992140a688bc076aa72c1712012c888b8b8b8b8a8a8a604051610ef19796959493929190611613565b60405180910390a3505050505050505050565b610f0c611100565b6040517f79cc6790000000000000000000000000000000000000000000000000000000008152336004820152608082013560248201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906379cc679090604401600060405180830381600087803b158015610f9d57600080fd5b505af1158015610fb1573d6000803e3d6000fd5b5050823591503390507f7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e432610fe96020860186611670565b6080870135604088013561100060608a018a611670565b61100d60a08c018c611670565b604051611022999897969594939291906116d5565b60405180910390a350565b61103561113d565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016109e8565b6110ad611100565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586110803390565b60005460ff1615610710576040517fd93c066500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60005460ff16610710576040517f8dfc202b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008083601f84011261118b57600080fd5b50813567ffffffffffffffff8111156111a357600080fd5b6020830191508360208285010111156111bb57600080fd5b9250929050565b803573ffffffffffffffffffffffffffffffffffffffff811681146111e657600080fd5b919050565b60008060008060008060008060c0898b03121561120757600080fd5b883567ffffffffffffffff81111561121e57600080fd5b61122a8b828c01611179565b9099509750506020890135955061124360408a016111c2565b945060608901359350608089013567ffffffffffffffff81111561126657600080fd5b6112728b828c01611179565b999c989b50969995989497949560a00135949350505050565b60006020828403121561129d57600080fd5b6112a6826111c2565b9392505050565b6000602082840312156112bf57600080fd5b5035919050565b600080600080600080600080600060e08a8c0312156112e457600080fd5b6112ed8a6111c2565b985060208a0135975060408a013567ffffffffffffffff81111561131057600080fd5b61131c8c828d01611179565b90985096505060608a0135945060808a0135935060a08a013567ffffffffffffffff81111561134a57600080fd5b6113568c828d01611179565b9a9d999c50979a9699959894979660c00135949350505050565b60006020828403121561138257600080fd5b813567ffffffffffffffff81111561139957600080fd5b820160c081850312156112a657600080fd5b6000602082840312156113bd57600080fd5b5051919050565b808201808211156113fe577f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b92915050565b6000815180845260005b8181101561142a5760208185018101518683018201520161140e565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b602081526000825160a0602084015261148460c0840182611404565b90506020840151604084015273ffffffffffffffffffffffffffffffffffffffff60408501511660608401526060840151608084015260808401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160a08501526114f38282611404565b95945050505050565b8183528181602085013750600060208284010152600060207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116840101905092915050565b6060815260006115596060830187896114fc565b85602084015282810360408401526115728185876114fc565b98975050505050505050565b6020815273ffffffffffffffffffffffffffffffffffffffff8251166020820152602082015160408201526000604083015160c060608401526115c460e0840182611404565b905060608401516080840152608084015160a084015260a08401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08483030160c08501526114f38282611404565b73ffffffffffffffffffffffffffffffffffffffff8816815286602082015260a06040820152600061164960a0830187896114fc565b85606084015282810360808401526116628185876114fc565b9a9950505050505050505050565b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126116a557600080fd5b83018035915067ffffffffffffffff8211156116c057600080fd5b6020019150368190038213156111bb57600080fd5b73ffffffffffffffffffffffffffffffffffffffff8a16815260c06020820152600061170560c083018a8c6114fc565b88604084015287606084015282810360808401526117248187896114fc565b905082810360a08401526117398185876114fc565b9c9b50505050505050505050505056fea264697066735822122056c9cd3efa46f596cf4d0799036a0a9c54115921510304afaebdd3aa568dd85464736f6c634300081a0033",
}

// MuseConnectorNonEthABI is the input ABI used to generate the binding from.
// Deprecated: Use MuseConnectorNonEthMetaData.ABI instead.
var MuseConnectorNonEthABI = MuseConnectorNonEthMetaData.ABI

// MuseConnectorNonEthBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MuseConnectorNonEthMetaData.Bin instead.
var MuseConnectorNonEthBin = MuseConnectorNonEthMetaData.Bin

// DeployMuseConnectorNonEth deploys a new Ethereum contract, binding an instance of MuseConnectorNonEth to it.
func DeployMuseConnectorNonEth(auth *bind.TransactOpts, backend bind.ContractBackend, museTokenAddress_ common.Address, tssAddress_ common.Address, tssAddressUpdater_ common.Address, pauserAddress_ common.Address) (common.Address, *types.Transaction, *MuseConnectorNonEth, error) {
	parsed, err := MuseConnectorNonEthMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MuseConnectorNonEthBin), backend, museTokenAddress_, tssAddress_, tssAddressUpdater_, pauserAddress_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MuseConnectorNonEth{MuseConnectorNonEthCaller: MuseConnectorNonEthCaller{contract: contract}, MuseConnectorNonEthTransactor: MuseConnectorNonEthTransactor{contract: contract}, MuseConnectorNonEthFilterer: MuseConnectorNonEthFilterer{contract: contract}}, nil
}

// MuseConnectorNonEth is an auto generated Go binding around an Ethereum contract.
type MuseConnectorNonEth struct {
	MuseConnectorNonEthCaller     // Read-only binding to the contract
	MuseConnectorNonEthTransactor // Write-only binding to the contract
	MuseConnectorNonEthFilterer   // Log filterer for contract events
}

// MuseConnectorNonEthCaller is an auto generated read-only Go binding around an Ethereum contract.
type MuseConnectorNonEthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseConnectorNonEthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MuseConnectorNonEthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseConnectorNonEthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MuseConnectorNonEthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MuseConnectorNonEthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MuseConnectorNonEthSession struct {
	Contract     *MuseConnectorNonEth // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MuseConnectorNonEthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MuseConnectorNonEthCallerSession struct {
	Contract *MuseConnectorNonEthCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// MuseConnectorNonEthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MuseConnectorNonEthTransactorSession struct {
	Contract     *MuseConnectorNonEthTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// MuseConnectorNonEthRaw is an auto generated low-level Go binding around an Ethereum contract.
type MuseConnectorNonEthRaw struct {
	Contract *MuseConnectorNonEth // Generic contract binding to access the raw methods on
}

// MuseConnectorNonEthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MuseConnectorNonEthCallerRaw struct {
	Contract *MuseConnectorNonEthCaller // Generic read-only contract binding to access the raw methods on
}

// MuseConnectorNonEthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MuseConnectorNonEthTransactorRaw struct {
	Contract *MuseConnectorNonEthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMuseConnectorNonEth creates a new instance of MuseConnectorNonEth, bound to a specific deployed contract.
func NewMuseConnectorNonEth(address common.Address, backend bind.ContractBackend) (*MuseConnectorNonEth, error) {
	contract, err := bindMuseConnectorNonEth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEth{MuseConnectorNonEthCaller: MuseConnectorNonEthCaller{contract: contract}, MuseConnectorNonEthTransactor: MuseConnectorNonEthTransactor{contract: contract}, MuseConnectorNonEthFilterer: MuseConnectorNonEthFilterer{contract: contract}}, nil
}

// NewMuseConnectorNonEthCaller creates a new read-only instance of MuseConnectorNonEth, bound to a specific deployed contract.
func NewMuseConnectorNonEthCaller(address common.Address, caller bind.ContractCaller) (*MuseConnectorNonEthCaller, error) {
	contract, err := bindMuseConnectorNonEth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthCaller{contract: contract}, nil
}

// NewMuseConnectorNonEthTransactor creates a new write-only instance of MuseConnectorNonEth, bound to a specific deployed contract.
func NewMuseConnectorNonEthTransactor(address common.Address, transactor bind.ContractTransactor) (*MuseConnectorNonEthTransactor, error) {
	contract, err := bindMuseConnectorNonEth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthTransactor{contract: contract}, nil
}

// NewMuseConnectorNonEthFilterer creates a new log filterer instance of MuseConnectorNonEth, bound to a specific deployed contract.
func NewMuseConnectorNonEthFilterer(address common.Address, filterer bind.ContractFilterer) (*MuseConnectorNonEthFilterer, error) {
	contract, err := bindMuseConnectorNonEth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthFilterer{contract: contract}, nil
}

// bindMuseConnectorNonEth binds a generic wrapper to an already deployed contract.
func bindMuseConnectorNonEth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MuseConnectorNonEthMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MuseConnectorNonEth *MuseConnectorNonEthRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MuseConnectorNonEth.Contract.MuseConnectorNonEthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MuseConnectorNonEth *MuseConnectorNonEthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.MuseConnectorNonEthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MuseConnectorNonEth *MuseConnectorNonEthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.MuseConnectorNonEthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MuseConnectorNonEth *MuseConnectorNonEthCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MuseConnectorNonEth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.contract.Transact(opts, method, params...)
}

// GetLockedAmount is a free data retrieval call binding the contract method 0x252bc886.
//
// Solidity: function getLockedAmount() view returns(uint256)
func (_MuseConnectorNonEth *MuseConnectorNonEthCaller) GetLockedAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MuseConnectorNonEth.contract.Call(opts, &out, "getLockedAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockedAmount is a free data retrieval call binding the contract method 0x252bc886.
//
// Solidity: function getLockedAmount() view returns(uint256)
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) GetLockedAmount() (*big.Int, error) {
	return _MuseConnectorNonEth.Contract.GetLockedAmount(&_MuseConnectorNonEth.CallOpts)
}

// GetLockedAmount is a free data retrieval call binding the contract method 0x252bc886.
//
// Solidity: function getLockedAmount() view returns(uint256)
func (_MuseConnectorNonEth *MuseConnectorNonEthCallerSession) GetLockedAmount() (*big.Int, error) {
	return _MuseConnectorNonEth.Contract.GetLockedAmount(&_MuseConnectorNonEth.CallOpts)
}

// MaxSupply is a free data retrieval call binding the contract method 0xd5abeb01.
//
// Solidity: function maxSupply() view returns(uint256)
func (_MuseConnectorNonEth *MuseConnectorNonEthCaller) MaxSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MuseConnectorNonEth.contract.Call(opts, &out, "maxSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxSupply is a free data retrieval call binding the contract method 0xd5abeb01.
//
// Solidity: function maxSupply() view returns(uint256)
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) MaxSupply() (*big.Int, error) {
	return _MuseConnectorNonEth.Contract.MaxSupply(&_MuseConnectorNonEth.CallOpts)
}

// MaxSupply is a free data retrieval call binding the contract method 0xd5abeb01.
//
// Solidity: function maxSupply() view returns(uint256)
func (_MuseConnectorNonEth *MuseConnectorNonEthCallerSession) MaxSupply() (*big.Int, error) {
	return _MuseConnectorNonEth.Contract.MaxSupply(&_MuseConnectorNonEth.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MuseConnectorNonEth *MuseConnectorNonEthCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _MuseConnectorNonEth.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) Paused() (bool, error) {
	return _MuseConnectorNonEth.Contract.Paused(&_MuseConnectorNonEth.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_MuseConnectorNonEth *MuseConnectorNonEthCallerSession) Paused() (bool, error) {
	return _MuseConnectorNonEth.Contract.Paused(&_MuseConnectorNonEth.CallOpts)
}

// PauserAddress is a free data retrieval call binding the contract method 0xf7fb869b.
//
// Solidity: function pauserAddress() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthCaller) PauserAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorNonEth.contract.Call(opts, &out, "pauserAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PauserAddress is a free data retrieval call binding the contract method 0xf7fb869b.
//
// Solidity: function pauserAddress() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) PauserAddress() (common.Address, error) {
	return _MuseConnectorNonEth.Contract.PauserAddress(&_MuseConnectorNonEth.CallOpts)
}

// PauserAddress is a free data retrieval call binding the contract method 0xf7fb869b.
//
// Solidity: function pauserAddress() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthCallerSession) PauserAddress() (common.Address, error) {
	return _MuseConnectorNonEth.Contract.PauserAddress(&_MuseConnectorNonEth.CallOpts)
}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthCaller) TssAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorNonEth.contract.Call(opts, &out, "tssAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) TssAddress() (common.Address, error) {
	return _MuseConnectorNonEth.Contract.TssAddress(&_MuseConnectorNonEth.CallOpts)
}

// TssAddress is a free data retrieval call binding the contract method 0x5b112591.
//
// Solidity: function tssAddress() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthCallerSession) TssAddress() (common.Address, error) {
	return _MuseConnectorNonEth.Contract.TssAddress(&_MuseConnectorNonEth.CallOpts)
}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthCaller) TssAddressUpdater(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorNonEth.contract.Call(opts, &out, "tssAddressUpdater")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) TssAddressUpdater() (common.Address, error) {
	return _MuseConnectorNonEth.Contract.TssAddressUpdater(&_MuseConnectorNonEth.CallOpts)
}

// TssAddressUpdater is a free data retrieval call binding the contract method 0x328a01d0.
//
// Solidity: function tssAddressUpdater() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthCallerSession) TssAddressUpdater() (common.Address, error) {
	return _MuseConnectorNonEth.Contract.TssAddressUpdater(&_MuseConnectorNonEth.CallOpts)
}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthCaller) MuseToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MuseConnectorNonEth.contract.Call(opts, &out, "museToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) MuseToken() (common.Address, error) {
	return _MuseConnectorNonEth.Contract.MuseToken(&_MuseConnectorNonEth.CallOpts)
}

// MuseToken is a free data retrieval call binding the contract method 0x21e093b1.
//
// Solidity: function museToken() view returns(address)
func (_MuseConnectorNonEth *MuseConnectorNonEthCallerSession) MuseToken() (common.Address, error) {
	return _MuseConnectorNonEth.Contract.MuseToken(&_MuseConnectorNonEth.CallOpts)
}

// OnReceive is a paid mutator transaction binding the contract method 0x29dd214d.
//
// Solidity: function onReceive(bytes museTxSenderAddress, uint256 sourceChainId, address destinationAddress, uint256 museValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactor) OnReceive(opts *bind.TransactOpts, museTxSenderAddress []byte, sourceChainId *big.Int, destinationAddress common.Address, museValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorNonEth.contract.Transact(opts, "onReceive", museTxSenderAddress, sourceChainId, destinationAddress, museValue, message, internalSendHash)
}

// OnReceive is a paid mutator transaction binding the contract method 0x29dd214d.
//
// Solidity: function onReceive(bytes museTxSenderAddress, uint256 sourceChainId, address destinationAddress, uint256 museValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) OnReceive(museTxSenderAddress []byte, sourceChainId *big.Int, destinationAddress common.Address, museValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.OnReceive(&_MuseConnectorNonEth.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, museValue, message, internalSendHash)
}

// OnReceive is a paid mutator transaction binding the contract method 0x29dd214d.
//
// Solidity: function onReceive(bytes museTxSenderAddress, uint256 sourceChainId, address destinationAddress, uint256 museValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorSession) OnReceive(museTxSenderAddress []byte, sourceChainId *big.Int, destinationAddress common.Address, museValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.OnReceive(&_MuseConnectorNonEth.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, museValue, message, internalSendHash)
}

// OnRevert is a paid mutator transaction binding the contract method 0x942a5e16.
//
// Solidity: function onRevert(address museTxSenderAddress, uint256 sourceChainId, bytes destinationAddress, uint256 destinationChainId, uint256 remainingMuseValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactor) OnRevert(opts *bind.TransactOpts, museTxSenderAddress common.Address, sourceChainId *big.Int, destinationAddress []byte, destinationChainId *big.Int, remainingMuseValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorNonEth.contract.Transact(opts, "onRevert", museTxSenderAddress, sourceChainId, destinationAddress, destinationChainId, remainingMuseValue, message, internalSendHash)
}

// OnRevert is a paid mutator transaction binding the contract method 0x942a5e16.
//
// Solidity: function onRevert(address museTxSenderAddress, uint256 sourceChainId, bytes destinationAddress, uint256 destinationChainId, uint256 remainingMuseValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) OnRevert(museTxSenderAddress common.Address, sourceChainId *big.Int, destinationAddress []byte, destinationChainId *big.Int, remainingMuseValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.OnRevert(&_MuseConnectorNonEth.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, destinationChainId, remainingMuseValue, message, internalSendHash)
}

// OnRevert is a paid mutator transaction binding the contract method 0x942a5e16.
//
// Solidity: function onRevert(address museTxSenderAddress, uint256 sourceChainId, bytes destinationAddress, uint256 destinationChainId, uint256 remainingMuseValue, bytes message, bytes32 internalSendHash) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorSession) OnRevert(museTxSenderAddress common.Address, sourceChainId *big.Int, destinationAddress []byte, destinationChainId *big.Int, remainingMuseValue *big.Int, message []byte, internalSendHash [32]byte) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.OnRevert(&_MuseConnectorNonEth.TransactOpts, museTxSenderAddress, sourceChainId, destinationAddress, destinationChainId, remainingMuseValue, message, internalSendHash)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorNonEth.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) Pause() (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.Pause(&_MuseConnectorNonEth.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorSession) Pause() (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.Pause(&_MuseConnectorNonEth.TransactOpts)
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactor) RenounceTssAddressUpdater(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorNonEth.contract.Transact(opts, "renounceTssAddressUpdater")
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) RenounceTssAddressUpdater() (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.RenounceTssAddressUpdater(&_MuseConnectorNonEth.TransactOpts)
}

// RenounceTssAddressUpdater is a paid mutator transaction binding the contract method 0x779e3b63.
//
// Solidity: function renounceTssAddressUpdater() returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorSession) RenounceTssAddressUpdater() (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.RenounceTssAddressUpdater(&_MuseConnectorNonEth.TransactOpts)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactor) Send(opts *bind.TransactOpts, input MuseInterfacesSendInput) (*types.Transaction, error) {
	return _MuseConnectorNonEth.contract.Transact(opts, "send", input)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) Send(input MuseInterfacesSendInput) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.Send(&_MuseConnectorNonEth.TransactOpts, input)
}

// Send is a paid mutator transaction binding the contract method 0xec026901.
//
// Solidity: function send((uint256,bytes,uint256,bytes,uint256,bytes) input) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorSession) Send(input MuseInterfacesSendInput) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.Send(&_MuseConnectorNonEth.TransactOpts, input)
}

// SetMaxSupply is a paid mutator transaction binding the contract method 0x6f8b44b0.
//
// Solidity: function setMaxSupply(uint256 maxSupply_) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactor) SetMaxSupply(opts *bind.TransactOpts, maxSupply_ *big.Int) (*types.Transaction, error) {
	return _MuseConnectorNonEth.contract.Transact(opts, "setMaxSupply", maxSupply_)
}

// SetMaxSupply is a paid mutator transaction binding the contract method 0x6f8b44b0.
//
// Solidity: function setMaxSupply(uint256 maxSupply_) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) SetMaxSupply(maxSupply_ *big.Int) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.SetMaxSupply(&_MuseConnectorNonEth.TransactOpts, maxSupply_)
}

// SetMaxSupply is a paid mutator transaction binding the contract method 0x6f8b44b0.
//
// Solidity: function setMaxSupply(uint256 maxSupply_) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorSession) SetMaxSupply(maxSupply_ *big.Int) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.SetMaxSupply(&_MuseConnectorNonEth.TransactOpts, maxSupply_)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MuseConnectorNonEth.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) Unpause() (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.Unpause(&_MuseConnectorNonEth.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorSession) Unpause() (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.Unpause(&_MuseConnectorNonEth.TransactOpts)
}

// UpdatePauserAddress is a paid mutator transaction binding the contract method 0x6128480f.
//
// Solidity: function updatePauserAddress(address pauserAddress_) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactor) UpdatePauserAddress(opts *bind.TransactOpts, pauserAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorNonEth.contract.Transact(opts, "updatePauserAddress", pauserAddress_)
}

// UpdatePauserAddress is a paid mutator transaction binding the contract method 0x6128480f.
//
// Solidity: function updatePauserAddress(address pauserAddress_) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) UpdatePauserAddress(pauserAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.UpdatePauserAddress(&_MuseConnectorNonEth.TransactOpts, pauserAddress_)
}

// UpdatePauserAddress is a paid mutator transaction binding the contract method 0x6128480f.
//
// Solidity: function updatePauserAddress(address pauserAddress_) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorSession) UpdatePauserAddress(pauserAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.UpdatePauserAddress(&_MuseConnectorNonEth.TransactOpts, pauserAddress_)
}

// UpdateTssAddress is a paid mutator transaction binding the contract method 0x9122c344.
//
// Solidity: function updateTssAddress(address tssAddress_) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactor) UpdateTssAddress(opts *bind.TransactOpts, tssAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorNonEth.contract.Transact(opts, "updateTssAddress", tssAddress_)
}

// UpdateTssAddress is a paid mutator transaction binding the contract method 0x9122c344.
//
// Solidity: function updateTssAddress(address tssAddress_) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthSession) UpdateTssAddress(tssAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.UpdateTssAddress(&_MuseConnectorNonEth.TransactOpts, tssAddress_)
}

// UpdateTssAddress is a paid mutator transaction binding the contract method 0x9122c344.
//
// Solidity: function updateTssAddress(address tssAddress_) returns()
func (_MuseConnectorNonEth *MuseConnectorNonEthTransactorSession) UpdateTssAddress(tssAddress_ common.Address) (*types.Transaction, error) {
	return _MuseConnectorNonEth.Contract.UpdateTssAddress(&_MuseConnectorNonEth.TransactOpts, tssAddress_)
}

// MuseConnectorNonEthMaxSupplyUpdatedIterator is returned from FilterMaxSupplyUpdated and is used to iterate over the raw logs and unpacked data for MaxSupplyUpdated events raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthMaxSupplyUpdatedIterator struct {
	Event *MuseConnectorNonEthMaxSupplyUpdated // Event containing the contract specifics and raw log

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
func (it *MuseConnectorNonEthMaxSupplyUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorNonEthMaxSupplyUpdated)
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
		it.Event = new(MuseConnectorNonEthMaxSupplyUpdated)
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
func (it *MuseConnectorNonEthMaxSupplyUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorNonEthMaxSupplyUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorNonEthMaxSupplyUpdated represents a MaxSupplyUpdated event raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthMaxSupplyUpdated struct {
	CallerAddress common.Address
	NewMaxSupply  *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterMaxSupplyUpdated is a free log retrieval operation binding the contract event 0x26843c619c87f9021bcc4ec5143177198076d9da3c13ce1bb2e941644c69d42e.
//
// Solidity: event MaxSupplyUpdated(address callerAddress, uint256 newMaxSupply)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) FilterMaxSupplyUpdated(opts *bind.FilterOpts) (*MuseConnectorNonEthMaxSupplyUpdatedIterator, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.FilterLogs(opts, "MaxSupplyUpdated")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthMaxSupplyUpdatedIterator{contract: _MuseConnectorNonEth.contract, event: "MaxSupplyUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxSupplyUpdated is a free log subscription operation binding the contract event 0x26843c619c87f9021bcc4ec5143177198076d9da3c13ce1bb2e941644c69d42e.
//
// Solidity: event MaxSupplyUpdated(address callerAddress, uint256 newMaxSupply)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) WatchMaxSupplyUpdated(opts *bind.WatchOpts, sink chan<- *MuseConnectorNonEthMaxSupplyUpdated) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.WatchLogs(opts, "MaxSupplyUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorNonEthMaxSupplyUpdated)
				if err := _MuseConnectorNonEth.contract.UnpackLog(event, "MaxSupplyUpdated", log); err != nil {
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

// ParseMaxSupplyUpdated is a log parse operation binding the contract event 0x26843c619c87f9021bcc4ec5143177198076d9da3c13ce1bb2e941644c69d42e.
//
// Solidity: event MaxSupplyUpdated(address callerAddress, uint256 newMaxSupply)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) ParseMaxSupplyUpdated(log types.Log) (*MuseConnectorNonEthMaxSupplyUpdated, error) {
	event := new(MuseConnectorNonEthMaxSupplyUpdated)
	if err := _MuseConnectorNonEth.contract.UnpackLog(event, "MaxSupplyUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorNonEthPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthPausedIterator struct {
	Event *MuseConnectorNonEthPaused // Event containing the contract specifics and raw log

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
func (it *MuseConnectorNonEthPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorNonEthPaused)
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
		it.Event = new(MuseConnectorNonEthPaused)
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
func (it *MuseConnectorNonEthPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorNonEthPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorNonEthPaused represents a Paused event raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) FilterPaused(opts *bind.FilterOpts) (*MuseConnectorNonEthPausedIterator, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthPausedIterator{contract: _MuseConnectorNonEth.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *MuseConnectorNonEthPaused) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorNonEthPaused)
				if err := _MuseConnectorNonEth.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) ParsePaused(log types.Log) (*MuseConnectorNonEthPaused, error) {
	event := new(MuseConnectorNonEthPaused)
	if err := _MuseConnectorNonEth.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorNonEthPauserAddressUpdatedIterator is returned from FilterPauserAddressUpdated and is used to iterate over the raw logs and unpacked data for PauserAddressUpdated events raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthPauserAddressUpdatedIterator struct {
	Event *MuseConnectorNonEthPauserAddressUpdated // Event containing the contract specifics and raw log

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
func (it *MuseConnectorNonEthPauserAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorNonEthPauserAddressUpdated)
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
		it.Event = new(MuseConnectorNonEthPauserAddressUpdated)
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
func (it *MuseConnectorNonEthPauserAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorNonEthPauserAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorNonEthPauserAddressUpdated represents a PauserAddressUpdated event raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthPauserAddressUpdated struct {
	CallerAddress common.Address
	NewTssAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterPauserAddressUpdated is a free log retrieval operation binding the contract event 0xd41d83655d484bdf299598751c371b2d92088667266fe3774b25a97bdd5d0397.
//
// Solidity: event PauserAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) FilterPauserAddressUpdated(opts *bind.FilterOpts) (*MuseConnectorNonEthPauserAddressUpdatedIterator, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.FilterLogs(opts, "PauserAddressUpdated")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthPauserAddressUpdatedIterator{contract: _MuseConnectorNonEth.contract, event: "PauserAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchPauserAddressUpdated is a free log subscription operation binding the contract event 0xd41d83655d484bdf299598751c371b2d92088667266fe3774b25a97bdd5d0397.
//
// Solidity: event PauserAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) WatchPauserAddressUpdated(opts *bind.WatchOpts, sink chan<- *MuseConnectorNonEthPauserAddressUpdated) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.WatchLogs(opts, "PauserAddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorNonEthPauserAddressUpdated)
				if err := _MuseConnectorNonEth.contract.UnpackLog(event, "PauserAddressUpdated", log); err != nil {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) ParsePauserAddressUpdated(log types.Log) (*MuseConnectorNonEthPauserAddressUpdated, error) {
	event := new(MuseConnectorNonEthPauserAddressUpdated)
	if err := _MuseConnectorNonEth.contract.UnpackLog(event, "PauserAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorNonEthTSSAddressUpdatedIterator is returned from FilterTSSAddressUpdated and is used to iterate over the raw logs and unpacked data for TSSAddressUpdated events raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthTSSAddressUpdatedIterator struct {
	Event *MuseConnectorNonEthTSSAddressUpdated // Event containing the contract specifics and raw log

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
func (it *MuseConnectorNonEthTSSAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorNonEthTSSAddressUpdated)
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
		it.Event = new(MuseConnectorNonEthTSSAddressUpdated)
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
func (it *MuseConnectorNonEthTSSAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorNonEthTSSAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorNonEthTSSAddressUpdated represents a TSSAddressUpdated event raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthTSSAddressUpdated struct {
	CallerAddress common.Address
	NewTssAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTSSAddressUpdated is a free log retrieval operation binding the contract event 0xe79965b5c67dcfb2cf5fe152715e4a7256cee62a3d5dd8484fd8a8539eb8beff.
//
// Solidity: event TSSAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) FilterTSSAddressUpdated(opts *bind.FilterOpts) (*MuseConnectorNonEthTSSAddressUpdatedIterator, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.FilterLogs(opts, "TSSAddressUpdated")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthTSSAddressUpdatedIterator{contract: _MuseConnectorNonEth.contract, event: "TSSAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchTSSAddressUpdated is a free log subscription operation binding the contract event 0xe79965b5c67dcfb2cf5fe152715e4a7256cee62a3d5dd8484fd8a8539eb8beff.
//
// Solidity: event TSSAddressUpdated(address callerAddress, address newTssAddress)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) WatchTSSAddressUpdated(opts *bind.WatchOpts, sink chan<- *MuseConnectorNonEthTSSAddressUpdated) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.WatchLogs(opts, "TSSAddressUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorNonEthTSSAddressUpdated)
				if err := _MuseConnectorNonEth.contract.UnpackLog(event, "TSSAddressUpdated", log); err != nil {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) ParseTSSAddressUpdated(log types.Log) (*MuseConnectorNonEthTSSAddressUpdated, error) {
	event := new(MuseConnectorNonEthTSSAddressUpdated)
	if err := _MuseConnectorNonEth.contract.UnpackLog(event, "TSSAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorNonEthTSSAddressUpdaterUpdatedIterator is returned from FilterTSSAddressUpdaterUpdated and is used to iterate over the raw logs and unpacked data for TSSAddressUpdaterUpdated events raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthTSSAddressUpdaterUpdatedIterator struct {
	Event *MuseConnectorNonEthTSSAddressUpdaterUpdated // Event containing the contract specifics and raw log

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
func (it *MuseConnectorNonEthTSSAddressUpdaterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorNonEthTSSAddressUpdaterUpdated)
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
		it.Event = new(MuseConnectorNonEthTSSAddressUpdaterUpdated)
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
func (it *MuseConnectorNonEthTSSAddressUpdaterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorNonEthTSSAddressUpdaterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorNonEthTSSAddressUpdaterUpdated represents a TSSAddressUpdaterUpdated event raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthTSSAddressUpdaterUpdated struct {
	CallerAddress        common.Address
	NewTssUpdaterAddress common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterTSSAddressUpdaterUpdated is a free log retrieval operation binding the contract event 0x5104c9abdc7d111c2aeb4ce890ac70274a4be2ee83f46a62551be5e6ebc82dd0.
//
// Solidity: event TSSAddressUpdaterUpdated(address callerAddress, address newTssUpdaterAddress)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) FilterTSSAddressUpdaterUpdated(opts *bind.FilterOpts) (*MuseConnectorNonEthTSSAddressUpdaterUpdatedIterator, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.FilterLogs(opts, "TSSAddressUpdaterUpdated")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthTSSAddressUpdaterUpdatedIterator{contract: _MuseConnectorNonEth.contract, event: "TSSAddressUpdaterUpdated", logs: logs, sub: sub}, nil
}

// WatchTSSAddressUpdaterUpdated is a free log subscription operation binding the contract event 0x5104c9abdc7d111c2aeb4ce890ac70274a4be2ee83f46a62551be5e6ebc82dd0.
//
// Solidity: event TSSAddressUpdaterUpdated(address callerAddress, address newTssUpdaterAddress)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) WatchTSSAddressUpdaterUpdated(opts *bind.WatchOpts, sink chan<- *MuseConnectorNonEthTSSAddressUpdaterUpdated) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.WatchLogs(opts, "TSSAddressUpdaterUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorNonEthTSSAddressUpdaterUpdated)
				if err := _MuseConnectorNonEth.contract.UnpackLog(event, "TSSAddressUpdaterUpdated", log); err != nil {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) ParseTSSAddressUpdaterUpdated(log types.Log) (*MuseConnectorNonEthTSSAddressUpdaterUpdated, error) {
	event := new(MuseConnectorNonEthTSSAddressUpdaterUpdated)
	if err := _MuseConnectorNonEth.contract.UnpackLog(event, "TSSAddressUpdaterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorNonEthUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthUnpausedIterator struct {
	Event *MuseConnectorNonEthUnpaused // Event containing the contract specifics and raw log

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
func (it *MuseConnectorNonEthUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorNonEthUnpaused)
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
		it.Event = new(MuseConnectorNonEthUnpaused)
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
func (it *MuseConnectorNonEthUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorNonEthUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorNonEthUnpaused represents a Unpaused event raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) FilterUnpaused(opts *bind.FilterOpts) (*MuseConnectorNonEthUnpausedIterator, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthUnpausedIterator{contract: _MuseConnectorNonEth.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *MuseConnectorNonEthUnpaused) (event.Subscription, error) {

	logs, sub, err := _MuseConnectorNonEth.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorNonEthUnpaused)
				if err := _MuseConnectorNonEth.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) ParseUnpaused(log types.Log) (*MuseConnectorNonEthUnpaused, error) {
	event := new(MuseConnectorNonEthUnpaused)
	if err := _MuseConnectorNonEth.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorNonEthMuseReceivedIterator is returned from FilterMuseReceived and is used to iterate over the raw logs and unpacked data for MuseReceived events raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthMuseReceivedIterator struct {
	Event *MuseConnectorNonEthMuseReceived // Event containing the contract specifics and raw log

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
func (it *MuseConnectorNonEthMuseReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorNonEthMuseReceived)
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
		it.Event = new(MuseConnectorNonEthMuseReceived)
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
func (it *MuseConnectorNonEthMuseReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorNonEthMuseReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorNonEthMuseReceived represents a MuseReceived event raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthMuseReceived struct {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) FilterMuseReceived(opts *bind.FilterOpts, sourceChainId []*big.Int, destinationAddress []common.Address, internalSendHash [][32]byte) (*MuseConnectorNonEthMuseReceivedIterator, error) {

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

	logs, sub, err := _MuseConnectorNonEth.contract.FilterLogs(opts, "MuseReceived", sourceChainIdRule, destinationAddressRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthMuseReceivedIterator{contract: _MuseConnectorNonEth.contract, event: "MuseReceived", logs: logs, sub: sub}, nil
}

// WatchMuseReceived is a free log subscription operation binding the contract event 0xf1302855733b40d8acb467ee990b6d56c05c80e28ebcabfa6e6f3f57cb50d698.
//
// Solidity: event MuseReceived(bytes museTxSenderAddress, uint256 indexed sourceChainId, address indexed destinationAddress, uint256 museValue, bytes message, bytes32 indexed internalSendHash)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) WatchMuseReceived(opts *bind.WatchOpts, sink chan<- *MuseConnectorNonEthMuseReceived, sourceChainId []*big.Int, destinationAddress []common.Address, internalSendHash [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _MuseConnectorNonEth.contract.WatchLogs(opts, "MuseReceived", sourceChainIdRule, destinationAddressRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorNonEthMuseReceived)
				if err := _MuseConnectorNonEth.contract.UnpackLog(event, "MuseReceived", log); err != nil {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) ParseMuseReceived(log types.Log) (*MuseConnectorNonEthMuseReceived, error) {
	event := new(MuseConnectorNonEthMuseReceived)
	if err := _MuseConnectorNonEth.contract.UnpackLog(event, "MuseReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorNonEthMuseRevertedIterator is returned from FilterMuseReverted and is used to iterate over the raw logs and unpacked data for MuseReverted events raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthMuseRevertedIterator struct {
	Event *MuseConnectorNonEthMuseReverted // Event containing the contract specifics and raw log

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
func (it *MuseConnectorNonEthMuseRevertedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorNonEthMuseReverted)
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
		it.Event = new(MuseConnectorNonEthMuseReverted)
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
func (it *MuseConnectorNonEthMuseRevertedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorNonEthMuseRevertedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorNonEthMuseReverted represents a MuseReverted event raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthMuseReverted struct {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) FilterMuseReverted(opts *bind.FilterOpts, destinationChainId []*big.Int, internalSendHash [][32]byte) (*MuseConnectorNonEthMuseRevertedIterator, error) {

	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	var internalSendHashRule []interface{}
	for _, internalSendHashItem := range internalSendHash {
		internalSendHashRule = append(internalSendHashRule, internalSendHashItem)
	}

	logs, sub, err := _MuseConnectorNonEth.contract.FilterLogs(opts, "MuseReverted", destinationChainIdRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthMuseRevertedIterator{contract: _MuseConnectorNonEth.contract, event: "MuseReverted", logs: logs, sub: sub}, nil
}

// WatchMuseReverted is a free log subscription operation binding the contract event 0x521fb0b407c2eb9b1375530e9b9a569889992140a688bc076aa72c1712012c88.
//
// Solidity: event MuseReverted(address museTxSenderAddress, uint256 sourceChainId, uint256 indexed destinationChainId, bytes destinationAddress, uint256 remainingMuseValue, bytes message, bytes32 indexed internalSendHash)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) WatchMuseReverted(opts *bind.WatchOpts, sink chan<- *MuseConnectorNonEthMuseReverted, destinationChainId []*big.Int, internalSendHash [][32]byte) (event.Subscription, error) {

	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	var internalSendHashRule []interface{}
	for _, internalSendHashItem := range internalSendHash {
		internalSendHashRule = append(internalSendHashRule, internalSendHashItem)
	}

	logs, sub, err := _MuseConnectorNonEth.contract.WatchLogs(opts, "MuseReverted", destinationChainIdRule, internalSendHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorNonEthMuseReverted)
				if err := _MuseConnectorNonEth.contract.UnpackLog(event, "MuseReverted", log); err != nil {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) ParseMuseReverted(log types.Log) (*MuseConnectorNonEthMuseReverted, error) {
	event := new(MuseConnectorNonEthMuseReverted)
	if err := _MuseConnectorNonEth.contract.UnpackLog(event, "MuseReverted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MuseConnectorNonEthMuseSentIterator is returned from FilterMuseSent and is used to iterate over the raw logs and unpacked data for MuseSent events raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthMuseSentIterator struct {
	Event *MuseConnectorNonEthMuseSent // Event containing the contract specifics and raw log

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
func (it *MuseConnectorNonEthMuseSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MuseConnectorNonEthMuseSent)
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
		it.Event = new(MuseConnectorNonEthMuseSent)
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
func (it *MuseConnectorNonEthMuseSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MuseConnectorNonEthMuseSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MuseConnectorNonEthMuseSent represents a MuseSent event raised by the MuseConnectorNonEth contract.
type MuseConnectorNonEthMuseSent struct {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) FilterMuseSent(opts *bind.FilterOpts, museTxSenderAddress []common.Address, destinationChainId []*big.Int) (*MuseConnectorNonEthMuseSentIterator, error) {

	var museTxSenderAddressRule []interface{}
	for _, museTxSenderAddressItem := range museTxSenderAddress {
		museTxSenderAddressRule = append(museTxSenderAddressRule, museTxSenderAddressItem)
	}
	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	logs, sub, err := _MuseConnectorNonEth.contract.FilterLogs(opts, "MuseSent", museTxSenderAddressRule, destinationChainIdRule)
	if err != nil {
		return nil, err
	}
	return &MuseConnectorNonEthMuseSentIterator{contract: _MuseConnectorNonEth.contract, event: "MuseSent", logs: logs, sub: sub}, nil
}

// WatchMuseSent is a free log subscription operation binding the contract event 0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4.
//
// Solidity: event MuseSent(address sourceTxOriginAddress, address indexed museTxSenderAddress, uint256 indexed destinationChainId, bytes destinationAddress, uint256 museValueAndGas, uint256 destinationGasLimit, bytes message, bytes museParams)
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) WatchMuseSent(opts *bind.WatchOpts, sink chan<- *MuseConnectorNonEthMuseSent, museTxSenderAddress []common.Address, destinationChainId []*big.Int) (event.Subscription, error) {

	var museTxSenderAddressRule []interface{}
	for _, museTxSenderAddressItem := range museTxSenderAddress {
		museTxSenderAddressRule = append(museTxSenderAddressRule, museTxSenderAddressItem)
	}
	var destinationChainIdRule []interface{}
	for _, destinationChainIdItem := range destinationChainId {
		destinationChainIdRule = append(destinationChainIdRule, destinationChainIdItem)
	}

	logs, sub, err := _MuseConnectorNonEth.contract.WatchLogs(opts, "MuseSent", museTxSenderAddressRule, destinationChainIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MuseConnectorNonEthMuseSent)
				if err := _MuseConnectorNonEth.contract.UnpackLog(event, "MuseSent", log); err != nil {
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
func (_MuseConnectorNonEth *MuseConnectorNonEthFilterer) ParseMuseSent(log types.Log) (*MuseConnectorNonEthMuseSent, error) {
	event := new(MuseConnectorNonEthMuseSent)
	if err := _MuseConnectorNonEth.contract.UnpackLog(event, "MuseSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
