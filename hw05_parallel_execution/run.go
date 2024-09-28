package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	tasksChan := make(chan Task)
	errCount := 0

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				task, ok := <-tasksChan
				if !ok {
					return
				}

				if err := task(); err != nil {
					mu.Lock()
					errCount++
					mu.Unlock()
				}
			}
		}()
	}

	for _, task := range tasks {
		mu.Lock()
		if errCount >= m {
			mu.Unlock()
			close(tasksChan)
			break
		}
		mu.Unlock()
		tasksChan <- task
	}

	mu.Lock()
	if errCount < m {
		close(tasksChan)
	}
	mu.Unlock()

	wg.Wait()

	if errCount > m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
