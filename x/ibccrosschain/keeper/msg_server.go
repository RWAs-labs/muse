package keeper

import "github.com/RWAs-labs/muse/x/ibccrosschain/types"

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper} //nolint:typecheck
}

var _ types.MsgServer = msgServer{} //nolint:typecheck
