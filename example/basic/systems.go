package main

import (
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/resource"
	"github.com/bits-engine/bits-ecs/world"
)

// Definition of systems
type tickSystem struct {
	// For access to gameState resource we need it's type with identifier.
	gameStateType *resource.ResourceType[gameState]
}

func (s *tickSystem) Access(fr world.FilterRegistry) *world.AccessConfig {
	cfg := &world.AccessConfig{}
	// Defining write access for gameState resource
	cfg.WritesRes(s.gameStateType)
	return cfg
}

func (s *tickSystem) Run(w *world.World) {
	// Getting reference for gameState resource
	state, _ := resource.Get(w.RS(), s.gameStateType)
	state.tick++
}

type stopSystem struct {
	gameStateType *resource.ResourceType[gameState]
}

func (s *stopSystem) Access(fr world.FilterRegistry) *world.AccessConfig {
	cfg := &world.AccessConfig{}
	return cfg.Exclusive(true) // Exclusive because this system manipulates world.
}

func (s *stopSystem) Run(w *world.World) {
	state, _ := resource.Get(w.RS(), s.gameStateType)
	if state.tick > state.stopAfterTick {
		log.Get().Info("stopping the simulation")
		w.Stop()
	}
}

type healSystem struct {
	healthType   *component.ComponentType[health]
	healthFilter component.FilterID
}

func (s *healSystem) Access(fr world.FilterRegistry) *world.AccessConfig {
	cfg := &world.AccessConfig{}
	cfg.Writes(s.healthType)
	// Register a new filter. It will be checked against AccessConfig for allowance for used components.
	s.healthFilter = fr.Register(cfg, component.NewFilter().Require(s.healthType))
	return cfg
}

func (s *healSystem) Run(w *world.World) {
	// Querying all entities with health component
	qr := component.Query(w.CS(), s.healthFilter)

	// Iterating over entities and using components
	for qr.Next() {
		health, _ := component.Field(qr, s.healthType)
		health.current = min(health.current+1, health.max)
	}
}

type healthPrintSystem struct {
	healthType *component.ComponentType[health]
	nameType   *component.ComponentType[name]
	filter     component.FilterID
}

func (s *healthPrintSystem) Access(fr world.FilterRegistry) *world.AccessConfig {
	cfg := &world.AccessConfig{}
	cfg.Reads(s.healthType).Reads(s.nameType)
	// Register a new filter. It will be checked against AccessConfig for allowance for used components.
	s.filter = fr.Register(cfg,
		component.NewFilter().
			Require(s.healthType).
			Require(s.nameType),
	)
	return cfg
}

func (s *healthPrintSystem) Run(w *world.World) {
	// Querying all entities with health component
	qr := component.Query(w.CS(), s.filter)

	// Iterating over entities and using components
	for qr.Next() {
		health, _ := component.Field(qr, s.healthType)
		name, _ := component.Field(qr, s.nameType)
		log.Get().Info("health print", "name", name.value, "health", health.current, "maxHealth", health.max)
	}
}
