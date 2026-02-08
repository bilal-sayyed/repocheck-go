package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "repocheck",
	Short: "A CLI tool to analyze source code repositories",
	Long: `repocheck-go is a CLI tool that analyzes a source code repository and reports on its 
readiness, hygiene, and risk level.

It runs locally and checks for things like README presence, dependency management, 
and basic configuration hygiene.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .repocheck.yaml in repo root)")
}
