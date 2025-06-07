package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func workerWithYieldContext(ctx context.Context, name string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s starting\n", name)

	for i := range 5 {
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

			// Simulate work
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Printf("%s completed normally\n", name)
}

func workerWithWaitContext(ctx context.Context, name string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s starting\n", name)

	// Try to wait for high priority to finish, with context
	fmt.Printf("%s waiting for high priority to finish...\n", name)
	err := yieldpoint.WaitIfActiveWithContext(ctx)
	if err != nil {
		fmt.Printf("%s wait cancelled: %v\n", name, err)
		return
	}

	fmt.Printf("%s high priority finished, proceeding with work\n", name)
	for i := range 3 {
		fmt.Printf("%s doing work %d\n", name, i)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("%s completed normally\n", name)
}

func main() {
	var wg sync.WaitGroup

	// Create contexts with different timeouts
	yieldCtx, yieldCancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer yieldCancel()

	waitCtx, waitCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer waitCancel()

	// Start a high-priority section
	yieldpoint.EnterHighPriority()
	defer yieldpoint.ExitHighPriority()

	// Start workers that use MaybeYieldWithContext
	wg.Add(2)
	go workerWithYieldContext(yieldCtx, "Yield Worker 1", &wg)
	go workerWithYieldContext(yieldCtx, "Yield Worker 2", &wg)

	// Start workers that use WaitIfActiveWithContext
	wg.Add(2)
	go workerWithWaitContext(waitCtx, "Wait Worker 1", &wg)
	go workerWithWaitContext(waitCtx, "Wait Worker 2", &wg)

	// Simulate some high-priority work
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Exiting high priority section...")
	yieldpoint.ExitHighPriority()

	// Wait for all workers to complete or timeout
	wg.Wait()
	fmt.Println("Main function completed")
}
