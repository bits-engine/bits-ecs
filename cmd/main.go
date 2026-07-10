package main

import (
	"fmt"

	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/world"
)

type Health struct {
	hp int
}
type Player struct{}
type Enemy struct{}
type HealEffect struct {
	multiplier float32
}

func main() {
	cs := component.NewComponentStorage()

	// Registration of components
	PlayerType := component.Register[Player](cs)
	EnemyType := component.Register[Enemy](cs)
	HealthType := component.Register[Health](cs)
	HealEffectType := component.Register[HealEffect](cs)

	// Creating entities and components
	enemy := world.Entity(0)
	player := world.Entity(1)

	component.Set(cs, enemy, EnemyType, Enemy{})
	component.Set(cs, enemy, HealthType, Health{})
	component.Set(cs, player, PlayerType, Player{})
	component.Set(cs, player, HealthType, Health{})
	component.Set(cs, player, HealEffectType, HealEffect{})

	// Getting exact components
	health, exists := component.Get(cs, player, HealthType)
	if exists {
		health.hp = 100
	}
	fmt.Println(health)
	healEffect, exists := component.Get(cs, player, HealEffectType)
	if exists {
		healEffect.multiplier = 1.5
	}
	fmt.Println(healEffect)

	// Creating filter
	healthAndHealFilter := component.RegisterFilter(
		cs,
		component.NewFilter().
		Require(HealthType.ID()).
		Require(HealEffectType.ID()),
	)
	healthFilter := component.RegisterFilter(
		cs,
		component.NewFilter().Require(HealthType.ID()),
	)

	// Querying
	qr := component.Query(cs, healthAndHealFilter)

	for qr.Next() {
		health, _ := component.Field(qr, HealthType)
		healEffect, _ := component.Field(qr, HealEffectType)

		health.hp = int(healEffect.multiplier * float32(health.hp))
		fmt.Println("changed", health)
	}

	qr = component.Query(cs, healthFilter)

	for qr.Next() {
		health, _ := component.Field(qr, HealthType)
		fmt.Println("health", qr.Entity(), health)
	}
}
