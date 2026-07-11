package resource

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

func Register[T any](rs *ResourceStorage) *ResourceType[T] {
	return &ResourceType[T]{
		id: rs.nextID(),
	}
}

func Set[T any](rs *ResourceStorage, typ *ResourceType[T], val T) {
	rs.resources[typ.ID()] = &val
}

func Has[T any](rs *ResourceStorage, typ *ResourceType[T]) bool {
	_, exists := rs.resources[typ.ID()]
	if !exists {
		return false
	}

	return true
}

func Get[T any](rs *ResourceStorage, typ *ResourceType[T]) (*T, bool) {
	res, exists := rs.resources[typ.ID()]
	if !exists {
		return nil, false
	}

	return res.(*T), true
}

func Remove[T any](rs *ResourceStorage, typ *ResourceType[T]) bool {
	if !Has(rs, typ) {
		return false
	}

	delete(rs.resources, typ.ID())
	return true
}
