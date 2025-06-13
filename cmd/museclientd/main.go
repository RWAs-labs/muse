package main

import (
	"context"
	"fmt"
	"os"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/app"
	"github.com/RWAs-labs/muse/pkg/constant"
	_ "github.com/RWAs-labs/muse/pkg/sdkconfig/default"
)

var (
	RootCmd = &cobra.Command{
		Use:   "museclientd",
		Short: "museclient cli & server",
	}
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "prints version",
		Run:   func(_ *cobra.Command, _ []string) { fmt.Print(constant.Version) },
	}

	InitializeConfigCmd = &cobra.Command{
		Use:     "init-config",
		Aliases: []string{"init"},
		Short:   "Initialize Museclient Configuration file",
		RunE:    InitializeConfig,
	}
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start MuseClient Observer",
		RunE:  Start,
	}

	TSSCmd        = &cobra.Command{Use: "tss", Short: "TSS commands"}
	TSSEncryptCmd = &cobra.Command{
		Use:   "encrypt [file-path] [secret-key]",
		Short: "Utility command to encrypt existing tss key-share file",
		Args:  cobra.ExactArgs(2),
		RunE:  TSSEncryptFile,
	}
	TSSGeneratePreParamsCmd = &cobra.Command{
		Use:   "gen-pre-params [path]",
		Short: "Generate pre parameters for TSS",
		Args:  cobra.ExactArgs(1),
		RunE:  TSSGeneratePreParams,
	}

	RelayerCmd          = &cobra.Command{Use: "relayer", Short: "Relayer commands"}
	RelayerImportKeyCmd = &cobra.Command{
		Use:   "import-key --network=<net> --private-key=<pk> --password=<pass> --relayer-key-path=<path>",
		Short: "Import a relayer private key",
		RunE:  RelayerImportKey,
	}
	RelayerShowAddressCmd = &cobra.Command{
		Use:   "show-address --network=<new> --password=<pass> --relayer-key-path=<path>",
		Short: "Show relayer address",
		RunE:  RelayerShowAddress,
	}
)

// globalOptions defines the global options for all commands.
type globalOptions struct {
	MusecoreHome string
}

var globalOpts globalOptions

func setupGlobalOptions() {
	globals := RootCmd.PersistentFlags()

	globals.StringVar(&globalOpts.MusecoreHome, tmcli.HomeFlag, app.DefaultNodeHome, "home path")
	// add more options here (e.g. verbosity, etc...)
}

func init() {
	// Setup options
	setupGlobalOptions()
	setupInitializeConfigOptions()
	setupRelayerOptions()

	// Define commands
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(StartCmd)
	RootCmd.AddCommand(InitializeConfigCmd)

	RootCmd.AddCommand(TSSCmd)
	TSSCmd.AddCommand(TSSEncryptCmd)
	TSSCmd.AddCommand(TSSGeneratePreParamsCmd)

	RootCmd.AddCommand(RelayerCmd)
	RelayerCmd.AddCommand(RelayerImportKeyCmd)
	RelayerCmd.AddCommand(RelayerShowAddressCmd)
}

func main() {
	ctx := context.Background()

	if err := RootCmd.ExecuteContext(ctx); err != nil {
		fmt.Printf("Error: %s. Exit code 1\n", err)
		os.Exit(1)
	}
}
