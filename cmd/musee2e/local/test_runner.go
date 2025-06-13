package local

import (
	musee2econfig "github.com/RWAs-labs/muse/cmd/musee2e/config"
	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/runner"
)

// initTestRunner initializes a runner form tests
// it creates a runner with an account and copy contracts from deployer runner
func initTestRunner(
	name string,
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	account config.Account,
	logger *runner.Logger,
	opts ...runner.E2ERunnerOption,
) (*runner.E2ERunner, error) {
	// initialize runner for test
	testRunner, err := musee2econfig.RunnerFromConfig(
		deployerRunner.Ctx,
		name,
		deployerRunner.CtxCancel,
		conf,
		account,
		logger,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	// copy timeouts from deployer runner
	testRunner.CctxTimeout = deployerRunner.ReceiptTimeout
	testRunner.ReceiptTimeout = deployerRunner.ReceiptTimeout
	testRunner.TestFilter = deployerRunner.TestFilter

	// copy contracts from deployer runner
	if err := testRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}

	return testRunner, nil
}
