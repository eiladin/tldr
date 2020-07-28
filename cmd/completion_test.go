package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompletionCmd(t *testing.T) {
	cmd := newCompletionCmd()
	assert.Len(t, cmd.cmd.Commands(), 3)
}
