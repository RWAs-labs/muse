package types

import (
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"

	coin "github.com/RWAs-labs/muse/pkg/coin"
)

// Validate checks that the RateLimiterFlags is valid
func (r RateLimiterFlags) Validate() error {
	// window must not be negative
	if r.Window < 0 {
		return fmt.Errorf("window must be positive: %d", r.Window)
	}

	seen := make(map[string]bool)
	for _, conversion := range r.Conversions {
		// check no duplicated conversion
		if _, ok := seen[conversion.Mrc20]; ok {
			return fmt.Errorf("duplicated conversion: %s", conversion.Mrc20)
		}
		seen[conversion.Mrc20] = true

		// check conversion is valid
		if conversion.Rate.IsNil() {
			return fmt.Errorf("rate is nil for conversion: %s", conversion.Mrc20)
		}

		// check address is valid
		if !ethcommon.IsHexAddress(conversion.Mrc20) {
			return fmt.Errorf("invalid mrc20 address (%s)", conversion.Mrc20)
		}
	}

	return nil
}

// GetConversionRate returns the conversion rate for the given mrc20
func (r RateLimiterFlags) GetConversionRate(mrc20 string) (sdkmath.LegacyDec, bool) {
	for _, conversion := range r.Conversions {
		if conversion.Mrc20 == mrc20 {
			return conversion.Rate, true
		}
	}
	return sdkmath.LegacyNewDec(0), false
}

// BuildAssetRateMapFromList builds maps (foreign chain id -> asset -> rate) from a list of gas and erc20 asset rates
// The 1st map: foreign chain id -> gas coin asset rate
// The 2nd map: foreign chain id -> erc20 asset -> erc20 coin asset rate
func BuildAssetRateMapFromList(assetRates []AssetRate) (map[int64]AssetRate, map[int64]map[string]AssetRate) {
	// the result maps
	gasAssetRateMap := make(map[int64]AssetRate)
	erc20AssetRateMap := make(map[int64]map[string]AssetRate)

	// loop through the asset rates to build the maps
	for _, assetRate := range assetRates {
		chainID := assetRate.ChainId
		if assetRate.CoinType == coin.CoinType_Gas {
			gasAssetRateMap[chainID] = assetRate
		} else {
			if _, found := erc20AssetRateMap[chainID]; !found {
				erc20AssetRateMap[chainID] = make(map[string]AssetRate)
			}
			erc20AssetRateMap[chainID][strings.ToLower(assetRate.Asset)] = assetRate
		}
	}
	return gasAssetRateMap, erc20AssetRateMap
}

// ConvertCctxValueToAmuse converts the value of the cctx to amuse using given conversion rates
func ConvertCctxValueToAmuse(
	chainID int64,
	cctx *CrossChainTx,
	gasAssetRateMap map[int64]AssetRate,
	erc20AssetRateMap map[int64]map[string]AssetRate,
) sdkmath.Int {
	var rate sdkmath.LegacyDec
	var decimals uint32
	switch cctx.InboundParams.CoinType {
	case coin.CoinType_Muse:
		// no conversion needed for MUSE
		return sdkmath.NewIntFromBigInt(cctx.GetCurrentOutboundParam().Amount.BigInt())
	case coin.CoinType_Gas:
		assetRate, found := gasAssetRateMap[chainID]
		if !found {
			// skip if no rate found for gas coin on this chainID
			return sdkmath.NewInt(0)
		}
		rate = assetRate.Rate
		decimals = assetRate.Decimals
	case coin.CoinType_ERC20:
		// get the ERC20 coin rate
		_, found := erc20AssetRateMap[chainID]
		if !found {
			// skip if no rate found for this chainID
			return sdkmath.NewInt(0)
		}
		assetRate := erc20AssetRateMap[chainID][strings.ToLower(cctx.InboundParams.Asset)]
		rate = assetRate.Rate
		decimals = assetRate.Decimals
	default:
		// skip CoinType_Cmd
		return sdkmath.NewInt(0)
	}
	// should not happen, return 0 to skip if it happens
	if rate.IsNil() || rate.LTE(sdkmath.LegacyNewDec(0)) {
		return sdkmath.NewInt(0)
	}

	// the whole coin amounts of muse and mrc20
	// given decimals = 6, the amount will be 10^6 = 1000000
	oneMuse := coin.AmusePerMuse()
	oneMrc20 := sdkmath.LegacyNewDec(10).Power(uint64(decimals))

	// convert cctx asset amount into amuse amount
	// given amountCctx = 2000000, rate = 0.8, decimals = 6
	// amountCctxDec: 2000000 * 0.8 = 1600000.0
	// amountAmuseDec: 1600000.0 * 10e18 / 10e6 = 1600000000000000000.0
	amountCctxDec := sdkmath.LegacyNewDecFromBigInt(cctx.GetCurrentOutboundParam().Amount.BigInt())
	amountAmuseDec := amountCctxDec.Mul(rate).Mul(oneMuse).Quo(oneMrc20)
	return amountAmuseDec.TruncateInt()
}

// RateLimitExceeded accumulates the cctx value and then checks if the rate limit is exceeded
// returns true if the rate limit is exceeded
func RateLimitExceeded(
	chainID int64,
	cctx *CrossChainTx,
	gasAssetRateMap map[int64]AssetRate,
	erc20AssetRateMap map[int64]map[string]AssetRate,
	currentCctxValue *sdkmath.Int,
	withdrawLimitInAmuse sdkmath.Int,
) bool {
	amountMuse := ConvertCctxValueToAmuse(chainID, cctx, gasAssetRateMap, erc20AssetRateMap)
	*currentCctxValue = currentCctxValue.Add(amountMuse)
	return currentCctxValue.GT(withdrawLimitInAmuse)
}
