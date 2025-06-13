package ratelimiter_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/museclient/ratelimiter"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/testutil/sample"
	crosschainkeeper "github.com/RWAs-labs/muse/x/crosschain/keeper"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
)

func Test_NewInput(t *testing.T) {
	// sample response
	response := crosschaintypes.QueryRateLimiterInputResponse{
		Height:                  10,
		CctxsMissed:             []*crosschaintypes.CrossChainTx{sample.CrossChainTx(t, "1-1")},
		CctxsPending:            []*crosschaintypes.CrossChainTx{sample.CrossChainTx(t, "1-2")},
		TotalPending:            7,
		PastCctxsValue:          sdkmath.NewInt(12345678).Mul(sdkmath.NewInt(1e18)).String(),
		PendingCctxsValue:       sdkmath.NewInt(4321).Mul(sdkmath.NewInt(1e18)).String(),
		LowestPendingCctxHeight: 2,
	}

	t.Run("should create a input from gRPC response", func(t *testing.T) {
		filterInput, ok := ratelimiter.NewInput(response)
		require.True(t, ok)
		require.Equal(t, response.Height, filterInput.Height)
		require.Equal(t, response.CctxsMissed, filterInput.CctxsMissed)
		require.Equal(t, response.CctxsPending, filterInput.CctxsPending)
		require.Equal(t, response.PastCctxsValue, filterInput.PastCctxsValue.String())
		require.Equal(t, response.PendingCctxsValue, filterInput.PendingCctxsValue.String())
		require.Equal(t, response.LowestPendingCctxHeight, filterInput.LowestPendingCctxHeight)
	})
	t.Run("should return false if past cctxs value is invalid", func(t *testing.T) {
		invalidResp := response
		invalidResp.PastCctxsValue = "invalid"
		filterInput, ok := ratelimiter.NewInput(invalidResp)
		require.False(t, ok)
		require.Nil(t, filterInput)
	})
	t.Run("should return false if pending cctxs value is invalid", func(t *testing.T) {
		invalidResp := response
		invalidResp.PendingCctxsValue = "invalid"
		filterInput, ok := ratelimiter.NewInput(invalidResp)
		require.False(t, ok)
		require.Nil(t, filterInput)
	})
}

func Test_IsRateLimiterUsable(t *testing.T) {
	tests := []struct {
		name     string
		flags    crosschaintypes.RateLimiterFlags
		expected bool
	}{
		{
			name: "rate limiter is enabled",
			flags: crosschaintypes.RateLimiterFlags{
				Enabled: true,
				Window:  100,
				Rate:    sdkmath.NewUint(1e18), // 1 MUSE/block
			},
			expected: true,
		},
		{
			name: "rate limiter is disabled",
			flags: crosschaintypes.RateLimiterFlags{
				Enabled: false,
			},
			expected: false,
		},
		{
			name: "rate limiter is enabled with 0 window",
			flags: crosschaintypes.RateLimiterFlags{
				Enabled: true,
				Window:  0,
			},
			expected: false,
		},
		{
			name: "rate limiter is enabled with nil rate",
			flags: crosschaintypes.RateLimiterFlags{
				Enabled: true,
				Window:  100,
				Rate:    sdkmath.Uint{},
			},
			expected: false,
		},
		{
			name: "rate limiter is enabled with zero rate",
			flags: crosschaintypes.RateLimiterFlags{
				Enabled: true,
				Window:  100,
				Rate:    sdkmath.NewUint(0),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usable := ratelimiter.IsRateLimiterUsable(tt.flags)
			require.Equal(t, tt.expected, usable)
		})
	}
}

