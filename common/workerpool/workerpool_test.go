package workerpool_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/bits-engine/bits-ecs/common/workerpool"
	"github.com/stretchr/testify/assert"
)

type res struct {
	value int
}

func newwp(count int) *workerpool.WorkerPool[res] {
	return workerpool.New[res](count)
}

func TestWP_3Tasks(t *testing.T) {
	wp := newwp(runtime.NumCPU())
	taskDuration := 100 * time.Millisecond
	startTime := time.Now()

	wp.Add(
		func() res {
			time.Sleep(taskDuration)
			return res{value: 0}
		},
		func() res {
			time.Sleep(taskDuration)
			return res{value: 1}
		},
		func() res {
			time.Sleep(taskDuration)
			return res{value: 3}
		},
	)

	resList := wp.WaitN(3)
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

	wp.AddTo(
		0,
		func() res {
			time.Sleep(taskDuration)
			return res{value: 0}
		},
		func() res {
			time.Sleep(taskDuration)
			return res{value: 1}
		},
		func() res {
			time.Sleep(taskDuration)
			return res{value: 3}
		},
	)

	resList := wp.WaitN(3)
	totalDuration := time.Since(startTime)

	assert.Len(t, resList, 3)

	minAllowedParallelTime := 300 * time.Millisecond
	assert.Greater(t, totalDuration, minAllowedParallelTime)

	wp.Stop()
}
