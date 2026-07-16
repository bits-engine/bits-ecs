package workerpool

import (
	"runtime"
)

type task[R any] func() R

type worker[R any] struct {
	input  chan task[R]
	output chan R
}

func newWorker[R any](taskBuffer int, output chan R) *worker[R] {
	input := make(chan task[R], taskBuffer)

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
	}()
}

func (w *worker[R]) Add(t task[R]) {
	w.input <- t
}

func (w *worker[R]) Out() <-chan R {
	return w.output
}

func (w *worker[R]) Stop() {
	close(w.input)
}

type WorkerPool[R any] struct {
	workers       []*worker[R]
	roundRobinIDX int
	output        chan R
}

func New[R any](workersCount int) *WorkerPool[R] {
	if workersCount <= 0 {
		panic("worker count can not be <= 0")
	}

	output := make(chan R, workersCount)
	workers := make([]*worker[R], 0, workersCount)
	for range workersCount {
		w := newWorker(workersCount, output)
		w.Start(true)
		workers = append(workers, w)
	}

	return &WorkerPool[R]{
		workers:       workers,
		roundRobinIDX: 0,
		output: output,
	}
}

func (wp *WorkerPool[R]) nextWorkerIDX() int {
	idx := wp.roundRobinIDX
	wp.roundRobinIDX++
	wp.roundRobinIDX = wp.roundRobinIDX % len(wp.workers)

	return idx
}

func (wp *WorkerPool[R]) Output() <-chan R {
	return wp.output
}

func (wp *WorkerPool[R]) Stop() {
	for _, w := range wp.workers {
		w.Stop()
	}

	close(wp.output)
}

func (wp *WorkerPool[R]) AddTo(workerID int, tasks ...task[R]) {
	for _, t := range tasks {
		wp.workers[workerID].Add(t)
	}
}

func (wp *WorkerPool[R]) Add(tasks ...task[R]) {
	for _, t := range tasks {
		idx := wp.nextWorkerIDX()
		wp.workers[idx].Add(t)
	}
}

func (wp *WorkerPool[R]) WaitN(resultsCount int) []R {
	result := make([]R, 0, resultsCount)

	for len(result) < resultsCount {
		result = append(result, <-wp.output)
	}

	return result
}
