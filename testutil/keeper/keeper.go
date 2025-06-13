package keeper

import (
	"math/rand"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	ethermint "github.com/RWAs-labs/ethermint/types"
	evmmodule "github.com/RWAs-labs/ethermint/x/evm"
	evmkeeper "github.com/RWAs-labs/ethermint/x/evm/keeper"
	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	feemarketkeeper "github.com/RWAs-labs/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/RWAs-labs/ethermint/x/feemarket/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitymodule "github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/testutil/sample"
	authoritymodule "github.com/RWAs-labs/muse/x/authority"
	authoritykeeper "github.com/RWAs-labs/muse/x/authority/keeper"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	crosschainmodule "github.com/RWAs-labs/muse/x/crosschain"
	crosschainkeeper "github.com/RWAs-labs/muse/x/crosschain/keeper"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	emissionsmodule "github.com/RWAs-labs/muse/x/emissions"
	emissionskeeper "github.com/RWAs-labs/muse/x/emissions/keeper"
	emissionstypes "github.com/RWAs-labs/muse/x/emissions/types"
	fungiblemodule "github.com/RWAs-labs/muse/x/fungible"
	fungiblekeeper "github.com/RWAs-labs/muse/x/fungible/keeper"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	ibccrosschainmodule "github.com/RWAs-labs/muse/x/ibccrosschain"
	ibccrosschainkeeper "github.com/RWAs-labs/muse/x/ibccrosschain/keeper"
	ibccrosschaintypes "github.com/RWAs-labs/muse/x/ibccrosschain/types"
	lightclientmodule "github.com/RWAs-labs/muse/x/lightclient"
	lightclientkeeper "github.com/RWAs-labs/muse/x/lightclient/keeper"
	lightclienttypes "github.com/RWAs-labs/muse/x/lightclient/types"
	observermodule "github.com/RWAs-labs/muse/x/observer"
	observerkeeper "github.com/RWAs-labs/muse/x/observer/keeper"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// NewContext creates a new sdk.Context for testing purposes with initialized header
func NewContext(stateStore storetypes.MultiStore) sdk.Context {
	header := tmproto.Header{
		Height:  1,
		ChainID: "test_7000-1",
		Time:    time.Now().UTC(),
		LastBlockId: tmproto.BlockID{
			Hash: tmhash.Sum([]byte("block_id")),
			PartSetHeader: tmproto.PartSetHeader{
				Total: 11,
				Hash:  tmhash.Sum([]byte("partset_header")),
			},
		},
		AppHash:            tmhash.Sum([]byte("app")),
		DataHash:           tmhash.Sum([]byte("data")),
		EvidenceHash:       tmhash.Sum([]byte("evidence")),
		ValidatorsHash:     tmhash.Sum([]byte("validators")),
		NextValidatorsHash: tmhash.Sum([]byte("next_validators")),
		ConsensusHash:      tmhash.Sum([]byte("consensus")),
		LastResultsHash:    tmhash.Sum([]byte("last_result")),
	}
	ctx := sdk.NewContext(stateStore, header, false, log.NewNopLogger())
	ctx = ctx.WithHeaderHash(tmhash.Sum([]byte("header")))
	return ctx
}

// SDKKeepers is a struct containing regular SDK module keepers for test purposes
type SDKKeepers struct {
	ParamsKeeper         paramskeeper.Keeper
	AuthKeeper           authkeeper.AccountKeeper
	BankKeeper           bankkeeper.Keeper
	StakingKeeper        stakingkeeper.Keeper
	SlashingKeeper       slashingkeeper.Keeper
	FeeMarketKeeper      feemarketkeeper.Keeper
	EvmKeeper            *evmkeeper.Keeper
	CapabilityKeeper     *capabilitykeeper.Keeper
	IBCKeeper            *ibckeeper.Keeper
	TransferKeeper       ibctransferkeeper.Keeper
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	DistributionKeeper   distrkeeper.Keeper
	EmissionsKeeper      *emissionskeeper.Keeper

	IBCRouter *porttypes.Router
}

// MuseKeepers is a struct containing Muse module keepers for test purposes
type MuseKeepers struct {
	AuthorityKeeper     *authoritykeeper.Keeper
	CrosschainKeeper    *crosschainkeeper.Keeper
	EmissionsKeeper     *emissionskeeper.Keeper
	FungibleKeeper      *fungiblekeeper.Keeper
	ObserverKeeper      *observerkeeper.Keeper
	LightclientKeeper   *lightclientkeeper.Keeper
	IBCCrosschainKeeper *ibccrosschainkeeper.Keeper
}

var moduleAccountPerms = map[string][]string{
	authtypes.FeeCollectorName:                      nil,
	distrtypes.ModuleName:                           nil,
	stakingtypes.BondedPoolName:                     {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:                  {authtypes.Burner, authtypes.Staking},
	evmtypes.ModuleName:                             {authtypes.Minter, authtypes.Burner},
	crosschaintypes.ModuleName:                      {authtypes.Minter, authtypes.Burner},
	fungibletypes.ModuleName:                        {authtypes.Minter, authtypes.Burner},
	emissionstypes.ModuleName:                       {authtypes.Minter},
	emissionstypes.UndistributedObserverRewardsPool: nil,
	emissionstypes.UndistributedTSSRewardsPool:      nil,
	ibctransfertypes.ModuleName:                     {authtypes.Minter, authtypes.Burner},
	ibccrosschaintypes.ModuleName:                   nil,
}

var (
	testStoreKeys = storetypes.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		evmtypes.StoreKey,
		consensustypes.StoreKey,
	)
	testTransientKeys = storetypes.NewTransientStoreKeys(evmtypes.TransientKey)
	testMemKeys       = storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		//ibctransfertypes.ModuleName:                     {authtypes.Minter, authtypes.Burner},
		crosschaintypes.ModuleName: {authtypes.Minter, authtypes.Burner},
		//ibccrosschaintypes.ModuleName:                   nil,
		evmtypes.ModuleName:                             {authtypes.Minter, authtypes.Burner},
		fungibletypes.ModuleName:                        {authtypes.Minter, authtypes.Burner},
		emissionstypes.ModuleName:                       nil,
		emissionstypes.UndistributedObserverRewardsPool: nil,
		emissionstypes.UndistributedTSSRewardsPool:      nil,
	}
)

