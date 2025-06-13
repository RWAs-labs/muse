package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/RWAs-labs/muse/x/fungible/types"
)

type (
	Keeper struct {
		cdc             codec.Codec
		storeKey        storetypes.StoreKey
		memKey          storetypes.StoreKey
		authKeeper      types.AccountKeeper
		evmKeeper       types.EVMKeeper
		bankKeeper      types.BankKeeper
		observerKeeper  types.ObserverKeeper
		authorityKeeper types.AuthorityKeeper
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	authKeeper types.AccountKeeper,
	evmKeeper types.EVMKeeper,
	bankKeeper types.BankKeeper,
	observerKeeper types.ObserverKeeper,
	authorityKeeper types.AuthorityKeeper,
) *Keeper {
	return &Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		memKey:          memKey,
		authKeeper:      authKeeper,
		evmKeeper:       evmKeeper,
		bankKeeper:      bankKeeper,
		observerKeeper:  observerKeeper,
		authorityKeeper: authorityKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetAuthKeeper() types.AccountKeeper {
	return k.authKeeper
}

func (k Keeper) GetCodec() codec.Codec {
	return k.cdc
}

func (k Keeper) GetEVMKeeper() types.EVMKeeper {
	return k.evmKeeper
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetObserverKeeper() types.ObserverKeeper {
	return k.observerKeeper
}

func (k Keeper) GetAuthorityKeeper() types.AuthorityKeeper {
	return k.authorityKeeper
}
