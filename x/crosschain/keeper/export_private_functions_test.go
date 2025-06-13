package keeper

import (
	"github.com/RWAs-labs/muse/x/crosschain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// These functions are exported for testing purposes

func (k Keeper) UpdateMuseAccounting(ctx sdk.Context, cctx types.CrossChainTx) {
	k.updateMuseAccounting(ctx, cctx)
}

func (k Keeper) UpdateInboundHashToCCTX(ctx sdk.Context, cctx types.CrossChainTx) {
	k.updateInboundHashToCCTX(ctx, cctx)
}

func (k Keeper) SetNonceToCCTX(ctx sdk.Context, cctx types.CrossChainTx, tssPubkey string) {
	k.setNonceToCCTX(ctx, cctx, tssPubkey)
}

func (k Keeper) GetNextCctxCounter(ctx sdk.Context) uint64 {
	return k.getNextCctxCounter(ctx)
}
