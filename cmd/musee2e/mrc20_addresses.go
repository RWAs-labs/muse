package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	musee2econfig "github.com/RWAs-labs/muse/cmd/musee2e/config"
	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/pkg/chains"
	fungibletypes "github.com/RWAs-labs/muse/x/fungible/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

const flagMRC20Network = "mrc20-network"
const flagMRC20Symbol = "mrc20-symbol"

func registerERC20Flags(cmd *cobra.Command) {
	cmd.Flags().String(flagMRC20Network, "", "network from /muse-chain/observer/supportedChains")
	cmd.Flags().String(flagMRC20Symbol, "", "symbol from /muse-chain/fungible/foreign_coins")
}

func processMRC20Flags(cmd *cobra.Command, conf *config.Config) error {
	mrc20ChainName, err := cmd.Flags().GetString(flagMRC20Network)
	if err != nil {
		return err
	}
	mrc20Symbol, err := cmd.Flags().GetString(flagMRC20Symbol)
	if err != nil {
		return err
	}
	if mrc20ChainName != "" && mrc20Symbol != "" {
		erc20Asset, mrc20ContractAddress, chain, err := findMRC20(
			cmd.Context(),
			conf,
			mrc20ChainName,
			mrc20Symbol,
		)
		if err != nil {
			return err
		}
		if chain.IsEVMChain() {
			conf.Contracts.EVM.ERC20 = config.DoubleQuotedString(erc20Asset)
			conf.Contracts.MEVM.ERC20MRC20Addr = config.DoubleQuotedString(mrc20ContractAddress)
		} else if chain.IsSolanaChain() {
			conf.Contracts.Solana.SPLAddr = config.DoubleQuotedString(erc20Asset)
			conf.Contracts.MEVM.SPLMRC20Addr = config.DoubleQuotedString(mrc20ContractAddress)
		}
	}
	return nil
}

// findMRC20 loads ERC20/SPL/etc addresses via gRPC given CLI flags
func findMRC20(
	ctx context.Context,
	conf *config.Config,
	networkName, mrc20Symbol string,
) (string, string, chains.Chain, error) {
	clients, err := musee2econfig.GetMusecoreClient(*conf)
	if err != nil {
		return "", "", chains.Chain{}, fmt.Errorf("get muse clients: %w", err)
	}

	supportedChainsRes, err := clients.Observer.SupportedChains(ctx, &observertypes.QuerySupportedChains{})
	if err != nil {
		return "", "", chains.Chain{}, fmt.Errorf("get chain params: %w", err)
	}

	chainID := int64(0)
	for _, chain := range supportedChainsRes.Chains {
		if strings.EqualFold(chain.Network.String(), networkName) {
			chainID = chain.ChainId
			break
		}
	}
	if chainID == 0 {
		return "", "", chains.Chain{}, fmt.Errorf("chain %s not found", networkName)
	}

	chain, ok := chains.GetChainFromChainID(chainID, nil)
	if !ok {
		return "", "", chains.Chain{}, fmt.Errorf("invalid/unknown chain ID %d", chainID)
	}

	foreignCoinsRes, err := clients.Fungible.ForeignCoinsAll(ctx, &fungibletypes.QueryAllForeignCoinsRequest{})
	if err != nil {
		return "", "", chain, fmt.Errorf("get foreign coins: %w", err)
	}

	for _, coin := range foreignCoinsRes.ForeignCoins {
		if coin.ForeignChainId != chainID {
			continue
		}
		// sometimes symbol is USDT, sometimes it's like USDT.SEPOLIA
		if strings.HasPrefix(coin.Symbol, mrc20Symbol) || strings.HasSuffix(coin.Symbol, mrc20Symbol) {
			return coin.Asset, coin.Mrc20ContractAddress, chain, nil
		}
	}
	return "", "", chain, fmt.Errorf("mrc20 %s not found on %s", mrc20Symbol, networkName)
}
