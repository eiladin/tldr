package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

type zshCmd struct {
	cmd *cobra.Command
}

func newZshCmd() *zshCmd {
	c := &zshCmd{}
	cmd := &cobra.Command{
		Use:   "zsh",
		Short: "Generates zsh completion scripts",
		Long: `To load completion for oh-my-zsh:
		copy output to $ZSH_CUSTOM/plugins/tldr/_tldr
	
		$ mkdir -p $ZSH_CUSTOM/plugins/tldr && \
		tldr completion zsh > $ZSH_CUSTOM/plugins/tldr/_tldr
	
		Then define it in .zshrc as a plugin:
		plugins=(git tldr)
		`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, args []string) {
			genZshCompletion(cmd, os.Stdout)
		},
	}

	c.cmd = cmd
	return c
}

func genZshCompletion(cmd *cobra.Command, w io.Writer) {
	cmd.Root().GenZshCompletion(w)
}
