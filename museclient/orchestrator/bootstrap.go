package orchestrator

import (
	"context"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	tontools "github.com/tonkeeper/tongo/ton"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/chains/bitcoin"
	"github.com/RWAs-labs/muse/museclient/chains/bitcoin/client"
	btcobserver "github.com/RWAs-labs/muse/museclient/chains/bitcoin/observer"
	btcsigner "github.com/RWAs-labs/muse/museclient/chains/bitcoin/signer"
	"github.com/RWAs-labs/muse/museclient/chains/evm"
	evmclient "github.com/RWAs-labs/muse/museclient/chains/evm/client"
	evmobserver "github.com/RWAs-labs/muse/museclient/chains/evm/observer"
	evmsigner "github.com/RWAs-labs/muse/museclient/chains/evm/signer"
	"github.com/RWAs-labs/muse/museclient/chains/solana"
	solbserver "github.com/RWAs-labs/muse/museclient/chains/solana/observer"
	solanasigner "github.com/RWAs-labs/muse/museclient/chains/solana/signer"
	"github.com/RWAs-labs/muse/museclient/chains/sui"
	suiclient "github.com/RWAs-labs/muse/museclient/chains/sui/client"
	suiobserver "github.com/RWAs-labs/muse/museclient/chains/sui/observer"
	suisigner "github.com/RWAs-labs/muse/museclient/chains/sui/signer"
	"github.com/RWAs-labs/muse/museclient/chains/ton"
	"github.com/RWAs-labs/muse/museclient/chains/ton/liteapi"
	tonobserver "github.com/RWAs-labs/muse/museclient/chains/ton/observer"
	tonsigner "github.com/RWAs-labs/muse/museclient/chains/ton/signer"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/museclient/db"
	"github.com/RWAs-labs/muse/museclient/keys"
	"github.com/RWAs-labs/muse/pkg/chains"
	suigateway "github.com/RWAs-labs/muse/pkg/contracts/sui"
	toncontracts "github.com/RWAs-labs/muse/pkg/contracts/ton"
)

const btcBlocksPerDay = 144

func (oc *Orchestrator) bootstrapBitcoin(ctx context.Context, chain zctx.Chain) (*bitcoin.Bitcoin, error) {
	// should not happen
	if !chain.IsBitcoin() {
		return nil, errors.New("chain is not bitcoin")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	cfg, found := app.Config().GetBTCConfig(chain.ID())
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find btc config")
	}

	rpcClient, err := client.New(cfg, chain.ID(), oc.logger.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create rpc client")
	}

	var (
		rawChain = chain.RawChain()
		dbName   = btcDatabaseFileName(*rawChain)
	)

	baseObserver, err := oc.newBaseObserver(chain, dbName)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	observer, err := btcobserver.New(*rawChain, baseObserver, rpcClient)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer")
	}

	baseSigner := oc.newBaseSigner(chain)
	signer := btcsigner.New(baseSigner, rpcClient)

	return bitcoin.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) bootstrapEVM(ctx context.Context, chain zctx.Chain) (*evm.EVM, error) {
	// should not happen
	if !chain.IsEVM() {
		return nil, errors.New("chain is not EVM")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	cfg, found := app.Config().GetEVMConfig(chain.ID())
	if !found || cfg.Empty() {
		return nil, errors.Wrap(errSkipChain, "unable to find evm config")
	}

	evmClient, err := evmclient.NewFromEndpoint(ctx, cfg.Endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create evm client (%s)", cfg.Endpoint)
	}

	baseObserver, err := oc.newBaseObserver(chain, chain.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	observer, err := evmobserver.New(baseObserver, evmClient)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer")
	}

	var (
		museConnectorAddress = ethcommon.HexToAddress(chain.Params().ConnectorContractAddress)
		erc20CustodyAddress  = ethcommon.HexToAddress(chain.Params().Erc20CustodyContractAddress)
		gatewayAddress       = ethcommon.HexToAddress(chain.Params().GatewayAddress)
	)

	signer, err := evmsigner.New(
		oc.newBaseSigner(chain),
		evmClient,
		museConnectorAddress,
		erc20CustodyAddress,
		gatewayAddress,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create signer")
	}

	return evm.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) bootstrapSolana(ctx context.Context, chain zctx.Chain) (*solana.Solana, error) {
	// should not happen
	if !chain.IsSolana() {
		return nil, errors.New("chain is not Solana")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	cfg, found := app.Config().GetSolanaConfig()
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find solana config")
	}

	baseObserver, err := oc.newBaseObserver(chain, chain.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	gwAddress := chain.Params().GatewayAddress

	rpcClient := solrpc.New(cfg.Endpoint)
	if rpcClient == nil {
		return nil, errors.New("unable to create rpc client")
	}

	observer, err := solbserver.New(baseObserver, rpcClient, gwAddress)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer")
	}

	// Try loading Solana relayer key if present
	// Note that relayerKey might be nil if the key is not present. It's okay.
	password := chain.RelayerKeyPassword()
	relayerKey, err := keys.LoadRelayerKey(app.Config().GetRelayerKeyPath(), chain.RawChain().Network, password)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load relayer key")
	}

	baseSigner := oc.newBaseSigner(chain)

	// create Solana signer
	signer, err := solanasigner.New(baseSigner, rpcClient, gwAddress, relayerKey)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create signer")
	}

	return solana.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) bootstrapSui(ctx context.Context, chain zctx.Chain) (*sui.Sui, error) {
	// should not happen
	if !chain.IsSui() {
		return nil, errors.New("chain is not sui")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	cfg, found := app.Config().GetSuiConfig()
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find sui config")
	}

	// note that gw address should be in format of `$packageID,$gatewayObjectID`
	gateway, err := suigateway.NewGatewayFromPairID(chain.Params().GatewayAddress)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create gateway")
	}

	suiClient := suiclient.NewFromEndpoint(cfg.Endpoint)

	baseObserver, err := oc.newBaseObserver(chain, chain.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	observer := suiobserver.New(baseObserver, suiClient, gateway)

	signer := suisigner.New(oc.newBaseSigner(chain), suiClient, gateway, oc.deps.Musecore)

	return sui.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) bootstrapTON(ctx context.Context, chain zctx.Chain) (*ton.TON, error) {
	// should not happen
	if !chain.IsTON() {
		return nil, errors.New("chain is not TON")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	cfg, found := app.Config().GetTONConfig()
	if !found {
		return nil, errors.Wrap(errSkipChain, "unable to find TON config")
	}

	gwAddress := chain.Params().GatewayAddress
	if gwAddress == "" {
		return nil, errors.New("gateway address is empty")
	}

	gatewayID, err := tontools.ParseAccountID(gwAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse gateway address %q", gwAddress)
	}

	gw := toncontracts.NewGateway(gatewayID)

	client, err := tonResolveClient(ctx, cfg.LiteClientConfigURL)
	if err != nil {
		return nil, errors.Wrap(err, "unable to resolve TON liteclient")
	}

	baseObserver, err := oc.newBaseObserver(chain, chain.Name())
	if err != nil {
		return nil, errors.Wrap(err, "unable to create base observer")
	}

	observer, err := tonobserver.New(baseObserver, client, gw)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create observer")
	}

	signer := tonsigner.New(oc.newBaseSigner(chain), client, gw)

	return ton.New(oc.scheduler, observer, signer), nil
}

