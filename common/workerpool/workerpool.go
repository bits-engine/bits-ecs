package workerpool

import (
	"runtime"
	"sync"
)

type Task func()

type job struct {
	t  Task
	wg *sync.WaitGroup
}

type worker struct {
	input chan job
}

func newWorker(taskBuffer int) *worker {
	input := make(chan job, taskBuffer)

	return &worker{
		input: input,
	}
}

func (w *worker) Start(lockOSThread bool) {
	go func() {
		if lockOSThread {
			runtime.LockOSThread()
		}

		for job := range w.input {
			job.t()
			if job.wg != nil {
				job.wg.Done()
			}
		}
	}()
}

func (w *worker) Add(t Task) {
	w.input <- job{t: t, wg: nil}
}

func (w *worker) AddWG(t Task, wg *sync.WaitGroup) {
	wg.Add(1)
	w.input <- job{t: t, wg: wg}
}

func (w *worker) Stop() {
	close(w.input)
}

type WorkerPool struct {
	workers       []*worker
	roundRobinIDX int
}

func New(workersCount int) *WorkerPool {
	if workersCount <= 0 {
		panic("worker count can not be <= 0")
	}

	workers := make([]*worker, 0, workersCount)
	for range workersCount {
		w := newWorker(workersCount)
		w.Start(true)
		workers = append(workers, w)
	}

	return &WorkerPool{
		workers:       workers,
		roundRobinIDX: 0,
	}
}

func (wp *WorkerPool) nextWorkerIDX() int {
	idx := wp.roundRobinIDX
	wp.roundRobinIDX++
	wp.roundRobinIDX = wp.roundRobinIDX % len(wp.workers)

	return idx
}

func (wp *WorkerPool) WorkerCount() int {
	return len(wp.workers)
}

func (wp *WorkerPool) Stop() {
	for _, w := range wp.workers {
		w.Stop()
	}
}

func (wp *WorkerPool) AddTo(workerID int, t Task, wg *sync.WaitGroup) {
	wp.workers[workerID].AddWG(t, wg)
}

func (wp *WorkerPool) Add(t Task, wg *sync.WaitGroup) {
	wp.workers[wp.nextWorkerIDX()].AddWG(t, wg)
}
