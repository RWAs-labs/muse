package cli

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/x/fungible/types"
)

func CmdUpdateMRC20LiquidityCap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-mrc20-liquidity-cap [mrc20] [liquidity-cap]",
		Short: "Broadcast message UpdateMRC20LiquidityCap",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			newCap := math.NewUintFromString(args[1])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgUpdateMRC20LiquidityCap(
				clientCtx.GetFromAddress().String(),
				args[0],
				newCap,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
