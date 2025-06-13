package precompiles

import (
	storetypes "cosmossdk.io/store/types"
	evmkeeper "github.com/RWAs-labs/ethermint/x/evm/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	ethparams "github.com/ethereum/go-ethereum/params"

	"github.com/RWAs-labs/muse/precompiles/bank"
	"github.com/RWAs-labs/muse/precompiles/prototype"
	"github.com/RWAs-labs/muse/precompiles/staking"
	fungiblekeeper "github.com/RWAs-labs/muse/x/fungible/keeper"
)

// EnabledStatefulContracts contains the list of all enabled stateful precompiles.
// This is useful for listing and reading from other packages, such as BlockedAddrs() function.
// Setting to false a contract here will disable it, not being included in the blockchain.
var EnabledStatefulContracts = map[common.Address]bool{
	prototype.ContractAddress: false,
	staking.ContractAddress:   false,
	bank.ContractAddress:      false,
}

// StatefulContracts returns all the registered precompiled contracts.
func StatefulContracts(
	fungibleKeeper *fungiblekeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	distributionKeeper distrkeeper.Keeper,
	cdc codec.Codec,
	gasConfig storetypes.GasConfig,
) (precompiledContracts []evmkeeper.CustomContractFn) {
	// Initialize at 0 the custom compiled contracts and the addresses.
	precompiledContracts = make([]evmkeeper.CustomContractFn, 0)

	// Define the prototype contract function.
	if EnabledStatefulContracts[prototype.ContractAddress] {
		prototypeContract := func(_ sdktypes.Context, _ ethparams.Rules) vm.StatefulPrecompiledContract {
			return prototype.NewIPrototypeContract(fungibleKeeper, cdc, gasConfig)
		}

		// Append the prototype contract to the precompiledContracts slice.
		precompiledContracts = append(precompiledContracts, prototypeContract)
	}

	// Define the staking contract function.
	if EnabledStatefulContracts[staking.ContractAddress] {
		stakingContract := func(ctx sdktypes.Context, _ ethparams.Rules) vm.StatefulPrecompiledContract {
			return staking.NewIStakingContract(
				ctx,
				stakingKeeper,
				*fungibleKeeper,
				bankKeeper,
				distributionKeeper,
				cdc,
				gasConfig,
			)
		}

		// Append the staking contract to the precompiledContracts slice.
		precompiledContracts = append(precompiledContracts, stakingContract)
	}

	if EnabledStatefulContracts[bank.ContractAddress] {
		bankContract := func(ctx sdktypes.Context, _ ethparams.Rules) vm.StatefulPrecompiledContract {
			return bank.NewIBankContract(ctx, bankKeeper, *fungibleKeeper, cdc, gasConfig)
		}

		// Append the staking contract to the precompiledContracts slice.
		precompiledContracts = append(precompiledContracts, bankContract)
	}

	return precompiledContracts
}
