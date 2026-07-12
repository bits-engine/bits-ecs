package world

type SystemID uint

type System interface {
	Access(fr FilterRegistry) *AccessConfig
	Run(w *World)
}

type sysConf struct {
	system System
	// System will be running AFTER every system in this list
	after  []SystemID
	// System will be running BEFORE every system in this list
	before []SystemID
}

func NewSysConf(sys System) *sysConf {
	return &sysConf{
		system: sys,
	}
}

// System will be running AFTER every system in this list
func (sc *sysConf) After(sid SystemID) *sysConf {
	sc.after = append(sc.after, sid)
	return sc
}

// System will be running BEFORE every system in this list
func (sc *sysConf) Before(sid SystemID) *sysConf {
	sc.after = append(sc.after, sid)
	return sc
}
