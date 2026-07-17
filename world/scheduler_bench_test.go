package world

import (
	"testing"
)

type DummySys struct{}

func (s *DummySys) Access(fr FilterRegistry) *AccessConfig {
	cfg := &AccessConfig{}
	return cfg
}
func (s *DummySys) Run(w *World) {
	counter := 0
	for range 100 {
		counter++
	}
}

func BenchmarkScheduler_Run2x2(b *testing.B) {
	w := New()

	w.Sched(ScheduleUpdate).Add(NewSysConf(&DummySys{}))
	sys2 := w.Sched(ScheduleUpdate).Add(NewSysConf(&DummySys{}))
	w.Sched(ScheduleUpdate).Add(NewSysConf(&DummySys{}).After(sys2))
	w.Sched(ScheduleUpdate).Add(NewSysConf(&DummySys{}).After(sys2))

	for b.Loop() {
		w.Sched(ScheduleUpdate).run()
	}
}

func BenchmarkScheduler_Run4x1(b *testing.B) {
	w := New()

	sys1 := w.Sched(ScheduleUpdate).Add(NewSysConf(&DummySys{}))
	sys2 := w.Sched(ScheduleUpdate).Add(NewSysConf(&DummySys{}).After(sys1))
	sys3 := w.Sched(ScheduleUpdate).Add(NewSysConf(&DummySys{}).After(sys2))
	w.Sched(ScheduleUpdate).Add(NewSysConf(&DummySys{}).After(sys3))

	for b.Loop() {
		w.Sched(ScheduleUpdate).run()
	}
}

func BenchmarkScheduler_Run1x2(b *testing.B) {
	w := New()

	w.Sched(ScheduleUpdate).Add(NewSysConf(&DummySys{}))
	w.Sched(ScheduleUpdate).Add(NewSysConf(&DummySys{}))

	for b.Loop() {
		w.Sched(ScheduleUpdate).run()
	}
}
