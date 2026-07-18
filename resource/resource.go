// resource package provides structures for get/set/remove resources in ECS world.
//
// In ECS resource represents global data-structure for one world.
//
// Main structure is [ResourceStorage]
package resource

type ResourceID uint
type ResourceType[T any] struct {
	id ResourceID
}

func (id ResourceID) ID() ResourceID {
	return id
}

func (t ResourceType[T]) ID() ResourceID {
	return t.id
}
