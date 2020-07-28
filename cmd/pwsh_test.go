package cmd

import (
	"testing"

	"github.com/eiladin/tldr/testdata"
	"github.com/stretchr/testify/assert"
)

func TestPwshCompletion(t *testing.T) {
	out := testdata.ReadStdOut(func() {
		var cmd = newPwshCmd()
		cmd.cmd.Execute()
	})

	assert.Contains(t, out, "Register-ArgumentCompleter -Native -CommandName 'pwsh'")
}
