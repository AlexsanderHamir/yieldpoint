package yieldpoint

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMaybeYield(t *testing.T) {
	// Test that MaybeYield doesn't block when no high priority is active
	done := make(chan struct{})
	go func() {
		MaybeYield()
		close(done)
	}()

	select {
	case <-done:
		// Success - MaybeYield didn't block
	case <-time.After(time.Second):
		t.Fatal("MaybeYield blocked when no high priority was active")
	}

	// Test that MaybeYield yields when high priority is active
	EnterHighPriority()
	defer ExitHighPriority()

	start := time.Now()
	MaybeYield()
	duration := time.Since(start)

	// runtime.Gosched() should take some time, but not too much
	// Note: runtime.Gosched() can be very fast on modern systems
	if duration > time.Second {
		t.Errorf("MaybeYield took too long: %v", duration)
	}
}

func TestWaitIfActive(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	// Start a goroutine that will wait for high priority to end
	go func() {
		defer wg.Done()
		WaitIfActive()
	}()

	// Enter high priority
	EnterHighPriority()

	// Give the waiting goroutine time to start waiting
	time.Sleep(time.Millisecond)

	// Exit high priority should unblock the waiting goroutine
	ExitHighPriority()

	// Wait for the goroutine to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success - WaitIfActive unblocked
	case <-time.After(time.Second):
		t.Fatal("WaitIfActive didn't unblock after high priority ended")
	}
}

func TestContextAwareFunctions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Test MaybeYieldWithContext with cancellation
	EnterHighPriority()
	defer ExitHighPriority()

	err := MaybeYieldWithContext(ctx)
	if err != nil {
		t.Errorf("MaybeYieldWithContext returned unexpected error: %v", err)
	}

	// Test WaitIfActiveWithContext with cancellation
	err = WaitIfActiveWithContext(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("WaitIfActiveWithContext expected DeadlineExceeded, got: %v", err)
	}
}

func TestPriorityLevels(t *testing.T) {
	// Test default priority
	if got := GetHighPriority(); got != false {
		t.Errorf("GetHighPriority() = %v, want false", got)
	}

	// Test setting and getting priority
	SetHighPriority(true)
	if got := GetHighPriority(); got != true {
		t.Errorf("GetHighPriority() = %v, want true", got)
	}

	SetHighPriority(false)
	if got := GetHighPriority(); got != false {
		t.Errorf("GetHighPriority() = %v, want false", got)
	}
}

func TestTracing(t *testing.T) {
	var events []YieldEvent
	SetTraceFunc(func(e YieldEvent) {
		events = append(events, e)
	})

	// Trigger some yield events
	MaybeYield()
	EnterHighPriority()
	MaybeYield()
	ExitHighPriority()

	// Verify that events were recorded
	if len(events) == 0 {
		t.Error("No yield events were recorded")
	}

	// Disable tracing
	SetTraceFunc(nil)
	MaybeYield()
	if len(events) == 0 {
		t.Error("Events should have been recorded before disabling tracing")
	}
}
