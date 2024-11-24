package main

import (
	"github.com/RIDOS/hr/internal/middleware/task_runner"
	"github.com/RIDOS/hr/internal/service/task_generator"
	"github.com/RIDOS/hr/internal/service/task_proccessor"
	"os"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных
// к этому моменту тасков (отдельно успешные и отдельно с ошибками).

func main() {
	// Инициализация компонентов.
	taskGenerator := task_generator.NewTaskGenerator(10)
	taskProcessor := task_proccessor.NewTaskProcessor()

	taskRunner := task_runner.NewTaskRunner(taskGenerator, taskProcessor)

	// Запуск обработки задач.
	taskRunner.Run(os.Stdout)
}
