package keeper_test

import (
	"testing"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/x/observer/types"
	"github.com/stretchr/testify/require"
)

func TestOperationalFlagsKeeper(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)

	// found should be false on first run
	// we do not store a genesis value for the flags since
	// all fields default to zero value
	_, found := k.GetOperationalFlags(ctx)
	require.False(t, found)

	restartHeight := int64(100)

	k.SetOperationalFlags(ctx, types.OperationalFlags{
		RestartHeight: restartHeight,
	})

	operationalFlags, found := k.GetOperationalFlags(ctx)
	require.True(t, found)
	require.Equal(t, restartHeight, operationalFlags.RestartHeight)
}