// ModuleAccountAddrs returns all the app's module account addresses.
func ModuleAccountAddrs(maccPerms map[string][]string) map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// ParamsKeeper instantiates a param keeper for testing purposes
// TODO: remove https://github.com/RWAs-labs/muse/issues/848
func ParamsKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) paramskeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(paramstypes.StoreKey)
	tkeys := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)

	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(tkeys, storetypes.StoreTypeTransient, db)

	return paramskeeper.NewKeeper(
		cdc,
		fungibletypes.Amino,
		storeKey,
		tkeys,
	)
}

func ConsensusKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) consensuskeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(consensustypes.StoreKey)

	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	return consensuskeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		runtime.EventService{},
	)
}

// AccountKeeper instantiates an account keeper for testing purposes
func AccountKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) authkeeper.AccountKeeper {
	storeKey := storetypes.NewKVStoreKey(authtypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	return authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		ethermint.ProtoAccount,
		moduleAccountPerms,
		authcodec.NewBech32Codec("muse"),
		"muse",
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

// BankKeeper instantiates a bank keeper for testing purposes
func BankKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	authKeeper authkeeper.AccountKeeper,
) bankkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(banktypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	blockedAddrs := make(map[string]bool)

	return bankkeeper.NewBaseKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		authKeeper,
		blockedAddrs,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		log.NewNopLogger(),
	)
}

// StakingKeeper instantiates a staking keeper for testing purposes
func StakingKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
) stakingkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(stakingtypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	return *stakingkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		authKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		address.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		address.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
}

