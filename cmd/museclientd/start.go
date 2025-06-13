package main

import (
	"context"
	"net/http"
	_ "net/http/pprof" // #nosec G108 -- pprof enablement is intentional
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/config"
	zctx "github.com/RWAs-labs/muse/museclient/context"
	"github.com/RWAs-labs/muse/museclient/maintenance"
	"github.com/RWAs-labs/muse/museclient/metrics"
	"github.com/RWAs-labs/muse/museclient/musecore"
	"github.com/RWAs-labs/muse/museclient/orchestrator"
	musetss "github.com/RWAs-labs/muse/museclient/tss"
	"github.com/RWAs-labs/muse/pkg/bg"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/constant"
	"github.com/RWAs-labs/muse/pkg/graceful"
	museos "github.com/RWAs-labs/muse/pkg/os"
	"github.com/RWAs-labs/muse/pkg/scheduler"
)

const (
	// enables posting blame data to core for failed TSS signatures
	envFlagPostBlame = "POST_BLAME"
	envPprofAddr     = "PPROF_ADDR"
)

// Start starts museclientd process
func Start(_ *cobra.Command, _ []string) error {
	// Load Config file given path
	cfg, err := config.Load(globalOpts.MusecoreHome)
	if err != nil {
		return errors.Wrap(err, "unable to load config")
	}

	dbPath, err := config.ResolveDBPath()
	if err != nil {
		return errors.Wrap(err, "unable to resolve db path")
	}

	// Configure logger (also overrides the default log level)
	logger, err := base.NewLogger(cfg)
	if err != nil {
		return errors.Wrap(err, "unable to create logger")
	}

	passes, err := promptPasswords()
	if err != nil {
		return errors.Wrap(err, "unable to prompt for passwords")
	}

	appContext := zctx.New(cfg, passes.relayerKeys(), logger.Std)
	ctx := zctx.WithAppContext(context.Background(), appContext)

	err = config.LoadRestrictedAddressesConfig(cfg, globalOpts.MusecoreHome)
	if err != nil {
		logger.Std.Err(err).Msg("loading restricted addresses config")
	} else {
		bg.Work(ctx, func(ctx context.Context) error {
			return config.WatchRestrictedAddressesConfig(ctx, cfg, globalOpts.MusecoreHome, logger.Std)
		}, bg.WithName("watch_restricted_addresses_config"), bg.WithLogger(logger.Std))
	}

	telemetry, err := startTelemetry(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "unable to start telemetry")
	}

	// musecore client is used for all communication to muse node.
	// it accumulates votes, and provides a source of truth for all clients
	//
	// This call crated client, ensured block production, then prepares the client
	musecoreClient, err := musecore.NewFromConfig(ctx, &cfg, passes.hotkey, logger.Std)
	if err != nil {
		return errors.Wrap(err, "unable to create musecore client from config")
	}

	// Initialize core parameters from musecore
	if err = orchestrator.UpdateAppContext(ctx, appContext, musecoreClient, logger.Std); err != nil {
		return errors.Wrap(err, "unable to update app context")
	}

	log.Debug().Msgf("Config is updated from musecore\n %s", cfg.StringMasked())

	granteePubKeyBech32, err := resolveObserverPubKeyBech32(cfg, passes.hotkey)
	if err != nil {
		return errors.Wrap(err, "unable to resolve observer pub key bech32")
	}

	isObserver, err := isObserverNode(ctx, musecoreClient)
	switch {
	case err != nil:
		return errors.Wrap(err, "unable to check if observer node")
	case !isObserver:
		logger.Std.Warn().Msg("This node is not an observer node. Exit 0")
		return nil
	}

	shutdownListener := maintenance.NewShutdownListener(musecoreClient, logger.Std)
	if err := shutdownListener.RunPreStartCheck(ctx); err != nil {
		return errors.Wrap(err, "pre start check failed")
	}

	tssSetupProps := musetss.SetupProps{
		Config:              cfg,
		Musecore:            musecoreClient,
		GranteePubKeyBech32: granteePubKeyBech32,
		HotKeyPassword:      passes.hotkey,
		TSSKeyPassword:      passes.tss,
		BitcoinChainIDs:     btcChainIDsFromContext(appContext),
		PostBlame:           isEnvFlagEnabled(envFlagPostBlame),
		Telemetry:           telemetry,
	}

	// This will start p2p communication so it should only happen after
	// preflight checks have completed
	tss, err := musetss.Setup(ctx, tssSetupProps, logger.Std)
	if err != nil {
		return errors.Wrap(err, "unable to setup TSS service")
	}

	graceful.AddStopper(tss.Stop)

	// Starts various background TSS listeners.
	// Shuts down museclientd if any is triggered.
	maintenance.NewTSSListener(musecoreClient, logger.Std).Listen(ctx, func() {
		logger.Std.Info().Msg("TSS listener received an action to shutdown museclientd.")
		graceful.ShutdownNow()
	})

	shutdownListener.Listen(ctx, func() {
		logger.Std.Info().Msg("Shutdown listener received an action to shutdown museclientd.")
		graceful.ShutdownNow()
	})

	// Orchestrator wraps the musecore client and adds the observers and signer maps to it.
	// This is the high level object used for CCTX interactions
	// It also handles background configuration updates from musecore
	taskScheduler := scheduler.New(logger.Std, 0)
	maestroDeps := &orchestrator.Dependencies{
		Musecore:  musecoreClient,
		TSS:       tss,
		DBPath:    dbPath,
		Telemetry: telemetry,
	}

	maestro, err := orchestrator.New(taskScheduler, maestroDeps, logger)
	if err != nil {
		return errors.Wrap(err, "unable to create orchestrator")
	}

	// Start orchestrator with all observers and signers
	graceful.AddService(ctx, maestro)

	// Block current routine until a shutdown signal is received
	graceful.WaitForShutdown()

	return nil
}