func (oc *Orchestrator) newBaseObserver(chain zctx.Chain, dbName string) (*base.Observer, error) {
	var (
		rawChain       = chain.RawChain()
		rawChainParams = chain.Params()
	)

	database, err := db.NewFromSqlite(oc.deps.DBPath, dbName, true)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open database %s", dbName)
	}

	blocksCacheSize := base.DefaultBlockCacheSize
	if chain.IsBitcoin() {
		blocksCacheSize = btcBlocksPerDay
	}

	return base.NewObserver(
		*rawChain,
		*rawChainParams,
		oc.deps.Musecore,
		oc.deps.TSS,
		blocksCacheSize,
		oc.deps.Telemetry,
		database,
		oc.logger.base,
	)
}

func (oc *Orchestrator) newBaseSigner(chain zctx.Chain) *base.Signer {
	return base.NewSigner(*chain.RawChain(), oc.deps.TSS, oc.logger.base)
}

func btcDatabaseFileName(chain chains.Chain) string {
	// legacyBTCDatabaseFilename is the Bitcoin database file name now used in mainnet and testnet3
	// so we keep using it here for backward compatibility
	const legacyBTCDatabaseFilename = "btc_chain_client"

	// For additional bitcoin networks, we use the chain name as the database file name
	switch chain.ChainId {
	case chains.BitcoinMainnet.ChainId, chains.BitcoinTestnet.ChainId:
		return legacyBTCDatabaseFilename
	default:
		return fmt.Sprintf("%s_%s", legacyBTCDatabaseFilename, chain.Name)
	}
}

type (
	tonClientCtxKey struct{}
	tonClient       interface {
		tonobserver.LiteClient
		tonsigner.LiteClient
	}
)

// tonResolveClient resolves lite-api from a source OR from the context.
// The latter is used in testing because it's challenging to mock the entire lite-api e2e
// as it relies on low-level encrypted connections (otherwise could be wrapped with `httptest`)
func tonResolveClient(ctx context.Context, configSource string) (tonClient, error) {
	client, ok := ctx.Value(tonClientCtxKey{}).(tonClient)
	if ok {
		return client, nil
	}

	client, err := liteapi.NewFromSource(ctx, configSource)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create TON liteapi from %q", configSource)
	}

	return client, nil
}
