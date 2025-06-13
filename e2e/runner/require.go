package runner

import (
	"fmt"
	"math/big"

	"github.com/RWAs-labs/protocol-contracts/pkg/mrc20.sol"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/testutil/sample"
	crosschaintypes "github.com/RWAs-labs/muse/x/crosschain/types"
	emissionstypes "github.com/RWAs-labs/muse/x/emissions/types"
	"github.com/RWAs-labs/muse/x/observer/types"
)

// EnsureNoTrackers ensures that there are no trackers left on musecore
func (r *E2ERunner) EnsureNoTrackers() {
	// get all trackers
	res, err := r.CctxClient.OutTxTrackerAll(
		r.Ctx,
		&crosschaintypes.QueryAllOutboundTrackerRequest{},
	)
	require.NoError(r, err)
	require.Empty(r, res.OutboundTracker, "there should be no trackers at the end of the test")
}

// EnsureZeroBalanceAddressMEVM ensures that the balance of the restricted address is zero in the MEVM
func (r *E2ERunner) EnsureZeroBalanceOnRestrictedAddressMEVM() {
	restrictedAddress := ethcommon.HexToAddress(sample.RestrictedEVMAddressTest)

	// ensure MUSE balance is zero
	balance, err := r.WMuse.BalanceOf(&bind.CallOpts{}, restrictedAddress)
	require.NoError(r, err)
	require.Zero(r, balance.Cmp(big.NewInt(0)), "the wMUSE balance of the address should be zero")

	// ensure MRC20 ETH balance is zero
	ensureMRC20ZeroBalance(r, r.ETHMRC20, restrictedAddress)

	// ensure MRC20 ERC20 balance is zero
	ensureMRC20ZeroBalance(r, r.ERC20MRC20, restrictedAddress)

	// ensure MRC20 BTC balance is zero
	ensureMRC20ZeroBalance(r, r.BTCMRC20, restrictedAddress)

	// ensure MRC20 SOL balance is zero
	ensureMRC20ZeroBalance(r, r.SOLMRC20, restrictedAddress)
}

// ensureMRC20ZeroBalance ensures that the balance of the MRC20 token is zero on given address
func ensureMRC20ZeroBalance(r *E2ERunner, mrc20 *mrc20.MRC20, address ethcommon.Address) {
	balance, err := mrc20.BalanceOf(&bind.CallOpts{}, address)
	require.NoError(r, err)

	mrc20Name, err := mrc20.Name(&bind.CallOpts{})
	require.NoError(r, err)
	require.Zero(
		r,
		balance.Cmp(big.NewInt(0)),
		fmt.Sprintf("the balance of address %s should be zero on MRC20: %s", address, mrc20Name),
	)
}

// EnsureNoStaleBallots ensures that there are no stale ballots left on the chain.
func (r *E2ERunner) EnsureNoStaleBallots() {
	ballotsRes, err := r.ObserverClient.Ballots(r.Ctx, &types.QueryBallotsRequest{})
	require.NoError(r, err)
	currentBlockHeight, err := r.Clients.Musecore.GetBlockHeight(r.Ctx)
	require.NoError(r, err)
	emissionsParams, err := r.EmissionsClient.Params(r.Ctx, &emissionstypes.QueryParamsRequest{})
	require.NoError(r, err)
	staleBlockStart := currentBlockHeight - emissionsParams.Params.BallotMaturityBlocks
	if len(ballotsRes.Ballots) < 1 {
		return
	}

	firstBallotCreationHeight := int64(0)

	for _, ballot := range ballotsRes.Ballots {
		if ballot.IsFinalized() {
			firstBallotCreationHeight = ballot.BallotCreationHeight
			break
		}
	}

	if firstBallotCreationHeight == 0 {
		return
	}
	// Log data for debugging
	if firstBallotCreationHeight < staleBlockStart {
		r.Logger.Error(
			"First finalized ballot creation height: %d is less than stale block start: %d",
			firstBallotCreationHeight,
			staleBlockStart,
		)
		for _, ballot := range ballotsRes.Ballots {
			r.Logger.Error(
				"Ballot: %s Creation Height %d BallotStatus %s\n",
				ballot.BallotIdentifier,
				ballot.BallotCreationHeight,
				ballot.BallotStatus,
			)
		}
		r.Logger.Error("Block Maturity Height: %d", emissionsParams.Params.BallotMaturityBlocks)
	}
	require.GreaterOrEqual(r, firstBallotCreationHeight, staleBlockStart, "there should be no stale ballots")
}
