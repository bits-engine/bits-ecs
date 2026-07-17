package main

import (
	"log/slog"
	"os"

	"github.com/bits-engine/bits-ecs/world"
)

func main() {
	// Setting up default logger
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// Creating new world
	w := world.New(&world.Conf{
		WorkersCount: 4,
	})

	// Adding some systems
	w.Sched(world.ScheduleStartup).Add(world.NewSysConf(&SysA{}))
	sysTerm := w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&SysTerminator{
		stopCountDown: 3,
	}))
	w.Sched(world.ScheduleUpdate).Add(world.NewSysConf(&SysA{}).Before(sysTerm))

	// Running world
	w.Run()
}
