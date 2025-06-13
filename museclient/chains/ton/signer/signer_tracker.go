package signer

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"

	"github.com/RWAs-labs/muse/museclient/chains/interfaces"
	"github.com/RWAs-labs/muse/museclient/chains/ton/liteapi"
	"github.com/RWAs-labs/muse/museclient/metrics"
	toncontracts "github.com/RWAs-labs/muse/pkg/contracts/ton"
)

// trackOutbound tracks sent external message and records it as outboundTracker.
// Explanation:
// Due to TON's nature, it's not possible to get tx hash before it's confirmed on-chain,
// So we need to poll from the latest account state (prevState) up to the most recent tx
// and search for desired tx hash. After it's found, we can record it as outboundTracker.
//
// Note that another museclient observers that scrolls Gateway's txs can publish this tracker concurrently.
func (s *Signer) trackOutbound(
	ctx context.Context,
	musecore interfaces.MusecoreClient,
	w *toncontracts.Withdrawal,
	prevState tlb.ShardAccount,
) error {
	metrics.NumTrackerReporters.WithLabelValues(s.Chain().Name).Inc()
	defer metrics.NumTrackerReporters.WithLabelValues(s.Chain().Name).Dec()

	const (
		timeout = 60 * time.Second
		tick    = time.Second
	)

	var (
		start   = time.Now()
		chainID = s.Chain().ChainId

		acc   = s.gateway.AccountID()
		lt    = prevState.LastTransLt
		hash  = ton.Bits256(prevState.LastTransHash)
		nonce = uint64(w.Seqno)

		filter = withdrawalFilter(w)
	)

	for time.Since(start) <= timeout {
		txs, err := s.client.GetTransactionsSince(ctx, acc, lt, hash)
		if err != nil {
			return errors.Wrapf(err, "unable to get transactions (lt %d, hash %s)", lt, hash.Hex())
		}

		results := s.gateway.ParseAndFilterMany(txs, filter)
		if len(results) == 0 {
			time.Sleep(tick)
			continue
		}

		tx := results[0].Transaction
		txHash := liteapi.TransactionToHashString(results[0].Transaction)

		if !tx.IsSuccess() {
			// should not happen
			return errors.Errorf("transaction %q is not successful", txHash)
		}

		// Note that this method has a check for noop
		_, err = musecore.PostOutboundTracker(ctx, chainID, nonce, txHash)
		if err != nil {
			return errors.Wrap(err, "unable to add outbound tracker")
		}

		return nil
	}

	return errors.Errorf("timeout exceeded (%s)", time.Since(start).String())
}

// creates a tx filter for this very withdrawal
func withdrawalFilter(w *toncontracts.Withdrawal) func(tx *toncontracts.Transaction) bool {
	return func(tx *toncontracts.Transaction) bool {
		if !tx.IsOutbound() || tx.Operation != toncontracts.OpWithdraw {
			return false
		}

		wd, err := tx.Withdrawal()
		if err != nil {
			return false
		}

		return wd.Seqno == w.Seqno && wd.Sig == w.Sig
	}
}
