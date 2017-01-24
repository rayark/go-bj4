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

func ExampleBJ4_SetScheduledTask_repeated() {
	sch := New(&Config{})

	counter := 0

	sch.SetScheduledTask("hello", func(task *Task) (result string, nextUpdate time.Time, err error) {
		counter++
		fmt.Println("counter: ", counter)
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
