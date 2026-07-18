// logging package provides helper for creating new slog loggers and some keys for structured logging.
//
// Basically uses lazy evaluation. With it, user can define their logger by slog.SetDefault()
//
// # Usage
//     var log = logging.New("my-origin")
//     
//     log.Get().Info("Hello world!")
package logging

import (
	"log/slog"

	"github.com/bits-engine/bits-ecs/common/lazyvalue"
)

// New creates new [lazyvalue.LazyValue] that gets [slog.Default] logger.
func New(originName string) *lazyvalue.LazyValue[*slog.Logger] {
	return lazyvalue.New(func() *slog.Logger {
		 return slog.Default().With(KeyOrigin, originName)
	})
}
