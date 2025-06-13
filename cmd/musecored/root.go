package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	ethermintclient "github.com/RWAs-labs/ethermint/client"
	"github.com/RWAs-labs/ethermint/crypto/hd"
	evmenc "github.com/RWAs-labs/ethermint/encoding"
	"github.com/RWAs-labs/ethermint/types"
	tmcfg "github.com/cometbft/cometbft/config"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/app"
	musecoredconfig "github.com/RWAs-labs/muse/cmd/musecored/config"
	musemempool "github.com/RWAs-labs/muse/pkg/mempool"
	mevmserver "github.com/RWAs-labs/muse/server"
	servercfg "github.com/RWAs-labs/muse/server/config"
)

const EnvPrefix = "musecore"

// NewRootCmd creates a new root command for wasmd. It is called once in the
// main function.

type emptyAppOptions struct{}

func (ao emptyAppOptions) Get(_ string) interface{} { return nil }

func NewRootCmd() (*cobra.Command, types.EncodingConfig) {
	encodingConfig := app.MakeEncodingConfig()

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(app.DefaultNodeHome).
		WithKeyringOptions(hd.EthSecp256k1Option()).
		WithViper(EnvPrefix)

	rootCmd := &cobra.Command{
		Use:   musecoredconfig.AppName,
		Short: "Musecore Daemon (server)",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			// This needs to go after ReadFromClientConfig, as that function
			// sets the RPC client needed for SIGN_MODE_TEXTUAL. This sign mode
			// is only available if the client is online.
			if !initClientCtx.Offline {
				enabledSignModes := slices.Clone(tx.DefaultSignModes)
				enabledSignModes = append(enabledSignModes, signing.SignMode_SIGN_MODE_TEXTUAL)
				txConfigOpts := tx.ConfigOptions{
					EnabledSignModes:           enabledSignModes,
					TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(initClientCtx),
				}
				txConfig, err := tx.NewTxConfigWithOptions(
					initClientCtx.Codec,
					txConfigOpts,
				)
				if err != nil {
					return err
				}

				initClientCtx = initClientCtx.WithTxConfig(txConfig)
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, initTmConfig())
		},
	}

	initRootCmd(rootCmd, encodingConfig)
	rootCmd.AddCommand(
		confixcmd.ConfigCommand(),
	)

	// need to create this app instance to get autocliopts
	tempApp := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		map[int64]bool{},
		"",
		0,
		evmenc.MakeConfig(),
		emptyAppOptions{},
	)
	autoCliOpts := tempApp.AutoCliOpts()
	autoCliOpts.ClientCtx = initClientCtx

	if err := autoCliOpts.EnhanceRootCommand(rootCmd); err != nil {
		panic(err)
	}

	return rootCmd, encodingConfig
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	return servercfg.AppConfig(musecoredconfig.BaseDenom)
}

// initTmConfig overrides the default Tendermint config
func initTmConfig() *tmcfg.Config {
	cfg := tmcfg.DefaultConfig()

	cfg.DBBackend = "pebbledb"

	return cfg
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig types.EncodingConfig) {
	ac := appCreator{
		encCfg: encodingConfig,
	}

	rootCmd.AddCommand(
		ethermintclient.ValidateChainID(
			genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		),
		genutilcli.CollectGenTxsCmd(
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
			genutiltypes.DefaultMessageValidator,
			encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec(),
		),
		genutilcli.GenTxCmd(
			app.ModuleBasics,
			encodingConfig.TxConfig,
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
			encodingConfig.TxConfig.SigningContext().ValidatorAddressCodec(),
		),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		AddObserverListCmd(),
		CmdParseGenesisFile(),
		GetPubKeyCmd(),
		CollectObserverInfoCmd(),
		AddrConversionCmd(),
		UpgradeHandlerVersionCmd(),
		tmcli.NewCompletionCmd(rootCmd, true),
		ethermintclient.NewTestnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}),

		debug.Cmd(),
		snapshot.Cmd(ac.newApp),
	)

	mevmserver.AddCommands(
		rootCmd,
		mevmserver.NewDefaultStartOptions(ac.newApp, app.DefaultNodeHome),
		ac.appExport,
		addModuleInitFlags,
	)

	// the ethermintserver one supercedes the sdk one
	//server.AddCommands(rootCmd, app.DefaultNodeHome, ac.newApp, ac.createSimappAndExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		queryCommand(),
		txCommand(),
		docsCommand(),
		ethermintclient.KeyCommands(app.DefaultNodeHome),
	)

	// replace the default hd-path for the key add command with Ethereum HD Path
	if err := SetEthereumHDPath(rootCmd); err != nil {
		fmt.Printf("warning: unable to set default HD path: %v\n", err)
	}
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.ValidatorCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockCmd(),
		server.QueryBlocksCmd(),
		server.QueryBlockResultsCmd(),
	)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

type appCreator struct {
	encCfg types.EncodingConfig
}

func (ac appCreator) newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)
	maxTxs := cast.ToInt(appOpts.Get(server.FlagMempoolMaxTxs))
	if maxTxs <= 0 {
		maxTxs = musemempool.DefaultMaxTxs
	}
	baseappOptions = append(baseappOptions, func(app *baseapp.BaseApp) {
		app.SetMempool(musemempool.NewPriorityMempool(musemempool.PriorityNonceWithMaxTx(maxTxs)))
	})
	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	return app.New(logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		//cosmoscmd.EncodingConfig(ac.encCfg),
		ac.encCfg,
		appOpts,
		baseappOptions...,
	)
}

// appExport is used to export the state of the application for a genesis file.
func (ac appCreator) appExport(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var museApp *app.App

	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	loadLatest := false
	if height == -1 {
		loadLatest = true
	}

	museApp = app.New(
		logger,
		db,
		traceStore,
		loadLatest,
		map[int64]bool{},
		homePath,
		uint(1),
		ac.encCfg,
		appOpts,
	)

	// If height is -1, it means we are using the latest height.
	// For all other cases, we load the specified height from the Store
	if !loadLatest {
		if err := museApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return museApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
