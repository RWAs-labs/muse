package local

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/e2etests"
	"github.com/RWAs-labs/muse/e2e/runner"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// tssMigrationTestRoutine runs TSS migration related e2e tests
func tssMigrationTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserMigration
		// initialize runner for migration test
		tssMigrationTestRunner, err := initTestRunner(
			"triggerTSSMigration",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiGreen, "migration"),
			runner.WithMuseTxServer(deployerRunner.MuseTxServer),
		)
		if err != nil {
			return err
		}

		tssMigrationTestRunner.Logger.Print("🏃 starting TSS migration tests")
		startTime := time.Now()

		if len(testNames) == 0 {
			tssMigrationTestRunner.Logger.Print("🍾 TSS migration tests completed in %s", time.Since(startTime).String())
			return nil
		}
		// run TSS migration test
		testsToRun, err := tssMigrationTestRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("TSS migration tests failed: %v", err)
		}

		if err := tssMigrationTestRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("TSS migration tests failed: %v", err)
		}
		if err := tssMigrationTestRunner.CheckBTCTSSBalance(); err != nil {
			return err
		}

		tssMigrationTestRunner.Logger.Print("🍾 TSS migration tests completed in %s", time.Since(startTime).String())

		return nil
	}
}

func triggerTSSMigration(deployerRunner *runner.E2ERunner, logger *runner.Logger, verbose bool, conf config.Config) {
	migrationStartTime := time.Now()
	logger.Print("🏁 starting tss migration")

	response, err := deployerRunner.CctxClient.LastMuseHeight(
		deployerRunner.Ctx,
		&crosschaintypes.QueryLastMuseHeightRequest{},
	)
	require.NoError(deployerRunner, err)
	err = deployerRunner.MuseTxServer.UpdateKeygen(response.Height)
	require.NoError(deployerRunner, err)

	// Generate new TSS
	noError(waitKeygenHeight(deployerRunner.Ctx, deployerRunner.CctxClient, deployerRunner.ObserverClient, logger, 0))

	// Run migration
	// migrationRoutine runs migration e2e test , which migrates funds from the older TSS to the new one
	// The museclient restarts required for this process are managed by the background workers in museclient (TSSListener)
	fn := tssMigrationTestRoutine(conf, deployerRunner, verbose, e2etests.TestMigrateTSSName)

	if err := fn(); err != nil {
		logger.Print("❌ %v", err)
		logger.Print("❌ tss migration failed")
		os.Exit(1)
	}

	// Update TSS address for contracts in connected chains
	// TODO : Update TSS address for other chains if necessary
	// https://github.com/RWAs-labs/muse/issues/3599
	deployerRunner.UpdateTSSAddressForConnector()
	deployerRunner.UpdateTSSAddressForERC20custody()
	deployerRunner.UpdateTSSAddressForGateway()
	deployerRunner.UpdateTSSAddressSolana(
		conf.Contracts.Solana.GatewayProgramID.String(),
		conf.AdditionalAccounts.UserSolana.SolanaPrivateKey.String())
	logger.Print("✅ migration completed in %s ", time.Since(migrationStartTime).String())
}
