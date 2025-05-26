/* Copyright (c) 2017, Rayark Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * * Neither the name of the copyright holder nor the names of its
 *   contributors may be used to endorse or promote products derived from
 *   this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package bj4

import (
	"fmt"
	"sync"
	"time"
)

// Task defines the scheduler task.
type Task struct {
	taskStatus
	lock      sync.RWMutex
	bj4       *BJ4
	function  TaskFunction
	errorChan chan error
}

// taskStatus holds internal task state with restricted access.
// Fields are private; access through Task's getter/setter methods only.
type taskStatus struct {
	name       string
	status     string
	nextUpdate time.Time
	completed  time.Time
	disabled   bool
}

// TaskStatus is a read-only snapshot of a task's state.
type TaskStatus struct {
	Name       string
	Status     string
	NextUpdate time.Time
	Completed  time.Time
	Disabled   bool
}

// TaskFunction defines the function of a task.
type TaskFunction func(task *Task) (result string, nextUpdate time.Time, err error)

// Name gets the name of the running task.
func (task *Task) Name() string {
	task.lock.RLock()
	defer task.lock.RUnlock()
	return task.name
}

// SetName sets the name of the running task.
func (task *Task) SetName(name string) {
	task.lock.Lock()
	defer task.lock.Unlock()
	task.name = name
}

// Status gets the status of the running task.
func (task *Task) Status() string {
	task.lock.RLock()
	defer task.lock.RUnlock()
	return task.status
}

// SetStatus sets the status of the running task.
func (task *Task) SetStatus(status string) {
	task.lock.Lock()
	defer task.lock.Unlock()
	task.status = status
}

// NextUpdate gets the next update time of the running task.
func (task *Task) NextUpdate() time.Time {
	task.lock.RLock()
	defer task.lock.RUnlock()
	return task.nextUpdate
}

// SetNextUpdate sets the next update time of the running task.
func (task *Task) SetNextUpdate(update time.Time) {
	task.lock.Lock()
	defer task.lock.Unlock()
	task.nextUpdate = update
}

// Completed gets the completed time of the running task.
func (task *Task) Completed() time.Time {
	task.lock.RLock()
	defer task.lock.RUnlock()
	return task.completed
}

// SetCompleted sets the completed time of the running task.
func (task *Task) SetCompleted(completed time.Time) {
	task.lock.Lock()
	defer task.lock.Unlock()
	task.completed = completed
}

// Disabled gets the disable status of the running task.
func (task *Task) Disabled() bool {
	task.lock.RLock()
	defer task.lock.RUnlock()
	return task.disabled
}

// SetDisabled sets the disable status of the running task.
func (task *Task) SetDisabled(disabled bool) {
	task.lock.Lock()
	defer task.lock.Unlock()
	task.disabled = disabled
}

// TaskStatus gets the task status of the running task.
func (task *Task) TaskStatus() TaskStatus {
	task.lock.RLock()
	defer task.lock.RUnlock()
	return TaskStatus{
		Name:       task.name,
		Status:     task.status,
		NextUpdate: task.nextUpdate,
		Completed:  task.completed,
		Disabled:   task.disabled,
	}
}

func (task *Task) run() {
	if task.Disabled() {
		ttl := task.bj4.taskTTL
		now := time.Now()
		if ttl > 0 && now.Sub(task.Completed()) >= ttl {
			task.bj4.removeTask(task.Name())
		}
		return
	}

	if time.Since(task.NextUpdate()) < 0 {
		return
	}

	// drain errorChan to prevent blocking
	for len(task.errorChan) > 0 {
		<-task.errorChan
	}

	task.SetStatus("running")
	task.bj4.logger.OnTaskStart(task)

	result, next, err := task.function(task)

	task.SetCompleted(time.Now())
	if next.IsZero() {
		task.SetDisabled(true)
		task.SetNextUpdate(time.Time{})
	} else {
		task.SetNextUpdate(next)
	}

	if err != nil {
		task.SetStatus(fmt.Sprintf("error: %s", err.Error()))
		task.bj4.logger.OnTaskError(task, err)
		task.errorChan <- err
	} else {
		task.SetStatus(fmt.Sprintf("completed: %s", result))
		task.bj4.logger.OnTaskComplete(task, result)
		task.errorChan <- nil
	}
}
