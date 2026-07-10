package component

import (
	"fmt"

	"github.com/bits-engine/bits-ecs/world"
)

type queryResult struct {
	cs *ComponentStorage
	entitiesList []world.Entity
	currentIDx   int
}

func newQueryResult(cs *ComponentStorage, entitiesList []world.Entity) *queryResult {
	return &queryResult{
		entitiesList: entitiesList,
		cs: cs,
		currentIDx:   -1,
	}
}

func (qr *queryResult) Next() bool {
	qr.currentIDx++

	if qr.currentIDx >= len(qr.entitiesList) {
		return false
	}

	return true
}

func (qr *queryResult) List() []world.Entity {
	return qr.entitiesList
}

func (qr *queryResult) Entity() world.Entity {
	return qr.entitiesList[qr.currentIDx]
}

func Query(cs *ComponentStorage, filterID FilterID) *queryResult {
	filter, exists := cs.filtersTable.filters[filterID]
	if !exists {
		panic(fmt.Sprintf("filter with id %v is not registered!", filterID))
	}

	entitiesList := cs.entityLookupTable.query(filter.required, filter.excluded)
	qr := newQueryResult(cs, entitiesList)
	return qr
}

func Field[T any](qr *queryResult, typ *ComponentType[T]) (*T, bool) {
	return Get(qr.cs, qr.entitiesList[qr.currentIDx], typ)
}