// SlashingKeeper instantiates a slashing keeper for testing purposes
func SlashingKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	stakingKeeper stakingkeeper.Keeper,
) slashingkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(slashingtypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	return slashingkeeper.NewKeeper(
		cdc,
		codec.NewLegacyAmino(),
		runtime.NewKVStoreService(storeKey),
		stakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

// DistributionKeeper instantiates a distribution keeper for testing purposes
func DistributionKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
) distrkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(distrtypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	return distrkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		authKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

// ProtocolVersionSetter mock
type ProtocolVersionSetter struct{}

func (vs ProtocolVersionSetter) SetProtocolVersion(uint64) {}

// UpgradeKeeper instantiates an upgrade keeper for testing purposes
func UpgradeKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) upgradekeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(upgradetypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	skipUpgradeHeights := make(map[int64]bool)
	vs := ProtocolVersionSetter{}

	return *upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(storeKey),
		cdc,
		"",
		vs,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

// FeeMarketKeeper instantiates a feemarket keeper for testing purposes
func FeeMarketKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) feemarketkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(feemarkettypes.StoreKey)
	transientKey := storetypes.NewTransientStoreKey(feemarkettypes.TransientKey)

	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(transientKey, storetypes.StoreTypeTransient, db)

	return feemarketkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		authtypes.NewModuleAddress(govtypes.ModuleName),
		storeKey,
		transientKey,
	)
}

// EVMKeeper instantiates an evm keeper for testing purposes
func EVMKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	feemarketKeeper feemarketkeeper.Keeper,
) *evmkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(evmtypes.StoreKey)
	transientKey := storetypes.NewTransientStoreKey(evmtypes.TransientKey)

	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(transientKey, storetypes.StoreTypeTransient, db)

	allKeys := make(map[string]storetypes.StoreKey, len(testStoreKeys)+len(testTransientKeys)+len(testMemKeys))
	for k, v := range testStoreKeys {
		allKeys[k] = v
	}
	for k, v := range testTransientKeys {
		allKeys[k] = v
	}
	for k, v := range testMemKeys {
		allKeys[k] = v
	}

	k := evmkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(storeKey),
		storeKey,
		transientKey,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		authKeeper,
		bankKeeper,
		stakingKeeper,
		feemarketKeeper,
		"",
		nil,
		allKeys,
	)

	return k
}

// NewSDKKeepersWithKeys instantiates regular Cosmos SDK keeper such as staking with local storage for testing purposes
func NewSDKKeepersWithKeys(
	cdc codec.Codec,
	keys map[string]*storetypes.KVStoreKey,
	memKeys map[string]*storetypes.MemoryStoreKey,
	tKeys map[string]*storetypes.TransientStoreKey,
	allKeys map[string]storetypes.StoreKey,
) SDKKeepers {
	authorityKeeper := authoritykeeper.NewKeeper(
		cdc,
		keys[authoritytypes.StoreKey],
		memKeys[authoritytypes.MemStoreKey],
		AuthorityGovAddress,
	)
	accountKeeper := authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		ethermint.ProtoAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(authtypes.ModuleName).String(),
	)
	paramsKeeper := paramskeeper.NewKeeper(
		cdc,
		fungibletypes.Amino,
		keys[paramstypes.StoreKey],
		tKeys[paramstypes.TStoreKey],
	)
	authKeeper := authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		ethermint.ProtoAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authtypes.NewModuleAddress(authtypes.ModuleName).String(),
	)
	blockedAddrs := make(map[string]bool)
	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		authKeeper,
		blockedAddrs,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		log.NewNopLogger(),
	)
	stakingKeeper := *stakingkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		authKeeper,
		bankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		address.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		address.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	feeMarketKeeper := feemarketkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keys[feemarkettypes.StoreKey]),
		authtypes.NewModuleAddress(govtypes.ModuleName),
		keys[feemarkettypes.StoreKey],
		tKeys[feemarkettypes.TransientKey],
	)
	evmKeeper := evmkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keys[evmtypes.StoreKey]),
		keys[evmtypes.StoreKey],
		tKeys[evmtypes.TransientKey],
		authtypes.NewModuleAddress(govtypes.ModuleName),
		authKeeper,
		bankKeeper,
		stakingKeeper,
		feeMarketKeeper,
		"",
		[]evmkeeper.CustomContractFn{},
		allKeys,
	)
	slashingKeeper := slashingkeeper.NewKeeper(
		cdc,
		codec.NewLegacyAmino(),
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		stakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	capabilityKeeper := capabilitykeeper.NewKeeper(
		cdc,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)
	dstrKeeper := distrkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(distrtypes.ModuleName).String(),
	)
	lightclientKeeper := lightclientkeeper.NewKeeper(
		cdc,
		keys[lightclienttypes.StoreKey],
		memKeys[lightclienttypes.MemStoreKey],
		authorityKeeper,
	)
	observerKeeper := observerkeeper.NewKeeper(
		cdc,
		keys[observertypes.StoreKey],
		memKeys[observertypes.MemStoreKey],
		stakingKeeper,
		slashingKeeper,
		authorityKeeper,
		lightclientKeeper,
		bankKeeper,
		authKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	emissionsKeeper := emissionskeeper.NewKeeper(
		cdc,
		keys[emissionstypes.StoreKey],
		memKeys[emissionstypes.MemStoreKey],
		authtypes.FeeCollectorName,
		bankKeeper,
		stakingKeeper,
		observerKeeper,
		authKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	return SDKKeepers{
		ParamsKeeper:       paramsKeeper,
		AuthKeeper:         authKeeper,
		BankKeeper:         bankKeeper,
		StakingKeeper:      stakingKeeper,
		FeeMarketKeeper:    feeMarketKeeper,
		EvmKeeper:          evmKeeper,
		SlashingKeeper:     slashingKeeper,
		CapabilityKeeper:   capabilityKeeper,
		DistributionKeeper: dstrKeeper,
		EmissionsKeeper:    emissionsKeeper,
	}
}

// CapabilityKeeper instantiates a capability keeper for testing purposes
func CapabilityKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) *capabilitykeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(capabilitytypes.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(capabilitytypes.MemStoreKey)

	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, nil)

	return capabilitykeeper.NewKeeper(cdc, storeKey, memKey)
}

