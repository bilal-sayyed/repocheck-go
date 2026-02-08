package cmd

import (
	"fmt"
	"os"

	"github.com/bilal-sayyed/repocheck-go/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Local configuration management",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize .repocheck.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		defaultCfg := config.Config{
			IgnorePaths: []string{"node_modules", "vendor"},
			StrictMode:  false,
		}

		data, err := yaml.Marshal(defaultCfg)
		if err != nil {
			fmt.Println("Error generating config:", err)
			return
		}

		err = os.WriteFile(".repocheck.yaml", data, 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
		fmt.Println("Created .repocheck.yaml")
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			fmt.Println("No config file found or error loading it.")
			return
		}

		data, _ := yaml.Marshal(cfg)
		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
}
