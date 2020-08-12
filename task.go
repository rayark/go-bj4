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
