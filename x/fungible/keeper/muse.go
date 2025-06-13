package keeper

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

// MUSEMaxSupplyStr is the maximum mintable MUSE in the fungible module
// 1.85 billion MUSE
const MUSEMaxSupplyStr = "1850000000000000000000000000"

// MintMuseToEVMAccount mints MUSE (gas token) to the given address
// NOTE: this method should be used with a temporary context, and it should not be committed if the method returns an error
func (k *Keeper) MintMuseToEVMAccount(ctx sdk.Context, to sdk.AccAddress, amount *big.Int) error {
	if err := k.validateMuseSupply(ctx, amount); err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewIntFromBigInt(amount)))
	// Mint coins
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	// Send minted coins to the receiver
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, coins)
}

func (k *Keeper) MintMuseToFungibleModule(ctx sdk.Context, amount *big.Int) error {
	if err := k.validateMuseSupply(ctx, amount); err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdkmath.NewIntFromBigInt(amount)))
	// Mint coins
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
}

// validateMuseSupply checks if the minted MUSE amount exceeds the maximum supply
func (k *Keeper) validateMuseSupply(ctx sdk.Context, amount *big.Int) error {
	museMaxSupply, ok := sdkmath.NewIntFromString(MUSEMaxSupplyStr)
	if !ok {
		return fmt.Errorf("failed to parse MUSE max supply: %s", MUSEMaxSupplyStr)
	}

	supply := k.bankKeeper.GetSupply(ctx, config.BaseDenom)
	if supply.Amount.Add(sdkmath.NewIntFromBigInt(amount)).GT(museMaxSupply) {
		return types.ErrMaxSupplyReached
	}
	return nil
}
