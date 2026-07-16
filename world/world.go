package world

import (
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/resource"
)

type World struct {
	cs         *component.ComponentStorage
	rs         *resource.ResourceStorage
	schedulers []*Scheduler
	isRunning  bool
}

func New() *World {
	w := &World{
		cs:        component.NewComponentStorage(),
		rs:        resource.NewResourceStorage(),
		isRunning: false,
	}

	w.schedulers = []*Scheduler{
		NewScheduler(w),
		NewScheduler(w),
		NewScheduler(w),
	}

	return w
}

func (w *World) CS() *component.ComponentStorage {
	return w.cs
}

func (w *World) RS() *resource.ResourceStorage {
	return w.rs
}

func (w *World) Sched(schedule Schedule) *Scheduler {
	return w.schedulers[schedule]
}

func (w *World) AddSystem(schedule Schedule, cfg *sysConf) SystemID {
	return w.schedulers[schedule].Add(cfg)
}

func (w *World) Run() {
	// Run startup scheduler...
	w.Sched(ScheduleStartup).run()
	w.Sched(ScheduleStartup).stop()

	w.isRunning = true
	for w.isRunning {
		// Run update scheduler in loop...
		w.Sched(ScheduleUpdate).run()
	}

	w.Sched(ScheduleUpdate).stop()

	// Run shutdown scheduler...
	w.Sched(ScheduleShutdown).run()
	w.Sched(ScheduleShutdown).stop()
}

func (w *World) Stop() {
	w.isRunning = false
}
