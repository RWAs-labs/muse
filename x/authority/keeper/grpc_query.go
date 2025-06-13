package keeper

import (
	"github.com/RWAs-labs/muse/x/authority/types"
)

var _ types.QueryServer = Keeper{}
