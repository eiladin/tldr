package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashCompletion(t *testing.T) {
	var b bytes.Buffer
	genBashCompletion(&b)
	assert.Contains(t, b.String(), "bash completion for tldr")
}

func TestZshCompletion(t *testing.T) {
	var b bytes.Buffer
	genZshCompletion(&b)
	assert.Contains(t, b.String(), "#compdef _tldr tldr")
}

func TestPwshCompletion(t *testing.T) {
	var b bytes.Buffer
	genPwshCompletion(&b)
	assert.Contains(t, b.String(), "Register-ArgumentCompleter -Native -CommandName 'tldr'")
}
