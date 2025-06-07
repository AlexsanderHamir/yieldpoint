package main

import (
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func nestedHighPriorityTask(level int, wg *sync.WaitGroup) {
	defer wg.Done()

	yieldpoint.EnterHighPriority()
	defer yieldpoint.ExitHighPriority()

	time.Sleep(100 * time.Millisecond)

	if level < 2 {
		wg.Add(1)
		go nestedHighPriorityTask(level+1, wg)
	}

	// Wait a bit before exiting this level
	time.Sleep(100 * time.Millisecond)
}

func backgroundWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	for range 5 {
		yieldpoint.MaybeYield()
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
}
