package types

import ethcommon "github.com/ethereum/go-ethereum/common"

// DefaultLiquidityCap is the default value set for the liquidity cap of a new MRC20 when deployed
// for security reason, this value is low. An arbitrary value should be set during the process of deploying a new MRC20
// The value is represented in the base unit of the MRC20, final value is calculated by multiplying this value by 10^decimals
const DefaultLiquidityCap = uint64(1000)

// MRC20Data represents the MRC4 token details used to map
// the token to a Cosmos Coin
type MRC20Data struct {
	Name     string
	Symbol   string
	Decimals uint8
}

// MRC20StringResponse defines the string value from the call response
type MRC20StringResponse struct {
	Value string
}

// MRC20Uint8Response defines the uint8 value from the call response
type MRC20Uint8Response struct {
	Value uint8
}

// MRC20BoolResponse defines the bool value from the call response
type MRC20BoolResponse struct {
	Value bool
}

// UniswapV2FactoryByte32Response defines the string value from the call response
type UniswapV2FactoryByte32Response struct {
	Value [32]byte
}

// SystemAddressResponse defines the address value from the call response
type SystemAddressResponse struct {
	Value ethcommon.Address
}

// NewMRC20Data creates a new MRC20Data instance
func NewMRC20Data(name, symbol string, decimals uint8) MRC20Data {
	return MRC20Data{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}
}
