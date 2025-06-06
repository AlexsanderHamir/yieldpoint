// Package yieldpoint provides cooperative goroutine yielding based on priority-aware scheduling.
package yieldpoint

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// HighPriorityCount tracks the number of active high-priority sections
var HighPriorityCount atomic.Int32

// Mu is the mutex used for efficient blocking in WaitIfActive
var Mu sync.Mutex

// Cond is the condition variable used for efficient blocking
var Cond = sync.NewCond(&Mu)

// DefaultYieldDuration is the default duration to sleep when yielding
var DefaultYieldDuration = 1 * time.Millisecond

// SpinWaitIterations is the number of iterations to spin-wait before falling back to mutex-based waiting
var SpinWaitIterations = 1000

// SetDefaultYieldDuration sets the default duration to sleep when yielding
func SetDefaultYieldDuration(d time.Duration) {
	DefaultYieldDuration = d
}

// SetSpinWaitIterations sets the number of iterations to spin-wait before falling back to mutex-based waiting
func SetSpinWaitIterations(n int) {
	SpinWaitIterations = n
}

// MaybeYield voluntarily yields the current goroutine if any high-priority sections are active.
// This is a non-blocking operation that uses runtime.Gosched() combined with a small sleep
// to ensure effective processor yielding.
func MaybeYield() {
	if HighPriorityCount.Load() > 0 {
		// First try to yield using runtime.Gosched()
		runtime.Gosched()

		// Then sleep for a small duration to ensure the processor is actually yielded
		time.Sleep(DefaultYieldDuration)
	}
}

// EnterHighPriority begins a high-priority section.
// Multiple calls are supported through reference counting.
func EnterHighPriority() {
	HighPriorityCount.Add(1)
}

// ExitHighPriority ends a high-priority section.
// If this is the last high-priority section, it will signal any waiting goroutines.
func ExitHighPriority() {
	count := HighPriorityCount.Add(-1)
	if count == 0 {
		Mu.Lock()
		Cond.Broadcast()
		Mu.Unlock()
	} else if count < 0 {
		// Reset to 0 if we somehow went negative
		HighPriorityCount.Store(0)
	}
}

// WaitIfActive blocks the current goroutine until no high-priority sections are active.
// This is an efficient blocking operation that uses sync.Cond to avoid busy waiting.
func WaitIfActive() {
	for HighPriorityCount.Load() > 0 {
		Mu.Lock()
		Cond.Wait()
		Mu.Unlock()
	}
}

// IsHighPriorityActive returns true if any high-priority sections are currently active.
func IsHighPriorityActive() bool {
	return HighPriorityCount.Load() > 0
}

// MaybeYieldFast is a high-performance version of MaybeYield that avoids time.Sleep
// and uses only runtime.Gosched() for minimal overhead. This is suitable for
// performance-critical code paths where the exact timing of yields is less important.
func MaybeYieldFast() {
	if HighPriorityCount.Load() > 0 {
		runtime.Gosched()
	}
}

// WaitIfActiveFast is a high-performance version of WaitIfActive that uses a spin-wait
// strategy before falling back to mutex-based waiting. This is suitable for
// performance-critical code paths where the wait time is expected to be very short.
func WaitIfActiveFast() {
	// First try spin-waiting
	for range SpinWaitIterations {
		if HighPriorityCount.Load() == 0 {
			return
		}
		runtime.Gosched()
	}

	// Only fall back to mutex-based waiting if spin-wait didn't succeed
	Mu.Lock()
	for HighPriorityCount.Load() > 0 {
		Cond.Wait()
	}
	Mu.Unlock()
}


// MaybeYieldWithContext is a context-aware version of MaybeYield
func MaybeYieldWithContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		MaybeYield()
		return nil
	}
}

// WaitIfActiveWithContext is a context-aware version of WaitIfActive
func WaitIfActiveWithContext(ctx context.Context) error {
	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if HighPriorityCount.Load() == 0 {
				return nil
			}
		}
	}
}
