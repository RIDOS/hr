package task_runner

import (
	"fmt"
	"github.com/RIDOS/hr/internal/modal"
	"github.com/RIDOS/hr/internal/service/task_generator"
	"github.com/RIDOS/hr/internal/service/task_proccessor"
	"io"
	"sync"
	"time"
)

var (
	GenerationTaskDelay = 10 * time.Second
	PrintDelay          = 3 * time.Second
)

type TaskRunner struct {
	generator *task_generator.TaskGenerator
	processor *task_proccessor.TaskProcessor
}

// NewTaskRunner создает новый TaskRunner.
func NewTaskRunner(generator *task_generator.TaskGenerator, processor *task_proccessor.TaskProcessor) *TaskRunner {
	return &TaskRunner{
		generator: generator,
		processor: processor,
	}
}

// Run запускает обработку задач в многопоточном режиме.
func (tr *TaskRunner) Run(output io.Writer) {
	tr.generator.Start()

	var wg sync.WaitGroup
	ticker := time.NewTicker(PrintDelay)
	defer ticker.Stop()

	done := make(chan struct{})

	go func() {
		for task := range tr.generator.TaskChan {
			wg.Add(1)
			go func(t modal.Task) {
				defer wg.Done()
				tr.processor.Process(t)
			}(task)
		}
		wg.Wait()
		close(done)
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				tr.PrintResults(output)
			case <-done:
				return
			}
		}
	}()

	time.Sleep(GenerationTaskDelay)
	tr.generator.Stop()
	<-done
}

// PrintResults выводит результаты обработки.
func (tr *TaskRunner) PrintResults(writer io.Writer) {
	results := tr.processor.GetResults()
	errors := tr.processor.GetErrors()

	fmt.Fprintln(writer, "Done Tasks:")
	for _, task := range results {
		fmt.Fprintf(writer, "service ID: %d, Created At: %s, Finished At: %s, Result: %s\n",
			task.ID, task.CreatedAt, task.FinishedAt, task.Result)
	}

	fmt.Fprintln(writer, "\nErrors:")
	for _, err := range errors {
		fmt.Fprintln(writer, err)
	}
}
