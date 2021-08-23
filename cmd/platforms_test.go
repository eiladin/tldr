package cmd

import (
	"strings"
	"testing"

	"github.com/eiladin/tldr/testdata"
	"github.com/stretchr/testify/assert"
)

func TestPlatformsCmd(t *testing.T) {
	cases := []string{"common", "linux", "osx", "windows", "sunos", "android"}

	out := testdata.ReadStdOut(func() {
		cmd := newPlatformsCmd()
		err := cmd.cmd.Execute()
		assert.NoError(t, err)
	})
	for _, c := range cases {
		assert.Contains(t, out, c)
	}
	p := strings.Split(strings.TrimSpace(out), "\n")
	assert.Len(t, p, len(cases))
}
