package legacy

import (
	"context"
	"fmt"
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// WithdrawType is the type of withdraw to perform in the test
type withdrawType string

const (
	withdrawTypeMUSE  withdrawType = "MUSE"
	withdrawTypeETH   withdrawType = "ETH"
	withdrawTypeERC20 withdrawType = "ERC20"

	rateLimiterWithdrawNumber = 5
)

func TestRateLimiter(r *runner.E2ERunner, _ []string) {
	r.Logger.Info("TestRateLimiter")

	// rateLimiterFlags are the rate limiter flags for the test
	rateLimiterFlags := crosschaintypes.RateLimiterFlags{
		Enabled: true,
		Rate:    sdkmath.NewUint(1e17).MulUint64(5), // 0.5 MUSE this value is used so rate is reached
		Window:  10,
		Conversions: []crosschaintypes.Conversion{
			{
				Mrc20: r.ETHMRC20Addr.Hex(),
				Rate:  sdkmath.LegacyNewDec(2), // 1 ETH = 2 MUSE
			},
			{
				Mrc20: r.ERC20MRC20Addr.Hex(),
				Rate:  sdkmath.LegacyNewDec(1).QuoInt64(2), // 2 USDC = 1 MUSE
			},
		},
	}

	// these are the amounts for the withdraws for the different types
	// currently these are arbitrary values that can be fine-tuned for manual testing of rate limiter
	// TODO: define more rigorous assertions with proper values
	// https://github.com/RWAs-labs/muse/issues/2090
	museAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(3))
	ethAmount := big.NewInt(1e18)
	erc20Amount := big.NewInt(1e6)

	// approve tokens for the tests
	require.NoError(r, approveTokens(r))

	// add liquidity in the pool to prevent high slippage in WMUSE/gas pair
	require.NoError(r, addMuseGasLiquidity(r))

	// Set the rate limiter to 0.5MUSE per 10 blocks
	// These rate limiter flags will only allow to process 1 withdraw per 10 blocks
	r.Logger.Info("setting up rate limiter flags")
	require.NoError(r, setupRateLimiterFlags(r, rateLimiterFlags))

	// Test with rate limiter
	// TODO: define proper assertion to check the rate limiter is working
	// https://github.com/RWAs-labs/muse/issues/2090
	r.Logger.Print("rate limiter enabled")
	require.NoError(r, createAndWaitWithdraws(r, withdrawTypeMUSE, museAmount))
	require.NoError(r, createAndWaitWithdraws(r, withdrawTypeETH, ethAmount))
	require.NoError(r, createAndWaitWithdraws(r, withdrawTypeERC20, erc20Amount))

	// Disable rate limiter
	r.Logger.Info("disabling rate limiter")
	require.NoError(r, setupRateLimiterFlags(r, crosschaintypes.RateLimiterFlags{Enabled: false}))

	// Test without rate limiter again and try again MUSE withdraws
	r.Logger.Print("rate limiter disabled")
	require.NoError(r, createAndWaitWithdraws(r, withdrawTypeMUSE, museAmount))
}

// createAndWaitWithdraws performs RateLimiterWithdrawNumber withdraws
func createAndWaitWithdraws(r *runner.E2ERunner, withdrawType withdrawType, withdrawAmount *big.Int) error {
	startTime := time.Now()

	r.Logger.Print("starting %d %s withdraws", rateLimiterWithdrawNumber, withdrawType)

	// Perform RateLimiterWithdrawNumber withdraws to log time for completion
	txs := make([]*ethtypes.Transaction, rateLimiterWithdrawNumber)
	for i := 0; i < rateLimiterWithdrawNumber; i++ {
		// create a new withdraw depending on the type
		switch withdrawType {
		case withdrawTypeMUSE:
			txs[i] = r.LegacyWithdrawMuse(withdrawAmount, true)
		case withdrawTypeETH:
			txs[i] = r.LegacyWithdrawEther(withdrawAmount)
		case withdrawTypeERC20:
			txs[i] = r.LegacyWithdrawERC20(withdrawAmount)
		default:
			return fmt.Errorf("invalid withdraw type: %s", withdrawType)
		}
	}

	// start a error group to wait for all the withdraws to be mined
	g, ctx := errgroup.WithContext(r.Ctx)
	for i, tx := range txs {
		// capture the loop variables
		tx, i := tx, i

		// start a goroutine to wait for the withdraw to be mined
		g.Go(func() error {
			return waitForWithdrawMined(ctx, r, tx, i, startTime)
		})
	}

	// wait for all the withdraws to be mined
	if err := g.Wait(); err != nil {
		return err
	}

	duration := time.Since(startTime).Seconds()
	block, err := r.MEVMClient.BlockNumber(r.Ctx)
	if err != nil {
		return fmt.Errorf("error getting block number: %w", err)
	}
	r.Logger.Print("all withdraws completed in %vs at block %d", duration, block)

	return nil
}

