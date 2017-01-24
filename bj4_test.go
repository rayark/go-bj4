package bj4

import (
	"fmt"
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
