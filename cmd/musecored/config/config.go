package config

import (
	sdkmath "cosmossdk.io/math"
	ethermint "github.com/RWAs-labs/ethermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DisplayDenom = "muse"
	BaseDenom    = "amuse"
	AppName      = "musecored"
)

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(DisplayDenom, sdkmath.LegacyOneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(BaseDenom, sdkmath.LegacyNewDecWithPrec(1, ethermint.BaseDenomUnit)); err != nil {
		panic(err)
	}
}
