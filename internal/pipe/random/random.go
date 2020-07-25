package random

import (
	"io/ioutil"
	"math/rand"
	"path"
	"strings"
	"time"

	"github.com/eiladin/tldr/internal/pipe"
	"github.com/eiladin/tldr/pkg/context"
)

type Pipe struct{}

func (Pipe) String() string {
	return "getting random page"
}

func (Pipe) Run(ctx *context.Context) error {
	if ctx.Random {
		commonPath := path.Join(ctx.Cache.Location, ctx.PagesDirectory, "common")
		platformPath := path.Join(ctx.Cache.Location, ctx.PagesDirectory, ctx.Platform)
		paths := []string{commonPath, platformPath}
		srcs := make([]string, 0)
		for _, p := range paths {
			files, err := ioutil.ReadDir(p)
			if err != nil {
				return err
			}
			for _, f := range files {
				if strings.HasSuffix(f.Name(), ctx.PageSuffix) {
					srcs = append(srcs, path.Join(p, f.Name()))
				}
			}
		}
		rand.Seed(time.Now().UTC().UnixNano())
		page := srcs[rand.Intn(len(srcs))]
		if strings.Contains(page, "common") {
			ctx.FoundPlatform = "common"
		} else {
			ctx.FoundPlatform = ctx.Platform
		}
		ctx.Page = page
	} else {
		return pipe.Skip("not serving a random page")
	}
	return nil
}
