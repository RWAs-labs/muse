package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/RWAs-labs/muse/cmd/musecored/config"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/x/emissions/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetReservesFactor(t *testing.T) {
	t.Run("successfully get reserves factor", func(t *testing.T) {
		//Arrange
		k, ctx, sdkK, _ := keepertest.EmissionsKeeper(t)
		amount := sdkmath.NewInt(100000000000000000)
		err := sdkK.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)))
		require.NoError(t, err)
		//Act
		reserveAmount := k.GetReservesFactor(ctx)
		//Assert
		require.Equal(t, amount.ToLegacyDec(), reserveAmount)
	})
}
