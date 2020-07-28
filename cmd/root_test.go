package cmd

import (
	"testing"

	"github.com/eiladin/tldr/pkg/context"
	"github.com/eiladin/tldr/testdata"
	"github.com/stretchr/testify/assert"
)

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		update    bool
		random    bool
		purge     bool
		args      []string
		shouldErr bool
	}{
		{false, false, false, []string{}, true},
		{true, false, false, []string{}, false},
		{false, true, false, []string{}, false},
		{false, false, true, []string{}, false},
		{false, false, false, []string{"arg1"}, false},
	}

	for _, test := range tests {
		opts := options{
			update: test.update,
			random: test.random,
			purge:  test.purge,
		}
		err := validateArgs(opts, test.args)
		if test.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestSetupContext(t *testing.T) {
	cases := []struct {
		platform string
		random   bool
		purge    bool
		color    bool
	}{
		{"linux", false, false, false},
		{"sunos", false, false, false},
		{"windows", false, false, false},
		{"osx", false, false, false},
		{"linux", true, false, false},
		{"linux", false, true, false},
		{"linux", false, false, true},
		{"linux", true, true, false},
		{"linux", false, true, true},
		{"linux", true, true, true},
	}

	for _, c := range cases {
		ctx := context.New()
		setupContext(ctx, options{
			platform: c.platform,
			random:   c.random,
			purge:    c.purge,
			color:    c.color,
		},
			"foo", "bar")
		assert.Equal(t, c.platform, ctx.Platform)
		assert.Equal(t, c.random, ctx.Random)
		assert.Equal(t, c.purge, ctx.PurgeCache)
		assert.Equal(t, c.color, ctx.Color)
		assert.Equal(t, "foo-bar", ctx.Args)
	}
}

func TestRootCmdVersion(t *testing.T) {
	out := testdata.ReadStdOut(func() {
		Execute("testing", []string{"--version"})
	})
	assert.Equal(t, "tldr version testing\n", out)
}

func TestFindPage(t *testing.T) {
	out := testdata.ReadStdOut(func() {
		findPage(options{
			platform: "linux",
		}, "git", "fetch")
	})

	assert.Contains(t, out, "git fetch")
}

func TestDebugMode(t *testing.T) {
	out := testdata.ReadStdOut(func() {
		cmd := newRootCmd("version")
		cmd.cmd.SetArgs([]string{"-d", "git", "fetch"})
		cmd.cmd.Execute()
	})

	assert.Contains(t, out, "debug logs enabled")
}
