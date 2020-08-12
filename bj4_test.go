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
	"reflect"
	"sync"
	"testing"
	"time"
)

func approxDuration(expected, actual time.Duration) bool {
	return actual < expected+5*time.Millisecond &&
		actual > expected-5*time.Millisecond
}

func ExampleBJ4_SetTask() {
	sch := New(&Config{})

	errChan := sch.SetTask("hello", func(task *Task) (result string, nextUpdate time.Time, err error) {
		fmt.Println("Hello World")
		result = "done"
		return
	})

	go sch.Start()

	<-errChan // Wait for the task to complete

	// Output: Hello World
}

func ExampleBJ4_SetScheduledTask() {
	sch := New(&Config{})

	errChan := sch.SetScheduledTask("hello", func(task *Task) (result string, nextUpdate time.Time, err error) {
		fmt.Println("Hello World")
		result = "done"
		return
	}, time.Now().Add(3*time.Second))

	go sch.Start()

	<-errChan // Wait for the task to complete

	// Output: Hello World
}

func ExampleBJ4_SetScheduledTask_repeated() {
	sch := New(&Config{})

	counter := 0

	sch.SetScheduledTask("hello", func(task *Task) (result string, nextUpdate time.Time, err error) {
		counter++
		fmt.Println("counter:", counter)
		result = "done"
		nextUpdate = time.Now().Add(2 * time.Second)
		return
	}, time.Now().Add(3*time.Second))

	go sch.Start()

	time.Sleep(10 * time.Second)

	// Output: counter: 1
	// counter: 2
	// counter: 3
	// counter: 4
}

func TestBJ4(t *testing.T) {
	var seq []int64

	sch := New(&Config{})
	go sch.Start()

	sch.SetScheduledTask("1", func(task *Task) (result string, nextUpdate time.Time, err error) {
		seq = append(seq, 1)
		return
	}, time.Now().Add(200*time.Millisecond))

	sch.SetScheduledTask("2", func(task *Task) (result string, nextUpdate time.Time, err error) {
		seq = append(seq, 2)
		return
	}, time.Now().Add(100*time.Millisecond))

	time.Sleep(500 * time.Millisecond)

	if !reflect.DeepEqual(seq, []int64{2, 1}) {
		t.Error("wrong sequence:", seq)
	}
}

func TestMinWaitTime(t *testing.T) {
	wt := 500 * time.Millisecond
	sch := New(&Config{MinWaitTime: wt})
	s := time.Now()
	sch.wait()
	d := time.Now().Sub(s)

	if !approxDuration(wt, d) {
		t.Error("wrong duration. expected:", wt, ", actual:", d)
	}
}

func TestWait(t *testing.T) {
	wt := 1000 * time.Millisecond
	tt := 500 * time.Millisecond
	sch := New(&Config{MinWaitTime: wt})

	sch.SetScheduledTask("1", func(task *Task) (result string, nextUpdate time.Time, err error) {
		return
	}, time.Now().Add(tt))

	s := time.Now()
	sch.wait()
	d := time.Now().Sub(s)

	if !approxDuration(tt, d) {
		t.Error("wrong duration. expected:", tt, ", actual:", d)
	}
}

func TestStop(t *testing.T) {
	var seq []int64

	sch := New(&Config{})
	go sch.Start()
	time.AfterFunc(200*time.Millisecond, func() {
		sch.Stop()
	})

	sch.SetScheduledTask("1", func(task *Task) (result string, nextUpdate time.Time, err error) {
		seq = append(seq, 1)
		return
	}, time.Now().Add(300*time.Millisecond))

	sch.SetScheduledTask("2", func(task *Task) (result string, nextUpdate time.Time, err error) {
		seq = append(seq, 2)
		return
	}, time.Now().Add(100*time.Millisecond))

	time.Sleep(500 * time.Millisecond)

	if !reflect.DeepEqual(seq, []int64{2}) {
		t.Error("wrong sequence:", seq)
	}
}

func TestRemove(t *testing.T) {
	var seq []int64

	sch := New(&Config{})
	go sch.Start()

	sch.SetScheduledTask("1", func(task *Task) (result string, nextUpdate time.Time, err error) {
		seq = append(seq, 1)
		return
	}, time.Now().Add(300*time.Millisecond))

	sch.SetScheduledTask("2", func(task *Task) (result string, nextUpdate time.Time, err error) {
		seq = append(seq, 2)
		return
	}, time.Now().Add(100*time.Millisecond))

	sch.RemoveTask("1")

	time.Sleep(500 * time.Millisecond)

	if !reflect.DeepEqual(seq, []int64{2}) {
		t.Error("wrong sequence:", seq)
	}
}

func TestStopWaitsCurrentTaskDone(t *testing.T) {
	var count int
	var wg sync.WaitGroup

	sch := New(&Config{})
	go sch.Start()

	wg.Add(1)
	// will fire after 400 ms
	sch.SetScheduledTask("1", func(task *Task) (result string, nextUpdate time.Time, err error) {
		defer wg.Done()
		time.Sleep(time.Millisecond * 300)
		count++
		return
	}, time.Now().Add(100*time.Millisecond))

	wg.Add(1)
	time.AfterFunc(200*time.Millisecond, func() {
		defer wg.Done()

		sch.Stop()

		// Stop should guarantee on-going Task is done
		if count != 1 {
			t.Error("Stop() doesn't wait for on-going Task done")
		}
	})

	wg.Wait()
}
