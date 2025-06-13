package config

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/showa-93/go-mask"
)

// KeyringBackend is the type of keyring backend to use for the hotkey
type KeyringBackend string

const (
	// KeyringBackendUndefined is undefined keyring backend
	KeyringBackendUndefined KeyringBackend = ""

	// KeyringBackendTest is the test Cosmos keyring backend
	KeyringBackendTest KeyringBackend = "test"

	// KeyringBackendFile is the file Cosmos keyring backend
	KeyringBackendFile KeyringBackend = "file"

	DefaultRelayerDir = "relayer-keys"

	// DefaultRelayerKeyPath is the default path that relayer keys are stored
	DefaultRelayerKeyPath = "~/.musecored/" + DefaultRelayerDir
)

// ClientConfiguration is a subset of museclient config that is used by musecore client
type ClientConfiguration struct {
	ChainHost       string `json:"chain_host"        mapstructure:"chain_host"`
	ChainRPC        string `json:"chain_rpc"         mapstructure:"chain_rpc"`
	ChainHomeFolder string `json:"chain_home_folder" mapstructure:"chain_home_folder"`
	SignerName      string `json:"signer_name"       mapstructure:"signer_name"`
	SignerPasswd    string `json:"signer_passwd"`
}

// EVMConfig is the config for EVM chain
type EVMConfig struct {
	Endpoint string `mask:"filled"`
}

// BTCConfig is the config for Bitcoin chain
type BTCConfig struct {
	// the following are rpcclient ConnConfig fields
	RPCUsername string `mask:"filled"`
	RPCPassword string `mask:"filled"`
	RPCHost     string `mask:"filled"`
	RPCParams   string // "regtest", "mainnet", "testnet3" , "signet", "testnet4"
}

// SolanaConfig is the config for Solana chain
type SolanaConfig struct {
	Endpoint string `mask:"filled"`
}

// SuiConfig is the config for Sui chain
type SuiConfig struct {
	Endpoint string `mask:"filled"`
}

// TONConfig is the config for TON chain
type TONConfig struct {
	// Can be either URL of local file path
	LiteClientConfigURL string `json:"liteClientConfigURL"`
}

// ComplianceConfig is the config for compliance
type ComplianceConfig struct {
	LogPath string `json:"LogPath"`
	// Deprecated: use the separate restricted addresses config
	RestrictedAddresses []string `json:"RestrictedAddresses" mask:"zero"`
}

// Config is the config for MuseClient
// TODO: use snake case for json fields
// https://github.com/RWAs-labs/muse/issues/1020
type Config struct {
	Peer                    string         `json:"Peer"`
	PublicIP                string         `json:"PublicIP"`
	LogFormat               string         `json:"LogFormat"`
	LogLevel                int8           `json:"LogLevel"`
	LogSampler              bool           `json:"LogSampler"`
	PreParamsPath           string         `json:"PreParamsPath"`
	MuseCoreHome            string         `json:"MuseCoreHome"`
	ChainID                 string         `json:"ChainID"`
	MuseCoreURL             string         `json:"MuseCoreURL"`
	AuthzGranter            string         `json:"AuthzGranter"`
	AuthzHotkey             string         `json:"AuthzHotkey"`
	P2PDiagnostic           bool           `json:"P2PDiagnostic"`
	ConfigUpdateTicker      uint64         `json:"ConfigUpdateTicker"`
	P2PDiagnosticTicker     uint64         `json:"P2PDiagnosticTicker"`
	TssPath                 string         `json:"TssPath"`
	TSSMaxPendingSignatures uint64         `json:"TSSMaxPendingSignatures"`
	TestTssKeysign          bool           `json:"TestTssKeysign"`
	KeyringBackend          KeyringBackend `json:"KeyringBackend"`
	RelayerKeyPath          string         `json:"RelayerKeyPath"`

	// chain configs
	EVMChainConfigs map[int64]EVMConfig `json:"EVMChainConfigs"`
	BTCChainConfigs map[int64]BTCConfig `json:"BTCChainConfigs"`
	// Deprecated: the 'BitcoinConfig' will be removed once the 'BTCChainConfigs' is fully adopted
	BitcoinConfig BTCConfig    `json:"BitcoinConfig"`
	SolanaConfig  SolanaConfig `json:"SolanaConfig"`
	SuiConfig     SuiConfig    `json:"SuiConfig"`
	TONConfig     TONConfig    `json:"TONConfig"`

	// compliance config
	ComplianceConfig ComplianceConfig `json:"ComplianceConfig"`

	mu *sync.RWMutex
}

