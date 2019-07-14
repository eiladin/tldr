package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// pwshCmd represents the pwsh command
var pwshCmd = &cobra.Command{
	Use:   "pwsh",
	Short: "Generates pwsh completion scripts",
	Long: `To load completion for powershell: 
  copy output to a file and load into your powershell profile
  For more information on powershell profiles see 'Get-Help about_Profiles'
	`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenPowerShellCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(pwshCmd)
}
