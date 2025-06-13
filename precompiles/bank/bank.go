package bank

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	precompiletypes "github.com/RWAs-labs/muse/precompiles/types"
	fungiblekeeper "github.com/RWAs-labs/muse/x/fungible/keeper"
)

var (
	ABI                 abi.ABI
	ContractAddress     = common.HexToAddress("0x0000000000000000000000000000000000000067")
	GasRequiredByMethod = map[[4]byte]uint64{}
	ViewMethod          = map[[4]byte]bool{}
)

func init() {
	initABI()
}

func initABI() {
	if err := ABI.UnmarshalJSON([]byte(IBankMetaData.ABI)); err != nil {
		panic(err)
	}

	GasRequiredByMethod = map[[4]byte]uint64{}
	for methodName := range ABI.Methods {
		var methodID [4]byte
		copy(methodID[:], ABI.Methods[methodName].ID[:4])
		switch methodName {
		case DepositMethodName:
			GasRequiredByMethod[methodID] = DepositMethodGas
		case WithdrawMethodName:
			GasRequiredByMethod[methodID] = WithdrawMethodGas
		case BalanceOfMethodName:
			GasRequiredByMethod[methodID] = BalanceOfGas
		default:
			GasRequiredByMethod[methodID] = DefaultGas
		}
	}
}

type Contract struct {
	precompiletypes.BaseContract

	bankKeeper     bank.Keeper
	fungibleKeeper fungiblekeeper.Keeper
	cdc            codec.Codec
	kvGasConfig    storetypes.GasConfig
}

func NewIBankContract(
	ctx sdk.Context,
	bankKeeper bank.Keeper,
	fungibleKeeper fungiblekeeper.Keeper,
	cdc codec.Codec,
	kvGasConfig storetypes.GasConfig,
) *Contract {
	accAddress := sdk.AccAddress(ContractAddress.Bytes())
	if !fungibleKeeper.GetAuthKeeper().HasAccount(ctx, accAddress) {
		acc := fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ctx, accAddress)
		fungibleKeeper.GetAuthKeeper().SetAccount(ctx, acc)
	}

	return &Contract{
		BaseContract:   precompiletypes.NewBaseContract(ContractAddress),
		bankKeeper:     bankKeeper,
		fungibleKeeper: fungibleKeeper,
		cdc:            cdc,
		kvGasConfig:    kvGasConfig,
	}
}

// Address() is required to implement the PrecompiledContract interface.
func (c *Contract) Address() common.Address {
	return ContractAddress
}

// Abi() is required to implement the PrecompiledContract interface.
func (c *Contract) Abi() abi.ABI {
	return ABI
}

// RequiredGas is required to implement the PrecompiledContract interface.
// The gas has to be calculated deterministically based on the input.
func (c *Contract) RequiredGas(input []byte) uint64 {
	// get methodID (first 4 bytes)
	var methodID [4]byte
	copy(methodID[:], input[:4])
	// base cost to prevent large input size
	baseCost := uint64(len(input)) * c.kvGasConfig.WriteCostPerByte
	if ViewMethod[methodID] {
		baseCost = uint64(len(input)) * c.kvGasConfig.ReadCostPerByte
	}

	if requiredGas, ok := GasRequiredByMethod[methodID]; ok {
		return requiredGas + baseCost
	}

	// Can not happen, but return 0 if the method is not found.
	return 0
}

// Run is the entrypoint of the precompiled contract, it switches over the input method,
// and execute them accordingly.
func (c *Contract) Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) ([]byte, error) {
	method, err := ABI.MethodById(contract.Input[:4])
	if err != nil {
		return nil, err
	}

	args, err := method.Inputs.Unpack(contract.Input[4:])
	if err != nil {
		return nil, err
	}

	stateDB := evm.StateDB.(precompiletypes.ExtStateDB)

	switch method.Name {
	// Deposit and Withdraw methods are both not allowed in read-only mode.
	case DepositMethodName, WithdrawMethodName:
		if readOnly {
			return nil, precompiletypes.ErrWriteMethod{
				Method: method.Name,
			}
		}

		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			if method.Name == DepositMethodName {
				res, err = c.deposit(ctx, evm, contract, method, args)
			} else if method.Name == WithdrawMethodName {
				res, err = c.withdraw(ctx, evm, contract, method, args)
			}
			return err
		})
		if execErr != nil {
			res, errPack := method.Outputs.Pack(false)
			if errPack != nil {
				return nil, errPack
			}

			// Return the proper result (true/false) and the error message.
			// This way we make bank compliant with smart contracts which would expect a true/false.
			// And also with Go bindings which would expect an error.
			return res, err
		}
		return res, nil

	case BalanceOfMethodName:
		var res []byte
		execErr := stateDB.ExecuteNativeAction(contract.Address(), nil, func(ctx sdk.Context) error {
			res, err = c.balanceOf(ctx, method, args)
			return err
		})
		if execErr != nil {
			return nil, err
		}
		return res, nil

	default:
		return nil, precompiletypes.ErrInvalidMethod{
			Method: method.Name,
		}
	}
}
