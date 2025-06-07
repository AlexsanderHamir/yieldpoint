package main

import (
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func main() {
	var wg sync.WaitGroup

	// Start a background task that uses WaitIfActive
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 5 {
			yieldpoint.WaitIfActive()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Start a second background task that uses MaybeYield
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 5 {
			yieldpoint.MaybeYield()
			time.Sleep(300 * time.Millisecond)
		}
	}()

	yieldpoint.EnterHighPriority()
	time.Sleep(800 * time.Millisecond)
	yieldpoint.ExitHighPriority()

	time.Sleep(400 * time.Millisecond)

	// Second high-priority task
	yieldpoint.EnterHighPriority()
	time.Sleep(600 * time.Millisecond)
	yieldpoint.ExitHighPriority()

	// Wait for all background tasks to complete
	wg.Wait()
}
