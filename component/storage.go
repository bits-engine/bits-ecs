package component

import (
	"github.com/bits-engine/bits-ecs/entity"
)

type pool interface {
	has(ent entity.Entity) bool
	remove(ent entity.Entity) bool
}

// ComponentStorage holds everything that needed to get/set/remove/query components.
//
// Uses implementation based on sparse-sets.
// For querying uses archetype bitsets to fastly retrieving all matched entities.
type ComponentStorage struct {
	pools             map[ComponentID]pool
	idCounter         ComponentID
	entityLookupTable *lookupTable
	filtersTable      *filtersTable
}

func NewComponentStorage() *ComponentStorage {
	return &ComponentStorage{
		pools:             make(map[ComponentID]pool),
		entityLookupTable: newLookupTable(),
		filtersTable:      newFiltersTable(),
	}
}

func (s *ComponentStorage) nextComponentID() ComponentID {
	id := s.idCounter
	s.idCounter++
	return id
}

func (cs *ComponentStorage) has(ent entity.Entity, cid ComponentID) bool {
	return cs.pools[cid].has(ent)
}

func (cs *ComponentStorage) remove(ent entity.Entity, cid ComponentID) bool {
	cs.entityLookupTable.remove(ent, cid)
	return cs.pools[cid].remove(ent)
}

func (cs *ComponentStorage) removeAll(ent entity.Entity) {
	for cid, pool := range cs.pools {
		cs.entityLookupTable.remove(ent, cid)
		pool.remove(ent)
	}
}

// Remove removes component for entity.
func Remove[T any](
	cs *ComponentStorage,
	ent entity.Entity,
	typ *ComponentType[T],
) bool {
	return cs.remove(ent, typ.id)
}

// Has checks if entity has component.
func Has[T any](
	cs *ComponentStorage,
	ent entity.Entity,
	typ *ComponentType[T],
) bool {
	return cs.has(ent, typ.id)
}

// Register registers new type of component in component storage.
func Register[T any](cs *ComponentStorage) *ComponentType[T] {
	id := cs.nextComponentID()

	newPool := newPool[T]()
	cs.pools[id] = newPool

	return &ComponentType[T]{id: id}
}

// Set sets new component value for entity
func Set[T any](
	cs *ComponentStorage,
	ent entity.Entity,
	typ *ComponentType[T],
	val T,
) {
	typedPool := cs.pools[typ.id].(*typedPool[T])
	typedPool.set(ent, val)

	cs.entityLookupTable.set(ent, typ.id)
}

type componentSetter struct {
	componentID ComponentID
	setter      func(cs *ComponentStorage, ent entity.Entity)
}

// With creates new componentSetter that can be used in [SetMany]
func With[T any](typ *ComponentType[T], val T) componentSetter {
	return componentSetter{
		componentID: typ.ID(),
		setter: func(cs *ComponentStorage, ent entity.Entity) {
			typedPool := cs.pools[typ.id].(*typedPool[T])
			typedPool.set(ent, val)
		},
	}
}

// SetMany sets multiple components for entity at once. Uses [With] to create componentSetters
func SetMany(cs *ComponentStorage, ent entity.Entity, setters ...componentSetter) {
	componentIDs := make([]ComponentID, 0, len(setters))
	for _, s := range setters {
		s.setter(cs, ent)
		componentIDs = append(componentIDs, s.componentID)
	}

	cs.entityLookupTable.set(ent, componentIDs...)
}

// Get gets reference of entity's component.
//
// Can be used to read/write component.
// 
// If component is not presented on entity, second return value will be false.
func Get[T any](
	cs *ComponentStorage,
	ent entity.Entity,
	typ *ComponentType[T],
) (*T, bool) {
	typedPool := cs.pools[typ.id].(*typedPool[T])
	return typedPool.get(ent)
}

// RemoveAll removes all components from entity.
func RemoveAll(
	cs *ComponentStorage,
	ent entity.Entity,
) {
	cs.removeAll(ent)
}
