package cmd

import (
	"log"

	"github.com/eiladin/tldr/internal/config"
	"github.com/eiladin/tldr/internal/pipeline"
	"github.com/eiladin/tldr/pkg/context"
	"github.com/spf13/cobra"
)

var listCommandsCmd = &cobra.Command{
	Use:   "commands",
	Short: "list all commands for the selected platform.",
	Long:  `list all commands for the selected platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := listCommands(args...)
		if err != nil {
			log.Fatalf(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(listCommandsCmd)
	listCommandsCmd.Flags().StringVarP(&opts.platform, "platform", "p", config.CurrentPlatform(), "Platform to show usage for (run 'tldr platforms' to see available platforms)")
}

func listCommands(args ...string) (*context.Context, error) {
	ctx := context.New()
	setupContext(ctx, args...)
	ctx.Operation = context.OperationListCommands
	return pipeline.Execute(ctx, pipeline.ListCommandsPipeline)
}
