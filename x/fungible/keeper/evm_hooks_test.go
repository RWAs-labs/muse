package keeper_test

import (
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
	"github.com/RWAs-labs/muse/x/fungible/types"
)

func TestKeeper_CheckPausedMRC20(t *testing.T) {
	addrUnpausedMRC20A, addrUnpausedMRC20B, addrUnpausedMRC20C, addrPausedMRC20 :=
		sample.EthAddress(),
		sample.EthAddress(),
		sample.EthAddress(),
		sample.EthAddress()

	tt := []struct {
		name    string
		receipt *ethtypes.Receipt
		wantErr bool
	}{
		{
			name:    "should pass if receipt is nil",
			receipt: nil,
			wantErr: false,
		},
		{
			name: "should pass if receipt is empty",
			receipt: &ethtypes.Receipt{
				Logs: []*ethtypes.Log{},
			},
			wantErr: false,
		},
		{
			name: "should pass if receipt contains unpaused MRC20 and non MRC20 addresses",
			receipt: &ethtypes.Receipt{
				Logs: []*ethtypes.Log{
					{
						Address: sample.EthAddress(),
					},
					{
						Address: addrUnpausedMRC20A,
					},
					{
						Address: addrUnpausedMRC20B,
					},
					{
						Address: addrUnpausedMRC20C,
					},
					{
						Address: addrUnpausedMRC20A,
					},
					{
						Address: addrUnpausedMRC20A,
					},
					nil,
					{
						Address: sample.EthAddress(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "should fail if receipt contains paused MRC20 and non MRC20 addresses",
			receipt: &ethtypes.Receipt{
				Logs: []*ethtypes.Log{
					{
						Address: sample.EthAddress(),
					},
					{
						Address: addrUnpausedMRC20A,
					},
					{
						Address: addrUnpausedMRC20B,
					},
					{
						Address: addrUnpausedMRC20C,
					},
					{
						Address: addrPausedMRC20,
					},
					{
						Address: sample.EthAddress(),
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			k, ctx, _, _ := keepertest.FungibleKeeper(t)

			assertUnpaused := func(mrc20 string) {
				fc, found := k.GetForeignCoins(ctx, mrc20)
				require.True(t, found)
				require.False(t, fc.Paused)
			}
			assertPaused := func(mrc20 string) {
				fc, found := k.GetForeignCoins(ctx, mrc20)
				require.True(t, found)
				require.True(t, fc.Paused)
			}

			// setup MRC20
			k.SetForeignCoins(ctx, sample.ForeignCoins(t, addrUnpausedMRC20A.Hex()))
			k.SetForeignCoins(ctx, sample.ForeignCoins(t, addrUnpausedMRC20B.Hex()))
			k.SetForeignCoins(ctx, sample.ForeignCoins(t, addrUnpausedMRC20C.Hex()))
			pausedMRC20 := sample.ForeignCoins(t, addrPausedMRC20.Hex())
			pausedMRC20.Paused = true
			k.SetForeignCoins(ctx, pausedMRC20)

			// check paused status
			assertUnpaused(addrUnpausedMRC20A.Hex())
			assertUnpaused(addrUnpausedMRC20B.Hex())
			assertUnpaused(addrUnpausedMRC20C.Hex())
			assertPaused(addrPausedMRC20.Hex())

			// process test
			h := k.EVMHooks()
			err := h.PostTxProcessing(ctx, nil, tc.receipt)
			if tc.wantErr {
				require.ErrorIs(t, err, types.ErrPausedMRC20)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
