package component

// ComponentID represents the unique identifier of a registered component in [ComponentStorage].
type ComponentID uint

// ComponentType holds ID of a registered component with its type T.
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
