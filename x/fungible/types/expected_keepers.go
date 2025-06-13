package types

import (
	"context"
	"math/big"

	"github.com/RWAs-labs/ethermint/x/evm/statedb"
	evmtypes "github.com/RWAs-labs/ethermint/x/evm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/RWAs-labs/muse/pkg/chains"
	authoritytypes "github.com/RWAs-labs/muse/x/authority/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	GetSequence(ctx context.Context, addr sdk.AccAddress) (uint64, error)
	GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI
	HasAccount(ctx context.Context, addr sdk.AccAddress) bool
	NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
}

type BankKeeper interface {
	SendCoinsFromModuleToAccount(
		ctx context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx context.Context, denom string) sdk.Coin
}

type ObserverKeeper interface {
	GetSupportedChains(ctx sdk.Context) []chains.Chain
}

type EVMKeeper interface {
	ChainID() *big.Int
	GetBlockBloomTransient(ctx sdk.Context) *big.Int
	GetLogSizeTransient(ctx sdk.Context) uint64
	WithChainID(ctx sdk.Context)
	SetBlockBloomTransient(ctx sdk.Context, bloom *big.Int)
	SetLogSizeTransient(ctx sdk.Context, logSize uint64)
	EstimateGas(c context.Context, req *evmtypes.EthCallRequest) (*evmtypes.EstimateGasResponse, error)
	ApplyMessage(
		ctx sdk.Context,
		msg *core.Message,
		tracer vm.EVMLogger,
		commit bool,
	) (*evmtypes.MsgEthereumTxResponse, error)
	GetAccount(ctx sdk.Context, addr ethcommon.Address) *statedb.Account
	GetCode(ctx sdk.Context, codeHash ethcommon.Hash) []byte
	SetAccount(ctx sdk.Context, addr ethcommon.Address, account statedb.Account) error
}

type AuthorityKeeper interface {
	CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error
	GetAdditionalChainList(ctx sdk.Context) (list []chains.Chain)
	GetPolicies(ctx sdk.Context) (val authoritytypes.Policies, found bool)
}
