package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// rootCmd is the root for all commands
var rootCmd = &cobra.Command{
	Use:   "devops-school-bot",
	Short: "DevOps School Bot",
	Long:  "DevOps School Bot manage all students progress and provide the report",
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
