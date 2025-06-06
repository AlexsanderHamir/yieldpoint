package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func worker(name string, highPriority bool, wg *sync.WaitGroup) {
	defer wg.Done()

	// Set this worker's priority
	yieldpoint.SetHighPriority(highPriority)
	fmt.Printf("%s starting with high priority: %v\n", name, highPriority)

	for i := 0; i < 3; i++ {
		fmt.Printf("%s iteration %d\n", name, i)

		// Workers will yield based on their priority
		yieldpoint.MaybeYield()

		// Simulate some work
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("%s completed\n", name)
}

func main() {
	var wg sync.WaitGroup

	// Start workers with different priorities
	wg.Add(2)
	go worker("Normal Priority Worker", false, &wg)
	go worker("High Priority Worker", true, &wg)

	// Wait for all workers to complete
	wg.Wait()
	fmt.Println("All workers completed")
}
