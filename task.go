package bj4

import (
	"fmt"
	"time"
)

// Task defines the scheduler task.
type Task struct {
	TaskStatus
	bj4       *BJ4
	function  TaskFunction
	errorChan chan error
}

// TaskStatus defines the status of a task.
type TaskStatus struct {
	Name       string
	Status     string
	NextUpdate time.Time
	Completed  time.Time
	Disabled   bool
}

// TaskFunction defines the function of a task.
type TaskFunction func(task *Task) (result string, nextUpdate time.Time, err error)

// SetStatus sets the status of the running task.
func (task *Task) SetStatus(status string) {
	task.Status = status
	task.bj4.logger.OnTaskStatusUpdate(task)
}

func (task *Task) run() {
	if task.Disabled {
		ttl := task.bj4.taskTTL
		now := time.Now()
		if ttl > 0 && now.Sub(task.Completed) >= ttl {
			task.bj4.removeTask(task.Name)
		}
		return
	}

	if time.Since(task.NextUpdate) < 0 {
		return
	}

	// drain errorChan to prevent blocking
	for len(task.errorChan) > 0 {
		<-task.errorChan
	}

	task.Status = "running"
	task.bj4.logger.OnTaskStart(task)

	result, next, err := task.function(task)

	task.Completed = time.Now()
	if next.IsZero() {
		task.Disabled = true
		task.NextUpdate = time.Time{}
	} else {
		task.NextUpdate = next
	}

	if err != nil {
		task.Status = fmt.Sprintf("error: %s", err.Error())
		task.bj4.logger.OnTaskError(task, err)
		task.errorChan <- err
	} else {
		task.Status = fmt.Sprintf("completed: %s", result)
		task.bj4.logger.OnTaskComplete(task, result)
		task.errorChan <- nil
	}
}
