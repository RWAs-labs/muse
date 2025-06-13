package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/app"
)

func UpgradeHandlerVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade-handler-version",
		Short: "Print the default upgrade handler version",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(app.GetDefaultUpgradeHandlerVersion())
		},
	}
}
