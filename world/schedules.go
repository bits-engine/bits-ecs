package world

// Represents one schedule in world
type Schedule uint

const (
	// Runs once at startup.
	ScheduleStartup Schedule = iota
	// Runs in loop after startup.
	ScheduleUpdate
	// Runs once after stopping world.
	ScheduleShutdown
)
