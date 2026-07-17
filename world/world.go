package world

import (
	"github.com/bits-engine/bits-ecs/common/workerpool"
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/resource"
)

type World struct {
	cs         *component.ComponentStorage
	rs         *resource.ResourceStorage
	schedulers []*Scheduler
	isRunning  bool
	wp         *workerpool.WorkerPool
}

func New(conf *Conf) *World {
	wp := workerpool.New(conf.WorkersCount)

	w := &World{
		cs:        component.NewComponentStorage(),
		rs:        resource.NewResourceStorage(),
		isRunning: false,
		wp:        wp,
	}

	w.schedulers = []*Scheduler{
		NewScheduler(w, wp),
		NewScheduler(w, wp),
		NewScheduler(w, wp),
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

	w.isRunning = true
	for w.isRunning {
		// Run update scheduler in loop...
		w.Sched(ScheduleUpdate).run()
	}

	// Run shutdown scheduler...
	w.Sched(ScheduleShutdown).run()

	// Stop workers
	w.wp.Stop()
}

func (w *World) Stop() {
	w.isRunning = false
}
