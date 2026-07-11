package world_test

import (
	"testing"

	"github.com/bits-engine/bits-ecs/world"
	"github.com/stretchr/testify/assert"
)

func TestWorld_Schedules(t *testing.T) {
	w := world.New()

	assert.NotNil(t, w.Sched(world.ScheduleStartup))
	assert.NotNil(t, w.Sched(world.ScheduleUpdate))
	assert.NotNil(t, w.Sched(world.ScheduleShutdown))
}
