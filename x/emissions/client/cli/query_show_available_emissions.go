package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/x/emissions/types"
)

func CmdShowAvailableEmissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-available-emissions [address]",
		Short: "Query show-available-emissions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryShowAvailableEmissionsRequest{

				Address: reqAddress,
			}

			res, err := queryClient.ShowAvailableEmissions(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
