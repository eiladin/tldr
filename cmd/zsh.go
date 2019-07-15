package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// zshCmd represents the zsh command
var zshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generates zsh completion scripts",
	Long: `To load completion for oh-my-zsh:
  copy output to $ZSH_CUSTOM/plugins/tldr/_tldr

  mkdir -p $ZSH_CUSTOM/plugins/tldr &&
  tldr completion zsh > $ZSH_CUSTOM/plugins/tldr/_tldr

  Then define it in .zshrc as a plugin:
  plugins=(git tldr)
	`,
	Run: func(cmd *cobra.Command, args []string) {
		genZshCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(zshCmd)
}

func genZshCompletion(w io.Writer) {
	rootCmd.GenZshCompletion(w)
}
