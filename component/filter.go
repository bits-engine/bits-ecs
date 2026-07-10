package component

import (
	"github.com/bits-engine/bits-ecs/common/bitset"
)

type Filter struct {
	required []ComponentID
	excluded []ComponentID
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) Require(cid ComponentID) *Filter {
	f.required = append(f.required, cid)
	return f
}

func (f *Filter) Exclude(cid ComponentID) *Filter {
	f.excluded = append(f.excluded, cid)
	return f
}

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

func RegisterFilter(cs *ComponentStorage, f *Filter) FilterID {
	return cs.filtersTable.compile(f)
}
