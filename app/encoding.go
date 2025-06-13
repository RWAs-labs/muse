package app

import (
	evidencetypes "cosmossdk.io/x/evidence/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	evmenc "github.com/RWAs-labs/ethermint/encoding"
	ethermint "github.com/RWAs-labs/ethermint/types"
	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	feemarkettypes "github.com/RWAs-labs/ethermint/x/feemarket/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group"
	proposaltypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	emissionstypes "github.com/RWAs-labs/muse/x/emissions/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	lightclienttypes "github.com/RWAs-labs/muse/x/lightclient/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// MakeEncodingConfig creates an EncodingConfig
func MakeEncodingConfig() ethermint.EncodingConfig {
	encodingConfig := evmenc.MakeConfig()
	registry := encodingConfig.InterfaceRegistry
	// TODO test if we need to register these interfaces again as MakeConfig already registers them
	// https://github.com/RWAs-labs/muse/issues/3003
	cryptocodec.RegisterInterfaces(registry)
	authtypes.RegisterInterfaces(registry)
	authz.RegisterInterfaces(registry)
	banktypes.RegisterInterfaces(registry)
	stakingtypes.RegisterInterfaces(registry)
	slashingtypes.RegisterInterfaces(registry)
	upgradetypes.RegisterInterfaces(registry)
	distrtypes.RegisterInterfaces(registry)
	evidencetypes.RegisterInterfaces(registry)
	crisistypes.RegisterInterfaces(registry)
	evmtypes.RegisterInterfaces(registry)
	ethermint.RegisterInterfaces(registry)
	authoritytypes.RegisterInterfaces(registry)
	crosschaintypes.RegisterInterfaces(registry)
	emissionstypes.RegisterInterfaces(registry)
	fungibletypes.RegisterInterfaces(registry)
	observertypes.RegisterInterfaces(registry)
	lightclienttypes.RegisterInterfaces(registry)
	groupmodule.RegisterInterfaces(registry)
	govtypesv1beta1.RegisterInterfaces(registry)
	govtypesv1.RegisterInterfaces(registry)
	proposaltypes.RegisterInterfaces(registry)
	feemarkettypes.RegisterInterfaces(registry)
	consensusparamtypes.RegisterInterfaces(registry)
	vestingtypes.RegisterInterfaces(registry)

	return encodingConfig
}
