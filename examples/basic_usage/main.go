package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 10 {
			fmt.Printf("Background task iteration %d\n", i)
			// This will yield once per iteration if high-priority tasks are active
			yieldpoint.MaybeYield()

			// Simulate work
			time.Sleep(100 * time.Millisecond)
		}
	}()

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
