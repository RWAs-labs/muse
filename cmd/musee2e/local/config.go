package local

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/RWAs-labs/muse/e2e/config"
)

// GetConfig returns config from file from the command line flag
func GetConfig(cmd *cobra.Command) (config.Config, error) {
	configFile, err := cmd.Flags().GetString(FlagConfigFile)
	if err != nil {
		return config.Config{}, fmt.Errorf("--config is a required parameter")
	}

	configFile, err = filepath.Abs(configFile)
	if err != nil {
		return config.Config{}, err
	}

	return config.ReadConfig(configFile, true)
}
