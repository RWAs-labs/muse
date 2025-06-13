package keeper

import (
	"github.com/RWAs-labs/muse/x/crosschain/types"
)

var _ types.QueryServer = Keeper{}
