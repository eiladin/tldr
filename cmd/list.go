package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/eiladin/tldr/cache"
	"github.com/eiladin/tldr/config"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all commands for the selected platform.",
	Long:  `list all commands for the selected platform.`,
	Run:   listPages,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("platform", "p", config.OSName(), "Platform to show usage for (run 'tldr platforms' to see available platforms)")
}

func listPages(cmd *cobra.Command, args []string) {
	platform, _ := cmd.Flags().GetString("platform")
	listPlatformPages(cache.DefaultSettings, platform, args...)
}

func listPlatformPages(settings cache.Cache, platform string, args ...string) {
	cache, err := cache.Create(settings.Remote, settings.TTL, settings.Location)
	if err != nil {
		log.Fatalf("ERROR: creating cache: %s", err)
	}
	pages, err := cache.ListPages(platform)
	if err != nil {
		log.Fatalf("ERROR: fetching pages for platform: %s", err)
	}
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()
	for i := 0; i < len(pages)-1; i += 2 {
		fmt.Fprintf(w, "%s\t%s\n", pages[i], pages[i+1])
	}
	if len(pages)%2 > 0 {
		fmt.Fprintf(w, "%s\n", pages[len(pages)-1])
	}
}
