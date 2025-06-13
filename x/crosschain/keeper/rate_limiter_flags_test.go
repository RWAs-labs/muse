package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
)

// createForeignCoinAndAssetRate creates foreign coin and corresponding asset rate
func createForeignCoinAndAssetRate(
	t *testing.T,
	mrc20Addr string,
	asset string,
	chainID int64,
	decimals uint32,
	coinType coin.CoinType,
	rate sdkmath.LegacyDec,
) (fungibletypes.ForeignCoins, types.AssetRate) {
	// create foreign coin
	foreignCoin := sample.ForeignCoins(t, mrc20Addr)
	foreignCoin.Asset = asset
	foreignCoin.ForeignChainId = chainID
	foreignCoin.Decimals = decimals
	foreignCoin.CoinType = coinType

	// create corresponding asset rate
	assetRate := sample.CustomAssetRate(
		foreignCoin.ForeignChainId,
		foreignCoin.Asset,
		foreignCoin.Decimals,
		foreignCoin.CoinType,
		rate,
	)

	return foreignCoin, assetRate
}

func TestKeeper_GetRateLimiterFlags(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)

	// not found
	_, found := k.GetRateLimiterFlags(ctx)
	require.False(t, found)

	flags := sample.RateLimiterFlags()

	k.SetRateLimiterFlags(ctx, flags)
	r, found := k.GetRateLimiterFlags(ctx)
	require.True(t, found)
	require.Equal(t, flags, r)
}

func TestKeeper_GetRateLimiterAssetRateList(t *testing.T) {
	k, ctx, _, zk := keepertest.CrosschainKeeper(t)

	// create test flags
	chainID := chains.GoerliLocalnet.ChainId
	mrc20GasAddr := sample.EthAddress().Hex()
	mrc20ERC20Addr1 := sample.EthAddress().Hex()
	mrc20ERC20Addr2 := sample.EthAddress().Hex()
	testflags := types.RateLimiterFlags{
		Rate: sdkmath.NewUint(100),
		Conversions: []types.Conversion{
			{
				Mrc20: mrc20GasAddr,
				Rate:  sdkmath.LegacyNewDec(1),
			},
			{
				Mrc20: mrc20ERC20Addr1,
				Rate:  sdkmath.LegacyNewDec(2),
			},
			{
				Mrc20: mrc20ERC20Addr2,
				Rate:  sdkmath.LegacyNewDec(3),
			},
		},
	}

	// asset rates not found before setting flags
	flags, assetRates, found := k.GetRateLimiterAssetRateList(ctx)
	require.False(t, found)
	require.Equal(t, types.RateLimiterFlags{}, flags)
	require.Nil(t, assetRates)

	// set flags
	k.SetRateLimiterFlags(ctx, testflags)

	// add gas coin
	gasCoin, gasAssetRate := createForeignCoinAndAssetRate(
		t,
		mrc20GasAddr,
		"",
		chainID,
		18,
		coin.CoinType_Gas,
		sdkmath.LegacyNewDec(1),
	)
	zk.FungibleKeeper.SetForeignCoins(ctx, gasCoin)

	// add 1st erc20 coin
	erc20Coin1, erc20AssetRate1 := createForeignCoinAndAssetRate(
		t,
		mrc20ERC20Addr1,
		sample.EthAddress().Hex(),
		chainID,
		8,
		coin.CoinType_ERC20,
		sdkmath.LegacyNewDec(2),
	)
	zk.FungibleKeeper.SetForeignCoins(ctx, erc20Coin1)

	// add 2nd erc20 coin
	erc20Coin2, erc20AssetRate2 := createForeignCoinAndAssetRate(
		t,
		mrc20ERC20Addr2,
		sample.EthAddress().Hex(),
		chainID,
		6,
		coin.CoinType_ERC20,
		sdkmath.LegacyNewDec(3),
	)
	zk.FungibleKeeper.SetForeignCoins(ctx, erc20Coin2)

	// get rates
	flags, assetRates, found = k.GetRateLimiterAssetRateList(ctx)
	require.True(t, found)
	require.Equal(t, testflags, flags)
	require.EqualValues(t, []types.AssetRate{gasAssetRate, erc20AssetRate1, erc20AssetRate2}, assetRates)
}
