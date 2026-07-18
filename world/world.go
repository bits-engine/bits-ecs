// world package provides ECS world, and also scheduling functionality.
//
// Main structure is [World]
package world

import (
	"github.com/bits-engine/bits-ecs/common/workerpool"
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/entity"
	"github.com/bits-engine/bits-ecs/logging"
	"github.com/bits-engine/bits-ecs/resource"
)

var worldlog = logging.New("World")

// World holds component, resource storages and schedulers.
//
// Used to setup ECS world, manipulate with it from systems, and run / stop simulation.
type World struct {
	entityCounter entity.Entity
	cs            *component.ComponentStorage
	rs            *resource.ResourceStorage
	schedulers    []*Scheduler
	isRunning     bool
	wp            *workerpool.WorkerPool
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

// Creates new entity.
//
// Needs to be used only in exclusive systems.
func (w *World) New() entity.Entity {
	idx := w.entityCounter
	w.entityCounter++
	component.SetMany(w.CS(), idx)
	return idx
}

// CS returns [component.ComponentStorage] of world.
func (w *World) CS() *component.ComponentStorage {
	return w.cs
}

// RS returns [resource.ResourceStorage] of world.
func (w *World) RS() *resource.ResourceStorage {
	return w.rs
}

// Sched returns selected scheduler.
func (w *World) Sched(schedule Schedule) *Scheduler {
	return w.schedulers[schedule]
}

// AddSystem used to add one system at time in scheduler.
//
// It is better to use Add methods on scheduler itself, for example:
//
//	w := world.World(...)
//	w.Sched(world.ScheduleUpdate).Add(...)
func (w *World) AddSystem(schedule Schedule, cfg *sysConf) SystemID {
	return w.schedulers[schedule].Add(cfg)
}

// Run runs schedulers in world.
//
// - Firstly, scheduler for world.ScheduleStartup is run once.
// - Secondly, scheduler for world.ScheduleUpdate is running in loop.
// - When World.Stop() is called, world.ScheduleShutdown is run once.
func (w *World) Run() {
	worldlog.Get().Info("Starting world...")
	// Run startup scheduler...
	w.Sched(ScheduleStartup).run()

	worldlog.Get().Info("ScheduleStartup finished")

	w.isRunning = true
	for w.isRunning {
		// Run update scheduler in loop...
		w.Sched(ScheduleUpdate).run()
	}
	worldlog.Get().Info("ScheduleUpdate loop finished")

	// Run shutdown scheduler...
	w.Sched(ScheduleShutdown).run()
	worldlog.Get().Info("ScheduleShutdown finished")

	// Stop workers
	w.wp.Stop()
	worldlog.Get().Info("Worker pool stopped")
}

func (w *World) Stop() {
	worldlog.Get().Info("Stopping...")
	w.isRunning = false
}
