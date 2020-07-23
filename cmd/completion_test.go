package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashCompletion(t *testing.T) {
	var b bytes.Buffer
	var cmd = newRootCmd("dev")
	genBashCompletion(cmd.cmd, &b)
	assert.Contains(t, b.String(), "bash completion for tldr")
}

func TestZshCompletion(t *testing.T) {
	var b bytes.Buffer
	var cmd = newRootCmd("dev")
	genZshCompletion(cmd.cmd, &b)
	assert.Contains(t, b.String(), "#compdef _tldr tldr")
}

func TestPwshCompletion(t *testing.T) {
	var b bytes.Buffer
	var cmd = newRootCmd("dev")
	genPwshCompletion(cmd.cmd, &b)
	assert.Contains(t, b.String(), "Register-ArgumentCompleter -Native -CommandName 'tldr'")
}
