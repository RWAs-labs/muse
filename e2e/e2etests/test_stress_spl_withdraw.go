package e2etests

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestStressSPLWithdraw tests the stressing withdrawal of SPL
func TestStressSPLWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	withdrawSPLAmount := utils.ParseBigInt(r, args[0])
	numWithdrawalsSPL := utils.ParseInt(r, args[1])

	// load deployer private key
	privKey := r.GetSolanaPrivKey()

	r.Logger.Print("starting stress test of %d SPL withdrawals", numWithdrawalsSPL)

	tx, err := r.SOLMRC20.Approve(r.MEVMAuth, r.SOLMRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve_sol")

	tx, err = r.SPLMRC20.Approve(r.MEVMAuth, r.SPLMRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve_spl")

	tx, err = r.SOLMRC20.Approve(r.MEVMAuth, r.SPLMRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve_spl_sol")

	// create a wait group to wait for all the withdrawals to complete
	var eg errgroup.Group

	// store durations as float64 seconds like prometheus
	withdrawDurations := []float64{}
	withdrawDurationsLock := sync.Mutex{}

	// send the withdrawals SPL
	for i := 0; i < numWithdrawalsSPL; i++ {
		i := i

		// execute the withdraw SPL transaction
		tx, err := r.SPLMRC20.Withdraw(r.MEVMAuth, []byte(privKey.PublicKey().String()), withdrawSPLAmount)
		require.NoError(r, err)

		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.MEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt)

		r.Logger.Print("index %d: starting SPL withdraw, tx hash: %s", i, tx.Hash().Hex())

		eg.Go(func() error {
			segmentStartTime := time.Now()
			cctxFirstHash := utils.WaitCctxByInboundHash(
				r.Ctx,
				r,
				tx.Hash().Hex(),
				r.CctxClient,
				utils.HasOutboundTxHash(),
			)
			r.Logger.Info("index %d: got first outbound hash: %s", i, cctxFirstHash.OutboundParams[0].Hash)
			timeToOutboundHash := time.Since(segmentStartTime)

			segmentStartTime = time.Now()
			cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.ReceiptTimeout)
			if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
				return fmt.Errorf(
					"index %d: withdraw cctx failed with status %s, message %s, cctx index %s",
					i,
					cctx.CctxStatus.Status,
					cctx.CctxStatus.StatusMessage,
					cctx.Index,
				)
			}
			timeToFinalized := time.Since(segmentStartTime)
			totalTime := timeToOutboundHash + timeToFinalized
			r.Logger.Print("index %d: withdraw SPL cctx success in %s (outbound hash: %s + finalized: %s)",
				i,
				formatDuration(totalTime),
				formatDuration(timeToOutboundHash),
				formatDuration(timeToFinalized),
			)

			withdrawDurationsLock.Lock()
			withdrawDurations = append(withdrawDurations, totalTime.Seconds())
			withdrawDurationsLock.Unlock()

			return nil
		})
	}

	err = eg.Wait()

	desc, descErr := stats.Describe(withdrawDurations, false, &[]float64{50.0, 75.0, 90.0, 95.0})
	if descErr != nil {
		r.Logger.Print("❌ failed to calculate latency report: %v", descErr)
	}

	r.Logger.Print("Latency report:")
	r.Logger.Print("min:  %.2f", desc.Min)
	r.Logger.Print("max:  %.2f", desc.Max)
	r.Logger.Print("mean: %.2f", desc.Mean)
	r.Logger.Print("std:  %.2f", desc.Std)
	for _, p := range desc.DescriptionPercentiles {
		r.Logger.Print("p%.0f:  %.2f", p.Percentile, p.Value)
	}

	require.NoError(r, err)
	r.Logger.Print("all SPL withdrawals completed")
}
