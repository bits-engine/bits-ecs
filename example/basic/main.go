package main

import (
	"fmt"
	"runtime"

	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/logging"
	"github.com/bits-engine/bits-ecs/resource"
	"github.com/bits-engine/bits-ecs/world"
)

var log = logging.New("basic-example")

func main() {
	// Creating new world
	w := world.New(&world.Conf{
		// It is better to select worker count based on avarage number of parallel systems executed
		WorkersCount: runtime.NumCPU(),
	})

	// Registering components and resources
	healthType := component.Register[health](w.CS())
	nameType := component.Register[name](w.CS())
	gameStateType := resource.Register[gameState](w.RS())

	// Setting resources and entities with components
	resource.Set(w.RS(), gameStateType, gameState{stopAfterTick: 10})

	// Create 5 entities
	for i := range 3 {
		ent := w.New()
		component.SetMany(w.CS(), ent,
			component.With(nameType, name{value: fmt.Sprintf("Entity#%d", ent)}),
			component.With(healthType, health{max: i * 5 + 5, current: i}),
		)
	}

	// Register systems
	// world.ScheduleUpdate used for main loop of ECS world.
	// Also available world.ScheduleStartup (runs once at startup) and world.ScheduleShutdown (runs once at shutdown)
	stopSystemType := w.Sched(world.ScheduleUpdate).Add(
		world.NewSysConf(&stopSystem{gameStateType: gameStateType}),
	)
	w.Sched(world.ScheduleUpdate).AddMany(
		world.NewSysConf(&tickSystem{gameStateType: gameStateType}).After(stopSystemType),
		world.NewSysConf(&healSystem{healthType: healthType}).Before(stopSystemType),
		world.NewSysConf(&healthPrintSystem{healthType: healthType, nameType: nameType}).Before(stopSystemType),
	)

	// Start the world
	w.Run()
}
