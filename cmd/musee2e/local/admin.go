package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/e2etests"
	"github.com/RWAs-labs/muse/e2e/runner"
)

// adminTestRoutine runs admin functions tests
func adminTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserAdmin
		// initialize runner for erc20 advanced test
		adminRunner, err := initTestRunner(
			"admin",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiGreen, "admin"),
			runner.WithMuseTxServer(deployerRunner.MuseTxServer),
		)
		if err != nil {
			return err
		}

		adminRunner.Logger.Print("🏃 starting admin tests")
		startTime := time.Now()

		// funding the account
		// we transfer around the total supply of Muse to the admin for the chain migration test
		txMuseSend := deployerRunner.LegacySendMuseOnEvm(account.EVMAddress(), 20_500_000_000)
		txERC20Send := deployerRunner.SendERC20OnEVM(account.EVMAddress(), 1000)
		adminRunner.WaitForTxReceiptOnEVM(txMuseSend)
		adminRunner.WaitForTxReceiptOnEVM(txERC20Send)

		// depositing the necessary tokens on MuseChain
		txMuseDeposit := adminRunner.LegacyDepositMuse()
		txEtherDeposit := adminRunner.DepositEtherDeployer()
		txERC20Deposit := adminRunner.DepositERC20Deployer()
		adminRunner.WaitForMinedCCTX(txMuseDeposit)
		adminRunner.WaitForMinedCCTX(txEtherDeposit)
		adminRunner.WaitForMinedCCTX(txERC20Deposit)

		// run erc20 advanced test
		testsToRun, err := adminRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("admin tests failed: %v", err)
		}

		if err := adminRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("admin tests failed: %v", err)
		}

		adminRunner.Logger.Print("🍾 admin tests completed in %s", time.Since(startTime).String())

		return err
	}
}
