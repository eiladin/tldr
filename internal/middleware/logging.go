package middleware

import (
	"github.com/eiladin/tldr/pkg/context"
	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

// Logging pretty prints the given action and its title.
func Logging(title string, next Action) Action {
	return func(ctx *context.Context) error {
		log.Debug(aurora.Sprintf(aurora.Bold(title)))
		return next(ctx)
	}
}
