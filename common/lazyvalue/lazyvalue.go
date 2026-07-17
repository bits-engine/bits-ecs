package lazyvalue

import "sync"

type LazyValue[T any] struct {
	once sync.Once
	fn func() T
	val T
}

func New[T any](fn func() T) *LazyValue[T] {
	return &LazyValue[T]{
		fn: fn,
	}
}

func (lv *LazyValue[T]) Get() T {
	lv.once.Do(func() {
		lv.val = lv.fn()
	})

	return lv.val
}
