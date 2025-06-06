package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func main() {
	var wg sync.WaitGroup
	yieldpoint.SetDefaultYieldDuration(10 * time.Millisecond)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 10 {
			fmt.Printf("Background task iteration %d\n", i)
			// This will block until no high-priority tasks are running
			yieldpoint.MaybeYield()

			// Simulate work
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Simulate a high-priority task
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Starting high-priority task...")

	yieldpoint.EnterHighPriority()
	// Do some high-priority work
	for i := range 10 {
		fmt.Printf("High-priority task iteration %d\n", i)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println("High-priority task completed")

	// Exit high priority before waiting for background task
	yieldpoint.ExitHighPriority()

	// Now wait for background task to complete
	wg.Wait()
}
