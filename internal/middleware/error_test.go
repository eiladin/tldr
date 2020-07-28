package middleware

import (
	"errors"
	"testing"

	"github.com/eiladin/tldr/internal/pipe"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNoError(t *testing.T) {
	err := ErrHandler(mockAction(nil))(ctx)
	assert.NoError(t, err)
}

func TestError(t *testing.T) {
	log.SetLevel(log.PanicLevel)
	err := ErrHandler(mockAction(errors.New("pipe errored")))(ctx)
	assert.Error(t, err)
}

func TestSkipped(t *testing.T) {
	err := ErrHandler(mockAction(pipe.Skip("skipped")))(ctx)
	assert.NoError(t, err)
}