func Test_ApplyRateLimiter(t *testing.T) {
	// define test chain ids
	ethChainID := chains.Ethereum.ChainId
	btcChainID := chains.BitcoinMainnet.ChainId
	museChainID := chains.MuseChainMainnet.ChainId

	// create 10 missed and 90 pending cctxs for eth chain, the coinType/amount does not matter for this test
	// but we still use a proper cctx value (0.5 MUSE) to make the test more realistic
	ethCctxsMissed := sample.CustomCctxsInBlockRange(
		t,
		1,
		10,
		museChainID,
		ethChainID,
		coin.CoinType_Gas,
		"",
		uint64(2e14),
		crosschaintypes.CctxStatus_PendingOutbound,
	)
	ethCctxsPending := sample.CustomCctxsInBlockRange(
		t,
		11,
		100,
		museChainID,
		ethChainID,
		coin.CoinType_Gas,
		"",
		uint64(2e14),
		crosschaintypes.CctxStatus_PendingOutbound,
	)
	ethCctxsAll := append(append([]*crosschaintypes.CrossChainTx{}, ethCctxsMissed...), ethCctxsPending...)

	// create 10 missed and 90 pending cctxs for btc chain, the coinType/amount does not matter for this test
	// but we still use a proper cctx value (0.5 MUSE) to make the test more realistic
	btcCctxsMissed := sample.CustomCctxsInBlockRange(
		t,
		1,
		10,
		museChainID,
		btcChainID,
		coin.CoinType_Gas,
		"",
		2000,
		crosschaintypes.CctxStatus_PendingOutbound,
	)
	btcCctxsPending := sample.CustomCctxsInBlockRange(
		t,
		11,
		100,
		museChainID,
		btcChainID,
		coin.CoinType_Gas,
		"",
		2000,
		crosschaintypes.CctxStatus_PendingOutbound,
	)
	btcCctxsAll := append(append([]*crosschaintypes.CrossChainTx{}, btcCctxsMissed...), btcCctxsPending...)

	// all missed cctxs and all pending cctxs across all chains
	allCctxsMissed := crosschainkeeper.SortCctxsByHeightAndChainID(
		append(append([]*crosschaintypes.CrossChainTx{}, ethCctxsMissed...), btcCctxsMissed...))
	allCctxsPending := crosschainkeeper.SortCctxsByHeightAndChainID(
		append(append([]*crosschaintypes.CrossChainTx{}, ethCctxsPending...), btcCctxsPending...))

	// define test cases
	tests := []struct {
		name   string
		window int64
		rate   sdkmath.Uint
		input  ratelimiter.Input
		output ratelimiter.Output
	}{
		{
			name:   "should return all missed and pending cctxs",
			window: 100,
			rate:   sdkmath.NewUint(1e18), // 1 MUSE/block
			input: ratelimiter.Input{
				Height:                  100,
				CctxsMissed:             allCctxsMissed,
				CctxsPending:            allCctxsPending,
				PastCctxsValue:          sdkmath.NewInt(10).Mul(sdkmath.NewInt(1e18)), // 10 * 1 MUSE
				PendingCctxsValue:       sdkmath.NewInt(90).Mul(sdkmath.NewInt(1e18)), // 90 * 1 MUSE
				LowestPendingCctxHeight: 11,
			},
			output: ratelimiter.Output{
				CctxsMap: map[int64][]*crosschaintypes.CrossChainTx{
					ethChainID: ethCctxsAll,
					btcChainID: btcCctxsAll,
				},
				CurrentWithdrawWindow: 100,                  // height [1, 100]
				CurrentWithdrawRate:   sdkmath.NewInt(1e18), // (10 + 90) / 100
				RateLimitExceeded:     false,
			},
		},
		{
			name:   "should monitor a wider window and adjust the total limit",
			window: 50,
			rate:   sdkmath.NewUint(1e18), // 1 MUSE/block
			input: ratelimiter.Input{
				Height:       100,
				CctxsMissed:  allCctxsMissed,
				CctxsPending: allCctxsPending,
				PastCctxsValue: sdkmath.NewInt(
					0,
				), // no past cctx in height range [51, 100]
				PendingCctxsValue:       sdkmath.NewInt(90).Mul(sdkmath.NewInt(1e18)), // 90 * 1 MUSE
				LowestPendingCctxHeight: 11,
			},
			output: ratelimiter.Output{
				CctxsMap: map[int64][]*crosschaintypes.CrossChainTx{
					ethChainID: ethCctxsAll,
					btcChainID: btcCctxsAll,
				},
				CurrentWithdrawWindow: 90,                   // [LowestPendingCctxHeight, Height] = [11, 100]
				CurrentWithdrawRate:   sdkmath.NewInt(1e18), // 90 / 90 = 1 MUSE/block
				RateLimitExceeded:     false,
			},
		},
		{
			name:   "rate limit is exceeded in given sliding window 100",
			window: 100,
			rate:   sdkmath.NewUint(1e18), // 1 MUSE/block
			input: ratelimiter.Input{
				Height:       100,
				CctxsMissed:  allCctxsMissed,
				CctxsPending: allCctxsPending,
				PastCctxsValue: sdkmath.NewInt(11).
					Mul(sdkmath.NewInt(1e18)),
				// 11 MUSE, increased value by 1 MUSE
				PendingCctxsValue:       sdkmath.NewInt(90).Mul(sdkmath.NewInt(1e18)), // 90 * 1 MUSE
				LowestPendingCctxHeight: 11,
			},
			output: ratelimiter.Output{ // should return missed cctxs only
				CctxsMap: map[int64][]*crosschaintypes.CrossChainTx{
					ethChainID: ethCctxsMissed,
					btcChainID: btcCctxsMissed,
				},
				CurrentWithdrawWindow: 100, // height [1, 100]
				CurrentWithdrawRate: sdkmath.NewInt(
					101e16,
				), // (11 + 90) / 100 = 1.01 MUSE/block (exceeds 0.99 MUSE/block)
				RateLimitExceeded: true,
			},
		},
		{
			name:   "rate limit is exceeded in wider window then the given sliding window 50",
			window: 50,
			rate:   sdkmath.NewUint(1e18), // 1 MUSE/block
			input: ratelimiter.Input{
				Height:       100,
				CctxsMissed:  allCctxsMissed,
				CctxsPending: allCctxsPending,
				PastCctxsValue: sdkmath.NewInt(
					0,
				), // no past cctx in height range [51, 100]
				PendingCctxsValue: sdkmath.NewInt(91).
					Mul(sdkmath.NewInt(1e18)),
				// 91 MUSE, increased value by 1 MUSE
				LowestPendingCctxHeight: 11,
			},
			output: ratelimiter.Output{
				CctxsMap: map[int64][]*crosschaintypes.CrossChainTx{
					ethChainID: ethCctxsMissed,
					btcChainID: btcCctxsMissed,
				},
				CurrentWithdrawWindow: 90, // [LowestPendingCctxHeight, Height] = [11, 100]
				CurrentWithdrawRate: sdkmath.NewInt(91).
					Mul(sdkmath.NewInt(1e18)).
					Quo(sdkmath.NewInt(90)),
				// 91 / 90 = 1.011111111111111111 MUSE/block
				RateLimitExceeded: true,
			},
		},
		{
			name:   "should not exceed rate limit if we wait for 1 more block",
			window: 50,
			rate:   sdkmath.NewUint(1e18), // 1 MUSE/block
			input: ratelimiter.Input{
				Height:       101,
				CctxsMissed:  allCctxsMissed,
				CctxsPending: allCctxsPending,
				PastCctxsValue: sdkmath.NewInt(
					0,
				), // no past cctx in height range [52, 101]
				PendingCctxsValue: sdkmath.NewInt(91).
					Mul(sdkmath.NewInt(1e18)),
				// 91 MUSE, increased value by 1 MUSE
				LowestPendingCctxHeight: 11,
			},
			output: ratelimiter.Output{
				CctxsMap: map[int64][]*crosschaintypes.CrossChainTx{
					ethChainID: ethCctxsAll,
					btcChainID: btcCctxsAll,
				},
				CurrentWithdrawWindow: 91, // [LowestPendingCctxHeight, Height] = [11, 101]
				CurrentWithdrawRate: sdkmath.NewInt(91).
					Mul(sdkmath.NewInt(1e18)).
					Quo(sdkmath.NewInt(91)),
				// 91 / 91 = 1.011 MUSE/block
				RateLimitExceeded: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ratelimiter.ApplyRateLimiter(&tt.input, tt.window, tt.rate)
			require.Equal(t, tt.output.CctxsMap, output.CctxsMap)
			require.Equal(t, tt.output.CurrentWithdrawWindow, output.CurrentWithdrawWindow)
			require.Equal(t, tt.output.CurrentWithdrawRate, output.CurrentWithdrawRate)
			require.Equal(t, tt.output.RateLimitExceeded, output.RateLimitExceeded)
		})
	}
}
