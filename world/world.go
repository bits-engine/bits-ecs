package world

import (
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/resource"
)

type World struct {
	cs *component.ComponentStorage
	rs *resource.ResourceStorage
}
