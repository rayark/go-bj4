// Package bj4 provides one-task-at-a-time task scheduling.
package bj4

import "time"

const (
	minWaitTime = 1 * time.Hour
)

// Config configures the scheduler
type Config struct {
	// Logger provides logging. Leaving it nil and BJ4 will not log at all.
	// If basic logging is needed, use BuiltinLogger, or if you want
	// fancier one, LogrusLogger.
	Logger Logger
	// MinWaitTime is the minimum wait time for scheduler.
	MinWaitTime time.Duration
	// TaskTTL is the timeout when a task is not scheduled anymore.  If not
	// set, all tasks will be kept.
	TaskTTL time.Duration
}

// BJ4 is the scheduler struct itself. Refer to its member functions for
// details.
type BJ4 struct {
	tasks       map[string]*Task
	taskAdded   chan *Task
	logger      Logger
	minWaitTime time.Duration
	taskTTL     time.Duration
}

// New initiates the scheduler
func New(config *Config) *BJ4 {
	if config.Logger == nil {
		config.Logger = &NilLogger{}
	}
	if config.MinWaitTime == 0 {
		config.MinWaitTime = minWaitTime
	}
	return &BJ4{
		tasks:       make(map[string]*Task),
		taskAdded:   make(chan *Task, 16),
		logger:      config.Logger,
		minWaitTime: config.MinWaitTime,
		taskTTL:     config.TaskTTL,
	}
}

// Start starts the scheduler
func (bj4 *BJ4) Start() {
	bj4.logger.OnStart()
	for {
		bj4.run()
		bj4.wait()
	}
}

func (bj4 *BJ4) run() {
	for _, task := range bj4.tasks {
		task.run()
	}
}

func (bj4 *BJ4) wait() {
	for {
		waitTime := bj4.getWaitTime()
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(waitTime)
			timeout <- true
		}()
		select {
		case task := <-bj4.taskAdded:
			bj4.enqueueTask(task)
		case <-timeout:
			return
		}
	}
}

func (bj4 *BJ4) enqueueTask(task *Task) {
	bj4.tasks[task.Name] = task
}

func (bj4 *BJ4) getWaitTime() time.Duration {
	wt := bj4.minWaitTime
	now := time.Now()
	for _, task := range bj4.tasks {
		if task.NextUpdate.IsZero() {
			continue
		}

		t := task.NextUpdate.Sub(now)
		if wt > t {
			wt = t
		}
	}

	return wt
}

// SetTask runs the task on the scheduler as soon as possible
func (bj4 *BJ4) SetTask(name string, fn TaskFunction) <-chan error {
	return bj4.SetScheduledTask(name, fn, time.Now())
}

// SetScheduledTask sets the task running on specific time
func (bj4 *BJ4) SetScheduledTask(name string, fn TaskFunction, nextUpdate time.Time) <-chan error {
	task := &Task{
		TaskStatus: TaskStatus{
			Name:       name,
			NextUpdate: nextUpdate,
			Status:     "added",
		},
		function:  fn,
		bj4:       bj4,
		errorChan: make(chan error, 1),
	}
	bj4.taskAdded <- task

	bj4.logger.OnTaskAdded(task)

	return task.errorChan
}

// GetTasks gets the tasks from the scheduler in slice format
func (bj4 *BJ4) GetTasks() []TaskStatus {
	taskStatus := make([]TaskStatus, len(bj4.tasks))
	idx := 0
	for _, task := range bj4.tasks {
		taskStatus[idx] = task.TaskStatus
		idx++
	}
	return taskStatus
}

func (bj4 *BJ4) removeTask(name string) {
	delete(bj4.tasks, name)
}
