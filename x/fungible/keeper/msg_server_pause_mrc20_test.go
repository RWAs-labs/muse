package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/fungible/keeper"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestKeeper_PauseMRC20(t *testing.T) {
	t.Run("can pause mrc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		assertUnpaused := func(mrc20 string) {
			fc, found := k.GetForeignCoins(ctx, mrc20)
			require.True(t, found)
			require.False(t, fc.Paused)
		}
		assertPaused := func(mrc20 string) {
			fc, found := k.GetForeignCoins(ctx, mrc20)
			require.True(t, found)
			require.True(t, fc.Paused)
		}

		// setup mrc20
		mrc20A, mrc20B, mrc20C := sample.EthAddress().
			String(),
			sample.EthAddress().
				String(),
			sample.EthAddress().
				String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, mrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, mrc20B))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, mrc20C))
		assertUnpaused(mrc20A)
		assertUnpaused(mrc20B)
		assertUnpaused(mrc20C)

		// can pause mrc20
		msg := types.NewMsgPauseMRC20(
			admin,
			[]string{
				mrc20A,
				mrc20B,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.PauseMRC20(ctx, msg)
		require.NoError(t, err)
		assertPaused(mrc20A)
		assertPaused(mrc20B)
		assertUnpaused(mrc20C)

		// can pause already paused mrc20
		msg = types.NewMsgPauseMRC20(
			admin,
			[]string{
				mrc20B,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.PauseMRC20(ctx, msg)
		require.NoError(t, err)
		assertPaused(mrc20A)
		assertPaused(mrc20B)
		assertUnpaused(mrc20C)

		// can pause all mrc20
		msg = types.NewMsgPauseMRC20(
			admin,
			[]string{
				mrc20A,
				mrc20B,
				mrc20C,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.PauseMRC20(ctx, msg)
		require.NoError(t, err)
		assertPaused(mrc20A)
		assertPaused(mrc20B)
		assertPaused(mrc20C)
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgPauseMRC20(
			admin,
			[]string{sample.EthAddress().String()},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.PauseMRC20(ctx, msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if mrc20 does not exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		mrc20A, mrc20B := sample.EthAddress().String(), sample.EthAddress().String()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, mrc20A))
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, mrc20B))

		msg := types.NewMsgPauseMRC20(
			admin,
			[]string{
				mrc20A,
				sample.EthAddress().String(),
				mrc20B,
			},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.PauseMRC20(ctx, msg)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})
}
