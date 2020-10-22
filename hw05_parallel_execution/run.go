package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	wg := sync.WaitGroup{}
	taskCh := make(chan Task)
	var errCount int32
	isErrorLimitsExist := M > 0

	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()

			for task := range taskCh {
				if err := task(); err != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if isErrorLimitsExist && atomic.LoadInt32(&errCount) >= int32(M) {
			break
		}
		taskCh <- task
	}
	close(taskCh)

	wg.Wait()

	if isErrorLimitsExist && atomic.LoadInt32(&errCount) >= int32(M) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
