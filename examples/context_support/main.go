package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func workerWithContext(ctx context.Context, name string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s starting\n", name)

	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("%s cancelled: %v\n", name, ctx.Err())
			return
		default:
			fmt.Printf("%s iteration %d\n", name, i)

			// Try to yield with context
			err := yieldpoint.MaybeYieldWithContext(ctx)
			if err != nil {
				fmt.Printf("%s yield cancelled: %v\n", name, err)
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Printf("%s completed normally\n", name)
}

func main() {
	var wg sync.WaitGroup

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Start a high-priority section
	yieldpoint.EnterHighPriority()
	defer yieldpoint.ExitHighPriority()

	// Start workers that will be affected by the context
	wg.Add(2)
	go workerWithContext(ctx, "Worker 1", &wg)
	go workerWithContext(ctx, "Worker 2", &wg)

	// Wait for workers to complete or timeout
	wg.Wait()
	fmt.Println("Main function completed")
}
