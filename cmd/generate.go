package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"fixtures_from_db/config"
)

var (
	configFilePath string
)

var GenerateCmd = &cobra.Command{
	Use: "generate",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadFromFile(configFilePath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Printf("%+v", cfg)

		return nil
	},
}

func init() {
	GenerateCmd.Flags().StringVarP(&configFilePath, "config-file-path", "f", "config.yaml", "Path to config file")
}
