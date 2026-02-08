package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show tool version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("repocheck-go %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
