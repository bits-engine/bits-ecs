package world

import (
	"fmt"
	"slices"
	"sync"

	"github.com/bits-engine/bits-ecs/common/workerpool"
)

type Scheduler struct {
	idCounter       SystemID
	systems         []*systemNode
	systemIDX       map[SystemID]int
	executionLayers []executionLayer
	w               *World
	wp              *workerpool.WorkerPool
	wg              *sync.WaitGroup
}

func NewScheduler(w *World, wp *workerpool.WorkerPool) *Scheduler {
	return &Scheduler{
		systemIDX: map[SystemID]int{},
		w:         w,
		wp:        wp,
		wg:        &sync.WaitGroup{},
	}
}

func (s *Scheduler) nextID() SystemID {
	id := s.idCounter
	s.idCounter++
	return id
}

func (s *Scheduler) addNoRebuild(conf *sysConf) SystemID {
	id := s.nextID()

	sysNode := &systemNode{
		id:   id,
		conf: conf,
		accessConfig: conf.system.Access(
			&csFilterRegistry{cs: s.w.CS()},
		).Compile(),
		wpTask: func() {
			conf.system.Run(s.w)
		},
	}

	s.systems = append(s.systems, sysNode)
	s.systemIDX[id] = len(s.systems) - 1

	return id
}

func (s *Scheduler) AddMany(confs ...*sysConf) []SystemID {
	res := make([]SystemID, 0, len(confs))

	for _, conf := range confs {
		res = append(res, s.addNoRebuild(conf))
	}

	s.rebuildExecutionLayers()

	return res
}

func (s *Scheduler) Add(conf *sysConf) SystemID {
	sysID := s.addNoRebuild(conf)
	s.rebuildExecutionLayers()
	return sysID
}

func (s *Scheduler) rebuildExecutionLayers() {
	graph := s.buildDependencyGraph()
	sortedSystems := s.toposortSystemGraph(graph)
	s.executionLayers = s.buildExecutionLayers(sortedSystems)
}

// map[BeforeSystem][]AfterSystem
func (s *Scheduler) buildDependencyGraph() map[SystemID][]SystemID {
	graph := map[SystemID][]SystemID{}

	for _, node := range s.systems {
		graph[node.id] = []SystemID{}
	}

	for _, node := range s.systems {
		// before
		for _, depID := range node.conf.before {
			if _, exists := s.systemByID(depID); !exists {
				panic(fmt.Sprintf("System %+v missing dependency from before: %d", node, depID))
			}

			if !slices.Contains(graph[node.id], depID) {
				graph[node.id] = append(graph[node.id], depID)
			}
		}

		// after
		for _, depID := range node.conf.after {
			if _, exists := s.systemByID(depID); !exists {
				panic(fmt.Sprintf("System %+v missing dependency from after: %d", node, depID))
			}

			if !slices.Contains(graph[depID], node.id) {
				graph[depID] = append(graph[depID], node.id)
			}
		}
	}

	return graph
}

func (s *Scheduler) toposortSystemGraph(graph map[SystemID][]SystemID) []SystemID {
	inDegree := make(map[SystemID]int, len(graph))
	for sysID := range graph {
		inDegree[sysID] = 0
	}

	for _, deps := range graph {
		for _, dep := range deps {
			inDegree[dep] = inDegree[dep] + 1
		}
	}

	result := make([]SystemID, 0, len(graph))
	remaining := make([]SystemID, 0)
	for _, sys := range s.systems {
		if inDegree[sys.id] == 0 {
			remaining = append(remaining, sys.id)
		}
	}

	for len(remaining) > 0 {
		sysID := remaining[0]
		remaining = remaining[1:]
		result = append(result, sysID)

		for _, dep := range graph[sysID] {
			inDegree[dep] = inDegree[dep] - 1

			if inDegree[dep] == 0 {
				remaining = append(remaining, dep)
			}
		}
	}

	if len(result) != len(graph) {
		panic("systems graph have cycle dependencies!")
	}

	return result
}

func (s *Scheduler) buildExecutionLayers(sorted []SystemID) []executionLayer {
	layers := make([]executionLayer, 0, len(sorted)/3)

	for _, sysID := range sorted {
		sysNode, _ := s.systemByID(sysID)
		minLayer := 0

		for i, layer := range layers {
			if layer.HaveBefore(s, sysNode) {
				minLayer = i + 1
			}
		}

		placed := false
		for i := minLayer; i < len(layers); i++ {
			if layers[i].HaveConflicting(s, sysNode) {
				continue
			}

			layers[i] = append(layers[i], sysID)
			placed = true
			break
		}

		if !placed {
			layers = append(layers, []SystemID{sysID})
		}
	}

	return layers
}

func (s *Scheduler) systemByID(sysID SystemID) (*systemNode, bool) {
	if idx, exists := s.systemIDX[sysID]; exists {
		return s.systems[idx], true
	}

	return nil, false
}

func (s *Scheduler) run() {
	for _, layer := range s.executionLayers {
		for _, sysID := range layer {
			sysNode, exists := s.systemByID(sysID)
			if !exists {
				panic(fmt.Sprintf("system %d not found during Run", sysID))
			}
			if len(layer) == 1 && !sysNode.conf.threadLocked {
				sysNode.wpTask()
				continue
			}

			if !sysNode.conf.threadLocked {
				s.wp.Add(sysNode.wpTask, s.wg)
				continue
			}

			s.wp.AddTo(sysNode.conf.lockedOnThread%s.wp.WorkerCount(), sysNode.wpTask, s.wg)
		}

		s.wg.Wait()
	}
}
