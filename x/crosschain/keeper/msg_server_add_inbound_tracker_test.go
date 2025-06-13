package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	"github.com/RWAs-labs/muse/x/crosschain/keeper"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

func TestMsgServer_AddInboundTracker(t *testing.T) {
	t.Run("fail normal user submit", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		nonAdmin := sample.AccAddress()

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)

		msg := types.MsgAddInboundTracker{
			Creator:  nonAdmin,
			ChainId:  chainID,
			TxHash:   txHash,
			CoinType: coin.CoinType_Muse,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.False(t, found)
	})

	t.Run("fail for unsupported chain id", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, false)

		msg := types.MsgAddInboundTracker{
			Creator:  sample.AccAddress(),
			ChainId:  chainID + 1,
			TxHash:   txHash,
			CoinType: coin.CoinType_Muse,
		}
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.False(t, found)
	})

	t.Run("admin add tx tracker", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(false)

		setSupportedChain(ctx, zk, chainID)

		msg := types.MsgAddInboundTracker{
			Creator:  admin,
			ChainId:  chainID,
			TxHash:   txHash,
			CoinType: coin.CoinType_Muse,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := msgServer.AddInboundTracker(ctx, &msg)

		require.NoError(t, err)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.True(t, found)
	})

	t.Run("observer add tx tracker", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
			UseObserverMock:  true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		txHash := "string"
		chainID := getValidEthChainID()

		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(chains.Chain{}, true)
		observerMock.On("IsNonTombstonedObserver", mock.Anything, mock.Anything).Return(true)

		msg := types.MsgAddInboundTracker{
			Creator:  admin,
			ChainId:  chainID,
			TxHash:   txHash,
			CoinType: coin.CoinType_Muse,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.AddInboundTracker(ctx, &msg)
		require.NoError(t, err)
		_, found := k.GetInboundTracker(ctx, chainID, txHash)
		require.True(t, found)
	})
}
