package main

import (
	"fmt"
	"time"
)

func main() {
	scheduler := NewTaskScheduler()

	scheduler.AddTask("Buy Groceries", time.Now().Add(time.Second*5), func() error {
		defer fmt.Println("groceries shopping done!")
		fmt.Println("buying groceries")
		return nil
	})

	// Schedule another task to run 10 seconds later
	scheduler.AddTask("Pay Bills", time.Now().Add(10*time.Second), func() error {
		fmt.Println("Paying bills...")
		return nil // No error, success
	})

	go scheduler.ScheduleWorker()

	select {}
}
