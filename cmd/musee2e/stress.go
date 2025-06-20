package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	musee2econfig "github.com/RWAs-labs/muse/cmd/musee2e/config"
	"github.com/RWAs-labs/muse/cmd/musee2e/local"
	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/txserver"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/testutil"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

const (
	StatInterval      = 5
	StressTestTimeout = 100 * time.Minute
)

var (
	mevmNonce = big.NewInt(1)
)

type stressArguments struct {
	network           string
	txnInterval       int64
	contractsDeployed bool
	config            string
}

var stressTestArgs = stressArguments{}

var noError = testutil.NoError

func NewStressTestCmd() *cobra.Command {
	var StressCmd = &cobra.Command{
		Use:   "stress",
		Short: "Run Stress Test",
		Run:   StressTest,
	}

	StressCmd.Flags().StringVar(&stressTestArgs.network, "network", "LOCAL", "--network TESTNET")
	StressCmd.Flags().
		Int64Var(&stressTestArgs.txnInterval, "tx-interval", 500, "--tx-interval [TIME_INTERVAL_MILLISECONDS]")
	StressCmd.Flags().
		BoolVar(&stressTestArgs.contractsDeployed, "contracts-deployed", false, "--contracts-deployed=false")
	StressCmd.Flags().StringVar(&stressTestArgs.config, local.FlagConfigFile, "", "config file to use for the E2E test")
	StressCmd.Flags().Bool(flagVerbose, false, "set to true to enable verbose logging")

	return StressCmd
}

func StressTest(cmd *cobra.Command, _ []string) {
	testStartTime := time.Now()
	defer func() {
		fmt.Println("E2E test took", time.Since(testStartTime))
	}()
	go func() {
		time.Sleep(StressTestTimeout)
		fmt.Println("E2E test timed out after", StressTestTimeout)
		os.Exit(1)
	}()

	// initialize E2E tests config
	conf := must(local.GetConfig(cmd))

	deployerAccount := conf.DefaultAccount

	// Initialize clients ----------------------------------------------------------------
	evmClient := must(ethclient.Dial(conf.RPCs.EVM))
	bal := must(evmClient.BalanceAt(context.TODO(), deployerAccount.EVMAddress(), nil))

	fmt.Printf("Deployer address: %s, balance: %d Wei\n", deployerAccount.EVMAddress().Hex(), bal)

	grpcConn := must(grpc.Dial(conf.RPCs.MuseCoreGRPC, grpc.WithInsecure()))

	cctxClient := crosschaintypes.NewQueryClient(grpcConn)
	// -----------------------------------------------------------------------------------

	// Wait for Genesis and keygen to be completed if network is local. ~ height 30
	if stressTestArgs.network == "LOCAL" {
		time.Sleep(20 * time.Second)
		for {
			time.Sleep(5 * time.Second)
			response, err := cctxClient.LastMuseHeight(
				context.Background(),
				&crosschaintypes.QueryLastMuseHeightRequest{},
			)
			if err != nil {
				fmt.Printf("cctxClient.LastMuseHeight error: %s", err)
				continue
			}
			if response.Height >= 30 {
				break
			}
			fmt.Printf("Last MuseHeight: %d\n", response.Height)
		}
	}

	// initialize context
	ctx, cancel := context.WithCancelCause(context.Background())

	verbose := must(cmd.Flags().GetBool(flagVerbose))
	logger := runner.NewLogger(verbose, color.FgWhite, "setup")

	// initialize E2E test runner
	e2eTest := must(musee2econfig.RunnerFromConfig(
		ctx,
		"deployer",
		cancel,
		conf,
		conf.DefaultAccount,
		logger,
	))

	// setup TSS addresses
	noError(e2eTest.SetTSSAddresses())
	e2eTest.LegacySetupEVM(stressTestArgs.contractsDeployed, false)

	// If stress test is running on local docker environment
	switch stressTestArgs.network {
	case "LOCAL":
		// deploy and set mevm contract
		e2eTest.SetupMEVMProtocolContracts()
		e2eTest.SetupMEVMMRC20s(txserver.MRC20Deployment{
			ERC20Addr: e2eTest.ERC20Addr,
			SPLAddr:   nil, // no stress tests for solana atm
		})

		// deposit on MuseChain
		e2eTest.DepositEtherDeployer()
		e2eTest.LegacyDepositMuse()
	case "TESTNET":
		ethMRC20Addr := must(e2eTest.SystemContract.GasCoinMRC20ByChainId(&bind.CallOpts{}, big.NewInt(5)))
		e2eTest.ETHMRC20Addr = ethMRC20Addr

		e2eTest.ETHMRC20 = must(mrc20.NewMRC20(e2eTest.ETHMRC20Addr, e2eTest.MEVMClient))
	default:
		noError(errors.New("invalid network argument: " + stressTestArgs.network))
	}

	// Check mrc20 balance of Deployer address
	ethMRC20Balance := must(e2eTest.ETHMRC20.BalanceOf(nil, deployerAccount.EVMAddress()))
	fmt.Printf("eth mrc20 balance: %s Wei \n", ethMRC20Balance.String())

	//Pre-approve ETH withdraw on MEVM
	fmt.Println("approving ETH MRC20...")
	ethMRC20 := e2eTest.ETHMRC20
	tx := must(ethMRC20.Approve(e2eTest.MEVMAuth, e2eTest.ETHMRC20Addr, big.NewInt(1e18)))

	receipt := utils.MustWaitForTxReceipt(e2eTest.Ctx, e2eTest.MEVMClient, tx, logger, e2eTest.ReceiptTimeout)
	fmt.Printf("eth mrc20 approve receipt: status %d\n", receipt.Status)

	// Get current nonce on mevm for DeployerAddress - Need to keep track of nonce at client level
	blockNum := must(e2eTest.MEVMClient.BlockNumber(ctx))

	// #nosec G115 e2eTest - always in range
	nonce := must(e2eTest.MEVMClient.NonceAt(ctx, deployerAccount.EVMAddress(), big.NewInt(int64(blockNum))))

	// #nosec G115 e2e - always in range
	mevmNonce = big.NewInt(int64(nonce))

	// -------------- TEST BEGINS ------------------

	fmt.Println("**** STRESS TEST BEGINS ****")
	fmt.Println("	1. Periodically Withdraw ETH from MEVM to EVM")
	fmt.Println("	2. Display Network metrics to monitor performance [Num Pending outbound tx], [Num Trackers]")

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		// Withdraw from MEVM to EVM
		WithdrawCCtx(e2eTest)
	}()

	go func() {
		defer wg.Done()

		// Display Network metrics periodically to monitor performance
		EchoNetworkMetrics(e2eTest)
	}()

	wg.Wait()
}

