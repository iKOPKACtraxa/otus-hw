package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	totalErrors := 0
	toWorkerCh := make(chan Task)
	errToReturn := error(nil)
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() { // Горутина i-го работника, который выполняет Task()
			defer wg.Done()
			for {
				task, ok := <-toWorkerCh
				if ok {
					err := task()
					if err != nil {
						mutex.Lock()
						totalErrors++
						mutex.Unlock()
					}
					if totalErrors >= m && m > 0 {
						errToReturn = ErrErrorsLimitExceeded
						<-toWorkerCh
						return
					}
				} else {
					return
				}
			}
		}()
	}
	for _, task := range tasks {
		if totalErrors < m || m <= 0 {
			toWorkerCh <- task
		}
	}
	close(toWorkerCh)
	wg.Wait()
	return errToReturn
}
