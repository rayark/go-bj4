package bj4

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func approxDuration(expected, actual time.Duration) bool {
	return actual < expected+1*time.Millisecond &&
		actual > expected-1*time.Millisecond
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