type passwords struct {
	hotkey           string
	tss              string
	solanaRelayerKey string
}

// promptPasswords prompts for Hotkey, TSS key-share and relayer key passwords
func promptPasswords() (passwords, error) {
	titles := []string{"HotKey", "TSS", "Solana Relayer Key"}

	res, err := museos.PromptPasswords(titles)
	if err != nil {
		return passwords{}, errors.Wrap(err, "unable to get passwords")
	}

	return passwords{
		hotkey:           res[0],
		tss:              res[1],
		solanaRelayerKey: res[2],
	}, nil
}

func (p passwords) relayerKeys() map[string]string {
	return map[string]string{
		chains.Network_solana.String(): p.solanaRelayerKey,
	}
}

func startTelemetry(ctx context.Context, cfg config.Config) (*metrics.TelemetryServer, error) {
	// 1. Init pprof http server
	pprofServer := func(_ context.Context) error {
		addr := os.Getenv(envPprofAddr)
		if addr == "" {
			addr = "localhost:6061"
		}

		log.Info().Str("addr", addr).Msg("starting pprof http server")

		// #nosec G114 -- timeouts unneeded
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			log.Error().Err(err).Msg("pprof http server error")
		}

		return nil
	}

	// 2. Init metrics server
	metricsServer, err := metrics.NewMetrics()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create metrics")
	}

	metrics.Info.WithLabelValues(constant.Version).Set(1)
	metrics.LastStartTime.SetToCurrentTime()

	// 3. Init telemetry server
	telemetry := metrics.NewTelemetryServer()
	telemetry.SetIPAddress(cfg.PublicIP)

	// 4. Add services to the process
	graceful.AddStarter(ctx, pprofServer)
	graceful.AddService(ctx, metricsServer)
	graceful.AddService(ctx, telemetry)

	return telemetry, nil
}
