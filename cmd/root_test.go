package cmd

import (
	"testing"

	"github.com/tj/assert"
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
