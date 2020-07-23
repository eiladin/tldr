package platforms

import (
	"bytes"
	"testing"

	"github.com/eiladin/tldr/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestPlaformList(t *testing.T) {
	var b bytes.Buffer
	ctx := context.New()
	ctx.Writer = &b
	ctx.Operation = context.OperationListPlatforms
	ctx.AvailablePlatforms = []string{"linux", "sunos", "windows", "osx", "common"}
	err := Pipe{}.Run(ctx)
	out := b.String()
	assert.NoError(t, err)
	assert.Contains(t, out, "common")
	assert.Contains(t, out, "linux")
	assert.Contains(t, out, "osx")
	assert.Contains(t, out, "sunos")
	assert.Contains(t, out, "windows")
}
