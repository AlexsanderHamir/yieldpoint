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

func TestConcurrentHighPriority(t *testing.T) {
	// Test multiple goroutines entering and exiting high priority simultaneously
	var wg sync.WaitGroup
	const numGoroutines = 10
	const iterations = 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
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

	for i := 0; i < numYielders; i++ {
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

func TestTracingEvents(t *testing.T) {
	var events []YieldEvent
	SetTraceFunc(func(e YieldEvent) {
		events = append(events, e)
	})

	// Test tracing for various operations
	EnterHighPriority()
	MaybeYield()
	MaybeYieldFast()
	WaitIfActiveFast()
	ExitHighPriority()

	// Verify events were recorded
	if len(events) == 0 {
		t.Fatal("No yield events were recorded")
	}

	// Verify event types
	eventTypes := make(map[string]bool)
	for _, e := range events {
		eventTypes[e.Reason] = true
	}

	expectedTypes := map[string]bool{
		"enter_high_priority":       true,
		"high_priority_active":      true,
		"high_priority_active_fast": true,
		"wait_complete_fast":        true,
		"exit_high_priority":        true,
	}

	for expected := range expectedTypes {
		if !eventTypes[expected] {
			t.Errorf("Expected event type %s not found in trace", expected)
		}
	}

	// Verify event timestamps are in order
	for i := 1; i < len(events); i++ {
		if events[i].Timestamp.Before(events[i-1].Timestamp) {
			t.Errorf("Events out of order: %v before %v", events[i].Timestamp, events[i-1].Timestamp)
		}
	}

	// Test disabling tracing
	SetTraceFunc(nil)
	events = nil
	MaybeYield()
	if len(events) > 0 {
		t.Error("Events were recorded after disabling tracing")
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

func TestRaceConditions(t *testing.T) {
	// Test for race conditions in concurrent operations
	const numGoroutines = 10
	const iterations = 1000
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Mix of operations that could potentially race
				EnterHighPriority()
				MaybeYield()
				MaybeYieldFast()
				IsHighPriorityActive()
				ExitHighPriority()
				WaitIfActiveFast()
			}
		}()
	}

	// Run with race detector enabled
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success - no races detected
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out, possible deadlock")
	}
}

func TestGoroutineIDTracking(t *testing.T) {
	// Test that goroutine IDs are correctly tracked in trace events
	var events []YieldEvent
	SetTraceFunc(func(e YieldEvent) {
		events = append(events, e)
	})

	// Create multiple goroutines that yield
	const numGoroutines = 5
	var wg sync.WaitGroup
	goroutineIDs := make(map[uint64]bool)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			EnterHighPriority()
			MaybeYield()
			ExitHighPriority()
		}()
	}

	wg.Wait()

	// Verify that different goroutine IDs were recorded
	for _, e := range events {
		if e.GoroutineID == 0 {
			t.Error("Zero goroutine ID recorded in trace event")
		}
		goroutineIDs[e.GoroutineID] = true
	}

	if len(goroutineIDs) < numGoroutines {
		t.Errorf("Expected at least %d different goroutine IDs, got %d", numGoroutines, len(goroutineIDs))
	}
}
