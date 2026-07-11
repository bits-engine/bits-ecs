package world

type Schedule uint

const (
	ScheduleStartup Schedule = iota
	ScheduleUpdate Schedule = iota
	ScheduleShutdown Schedule = iota
)
