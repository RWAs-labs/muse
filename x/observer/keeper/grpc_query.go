package keeper

import (
	"github.com/RWAs-labs/muse/x/observer/types"
)

var _ types.QueryServer = Keeper{}
