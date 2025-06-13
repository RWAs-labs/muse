package keeper

import (
	"github.com/RWAs-labs/muse/x/ibccrosschain/types"
)

var _ types.QueryServer = Keeper{}
