package component_test

import (
	"testing"

	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/entity"
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

	ent1 := entity.Entity(0)
	ent2 := entity.Entity(1)

	component.Set(cs, ent1, CAT, CA{value: 67})
	component.Set(cs, ent1, CBT, CB{value: 4.2})
	component.SetMany(cs, ent2,
		component.With(CAT, CA{value: 78}),
		component.With(CDT, CD{value: "ent2"}),
	)

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

func TestStorage_SetGetSet(t *testing.T) {
	cs := component.NewComponentStorage()
	CAT := component.Register[CA](cs)

	ent1 := entity.Entity(0)

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

func TestStorage_Querying(t *testing.T) {
	cs := component.NewComponentStorage()
	CAT := component.Register[CA](cs)
	CBT := component.Register[CB](cs)
	CDT := component.Register[CD](cs)

	ent1 := entity.Entity(0)
	ent2 := entity.Entity(1)
	ent3 := entity.Entity(2)
	ent4 := entity.Entity(3)

	component.SetMany(cs, ent1,
		component.With(CAT, CA{value: 78}),
		component.With(CBT, CB{value: 42.2}),
	)

	component.SetMany(cs, ent2,
		component.With(CAT, CA{value: 78}),
		component.With(CBT, CB{value: 42.2}),
		component.With(CDT, CD{value: "ent2"}),
	)

	component.SetMany(cs, ent3,
		component.With(CBT, CB{value: 42.2}),
		component.With(CDT, CD{value: "ent3"}),
	)

	component.SetMany(cs, ent4)

	f1 := component.RegisterFilter(cs, component.NewFilter().
		Require(CAT).
		Require(CBT),
	)
	f2 := component.RegisterFilter(cs, component.NewFilter().
		Require(CBT).
		Exclude(CAT),
	)
	f3 := component.RegisterFilter(cs, component.NewFilter().
		Exclude(CBT).
		Exclude(CDT),
	)

	qr1 := component.Query(cs, f1)
	qr2 := component.Query(cs, f2)
	qr3 := component.Query(cs, f3)

	assert.Contains(t, qr1.List(), ent1)
	assert.Contains(t, qr1.List(), ent2)
	assert.NotContains(t, qr1.List(), ent3)
	assert.NotContains(t, qr1.List(), ent4)

	assert.Contains(t, qr2.List(), ent3)
	assert.NotContains(t, qr2.List(), ent1)
	assert.NotContains(t, qr2.List(), ent2)
	assert.NotContains(t, qr2.List(), ent4)

	assert.Contains(t, qr3.List(), ent4)
	assert.NotContains(t, qr3.List(), ent1)
	assert.NotContains(t, qr3.List(), ent2)
	assert.NotContains(t, qr3.List(), ent3)

	component.Remove(cs, ent1, CBT)

	qr3 = component.Query(cs, f3)
	assert.Contains(t, qr3.List(), ent4)
	assert.Contains(t, qr3.List(), ent1)
	assert.NotContains(t, qr3.List(), ent2)
	assert.NotContains(t, qr3.List(), ent3)
}

func TestStorage_Field(t *testing.T) {
	cs := component.NewComponentStorage()
	CAT := component.Register[CA](cs)
	CBT := component.Register[CB](cs)
	CDT := component.Register[CD](cs)

	ent1 := entity.Entity(0)
	ent2 := entity.Entity(1)
	ent3 := entity.Entity(2)
	ent4 := entity.Entity(3)

	component.SetMany(cs, ent1,
		component.With(CAT, CA{value: 78}),
		component.With(CBT, CB{value: 42.2}),
	)

	component.SetMany(cs, ent2,
		component.With(CAT, CA{value: 78}),
		component.With(CBT, CB{value: 42.2}),
		component.With(CDT, CD{value: "ent2"}),
	)

	component.SetMany(cs, ent3,
		component.With(CBT, CB{value: 42.2}),
		component.With(CDT, CD{value: "ent3"}),
	)

	component.SetMany(cs, ent4)

	f1 := component.RegisterFilter(cs, component.NewFilter().
		Require(CAT).
		Require(CBT),
	)

	qr := component.Query(cs, f1)
	assert.Equal(t, len(qr.List()), 2)

	for qr.Next() {
		ca, _ := component.Field(qr, CAT)
		ca.value = int(qr.Entity())

		cb, _ := component.Field(qr, CBT)
		cb.value *= float32(ca.value)
	}

	qr = component.Query(cs, f1)

	assert.Equal(t, len(qr.List()), 2)

	for qr.Next() {
		ca, _ := component.Field(qr, CAT)
		assert.Equal(t, ca.value, int(qr.Entity()))
	}
}

func TestStorage_RemoveAll(t *testing.T) {
	cs := component.NewComponentStorage()
	CAT := component.Register[CA](cs)
	CBT := component.Register[CB](cs)
	CDT := component.Register[CD](cs)

	ent1 := entity.Entity(0)
	ent2 := entity.Entity(1)

	component.Set(cs, ent1, CAT, CA{value: 67})
	component.Set(cs, ent1, CBT, CB{value: 4.2})
	component.SetMany(cs, ent2,
		component.With(CAT, CA{value: 78}),
		component.With(CDT, CD{value: "ent2"}),
	)

	assert.True(t, component.Has(cs, ent1, CAT))
	assert.True(t, component.Has(cs, ent1, CBT))
	assert.False(t, component.Has(cs, ent1, CDT))

	assert.True(t, component.Has(cs, ent2, CAT))
	assert.True(t, component.Has(cs, ent2, CDT))
	assert.False(t, component.Has(cs, ent2, CBT))

	component.RemoveAll(cs, ent1)

	assert.False(t, component.Has(cs, ent1, CAT))
	assert.False(t, component.Has(cs, ent1, CBT))
	assert.False(t, component.Has(cs, ent1, CDT))

	assert.True(t, component.Has(cs, ent2, CAT))
	assert.True(t, component.Has(cs, ent2, CDT))
	assert.False(t, component.Has(cs, ent2, CBT))
}
