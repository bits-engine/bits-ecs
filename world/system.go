package world

// SystemID is a unique identifier of system registered in one scheduler of world.
// 
// Schedules (world.ScheduleStartup, world.ScheduleUpdate, world.ScheduleShutdown) DOES NOT SHARE SYSTEMS.
type SystemID uint

// System represents runnable system for scheduler
type System interface {
	// Access defines component/resource read/write dependencies. Also used to define component filters for querying.
	Access(fr FilterRegistry) *AccessConfig
	// Runs on every scheduler run.
	Run(w *World)
}

// sysConf is ocnfiguring system dependencies list.
type sysConf struct {
	system System
	// System will be running AFTER every system in this list
	after []SystemID
	// System will be running BEFORE every system in this list
	before     []SystemID
	threadLocked bool
	lockedOnThread int
}

func NewSysConf(sys System) *sysConf {
	return &sysConf{
		system: sys,
		threadLocked: false,
		lockedOnThread: 0,
	}
}

// Before requires that system will be running AFTER passed system
func (sc *sysConf) After(sid SystemID) *sysConf {
	sc.after = append(sc.after, sid)
	return sc
}

// Before requires that system will be running BEFORE passed system
func (sc *sysConf) Before(sid SystemID) *sysConf {
	sc.before = append(sc.before, sid)
	return sc
}

// LockThread locks system to one thread.
// 
// ThreadID is basically a user-defined integer.
// 
// Can be used to do some functionality with CGO.
// 
// If set, Scheduler will be picking the same thread-locked worker every time for this system.
func (sc *sysConf) LockThread(threadHash int) *sysConf {
	sc.threadLocked = true
	sc.lockedOnThread = threadHash
	return sc
}
