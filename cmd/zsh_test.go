package cmd

import (
	"testing"

	"github.com/eiladin/tldr/testdata"
	"github.com/stretchr/testify/assert"
)

func TestZshCompletion(t *testing.T) {
	out := testdata.ReadStdOut(func() {
		var cmd = newZshCmd()
		cmd.cmd.Execute()
	})

	assert.Contains(t, out, "#compdef _zsh zsh")
}
