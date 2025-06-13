package chains

import "fmt"

var (
	/**
	* Mainnet chains
	 */

	// MuseChainMainnet is the mainnet chain for Muse
	MuseChainMainnet = Chain{
		ChainName:   ChainName_muse_mainnet,
		ChainId:     7000,
		Network:     Network_muse,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_tendermint,
		IsExternal:  false,
		CctxGateway: CCTXGateway_mevm,
		Name:        "muse_mainnet",
	}

	// Ethereum is Ethereum mainnet
	Ethereum = Chain{
		ChainName:   ChainName_eth_mainnet,
		ChainId:     1,
		Network:     Network_eth,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "eth_mainnet",
	}

	// BscMainnet is Binance Smart Chain mainnet
	BscMainnet = Chain{
		ChainName:   ChainName_bsc_mainnet,
		ChainId:     56,
		Network:     Network_bsc,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "bsc_mainnet",
	}

	// BitcoinMainnet is Bitcoin mainnet
	BitcoinMainnet = Chain{
		ChainName:   ChainName_btc_mainnet,
		ChainId:     8332,
		Network:     Network_btc,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_no_vm,
		Consensus:   Consensus_bitcoin,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "btc_mainnet",
	}

	// Polygon is Polygon mainnet
	Polygon = Chain{
		ChainName:   ChainName_polygon_mainnet,
		ChainId:     137,
		Network:     Network_polygon,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "polygon_mainnet",
	}

	// OptimismMainnet is Optimism mainnet
	OptimismMainnet = Chain{
		ChainName:   ChainName_optimism_mainnet,
		ChainId:     10,
		Network:     Network_optimism,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_op_stack,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "optimism_mainnet",
	}

	// BaseMainnet is Base mainnet
	BaseMainnet = Chain{
		ChainName:   ChainName_base_mainnet,
		ChainId:     8453,
		Network:     Network_base,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_op_stack,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "base_mainnet",
	}

	// AvalancheMainnet is the Avalanche mainnet C-Chain
	AvalancheMainnet = Chain{
		ChainId:     43114,
		Network:     Network_avalanche,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_snowman,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "avalanche_mainnet",
	}

	// ArbitrumMainnet is the Arbitrum mainnet
	ArbitrumMainnet = Chain{
		ChainId:     42161,
		Network:     Network_arbitrum,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_arbitrum_nitro,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "arbitrum_mainnet",
	}

	// SuiMainnet is the Sui mainnet
	// TODO: this value should be set to 101 but currently conflicts with MuseChain localnet
	// https://github.com/RWAs-labs/muse/issues/3491
	SuiMainnet = Chain{
		ChainId:     105,
		Network:     Network_sui,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_mvm_sui,
		Consensus:   Consensus_sui_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "sui_mainnet",
	}

	// WorldMainnet is the World Chain mainnet
	WorldMainnet = Chain{
		ChainId:     480,
		Network:     Network_worldchain,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_op_stack,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "world_mainnet",
	}

	// SolanaMainnet is Solana mainnet
	// TODO: define final chain ID
	// https://github.com/RWAs-labs/muse/issues/2421
	SolanaMainnet = Chain{
		ChainName:   ChainName_solana_mainnet,
		ChainId:     900,
		Network:     Network_solana,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_svm,
		Consensus:   Consensus_solana_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "solana_mainnet",
	}

	TONMainnet = Chain{
		// T[20] O[15] N[14] mainnet[0] :)
		ChainId:     2015140,
		Network:     Network_ton,
		NetworkType: NetworkType_mainnet,
		Vm:          Vm_tvm,
		Consensus:   Consensus_catchain_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "ton_mainnet",
	}

	/**
	* Testnet chains
	 */

	// MuseChainTestnet is the testnet chain for Muse
	MuseChainTestnet = Chain{
		ChainName:   ChainName_muse_testnet,
		ChainId:     7001,
		Network:     Network_muse,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_tendermint,
		IsExternal:  false,
		CctxGateway: CCTXGateway_mevm,
		Name:        "muse_testnet",
	}

	// Sepolia is Ethereum sepolia testnet
	Sepolia = Chain{
		ChainName:   ChainName_sepolia_testnet,
		ChainId:     11155111,
		Network:     Network_eth,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "sepolia_testnet",
	}

	// BscTestnet is Binance Smart Chain testnet
	BscTestnet = Chain{
		ChainName:   ChainName_bsc_testnet,
		ChainId:     97,
		Network:     Network_bsc,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "bsc_testnet",
	}

	// BitcoinTestnet is Bitcoin testnet3
	BitcoinTestnet = Chain{
		ChainName:   ChainName_btc_testnet,
		ChainId:     18332,
		Network:     Network_btc,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_no_vm,
		Consensus:   Consensus_bitcoin,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "btc_testnet",
	}

	// BitcoinSignetTestnet is Bitcoin Signet testnet
	BitcoinSignetTestnet = Chain{
		ChainId:     18333,
		Network:     Network_btc,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_no_vm,
		Consensus:   Consensus_bitcoin,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "btc_signet_testnet",
	}

	// BitcoinTestnet4 is Bitcoin testnet4
	BitcoinTestnet4 = Chain{
		ChainId:     18334,
		Network:     Network_btc,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_no_vm,
		Consensus:   Consensus_bitcoin,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "btc_testnet4",
	}

	// Amoy is Polygon amoy testnet
	Amoy = Chain{
		ChainName:   ChainName_amoy_testnet,
		ChainId:     80002,
		Network:     Network_polygon,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "amoy_testnet",
	}

	// OptimismSepolia is Optimism sepolia testnet
	OptimismSepolia = Chain{
		ChainName:   ChainName_optimism_sepolia,
		ChainId:     11155420,
		Network:     Network_optimism,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_op_stack,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "optimism_sepolia",
	}

	// BaseSepolia is Base sepolia testnet
	BaseSepolia = Chain{
		ChainName:   ChainName_base_sepolia,
		ChainId:     84532,
		Network:     Network_base,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_op_stack,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "base_sepolia",
	}

	// AvalancheTestnet is the Avalanche testnet C-Chain
	AvalancheTestnet = Chain{
		ChainId:     43113,
		Network:     Network_avalanche,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_snowman,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "avalanche_testnet",
	}

	// ArbitrumSepolia is the Arbitrum sepolia testnet
	ArbitrumSepolia = Chain{
		ChainId:     421614,
		Network:     Network_arbitrum,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_arbitrum_nitro,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "arbitrum_sepolia",
	}

	// WorldTestnet is the World Chain testnet
	WorldTestnet = Chain{
		ChainId:     4801,
		Network:     Network_worldchain,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_op_stack,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "world_testnet",
	}

	// SuiTestnet is the Sui testnet
	SuiTestnet = Chain{
		ChainId:     103,
		Network:     Network_sui,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_mvm_sui,
		Consensus:   Consensus_sui_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "sui_testnet",
	}

	// SolanaDevnet is Solana devnet
	// NOTE: Solana devnet refers to Solana testnet in our terminology
	// Solana uses devnet denomitation for network for development
	// TODO: define final chain ID
	// https://github.com/RWAs-labs/muse/issues/2421
	SolanaDevnet = Chain{
		ChainName:   ChainName_solana_devnet,
		ChainId:     901,
		Network:     Network_solana,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_svm,
		Consensus:   Consensus_solana_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "solana_devnet",
	}

	TONTestnet = Chain{
		ChainId:     2015141,
		Network:     Network_ton,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_tvm,
		Consensus:   Consensus_catchain_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "ton_testnet",
	}

	/**
	* Devnet chains
	 */

	// MuseChainDevnet is the devnet chain for Muse
	// used as live testing environment
	MuseChainDevnet = Chain{
		ChainName:   ChainName_muse_mainnet,
		ChainId:     70000,
		Network:     Network_muse,
		NetworkType: NetworkType_devnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_tendermint,
		IsExternal:  false,
		CctxGateway: CCTXGateway_mevm,
		Name:        "muse_mainnet",
	}

	/**
	* Privnet chains
	 */

	// MuseChainPrivnet is the privnet chain for Muse (localnet)
	MuseChainPrivnet = Chain{
		ChainName:   ChainName_muse_mainnet,
		ChainId:     101,
		Network:     Network_muse,
		NetworkType: NetworkType_privnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_tendermint,
		IsExternal:  false,
		CctxGateway: CCTXGateway_mevm,
		Name:        "muse_mainnet",
	}

	// BitcoinRegtest is Bitcoin regtest (localnet)
	BitcoinRegtest = Chain{
		ChainName:   ChainName_btc_regtest,
		ChainId:     18444,
		Network:     Network_btc,
		NetworkType: NetworkType_privnet,
		Vm:          Vm_no_vm,
		Consensus:   Consensus_bitcoin,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "btc_regtest",
	}

	// GoerliLocalnet is Ethereum local goerli (localnet)
	GoerliLocalnet = Chain{
		ChainName:   ChainName_goerli_localnet,
		ChainId:     1337,
		Network:     Network_eth,
		NetworkType: NetworkType_privnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "goerli_localnet",
	}

	// SolanaLocalnet is Solana localnet
	// TODO: define final chain ID
	// https://github.com/RWAs-labs/muse/issues/2421
	SolanaLocalnet = Chain{
		ChainName:   ChainName_solana_localnet,
		ChainId:     902,
		Network:     Network_solana,
		NetworkType: NetworkType_privnet,
		Vm:          Vm_svm,
		Consensus:   Consensus_solana_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "solana_localnet",
	}

	TONLocalnet = Chain{
		ChainId:     2015142,
		Network:     Network_ton,
		NetworkType: NetworkType_privnet,
		Vm:          Vm_tvm,
		Consensus:   Consensus_catchain_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "ton_localnet",
	}

	SuiLocalnet = Chain{
		ChainId:     104,
		Network:     Network_sui,
		NetworkType: NetworkType_privnet,
		Vm:          Vm_mvm_sui,
		Consensus:   Consensus_sui_consensus,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "sui_localnet",
	}

	/**
	* Deprecated chains
	 */

	// Goerli is Ethereum goerli testnet (deprecated for sepolia)
	Goerli = Chain{
		ChainName:   ChainName_goerli_testnet,
		ChainId:     5,
		Network:     Network_eth,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "goerli_testnet",
	}

	// Mumbai is Polygon mumbai testnet (deprecated for amoy)
	Mumbai = Chain{
		ChainName:   ChainName_mumbai_testnet,
		ChainId:     80001,
		Network:     Network_polygon,
		NetworkType: NetworkType_testnet,
		Vm:          Vm_evm,
		Consensus:   Consensus_ethereum,
		IsExternal:  true,
		CctxGateway: CCTXGateway_observers,
		Name:        "mumbai_testnet",
	}
)