// GetEVMConfig returns the EVM config for the given chain ID
func (c Config) GetEVMConfig(chainID int64) (EVMConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	evmCfg := c.EVMChainConfigs[chainID]
	return evmCfg, !evmCfg.Empty()
}

// GetAllEVMConfigs returns a map of all EVM configs
func (c Config) GetAllEVMConfigs() map[int64]EVMConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// deep copy evm configs
	copied := make(map[int64]EVMConfig, len(c.EVMChainConfigs))
	for chainID, evmConfig := range c.EVMChainConfigs {
		copied[chainID] = evmConfig
	}
	return copied
}

// GetBTCConfig returns the BTC config for the given chain ID
func (c Config) GetBTCConfig(chainID int64) (BTCConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// we prefer 'BTCChainConfigs' over 'BitcoinConfig' but still fallback to be backward compatible
	// this will allow new 'museclientd' binary to work with old config file
	btcCfg, found := c.BTCChainConfigs[chainID]
	if !found || btcCfg.Empty() {
		btcCfg = c.BitcoinConfig
	}

	return btcCfg, !btcCfg.Empty()
}

// GetSolanaConfig returns the Solana config
func (c Config) GetSolanaConfig() (SolanaConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.SolanaConfig, c.SolanaConfig != (SolanaConfig{})
}

// GetSuiConfig returns the Sui config
func (c Config) GetSuiConfig() (SuiConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.SuiConfig, c.SuiConfig != (SuiConfig{})
}

// GetTONConfig returns the TONConfig and a bool indicating if it's present.
func (c Config) GetTONConfig() (TONConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.TONConfig, c.TONConfig != TONConfig{}
}

// StringMasked returns the string representation of the config with sensitive fields masked.
// Currently only the endpoints and bitcoin credentials are masked.
func (c Config) StringMasked() string {
	// create a masker
	masker := mask.NewMasker()
	masker.RegisterMaskStringFunc(mask.MaskTypeFilled, masker.MaskFilledString)
	masker.RegisterMaskAnyFunc(mask.MaskTypeFilled, masker.MaskZero)

	// mask the config
	masked, err := masker.Mask(c)
	if err != nil {
		return ""
	}

	s, err := json.MarshalIndent(masked, "", "\t")
	if err != nil {
		return ""
	}
	return string(s)
}

// GetRestrictedAddressBook returns a map of restricted addresses
// Note: the restricted address book contains both ETH and BTC addresses
func (c Config) GetRestrictedAddressBook() map[string]bool {
	restrictedAddresses := make(map[string]bool)
	for _, address := range c.ComplianceConfig.RestrictedAddresses {
		if address != "" {
			restrictedAddresses[strings.ToLower(address)] = true
		}
	}
	return restrictedAddresses
}

// GetKeyringBackend returns the keyring backend
func (c Config) GetKeyringBackend() KeyringBackend {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.KeyringBackend
}

// GetRelayerKeyPath returns the relayer key path
func (c Config) GetRelayerKeyPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// use default path if not configured
	if c.RelayerKeyPath == "" {
		return DefaultRelayerKeyPath
	}
	return c.RelayerKeyPath
}

func (c EVMConfig) Empty() bool {
	return c.Endpoint == ""
}

func (c BTCConfig) Empty() bool {
	return c.RPCHost == ""
}
