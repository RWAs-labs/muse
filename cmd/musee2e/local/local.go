package local

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	musee2econfig "github.com/RWAs-labs/muse/cmd/musee2e/config"
	"github.com/RWAs-labs/muse/e2e/config"
	"github.com/RWAs-labs/muse/e2e/e2etests"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/txserver"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/errgroup"
	"github.com/RWAs-labs/muse/testutil"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

const (
	flagContractsDeployed = "deployed"
	flagWaitForHeight     = "wait-for"
	FlagConfigFile        = "config"
	flagConfigOut         = "config-out"
	flagVerbose           = "verbose"
	flagTestAdmin         = "test-admin"
	flagTestEthStress     = "test-stress-eth"
	flagTestSolanaStress  = "test-stress-solana"
	flagTestSuiStress     = "test-stress-sui"
	flagIterations        = "iterations"
	flagTestSolana        = "test-solana"
	flagTestTON           = "test-ton"
	flagTestSui           = "test-sui"
	flagSkipRegular       = "skip-regular"
	flagLight             = "light"
	flagSetupOnly         = "setup-only"
	flagSkipSetup         = "skip-setup"
	flagTestTSSMigration  = "test-tss-migration"
	flagSkipBitcoinSetup  = "skip-bitcoin-setup"
	flagSkipHeaderProof   = "skip-header-proof"
	flagTestLegacy        = "test-legacy"
	flagSkipTrackerCheck  = "skip-tracker-check"
	flagSkipPrecompiles   = "skip-precompiles"
	flagUpgradeContracts  = "upgrade-contracts"
	flagTestFilter        = "test-filter"
	flagTestStaking       = "test-staking"
)

var (
	TestTimeout        = 20 * time.Minute
	ErrTopLevelTimeout = errors.New("top level test timeout")
	noError            = testutil.NoError
)

// NewLocalCmd returns the local command
// which runs the E2E tests locally on the machine with localnet for each blockchain
func NewLocalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Run Local E2E tests",
		Run:   localE2ETest,
	}
	cmd.Flags().Bool(flagContractsDeployed, false, "set to to true if running tests again with existing state")
	cmd.Flags().Int64(flagWaitForHeight, 1, "block height for tests to begin, ex. --wait-for 100")
	cmd.Flags().String(FlagConfigFile, "", "config file to use for the tests")
	cmd.Flags().Bool(flagVerbose, false, "set to true to enable verbose logging")
	cmd.Flags().Bool(flagTestAdmin, false, "set to true to run admin tests")
	cmd.Flags().Bool(flagTestEthStress, false, "set to true to run eth stress tests")
	cmd.Flags().Bool(flagTestSolanaStress, false, "set to true to run solana stress tests")
	cmd.Flags().Bool(flagTestSuiStress, false, "set to true to run sui stress tests")
	cmd.Flags().Int(flagIterations, 100, "number of iterations to run each performance test")
	cmd.Flags().Bool(flagTestSolana, false, "set to true to run solana tests")
	cmd.Flags().Bool(flagTestTON, false, "set to true to run TON tests")
	cmd.Flags().Bool(flagTestSui, false, "set to true to run Sui tests")
	cmd.Flags().Bool(flagSkipRegular, false, "set to true to skip regular tests")
	cmd.Flags().Bool(flagLight, false, "run the most basic regular tests, useful for quick checks")
	cmd.Flags().Bool(flagSetupOnly, false, "set to true to only setup the networks")
	cmd.Flags().String(flagConfigOut, "", "config file to write the deployed contracts from the setup")
	cmd.Flags().Bool(flagSkipSetup, false, "set to true to skip setup")
	cmd.Flags().Bool(flagSkipBitcoinSetup, false, "set to true to skip bitcoin wallet setup")
	cmd.Flags().Bool(flagSkipHeaderProof, false, "set to true to skip header proof tests")
	cmd.Flags().Bool(flagTestTSSMigration, false, "set to true to include a migration test at the end")
	cmd.Flags().Bool(flagTestLegacy, false, "set to true to run legacy EVM tests")
	cmd.Flags().Bool(flagSkipTrackerCheck, false, "set to true to skip tracker check at the end of the tests")
	cmd.Flags().Bool(flagSkipPrecompiles, true, "set to true to skip stateful precompiled contracts test")
	cmd.Flags().
		Bool(flagUpgradeContracts, false, "set to true to upgrade Gateways and ERC20Custody contracts during setup for MEVM and EVM")
	cmd.Flags().String(flagTestFilter, "", "regexp filter to limit which test to run")
	cmd.Flags().Bool(flagTestStaking, false, "set to true to run staking tests")

	cmd.AddCommand(NewGetMuseclientBootstrap())

	return cmd
}

