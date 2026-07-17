package main

import (
	"log/slog"

	"github.com/bits-engine/bits-ecs/common/lazyvalue"
	"github.com/bits-engine/bits-ecs/logging"
	"github.com/bits-engine/bits-ecs/world"
)

// Setting up logger holder with origin name
var sysALog = logging.New("SysA")

// If you want to customize logger (for example: add different default fields),
// You can use lazyvalue.New() with custom function for logger initialization
var sysTermLog = lazyvalue.New(func() *slog.Logger {
	return slog.Default().With(logging.KeyOrigin, "SysTerminator", "another-field", "msg")
})

type SysA struct{}

func (s *SysA) Access(fr world.FilterRegistry) *world.AccessConfig {
	// Getting logger and using it
	sysALog.Get().Info("called Access method")
	cfg := &world.AccessConfig{}
	return cfg
}

func (s *SysA) Run(w *world.World) {
	sysALog.Get().Info("called Run method")
}

type SysTerminator struct {
	stopCountDown int
}

func (s *SysTerminator) Access(fr world.FilterRegistry) *world.AccessConfig {
	// Getting logger and using it
	sysTermLog.Get().Info("called Access method")
	cfg := &world.AccessConfig{}
	return cfg.Exclusive(true)
}

func (s *SysTerminator) Run(w *world.World) {
	if s.stopCountDown > 0 {
		sysTermLog.Get().Info("called Run method")
		s.stopCountDown--
		return
	}

	sysTermLog.Get().Info("stopping world")
	w.Stop()
}
