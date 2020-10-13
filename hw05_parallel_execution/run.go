package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrErrorsJobFailed = errors.New("job failed")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, N int, M int) error {
	wg := sync.WaitGroup{}
	taskChanel := make(chan Task)
	doneCh := make(chan struct{})
	//errorCount := 0
	//mu := sync.Mutex{}
	//resultChanel := make(chan error)

	worker := func() {
		for {
			select {
			case <-doneCh:
				fmt.Println("Воркер завершен")
				wg.Done()
				return
			case task := <-taskChanel:
				fmt.Println("Задача взята в работу")
				task()
			}
		}
	}

	for i := 0; i < N; i++ {
		wg.Add(1)
		go worker()
	}

	for _, task := range tasks {
		taskChanel <- task
	}
	fmt.Println("Все задачи в работе")

	go func() {
		select {
		case <-time.After(time.Second * 1):
			close(doneCh)
		}
	}()

	wg.Wait()

	return nil
}