// IBCKeeper instantiates an ibc keeper for testing purposes
func IBCKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	paramKeeper paramskeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	uppgradeKeeper upgradekeeper.Keeper,
	capabilityKeeper capabilitykeeper.Keeper,
) *ibckeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(ibcexported.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibcexported.ModuleName)

	return ibckeeper.NewKeeper(
		cdc,
		storeKey,
		paramKeeper.Subspace(ibcexported.ModuleName),
		stakingKeeper,
		uppgradeKeeper,
		scopedIBCKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
}

// TransferKeeper instantiates an ibc transfer keeper for testing purposes
func TransferKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	paramKeeper paramskeeper.Keeper,
	ibcKeeper *ibckeeper.Keeper,
	authKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper,
	capabilityKeeper capabilitykeeper.Keeper,
	ibcRouter *porttypes.Router,
) ibctransferkeeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(ibctransfertypes.StoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)

	scopedTransferKeeper := capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	transferKeeper := ibctransferkeeper.NewKeeper(
		cdc,
		storeKey,
		paramKeeper.Subspace(ibctransfertypes.ModuleName),
		ibcKeeper.ChannelKeeper,
		ibcKeeper.ChannelKeeper,
		ibcKeeper.PortKeeper,
		authKeeper,
		bankKeeper,
		scopedTransferKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// create IBC module from bottom to top of stack
	transferStack := transfer.NewIBCModule(transferKeeper)

	// Add transfer stack to IBC Router
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)

	return transferKeeper
}

