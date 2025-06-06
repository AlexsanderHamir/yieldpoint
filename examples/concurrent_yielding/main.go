package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

// workerUsingWaitIfActive demonstrates a worker that blocks until high-priority tasks complete
func workerUsingWaitIfActive(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s starting (using WaitIfActive)\n", name)

	for i := range 5 {
		fmt.Printf("%s iteration %d\n", name, i)

		// This will block until no high-priority tasks are active
		fmt.Printf("%s: Checking for high-priority tasks...\n", name)
		yieldpoint.WaitIfActive()
		fmt.Printf("%s: No high-priority tasks active, proceeding...\n", name)

		// Simulate work
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("%s completed\n", name)
}

// workerUsingMaybeYield demonstrates a worker that yields but continues if high-priority tasks are active
func workerUsingMaybeYield(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s starting (using MaybeYield)\n", name)

	for i := range 5 {
		fmt.Printf("%s iteration %d\n", name, i)

		// This will yield if high-priority tasks are active, but continue anyway
		fmt.Printf("%s: Yielding if high-priority tasks are active...\n", name)
		yieldpoint.MaybeYield()
		fmt.Printf("%s: Continuing work...\n", name)

		// Simulate work
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("%s completed\n", name)
}

// highPriorityTask simulates a high-priority task that runs for a specified duration
func highPriorityTask(name string, duration time.Duration) {
	fmt.Printf("\nStarting high-priority task: %s\n", name)
	yieldpoint.EnterHighPriority()
	defer yieldpoint.ExitHighPriority()

	fmt.Printf("%s: Doing high-priority work...\n", name)
	time.Sleep(duration)
	fmt.Printf("%s: High-priority work completed\n", name)
}

func main() {
	// Set up tracing to see what's happening
	yieldpoint.SetTraceFunc(func(e yieldpoint.YieldEvent) {
		fmt.Printf("TRACE: Goroutine %d - %s (duration: %v)\n",
			e.GoroutineID, e.Reason, e.Duration)
	})

	var wg sync.WaitGroup

	// Start two workers using WaitIfActive
	wg.Add(2)
	go workerUsingWaitIfActive("Worker 1", &wg)
	go workerUsingWaitIfActive("Worker 2", &wg)

	// Start two workers using MaybeYield
	wg.Add(2)
	go workerUsingMaybeYield("Worker 3", &wg)
	go workerUsingMaybeYield("Worker 4", &wg)

	// Let the workers run for a bit
	time.Sleep(300 * time.Millisecond)

	// Start a short high-priority task
	go highPriorityTask("Short Task", 400*time.Millisecond)

	// Let some work happen
	time.Sleep(500 * time.Millisecond)

	// Start a longer high-priority task
	go highPriorityTask("Long Task", 800*time.Millisecond)

	// Wait for all workers to complete
	wg.Wait()
	fmt.Println("\nAll workers completed!")
}
