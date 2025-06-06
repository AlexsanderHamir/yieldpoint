package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func nestedHighPriorityTask(level int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Entering high-priority level %d\n", level)
	yieldpoint.EnterHighPriority()
	defer yieldpoint.ExitHighPriority()

	// Simulate some work at this priority level
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Completed work at level %d\n", level)

	// If we haven't reached the maximum nesting level, create another nested section
	if level < 2 {
		wg.Add(1)
		go nestedHighPriorityTask(level+1, wg)
	}

	// Wait a bit before exiting this level
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Exiting high-priority level %d\n", level)
}

func backgroundWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range 5 {
		fmt.Printf("Background worker iteration %d\n", i)
		yieldpoint.MaybeYield()
		
		// Simulate work
		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	var wg sync.WaitGroup

	// Start a background worker
	wg.Add(1)
	go backgroundWorker(&wg)

	// Start the nested high-priority tasks
	wg.Add(1)
	go nestedHighPriorityTask(0, &wg)

	// Wait for all tasks to complete
	wg.Wait()
	fmt.Println("All tasks completed")
}
