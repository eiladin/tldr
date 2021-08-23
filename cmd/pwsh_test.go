package cmd

import (
	"testing"

	"github.com/eiladin/tldr/testdata"
	"github.com/stretchr/testify/assert"
)

func TestPwshCompletion(t *testing.T) {
	out := testdata.ReadStdOut(func() {
		var cmd = newPwshCmd()
		err := cmd.cmd.Execute()
		assert.NoError(t, err)
	})

	assert.Contains(t, out, "Register-ArgumentCompleter -Native -CommandName 'pwsh'")
}
