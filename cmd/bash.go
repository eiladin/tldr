package cmd

import (
	"io"
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
		genBashCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(bashCmd)
}

func genBashCompletion(w io.Writer) {
	err := rootCmd.GenBashCompletion(w)
	if err != nil {
		logFatalf("ERROR: generating bash completion: %s", err)
	}
}
