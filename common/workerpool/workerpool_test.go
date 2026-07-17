package workerpool_test

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/bits-engine/bits-ecs/common/workerpool"
	"github.com/stretchr/testify/assert"
)

type res struct {
	value int
}

func newwp(count int) *workerpool.WorkerPool {
	return workerpool.New(count)
}

func TestWP_3Tasks(t *testing.T) {
	wp := newwp(runtime.NumCPU())
	taskDuration := 100 * time.Millisecond
	startTime := time.Now()

	out := make(chan res, 3)
	wg := &sync.WaitGroup{}

	wp.Add(
		func() {
			time.Sleep(taskDuration)
			out <- res{value: 0}
		},
		wg,
	)
	wp.Add(
		func() {
			time.Sleep(taskDuration)
			out <- res{value: 1}
		},
		wg,
	)
	wp.Add(
		func() {
			time.Sleep(taskDuration)
			out <- res{value: 3}
		},
		wg,
	)

	wg.Wait()
	close(out)
	resList := make([]res, 0, 3)
	for msg := range out {
		resList = append(resList, msg)
	}

	totalDuration := time.Since(startTime)

	assert.Len(t, resList, 3)

	maxAllowedParallelTime := 150 * time.Millisecond
	assert.Less(t, totalDuration, maxAllowedParallelTime,
		"Tasks ran serially! Total time %v exceeded parallel threshold of %v",
		totalDuration, maxAllowedParallelTime,
	)

	wp.Stop()
}

func TestWP_3TasksOnOneWorker(t *testing.T) {
	wp := newwp(runtime.NumCPU())
	taskDuration := 100 * time.Millisecond
	startTime := time.Now()

	out := make(chan res, 3)
	wg := &sync.WaitGroup{}

	wp.AddTo(
		0,
		func() {
			time.Sleep(taskDuration)
			out <- res{value: 0}
		},
		wg,
	)
	wp.AddTo(
		0,
		func() {
			time.Sleep(taskDuration)
			out <- res{value: 1}
		},
		wg,
	)
	wp.AddTo(
		0,
		func() {
			time.Sleep(taskDuration)
			out <- res{value: 3}
		},
		wg,
	)

	wg.Wait()
	close(out)
	resList := make([]res, 0, 3)
	for msg := range out {
		resList = append(resList, msg)
	}

	totalDuration := time.Since(startTime)

	assert.Len(t, resList, 3)

	maxAllowedParallelTime := 300 * time.Millisecond
	assert.Greater(t, totalDuration, maxAllowedParallelTime)

	wp.Stop()
}
