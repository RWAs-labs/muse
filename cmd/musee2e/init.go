package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/cmd/musee2e/local"
	"github.com/RWAs-labs/muse/e2e/config"
)

var initConf = config.DefaultConfig()
var configFile = ""

func NewInitCmd() *cobra.Command {
	var InitCmd = &cobra.Command{
		Use:   "init",
		Short: "initialize config file for e2e tests",
		RunE:  initConfig,
	}

	InitCmd.Flags().StringVar(&initConf.RPCs.EVM, "ethURL", initConf.RPCs.EVM, "--ethURL http://eth:8545")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.MuseCoreGRPC, "grpcURL", initConf.RPCs.MuseCoreGRPC, "--grpcURL musecore0:9090")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.MuseCoreRPC, "rpcURL", initConf.RPCs.MuseCoreRPC, "--rpcURL http://musecore0:26657")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.Mevm, "mevmURL", initConf.RPCs.Mevm, "--mevmURL http://musecore0:8545")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.Bitcoin.Host, "btcURL", initConf.RPCs.Bitcoin.Host, "--btcURL bitcoin:18443")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.Solana, "solanaURL", initConf.RPCs.Solana, "--solanaURL http://solana:8899")
	InitCmd.Flags().
		StringVar(&initConf.RPCs.TON, "tonURL", initConf.RPCs.TON, "--tonURL http://ton:8000/lite-client.json")
	InitCmd.Flags().StringVar(&initConf.MuseChainID, "chainID", initConf.MuseChainID, "--chainID athens_101-1")
	InitCmd.Flags().StringVar(&configFile, local.FlagConfigFile, "e2e.config", "--cfg ./e2e.config")

	return InitCmd
}

func initConfig(_ *cobra.Command, _ []string) error {
	err := initConf.GenerateKeys()
	if err != nil {
		return fmt.Errorf("generating keys: %w", err)
	}
	err = config.WriteConfig(configFile, initConf)
	if err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}
