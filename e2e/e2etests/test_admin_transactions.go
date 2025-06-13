package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/e2e/runner"
	"github.com/RWAs-labs/muse/e2e/utils"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	"github.com/RWAs-labs/muse/testutil/sample"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	observertypes "github.com/RWAs-labs/muse/x/observer/types"
)

// TestCriticalAdminTransactions tests critical admin transactions that are the most used on mainnet .
// The complete list is
// MsgUpdateChainParams
// MsgRefundAbortedCCTX
// MsgEnableCCTX
// MsgDisableCCTX
// MsgUpdateGasPriceIncreaseFlags
// MsgAddInboundTracker
// MsgUpdateMRC20LiquidityCap
// MsgDeploySystemContracts
// MsgWhitelistERC20
// MsgPauseMRC20
// MsgMigrateTssFunds
// MsgUpdateTssAddress
//
//	However, the transactions other than `AddToInboundTracker` and `UpdateGasPriceIncreaseFlags` have already been used in other tests.
func TestCriticalAdminTransactions(r *runner.E2ERunner, _ []string) {
	TestAddToInboundTracker(r)
	TestUpdateGasPriceIncreaseFlags(r)
}

func TestUpdateGasPriceIncreaseFlags(r *runner.E2ERunner) {
	// Set default flags on musecore
	defaultFlags := observertypes.DefaultGasPriceIncreaseFlags
	msgGasPriceFlags := observertypes.NewMsgUpdateGasPriceIncreaseFlags(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		defaultFlags,
	)
	_, err := r.MuseTxServer.BroadcastTx(utils.OperationalPolicyName, msgGasPriceFlags)
	require.NoError(r, err)

	// create a new set of flag values by incrementing the epoch length by 1
	defaultFlagsUpdated := defaultFlags
	defaultFlagsUpdated.EpochLength = defaultFlags.EpochLength + 1

	// Update the flags on musecore with the new values
	msgGasPriceFlags = observertypes.NewMsgUpdateGasPriceIncreaseFlags(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		defaultFlagsUpdated,
	)
	_, err = r.MuseTxServer.BroadcastTx(utils.OperationalPolicyName, msgGasPriceFlags)
	require.NoError(r, err)

	r.WaitForBlocks(1)

	// Verify that the flags have been updated
	flags, err := r.ObserverClient.CrosschainFlags(r.Ctx, &observertypes.QueryGetCrosschainFlagsRequest{})
	require.NoError(r, err)
	require.Equal(r, defaultFlagsUpdated.EpochLength, flags.CrosschainFlags.GasPriceIncreaseFlags.EpochLength)
}

func TestAddToInboundTracker(r *runner.E2ERunner) {
	chainEth := chains.GoerliLocalnet
	chainBtc := chains.BitcoinRegtest
	msgEth := crosschaintypes.NewMsgAddInboundTracker(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		chainEth.ChainId,
		coin.CoinType_Gas,
		sample.Hash().Hex(),
	)
	_, err := r.MuseTxServer.BroadcastTx(utils.EmergencyPolicyName, msgEth)
	require.NoError(r, err)

	msgBtc := crosschaintypes.NewMsgAddInboundTracker(
		r.MuseTxServer.MustGetAccountAddressFromName(utils.EmergencyPolicyName),
		chainBtc.ChainId,
		coin.CoinType_Gas,
		sample.BtcHash().String(),
	)

	_, err = r.MuseTxServer.BroadcastTx(utils.EmergencyPolicyName, msgBtc)
	require.NoError(r, err)

	r.WaitForBlocks(1)

	tracker, err := r.CctxClient.InboundTracker(r.Ctx, &crosschaintypes.QueryInboundTrackerRequest{
		ChainId: msgEth.ChainId,
		TxHash:  msgEth.TxHash,
	})
	require.NoError(r, err)
	require.NotNil(r, tracker)
	require.Equal(r, msgEth.TxHash, tracker.InboundTracker.TxHash)

	tracker, err = r.CctxClient.InboundTracker(r.Ctx, &crosschaintypes.QueryInboundTrackerRequest{
		ChainId: msgBtc.ChainId,
		TxHash:  msgBtc.TxHash,
	})
	require.NoError(r, err)
	require.NotNil(r, tracker)
	require.Equal(r, msgBtc.TxHash, tracker.InboundTracker.TxHash)
}
