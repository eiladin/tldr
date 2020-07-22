package cmd

import (
	"github.com/spf13/cobra"
)

type completionCmd struct {
	cmd *cobra.Command
}

func newCompletionCmd() *completionCmd {
	var c = &completionCmd{}
	var cmd = &cobra.Command{
		Use:           "completion",
		Short:         "Generates completion scripts",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(
		newBashCmd().cmd,
		newPwshCmd().cmd,
		newZshCmd().cmd,
	)

	c.cmd = cmd
	return c
}
