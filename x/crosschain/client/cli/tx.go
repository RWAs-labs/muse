package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/x/crosschain/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdAddOutboundTracker(),
		CmdVoteGasPrice(),
		CmdVoteInbound(),
		CmdVoteOutbound(),
		CmdRemoveOutboundTracker(),
		CmdUpdateTss(),
		CmdMigrateTssFunds(),
		CmdAddInboundTracker(),
		CmdWhitelistERC20(),
		CmdAbortStuckCCTX(),
		CmdRefundAborted(),
		CmdRemoveInboundTracker(),
	)

	return cmd
}
