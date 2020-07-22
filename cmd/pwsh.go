package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

type pwshCmd struct {
	cmd *cobra.Command
}

func newPwshCmd() *pwshCmd {
	c := &pwshCmd{}
	cmd := &cobra.Command{
		Use:   "pwsh",
		Short: "Generates pwsh completion scripts",
		Long: `To load completion for powershell: 
		copy output to a file and load into your powershell profile
		For more information on powershell profiles see 'Get-Help about_Profiles'
		`,
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, args []string) {
			genPwshCompletion(cmd, os.Stdout)
		},
	}
	c.cmd = cmd
	return c
}

func genPwshCompletion(cmd *cobra.Command, w io.Writer) {
	cmd.Root().GenPowerShellCompletion(w)
}
