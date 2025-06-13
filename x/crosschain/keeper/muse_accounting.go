package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/RWAs-labs/muse/x/crosschain/types"
)

func (k Keeper) SetMuseAccounting(ctx sdk.Context, abortedMuseAmount types.MuseAccounting) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&abortedMuseAmount)
	store.Set([]byte(types.MuseAccountingKey), b)
}

func (k Keeper) GetMuseAccounting(ctx sdk.Context) (val types.MuseAccounting, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.MuseAccountingKey))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) AddMuseAbortedAmount(ctx sdk.Context, amount sdkmath.Uint) {
	museAccounting, found := k.GetMuseAccounting(ctx)
	if !found {
		museAccounting = types.MuseAccounting{
			AbortedMuseAmount: amount,
		}
	} else {
		museAccounting.AbortedMuseAmount = museAccounting.AbortedMuseAmount.Add(amount)
	}
	k.SetMuseAccounting(ctx, museAccounting)
}

func (k Keeper) RemoveMuseAbortedAmount(ctx sdk.Context, amount sdkmath.Uint) error {
	museAccounting, found := k.GetMuseAccounting(ctx)
	if !found {
		return types.ErrUnableToFindMuseAccounting
	}
	if museAccounting.AbortedMuseAmount.LT(amount) {
		return types.ErrInsufficientMuseAmount
	}
	museAccounting.AbortedMuseAmount = museAccounting.AbortedMuseAmount.Sub(amount)
	k.SetMuseAccounting(ctx, museAccounting)
	return nil
}
