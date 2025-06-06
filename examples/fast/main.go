package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

// ExampleHighPriorityTask demonstrates a high-priority task that needs to run quickly
func ExampleHighPriorityTask() {
	// Start a high-priority section
	yieldpoint.EnterHighPriority()
	defer yieldpoint.ExitHighPriority()

	// Perform critical work that needs to run quickly
	fmt.Println("Running high-priority task")
	// ... do critical work ...
}

// ExampleBackgroundTask demonstrates a background task that should yield to high-priority tasks
func ExampleBackgroundTask() {
	// This task will automatically yield when high-priority tasks are active
	for i := 0; i < 10; i++ {
		// Regular version - good for most cases
		yieldpoint.MaybeYield()

		// Or use the fast version for performance-critical sections
		// yieldpoint.MaybeYieldFast()

		// ... do background work ...
	}
}

// ExampleConcurrentProcessing demonstrates concurrent processing with priority
func ExampleConcurrentProcessing() {
	var wg sync.WaitGroup

	// Start some background workers
	for i := 0; i < 3; i++ {
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

// ExampleTightLoop demonstrates using MaybeYieldFast in a performance-critical tight loop
func ExampleTightLoop() {
	// This is a tight loop where we want minimal overhead
	for i := 0; i < 1000000; i++ {
		// Use MaybeYieldFast for minimal overhead
		yieldpoint.MaybeYieldFast()

		// ... do critical work ...
	}
}

// ExampleMixedPriorityWorkload demonstrates a mixed workload with both high and low priority tasks
func ExampleMixedPriorityWorkload() {
	var wg sync.WaitGroup

	// Start some low-priority background tasks
	for i := range 2 {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				// Regular MaybeYield for background tasks
				yieldpoint.MaybeYield()

				// ... do background work ...
				time.Sleep(50 * time.Millisecond)
			}
		}(i)
	}

	// Start a high-priority task that runs periodically
	go func() {
		for {
			yieldpoint.EnterHighPriority()

			// Critical work that needs to run quickly
			time.Sleep(10 * time.Millisecond)

			yieldpoint.ExitHighPriority()

			// Wait before next high-priority task
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Run for a while to see the behavior
	time.Sleep(1 * time.Second)
}

// ExampleRealTimeProcessing demonstrates a real-time processing scenario
func ExampleRealTimeProcessing() {
	// Start a real-time processing loop
	for {
		// Check if we need to wait for high-priority tasks
		if yieldpoint.IsHighPriorityActive() {
			// Use WaitIfActiveFast for real-time processing
			// where we expect waits to be very short
			yieldpoint.WaitIfActiveFast()
		}

		// Process real-time data
		// ... do real-time processing ...

		// Use MaybeYieldFast in the tight loop
		yieldpoint.MaybeYieldFast()
	}
}

func main() {
	// Run the mixed priority workload example to demonstrate the package's functionality
	ExampleRealTimeProcessing()
}
