package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/cmd/musetool/cli"
	"github.com/RWAs-labs/muse/cmd/musetool/config"
)

var rootCmd = &cobra.Command{
	Use:   "musetool",
	Short: "utility tool for muse-chain",
}

func init() {
	rootCmd.AddCommand(cli.NewGetInboundBallotCMD())
	rootCmd.AddCommand(cli.NewTrackCCTXCMD())
	rootCmd.AddCommand(cli.NewApplicationDBStatsCMD())
	rootCmd.PersistentFlags().String(config.FlagConfig, "", "custom config file: --config filename.json")
	rootCmd.PersistentFlags().
		Bool(config.FlagDebug, false, "enable debug mode, to show more details on why the command might be failing")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		os.Exit(1)
	}
}
