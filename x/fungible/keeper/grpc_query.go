package keeper

import (
	"github.com/RWAs-labs/muse/x/fungible/types"
)

var _ types.QueryServer = Keeper{}
