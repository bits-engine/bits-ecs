package resource

// ResourceStorage holds every resource in world.
type ResourceStorage struct {
	resources map[ResourceID]any
	idCounter ResourceID
}

func NewResourceStorage() *ResourceStorage {
	return &ResourceStorage{
		resources: map[ResourceID]any{},
		idCounter: 0,
	}
}

func (s *ResourceStorage) nextID() ResourceID {
	id := s.idCounter
	s.idCounter++
	return id
}

// Register registers new resource type.
func Register[T any](rs *ResourceStorage) *ResourceType[T] {
	return &ResourceType[T]{
		id: rs.nextID(),
	}
}

// Set sets resource value in world
func Set[T any](rs *ResourceStorage, typ *ResourceType[T], val T) {
	rs.resources[typ.ID()] = &val
}

// Has checks if resource of type typ is presented in world
func Has[T any](rs *ResourceStorage, typ *ResourceType[T]) bool {
	_, exists := rs.resources[typ.ID()]
	if !exists {
		return false
	}

	return true
}

// Get gets reference of resource in world.
//
// Can be used to read/write resource.
// 
// If component is not presented, second return value will be false.
func Get[T any](rs *ResourceStorage, typ *ResourceType[T]) (*T, bool) {
	res, exists := rs.resources[typ.ID()]
	if !exists {
		return nil, false
	}

	return res.(*T), true
}

// Remove removes resource from world.
func Remove[T any](rs *ResourceStorage, typ *ResourceType[T]) bool {
	if !Has(rs, typ) {
		return false
	}

	delete(rs.resources, typ.ID())
	return true
}
