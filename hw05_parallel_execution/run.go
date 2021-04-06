package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup

	for _, v := range tasks {
		wg.Add(1)
		// time.Sleep(1 * time.Millisecond)
		go func() {
			defer wg.Done()
			err := v()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	// go func() {
	wg.Wait()
	fmt.Println("this is the end")
	// }()
	return nil
}
