package comp

import (
	"github.com/bits-engine/bits-ecs/world"
)

type pool interface {
	has(ent world.Entity) bool
	remove(ent world.Entity) bool
}

type ComponentStorage struct {
	pools     map[ComponentID]pool
	idCounter ComponentID
	entityLookupTable *lookupTable
	filtersTable *filtersTable
}

func NewComponentStorage() *ComponentStorage {
	return &ComponentStorage{
		pools: make(map[ComponentID]pool),
		entityLookupTable: newLookupTable(),
		filtersTable: newFiltersTable(),
	}
}

func (s *ComponentStorage) nextComponentID() ComponentID {
	id := s.idCounter
	s.idCounter++
	return id
}

func (cs *ComponentStorage) has(ent world.Entity, cid ComponentID) bool {
	return cs.pools[cid].has(ent)
}

func (cs *ComponentStorage) remove(ent world.Entity, cid ComponentID) bool {
	cs.entityLookupTable.remove(ent, cid)
	return cs.pools[cid].remove(ent)
}

func Remove[T any](
	cs *ComponentStorage,
	ent world.Entity,
	typ *ComponentType[T],
) bool {
	return cs.remove(ent, typ.id)
}

func Has[T any](
	cs *ComponentStorage,
	ent world.Entity,
	typ *ComponentType[T],
) bool {
	return cs.has(ent, typ.id)
}

func Register[T any](cs *ComponentStorage) *ComponentType[T] {
	id := cs.nextComponentID()

	newPool := newPool[T]()
	cs.pools[id] = newPool

	return &ComponentType[T]{id: id}
}

func Set[T any](
	cs *ComponentStorage,
	ent world.Entity,
	typ *ComponentType[T],
	val T,
) {
	typedPool := cs.pools[typ.id].(*typedPool[T])
	typedPool.set(ent, val)

	cs.entityLookupTable.set(ent, typ.id)
}

func Get[T any](
	cs *ComponentStorage,
	ent world.Entity,
	typ *ComponentType[T],
) (*T, bool) {
	typedPool := cs.pools[typ.id].(*typedPool[T])
	return typedPool.get(ent)
}
