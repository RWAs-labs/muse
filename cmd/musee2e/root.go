package main

import (
	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/cmd/musee2e/local"
)

const banner = `
         _             ____      
 _______| |_ __ _  ___|___ \ ___ 
|_  / _ \ __/ _  |/ _ \ __) / _ \
 / /  __/ || (_| |  __// __/  __/
/___\___|\__\__,_|\___|_____\___|
`

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "musee2e",
		Short: banner,
	}
	cmd.AddCommand(
		NewRunCmd(),
		NewBalancesCmd(),
		NewBitcoinAddressCmd(),
		NewListTestsCmd(),
		NewShowTSSCmd(),
		local.NewLocalCmd(),
		NewStressTestCmd(),
		NewInitCmd(),
		NewSetupBitcoinCmd(),
		NewPopulateAddressesCmd(),
	)

	return cmd
}
