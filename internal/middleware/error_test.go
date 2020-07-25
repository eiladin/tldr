package middleware

import (
	"fmt"
	"testing"

	"github.com/eiladin/tldr/internal/pipe"
	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		require.NoError(t, ErrHandler(mockAction(nil))(ctx))
	})

	t.Run("some err", func(t *testing.T) {
		require.Error(t, ErrHandler(mockAction(fmt.Errorf("pipe errored")))(ctx))
	})

	t.Run("skipped", func(t *testing.T) {
		require.NoError(t, ErrHandler(mockAction(pipe.Skip("test")))(ctx))
	})
}
