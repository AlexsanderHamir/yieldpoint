package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

type Job struct {
	Name         string
	HighPriority bool
	Duration     time.Duration
}

func processJob(job Job, wg *sync.WaitGroup) {
	defer wg.Done()

	// Set the job's priority
	yieldpoint.SetHighPriority(job.HighPriority)
	fmt.Printf("Starting job: %s (High Priority: %v)\n", job.Name, job.HighPriority)

	// Process the job
	for i := 0; i < 3; i++ {
		fmt.Printf("Job %s: step %d\n", job.Name, i+1)

		// Yield to higher priority jobs if needed
		yieldpoint.MaybeYield()

		// Simulate work
		time.Sleep(job.Duration)
	}

	fmt.Printf("Completed job: %s\n", job.Name)
}

func main() {
	// Define some jobs with different priorities
	jobs := []Job{
		{Name: "System Update", HighPriority: true, Duration: 50 * time.Millisecond},
		{Name: "User Request", HighPriority: true, Duration: 100 * time.Millisecond},
		{Name: "Background Sync", HighPriority: false, Duration: 150 * time.Millisecond},
		{Name: "Log Cleanup", HighPriority: false, Duration: 200 * time.Millisecond},
	}

	var wg sync.WaitGroup

	// Process all jobs concurrently
	for _, job := range jobs {
		wg.Add(1)
		go processJob(job, &wg)
	}

	// Wait for all jobs to complete
	wg.Wait()
	fmt.Println("All jobs completed")
}
