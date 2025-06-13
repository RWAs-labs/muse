package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	"github.com/stretchr/testify/require"

	ibccrosschainmocks "github.com/RWAs-labs/muse/testutil/keeper/mocks/ibccrosschain"
	"github.com/RWAs-labs/muse/x/ibccrosschain/keeper"
	"github.com/RWAs-labs/muse/x/ibccrosschain/types"
)

type IBCCroscchainMockOptions struct {
	UseCrosschainMock  bool
	UseIBCTransferMock bool
}

var (
	IBCCrosschainMocksAll = IBCCroscchainMockOptions{
		UseCrosschainMock:  true,
		UseIBCTransferMock: true,
	}
	IBCCrosschainNoMocks = IBCCroscchainMockOptions{}
)

func initIBCCrosschainKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	crosschainKeeper types.CrosschainKeeper,
	ibcTransferKeeper types.IBCTransferKeeper,
	capabilityKeeper capabilitykeeper.Keeper,
) *keeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, db)

	capabilityKeeper.ScopeToModule(types.ModuleName)

	return keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		crosschainKeeper,
		ibcTransferKeeper,
	)
}

func IBCCrosschainKeeperWithMocks(
	t testing.TB,
	mockOptions IBCCroscchainMockOptions,
) (*keeper.Keeper, sdk.Context, SDKKeepers, MuseKeepers) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	cdc := NewCodec()

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create muse keepers
	authorityKeeper := initAuthorityKeeper(cdc, stateStore)
	lightclientKeeper := initLightclientKeeper(cdc, stateStore, authorityKeeper)
	observerKeeper := initObserverKeeper(
		cdc,
		stateStore,
		sdkKeepers.StakingKeeper,
		sdkKeepers.SlashingKeeper,
		authorityKeeper,
		sdkKeepers.BankKeeper,
		sdkKeepers.AuthKeeper,
		lightclientKeeper,
	)
	fungibleKeeper := initFungibleKeeper(
		cdc,
		stateStore,
		sdkKeepers.AuthKeeper,
		sdkKeepers.BankKeeper,
		sdkKeepers.EvmKeeper,
		observerKeeper,
		authorityKeeper,
	)
	crosschainKeeperTmp := initCrosschainKeeper(
		cdc,
		db,
		stateStore,
		sdkKeepers.StakingKeeper,
		sdkKeepers.AuthKeeper,
		sdkKeepers.BankKeeper,
		observerKeeper,
		fungibleKeeper,
		authorityKeeper,
		lightclientKeeper,
	)

	museKeepers := MuseKeepers{
		ObserverKeeper:    observerKeeper,
		FungibleKeeper:    fungibleKeeper,
		AuthorityKeeper:   &authorityKeeper,
		LightclientKeeper: &lightclientKeeper,
		CrosschainKeeper:  crosschainKeeperTmp,
	}

	var crosschainKeeper types.CrosschainKeeper = crosschainKeeperTmp
	var ibcTransferKeeper types.IBCTransferKeeper = sdkKeepers.TransferKeeper

	// Create the ibccrosschain keeper
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := NewContext(stateStore)

	// Initialize modules genesis
	sdkKeepers.InitGenesis(ctx)
	museKeepers.InitGenesis(ctx)

	// Add a proposer to the context
	ctx = sdkKeepers.InitBlockProposer(t, ctx)

	// Initialize mocks for mocked keepers
	if mockOptions.UseCrosschainMock {
		crosschainKeeper = ibccrosschainmocks.NewLightclientCrosschainKeeper(t)
	}
	if mockOptions.UseIBCTransferMock {
		ibcTransferKeeper = ibccrosschainmocks.NewLightclientTransferKeeper(t)
	}

	sdkKeepers.CapabilityKeeper.ScopeToModule(types.ModuleName)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		crosschainKeeper,
		ibcTransferKeeper,
	)

	// seal the IBC router
	sdkKeepers.IBCKeeper.SetRouter(sdkKeepers.IBCRouter)

	return k, ctx, sdkKeepers, museKeepers
}

// IBCCrosschainKeeperAllMocks creates a new ibccrosschain keeper with all mocked keepers
func IBCCrosschainKeeperAllMocks(t testing.TB) (*keeper.Keeper, sdk.Context) {
	k, ctx, _, _ := IBCCrosschainKeeperWithMocks(t, IBCCrosschainMocksAll)
	return k, ctx
}

// IBCCrosschainKeeper creates a new ibccrosschain keeper with no mocked keepers
func IBCCrosschainKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, SDKKeepers, MuseKeepers) {
	return IBCCrosschainKeeperWithMocks(t, IBCCrosschainNoMocks)
}
