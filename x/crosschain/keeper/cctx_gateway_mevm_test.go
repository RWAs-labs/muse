package keeper_test

import (
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/crosschain/keeper"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestKeeper_InitiateOutboundMEVM(t *testing.T) {
	t.Run("initiate outbound from gateway MEVM successfully", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		gatewayMEVM := keeper.NewCCTXGatewayMEVM(*k)

		// setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// mock up MUSEDepositAndCallContract
		fungibleMock.On("MUSEDepositAndCallContract", mock.Anything,
			mock.Anything,
			receiver, int64(1), amount, mock.Anything, mock.Anything).Return(nil, nil)

		// mock up CCTX data
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
		cctx.GetCurrentOutboundParam().Receiver = receiver.String()
		cctx.GetCurrentOutboundParam().ReceiverChainId = chains.MuseChainPrivnet.ChainId
		cctx.GetInboundParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.InboundParams.CoinType = coin.CoinType_Muse
		cctx.GetInboundParams().SenderChainId = 1

		// ACT
		// call InitiateOutbound
		newStatus, err := gatewayMEVM.InitiateOutbound(
			ctx,
			keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
		)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_OutboundMined, newStatus)
	})

	t.Run("should return aborted status on insufficient depositor fee", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		gatewayMEVM := keeper.NewCCTXGatewayMEVM(*k)

		// mock up CCTX data
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
		cctx.InboundParams.Status = types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE

		// ACT
		// call InitiateOutbound
		newStatus, err := gatewayMEVM.InitiateOutbound(
			ctx,
			keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
		)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_Aborted, newStatus)
	})

	t.Run("should return aborted status on unknown inbound status", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		gatewayMEVM := keeper.NewCCTXGatewayMEVM(*k)

		// mock up CCTX data
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
		cctx.InboundParams.Status = types.InboundStatus(1000)

		// ACT
		// call InitiateOutbound
		newStatus, err := gatewayMEVM.InitiateOutbound(
			ctx,
			keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
		)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_Aborted, newStatus)
	})

	t.Run("should return reverted status on invalid memo inbound status", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})
		gatewayMEVM := keeper.NewCCTXGatewayMEVM(*k)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     1,
			MedianIndex: 0,
			Prices:      []uint64{1},
		})

		//mock necessary calls made during creation of the revert outbound
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(sample.Chain(1), true)
		fungibleMock.On("GetGasCoinForForeignCoin", mock.Anything, mock.Anything).
			Return(fungibletypes.ForeignCoins{Mrc20ContractAddress: "0x"}, true)
		fungibleMock.On("QueryGasLimit", mock.Anything, mock.Anything).Return(big.NewInt(1), nil)
		fungibleMock.On("QuerySystemContractGasCoinMRC20", mock.Anything, mock.Anything).
			Return(ethcommon.Address{}, nil)
		fungibleMock.On("QuerySystemContractGasCoinMRC20", mock.Anything, mock.Anything).
			Return(ethcommon.Address{}, nil)
		fungibleMock.On("QueryProtocolFlatFee", mock.Anything, mock.Anything).Return(big.NewInt(1), nil)
		observerMock.On("GetChainNonces", mock.Anything, mock.Anything).
			Return(observertypes.ChainNonces{Nonce: 1}, true)
		observerMock.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, true)
		observerMock.On("GetPendingNonces", mock.Anything, mock.Anything, mock.Anything).
			Return(observertypes.PendingNonces{NonceHigh: int64(1)}, true)
		observerMock.On("SetChainNonces", mock.Anything, mock.Anything)
		observerMock.On("SetPendingNonces", mock.Anything, mock.Anything)

		// mock up CCTX data
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
		cctx.InboundParams.Status = types.InboundStatus_INVALID_MEMO
		cctx.InboundParams.SenderChainId = 1
		cctx.OutboundParams = []*types.OutboundParams{cctx.OutboundParams[0]}

		// ACT
		// call InitiateOutbound
		newStatus, err := gatewayMEVM.InitiateOutbound(
			ctx,
			keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
		)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, types.CctxStatus_PendingRevert, cctx.CctxStatus.Status)
		require.Equal(t, types.CctxStatus_PendingRevert, newStatus)
	})

	t.Run(
		"should return aborted status on 'error during deposit that is not smart contract revert'",
		func(t *testing.T) {
			// ARRANGE
			k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
				UseFungibleMock: true,
			})
			gatewayMEVM := keeper.NewCCTXGatewayMEVM(*k)

			// setup mock data
			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
			receiver := sample.EthAddress()
			amount := big.NewInt(42)

			// mock up MUSEDepositAndCallContract
			fungibleMock.On("MUSEDepositAndCallContract", mock.Anything,
				mock.Anything,
				receiver, int64(1), amount, mock.Anything, mock.Anything).Return(nil, sample.ErrSample)

			// mock up CCTX data
			cctx := sample.CrossChainTx(t, "test")
			cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingOutbound}
			cctx.GetCurrentOutboundParam().Receiver = receiver.String()
			cctx.GetCurrentOutboundParam().ReceiverChainId = chains.MuseChainPrivnet.ChainId
			cctx.GetInboundParams().Amount = sdkmath.NewUintFromBigInt(amount)
			cctx.InboundParams.CoinType = coin.CoinType_Muse
			cctx.GetInboundParams().SenderChainId = 1

			// ACT
			// call InitiateOutbound
			newStatus, err := gatewayMEVM.InitiateOutbound(
				ctx,
				keeper.InitiateOutboundConfig{CCTX: cctx, ShouldPayGas: true},
			)

			// ASSERT
			require.Error(t, err)
			require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
			require.Equal(t, types.CctxStatus_Aborted, newStatus)
			require.Contains(
				t,
				cctx.CctxStatus.StatusMessage,
				"outbound failed but the universal contract did not revert",
			)
			require.True(t, strings.Contains(cctx.CctxStatus.ErrorMessage, sample.ErrSample.Error()))
		},
	)
}
