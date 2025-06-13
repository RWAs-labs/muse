package config

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/runner"
	tonrunner "github.com/RWAs-labs/muse/e2e/runner/ton"
	btcclient "github.com/RWAs-labs/muse/museclient/chains/bitcoin/client"
	tonconfig "github.com/RWAs-labs/muse/museclient/chains/ton/config"
	museclientconfig "github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/retry"
	musecore_rpc "github.com/RWAs-labs/muse/pkg/rpc"
)

// getClientsFromConfig get clients from config
func getClientsFromConfig(ctx context.Context, conf config.Config, account config.Account) (runner.Clients, error) {
	btcRPCClient, err := getBtcClient(conf.RPCs.Bitcoin)
	if err != nil {
		return runner.Clients{}, fmt.Errorf("failed to get btc client: %w", err)
	}

	evmClient, evmAuth, err := getEVMClient(ctx, conf.RPCs.EVM, account)
	if err != nil {
		return runner.Clients{}, fmt.Errorf("failed to get evm client: %w", err)
	}

	var solanaClient *rpc.Client
	if conf.RPCs.Solana != "" {
		if solanaClient = rpc.New(conf.RPCs.Solana); solanaClient == nil {
			return runner.Clients{}, fmt.Errorf("failed to get solana client")
		}
	}

	var tonClient *tonrunner.Client
	if conf.RPCs.TON != "" {
		c, err := getTONClient(ctx, conf.RPCs.TON)
		if err != nil {
			return runner.Clients{}, fmt.Errorf("failed to get ton client: %w", err)
		}
		tonClient = c
	}

	var suiClient sui.ISuiAPI
	if conf.RPCs.Sui != "" {
		suiClient = sui.NewSuiClient(conf.RPCs.Sui)
	}

	museCoreClients, err := GetMusecoreClient(conf)
	if err != nil {
		return runner.Clients{}, fmt.Errorf("failed to get musecore client: %w", err)
	}

	mevmClient, mevmAuth, err := getEVMClient(ctx, conf.RPCs.Mevm, account)
	if err != nil {
		return runner.Clients{}, fmt.Errorf("failed to get mevm client: %w", err)
	}

	return runner.Clients{
		Musecore:          museCoreClients,
		BtcRPC:            btcRPCClient,
		Solana:            solanaClient,
		TON:               tonClient,
		Sui:               suiClient,
		Evm:               evmClient,
		EvmAuth:           evmAuth,
		Mevm:              mevmClient,
		MevmAuth:          mevmAuth,
		MuseclientMetrics: &runner.MetricsClient{URL: conf.RPCs.MuseclientMetrics},
	}, nil
}

// getBtcClient get btc client
func getBtcClient(e2eConfig config.BitcoinRPC) (*btcclient.Client, error) {
	cfg := museclientconfig.BTCConfig{
		RPCUsername: e2eConfig.User,
		RPCPassword: e2eConfig.Pass,
		RPCHost:     e2eConfig.Host,
		RPCParams:   string(e2eConfig.Params),
	}

	var chain chains.Chain
	switch e2eConfig.Params {
	case config.Regnet:
		chain = chains.BitcoinRegtest
	case config.Testnet3:
		chain = chains.BitcoinTestnet
	case config.Mainnet:
		chain = chains.BitcoinMainnet
	default:
		return nil, fmt.Errorf("invalid bitcoin params %s", e2eConfig.Params)
	}

	return btcclient.New(cfg, chain.ChainId, zerolog.Nop())
}

// getEVMClient get evm client
func getEVMClient(
	ctx context.Context,
	rpc string,
	account config.Account,
) (*ethclient.Client, *bind.TransactOpts, error) {
	evmClient, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial evm client: %w", err)
	}

	chainid, err := evmClient.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chain id: %w", err)
	}
	privKey, err := account.PrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get deployer privkey: %w", err)
	}
	evmAuth, err := bind.NewKeyedTransactorWithChainID(privKey, chainid)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get keyed transactor: %w", err)
	}

	return evmClient, evmAuth, nil
}

// getTONClient resolved tonrunner based on lite-server config (path or url)
func getTONClient(ctx context.Context, configURLOrPath string) (*tonrunner.Client, error) {
	if configURLOrPath == "" {
		return nil, fmt.Errorf("config is empty")
	}

	// It might take some time to bootstrap the sidecar
	cfg, err := retry.DoTypedWithRetry(
		func() (*tonconfig.GlobalConfigurationFile, error) {
			return tonconfig.FromSource(ctx, configURLOrPath)
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get ton config: %w", err)
	}

	return tonrunner.NewClient(cfg)
}

func GetMusecoreClient(conf config.Config) (musecore_rpc.Clients, error) {
	if conf.RPCs.MuseCoreGRPC != "" {
		return musecore_rpc.NewGRPCClients(
			conf.RPCs.MuseCoreGRPC,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	}
	if conf.RPCs.MuseCoreRPC != "" {
		return musecore_rpc.NewCometBFTClients(conf.RPCs.MuseCoreRPC)
	}
	return musecore_rpc.Clients{}, fmt.Errorf("no MuseCore gRPC or RPC specified")
}