// ErrNotMuseChain is the error for chain not being a MuseChain chain
var ErrNotMuseChain = fmt.Errorf("chain is not a MuseChain chain")

// BtcNonceMarkOffset is the offset satoshi amount to calculate the nonce mark output
func BtcNonceMarkOffset() int64 {
	return 2000
}

// DefaultChainsList returns a list of default chains
func DefaultChainsList() []Chain {
	return []Chain{
		BitcoinMainnet,
		BscMainnet,
		Ethereum,
		BitcoinTestnet,
		BitcoinSignetTestnet,
		BitcoinTestnet4,
		Mumbai,
		Amoy,
		BscTestnet,
		Goerli,
		Sepolia,
		BitcoinRegtest,
		GoerliLocalnet,
		MuseChainMainnet,
		MuseChainTestnet,
		MuseChainDevnet,
		MuseChainPrivnet,
		Polygon,
		OptimismMainnet,
		OptimismSepolia,
		BaseMainnet,
		BaseSepolia,
		SolanaMainnet,
		SolanaDevnet,
		SolanaLocalnet,
		TONMainnet,
		TONTestnet,
		TONLocalnet,
		AvalancheMainnet,
		AvalancheTestnet,
		ArbitrumMainnet,
		ArbitrumSepolia,
		WorldMainnet,
		WorldTestnet,
		SuiMainnet,
		SuiTestnet,
		SuiLocalnet,
	}
}