// WithdrawCCtx withdraw ETHMRC20 from MEVM to EVM
func WithdrawCCtx(runner *runner.E2ERunner) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(stressTestArgs.txnInterval))
	for {
		select {
		case <-ticker.C:
			WithdrawETHMRC20(runner)
		}
	}
}

func EchoNetworkMetrics(r *runner.E2ERunner) {
	var (
		ticker            = time.NewTicker(time.Second * StatInterval)
		queue             = make([]uint64, 0)
		numTicks          int
		totalMinedTxns    uint64
		previousMinedTxns uint64
		chainID           = must(getChainID(r.EVMClient))
	)

	for {
		select {
		case <-ticker.C:
			numTicks++
			// Get all pending outbound transactions
			cctxResp, err := r.CctxClient.ListPendingCctx(
				context.Background(),
				&crosschaintypes.QueryListPendingCctxRequest{
					ChainId: chainID.Int64(),
				},
			)
			if err != nil {
				continue
			}
			sends := cctxResp.CrossChainTx
			sort.Slice(sends, func(i, j int) bool {
				return sends[i].GetCurrentOutboundParam().TssNonce < sends[j].GetCurrentOutboundParam().TssNonce
			})
			if len(sends) > 0 {
				fmt.Printf(
					"pending nonces %d to %d\n",
					sends[0].GetCurrentOutboundParam().TssNonce,
					sends[len(sends)-1].GetCurrentOutboundParam().TssNonce,
				)
			} else {
				continue
			}
			//
			// Get all trackers
			trackerResp, err := r.CctxClient.OutboundTrackerAll(
				context.Background(),
				&crosschaintypes.QueryAllOutboundTrackerRequest{},
			)
			if err != nil {
				continue
			}

			currentMinedTxns := sends[0].GetCurrentOutboundParam().TssNonce
			newMinedTxCnt := currentMinedTxns - previousMinedTxns
			previousMinedTxns = currentMinedTxns

			// Add new mined txn count to queue and remove the oldest entry
			queue = append(queue, newMinedTxCnt)
			if numTicks > 60/StatInterval {
				totalMinedTxns -= queue[0]
				queue = queue[1:]
				numTicks = 60/StatInterval + 1 //prevent overflow
			}

			//Calculate rate -> tx/min
			totalMinedTxns += queue[len(queue)-1]
			rate := totalMinedTxns

			numPending := len(cctxResp.CrossChainTx)
			numTrackers := len(trackerResp.OutboundTracker)

			fmt.Println(
				"Network Stat => Num of Pending cctx: ",
				numPending,
				"Num active trackers: ",
				numTrackers,
				"Tx Rate: ",
				rate,
				" tx/min",
			)
		}
	}
}

func WithdrawETHMRC20(r *runner.E2ERunner) {
	defer func() {
		mevmNonce.Add(mevmNonce, big.NewInt(1))
	}()

	ethMRC20 := r.ETHMRC20
	r.MEVMAuth.Nonce = mevmNonce

	must(ethMRC20.Withdraw(r.MEVMAuth, r.EVMAddress().Bytes(), big.NewInt(100)))
}

// Get ETH based chain ID
func getChainID(client *ethclient.Client) (*big.Int, error) {
	return client.ChainID(context.Background())
}

func must[T any](v T, err error) T {
	return testutil.Must(v, err)
}
