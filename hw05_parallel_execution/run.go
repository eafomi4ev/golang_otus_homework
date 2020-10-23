package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n int, m int) error {
	if n <= 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	taskCh := make(chan Task)
	var errCount int32
	isErrorLimitsExist := m > 0

	wg.Add(n)
	for i := 0; i < n; i++ {
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
		if isErrorLimitsExist && atomic.LoadInt32(&errCount) >= int32(m) {
			break
		}
		taskCh <- task
	}
	close(taskCh)

	wg.Wait()

	if isErrorLimitsExist && atomic.LoadInt32(&errCount) >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
