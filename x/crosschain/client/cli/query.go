package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(_ string) *cobra.Command {
	// Group crosschain queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdListOutboundTracker(),
		CmdShowOutboundTracker(),
		CmdListGasPrice(),
		CmdShowGasPrice(),

		CmdListSend(),
		CmdShowSend(),
		CmdLastMuseHeight(),
		CmdInboundHashToCctxData(),
		CmdListInboundHashToCctx(),
		CmdShowInboundHashToCctx(),

		CmdPendingCctx(),
		CmdListInboundTrackerByChain(),
		CmdListInboundTrackers(),
		CmdShowInboundTracker(),
		CmdGetMuseAccounting(),
		CmdListPendingCCTXWithinRateLimit(),

		CmdShowUpdateRateLimiterFlags(),
	)

	return cmd
}
