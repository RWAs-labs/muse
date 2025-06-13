package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	lightclientmocks "github.com/RWAs-labs/muse/testutil/keeper/mocks/lightclient"
	"github.com/RWAs-labs/muse/x/lightclient/keeper"
	"github.com/RWAs-labs/muse/x/lightclient/types"
)

// LightclientMockOptions represents options for instantiating a lightclient keeper with mocks
type LightclientMockOptions struct {
	UseAuthorityMock bool
}

var (
	LightclientMocksAll = LightclientMockOptions{
		UseAuthorityMock: true,
	}
	LightclientNoMocks = LightclientMockOptions{}
)

func initLightclientKeeper(
	cdc codec.Codec,
	ss store.CommitMultiStore,
	authorityKeeper types.AuthorityKeeper,
) keeper.Keeper {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, nil)

	return keeper.NewKeeper(cdc, storeKey, memKey, authorityKeeper)
}

// LightclientKeeperWithMocks instantiates a lightclient keeper for testing purposes with the option to mock specific keepers
func LightclientKeeperWithMocks(
	t testing.TB,
	mockOptions LightclientMockOptions,
) (*keeper.Keeper, sdk.Context, SDKKeepers, MuseKeepers) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := rootmulti.NewStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	cdc := NewCodec()

	authorityKeeperTmp := initAuthorityKeeper(cdc, stateStore)

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create the observer keeper
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, nil)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := NewContext(stateStore)

	// Initialize modules genesis
	sdkKeepers.InitGenesis(ctx)

	// Add a proposer to the context
	ctx = sdkKeepers.InitBlockProposer(t, ctx)

	// Initialize mocks for mocked keepers
	var authorityKeeper types.AuthorityKeeper = authorityKeeperTmp
	if mockOptions.UseAuthorityMock {
		authorityKeeper = lightclientmocks.NewLightclientAuthorityKeeper(t)
	}

	k := keeper.NewKeeper(cdc, storeKey, memStoreKey, authorityKeeper)

	return &k, ctx, sdkKeepers, MuseKeepers{
		AuthorityKeeper: &authorityKeeperTmp,
	}
}

// LightclientKeeper instantiates an lightclient keeper for testing purposes
func LightclientKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, SDKKeepers, MuseKeepers) {
	return LightclientKeeperWithMocks(t, LightclientNoMocks)
}

// GetLightclientAuthorityMock returns a new lightclient authority keeper mock
func GetLightclientAuthorityMock(t testing.TB, keeper *keeper.Keeper) *lightclientmocks.LightclientAuthorityKeeper {
	cok, ok := keeper.GetAuthorityKeeper().(*lightclientmocks.LightclientAuthorityKeeper)
	require.True(t, ok)
	return cok
}
