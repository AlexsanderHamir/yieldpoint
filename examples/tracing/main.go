package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func tracedWorker(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s starting\n", name)

	for i := range 3 {
		fmt.Printf("%s iteration %d\n", name, i)

		// This yield will be traced
		yieldpoint.MaybeYield()

		// Simulate work
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("%s completed\n", name)
}

func main() {
	// Set up tracing
	yieldpoint.SetTraceFunc(func(e yieldpoint.YieldEvent) {
		fmt.Printf("TRACE: Goroutine %d yielded at %v (duration: %v)\n",
			e.GoroutineID, e.Timestamp.Format("15:04:05.000"), e.Duration)
	})

	var wg sync.WaitGroup

	// Start a high-priority section
	yieldpoint.EnterHighPriority()
	defer yieldpoint.ExitHighPriority()

	// Start some workers that will yield
	wg.Add(2)
	go tracedWorker("Worker 1", &wg)
	go tracedWorker("Worker 2", &wg)

	// Wait for workers to complete
	wg.Wait()
	fmt.Println("All workers completed")
}
