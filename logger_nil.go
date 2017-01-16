package bj4

// NilLogger implements Logger and does nothing.  This is the default logger if
// the logger in BJ4 config is left nil.
type NilLogger struct{}

func (lgr *NilLogger) OnStart() {
}

func (lgr *NilLogger) OnTaskAdded(task *Task) {
}

func (lgr *NilLogger) OnTaskStart(task *Task) {
}

func (lgr *NilLogger) OnTaskStatusUpdate(task *Task) {
}

func (lgr *NilLogger) OnTaskComplete(task *Task, result string) {
}

func (lgr *NilLogger) OnTaskError(task *Task, err error) {
}
