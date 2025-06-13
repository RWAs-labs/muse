package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

func TestKeeper_ConvertGasToMuse(t *testing.T) {
	t.Run("should err if chain not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.ConvertGasToMuse(ctx, &types.QueryConvertGasToMuseRequest{
			ChainId: 987,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should err if median price not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		res, err := k.ConvertGasToMuse(ctx, &types.QueryConvertGasToMuseRequest{
			ChainId: 5,
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should err if mrc20 not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("QuerySystemContractGasCoinMRC20", mock.Anything, mock.Anything).
			Return(common.Address{}, errors.New("err"))

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     5,
			MedianIndex: 0,
			Prices:      []uint64{2},
		})

		res, err := k.ConvertGasToMuse(ctx, &types.QueryConvertGasToMuseRequest{
			ChainId:  5,
			GasLimit: "10",
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should err if uniswap2router not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("QuerySystemContractGasCoinMRC20", mock.Anything, mock.Anything).
			Return(sample.EthAddress(), nil)

		fungibleMock.On("QueryUniswapV2RouterGetMuseAmountsIn", mock.Anything, mock.Anything, mock.Anything).
			Return(big.NewInt(0), errors.New("err"))

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     5,
			MedianIndex: 0,
			Prices:      []uint64{2},
		})

		res, err := k.ConvertGasToMuse(ctx, &types.QueryConvertGasToMuseRequest{
			ChainId:  5,
			GasLimit: "10",
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if all is set", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		fungibleMock.On("QuerySystemContractGasCoinMRC20", mock.Anything, mock.Anything).
			Return(sample.EthAddress(), nil)

		fungibleMock.On("QueryUniswapV2RouterGetMuseAmountsIn", mock.Anything, mock.Anything, mock.Anything).
			Return(big.NewInt(5), nil)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     5,
			MedianIndex: 0,
			Prices:      []uint64{2},
		})

		res, err := k.ConvertGasToMuse(ctx, &types.QueryConvertGasToMuseRequest{
			ChainId:  5,
			GasLimit: "10",
		})
		require.NoError(t, err)
		require.Equal(t, &types.QueryConvertGasToMuseResponse{
			OutboundGasInMuse: "5",
			ProtocolFeeInMuse: types.GetProtocolFee().String(),
			// #nosec G115 always positive
			MuseBlockHeight: uint64(ctx.BlockHeight()),
		}, res)
	})
}

func TestKeeper_ProtocolFee(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	res, err := k.ProtocolFee(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, &types.QueryMessagePassingProtocolFeeResponse{
		FeeInMuse: types.GetProtocolFee().String(),
	}, res)
}
