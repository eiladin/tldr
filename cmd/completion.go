package cmd

import (
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates completion scripts",
	Long:  `Generate shell completion scripts`,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
