package cmd

import (
	"testing"

	"github.com/eiladin/tldr/testdata"
	"github.com/stretchr/testify/assert"
)

func TestBashCompletion(t *testing.T) {
	out := testdata.ReadStdOut(func() {
		var cmd = newBashCmd()
		cmd.cmd.Execute()
	})

	assert.Contains(t, out, "bash completion")
}
