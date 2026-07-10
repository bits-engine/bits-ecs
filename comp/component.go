package comp

type ComponentID = uint

type ComponentType[T any] struct {
	id ComponentID
}