// ChainListByNetworkType returns a list of chains by network type
func ChainListByNetworkType(networkType NetworkType, additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.NetworkType == networkType {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ChainListByNetwork returns a list of chains by network
func ChainListByNetwork(network Network, additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.Network == network {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ExternalChainList returns a list chains that are not Muse
func ExternalChainList(additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.IsExternal {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// ChainListByConsensus returns a list of chains by consensus
func ChainListByConsensus(consensus Consensus, additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.Consensus == consensus {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

func ChainListByGateway(gateway CCTXGateway, additionalChains []Chain) []Chain {
	var chainList []Chain
	for _, chain := range CombineDefaultChainsList(additionalChains) {
		if chain.CctxGateway == gateway {
			chainList = append(chainList, chain)
		}
	}
	return chainList
}

// MuseChainFromCosmosChainID returns a MuseChain chain object from a Cosmos chain ID
func MuseChainFromCosmosChainID(chainID string) (Chain, error) {
	ethChainID, err := CosmosToEthChainID(chainID)
	if err != nil {
		return Chain{}, err
	}

	return MuseChainFromChainID(ethChainID)
}

// MuseChainFromChainID returns a MuseChain chain object from a chain ID
func MuseChainFromChainID(chainID int64) (Chain, error) {
	switch chainID {
	case MuseChainPrivnet.ChainId:
		return MuseChainPrivnet, nil
	case MuseChainMainnet.ChainId:
		return MuseChainMainnet, nil
	case MuseChainTestnet.ChainId:
		return MuseChainTestnet, nil
	case MuseChainDevnet.ChainId:
		return MuseChainDevnet, nil
	default:
		return Chain{}, ErrNotMuseChain
	}
}

// CombineDefaultChainsList combines the default chains list with a list of chains
// duplicated chain ID are overwritten by the second list
func CombineDefaultChainsList(chains []Chain) []Chain {
	return CombineChainList(DefaultChainsList(), chains)
}

// CombineChainList combines a list of chains with a list of chains
// duplicated chain ID are overwritten by the second list
func CombineChainList(base []Chain, additional []Chain) []Chain {
	combined := make([]Chain, 0, len(base)+len(additional))
	combined = append(combined, base...)

	// map chain ID in combined to index in the list
	chainIDIndexMap := make(map[int64]int)
	for i, chain := range combined {
		chainIDIndexMap[chain.ChainId] = i
	}

	// add chains2 to combined
	// if chain ID already exists in chains1, overwrite it
	for _, chain := range additional {
		if index, ok := chainIDIndexMap[chain.ChainId]; ok {
			combined[index] = chain
		} else {
			combined = append(combined, chain)
		}
	}

	return combined
}
