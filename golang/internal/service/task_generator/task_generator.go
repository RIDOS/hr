package task_generator

import (
	"github.com/RIDOS/hr/internal/modal"
	"time"
)

type TaskGenerator struct {
	TaskChan chan modal.Task
	StopChan chan struct{}
}

func NewTaskGenerator(bufferSize int) *TaskGenerator {
	return &TaskGenerator{
		TaskChan: make(chan modal.Task, bufferSize),
		StopChan: make(chan struct{}),
	}
}

func (tg *TaskGenerator) Start() {
	go func() {
		for {
			select {
			case <-tg.StopChan:
				close(tg.TaskChan)
				return
			default:
				now := time.Now()
				createdAt := now.Format(time.RFC3339)
				if now.Nanosecond()%2 > 0 {
					createdAt = "Some error occurred"
				}
				task := modal.Task{ID: int(now.UnixNano()), CreatedAt: createdAt}
				tg.TaskChan <- task
				time.Sleep(100 * time.Millisecond) // Задержка между задачами
			}
		}
	}()
}

func (tg *TaskGenerator) Stop() {
	close(tg.StopChan)
}
