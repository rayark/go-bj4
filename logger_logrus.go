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
