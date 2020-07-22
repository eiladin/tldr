package middleware

import "github.com/eiladin/tldr/pkg/context"

var ctx = &context.Context{}

func mockAction(err error) Action {
	return func(ctx *context.Context) error {
		return err
	}
}
