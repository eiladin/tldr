package platforms

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/eiladin/tldr/pkg/context"
)

var (
	errReadingPagesDir = errors.New("unable to read pages folder")
)

type Pipe struct{}

func (Pipe) String() string {
	return "checking platform"
}

func (Pipe) Run(ctx *context.Context) error {
	var platforms []string
	available, err := ioutil.ReadDir(path.Join(ctx.Cache.Location, ctx.PagesDirectory))
	if err != nil {
		return fmt.Errorf("cache: %s: %s", err, errReadingPagesDir.Error())
	}

	for _, f := range available {
		platform := f.Name()
		if f.IsDir() {
			platforms = append(platforms, platform)
		}
	}
	ctx.AvailablePlatforms = platforms
	return nil
}
