package world

import (
	"github.com/bits-engine/bits-ecs/common/workerpool"
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/logging"
	"github.com/bits-engine/bits-ecs/resource"
)

var log = logging.New("World")

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
	log.Get().Info("Starting world...")
	// Run startup scheduler...
	w.Sched(ScheduleStartup).run()

	log.Get().Info("ScheduleStartup finished")

	w.isRunning = true
	for w.isRunning {
		// Run update scheduler in loop...
		w.Sched(ScheduleUpdate).run()
	}
	log.Get().Info("ScheduleUpdate loop finished")

	// Run shutdown scheduler...
	w.Sched(ScheduleShutdown).run()
	log.Get().Info("ScheduleShutdown finished")

	// Stop workers
	w.wp.Stop()
	log.Get().Info("Worker pool stopped")
}

func (w *World) Stop() {
	log.Get().Info("Stopping...")
	w.isRunning = false
}
