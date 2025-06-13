package main

import (
	"os"

	"github.com/fatih/color"

	_ "github.com/RWAs-labs/muse/pkg/sdkconfig/default"
)

func main() {
	// enable color output
	color.NoColor = false

	// initialize root command
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
