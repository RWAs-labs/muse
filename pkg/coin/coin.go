package coin

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
)

func AmusePerMuse() sdkmath.LegacyDec {
	return sdkmath.LegacyNewDec(1e18)
}

func GetCoinType(coin string) (CoinType, error) {
	coinInt, err := strconv.ParseInt(coin, 10, 32)
	if err != nil {
		return CoinType_Cmd, err
	}

	// check boundaries of the enum
	if coinInt < 0 || coinInt > int64(len(CoinType_name)) {
		return CoinType_Cmd, fmt.Errorf("invalid coin type %d", coinInt)
	}

	// #nosec G115 always in range
	return CoinType(coinInt), nil
}

func GetAmuseDecFromAmountInMuse(museAmount string) (sdkmath.LegacyDec, error) {
	museDec, err := sdkmath.LegacyNewDecFromStr(museAmount)
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	museToAmuseConvertionFactor := sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(1000000000000000000))
	return museDec.Mul(museToAmuseConvertionFactor), nil
}

func (c CoinType) SupportsRefund() bool {
	return c == CoinType_ERC20 || c == CoinType_Gas || c == CoinType_Muse
}

// IsAsset returns true if the coin type represents transport of asset.
// CoinType_Cmd and CoinType_NoAssetCall are not transport of asset.
func (c CoinType) IsAsset() bool {
	return c == CoinType_ERC20 || c == CoinType_Gas || c == CoinType_Muse
}
