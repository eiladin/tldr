package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

type bashCmd struct {
	cmd *cobra.Command
}

func newBashCmd() *bashCmd {
	c := &bashCmd{}
	cmd := &cobra.Command{
		Use:   "bash",
		Short: "Generates bash completion scripts",
		Long: `To load completion run:
		$ source <(tldr completion)
		
		To configure your bash shell to load completions for each session add to your bashrc
		
		# ~/.bashrc or ~/.profile
		source <(tldr completion bash)
		`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, args []string) {
			genBashCompletion(cmd, os.Stdout)
		},
	}
	c.cmd = cmd
	return c
}

func genBashCompletion(cmd *cobra.Command, w io.Writer) {
	cmd.Root().GenBashCompletion(w)
}
