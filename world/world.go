package world

import (
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/resource"
	"github.com/bits-engine/bits-ecs/scheduler"
)

type World struct {
	cs         *component.ComponentStorage
	rs         *resource.ResourceStorage
	schedulers []*scheduler.Scheduler
	isRunning  bool
}

func New() *World {
	return &World{
		cs:        component.NewComponentStorage(),
		rs:        resource.NewResourceStorage(),
		isRunning: false,
		schedulers: []*scheduler.Scheduler{
			scheduler.New(),
			scheduler.New(),
			scheduler.New(),
		},
	}
}

func (w *World) CS() *component.ComponentStorage {
	return w.cs
}

func (w *World) RS() *resource.ResourceStorage {
	return w.rs
}

func (w *World) Sched(schedule Schedule) *scheduler.Scheduler {
	return w.schedulers[schedule]
}
