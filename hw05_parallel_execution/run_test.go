package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 100
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunAdditional(t *testing.T) {
	tests := []struct {
		tasksCount     int
		workersCount   int
		maxErrorsCount int
		err            error
	}{
		{1, 1, 1, ErrErrorsLimitExceeded},
		{1, 1, 5, nil},
		{1, 5, 1, ErrErrorsLimitExceeded},
		{1, 5, 5, nil},
		{5, 1, 1, ErrErrorsLimitExceeded},
		{5, 1, 5, ErrErrorsLimitExceeded},
		{5, 5, 1, ErrErrorsLimitExceeded},
		{5, 5, 5, ErrErrorsLimitExceeded},
	}
	for _, tc := range tests {
		tc := tc
		defer goleak.VerifyNone(t)
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			tasks := make([]Task, 0, tc.tasksCount)

			var runTasksCount int32

			for i := 0; i < tc.tasksCount; i++ {
				err := fmt.Errorf("error from task %d", i)
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					atomic.AddInt32(&runTasksCount, 1)
					return err
				})
			}

			err := Run(tasks, tc.workersCount, tc.maxErrorsCount)

			require.Truef(t, errors.Is(err, tc.err), "actual err - %v", err)
			require.LessOrEqual(t, runTasksCount, int32(tc.workersCount+tc.maxErrorsCount), "extra tasks were started")
		})
	}
}

func TestRunAdditional2(t *testing.T) {
	tests := []struct {
		tasksCount     int
		workersCount   int
		maxErrorsCount int
		err            error
	}{
		// простые тесты
		{100, 3, 1, nil},
		{1000000, 100000, 1, nil},
		// далее тесты при m<=0: трактуется как игнорировать ошибки
		{5, 5, 0, nil},
		{5, 5, -1, nil},
		{100, 10, 0, nil},
		{10, 100, 0, nil},
		{100, 10, -1, nil},
		{10, 100, -1, nil},
	}
	for _, tc := range tests {
		tc := tc
		defer goleak.VerifyNone(t)
		t.Run("tasks without errors (for m<=0)", func(t *testing.T) {
			tasks := make([]Task, 0, tc.tasksCount)

			var runTasksCount int32
			var sumTime time.Duration

			for i := 0; i < tc.tasksCount; i++ {
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
				sumTime += taskSleep

				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}
			start := time.Now()
			err := Run(tasks, tc.workersCount, tc.maxErrorsCount)
			elapsedTime := time.Since(start)
			require.NoError(t, err)

			require.Equal(t, runTasksCount, int32(tc.tasksCount), "not all tasks were completed")
			require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
		})
	}
}
