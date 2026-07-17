package logging

import (
	"log/slog"

	"github.com/bits-engine/bits-ecs/common/lazyvalue"
)

func New(originName string) *lazyvalue.LazyValue[*slog.Logger] {
	return lazyvalue.New(func() *slog.Logger {
		 return slog.Default().With(KeyOrigin, originName)
	})
}
