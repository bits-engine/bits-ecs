package component

type ComponentID uint

type ComponentType[T any] struct {
	id ComponentID
}

type componentIdentifier interface {
	ID() ComponentID
}

func (t *ComponentType[T]) ID() ComponentID {
	return t.id
}

func (id ComponentID) ID() ComponentID {
	return id
}
