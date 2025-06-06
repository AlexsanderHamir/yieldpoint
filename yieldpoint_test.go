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


func TestConcurrentHighPriority(t *testing.T) {
	// Test multiple goroutines entering and exiting high priority simultaneously
	var wg sync.WaitGroup
	const numGoroutines = 10
	const iterations = 100

	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for range iterations {
				EnterHighPriority()
				// Verify we can detect high priority is active
				if !IsHighPriorityActive() {
					t.Errorf("High priority not detected after EnterHighPriority in goroutine %d", id)
				}
				ExitHighPriority()
			}
		}(i)
	}

	wg.Wait()

	// Verify high priority count is back to 0
	if IsHighPriorityActive() {
		t.Error("High priority still active after all goroutines finished")
	}
}

func TestNestedHighPriority(t *testing.T) {
	// Test nested high priority sections
	EnterHighPriority()
	EnterHighPriority()
	EnterHighPriority()

	// Should still be active after multiple enters
	if !IsHighPriorityActive() {
		t.Error("High priority not active after multiple EnterHighPriority calls")
	}

	ExitHighPriority()
	ExitHighPriority()
	ExitHighPriority()

	// Should be inactive after matching exits
	if IsHighPriorityActive() {
		t.Error("High priority still active after matching ExitHighPriority calls")
	}

	// Test that exiting more times than entering doesn't cause issues
	ExitHighPriority()
	ExitHighPriority()
	if IsHighPriorityActive() {
		t.Error("High priority active after extra ExitHighPriority calls")
	}
}

func TestMaybeYieldFast(t *testing.T) {
	// Test that MaybeYieldFast doesn't block when no high priority is active
	done := make(chan struct{})
	go func() {
		MaybeYieldFast()
		close(done)
	}()

	select {
	case <-done:
		// Success - MaybeYieldFast didn't block
	case <-time.After(time.Second):
		t.Fatal("MaybeYieldFast blocked when no high priority was active")
	}

	// Test that MaybeYieldFast yields when high priority is active
	EnterHighPriority()
	defer ExitHighPriority()

	start := time.Now()
	MaybeYieldFast()
	duration := time.Since(start)

	// runtime.Gosched() should be very fast
	if duration > time.Millisecond {
		t.Errorf("MaybeYieldFast took too long: %v", duration)
	}
}

func TestConcurrentWaitAndYield(t *testing.T) {
	const numWaiters = 5
	const numYielders = 5
	var wg sync.WaitGroup

	// Start multiple waiters and yielders
	for i := 0; i < numWaiters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			WaitIfActive()
		}(i)
	}

	for i := range numYielders {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			MaybeYield()
		}(i)
	}

	// Enter high priority
	EnterHighPriority()

	// Give time for goroutines to start
	time.Sleep(time.Millisecond)

	// Exit high priority
	ExitHighPriority()

	// Wait for all goroutines to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success - all goroutines completed
	case <-time.After(time.Second):
		t.Fatal("Not all goroutines completed after high priority ended")
	}
}

func TestContextCancellation(t *testing.T) {
	// Test context cancellation during MaybeYieldWithContext
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := MaybeYieldWithContext(ctx)
	if err != context.Canceled {
		t.Errorf("MaybeYieldWithContext expected Canceled error, got: %v", err)
	}

	// Test context cancellation during WaitIfActiveWithContext
	ctx, cancel = context.WithCancel(context.Background())
	EnterHighPriority()
	defer ExitHighPriority()

	done := make(chan struct{})
	go func() {
		err := WaitIfActiveWithContext(ctx)
		if err != context.Canceled {
			t.Errorf("WaitIfActiveWithContext expected Canceled error, got: %v", err)
		}
		close(done)
	}()

	// Give time for goroutine to start waiting
	time.Sleep(time.Millisecond)

	// Cancel the context
	cancel()

	select {
	case <-done:
		// Success - WaitIfActiveWithContext responded to cancellation
	case <-time.After(time.Second):
		t.Fatal("WaitIfActiveWithContext didn't respond to context cancellation")
	}
}

func TestConfigurationChanges(t *testing.T) {
	// Test changing configuration values
	originalYieldDuration := DefaultYieldDuration
	originalSpinIterations := SpinWaitIterations

	// Set new values
	SetDefaultYieldDuration(2 * time.Millisecond)
	SetSpinWaitIterations(2000)

	// Verify changes took effect
	if DefaultYieldDuration != 2*time.Millisecond {
		t.Errorf("DefaultYieldDuration not updated, got: %v", DefaultYieldDuration)
	}
	if SpinWaitIterations != 2000 {
		t.Errorf("SpinWaitIterations not updated, got: %v", SpinWaitIterations)
	}

	// Test that changes affect behavior
	EnterHighPriority()
	start := time.Now()
	MaybeYield()
	duration := time.Since(start)

	// Should sleep for approximately the new duration
	if duration < DefaultYieldDuration {
		t.Errorf("MaybeYield slept for less than configured duration: %v", duration)
	}

	// Restore original values
	SetDefaultYieldDuration(originalYieldDuration)
	SetSpinWaitIterations(originalSpinIterations)
	ExitHighPriority()
}
