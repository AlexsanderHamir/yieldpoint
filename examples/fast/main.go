package main

import (
	"testing"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

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

