package world

import (
	"slices"
	"testing"

	"github.com/bits-engine/bits-ecs/component"
)

// ---------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------

// newTestNode builds a systemNode directly, bypassing sysConf builder
// methods and the System interface entirely — we only care about the
// scheduler's graph/layering logic here, not about actual system execution.
func newTestNode(id SystemID, before, after []SystemID, access *compiledAccessConfig) *systemNode {
	if access == nil {
		access = noConflictAccess()
	}
	return &systemNode{
		id: id,
		conf: &sysConf{
			before: before,
			after:  after,
		},
		accessConfig: access,
	}
}

// newTestScheduler registers nodes directly (bypassing Add/AddMany, which
// require a *World) and keeps systemIDX in sync, exactly like addNoRebuild does.
func newTestScheduler(nodes ...*systemNode) *Scheduler {
	s := NewScheduler()
	for _, n := range nodes {
		s.systems = append(s.systems, n)
		s.systemIDX[n.id] = len(s.systems) - 1
	}
	return s
}

func noConflictAccess() *compiledAccessConfig {
	return (&AccessConfig{}).Compile()
}

// conflictingAccessPair returns two access configs that conflict with each
// other (both write the same component), so HaveConflicting will report true.
func conflictingAccessPair(compID component.ComponentID) (*compiledAccessConfig, *compiledAccessConfig) {
	a := (&AccessConfig{compWrite: []component.ComponentID{compID}}).Compile()
	b := (&AccessConfig{compWrite: []component.ComponentID{compID}}).Compile()
	return a, b
}

func mustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	f()
}

func containsID(list []SystemID, id SystemID) bool {
	return slices.Contains(list, id)
}

func idsEqual(a, b []SystemID) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func wantOrderBefore(t *testing.T, sorted []SystemID, before, after SystemID) {
	t.Helper()
	bi, ai := -1, -1
	for i, id := range sorted {
		if id == before {
			bi = i
		}
		if id == after {
			ai = i
		}
	}
	if bi == -1 || ai == -1 {
		t.Fatalf("ids not found in sorted result: %v", sorted)
	}
	if bi >= ai {
		t.Errorf("expected system %d before system %d, got order %v", before, after, sorted)
	}
}

// ---------------------------------------------------------------------
// buildDependencyGraph
// ---------------------------------------------------------------------

func TestBuildDependencyGraph_BeforeAfterEdges(t *testing.T) {
	// system 1: before -> [2]   => edge 1->2
	// system 3: after  -> [1]   => edge 1->3
	a := newTestNode(1, []SystemID{2}, nil, nil)
	b := newTestNode(2, nil, nil, nil)
	c := newTestNode(3, nil, []SystemID{1}, nil)

	s := newTestScheduler(a, b, c)
	graph := s.buildDependencyGraph()

	if len(graph) != 3 {
		t.Fatalf("expected 3 nodes in graph, got %d", len(graph))
	}
	if !containsID(graph[1], 2) {
		t.Errorf("expected edge 1->2 (before), got %v", graph[1])
	}
	if !containsID(graph[1], 3) {
		t.Errorf("expected edge 1->3 (after), got %v", graph[1])
	}
	if len(graph[2]) != 0 {
		t.Errorf("expected no outgoing edges from 2, got %v", graph[2])
	}
	if len(graph[3]) != 0 {
		t.Errorf("expected no outgoing edges from 3, got %v", graph[3])
	}
}

func TestBuildDependencyGraph_DeduplicatesEquivalentEdges(t *testing.T) {
	// node 1 declares "before 2", node 2 declares "after 1" — same logical
	// edge (1->2) expressed from both sides, should not be duplicated.
	a := newTestNode(1, []SystemID{2}, nil, nil)
	b := newTestNode(2, nil, []SystemID{1}, nil)

	s := newTestScheduler(a, b)
	graph := s.buildDependencyGraph()

	if len(graph[1]) != 1 {
		t.Errorf("expected single deduplicated edge 1->2, got %v", graph[1])
	}
}

func TestBuildDependencyGraph_MissingBeforeDependencyPanics(t *testing.T) {
	a := newTestNode(1, []SystemID{999}, nil, nil) // 999 is not registered
	s := newTestScheduler(a)

	mustPanic(t, func() {
		s.buildDependencyGraph()
	})
}

