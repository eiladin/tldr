package cmd

import (
	"testing"

	"github.com/eiladin/tldr/testdata"
	"github.com/stretchr/testify/assert"
)

func TestZshCompletion(t *testing.T) {
	out := testdata.ReadStdOut(func() {
		var cmd = newZshCmd()
		err := cmd.cmd.Execute()
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "#compdef _zsh zsh")
}
