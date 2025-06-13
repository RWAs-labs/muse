package keeper

import (
	"github.com/RWAs-labs/muse/x/emissions/types"
)

var _ types.QueryServer = Keeper{}
