package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/e2etests"
	"github.com/RWAs-labs/muse/e2e/runner"
)

// legacyERC20TestRoutine runs erc20 related e2e tests
func legacyERC20TestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserLegacyERC20
		// initialize runner for erc20 test
		erc20Runner, err := initTestRunner(
			"erc20",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgGreen, "erc20"),
			runner.WithMuseTxServer(deployerRunner.MuseTxServer),
		)
		if err != nil {
			return err
		}

		erc20Runner.Logger.Print("🏃 starting erc20 tests")
		startTime := time.Now()

		// funding the account
		txERC20Send := deployerRunner.SendERC20OnEVM(account.EVMAddress(), 10000)
		erc20Runner.WaitForTxReceiptOnEVM(txERC20Send)

		// depositing the necessary tokens on MuseChain
		txEtherDeposit := erc20Runner.LegacyDepositEther()
		txERC20Deposit := erc20Runner.LegacyDepositERC20()
		erc20Runner.WaitForMinedCCTX(txEtherDeposit)
		erc20Runner.WaitForMinedCCTX(txERC20Deposit)

		// run erc20 test
		testsToRun, err := erc20Runner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("erc20 tests failed: %v", err)
		}

		if err := erc20Runner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("erc20 tests failed: %v", err)
		}

		erc20Runner.Logger.Print("🍾 erc20 tests completed in %s", time.Since(startTime).String())

		return err
	}
}

// legacyEthereumTestRoutine runs Ethereum related e2e tests
func legacyEthereumTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		// initialize runner for ether test
		ethereumRunner, err := initTestRunner(
			"ether",
			conf,
			deployerRunner,
			conf.AdditionalAccounts.UserLegacyEther,
			runner.NewLogger(verbose, color.FgMagenta, "ether"),
		)
		if err != nil {
			return err
		}

		ethereumRunner.Logger.Print("🏃 starting Ethereum tests")
		startTime := time.Now()

		// depositing the necessary tokens on MuseChain
		txEtherDeposit := ethereumRunner.LegacyDepositEther()
		ethereumRunner.WaitForMinedCCTX(txEtherDeposit)

		// run ethereum test
		// Note: due to the extensive block generation in Ethereum localnet, block header test is run first
		// to make it faster to catch up with the latest block header
		testsToRun, err := ethereumRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("ethereum tests failed: %v", err)
		}

		if err := ethereumRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("ethereum tests failed: %v", err)
		}

		ethereumRunner.Logger.Print("🍾 Ethereum tests completed in %s", time.Since(startTime).String())

		return err
	}
}

// legacyMEVMMPTestRoutine runs MEVM message passing related e2e tests
func legacyMEVMMPTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserLegacyMEVMMP
		// initialize runner for mevm mp test
		mevmMPRunner, err := initTestRunner(
			"mevm_mp",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiRed, "mevm_mp"),
		)
		if err != nil {
			return err
		}

		mevmMPRunner.Logger.Print("🏃 starting MEVM Message Passing tests")
		startTime := time.Now()

		// funding the account
		txMuseSend := deployerRunner.LegacySendMuseOnEvm(account.EVMAddress(), 1000)
		mevmMPRunner.WaitForTxReceiptOnEVM(txMuseSend)

		// depositing the necessary tokens on MuseChain
		txMuseDeposit := mevmMPRunner.LegacyDepositMuse()
		txEtherDeposit := mevmMPRunner.LegacyDepositEther()
		mevmMPRunner.WaitForMinedCCTX(txMuseDeposit)
		mevmMPRunner.WaitForMinedCCTX(txEtherDeposit)

		// run mevm message passing test
		testsToRun, err := mevmMPRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("mevm message passing tests failed: %v", err)
		}

		if err := mevmMPRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("mevm message passing tests failed: %v", err)
		}

		mevmMPRunner.Logger.Print("🍾 MEVM message passing tests completed in %s", time.Since(startTime).String())

		return err
	}
}

// legacyMUSETestRoutine runs Muse transfer and message passing related e2e tests
func legacyMUSETestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserLegacyMuse
		// initialize runner for muse test
		museRunner, err := initTestRunner(
			"muse",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgBlue, "muse"),
		)
		if err != nil {
			return err
		}

		museRunner.Logger.Print("🏃 starting Muse tests")
		startTime := time.Now()

		// funding the account
		txMuseSend := deployerRunner.LegacySendMuseOnEvm(account.EVMAddress(), 1000)
		museRunner.WaitForTxReceiptOnEVM(txMuseSend)

		// depositing the necessary tokens on MuseChain
		txMuseDeposit := museRunner.LegacyDepositMuse()
		txEtherDeposit := museRunner.LegacyDepositEther()
		museRunner.WaitForMinedCCTX(txMuseDeposit)
		museRunner.WaitForMinedCCTX(txEtherDeposit)

		// run muse test
		testsToRun, err := museRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("muse tests failed: %v", err)
		}

		if err := museRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("muse tests failed: %v", err)
		}

		museRunner.Logger.Print("🍾 Muse tests completed in %s", time.Since(startTime).String())

		return err
	}
}
