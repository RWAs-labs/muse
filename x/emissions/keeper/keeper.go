package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/RWAs-labs/muse/x/emissions/types"
)

type (
	Keeper struct {
		cdc              codec.Codec
		storeKey         storetypes.StoreKey
		memKey           storetypes.StoreKey
		feeCollectorName string
		bankKeeper       types.BankKeeper
		stakingKeeper    types.StakingKeeper
		observerKeeper   types.ObserverKeeper
		authKeeper       types.AccountKeeper
		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	feeCollectorName string,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	observerKeeper types.ObserverKeeper,
	authKeeper types.AccountKeeper,
	authority string,
) *Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(err)
	}

	return &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		memKey:           memKey,
		feeCollectorName: feeCollectorName,
		bankKeeper:       bankKeeper,
		stakingKeeper:    stakingKeeper,
		observerKeeper:   observerKeeper,
		authKeeper:       authKeeper,
		authority:        authority,
	}
}

func (k Keeper) GetCodec() codec.Codec {
	return k.cdc
}

func (k Keeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetFeeCollector() string {
	return k.feeCollectorName
}

func (k Keeper) GetBankKeeper() types.BankKeeper {
	return k.bankKeeper
}

func (k Keeper) GetStakingKeeper() types.StakingKeeper {
	return k.stakingKeeper
}

func (k Keeper) GetObserverKeeper() types.ObserverKeeper {
	return k.observerKeeper
}

func (k Keeper) GetAuthKeeper() types.AccountKeeper {
	return k.authKeeper
}

func (k Keeper) GetAuthority() string {
	return k.authority
}
