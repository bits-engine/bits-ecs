package component_test

import (
	"testing"

	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/world"
	"github.com/stretchr/testify/assert"
)

type CA struct {
	value int
}
type CB struct {
	value float32
}
type CD struct {
	value string
}

// Registering components
func TestStorage_Register(t *testing.T) {
	cs := component.NewComponentStorage()
	CAT := component.Register[CA](cs)
	CBT := component.Register[CB](cs)

	assert.Equal(t, CAT.ID(), component.ComponentID(0))
	assert.Equal(t, CBT.ID(), component.ComponentID(1))
}

func TestStorage_SetHasRemove(t *testing.T) {
	cs := component.NewComponentStorage()
	CAT := component.Register[CA](cs)
	CBT := component.Register[CB](cs)
	CDT := component.Register[CD](cs)

	ent1 := world.Entity(0)
	ent2 := world.Entity(1)

	component.Set(cs, ent1, CAT, CA{value: 67})
	component.Set(cs, ent1, CBT, CB{value: 4.2})
	component.Set(cs, ent2, CAT, CA{value: 78})
	component.Set(cs, ent2, CDT, CD{value: "ent2"})

	assert.True(t, component.Has(cs, ent1, CAT))
	assert.True(t, component.Has(cs, ent1, CBT))
	assert.False(t, component.Has(cs, ent1, CDT))

	assert.True(t, component.Has(cs, ent2, CAT))
	assert.True(t, component.Has(cs, ent2, CDT))
	assert.False(t, component.Has(cs, ent2, CBT))

	component.Remove(cs, ent1, CAT)
	component.Remove(cs, ent2, CDT)
	assert.False(t, component.Has(cs, ent1, CAT))
	assert.True(t, component.Has(cs, ent1, CBT))
	assert.False(t, component.Has(cs, ent1, CDT))

	assert.True(t, component.Has(cs, ent2, CAT))
	assert.False(t, component.Has(cs, ent2, CDT))
	assert.False(t, component.Has(cs, ent2, CBT))
}

// getting, changing, setting again
func TestStorage_SetGetSet(t *testing.T) {
	cs := component.NewComponentStorage()
	CAT := component.Register[CA](cs)

	ent1 := world.Entity(0)

	component.Set(cs, ent1, CAT, CA{value: 67})

	ca, exists := component.Get(cs, ent1, CAT)
	assert.True(t, exists)
	assert.Equal(t, ca.value, 67)

	ca.value = 100

	ca, exists = component.Get(cs, ent1, CAT)
	assert.True(t, exists)
	assert.Equal(t, ca.value, 100)

	component.Remove(cs, ent1, CAT)
	ca, exists = component.Get(cs, ent1, CAT)
	assert.False(t, exists)

	component.Set(cs, ent1, CAT, CA{value: 42})

	ca, exists = component.Get(cs, ent1, CAT)
	assert.True(t, exists)
	assert.Equal(t, ca.value, 42)
}

// Filter creation, Filter usage
func TestStorage_Querying(t *testing.T) {

}
