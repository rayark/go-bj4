package bj4

import "log"

// BuiltinLogger implements Logger. It uses log.Printf and log.Println for
// logging.
type BuiltinLogger struct{}

func (lgr *BuiltinLogger) OnStart() {
	log.Println("bj4 is starting")
}

func (lgr *BuiltinLogger) OnTaskAdded(task *Task) {
	log.Printf("task \"%s\" added\n", task.Name)
}

func (lgr *BuiltinLogger) OnTaskStart(task *Task) {
	log.Printf("task \"%s\" starts running\n", task.Name)
}

func (lgr *BuiltinLogger) OnTaskStatusUpdate(task *Task) {
	log.Printf("task \"%s\" update: %s\n", task.Name, task.Status)
}

func (lgr *BuiltinLogger) OnTaskComplete(task *Task, result string) {
	log.Printf("task \"%s\" complete: %s\n", task.Name, result)
}

func (lgr *BuiltinLogger) OnTaskError(task *Task, err error) {
	log.Printf("task \"%s\" error: %s\n", task.Name, err.Error())
}
