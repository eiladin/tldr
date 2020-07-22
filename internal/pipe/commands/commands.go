package commands

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/eiladin/tldr/pkg/context"
)

var errListingPages = errors.New("unable to list pages")

type Pipe struct{}

func (Pipe) String() string {
	return "listing commands"
}

func (Pipe) Run(ctx *context.Context) error {
	if strings.ToLower(ctx.Platform) == "all" {
		pgs := make([]string, 0)
		for _, p := range ctx.AvailablePlatforms {
			ppg, err := getPages(ctx, p)
			if err != nil {
				return err
			}
			for _, pg := range ppg {
				pgs = append(pgs, pg)
			}
		}
		pgs = unique(pgs)
		sort.Strings(pgs)
		print(ctx, pgs)
		return nil
	} else {
		pgs, err := getPages(ctx, ctx.Platform)
		if err != nil {
			return err
		}
		print(ctx, pgs)
		return nil
	}
}

func print(ctx *context.Context, pages []string) {
	w := tabwriter.NewWriter(ctx.Writer, 8, 8, 0, '\t', 0)
	defer w.Flush()
	for i := 0; i < len(pages); i++ {
		fmt.Fprintf(w, "%s\n", pages[i])
	}
}

func getPages(ctx *context.Context, platform string) ([]string, error) {
	dir := path.Join(ctx.Cache.Location, ctx.PagesDirectory, platform)
	pages := []os.FileInfo{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ctx.PageSuffix) {
			pages = append(pages, f)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cache: %s: %s", err, errListingPages.Error())
	}

	names := make([]string, len(pages))
	for i, page := range pages {
		name := page.Name()
		names[i] = name[:len(name)-3]
	}
	return names, nil
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	r := []string{}
	for _, s := range stringSlice {
		if _, v := keys[s]; !v {
			keys[s] = true
			r = append(r, s)
		}
	}
	return r
}
