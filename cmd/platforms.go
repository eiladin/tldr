package cmd

import (
	"log"

	"github.com/eiladin/tldr/internal/pipeline"
	"github.com/eiladin/tldr/pkg/context"
	"github.com/spf13/cobra"
)

// platformsCmd represents the platforms command
var listPlatformsCmd = &cobra.Command{
	Use:   "platforms",
	Short: "List available platforms.",
	Long:  `List available platforms.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := listPlatforms(args...)
		if err != nil {
			log.Fatalf(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(listPlatformsCmd)
}

func listPlatforms(args ...string) (*context.Context, error) {
	ctx := context.New()
	setupContext(ctx, args...)
	ctx.Operation = context.OperationListPlatforms
	return pipeline.Execute(ctx, pipeline.ListPlatformsPipeline)
}
