package middleware

import "github.com/eiladin/tldr/pkg/context"

type Action func(ctx *context.Context) error
