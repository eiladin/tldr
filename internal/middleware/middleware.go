package middleware

import "github.com/eiladin/tldr/pkg/context"

// Action is a function that takes a context and returns an error.
type Action func(ctx *context.Context) error