// NewSDKKeepers instantiates regular Cosmos SDK keeper such as staking with local storage for testing purposes
func NewSDKKeepers(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) SDKKeepers {
	ibcRouter := porttypes.NewRouter()

	capabilityKeeper := CapabilityKeeper(cdc, db, ss)
	paramsKeeper := ParamsKeeper(cdc, db, ss)
	authKeeper := AccountKeeper(cdc, db, ss)
	bankKeeper := BankKeeper(cdc, db, ss, authKeeper)
	stakingKeeper := StakingKeeper(cdc, db, ss, authKeeper, bankKeeper)
	feeMarketKeeper := FeeMarketKeeper(cdc, db, ss)
	evmKeeper := EVMKeeper(
		cdc,
		db,
		ss,
		authKeeper,
		bankKeeper,
		stakingKeeper,
		feeMarketKeeper,
	)
	slashingKeeper := SlashingKeeper(cdc, db, ss, stakingKeeper)

	ibcKeeper := IBCKeeper(cdc, db, ss, paramsKeeper, stakingKeeper, UpgradeKeeper(cdc, db, ss), *capabilityKeeper)
	transferKeeper := TransferKeeper(
		cdc,
		db,
		ss,
		paramsKeeper,
		ibcKeeper,
		authKeeper,
		bankKeeper,
		*capabilityKeeper,
		ibcRouter,
	)

	return SDKKeepers{
		CapabilityKeeper: capabilityKeeper,
		ParamsKeeper:     paramsKeeper,
		AuthKeeper:       authKeeper,
		BankKeeper:       bankKeeper,
		StakingKeeper:    stakingKeeper,
		FeeMarketKeeper:  feeMarketKeeper,
		EvmKeeper:        evmKeeper,
		SlashingKeeper:   slashingKeeper,
		IBCKeeper:        ibcKeeper,
		TransferKeeper:   transferKeeper,
		IBCRouter:        ibcRouter,
	}
}

// InitGenesis initializes the test modules genesis state
func (sdkk SDKKeepers) InitGenesis(ctx sdk.Context) {
	capabilitymodule.InitGenesis(ctx, *sdkk.CapabilityKeeper, *capabilitytypes.DefaultGenesis())
	sdkk.AuthKeeper.InitGenesis(ctx, *authtypes.DefaultGenesisState())
	sdkk.BankKeeper.InitGenesis(ctx, banktypes.DefaultGenesisState())
	sdkk.StakingKeeper.InitGenesis(ctx, stakingtypes.DefaultGenesisState())
	evmGenesis := *evmtypes.DefaultGenesisState()
	evmGenesis.Params.EvmDenom = "amuse"
	evmmodule.InitGenesis(ctx, sdkk.EvmKeeper, sdkk.AuthKeeper, evmGenesis)
}

// InitBlockProposer initialize the block proposer for test purposes with an associated validator
func (sdkk SDKKeepers) InitBlockProposer(t testing.TB, ctx sdk.Context) sdk.Context {
	// #nosec G404 test purpose - weak randomness is not an issue here
	r := rand.New(rand.NewSource(42))

	// Set validator in the store
	validator := sample.Validator(t, r)
	err := sdkk.StakingKeeper.SetValidator(ctx, validator)
	require.NoError(t, err)
	err = sdkk.StakingKeeper.SetValidatorByConsAddr(ctx, validator)
	require.NoError(t, err)

	// Validator is proposer
	consAddr, err := validator.GetConsAddr()
	require.NoError(t, err)
	return ctx.WithProposer(consAddr)
}

// InitGenesis initializes the test modules genesis state for defined Muse modules
func (zk MuseKeepers) InitGenesis(ctx sdk.Context) {
	if zk.AuthorityKeeper != nil {
		authoritymodule.InitGenesis(ctx, *zk.AuthorityKeeper, *authoritytypes.DefaultGenesis())
	}
	if zk.CrosschainKeeper != nil {
		crosschainmodule.InitGenesis(ctx, *zk.CrosschainKeeper, *crosschaintypes.DefaultGenesis())
	}
	if zk.EmissionsKeeper != nil {
		emissionsmodule.InitGenesis(ctx, *zk.EmissionsKeeper, *emissionstypes.DefaultGenesis())
	}
	if zk.FungibleKeeper != nil {
		fungiblemodule.InitGenesis(ctx, *zk.FungibleKeeper, *fungibletypes.DefaultGenesis())
	}
	if zk.ObserverKeeper != nil {
		observermodule.InitGenesis(ctx, *zk.ObserverKeeper, *observertypes.DefaultGenesis())
	}
	if zk.LightclientKeeper != nil {
		lightclientmodule.InitGenesis(ctx, *zk.LightclientKeeper, *lightclienttypes.DefaultGenesis())
	}
	if zk.IBCCrosschainKeeper != nil {
		ibccrosschainmodule.InitGenesis(ctx, *zk.IBCCrosschainKeeper, *ibccrosschaintypes.DefaultGenesis())
	}
}
