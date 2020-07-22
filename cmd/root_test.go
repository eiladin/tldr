package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// var settings = cache.Cache{
// 	Remote:   "http://tldr-pages.github.com/assets/tldr.zip",
// 	TTL:      time.Minute,
// 	Location: "./tldr-cmd-test",
// }

func getCobraCommand() *cobra.Command {
	c := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	c.Flags().BoolP("update", "u", false, "")
	c.Flags().BoolP("random", "r", false, "")
	c.Flags().BoolP("purge", "", false, "")
	c.Flags().StringP("platform", "p", "linux", "")
	c.Flags().BoolP("color", "c", false, "")
	return c
}

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		update    string
		random    string
		purge     string
		args      []string
		shouldErr bool
	}{
		{"false", "false", "false", []string{}, true},
		{"true", "false", "false", []string{}, false},
		{"false", "true", "false", []string{}, false},
		{"false", "false", "true", []string{}, false},
		{"false", "false", "false", []string{"arg1"}, false},
	}

	c := getCobraCommand()

	for _, test := range tests {
		c.Flag("update").Value.Set(test.update) // nolint: errcheck
		c.Flag("random").Value.Set(test.random) // nolint: errcheck
		c.Flag("purge").Value.Set(test.purge)   // nolint: errcheck
		err := validateArgs(c, test.args)
		if test.shouldErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