func TestBuildDependencyGraph_MissingAfterDependencyPanics(t *testing.T) {
	a := newTestNode(1, nil, []SystemID{999}, nil) // 999 is not registered
	s := newTestScheduler(a)

	mustPanic(t, func() {
		s.buildDependencyGraph()
	})
}

// ---------------------------------------------------------------------
// toposortSystemGraph
// ---------------------------------------------------------------------

func TestToposort_SimpleChain(t *testing.T) {
	// 1 -> 2 -> 3
	a := newTestNode(1, []SystemID{2}, nil, nil)
	b := newTestNode(2, []SystemID{3}, nil, nil)
	c := newTestNode(3, nil, nil, nil)

	s := newTestScheduler(a, b, c)
	graph := s.buildDependencyGraph()
	sorted := s.toposortSystemGraph(graph)

	wantOrderBefore(t, sorted, 1, 2)
	wantOrderBefore(t, sorted, 2, 3)
}

func TestToposort_IndependentNodesKeepRegistrationOrder(t *testing.T) {
	a := newTestNode(1, nil, nil, nil)
	b := newTestNode(2, nil, nil, nil)
	c := newTestNode(3, nil, nil, nil)

	s := newTestScheduler(a, b, c)
	graph := s.buildDependencyGraph()
	sorted := s.toposortSystemGraph(graph)

	want := []SystemID{1, 2, 3}
	if !idsEqual(sorted, want) {
		t.Errorf("expected deterministic order %v, got %v", want, sorted)
	}
}

func TestToposort_DeterministicAcrossRuns(t *testing.T) {
	// Regression test: toposort must not depend on Go's randomized map
	// iteration order. Same graph, called repeatedly, must always produce
	// the exact same result.
	a := newTestNode(1, nil, nil, nil)
	b := newTestNode(2, nil, nil, nil)
	c := newTestNode(3, nil, nil, nil)
	d := newTestNode(4, nil, nil, nil)

	s := newTestScheduler(a, b, c, d)
	graph := s.buildDependencyGraph()

	first := s.toposortSystemGraph(graph)
	for i := range 20 {
		got := s.toposortSystemGraph(graph)
		if !idsEqual(first, got) {
			t.Fatalf("toposort is not deterministic: run0=%v run%d=%v", first, i+1, got)
		}
	}
}

func TestToposort_CycleDetectionPanics(t *testing.T) {
	// 1 -> 2 -> 1 (cycle)
	a := newTestNode(1, []SystemID{2}, nil, nil)
	b := newTestNode(2, []SystemID{1}, nil, nil)

	s := newTestScheduler(a, b)
	graph := s.buildDependencyGraph()

	mustPanic(t, func() {
		s.toposortSystemGraph(graph)
	})
}

// ---------------------------------------------------------------------
// buildExecutionLayers
// ---------------------------------------------------------------------

func TestBuildExecutionLayers_IndependentNonConflictingSystemsShareLayer(t *testing.T) {
	a := newTestNode(1, nil, nil, nil)
	b := newTestNode(2, nil, nil, nil)
	c := newTestNode(3, nil, nil, nil)

	s := newTestScheduler(a, b, c)
	sorted := s.toposortSystemGraph(s.buildDependencyGraph())
	layers := s.buildExecutionLayers(sorted)

	if len(layers) != 1 {
		t.Fatalf("expected 1 layer, got %d: %v", len(layers), layers)
	}
	if len(layers[0]) != 3 {
		t.Fatalf("expected all 3 systems packed into a single layer, got %v", layers[0])
	}
}

func TestBuildExecutionLayers_ConflictingSystemsAreSplit(t *testing.T) {
	accA, accB := conflictingAccessPair(component.ComponentID(1))

	a := newTestNode(1, nil, nil, accA)
	b := newTestNode(2, nil, nil, accB)

	s := newTestScheduler(a, b)
	sorted := s.toposortSystemGraph(s.buildDependencyGraph())
	layers := s.buildExecutionLayers(sorted)

	if len(layers) != 2 {
		t.Fatalf("expected 2 layers due to resource conflict, got %d: %v", len(layers), layers)
	}
	if !containsID(layers[0], 1) || !containsID(layers[1], 2) {
		t.Errorf("expected system 1 in layer 0 and system 2 in layer 1, got %v", layers)
	}
}

