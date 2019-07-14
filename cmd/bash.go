package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// bashCmd represents the bash command
var bashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion scripts",
	Long: `To load completion run:
  . <(tldr completion)
  
  To configure your bash shell to load completions for each session add to your bashrc
  
  # ~/.bashrc or ~/.profile
  . <(tldr completion bash)
	`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(bashCmd)
}
