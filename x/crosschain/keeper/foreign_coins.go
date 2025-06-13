package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fungibleModuleTypes "github.com/RWAs-labs/muse/x/fungible/types"
)

func (k Keeper) GetAllForeignCoins(ctx sdk.Context) []fungibleModuleTypes.ForeignCoins {
	chains := k.museObserverKeeper.GetSupportedChains(ctx)
	var fCoins []fungibleModuleTypes.ForeignCoins
	for _, chain := range chains {
		fCoins = append(fCoins, k.fungibleKeeper.GetAllForeignCoinsForChain(ctx, chain.ChainId)...)
	}
	return fCoins
}
