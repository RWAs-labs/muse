package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/RWAs-labs/muse/testutil/keeper"
	"github.com/RWAs-labs/muse/testutil/sample"
)

func TestKeeper_SetPolicies(t *testing.T) {
	k, ctx := keepertest.AuthorityKeeper(t)
	policies := sample.Policies()

	_, found := k.GetPolicies(ctx)
	require.False(t, found)

	k.SetPolicies(ctx, policies)

	// Check policy is set
	got, found := k.GetPolicies(ctx)
	require.True(t, found)
	require.Equal(t, policies, got)

	// Can set policies again
	newPolicies := sample.Policies()
	require.NotEqual(t, policies, newPolicies)
	k.SetPolicies(ctx, newPolicies)
	got, found = k.GetPolicies(ctx)
	require.True(t, found)
	require.Equal(t, newPolicies, got)
}
