package component

import (
	"fmt"

	"github.com/bits-engine/bits-ecs/entity"
)

// queryResult represents a filtered list of entities that matched filter.
type queryResult struct {
	cs *ComponentStorage
	entitiesList []entity.Entity
	currentIDx   int
}

func newQueryResult(cs *ComponentStorage, entitiesList []entity.Entity) *queryResult {
	return &queryResult{
		entitiesList: entitiesList,
		cs: cs,
		currentIDx:   -1,
	}
}

// Next selects next entity in list. If there is next entity it returns true, else returns false.
//
// Needs to be called at least once before using [Field]
func (qr *queryResult) Next() bool {
	qr.currentIDx++

	if qr.currentIDx >= len(qr.entitiesList) {
		return false
	}

	return true
}

// List returns entity list from query result.
func (qr *queryResult) List() []entity.Entity {
	return qr.entitiesList
}

// Entity returns current entity from query result.
//
// Can be skipped with [queryResult.Next]
func (qr *queryResult) Entity() entity.Entity {
	return qr.entitiesList[qr.currentIDx]
}

// Query searches all matching entities within passed filter by filterID.
func Query(cs *ComponentStorage, filterID FilterID) *queryResult {
	filter, exists := cs.filtersTable.filters[filterID]
	if !exists {
		panic(fmt.Sprintf("filter with id %v is not registered!", filterID))
	}

	entitiesList := cs.entityLookupTable.query(filter.required, filter.excluded)
	qr := newQueryResult(cs, entitiesList)
	return qr
}

// Field retrieves component from current entity in [queryResult]
func Field[T any](qr *queryResult, typ *ComponentType[T]) (*T, bool) {
	return Get(qr.cs, qr.entitiesList[qr.currentIDx], typ)
}
