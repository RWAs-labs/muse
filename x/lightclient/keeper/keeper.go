package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/RWAs-labs/muse/x/lightclient/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc             codec.Codec
	storeKey        storetypes.StoreKey
	memKey          storetypes.StoreKey
	authorityKeeper types.AuthorityKeeper
}

// NewKeeper creates new instances of the lightclient Keeper
func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	authorityKeeper types.AuthorityKeeper,
) Keeper {
	return Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		memKey:          memKey,
		authorityKeeper: authorityKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetStoreKey returns the key to the store for lightclient
func (k Keeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

// GetMemKey returns the mem key to the store for lightclient
func (k Keeper) GetMemKey() storetypes.StoreKey {
	return k.memKey
}

// GetCodec returns the codec for lightclient
func (k Keeper) GetCodec() codec.Codec {
	return k.cdc
}

// GetAuthorityKeeper returns the authority keeper
func (k Keeper) GetAuthorityKeeper() types.AuthorityKeeper {
	return k.authorityKeeper
}
