package task_proccessor

import (
	"fmt"
	"github.com/RIDOS/hr/internal/modal"
	"sync"
	"time"
)

type TaskProcessor struct {
	results      []modal.Task
	errors       []error
	resultsMutex sync.Mutex
	errorsMutex  sync.Mutex
}

// NewTaskProcessor создает новый процессор задач.
func NewTaskProcessor() *TaskProcessor {
	return &TaskProcessor{}
}

// Process обрабатывает задачу и сохраняет результат или ошибку.
func (tp *TaskProcessor) Process(task modal.Task) {
	finishedAt := time.Now().Format(time.RFC3339Nano)
	var result string
	if _, err := time.Parse(time.RFC3339, task.CreatedAt); err == nil {
		result = "service has been succeeded"
	} else {
		result = "service has failed"
	}

	task.FinishedAt = finishedAt
	task.Result = result

	tp.resultsMutex.Lock()
	defer tp.resultsMutex.Unlock()

	if result == "service has failed" {
		tp.errorsMutex.Lock()
		tp.errors = append(tp.errors, fmt.Errorf("task_generator ID %d failed: %s", task.ID, task.CreatedAt))
		tp.errorsMutex.Unlock()
	} else {
		tp.results = append(tp.results, task)
	}
}

// GetResults возвращает обработанные задачи.
func (tp *TaskProcessor) GetResults() []modal.Task {
	tp.resultsMutex.Lock()
	defer tp.resultsMutex.Unlock()
	return tp.results
}

// GetErrors возвращает ошибки обработки.
func (tp *TaskProcessor) GetErrors() []error {
	tp.errorsMutex.Lock()
	defer tp.errorsMutex.Unlock()
	return tp.errors
}
