package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/pkg/coin"
	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/observer/types"
)

func TestKeeper_VoteOnInboundBallot(t *testing.T) {

	t.Run("fail if inbound not enabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: false,
		})

		_, _, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			chains.MuseChainPrivnet.ChainId,
			coin.CoinType_ERC20,
			sample.AccAddress(),
			"index",
			"inTxHash",
		)

		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInboundDisabled)
	})

	t.Run("fail if sender chain not supported", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{})

		_, _, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			chains.MuseChainPrivnet.ChainId,
			coin.CoinType_ERC20,
			sample.AccAddress(),
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrSupportedChains)

		// set the chain but not supported
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: false,
				},
			},
		})

		_, _, err = k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			chains.MuseChainPrivnet.ChainId,
			coin.CoinType_ERC20,
			sample.AccAddress(),
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{})

		_, _, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			chains.MuseChainPrivnet.ChainId,
			coin.CoinType_ERC20,
			sample.AccAddress(),
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrNotObserver)
	})

	t.Run("fail if receiver chain not supported", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		_, _, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			chains.MuseChainPrivnet.ChainId,
			coin.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrSupportedChains)

		// set the chain but not supported
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
				{
					ChainId:     chains.MuseChainPrivnet.ChainId,
					IsSupported: false,
				},
			},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		_, _, err = k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			chains.MuseChainPrivnet.ChainId,
			coin.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("fail if inbound contain MUSE but receiver chain doesn't support MUSE", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
				{
					ChainId:                  getValidEthChainIDWithIndex(t, 1),
					IsSupported:              true,
					MuseTokenContractAddress: "",
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		_, _, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			getValidEthChainIDWithIndex(t, 1),
			coin.CoinType_Muse,
			observer,
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrInvalidMuseCoinTypes)
	})

	t.Run("can add vote and create ballot", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
				{
					ChainId:     getValidEthChainIDWithIndex(t, 1),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		isFinalized, isNew, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			getValidEthChainIDWithIndex(t, 1),
			coin.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.NoError(t, err)

		// ballot should be finalized since there is only one observer
		require.True(t, isFinalized)
		require.True(t, isNew)
	})

	t.Run("fail if can not add vote", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
				{
					ChainId:     getValidEthChainIDWithIndex(t, 1),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)
		ballot := types.Ballot{
			Index:            "index",
			BallotIdentifier: "index",
			VoterList:        []string{observer},
			// already voted
			Votes:           []types.VoteType{types.VoteType_SuccessObservation},
			BallotStatus:    types.BallotStatus_BallotInProgress,
			BallotThreshold: sdkmath.LegacyNewDec(2),
		}
		k.SetBallot(ctx, &ballot)
		isFinalized, isNew, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			getValidEthChainIDWithIndex(t, 1),
			coin.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.Error(t, err)
		require.False(t, isFinalized)
		require.False(t, isNew)
	})

	t.Run("can add vote and create ballot without finalizing ballot", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		// threshold high enough to not finalize ballot
		threshold, err := sdkmath.LegacyNewDecFromStr("0.7")
		require.NoError(t, err)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:         getValidEthChainIDWithIndex(t, 0),
					IsSupported:     true,
					BallotThreshold: threshold,
				},
				{
					ChainId:         getValidEthChainIDWithIndex(t, 1),
					IsSupported:     true,
					BallotThreshold: threshold,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{
				observer,
				sample.AccAddress(),
			},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		isFinalized, isNew, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			getValidEthChainIDWithIndex(t, 1),
			coin.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.NoError(t, err)

		// ballot should be finalized since there is only one observer
		require.False(t, isFinalized)
		require.True(t, isNew)
	})

	t.Run("can add vote to an existing ballot", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
				{
					ChainId:     getValidEthChainIDWithIndex(t, 1),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		// set a ballot
		threshold, err := sdkmath.LegacyNewDecFromStr("0.7")
		require.NoError(t, err)
		ballot := types.Ballot{
			Index:            "index",
			BallotIdentifier: "index",
			VoterList: []string{
				sample.AccAddress(),
				sample.AccAddress(),
				observer,
				sample.AccAddress(),
				sample.AccAddress(),
			},
			Votes:           types.CreateVotes(5),
			ObservationType: types.ObservationType_InboundTx,
			BallotThreshold: threshold,
			BallotStatus:    types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		isFinalized, isNew, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			getValidEthChainIDWithIndex(t, 1),
			coin.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.NoError(t, err)

		// ballot should not be finalized as the threshold is not reached
		require.False(t, isFinalized)
		require.False(t, isNew)
	})

	t.Run("can add vote to an existing ballot and finalize ballot", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMocksAll)

		observer := sample.AccAddress()
		stakingMock := keepertest.GetObserverStakingMock(t, k)
		slashingMock := keepertest.GetObserverSlashingMock(t, k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled: true,
		})
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				{
					ChainId:     getValidEthChainIDWithIndex(t, 0),
					IsSupported: true,
				},
				{
					ChainId:     getValidEthChainIDWithIndex(t, 1),
					IsSupported: true,
				},
			},
		})
		k.SetObserverSet(ctx, types.ObserverSet{
			ObserverList: []string{observer},
		})
		stakingMock.MockGetValidator(sample.Validator(t, sample.Rand()))
		slashingMock.MockIsTombstoned(false)

		// set a ballot
		threshold, err := sdkmath.LegacyNewDecFromStr("0.1")
		require.NoError(t, err)
		ballot := types.Ballot{
			Index:            "index",
			BallotIdentifier: "index",
			VoterList: []string{
				observer,
				sample.AccAddress(),
				sample.AccAddress(),
			},
			Votes:           types.CreateVotes(3),
			ObservationType: types.ObservationType_InboundTx,
			BallotThreshold: threshold,
			BallotStatus:    types.BallotStatus_BallotInProgress,
		}
		k.SetBallot(ctx, &ballot)

		isFinalized, isNew, err := k.VoteOnInboundBallot(
			ctx,
			getValidEthChainIDWithIndex(t, 0),
			getValidEthChainIDWithIndex(t, 1),
			coin.CoinType_ERC20,
			observer,
			"index",
			"inTxHash",
		)
		require.NoError(t, err)

		// ballot should not be finalized as the threshold is not reached
		require.True(t, isFinalized)
		require.False(t, isNew)
	})
}
