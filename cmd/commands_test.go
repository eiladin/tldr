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
			cmd.cmd.Execute()
		})
		for _, e := range c.expected {
			assert.Contains(t, out, e)
		}
	}
}

func TestCommandsCmdError(t *testing.T) {
	out := testdata.ReadStdOut(func() {
		listCommands(commandsOptions{
			platform: "fake",
		})
	})
	assert.Contains(t, out, "ERROR: platform fake not found")
}
