package component

import (
	"github.com/bits-engine/bits-ecs/common/bitset"
	"github.com/bits-engine/bits-ecs/entity"
)

type archetypeID = string
type archetype struct {
	components *bitset.BitSet
	cachedID   archetypeID
}

type entityArchetypeRecord struct {
	idxInArchetype int
	archetypeID    archetypeID
}

func (at *archetype) id() archetypeID {
	return at.cachedID
}

func (at *archetype) recomputedID() archetypeID {
	at.cachedID = at.components.Key()
	return at.id()
}

type lookupTable struct {
	archetypes        map[archetypeID]*archetype
	archetypeToEntity map[archetypeID][]entity.Entity
	entityToArchetype *typedPool[entityArchetypeRecord]
}

func newLookupTable() *lookupTable {
	return &lookupTable{
		archetypes:        make(map[archetypeID]*archetype),
		archetypeToEntity: make(map[archetypeID][]entity.Entity),
		entityToArchetype: &typedPool[entityArchetypeRecord]{},
	}
}

func (t *lookupTable) query(required *bitset.BitSet, excluded *bitset.BitSet) []entity.Entity {
	res := make([]entity.Entity, 0)

	for archID, arch := range t.archetypes {
		if !arch.components.HasAll(required) || !arch.components.HasNone(excluded) {
			continue
		}

		res = append(res, t.archetypeToEntity[archID]...)
	}

	return res
}

func (t *lookupTable) removeEntityArchetype(ent entity.Entity, archID archetypeID) {
	rec, _ := t.entityToArchetype.get(ent)
	idx := rec.idxInArchetype
	ln := len(t.archetypeToEntity[archID])
	movedEnt := t.archetypeToEntity[archID][ln-1]
	t.archetypeToEntity[archID][idx] = t.archetypeToEntity[archID][ln-1]
	t.archetypeToEntity[archID] = t.archetypeToEntity[archID][:ln-1]

	if idx != ln-1 {
		movedEntRec, _ := t.entityToArchetype.get(movedEnt)
		movedEntRec.idxInArchetype = idx
	}

	if ln-1 == 0 {
		delete(t.archetypeToEntity, archID)
		delete(t.archetypes, archID)
	}

	t.entityToArchetype.remove(ent)
}

func (t *lookupTable) addEntityArchetype(ent entity.Entity, newArch archetype) {
	if _, exists := t.archetypes[newArch.recomputedID()]; !exists {
		t.archetypes[newArch.id()] = &newArch
	}

	if _, exists := t.archetypeToEntity[newArch.id()]; !exists {
		t.archetypeToEntity[newArch.id()] = make([]entity.Entity, 0, 1)
	}

	t.archetypeToEntity[newArch.id()] = append(t.archetypeToEntity[newArch.id()], ent)
	rec := entityArchetypeRecord{
		idxInArchetype: len(t.archetypeToEntity[newArch.id()]) - 1,
		archetypeID:    newArch.id(),
	}
	t.entityToArchetype.set(ent, rec)
}

func (t *lookupTable) set(ent entity.Entity, cids ...ComponentID) {
	newArch := archetype{components: &bitset.BitSet{}}

	if rec, exists := t.entityToArchetype.get(ent); exists {
		newArch.components = t.archetypes[rec.archetypeID].components.Clone()
		t.removeEntityArchetype(ent, rec.archetypeID)
	}

	for _, cid := range cids {
		newArch.components.Set(uint64(cid))
	}

	t.addEntityArchetype(ent, newArch)
}

func (t *lookupTable) remove(ent entity.Entity, cids ...ComponentID) {
	newArch := archetype{components: &bitset.BitSet{}}

	if rec, exists := t.entityToArchetype.get(ent); exists {
		newArch.components = t.archetypes[rec.archetypeID].components.Clone()
		t.removeEntityArchetype(ent, rec.archetypeID)
	}

	for _, cid := range cids {
		newArch.components.Clear(uint64(cid))
	}

	t.addEntityArchetype(ent, newArch)
}
