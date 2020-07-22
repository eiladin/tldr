package middleware

import (
	"github.com/apex/log"
	"github.com/eiladin/tldr/internal/pipe"
	"github.com/eiladin/tldr/pkg/context"
)

// ErrHandler handles an action error, ignoring and logging pipe skipped
// errors.
func ErrHandler(action Action) Action {
	return func(ctx *context.Context) error {
		var err = action(ctx)
		if err == nil {
			return nil
		}
		if pipe.IsSkip(err) {
			log.WithFields(log.Fields{"message": err}).Debug("pipe skipped")
			return nil
		}
		return err
	}
}
