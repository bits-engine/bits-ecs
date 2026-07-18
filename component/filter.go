package component

import (
	"github.com/bits-engine/bits-ecs/common/bitset"
)

// Filter represents filter for querying components from [ComponentStorage].
type Filter struct {
	required []ComponentID
	excluded []ComponentID
}

func NewFilter() *Filter {
	return &Filter{}
}

// Require requires that every entity in query result will have this component.
func (f *Filter) Require(cid componentIdentifier) *Filter {
	f.required = append(f.required, cid.ID())
	return f
}

// Exclude requires that none of entities in query result will have this component.
func (f *Filter) Exclude(cid componentIdentifier) *Filter {
	f.excluded = append(f.excluded, cid.ID())
	return f
}

// Required lists all required components by their [ComponentID].
func (f *Filter) Required() []ComponentID {
	return f.required
}

// Required lists all excluded components by their [ComponentID].
func (f *Filter) Excluded() []ComponentID {
	return f.excluded
}

// FilterID represents first
type FilterID = uint

type compiledFilter struct {
	required *bitset.BitSet
	excluded *bitset.BitSet
}

func compileFilter(f *Filter) *compiledFilter {
	required := &bitset.BitSet{}
	excluded := &bitset.BitSet{}

	for _, cid := range f.required {
		required.Set(uint64(cid))
	}

	for _, cid := range f.excluded {
		excluded.Set(uint64(cid))
	}

	return &compiledFilter{
		required: required,
		excluded: excluded,
	}
}

type filtersTable struct {
	filters map[FilterID]*compiledFilter
	idCounter FilterID
}

func newFiltersTable() *filtersTable {
	return &filtersTable{filters: map[FilterID]*compiledFilter{}}
}

func (t *filtersTable) nextID() FilterID {
	res := t.idCounter
	t.idCounter++
	return res
}

func (t *filtersTable) compile(f *Filter) FilterID {
	id := t.nextID()
	t.filters[id] = compileFilter(f)
	return id
}

// RegisterFilter registrates a new [Filter] inside [ComponentStorage].
//
// Returned [FilterID] is used to actually query with this filter.
func RegisterFilter(cs *ComponentStorage, f *Filter) FilterID {
	return cs.filtersTable.compile(f)
}
