package component_test

import (
	"testing"

	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/entity"
)

type BenchA struct {
	X int
	Y float64
}

type BenchB struct {
	X float32
}

type BenchC struct {
	Name string
}

const benchmarkEntities = 100000

func prepareStorage() (*component.ComponentStorage, *component.ComponentType[BenchA], *component.ComponentType[BenchB], *component.ComponentType[BenchC]) {
	cs := component.NewComponentStorage()

	a := component.Register[BenchA](cs)
	b := component.Register[BenchB](cs)
	c := component.Register[BenchC](cs)

	for i := range benchmarkEntities {
		component.SetMany(cs, entity.Entity(i),
			component.With(a, BenchA{X: i}),
			component.With(b, BenchB{X: float32(i)}),
		)

		if i%2 == 0 {
			component.Set(cs, entity.Entity(i), c, BenchC{Name: "entity"})
		}
	}

	return cs, a, b, c
}

func BenchmarkSet(b *testing.B) {
	cs := component.NewComponentStorage()
	a := component.Register[BenchA](cs)

	for i := 0; b.Loop(); i++ {
		component.Set(cs, entity.Entity(i), a, BenchA{X: i})
	}
}

func BenchmarkGet(b *testing.B) {
	cs, a, _, _ := prepareStorage()

	for i := 0; b.Loop(); i++ {
		component.Get(cs, entity.Entity(i%benchmarkEntities), a)
	}
}

func BenchmarkHas(b *testing.B) {
	cs, a, _, _ := prepareStorage()

	for i := 0; b.Loop(); i++ {
		component.Has(cs, entity.Entity(i%benchmarkEntities), a)
	}
}

func BenchmarkRemoveSet(b *testing.B) {
	cs := component.NewComponentStorage()
	a := component.Register[BenchA](cs)

	for i := range benchmarkEntities {
		component.Set(cs, entity.Entity(i), a, BenchA{X: i})
	}

	for i := 0; b.Loop(); i++ {
		e := entity.Entity(i % benchmarkEntities)
		component.Remove(cs, e, a)
		component.Set(cs, e, a, BenchA{X: i})
	}
}

func BenchmarkSetMany(b *testing.B) {
	cs := component.NewComponentStorage()

	a := component.Register[BenchA](cs)
	c := component.Register[BenchC](cs)

	for i := 0; b.Loop(); i++ {
		component.SetMany(cs, entity.Entity(i),
			component.With(a, BenchA{X: i}),
			component.With(c, BenchC{Name: "hello"}),
		)
	}
}

func BenchmarkQuery(b *testing.B) {
	cs, a, bb, _ := prepareStorage()

	filter := component.RegisterFilter(cs,
		component.NewFilter().
			Require(a).
			Require(bb),
	)

	for b.Loop() {
		q := component.Query(cs, filter)
		_ = q.List()
	}
}

func BenchmarkQueryIterate(b *testing.B) {
	cs, a, bb, _ := prepareStorage()

	filter := component.RegisterFilter(cs,
		component.NewFilter().
			Require(a).
			Require(bb),
	)

	for b.Loop() {
		q := component.Query(cs, filter)

		for q.Next() {
			c, _ := component.Field(q, a)
			c.X++
		}
	}
}