// TODO: simplify this file: put the different type of tests in separate files
// https://github.com/RWAs-labs/muse/issues/2762
func localE2ETest(cmd *cobra.Command, _ []string) {
	// fetch flags
	var (
		waitForHeight     = must(cmd.Flags().GetInt64(flagWaitForHeight))
		contractsDeployed = must(cmd.Flags().GetBool(flagContractsDeployed))
		verbose           = must(cmd.Flags().GetBool(flagVerbose))
		configOut         = must(cmd.Flags().GetString(flagConfigOut))
		testAdmin         = must(cmd.Flags().GetBool(flagTestAdmin))
		testEthStress     = must(cmd.Flags().GetBool(flagTestEthStress))
		testSolanaStress  = must(cmd.Flags().GetBool(flagTestSolanaStress))
		testSuiStress     = must(cmd.Flags().GetBool(flagTestSuiStress))
		iterations        = must(cmd.Flags().GetInt(flagIterations))
		testSolana        = must(cmd.Flags().GetBool(flagTestSolana))
		testTON           = must(cmd.Flags().GetBool(flagTestTON))
		testSui           = must(cmd.Flags().GetBool(flagTestSui))
		skipRegular       = must(cmd.Flags().GetBool(flagSkipRegular))
		light             = must(cmd.Flags().GetBool(flagLight))
		setupOnly         = must(cmd.Flags().GetBool(flagSetupOnly))
		skipSetup         = must(cmd.Flags().GetBool(flagSkipSetup))
		skipBitcoinSetup  = must(cmd.Flags().GetBool(flagSkipBitcoinSetup))
		skipHeaderProof   = must(cmd.Flags().GetBool(flagSkipHeaderProof))
		skipTrackerCheck  = must(cmd.Flags().GetBool(flagSkipTrackerCheck))
		testTSSMigration  = must(cmd.Flags().GetBool(flagTestTSSMigration))
		testLegacy        = must(cmd.Flags().GetBool(flagTestLegacy))
		skipPrecompiles   = must(cmd.Flags().GetBool(flagSkipPrecompiles))
		upgradeContracts  = must(cmd.Flags().GetBool(flagUpgradeContracts))
		testStress        = testEthStress || testSolanaStress || testSuiStress
		setupSolana       = testSolana || testStress
		setupSui          = testSui || testStress
		testFilterStr     = must(cmd.Flags().GetString(flagTestFilter))
		testStaking       = must(cmd.Flags().GetBool(flagTestStaking))
	)

	testFilter := regexp.MustCompile(testFilterStr)

	logger := runner.NewLogger(verbose, color.FgWhite, "setup")

	testStartTime := time.Now()

	logger.Print("starting E2E tests")

	if verbose {
		logger.Info("Flags")
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			logger.Info(`--%s="%s"`, flag.Name, flag.Value.String())
		})
	}

	if testAdmin {
		logger.Print("⚠️ admin tests enabled")
	}

	// skip regular tests if stress tests are enabled
	if testStress {
		logger.Print("⚠️ performance tests enabled, regular tests will be skipped")
		skipRegular = true
		skipPrecompiles = true

		if iterations > 100 {
			TestTimeout = time.Hour
		}
	}

	// initialize tests config
	conf, err := GetConfig(cmd)
	noError(err)

	// initialize context
	ctx, timeoutCancel := context.WithTimeoutCause(context.Background(), TestTimeout, ErrTopLevelTimeout)
	defer timeoutCancel()
	ctx, cancel := context.WithCancelCause(ctx)

	// route os signals to context cancellation.
	// using NotifyContext will ensure that the second signal
	// will not be handled and should kill the process.
	go func() {
		notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		<-notifyCtx.Done()
		cancel(fmt.Errorf("notify context: %w", notifyCtx.Err()))
		stop()
	}()

	// wait for a specific height on MuseChain
	noError(utils.WaitForBlockHeight(ctx, waitForHeight, conf.RPCs.MuseCoreRPC, logger))

	museTxServer, err := txserver.NewMuseTxServer(
		conf.RPCs.MuseCoreRPC,
		[]string{
			utils.EmergencyPolicyName,
			utils.OperationalPolicyName,
			utils.AdminPolicyName,
			utils.UserEmissionsWithdrawName,
		},
		[]string{
			conf.PolicyAccounts.EmergencyPolicyAccount.RawPrivateKey.String(),
			conf.PolicyAccounts.OperationalPolicyAccount.RawPrivateKey.String(),
			conf.PolicyAccounts.AdminPolicyAccount.RawPrivateKey.String(),
			conf.AdditionalAccounts.UserEmissionsWithdraw.RawPrivateKey.String(),
		},
		conf.MuseChainID,
	)
	noError(err)

	// Drop this cond after TON e2e is included in the default suite
	if !testTON {
		conf.RPCs.TON = ""
	}

	// initialize deployer runner with config
	deployerRunner, err := musee2econfig.RunnerFromConfig(
		ctx,
		"deployer",
		cancel,
		conf,
		conf.DefaultAccount,
		logger,
		runner.WithMuseTxServer(museTxServer),
		runner.WithTestFilter(testFilter),
	)
	noError(err)

	// monitor block production to ensure we fail fast if there are consensus failures
	go monitorBlockProductionCancel(ctx, cancel, conf)

	// set the authority client to the muse tx server to be able to query message permissions
	deployerRunner.MuseTxServer.SetAuthorityClient(deployerRunner.AuthorityClient)
	if !skipSetup {
		// run setup steps that do not require tss
		noError(deployerRunner.FundEmissionsPool())

		// wait for keygen to be completed
		// if setup is skipped, we assume that the keygen is already completed
		noError(waitKeygenHeight(ctx, deployerRunner.CctxClient, deployerRunner.ObserverClient, logger, 10))
	}

	// query and set the TSS
	noError(deployerRunner.SetTSSAddresses())

	if !skipHeaderProof {
		noError(deployerRunner.EnableHeaderVerification([]int64{
			chains.GoerliLocalnet.ChainId,
			chains.BitcoinRegtest.ChainId,
		}))
	}

	e2eStartHeight, err := deployerRunner.Clients.Musecore.GetBlockHeight(ctx)
	noError(err)

	// setting up the networks
	if !skipSetup {
		logger.Print("⚙️ setting up networks")
		startTime := time.Now()

		// setup TSS address and setup deployer wallet
		deployerRunner.SetupBitcoinAccounts(true)

		//setup protocol contracts v1 as they are still supported for now
		deployerRunner.LegacySetupEVM(contractsDeployed, testLegacy)

		// setup protocol contracts on the connected EVM chain
		deployerRunner.SetupEVM()

		if setupSolana {
			deployerRunner.SetupSolana(
				conf.Contracts.Solana.GatewayProgramID.String(),
				conf.AdditionalAccounts.UserSolana.SolanaPrivateKey.String(),
				conf.AdditionalAccounts.UserSPL.SolanaPrivateKey.String(),
			)
		}

		deployerRunner.SetupMEVMProtocolContracts()
		deployerRunner.SetupLegacyMEVMContracts()

		mrc20Deployment := txserver.MRC20Deployment{
			ERC20Addr: deployerRunner.ERC20Addr,
			SPLAddr:   nil,
		}
		if setupSolana {
			mrc20Deployment.SPLAddr = deployerRunner.SPLAddr.ToPointer()
		}
		deployerRunner.SetupMEVMMRC20s(mrc20Deployment)

		// Update the chain params to contains protocol contract addresses
		deployerRunner.UpdateProtocolContractsInChainParams()

		if testTON {
			deployerRunner.SetupTON(
				conf.RPCs.TONFaucet,
				conf.AdditionalAccounts.UserTON,
			)
		}

		if setupSui {
			deployerRunner.SetupSui(conf.RPCs.SuiFaucet)
		}

		logger.Print("✅ setup completed in %s", time.Since(startTime))
	}

	// if a config output is specified, write the config
	if configOut != "" {
		newConfig := musee2econfig.ExportContractsFromRunner(deployerRunner, conf)

		// write config into stdout
		configOut, err := filepath.Abs(configOut)
		noError(err)

		noError(config.WriteConfig(configOut, newConfig))

		logger.Print("✅ config file written in %s", configOut)
	}

	deployerRunner.PrintContractAddresses()

	// if setup only, quit
	if setupOnly {
		logger.Print("✅ the localnet has been setup")
		os.Exit(0)
	}

	if upgradeContracts {
		deployerRunner.UpgradeGatewaysAndERC20Custody()
	}
	// always mint ERC20 before every test execution
	deployerRunner.MintERC20OnEVM(1e10)

	// Run the proposals under the start sequence(proposals_e2e_start folder)
	if !skipRegular {
		noError(deployerRunner.CreateGovProposals(runner.StartOfE2E))
	}

	// run tests
	var eg errgroup.Group

	if !skipRegular {
		// start the EVM tests
		startEVMTests(&eg, conf, deployerRunner, verbose)
		startBitcoinTests(&eg, conf, deployerRunner, verbose, light, skipBitcoinSetup)
	}
	if !skipPrecompiles {
		precompiledContractTests := []string{
			//e2etests.TestPrecompilesPrototypeName,
			//e2etests.TestPrecompilesPrototypeThroughContractName,
			//// Disabled until further notice, check https://github.com/RWAs-labs/muse/issues/3005.
			//// e2etests.TestPrecompilesStakingThroughContractName,
			//e2etests.TestPrecompilesBankName,
			//e2etests.TestPrecompilesBankFailName,
			//e2etests.TestPrecompilesBankThroughContractName,
		}
		if e2eStartHeight < 100 {
			// these tests require a clean system
			// since unstaking has an unbonding period
			//precompiledContractTests = append(precompiledContractTests,
			//	e2etests.TestPrecompilesStakingName,
			//	e2etests.TestPrecompilesDistributeName,
			//	e2etests.TestPrecompilesDistributeNonMRC20Name,
			//	e2etests.TestPrecompilesDistributeThroughContractName,
			//)
			// prevent lint error
			_ = precompiledContractTests
		} else {
			logger.Print("⚠️ partial precompiled run (unclean state)")
		}
		eg.Go(statefulPrecompilesTestRoutine(conf, deployerRunner, verbose, precompiledContractTests...))
	}

	if testAdmin {
		eg.Go(adminTestRoutine(conf, deployerRunner, verbose,
			e2etests.TestUpdateMRC20NameName,
			e2etests.TestMuseclientSignerOffsetName,
			e2etests.TestMuseclientRestartHeightName,
			e2etests.TestWhitelistERC20Name,
			e2etests.TestPauseMRC20Name,
			e2etests.TestUpdateBytecodeMRC20Name,
			e2etests.TestUpdateBytecodeConnectorName,
			e2etests.TestDepositEtherLiquidityCapName,
			e2etests.TestCriticalAdminTransactionsName,
			e2etests.TestPauseERC20CustodyName,
			e2etests.TestMigrateERC20CustodyFundsName,
			e2etests.TestUpdateOperationalChainParamsName,

			// Currently this test doesn't work with Anvil because pre-EIP1559 txs are not supported
			// See issue below for details
			// TODO: reenable this test as per the issue below
			// https://github.com/RWAs-labs/muse/issues/1980
			// e2etests.TestMigrateChainSupportName,
		))
	}

	// stress tests
	// TODO: add btc stress tests goroutines
	// https://github.com/RWAs-labs/muse/issues/3909
	if testStress {
		runE2EStressTests(
			conf,
			deployerRunner,
			verbose,
			iterations,
			testEthStress,
			testSolanaStress,
			testSuiStress,
			&eg,
		)
	}

	if testSolana {
		if deployerRunner.SolanaClient == nil {
			logger.Print("❌ solana client is nil, maybe solana rpc is not set")
			os.Exit(1)
		}
		// Run only basic solana tests if an upgrade is in progress
		// This is done to avoid running the tests that take too long to complete
		// Related : https://github.com/RWAs-labs/muse/issues/3666
		solanaTests := []string{
			e2etests.TestSolanaDepositName,
			e2etests.TestSolanaWithdrawName,
		}

		splTests := []string{
			e2etests.TestSolanaDepositName,
			e2etests.TestSPLDepositName,
		}

		if !deployerRunner.IsRunningUpgrade() {
			solanaTests = append(solanaTests, []string{
				e2etests.TestSolanaDepositThroughProgramName,
				e2etests.TestSolanaDepositAndCallName,
				e2etests.TestSolanaWithdrawAndCallName,
				e2etests.TestSolanaWithdrawRevertExecutableReceiverName,
				e2etests.TestSolanaWithdrawAndCallInvalidMsgEncodingName,
				e2etests.TestMEVMToSolanaCallName,
				e2etests.TestSolanaWithdrawAndCallRevertWithCallName,
				e2etests.TestSolanaDepositAndCallRevertName,
				e2etests.TestSolanaDepositAndCallRevertWithCallName,
				e2etests.TestSolanaDepositAndCallRevertWithCallThatRevertsName,
				e2etests.TestSolanaDepositAndCallRevertWithDustName,
				e2etests.TestSolanaDepositRestrictedName,
				e2etests.TestSolanaToMEVMCallName,
				e2etests.TestSolanaToMEVMCallAbortName,
				e2etests.TestSolanaWithdrawRestrictedName,
			}...)

			splTests = append(splTests, []string{
				e2etests.TestSPLDepositAndCallName,
				e2etests.TestSPLDepositAndCallRevertName,
				e2etests.TestSPLDepositAndCallRevertWithCallName,
				e2etests.TestSPLDepositAndCallRevertWithCallThatRevertsName,
				e2etests.TestSPLWithdrawName,
				e2etests.TestSPLWithdrawAndCallName,
				e2etests.TestSPLWithdrawAndCallRevertName,
				e2etests.TestSPLWithdrawAndCreateReceiverAtaName,
				// TODO move under admin tests
				// https://github.com/RWAs-labs/muse/issues/3085
				e2etests.TestSolanaWhitelistSPLName,
			}...)
		}

		eg.Go(
			solanaTestRoutine(
				conf,
				"solana",
				conf.AdditionalAccounts.UserSolana,
				deployerRunner,
				verbose,
				solanaTests...),
		)
		eg.Go(
			solanaTestRoutine(
				conf,
				"spl",
				conf.AdditionalAccounts.UserSPL,
				deployerRunner,
				verbose,
				splTests...),
		)
	}

	if testSui {
		suiTests := []string{
			e2etests.TestSuiDepositName,
			e2etests.TestSuiDepositAndCallRevertName,
			e2etests.TestSuiDepositAndCallName,
			e2etests.TestSuiTokenDepositName,
			e2etests.TestSuiTokenDepositAndCallName,
			e2etests.TestSuiTokenDepositAndCallRevertName,
			e2etests.TestSuiWithdrawName,
			e2etests.TestSuiWithdrawAndCallName,
			e2etests.TestSuiWithdrawRevertWithCallName,
			e2etests.TestSuiWithdrawAndCallRevertWithCallName,
			e2etests.TestSuiTokenWithdrawName,
			e2etests.TestSuiTokenWithdrawAndCallName,
			e2etests.TestSuiTokenWithdrawAndCallRevertWithCallName,
			e2etests.TestSuiDepositRestrictedName,
			e2etests.TestSuiWithdrawRestrictedName,
			e2etests.TestSuiWithdrawInvalidReceiverName,
		}
		eg.Go(suiTestRoutine(conf, deployerRunner, verbose, suiTests...))
	}

	if testTON {
		if deployerRunner.Clients.TON == nil {
			logger.Print("❌ TON client is nil, maybe TON lite-server config is not set")
			os.Exit(1)
		}

		tonTests := []string{
			e2etests.TestTONDepositName,
			e2etests.TestTONDepositAndCallName,
			e2etests.TestTONDepositAndCallRefundName,
			e2etests.TestTONDepositRestrictedName,
			e2etests.TestTONWithdrawName,
			e2etests.TestTONWithdrawConcurrentName,
		}

		eg.Go(tonTestRoutine(conf, deployerRunner, verbose, tonTests...))
	}

	if testLegacy {
		eg.Go(legacyERC20TestRoutine(conf, deployerRunner, verbose,
			e2etests.TestLegacyERC20WithdrawName,
			e2etests.TestLegacyMultipleERC20WithdrawsName,
			e2etests.TestLegacyERC20DepositAndCallRefundName))
		eg.Go(legacyMUSETestRoutine(conf, deployerRunner, verbose,
			e2etests.TestLegacyMuseWithdrawName,
			e2etests.TestLegacyMessagePassingExternalChainsName,
			e2etests.TestLegacyMessagePassingRevertFailExternalChainsName,
			e2etests.TestLegacyMessagePassingRevertSuccessExternalChainsName,
			e2etests.TestLegacyMuseDepositRestrictedName,
			e2etests.TestLegacyMuseDepositName,
			e2etests.TestLegacyMuseDepositNewAddressName,
		))
		eg.Go(legacyMEVMMPTestRoutine(conf, deployerRunner, verbose,
			e2etests.TestLegacyMessagePassingMEVMToEVMName,
			e2etests.TestLegacyMessagePassingEVMtoMEVMName,
			e2etests.TestLegacyMessagePassingEVMtoMEVMRevertName,
			e2etests.TestLegacyMessagePassingMEVMtoEVMRevertName,
			e2etests.TestLegacyMessagePassingMEVMtoEVMRevertFailName,
			e2etests.TestLegacyMessagePassingEVMtoMEVMRevertFailName,
		))
		eg.Go(legacyEthereumTestRoutine(conf, deployerRunner, verbose,
			e2etests.TestLegacyEtherWithdrawName,
			e2etests.TestLegacyEtherDepositAndCallName,
			e2etests.TestLegacyEtherDepositAndCallRefundName,
		))
	}

	// while tests are executed, monitor blocks in parallel to check if system txs are on top and they have biggest priority
	txPriorityErrCh := make(chan error, 1)
	ctx, monitorPriorityCancel := context.WithCancel(context.Background())
	go monitorTxPriorityInBlocks(ctx, conf, txPriorityErrCh)

	if err := eg.Wait(); err != nil {
		deployerRunner.CtxCancel(err)
		monitorPriorityCancel()
		logger.Print("❌ %v", err)
		logger.Print("❌ e2e tests failed after %s", time.Since(testStartTime).String())
		os.Exit(1)
	}

	// Default ballot maturity is set to 30 blocks.
	// We can wait for 31 blocks to ensure that all ballots created during the test are matured, as emission rewards may be slashed for some of the observers based on their vote.
	// This seems to be a problem only in performance tests where we are creating a lot of ballots in a short time. We do not need to slow down regular tests for this check as we expect all observers to vote correctly.
	if testStress {
		deployerRunner.WaitForBlocks(31)
	}

	noError(deployerRunner.WithdrawEmissions())

	// Run gateway upgrade tests for external chains
	deployerRunner.RunGatewayUpgradeTestsExternalChains(conf, runner.UpgradeGatewayOptions{
		TestSolana: testSolana,
		TestSui:    testSui,
	})

	// if all tests pass, cancel txs priority monitoring and check if tx priority is not correct in some blocks
	logger.Print("⏳ e2e tests passed, checking tx priority")
	monitorPriorityCancel()
	if err := <-txPriorityErrCh; err != nil && errors.Is(err, errWrongTxPriority) {
		logger.Print("❌ %v", err)
		logger.Print("❌ e2e tests failed after %s", time.Since(testStartTime).String())
		os.Exit(1)
	}
	if !skipRegular {
		noError(deployerRunner.CreateGovProposals(runner.EndOfE2E))
	}

	logger.Print("✅ e2e tests completed in %s", time.Since(testStartTime).String())

	if testTSSMigration {
		addNewObserver(deployerRunner)
		triggerTSSMigration(deployerRunner, logger, verbose, conf)
	}

	// Verify that there are no trackers left over after tests complete
	if !skipTrackerCheck {
		deployerRunner.EnsureNoTrackers()
	}

	// Verify that the balance of restricted address is zero
	deployerRunner.EnsureZeroBalanceOnRestrictedAddressMEVM()

	if !deployerRunner.IsRunningUpgrade() {
		// Verify that there are no stale ballots left over after tests complete
		deployerRunner.EnsureNoStaleBallots()
	}

	// This should only be run at the end to the test as it would remove the observer.
	if testStaking {
		e2etests.UndelegateToBelowMinimumObserverDelegation(deployerRunner, []string{})
	}
	// print and validate report
	networkReport, err := deployerRunner.GenerateNetworkReport()
	if err != nil {
		logger.Print("❌ failed to generate network report %v", err)
	}
	deployerRunner.PrintNetworkReport(networkReport)
	if err := networkReport.Validate(); err != nil {
		logger.Print("❌ network report validation failed %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

// runE2EStressTests runs the appropriate stress tests based on the provided flags
func runE2EStressTests(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	iterations int,
	testEthStress bool,
	testSolanaStress bool,
	testSuiStress bool,
	eg *errgroup.Group,
) {
	if testEthStress {
		eg.Go(
			ethereumDepositPerformanceRoutine(
				conf,
				deployerRunner,
				verbose,
				[]string{e2etests.TestStressEtherDepositName},
				iterations,
			),
		)
		eg.Go(
			ethereumWithdrawPerformanceRoutine(
				conf,
				deployerRunner,
				verbose,
				[]string{e2etests.TestStressEtherWithdrawName},
				iterations,
			),
		)
	}

	if testSolanaStress {
		eg.Go(
			solanaDepositPerformanceRoutine(
				conf,
				"perf_sol_deposit",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSolana,
				[]string{e2etests.TestStressSolanaDepositName},
			),
		)
		eg.Go(
			solanaDepositPerformanceRoutine(
				conf,
				"perf_spl_deposit",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSPL,
				[]string{e2etests.TestStressSPLDepositName},
			),
		)
		eg.Go(
			solanaWithdrawPerformanceRoutine(
				conf,
				"perf_sol_withdraw",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSolana,
				[]string{e2etests.TestStressSolanaWithdrawName},
			),
		)
		eg.Go(
			solanaWithdrawPerformanceRoutine(
				conf,
				"perf_spl_withdraw",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSPL,
				[]string{e2etests.TestStressSPLWithdrawName},
			),
		)
	}

	if testSuiStress {
		eg.Go(
			suiDepositPerformanceRoutine(
				conf,
				"perf_sui_deposit",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSui,
				[]string{e2etests.TestStressSuiDepositName},
			),
		)
		eg.Go(
			suiWithdrawPerformanceRoutine(
				conf,
				"perf_sui_withdraw",
				deployerRunner,
				verbose,
				// use different account to avoid race conditions on the SUI coin objects
				deployerRunner.Account,
				[]string{e2etests.TestStressSuiWithdrawName},
			),
		)
	}
}

// waitKeygenHeight waits for keygen height
func waitKeygenHeight(
	ctx context.Context,
	cctxClient crosschaintypes.QueryClient,
	observerClient observertypes.QueryClient,
	logger *runner.Logger,
	bufferBlocks int64,
) error {
	// wait for keygen to be completed
	resp, err := observerClient.Keygen(ctx, &observertypes.QueryGetKeygenRequest{})

	switch {
	case err != nil:
		return errors.Wrap(err, "observerClient.Keygen error")
	case resp.Keygen == nil:
		return errors.New("keygen is nil")
	case resp.Keygen.Status != observertypes.KeygenStatus_PendingKeygen:
		return errors.Errorf("keygen is not pending (status: %s)", resp.Keygen.Status.String())
	}

	keygenHeight := resp.Keygen.BlockNumber
	logger.Print("⏳ wait height %v for keygen to be completed", keygenHeight)

	for {
		time.Sleep(2 * time.Second)

		response, err := cctxClient.LastMuseHeight(ctx, &crosschaintypes.QueryLastMuseHeightRequest{})
		if err != nil {
			logger.Error("cctxClient.LastMuseHeight error: %s", err)
			continue
		}

		logger.Info("Last MuseHeight: %d", response.Height)

		if response.Height >= keygenHeight+bufferBlocks {
			return nil
		}
	}
}

func must[T any](v T, err error) T {
	return testutil.Must(v, err)
}
