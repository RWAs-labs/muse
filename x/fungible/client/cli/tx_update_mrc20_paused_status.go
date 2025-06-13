package cli

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/x/fungible/types"
)

func CmdPauseMRC20() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pause-mrc20 [contractAddress1, contractAddress2, ...]",
		Short:   "Broadcast message PauseMRC20",
		Example: `musecored tx fungible pause-mrc20 "0xece40cbB54d65282c4623f141c4a8a0bE7D6AdEc, 0xece40cbB54d65282c4623f141c4a8a0bEjgksncf" `,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddressList := strings.Split(strings.TrimSpace(args[0]), ",")

			msg := types.NewMsgPauseMRC20(
				clientCtx.GetFromAddress().String(),
				contractAddressList,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUnpauseMRC20() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unpause-mrc20 [contractAddress1, contractAddress2, ...]",
		Short:   "Broadcast message UnpauseMRC20",
		Example: `musecored tx fungible unpause-mrc20 "0xece40cbB54d65282c4623f141c4a8a0bE7D6AdEc, 0xece40cbB54d65282c4623f141c4a8a0bEjgksncf" `,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			contractAddressList := strings.Split(strings.TrimSpace(args[0]), ",")

			msg := types.NewMsgUnpauseMRC20(
				clientCtx.GetFromAddress().String(),
				contractAddressList,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
