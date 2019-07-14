package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/eiladin/tldr/cache"
	"github.com/spf13/cobra"
)

// platformsCmd represents the platforms command
var listPlatformsCmd = &cobra.Command{
	Use:   "list-platforms",
	Short: "List available platforms.",
	Long:  `List available platforms.`,
	Run:   listPlatforms,
}

func init() {
	rootCmd.AddCommand(listPlatformsCmd)
}

func listPlatforms(cmd *cobra.Command, args []string) {
	listAvailablePlatforms(os.Stdout, cache.DefaultSettings, args...)
}

func listAvailablePlatforms(writer io.Writer, settings cache.Cache, args ...string) {
	cache, err := cache.Create(settings.Remote, settings.TTL, settings.Location)
	if err != nil {
		log.Fatalf("ERROR: Creating cache: %s", err)
	}
	platforms, err := cache.AvailablePlatforms()
	if err != nil {
		log.Fatalf("ERROR: Getting platforms: %s", err)
	}
	platformList := strings.Join(platforms, ", ")
	fmt.Fprintf(writer, "Available Platforms: %s\n", platformList)
}