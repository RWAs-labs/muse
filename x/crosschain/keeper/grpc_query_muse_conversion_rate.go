package keeper

import (
	"context"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/x/crosschain/types"
	museObserverTypes "github.com/RWAs-labs/muse/x/observer/types"
)

func (k Keeper) ConvertGasToMuse(
	context context.Context,
	request *types.QueryConvertGasToMuseRequest,
) (*types.QueryConvertGasToMuseResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)

	chain, found := chains.GetChainFromChainID(request.ChainId, k.GetAuthorityKeeper().GetAdditionalChainList(ctx))
	if !found {
		return nil, museObserverTypes.ErrSupportedChains
	}

	medianGasPrice, _, isFound := k.GetMedianGasValues(ctx, chain.ChainId)
	if !isFound {
		return nil, status.Error(codes.InvalidArgument, "invalid request: param chain")
	}

	gasLimit := math.NewUintFromString(request.GasLimit)
	outTxGasFee := medianGasPrice.Mul(gasLimit)
	mrc20, err := k.fungibleKeeper.QuerySystemContractGasCoinMRC20(ctx, big.NewInt(chain.ChainId))
	if err != nil {
		return nil, status.Error(codes.NotFound, "mrc20 not found")
	}

	outTxGasFeeInMuse, err := k.fungibleKeeper.QueryUniswapV2RouterGetMuseAmountsIn(ctx, outTxGasFee.BigInt(), mrc20)
	if err != nil {
		return nil, status.Error(codes.Internal, "zQueryUniswapv2RouterGetAmountsIn failed")
	}

	return &types.QueryConvertGasToMuseResponse{
		OutboundGasInMuse: outTxGasFeeInMuse.String(),
		ProtocolFeeInMuse: types.GetProtocolFee().String(),
		// #nosec G115 always positive
		MuseBlockHeight: uint64(ctx.BlockHeight()),
	}, nil
}

func (k Keeper) ProtocolFee(
	_ context.Context,
	_ *types.QueryMessagePassingProtocolFeeRequest,
) (*types.QueryMessagePassingProtocolFeeResponse, error) {
	return &types.QueryMessagePassingProtocolFeeResponse{
		FeeInMuse: types.GetProtocolFee().String(),
	}, nil
}
