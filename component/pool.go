package component

import "github.com/bits-engine/bits-ecs/entity"

const noEntity = -1

type typedPool[T any] struct {
	dense         []T
	denseEntities []entity.Entity
	sparse        []int // referencing idx in denseEntities and dense like Entity -> denseEntitiesIdx and Entity -> denseIdx
}

func newPool[T any]() *typedPool[T] {
	return &typedPool[T]{}
}

func (p *typedPool[T]) set(ent entity.Entity, c T) {
	for len(p.sparse) <= int(ent) {
		p.sparse = append(p.sparse, noEntity)
	}

	if p.sparse[ent] != noEntity {
		p.dense[p.sparse[ent]] = c
		return
	}

	p.sparse[ent] = len(p.denseEntities)

	p.denseEntities = append(p.denseEntities, ent)
	p.dense = append(p.dense, c)
}

func (p *typedPool[T]) has(ent entity.Entity) bool {
	if len(p.sparse) <= int(ent) {
		return false
	}

	return p.sparse[ent] != noEntity
}

func (p *typedPool[T]) get(ent entity.Entity) (*T, bool) {
	if !p.has(ent) {
		return nil, false
	}

	return &p.dense[p.sparse[ent]], true
}

func (p *typedPool[T]) remove(ent entity.Entity) bool {
	exists := p.has(ent)
	if !exists {
		return false
	}

	denseIDx := p.sparse[ent]
	lastEnt := p.denseEntities[len(p.denseEntities)-1]

	p.dense[denseIDx] = p.dense[len(p.dense)-1]
	p.dense = p.dense[:len(p.dense)-1]

	p.denseEntities[denseIDx] = p.denseEntities[len(p.denseEntities)-1]
	p.denseEntities = p.denseEntities[:len(p.denseEntities)-1]

	p.sparse[lastEnt] = denseIDx
	p.sparse[ent] = noEntity

	return true
}
