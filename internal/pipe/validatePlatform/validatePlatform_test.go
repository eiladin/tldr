package validatePlatform

import (
	"testing"
	"time"

	"github.com/eiladin/tldr/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	ctx := context.New()
	ctx.Cache.Location = "./test-cache"
	ctx.Cache.TTL = time.Hour
	ctx.Platform = "all"
	ctx.Operation = context.OperationListCommands
	err := Pipe{}.Run(ctx)
	assert.NoError(t, err)
}

func TestPlatforms(t *testing.T) {
	ctx := context.New()
	ctx.Cache.Location = "./test-cache"
	ctx.Cache.TTL = time.Hour
	ctx.AvailablePlatforms = []string{"linux", "osx", "sunos", "windows"}

	ctx.Platform = "linux"
	err := Pipe{}.Run(ctx)
	assert.NoError(t, err)

	ctx.Platform = "sunos"
	err = Pipe{}.Run(ctx)
	assert.NoError(t, err)

	ctx.Platform = "osx"
	err = Pipe{}.Run(ctx)
	assert.NoError(t, err)

	ctx.Platform = "windows"
	err = Pipe{}.Run(ctx)
	assert.NoError(t, err)

	ctx.Platform = "error"
	err = Pipe{}.Run(ctx)
	assert.Error(t, err)

}
