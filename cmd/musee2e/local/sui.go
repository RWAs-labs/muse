package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/e2etests"
	"github.com/RWAs-labs/muse/e2e/runner"
)

// suiTestRoutine runs Sui related e2e tests
func suiTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		// initialize runner for sui test
		suiRunner, err := initTestRunner(
			"sui",
			conf,
			deployerRunner,
			conf.AdditionalAccounts.UserSui,
			runner.NewLogger(verbose, color.FgHiCyan, "sui"),
			runner.WithMuseTxServer(deployerRunner.MuseTxServer),
		)
		if err != nil {
			return err
		}

		suiRunner.Logger.Print("🏃 starting Sui tests")
		startTime := time.Now()

		suiRunnerSigner, err := suiRunner.Account.SuiSigner()
		if err != nil {
			return err
		}

		// get tokens for the account
		suiRunner.RequestSuiFromFaucet(conf.RPCs.SuiFaucet, suiRunnerSigner.Address())

		// mint fungible tokens to the account
		txRes := deployerRunner.SuiMintUSDC("100000000000", suiRunnerSigner.Address())

		deployerRunner.Logger.Info("Sui USDC mint tx: %s", txRes.Digest)

		// run sui test
		testsToRun, err := suiRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("sui tests failed: %v", err)
		}

		if err := suiRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("sui tests failed: %v", err)
		}

		suiRunner.Logger.Print("🍾 sui tests completed in %s", time.Since(startTime).String())

		return err
	}
}
