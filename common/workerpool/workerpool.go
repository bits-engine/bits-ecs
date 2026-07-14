package workerpool

import "runtime"

// Many workers
// Thread binded
// Run(func()...) []R
// RunLocked(func(), threadHash)

type task[R any] func() R

type worker[R any] struct {
	input  chan task[R]
	output chan R
}

func newWorker[R any](taskBuffer int) *worker[R] {
	input := make(chan task[R], taskBuffer)
	output := make(chan R, taskBuffer)

	return &worker[R]{
		input:  input,
		output: output,
	}
}

func (w *worker[R]) Start(lockOSThread bool) {
	go func() {
		if lockOSThread {
			runtime.LockOSThread()
		}

		for task := range w.input {
			w.output <- task()
		}

		close(w.output)
	}()
}

func (w *worker[R]) AddTask(t task[R]) {
	w.input <- t
}

func (w *worker[R]) Out() <-chan R {
	return w.output
}

func (w *worker[R]) Stop() {
	close(w.input)
}

type WorkerPool[R any] struct {
	workers []*worker[R]
	queue []
}

func New[R any](workersCount int) *WorkerPool[R] {
	if workersCount <= 0 {
		panic("worker count can not be <= 0")
	}

	workers := make([]*worker[R], 0, workersCount)
	for _ = range workersCount {
		w := newWorker[R](workersCount)
		w.Start(true)
		workers = append(workers, w)
	}

	return &WorkerPool[R]{
		workers: workers,
	}
}

func Run(task)
