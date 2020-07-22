package cmd

import (
	"io"
	"log"
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
		genPwshCompletion(os.Stdout)
	},
}

func init() {
	completionCmd.AddCommand(pwshCmd)
}

func genPwshCompletion(w io.Writer) {
	err := rootCmd.GenPowerShellCompletion(w)
	if err != nil {
		log.Fatalf("ERROR: generating pwsh completion: %s", err)
	}
}
