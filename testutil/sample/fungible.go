package sample

import (
	"testing"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func ForeignCoins(t *testing.T, address string) types.ForeignCoins {
	r := newRandFromStringSeed(t, address)

	return types.ForeignCoins{
		Mrc20ContractAddress: address,
		Asset:                EthAddress().String(),
		ForeignChainId:       r.Int63(),
		Decimals:             uint32(r.Uint64()),
		Name:                 StringRandom(r, 32),
		Symbol:               StringRandom(r, 32),
		CoinType:             coin.CoinType_ERC20,
		GasLimit:             r.Uint64(),
		LiquidityCap:         UintInRange(0, 10000000000),
	}
}

func ForeignCoinList(t *testing.T, mrc20ETH, mrc20BTC, mrc20ERC20, erc20Asset string) []types.ForeignCoins {
	// eth and btc chain id
	ethChainID := chains.GoerliLocalnet.ChainId
	btcChainID := chains.BitcoinRegtest.ChainId

	// add mrc20 ETH
	fcGas := ForeignCoins(t, mrc20ETH)
	fcGas.Asset = ""
	fcGas.ForeignChainId = ethChainID
	fcGas.Decimals = 18
	fcGas.CoinType = coin.CoinType_Gas

	// add mrc20 BTC
	fcBTC := ForeignCoins(t, mrc20BTC)
	fcBTC.Asset = ""
	fcBTC.ForeignChainId = btcChainID
	fcBTC.Decimals = 8
	fcBTC.CoinType = coin.CoinType_Gas

	// add mrc20 ERC20
	fcERC20 := ForeignCoins(t, mrc20ERC20)
	fcERC20.Asset = erc20Asset
	fcERC20.ForeignChainId = ethChainID
	fcERC20.Decimals = 6
	fcERC20.CoinType = coin.CoinType_ERC20

	return []types.ForeignCoins{fcGas, fcBTC, fcERC20}
}

func SystemContract() *types.SystemContract {
	return &types.SystemContract{
		SystemContract: EthAddress().String(),
		ConnectorMevm:  EthAddress().String(),
	}
}
