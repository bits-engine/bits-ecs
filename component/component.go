// component package holds everything needed to manipulate with components.
//
// [ComponentStorage] is the main structure. It holds every component and can be used to set/get/remove/query components.
// 
// Uses sparse-set implementation
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
