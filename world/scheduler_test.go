package world

import (
	"testing"

	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/resource"
	"github.com/stretchr/testify/assert"
)

type SysA struct {
	comp componentIdentifier
	res  resourceIdentifier
	out  chan string
}
type SysB struct {
	comp componentIdentifier
	res  resourceIdentifier
	out  chan string
}
type SysC struct {
	comp componentIdentifier
	res  resourceIdentifier
	out  chan string
}
type SysExclusive struct {
	comp componentIdentifier
	res  resourceIdentifier
	out  chan string
}

func (s *SysA) Access(fr FilterRegistry) *AccessConfig {
	cfg := &AccessConfig{}
	cfg.Reads(s.comp)
	return cfg
}
func (s *SysA) Run(w *World) {
	s.out <- "SysA"
}

func (s *SysB) Access(fr FilterRegistry) *AccessConfig {
	cfg := &AccessConfig{}
	cfg.Writes(s.comp).WritesRes(s.res)
	return cfg
}
func (s *SysB) Run(w *World) {
	s.out <- "SysB"
}

func (s *SysC) Access(fr FilterRegistry) *AccessConfig {
	cfg := &AccessConfig{}
	cfg.Reads(s.comp).WritesRes(s.res)
	return cfg
}
func (s *SysC) Run(w *World) {
	s.out <- "SysC"
}

func (s *SysExclusive) Access(fr FilterRegistry) *AccessConfig {
	cfg := &AccessConfig{}
	cfg.Exclusive(true)
	return cfg
}
func (s *SysExclusive) Run(w *World) {
	s.out <- "SysExclusive"
}

func TestScheduler_AddSystem(t *testing.T) {
	w := New(&Conf{WorkersCount: 2})
	cs := w.CS()
	rs := w.RS()

	comp := component.Register[int](cs)
	res := resource.Register[int](rs)

	sys1 := w.AddSystem(ScheduleUpdate, NewSysConf(&SysB{comp: comp, res: res}))
	sys2 := w.AddSystem(ScheduleUpdate, NewSysConf(&SysExclusive{comp: comp, res: res}))
	sys3 := w.AddSystem(ScheduleUpdate, NewSysConf(&SysC{comp: comp, res: res}).After(sys2))
	sys4 := w.AddSystem(ScheduleUpdate, NewSysConf(&SysA{comp: comp, res: res}))

	assert.Equal(t, w.Sched(ScheduleUpdate).systems[0].id, sys1)
	assert.Equal(t, w.Sched(ScheduleUpdate).systems[1].id, sys2)
	assert.Equal(t, w.Sched(ScheduleUpdate).systems[2].id, sys3)
	assert.Equal(t, w.Sched(ScheduleUpdate).systems[3].id, sys4)

	assert.Len(t, w.Sched(ScheduleUpdate).executionLayers[0], 1) // sysB
	assert.Len(t, w.Sched(ScheduleUpdate).executionLayers[1], 1) // sysExclusive
	assert.Len(t, w.Sched(ScheduleUpdate).executionLayers[2], 2) // sysC, sysA

	assert.Contains(t, w.Sched(ScheduleUpdate).executionLayers[0], sys1)
	assert.Contains(t, w.Sched(ScheduleUpdate).executionLayers[1], sys2)
	assert.Contains(t, w.Sched(ScheduleUpdate).executionLayers[2], sys3)
	assert.Contains(t, w.Sched(ScheduleUpdate).executionLayers[2], sys4)
}

func TestScheduler_Run(t *testing.T) {
	w := New(&Conf{WorkersCount: 2})
	cs := w.CS()
	rs := w.RS()

	comp := component.Register[int](cs)
	res := resource.Register[int](rs)

	out := make(chan string, 4)

	systems := w.Sched(ScheduleUpdate).AddMany(
		NewSysConf(&SysB{comp: comp, res: res, out: out}).LockThread(1),
		NewSysConf(&SysExclusive{comp: comp, res: res, out: out}).LockThread(1),
	)
	w.AddSystem(ScheduleUpdate, NewSysConf(&SysC{comp: comp, res: res, out: out}).After(systems[1]))
	w.AddSystem(ScheduleUpdate, NewSysConf(&SysA{comp: comp, res: res, out: out}))

	w.Sched(ScheduleUpdate).run()

	result := make([]string, 0, 4)

	for range cap(result) {
		result = append(result, <-out)
	}

	assert.Equal(t, result[0], "SysB", result)
	assert.Equal(t, result[1], "SysExclusive", result)
	assert.Contains(t, result[2:], "SysA", result)
	assert.Contains(t, result[2:], "SysC", result)
}
