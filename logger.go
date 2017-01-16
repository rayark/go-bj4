package bj4

// Logger is an interface for bj4 to log.
type Logger interface {
	// OnStart will run when bj4 is started.
	OnStart()

	// OnTaskAdded will run when a task is added into bj4.
	OnTaskAdded(task *Task)

	// OnTaskStart will run when a task is started to run.
	OnTaskStart(task *Task)

	// OnTaskStatusUpdate will run when a running task updates the status
	// using SetStatus.
	OnTaskStatusUpdate(task *Task)

	// OnTaskComplete will run when a task is completed and returns nil
	// error.
	OnTaskComplete(task *Task, result string)

	// OnTaskError will run when a task returns an error.
	OnTaskError(task *Task, err error)
}
