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

type SysA struct {
	out chan string
}
type SysB struct {
	count int
	out   chan string
}
type SysC struct {
	out chan string
}

func (s *SysA) Access(fr world.FilterRegistry) *world.AccessConfig {
	cfg := &world.AccessConfig{}
	return cfg
}
func (s *SysA) Run(w *world.World) {
	s.out <- "SysA"
}

func (s *SysB) Access(fr world.FilterRegistry) *world.AccessConfig {
	cfg := &world.AccessConfig{}
	return cfg
}
func (s *SysB) Run(w *world.World) {
	if s.count < 3 {
		s.out <- "SysB"
		s.count++
		return
	}

	w.Stop()
	s.out <- "SysBStop"
}

func (s *SysC) Access(fr world.FilterRegistry) *world.AccessConfig {
	cfg := &world.AccessConfig{}
	return cfg
}
func (s *SysC) Run(w *world.World) {
	s.out <- "SysC"
}

func TestScheduler_Run(t *testing.T) {
	w := world.New()

	out := make(chan string, 6)

	w.AddSystem(world.ScheduleStartup, world.NewSysConf(&SysA{out: out}))
	w.AddSystem(world.ScheduleUpdate, world.NewSysConf(&SysB{out: out}))
	w.AddSystem(world.ScheduleShutdown, world.NewSysConf(&SysC{out: out}))

	w.Run()

	result := make([]string, 0, 6)

	for range cap(result) {
		result = append(result, <-out)
	}

	assert.Equal(t, result[0], "SysA")
	assert.Equal(t, result[1], "SysB")
	assert.Equal(t, result[2], "SysB")
	assert.Equal(t, result[3], "SysB")
	assert.Equal(t, result[4], "SysBStop")
	assert.Equal(t, result[5], "SysC")
}
