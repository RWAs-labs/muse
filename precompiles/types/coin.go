package types

import (
	"math/big"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/RWAs-labs/muse/cmd/musecored/config"
)

// MRC20ToCosmosDenom returns the cosmos coin address for a given MRC20 address.
// This is converted to "mrc20/{MRC20Address}".
func MRC20ToCosmosDenom(MRC20Address common.Address) string {
	return config.MRC20DenomPrefix + MRC20Address.String()
}

func CreateMRC20CoinSet(mrc20address common.Address, amount *big.Int) (sdk.Coins, error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	if (mrc20address == common.Address{}) {
		return nil, &ErrInvalidAddr{
			Got:    mrc20address.String(),
			Reason: "empty address",
		}
	}

	denom := MRC20ToCosmosDenom(mrc20address)

	coin := sdk.NewCoin(denom, math.NewIntFromBigInt(amount))
	if !coin.IsValid() {
		return nil, &ErrInvalidCoin{
			Got:      coin.GetDenom(),
			Negative: coin.IsNegative(),
			Nil:      coin.IsNil(),
		}
	}

	// A sdk.Coins (type []sdk.Coin) has to be created because it's the type expected by MintCoins
	// and SendCoinsFromModuleToAccount.
	// But coinSet will only contain one coin, always.
	coinSet := sdk.NewCoins(coin)
	if !coinSet.IsValid() || coinSet.Empty() || coinSet.IsAnyNil() || coinSet == nil {
		return nil, &ErrInvalidCoin{
			Got:      coinSet.String(),
			Negative: coinSet.IsAnyNegative(),
			Nil:      coinSet.IsAnyNil(),
			Empty:    coinSet.Empty(),
		}
	}

	return coinSet, nil
}

// CoinIsMRC20 checks if a given coin is a MRC20 coin based on its denomination.
func CoinIsMRC20(denom string) bool {
	// Fail fast if the prefix is not set.
	if !strings.HasPrefix(denom, config.MRC20DenomPrefix) {
		return false
	}

	// Prefix is correctly set, extract the mrc20 address.
	mrc20Addr := strings.TrimPrefix(denom, config.MRC20DenomPrefix)

	// Return true only if address is not empty and is a valid hex address.
	return common.HexToAddress(mrc20Addr) != common.Address{} && common.IsHexAddress(mrc20Addr)
}
