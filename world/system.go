package world

import (
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/resource"
)

type System interface {
	Run(w *World)
}

type AccessConfig struct {
	compRead []component.ComponentID
	compWrite []component.ComponentID
	resRead []resource.ResourceID
	resWrite []resource.ResourceID
	exclusive bool
}
