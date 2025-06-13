package base_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RWAs-labs/muse/museclient/chains/base"
	"github.com/RWAs-labs/muse/museclient/config"
	"github.com/RWAs-labs/muse/museclient/testutils/mocks"
	"github.com/RWAs-labs/muse/pkg/chains"
	"github.com/RWAs-labs/muse/testutil/sample"
)

// createSigner creates a new signer for testing
func createSigner(t *testing.T) *base.Signer {
	// constructor parameters
	chain := chains.Ethereum
	tss := mocks.NewTSS(t)
	logger := base.DefaultLogger()

	// create signer
	return base.NewSigner(chain, tss, logger)
}

func TestNewSigner(t *testing.T) {
	signer := createSigner(t)
	require.NotNil(t, signer)
}

func Test_BeingReportedFlag(t *testing.T) {
	signer := createSigner(t)

	// hash to be reported
	hash := "0x1234"
	alreadySet := signer.SetBeingReportedFlag(hash)
	require.False(t, alreadySet)

	// set reported outbound again and check
	alreadySet = signer.SetBeingReportedFlag(hash)
	require.True(t, alreadySet)

	// clear reported outbound and check again
	signer.ClearBeingReportedFlag(hash)
	alreadySet = signer.SetBeingReportedFlag(hash)
	require.False(t, alreadySet)
}

func Test_PassesCompliance(t *testing.T) {
	signer := createSigner(t)

	// create config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return false for restricted CCTX", func(t *testing.T) {
		cctx := sample.CrossChainTxV2(t, "abcd")
		cfg.ComplianceConfig.RestrictedAddresses = []string{cctx.InboundParams.Sender}
		config.SetRestrictedAddressesFromConfig(cfg)

		require.False(t, signer.PassesCompliance(cctx))
	})
	t.Run("should return true for non restricted CCTX", func(t *testing.T) {
		cctx := sample.CrossChainTxV2(t, "abcd")
		cfg.ComplianceConfig.RestrictedAddresses = []string{sample.EthAddress().Hex()}
		config.SetRestrictedAddressesFromConfig(cfg)

		require.True(t, signer.PassesCompliance(cctx))
	})
}
