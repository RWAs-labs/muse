package main

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func AddrConversionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addr-conversion [muse address]",
		Short: "convert a muse1xxx address to validator operator address musevaloper1xxx",
		Long: `
read a muse1xxx or musevaloper1xxx address and convert it to the other type;
it always outputs three lines; the first line is the muse1xxx address, the second line is the musevaloper1xxx address
and the third line is the ethereum address.
			`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err == nil {
				valAddr := sdk.ValAddress(addr.Bytes())
				fmt.Printf("%s\n", addr.String())
				fmt.Printf("%s\n", valAddr.String())
				fmt.Printf("%s\n", ethcommon.BytesToAddress(addr.Bytes()).String())
				return nil
			}
			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err == nil {
				addr := sdk.AccAddress(valAddr.Bytes())
				fmt.Printf("%s\n", addr.String())
				fmt.Printf("%s\n", valAddr.String())
				fmt.Printf("%s\n", ethcommon.BytesToAddress(addr.Bytes()).String())
				return nil
			}
			ethAddr := ethcommon.HexToAddress(args[0])
			if ethAddr != (ethcommon.Address{}) {
				addr := sdk.AccAddress(ethAddr.Bytes())
				valAddr := sdk.ValAddress(addr.Bytes())
				fmt.Printf("%s\n", addr.String())
				fmt.Printf("%s\n", valAddr.String())
				fmt.Printf("%s\n", ethAddr.String())
				return nil
			}
			return fmt.Errorf("invalid address: %s", args[0])
		},
	}
	return cmd
}
