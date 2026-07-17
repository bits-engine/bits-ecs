package world

import (
	"slices"

	"github.com/bits-engine/bits-ecs/common/workerpool"
)

type systemNode struct {
	id           SystemID
	conf         *sysConf
	accessConfig *compiledAccessConfig
	wpTask       workerpool.Task
}

func (sn *systemNode) isBefore(o *systemNode) bool {
	if slices.Contains(sn.conf.before, o.id) || slices.Contains(o.conf.after, sn.id) {
		return true
	}

	return false
}

type executionLayer []SystemID

func (e executionLayer) HaveConflicting(s *Scheduler, sysNode *systemNode) bool {
	for _, sysID := range e {
		sys, _ := s.systemByID(sysID)
		if sys.accessConfig.Conflicts(sysNode.accessConfig) {
			return true
		}
	}

	return false
}

func (e executionLayer) HaveBefore(s *Scheduler, sysNode *systemNode) bool {
	for _, sysID := range e {
		sys, _ := s.systemByID(sysID)
		if sys.isBefore(sysNode) {
			return true
		}
	}

	return false
}
