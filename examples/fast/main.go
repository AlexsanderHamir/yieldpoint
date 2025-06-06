package main

import (
	"sync"
	"testing"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)


// ExampleConcurrentProcessing demonstrates concurrent processing with priority
func ExampleConcurrentProcessing() {
	var wg sync.WaitGroup

	yieldpoint.SetSpinWaitIterations(1000)
	yieldpoint.SetDefaultYieldDuration(10 * time.Millisecond)

	// Start some background workers
	for i := range 3 {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				// Use MaybeYieldFast in tight loops for better performance
				yieldpoint.MaybeYieldFast()

				// ... do background processing ...

				// Simulate some work
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Start a high-priority task
	go func() {
		yieldpoint.EnterHighPriority()
		defer yieldpoint.ExitHighPriority()

		// Critical work that needs to run quickly
		time.Sleep(100 * time.Millisecond)
	}()

	// Wait for some time to see the behavior
	time.Sleep(1 * time.Second)
}

// TestWaitIfActive demonstrates the difference between WaitIfActive and WaitIfActiveFast
func TestWaitIfActive(t *testing.T) {
	// Start a high-priority section
	yieldpoint.EnterHighPriority()

	// Start a goroutine that will wait
	done := make(chan struct{})
	go func() {
		// Use WaitIfActiveFast for short expected waits
		yieldpoint.WaitIfActiveFast()
		close(done)
	}()

	// Simulate some high-priority work
	time.Sleep(10 * time.Millisecond)

	// End the high-priority section
	yieldpoint.ExitHighPriority()

	// Wait for the waiting goroutine to complete
	select {
	case <-done:
		// Success
	case <-time.After(time.Second):
		t.Fatal("WaitIfActiveFast timed out")
	}
}

