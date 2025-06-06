package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func main() {
	var wg sync.WaitGroup

	// Start a background task that uses WaitIfActive
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 5 {
			fmt.Printf("Background task iteration %d\n", i)

			// Wait for any active high-priority tasks to complete
			fmt.Println("Background task: Waiting for any active high-priority tasks...")
			yieldpoint.WaitIfActive()
			fmt.Println("Background task: No high-priority tasks active, proceeding...")

			// Simulate work
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Start a second background task that uses MaybeYield
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 5 {
			fmt.Printf("Second background task iteration %d\n", i)

			// This will yield if there are high-priority tasks
			yieldpoint.MaybeYield()

			// Simulate work
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// Simulate some high-priority tasks at different times
	time.Sleep(200 * time.Millisecond)

	// First high-priority task
	fmt.Println("\nStarting first high-priority task...")
	yieldpoint.EnterHighPriority()
	fmt.Println("First high-priority task: Doing some work...")
	time.Sleep(800 * time.Millisecond)
	fmt.Println("First high-priority task completed")
	yieldpoint.ExitHighPriority()

	time.Sleep(400 * time.Millisecond)

	// Second high-priority task
	fmt.Println("\nStarting second high-priority task...")
	yieldpoint.EnterHighPriority()
	fmt.Println("Second high-priority task: Doing some work...")
	time.Sleep(600 * time.Millisecond)
	fmt.Println("Second high-priority task completed")
	yieldpoint.ExitHighPriority()

	// Wait for all background tasks to complete
	wg.Wait()
	fmt.Println("\nAll tasks completed!")
}