// waitForWithdrawMined waits for a withdraw to be mined
// we first wait to get the receipt
// NOTE: this could be a more general function but we define it here for this test because we emit in the function logs specific to this test
func waitForWithdrawMined(
	ctx context.Context,
	r *runner.E2ERunner,
	tx *ethtypes.Transaction,
	index int,
	startTime time.Time,
) error {
	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"expected cctx status to be %s; got %s, message %s",
			crosschaintypes.CctxStatus_OutboundMined,
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage,
		)
	}

	// record the time for completion
	duration := time.Since(startTime).Seconds()
	block, err := r.MEVMClient.BlockNumber(ctx)
	if err != nil {
		return err
	}
	r.Logger.Print("cctx %d mined in %vs at block %d", index, duration, block)

	return nil
}

// setupRateLimiterFlags sets up the rate limiter flags with flags defined in the test
func setupRateLimiterFlags(r *runner.E2ERunner, flags crosschaintypes.RateLimiterFlags) error {
	adminAddr, err := r.MuseTxServer.GetAccountAddressFromName(utils.OperationalPolicyName)
	if err != nil {
		return err
	}
	_, err = r.MuseTxServer.BroadcastTx(utils.OperationalPolicyName, crosschaintypes.NewMsgUpdateRateLimiterFlags(
		adminAddr,
		flags,
	))
	if err != nil {
		return err
	}

	return nil
}

// addMuseGasLiquidity adds liquidity to the MUSE/gas pool
func addMuseGasLiquidity(r *runner.E2ERunner) error {
	// use 10 MUSE and 10 ETH for the liquidity
	// this will be sufficient for the tests
	amount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(10))
	approveAmount := big.NewInt(0).Mul(amount, big.NewInt(10))

	// approve uniswap router to spend gas
	txETHMRC20Approve, err := r.ETHMRC20.Approve(r.MEVMAuth, r.UniswapV2RouterAddr, approveAmount)
	if err != nil {
		return fmt.Errorf("error approving MUSE: %w", err)
	}

	// wait for the tx to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, txETHMRC20Approve, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		return fmt.Errorf("approve failed")
	}

	// approve uniswap router to spend MUSE
	txMUSEApprove, err := r.WMuse.Approve(r.MEVMAuth, r.UniswapV2RouterAddr, approveAmount)
	if err != nil {
		return fmt.Errorf("error approving MUSE: %w", err)
	}

	// wait for the tx to be mined
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, txMUSEApprove, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		return fmt.Errorf("approve failed")
	}

	// add liquidity in the pool to prevent high slippage in WMUSE/gas pair
	r.MEVMAuth.Value = amount
	txAddLiquidity, err := r.UniswapV2Router.AddLiquidityETH(
		r.MEVMAuth,
		r.ETHMRC20Addr,
		amount,
		big.NewInt(1e18),
		big.NewInt(1e18),
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	if err != nil {
		return fmt.Errorf("error adding liquidity: %w", err)
	}
	r.MEVMAuth.Value = big.NewInt(0)

	// wait for the tx to be mined
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, txAddLiquidity, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		return fmt.Errorf("add liquidity failed")
	}

	return nil
}

// approveTokens approves the tokens for the tests
func approveTokens(r *runner.E2ERunner) error {
	// deposit and approve 50 WMUSE for the tests
	approveAmount := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(50))
	r.LegacyDepositAndApproveWMuse(approveAmount)

	// approve ETH for withdraws
	tx, err := r.ETHMRC20.Approve(r.MEVMAuth, r.ETHMRC20Addr, approveAmount)
	if err != nil {
		return fmt.Errorf("error approving ETH: %w", err)
	}
	r.Logger.EVMTransaction(*tx, "approve")

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		return fmt.Errorf("eth approve failed")
	}
	r.Logger.EVMReceipt(*receipt, "approve")

	// approve ETH for ERC20 withdraws (this is for the gas fees)
	tx, err = r.ETHMRC20.Approve(r.MEVMAuth, r.ERC20MRC20Addr, approveAmount)
	if err != nil {
		return fmt.Errorf("error approving ERC20: %w", err)
	}

	r.Logger.EVMTransaction(*tx, "approve")

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		return fmt.Errorf("erc 20 approve failed")
	}
	r.Logger.EVMReceipt(*receipt, "approve")

	return nil
}
