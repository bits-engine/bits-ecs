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
