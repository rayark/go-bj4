package bj4

import log "github.com/sirupsen/logrus"

// LogrusLogger implements Logger and uses sirupsen/logrus to log. This logger
// provides more verbose information than BuiltinLogger.
type LogrusLogger struct{}

func (lgr *LogrusLogger) OnStart() {
	log.WithFields(log.Fields{
		"pool": "bj4",
	}).Infof("task scheduler is starting")
}

func (lgr *LogrusLogger) OnTaskAdded(task *Task) {
	log.WithFields(log.Fields{
		"task": task.TaskStatus,
		"pool": "bj4",
	}).Infof("adding task \"%s\"", task.Name)
}

func (lgr *LogrusLogger) OnTaskStart(task *Task) {
	log.WithFields(log.Fields{
		"task": task.TaskStatus,
		"pool": "bj4",
	}).Infof("task \"%s\" starts running", task.Name)
}

func (lgr *LogrusLogger) OnTaskStatusUpdate(task *Task) {
	log.WithFields(log.Fields{
		"task": task.TaskStatus,
		"pool": "bj4",
	}).Infof("task \"%s\" update: %s", task.Name, task.Status)
}

func (lgr *LogrusLogger) OnTaskComplete(task *Task, result string) {
	log.WithFields(log.Fields{
		"task": task.TaskStatus,
		"pool": "bj4",
	}).Infof("task \"%s\" complete: %s", task.Name, result)
}

func (lgr *LogrusLogger) OnTaskError(task *Task, err error) {
	log.WithFields(log.Fields{
		"task": task.TaskStatus,
		"pool": "bj4",
	}).Errorf("task \"%s\" error: %s", task.Name, err.Error())
}
