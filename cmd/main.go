package main

import (
	"github.com/bits-engine/bits-ecs/comp"
	"github.com/bits-engine/bits-ecs/world"
)

type Health struct {
	hp int
}
type Player struct {}
type Enemy struct {}
type HealEffect struct {
	multiplier float32
}

func main() {
	cs := comp.NewComponentStorage()
	PlayerType := comp.Register[Player](cs)
	EnemyType := comp.Register[Enemy](cs)
	HealthType := comp.Register[Health](cs)
	HealEffectType := comp.Register[HealEffect](cs)

	enemy := world.Entity(0)
	player := world.Entity(1)

	comp.Set(cs, enemy, EnemyType, Enemy{})
	comp.Set(cs, enemy, HealEffectType, HealEffect{})
}
