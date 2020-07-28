package pipeline

import (
	"errors"
	"testing"

	"github.com/eiladin/tldr/pkg/context"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type TestErrorPipe struct{}

func (TestErrorPipe) String() string {
	return "test"
}
func (TestErrorPipe) Run(ctx *context.Context) error {
	ctx.Args = "I ran"
	return errors.New("error")
}

type TestNoErrorPipe struct{}

func (TestNoErrorPipe) String() string {
	return "test"
}
func (TestNoErrorPipe) Run(ctx *context.Context) error {
	ctx.Args = "I ran"
	return nil
}

func TestPipeline(t *testing.T) {
	p := []Piper{
		TestNoErrorPipe{},
	}
	c := context.New()
	log.SetLevel(log.PanicLevel)
	ctx, err := Execute(c, p)
	assert.NoError(t, err)
	assert.Equal(t, "I ran", ctx.Args)
}

func TestErrorPipeline(t *testing.T) {
	p := []Piper{
		TestErrorPipe{},
	}
	c := context.New()
	log.SetLevel(log.PanicLevel)
	ctx, err := Execute(c, p)
	assert.Error(t, err)
	assert.Equal(t, "I ran", ctx.Args)
}
