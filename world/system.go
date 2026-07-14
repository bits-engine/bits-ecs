package world

type SystemID uint

type System interface {
	Access(fr FilterRegistry) *AccessConfig
	Run(w *World)
}

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

// System will be running AFTER every system in this list
func (sc *sysConf) After(sid SystemID) *sysConf {
	sc.after = append(sc.after, sid)
	return sc
}

// System will be running BEFORE every system in this list
func (sc *sysConf) Before(sid SystemID) *sysConf {
	sc.before = append(sc.before, sid)
	return sc
}

// ThreadID is basically a user-defined integer.
// 
// If set, Scheduler will be picking the same thread-locked worker every time for this system.
func (sc *sysConf) LockThread(threadHash int) *sysConf {
	sc.threadLocked = true
	sc.lockedOnThread = threadHash
	return sc
}
