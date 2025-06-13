package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/crosschain/keeper"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestMsgServer_RemoveFromOutboundTracker(t *testing.T) {
	t.Run("should error if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId: 1,
			Nonce:   1,
		})

		admin := sample.AccAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)

		msg := types.MsgRemoveOutboundTracker{
			Creator: admin,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		res, err := msgServer.RemoveOutboundTracker(ctx, &msg)
		require.Error(t, err)
		require.Empty(t, res)

		_, found := k.GetOutboundTracker(ctx, 1, 1)
		require.True(t, found)
	})

	t.Run("should remove if authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			ChainId: 1,
			Nonce:   1,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)

		msg := types.MsgRemoveOutboundTracker{
			Creator: admin,
			ChainId: 1,
			Nonce:   1,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		res, err := msgServer.RemoveOutboundTracker(ctx, &msg)
		require.NoError(t, err)
		require.Empty(t, res)

		_, found := k.GetOutboundTracker(ctx, 1, 1)
		require.False(t, found)
	})
}
