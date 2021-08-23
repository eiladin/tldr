package cmd

import (
	"testing"

	"github.com/eiladin/tldr/testdata"
	"github.com/stretchr/testify/assert"
)

func TestCommandsCmd(t *testing.T) {
	cases := []struct {
		platform string
		expected []string
	}{
		{"linux", []string{"pacman", "apt", "rpm", "zypper"}},
		{"windows", []string{"choco", "cmd"}},
		{"osx", []string{"brew", "sed"}},
		{"sunos", []string{"snoop", "dmesg"}},
		{"common", []string{"curl", "nslookup", "ping", "dig"}},
	}

	for _, c := range cases {
		out := testdata.ReadStdOut(func() {
			cmd := newCommandsCmd()
			cmd.cmd.SetArgs([]string{"-p", c.platform})
			err := cmd.cmd.Execute()
			assert.NoError(t, err)
		})
		for _, e := range c.expected {
			assert.Contains(t, out, e)
		}
	}
}

func TestCommandsCmdError(t *testing.T) {
	out := testdata.ReadStdOut(func() {
		_, err := listCommands(commandsOptions{
			platform: "fake",
		})
		assert.Error(t, err, "ERROR: platform fake not found")
	})
	assert.Contains(t, out, "ERROR: platform fake not found")
}
