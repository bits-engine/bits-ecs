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
	return w.schedulers[schedule].Add(w, cfg)
}

func (w *World) Run() {
	// Run startup scheduler...
	// w.Sched(ScheduleStartup).Run(w)

	w.isRunning = true
	for w.isRunning {
		// Run update scheduler in loop...
		// w.Sched(ScheduleUpdate).Run(w)
	}

	// Run shutdown scheduler...
	// w.Sched(ScheduleShutdown).Run(w)
}

func (w *World) Stop() {
	w.isRunning = false
}
