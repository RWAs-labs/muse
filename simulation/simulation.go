package simulation

import (
	"encoding/json"
	"fmt"
	"os"

	"cosmossdk.io/log"
	"github.com/RWAs-labs/ethermint/app"
	evmante "github.com/RWAs-labs/ethermint/app/ante"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	museapp "github.com/RWAs-labs/muse/app"
	"github.com/RWAs-labs/muse/app/ante"
)

func NewSimApp(
	logger log.Logger,
	db dbm.DB,
	appOptions servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) (*museapp.App, error) {
	encCdc := museapp.MakeEncodingConfig()

	// Set load latest version to false as we manually set it later.
	museApp := museapp.New(
		logger,
		db,
		nil,
		false,
		map[int64]bool{},
		app.DefaultNodeHome,
		5,
		encCdc,
		appOptions,
		baseAppOptions...,
	)

	// use muse antehandler
	options := ante.HandlerOptions{
		AccountKeeper:   museApp.AccountKeeper,
		BankKeeper:      museApp.BankKeeper,
		EvmKeeper:       museApp.EvmKeeper,
		FeeMarketKeeper: museApp.FeeMarketKeeper,
		SignModeHandler: encCdc.TxConfig.SignModeHandler(),
		SigGasConsumer:  evmante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:  0,
		ObserverKeeper:  museApp.ObserverKeeper,
	}

	anteHandler, err := ante.NewAnteHandler(options)
	if err != nil {
		panic(err)
	}

	museApp.SetAnteHandler(anteHandler)
	if err := museApp.LoadLatestVersion(); err != nil {
		return nil, err
	}
	return museApp, nil
}

// PrintStats prints the corresponding statistics from the app DB.
func PrintStats(db dbm.DB) {
	fmt.Println("\nDB Stats")
	fmt.Println(db.Stats()["leveldb.stats"])
	fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
}

// CheckExportSimulation exports the app state and simulation parameters to JSON
// if the export paths are defined.
func CheckExportSimulation(app runtime.AppI, config simtypes.Config, params simtypes.Params) error {
	if config.ExportStatePath != "" {
		exported, err := app.ExportAppStateAndValidators(false, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to export app state: %w", err)
		}

		if err := os.WriteFile(config.ExportStatePath, exported.AppState, 0o600); err != nil {
			return err
		}
	}

	if config.ExportParamsPath != "" {
		paramsBz, err := json.MarshalIndent(params, "", " ")
		if err != nil {
			return fmt.Errorf("failed to write app state to %s: %w", config.ExportStatePath, err)
		}

		if err := os.WriteFile(config.ExportParamsPath, paramsBz, 0o600); err != nil {
			return err
		}
	}
	return nil
}
