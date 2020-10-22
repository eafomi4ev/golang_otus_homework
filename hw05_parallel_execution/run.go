package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

func worker(wg *sync.WaitGroup, doneCh chan struct{}, taskCh chan Job, resultCh chan error, number int) {

	for {
		select {
		case <-doneCh:
			fmt.Printf("Worker %d ОСТАНОВЛЕН\n", number)
			wg.Done()
			return
		case job := <-taskCh:
			fmt.Printf("Worker %d начал работу над задачей %d\n", number, job.number)
			result := job.task()
			fmt.Printf("Worker %d закончил работу над задачей %d\n", number, job.number)
			resultCh <- result
		}
	}
}

type Task func() error

type Job struct {
	task   Task
	number int
}

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	wg := sync.WaitGroup{}
	taskCh := make(chan Job)
	doneCh := make(chan struct{})
	resultCh := make(chan error)
	isResultChClosed := false
	mu := sync.Mutex{}

	var errorCount int32
	var tasksDone int32

	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()

			for job := range taskCh {
				if err := job.task(); err != nil {
					atomic.AddInt32(&errorCount, 1)
				}
			}
		}()
		//go worker(&wg, doneCh, taskCh, resultCh, i+1)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer mu.Unlock()

		for {
			select {
			case ok := <-resultCh:
				atomic.AddInt32(&tasksDone, 1)

				if ok != nil {
					atomic.AddInt32(&errorCount, 1)
					if atomic.LoadInt32(&errorCount) >= int32(M) {
						mu.Lock()
						if !isResultChClosed {
							close(doneCh)
							isResultChClosed = true
						}
						mu.Unlock()
					}
				}
			}
		}
	}()

	for i, task := range tasks {
		mu.Lock()
		if !isResultChClosed {
			taskCh <- Job{task, i + 1}
		}
		mu.Unlock()
	}

	if !isResultChClosed {
		close(doneCh)
		isResultChClosed = true
	}

	wg.Wait()

	//fmt.Printf("Задач сделано: %d\n", tasksDone)

	if atomic.LoadInt32(&errorCount) >= int32(M) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
