package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/eiladin/tldr/cache"
	"github.com/eiladin/tldr/config"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCommandsCmd = &cobra.Command{
	Use:   "list-commands",
	Short: "list all commands for the selected platform.",
	Long:  `list all commands for the selected platform.`,
	Run:   listPages,
}

func init() {
	rootCmd.AddCommand(listCommandsCmd)
	listCommandsCmd.Flags().StringP("platform", "p", config.OSName(), "Platform to show usage for (run 'tldr platforms' to see available platforms)")
}

func listPages(cmd *cobra.Command, args []string) {
	platform, _ := cmd.Flags().GetString("platform")
	listPlatformPages(os.Stdout, cache.DefaultSettings, platform, args...)
}

func listPlatformPages(writer io.Writer, settings cache.Cache, platform string, args ...string) {
	cache, err := cache.Create(settings.Remote, settings.TTL, settings.Location)
	if err != nil {
		log.Fatalf("ERROR: creating cache: %s", err)
	}
	platformValid := cache.IsPlatformValid(platform)
	if !platformValid {
		availablePlatforms, _ := cache.AvailablePlatforms()
		log.Fatalf("ERROR: platform %s not found\nAvailable platforms: %s", platform, strings.Join(availablePlatforms, ", "))
	}
	pages, err := cache.ListPages(platform)
	if err != nil {
		log.Fatalf("ERROR: fetching pages for platform: %s", err)
	}
	w := tabwriter.NewWriter(writer, 8, 8, 0, '\t', 0)
	defer w.Flush()
	for i := 0; i < len(pages)-1; i += 2 {
		fmt.Fprintf(w, "%s\t%s\n", pages[i], pages[i+1])
	}
	if len(pages)%2 > 0 {
		fmt.Fprintf(w, "%s\n", pages[len(pages)-1])
	}
}