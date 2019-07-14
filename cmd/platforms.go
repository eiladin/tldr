package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/eiladin/tldr/cache"
	"github.com/spf13/cobra"
)

// platformsCmd represents the platforms command
var platformsCmd = &cobra.Command{
	Use:   "platforms",
	Short: "List available platforms.",
	Long:  `List available platforms.`,
	Run:   listPlatforms,
}

func init() {
	rootCmd.AddCommand(platformsCmd)
}

func listPlatforms(cmd *cobra.Command, args []string) {
	listAvailablePlatforms(cache.DefaultSettings, args...)
}

func listAvailablePlatforms(settings cache.Cache, args ...string) {
	cache, err := cache.Create(settings.Remote, settings.TTL, settings.Location)
	if err != nil {
		log.Fatalf("ERROR: Creating cache: %s", err)
	}
	platforms, err := cache.AvailablePlatforms()
	if err != nil {
		log.Fatalf("ERROR: Getting platforms: %s", err)
	}
	platformList := strings.Join(platforms, ", ")
	fmt.Printf("Available Platforms: %s\n", platformList)
}
