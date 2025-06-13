package keeper

import (
	"github.com/RWAs-labs/muse/x/lightclient/types"
)

var _ types.QueryServer = Keeper{}
