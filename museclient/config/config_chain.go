package config

import (
	"sync"

	"github.com/RWAs-labs/muse/pkg/chains"
)

const (
	// MaxBlocksPerScan is the maximum number of blocks to scan in one ticker
	MaxBlocksPerScan = 100
)

// New constructs Config optionally with default values.
func New(setDefaults bool) Config {
	cfg := Config{
		EVMChainConfigs: make(map[int64]EVMConfig),
		BTCChainConfigs: make(map[int64]BTCConfig),

		mu: &sync.RWMutex{},
	}

	if setDefaults {
		cfg.EVMChainConfigs = evmChainsConfigs()
		cfg.BTCChainConfigs = btcChainsConfigs()
		cfg.SolanaConfig = solanaConfigLocalnet()
		cfg.SuiConfig = suiConfigLocalnet()
		cfg.TONConfig = tonConfigLocalnet()
	}

	return cfg
}

// bitcoinConfigRegnet contains Bitcoin config for regnet
func bitcoinConfigRegnet() BTCConfig {
	return BTCConfig{
		// `smoketest` is the previous name for E2E test,
		// we keep this name for compatibility between client versions in upgrade test
		RPCUsername: "smoketest",
		RPCPassword: "123",
		RPCHost:     "bitcoin:18443",
		RPCParams:   "regtest",
	}
}

// solanaConfigLocalnet contains config for Solana localnet
func solanaConfigLocalnet() SolanaConfig {
	return SolanaConfig{
		Endpoint: "http://solana:8899",
	}
}

func suiConfigLocalnet() SuiConfig {
	return SuiConfig{
		Endpoint: "http://sui:9000",
	}
}

func tonConfigLocalnet() TONConfig {
	return TONConfig{
		LiteClientConfigURL: "http://ton:8000/lite-client.json",
	}
}

// evmChainsConfigs contains EVM chain configs
// it contains list of EVM chains with empty endpoint except for localnet
func evmChainsConfigs() map[int64]EVMConfig {
	return map[int64]EVMConfig{
		chains.GoerliLocalnet.ChainId: {
			Endpoint: "http://eth:8545",
		},
	}
}

// btcChainsConfigs contains BTC chain configs
func btcChainsConfigs() map[int64]BTCConfig {
	return map[int64]BTCConfig{
		chains.BitcoinRegtest.ChainId: bitcoinConfigRegnet(),
	}
}
