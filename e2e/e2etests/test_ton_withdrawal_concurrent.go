package e2etests

import (
	"math/rand"
	"sync"

	"cosmossdk.io/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	toncontracts "github.com/RWAs-labs/muse/pkg/contracts/ton"
	"github.com/RWAs-labs/muse/testutil/sample"
	cc "github.com/RWAs-labs/muse/x/crosschain/types"
)

// TestTONWithdrawConcurrent makes sure that multiple concurrent
// withdrawals will be eventually processed by sequentially increasing Gateway nonce
// and that museclient tolerates "invalid nonce" error from RPC.
func TestTONWithdrawConcurrent(r *runner.E2ERunner, _ []string) {
	// ARRANGE
	const recipientsCount = 10

	// Fire withdrawals. Note that mevm sender is r.MEVMAuth
	var wg sync.WaitGroup
	for i := 0; i < recipientsCount; i++ {
		// ARRANGE
		// Given multiple recipients WITHOUT deployed wallet-contracts
		// and withdrawal amounts between 1 and 5 TON
		var (
			// #nosec G404: it's a test
			amountCoins = 1 + rand.Intn(5)
			// #nosec G115 test - always in range
			amount    = toncontracts.Coins(uint64(amountCoins))
			recipient = sample.GenerateTONAccountID()
		)

		// ACT
		r.Logger.Info(
			"Withdrawal #%d: sending %s to %s",
			i+1,
			toncontracts.FormatCoins(amount),
			recipient.ToRaw(),
		)

		approvedAmount := amount.Add(toncontracts.Coins(1))
		tx := r.SendWithdrawTONMRC20(recipient, amount.BigInt(), approvedAmount.BigInt())

		wg.Add(1)

		go func(number int, recipient ton.AccountID, amount math.Uint, tx *ethtypes.Transaction) {
			defer wg.Done()

			// wait for the cctx to be mined
			cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

			// ASSERT
			utils.RequireCCTXStatus(r, cctx, cc.CctxStatus_OutboundMined)
			r.Logger.Info("Withdrawal #%d complete! cctx index: %s", number, cctx.Index)

			// Check recipient's balance ON TON
			balance, err := r.Clients.TON.GetBalanceOf(r.Ctx, recipient, false)
			require.NoError(r, err, "failed to get balance of %s", recipient.ToRaw())
			require.Equal(r, amount.Uint64(), balance.Uint64())
		}(i+1, recipient, amount, tx)
	}

	wg.Wait()
}
