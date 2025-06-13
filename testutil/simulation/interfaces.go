package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/RWAs-labs/muse/pkg/chains"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

type ObserverKeeper interface {
	GetObserverSet(ctx sdk.Context) (val observertypes.ObserverSet, found bool)
	IsNonTombstonedObserver(ctx sdk.Context, address string) bool
	GetSupportedChains(ctx sdk.Context) []chains.Chain
	GetNodeAccount(ctx sdk.Context, address string) (observertypes.NodeAccount, bool)
	GetAllNodeAccount(ctx sdk.Context) []observertypes.NodeAccount
}

type AuthorityKeeper interface {
	CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error
	GetAdditionalChainList(ctx sdk.Context) (list []chains.Chain)
	GetPolicies(ctx sdk.Context) (val authoritytypes.Policies, found bool)
}

type FungibleKeeper interface {
	GetForeignCoins(ctx sdk.Context, mrc20Addr string) (val fungibletypes.ForeignCoins, found bool)
	GetAllForeignCoins(ctx sdk.Context) (list []fungibletypes.ForeignCoins)
}
