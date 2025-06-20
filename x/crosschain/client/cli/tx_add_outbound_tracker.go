package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/x/crosschain/types"
)

func CmdAddOutboundTracker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-outbound-tracker [chain] [nonce] [tx-hash]",
		Short: "Add an outbound tracker",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argChain, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			argNonce, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			argTxHash := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress().String()
			msg := types.NewMsgAddOutboundTracker(creator, argChain, argNonce, argTxHash)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