func TestBuildExecutionLayers_NonConflictingSystemBackfillsEarliestLayer(t *testing.T) {
	// 1 and 2 conflict on the same resource -> forced into separate layers.
	// 3 conflicts with neither -> must be packed into the EARLIEST layer
	// that has room (layer 0, alongside 1), not appended as a new layer,
	// and NOT placed into every non-conflicting layer it scans past.
	accA, accB := conflictingAccessPair(component.ComponentID(1))

	a := newTestNode(1, nil, nil, accA)
	b := newTestNode(2, nil, nil, accB)
	c := newTestNode(3, nil, nil, noConflictAccess())

	s := newTestScheduler(a, b, c)
	sorted := s.toposortSystemGraph(s.buildDependencyGraph())
	layers := s.buildExecutionLayers(sorted)

	if len(layers) != 2 {
		t.Fatalf("expected 2 layers, got %d: %v", len(layers), layers)
	}
	if !containsID(layers[0], 1) || !containsID(layers[0], 3) {
		t.Errorf("expected systems 1 and 3 packed into layer 0, got %v", layers[0])
	}
	if !containsID(layers[1], 2) {
		t.Errorf("expected system 2 in layer 1, got %v", layers[1])
	}
	// Regression guard: system 3 must appear in exactly ONE layer.
	// Without the `break` after a successful placement, it would also
	// get appended to layer 1 (since it doesn't conflict with 2 either).
	if containsID(layers[1], 3) {
		t.Errorf("system 3 was placed into layer 1 as well as layer 0 (missing break?), got %v", layers)
	}
}

func TestBuildExecutionLayers_ExplicitOrderingForcesLaterLayer(t *testing.T) {
	// No resource conflict at all between 1 and 2, but 1 declares
	// "before 2" explicitly, so they still must land in different layers.
	a := newTestNode(1, []SystemID{2}, nil, nil)
	b := newTestNode(2, nil, nil, nil)

	s := newTestScheduler(a, b)
	sorted := s.toposortSystemGraph(s.buildDependencyGraph())
	layers := s.buildExecutionLayers(sorted)

	if len(layers) != 2 {
		t.Fatalf("expected 2 layers due to explicit before/after ordering, got %d: %v", len(layers), layers)
	}
	if !containsID(layers[0], 1) {
		t.Errorf("expected system 1 in layer 0, got %v", layers[0])
	}
	if !containsID(layers[1], 2) {
		t.Errorf("expected system 2 in layer 1, got %v", layers[1])
	}
}

func TestBuildExecutionLayers_TransitiveDependencyChainProducesStrictLayers(t *testing.T) {
	// 1 -> before 2 -> before 3, no resource conflicts anywhere.
	// Each system depends (directly or transitively) on the previous one,
	// so they must end up in three strictly increasing layers.
	a := newTestNode(1, []SystemID{2}, nil, nil)
	b := newTestNode(2, []SystemID{3}, nil, nil)
	c := newTestNode(3, nil, nil, nil)

	s := newTestScheduler(a, b, c)
	sorted := s.toposortSystemGraph(s.buildDependencyGraph())
	layers := s.buildExecutionLayers(sorted)

	if len(layers) != 3 {
		t.Fatalf("expected 3 separate layers for a strict chain, got %d: %v", len(layers), layers)
	}
	if !containsID(layers[0], 1) || !containsID(layers[1], 2) || !containsID(layers[2], 3) {
		t.Errorf("expected strict layer ordering [1],[2],[3], got %v", layers)
	}
}

func TestBuildExecutionLayers_EachSystemPlacedExactlyOnce(t *testing.T) {
	// Broader regression check across a mixed graph: whatever the layering
	// looks like, every system must appear in the final layers exactly once.
	accA, accB := conflictingAccessPair(component.ComponentID(1))

	a := newTestNode(1, nil, nil, accA)
	b := newTestNode(2, nil, nil, accB)
	c := newTestNode(3, []SystemID{4}, nil, noConflictAccess())
	d := newTestNode(4, nil, nil, noConflictAccess())
	e := newTestNode(5, nil, nil, noConflictAccess())

	s := newTestScheduler(a, b, c, d, e)
	sorted := s.toposortSystemGraph(s.buildDependencyGraph())
	layers := s.buildExecutionLayers(sorted)

	seen := map[SystemID]int{}
	for _, layer := range layers {
		for _, id := range layer {
			seen[id]++
		}
	}
	for _, id := range []SystemID{1, 2, 3, 4, 5} {
		if seen[id] != 1 {
			t.Errorf("system %d placed in %d layers, expected exactly 1 (layers: %v)", id, seen[id], layers)
		}
	}
}
