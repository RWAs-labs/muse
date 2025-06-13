package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/x/observer/types"
)

func CmdGetTssAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-tss-address [bitcoinChainId]]",
		Short: "Query current tss address",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryGetTssAddressRequest{}
			if len(args) == 1 {
				bitcoinChainID, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return err
				}
				params.BitcoinChainId = bitcoinChainID
			}

			res, err := queryClient.GetTssAddress(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetTssAddressByFinalizedMuseHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-historical-tss-address [finalizedMuseHeight] [bitcoinChainId]",
		Short: "Query tss address by finalized muse height (for historical tss addresses)",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			finalizedMuseHeight, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			params := &types.QueryGetTssAddressByFinalizedHeightRequest{
				FinalizedMuseHeight: finalizedMuseHeight,
			}
			if len(args) == 2 {
				bitcoinChainID, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					return err
				}
				params.BitcoinChainId = bitcoinChainID
			}

			res, err := queryClient.GetTssAddressByFinalizedHeight(cmd.Context(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
