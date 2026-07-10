package component

type ComponentID = uint

type ComponentType[T any] struct {
	id ComponentID
}

func (t *ComponentType[T]) ID() ComponentID {
	return t.id
}
