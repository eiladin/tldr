package page

import (
	"fmt"
	"os"
	"path"

	"github.com/eiladin/tldr/internal/pipe"
	"github.com/eiladin/tldr/pkg/context"
)

type Pipe struct{}

func (Pipe) String() string {
	return "getting page"
}

func (Pipe) Run(ctx *context.Context) error {
	if !ctx.Random {
		platforms := []string{ctx.Platform, "common"}
		if err := getPage(platforms, ctx); err != nil {
			return err
		}
	} else {
		pipe.Skip("serving a random page")
	}
	return nil
}

type platformPath struct {
	Platform string
	Path     string
}

func makePath(ctx *context.Context, platform string) string {
	return path.Join(ctx.Cache.Location, ctx.PagesDirectory, platform, ctx.Args+ctx.PageSuffix)
}

func getPage(platforms []string, ctx *context.Context) error {
	ps := make([]platformPath, 0)
	for _, p := range platforms {
		ps = append(ps, platformPath{Platform: p, Path: makePath(ctx, p)})
	}

	for _, pp := range ps {
		if _, err := os.Stat(pp.Path); os.IsNotExist(err) {
			continue
		} else {
			ctx.FoundPlatform = pp.Platform
			ctx.Page = pp.Path
			return nil
		}
	}

	return fmt.Errorf("This page (" + ctx.Args + ") does not exist yet!\n" +
		"Submit new pages here: https://github.com/tldr-pages/tldr/issues/new?title=page%20request:%20" + ctx.Args)
}
