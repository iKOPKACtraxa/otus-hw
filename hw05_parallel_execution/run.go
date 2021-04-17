package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var totalErrors int32
	toWorkerCh := make(chan Task)
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() { // Горутина i-го работника, который выполняет Task()
			defer wg.Done()
			for task := range toWorkerCh {
				err := task()
				if err != nil {
					atomic.AddInt32(&totalErrors, 1)
				}
			}
		}()
	}
	for _, task := range tasks {
		if int(atomic.LoadInt32(&totalErrors)) < m || m <= 0 {
			toWorkerCh <- task
		} else {
			break
		}
	}
	close(toWorkerCh)
	wg.Wait()
	if int(totalErrors) >= m && m > 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
