package keeper

import (
	"context"
	"sort"

	"github.com/btcsuite/btcd/chaincfg"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/crypto"
	"github.com/RWAs-labs/muse/x/observer/types"
)

// TSS returns the tss address for the current tss only
func (k Keeper) TSS(c context.Context, req *types.QueryGetTSSRequest) (*types.QueryGetTSSResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTSS(ctx)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetTSSResponse{TSS: val}, nil
}

// TssHistory Query historical list of TSS information
func (k Keeper) TssHistory(c context.Context, _ *types.QueryTssHistoryRequest) (*types.QueryTssHistoryResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	tssList := k.GetAllTSS(ctx)
	sort.SliceStable(tssList, func(i, j int) bool {
		return tssList[i].FinalizedMuseHeight < tssList[j].FinalizedMuseHeight
	})
	return &types.QueryTssHistoryResponse{TssList: tssList}, nil
}

func (k Keeper) GetTssAddress(
	goCtx context.Context,
	req *types.QueryGetTssAddressRequest,
) (*types.QueryGetTssAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "current tss not set")
	}

	ethAddress, err := crypto.GetTSSAddrEVM(tss.TssPubkey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	bitcoinParams := &chaincfg.RegressionNetParams
	if req.BitcoinChainId != 0 {
		bitcoinParams, err = chains.BitcoinNetParamsFromChainID(req.BitcoinChainId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	btcAddress, err := crypto.GetTSSAddrBTC(tss.TssPubkey, bitcoinParams)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	suiAddress, err := crypto.GetTSSAddrSui(tss.TssPubkey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetTssAddressResponse{
		Eth: ethAddress.String(),
		Btc: btcAddress,
		Sui: suiAddress,
	}, nil
}

func (k Keeper) GetTssAddressByFinalizedHeight(
	goCtx context.Context,
	req *types.QueryGetTssAddressByFinalizedHeightRequest,
) (*types.QueryGetTssAddressByFinalizedHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	tss, found := k.GetHistoricalTssByFinalizedHeight(ctx, req.FinalizedMuseHeight)
	if !found {
		return nil, status.Error(codes.NotFound, "tss not found")
	}

	ethAddress, err := crypto.GetTSSAddrEVM(tss.TssPubkey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	bitcoinParams := &chaincfg.RegressionNetParams
	if req.BitcoinChainId != 0 {
		bitcoinParams, err = chains.BitcoinNetParamsFromChainID(req.BitcoinChainId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	btcAddress, err := crypto.GetTSSAddrBTC(tss.TssPubkey, bitcoinParams)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	suiAddress, err := crypto.GetTSSAddrSui(tss.TssPubkey)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetTssAddressByFinalizedHeightResponse{
		Eth: ethAddress.String(),
		Btc: btcAddress,
		Sui: suiAddress,
	}, nil
}
