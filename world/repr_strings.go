package world

import (
	"fmt"
	"reflect"
)

func (sn *systemNode) repr() string {
	return fmt.Sprintf(
		"systemNode(id=%d; system=%s)",
		sn.id,
		reflect.TypeOf(sn.conf.system).String(),
	)
}

func reprExecutionLayers(s *Scheduler) [][]string {
	layers := make([][]string, 0, len(s.executionLayers))
	for _, layer := range s.executionLayers {
		newLayer := make([]string, 0, len(layer))
		for _, sysID := range layer {
			sys, _ := s.systemByID(sysID)
			newLayer = append(newLayer, sys.repr())
		}

		layers = append(layers, newLayer)
	}

	return layers
}
