package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/e2etests"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/pkg/errgroup"
)

// startEVMTests starts EVM chains related tests in parallel
func startEVMTests(eg *errgroup.Group, conf config.Config, deployerRunner *runner.E2ERunner, verbose bool) {
	// Test happy paths for gas token workflow
	eg.Go(evmTestRoutine(conf, "eth", conf.AdditionalAccounts.UserEther, color.FgHiGreen, deployerRunner, verbose,
		e2etests.TestETHDepositName,
		e2etests.TestETHDepositAndCallName,
		e2etests.TestETHDepositFastConfirmationName,
		e2etests.TestETHWithdrawName,
		e2etests.TestETHWithdrawAndArbitraryCallName,
		e2etests.TestETHWithdrawAndCallName,
		e2etests.TestETHWithdrawAndCallThroughContractName,
		e2etests.TestMEVMToEVMArbitraryCallName,
		e2etests.TestMEVMToEVMCallName,
		e2etests.TestMEVMToEVMCallThroughContractName,
		e2etests.TestEVMToMEVMCallName,
		e2etests.TestETHDepositAndCallNoMessageName,
		e2etests.TestETHWithdrawAndCallNoMessageName,
		e2etests.TestEtherWithdrawRestrictedName,
	))

	// Test happy paths for erc20 token workflow
	eg.Go(evmTestRoutine(conf, "erc20", conf.AdditionalAccounts.UserERC20, color.FgHiBlue, deployerRunner, verbose,
		e2etests.TestETHDepositName, // necessary to pay fees on MEVM
		e2etests.TestERC20DepositName,
		e2etests.TestERC20DepositAndCallName,
		e2etests.TestERC20WithdrawName,
		e2etests.TestERC20WithdrawAndArbitraryCallName,
		e2etests.TestERC20WithdrawAndCallName,
		e2etests.TestERC20DepositAndCallNoMessageName,
		e2etests.TestERC20WithdrawAndCallNoMessageName,
		e2etests.TestDepositAndCallSwapName,
		e2etests.TestERC20DepositRestrictedName,
	))

	// Test revert cases for gas token workflow
	eg.Go(
		evmTestRoutine(
			conf,
			"eth-revert",
			conf.AdditionalAccounts.UserEtherRevert,
			color.FgHiYellow,
			deployerRunner,
			verbose,
			e2etests.TestETHDepositName, // necessary to pay fees on MEVM and withdraw
			e2etests.TestETHDepositAndCallRevertName,
			e2etests.TestETHDepositAndCallRevertWithCallName,
			e2etests.TestETHDepositRevertAndAbortName,
			e2etests.TestETHWithdrawAndCallRevertName,
			e2etests.TestETHWithdrawAndCallRevertWithCallName,
			e2etests.TestETHWithdrawRevertAndAbortName,
			e2etests.TestETHWithdrawAndCallRevertWithWithdrawName,
			e2etests.TestDepositAndCallOutOfGasName,
			e2etests.TestMEVMToEVMCallRevertName,
			e2etests.TestMEVMToEVMCallRevertAndAbortName,
			e2etests.TestEVMToMEVMCallAbortName,
		),
	)

	// Test revert cases for erc20 token workflow
	eg.Go(
		evmTestRoutine(
			conf,
			"erc20-revert",
			conf.AdditionalAccounts.UserERC20Revert,
			color.FgHiRed,
			deployerRunner,
			verbose,
			e2etests.TestETHDepositName,   // necessary to pay fees on MEVM
			e2etests.TestERC20DepositName, // necessary to have assets to withdraw
			e2etests.TestOperationAddLiquidityETHName, // liquidity with gas and ERC20 are necessary for reverts
			e2etests.TestOperationAddLiquidityERC20Name,
			e2etests.TestERC20DepositAndCallRevertName,
			e2etests.TestERC20DepositAndCallRevertWithCallName,
			e2etests.TestERC20DepositRevertAndAbortName,
			e2etests.TestERC20WithdrawAndCallRevertName,
			e2etests.TestERC20WithdrawAndCallRevertWithCallName,
			e2etests.TestERC20WithdrawRevertAndAbortName,
		),
	)
}

// evmTestRoutine runs EVM chain related e2e tests
func evmTestRoutine(
	conf config.Config,
	name string,
	account config.Account,
	color color.Attribute,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		name = "v2-" + name

		// initialize runner for erc20 test
		v2Runner, err := initTestRunner(
			name,
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color, name),
			runner.WithMuseTxServer(deployerRunner.MuseTxServer),
		)
		if err != nil {
			return err
		}

		v2Runner.Logger.Print("🏃 starting %s tests", name)
		startTime := time.Now()

		// funding the account
		txERC20Send := deployerRunner.SendERC20OnEVM(account.EVMAddress(), 10000)
		v2Runner.WaitForTxReceiptOnEVM(txERC20Send)

		// run erc20 test
		testsToRun, err := v2Runner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("%s tests failed: %v", name, err)
		}

		if err := v2Runner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("%s tests failed: %v", name, err)
		}

		v2Runner.Logger.Print("🍾 %s tests completed in %s", name, time.Since(startTime).String())

		return err
	}
}
