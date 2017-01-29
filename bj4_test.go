package bj4

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

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
	}, time.Now().Add(3*time.Second))

	sch.SetScheduledTask("2", func(task *Task) (result string, nextUpdate time.Time, err error) {
		seq = append(seq, 2)
		return
	}, time.Now().Add(2*time.Second))

	time.Sleep(4 * time.Second)

	if !reflect.DeepEqual(seq, []int64{2, 1}) {
		t.Error("wrong sequence:", seq)
	}
}
