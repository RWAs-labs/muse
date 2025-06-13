package e2etests

import (
	"time"

	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

const (
	startTimestampMetricName        = "museclient_last_start_timestamp_seconds"
	blockTimeLatencyMetricName      = "museclient_core_block_latency"
	blockTimeLatencySleepMetricName = "museclient_core_block_latency_sleep"
)

// TestMuseclientRestartHeight tests scheduling a museclient restart via operational flags
func TestMuseclientRestartHeight(r *runner.E2ERunner, _ []string) {
	_, err := r.Clients.Musecore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)

	currentHeight, err := r.Clients.Musecore.GetBlockHeight(r.Ctx)
	require.NoError(r, err)

	// schedule a restart for 5 blocks in the future
	restartHeight := currentHeight + 5
	updateMsg := observertypes.NewMsgUpdateOperationalFlags(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		observertypes.OperationalFlags{
			RestartHeight: restartHeight,
		},
	)

	_, err = r.MuseTxServer.BroadcastTx(utils.OperationalPolicyName, updateMsg)
	require.NoError(r, err)

	operationalFlagsRes, err := r.Clients.Musecore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)
	require.Equal(r, restartHeight, operationalFlagsRes.OperationalFlags.RestartHeight)

	originalStartTime, err := r.Clients.MuseclientMetrics.FetchGauge(startTimestampMetricName)
	require.NoError(r, err, "fetching museclient metric name")

	// wait for height above restart height
	// wait for a few extra block to account for shutdown and startup time
	require.Eventually(r, func() bool {
		height, err := r.Clients.Musecore.GetBlockHeight(r.Ctx)
		require.NoError(r, err)
		return height > restartHeight+3
	}, time.Minute, time.Second)

	currentStartTime, err := r.Clients.MuseclientMetrics.FetchGauge(startTimestampMetricName)
	require.NoError(r, err)

	require.Greater(r, currentStartTime, originalStartTime+1)
}

// TestMuseclientSignerOffset tests scheduling a museclient restart via operational flags
func TestMuseclientSignerOffset(r *runner.E2ERunner, _ []string) {
	startBlockTimeLatencySleep, err := r.Clients.MuseclientMetrics.FetchGauge(blockTimeLatencySleepMetricName)
	require.NoError(r, err)
	require.InDelta(r, 0, startBlockTimeLatencySleep, .01, "start block time latency should be 0")

	// get starting block time latency.
	// we need to ensure it's not zero (if museclient just finished a restart)
	var startBlockTimeLatency float64
	require.Eventually(r, func() bool {
		startBlockTimeLatency, err = r.Clients.MuseclientMetrics.FetchGauge(blockTimeLatencyMetricName)
		require.NoError(r, err)
		return startBlockTimeLatency > 1
	}, time.Second*15, time.Millisecond*100)

	desiredSignerBlockTimeOffset := time.Duration(startBlockTimeLatency*float64(time.Second)) + time.Millisecond*200

	updateMsg := observertypes.NewMsgUpdateOperationalFlags(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		observertypes.OperationalFlags{
			SignerBlockTimeOffset: &desiredSignerBlockTimeOffset,
		},
	)

	_, err = r.MuseTxServer.BroadcastTx(utils.OperationalPolicyName, updateMsg)
	require.NoError(r, err)

	operationalFlagsRes, err := r.Clients.Musecore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)
	require.InDelta(r, desiredSignerBlockTimeOffset, *(operationalFlagsRes.OperationalFlags.SignerBlockTimeOffset), .01)

	require.Eventually(r, func() bool {
		blockTimeLatencySleep, err := r.Clients.MuseclientMetrics.FetchGauge(blockTimeLatencySleepMetricName)
		if err != nil {
			return false
		}
		return blockTimeLatencySleep > .05
	}, time.Second*20, time.Second*1)
}
