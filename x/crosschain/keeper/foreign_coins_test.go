package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
)

func TestKeeper_GetAllForeignCoins(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	fc := sample.ForeignCoins(t, sample.EthAddress().Hex())
	fc.ForeignChainId = 101
	k.GetFungibleKeeper().SetForeignCoins(ctx, fc)

	res := k.GetAllForeignCoins(ctx)
	require.Equal(t, 1, len(res))
}
