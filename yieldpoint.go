// Package yieldpoint provides cooperative goroutine yielding based on priority-aware scheduling.
package yieldpoint

import (
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

// SetDefaultYieldDuration sets the default duration to sleep when yielding
func SetDefaultYieldDuration(d time.Duration) {
	DefaultYieldDuration = d
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

		traceYieldEvent("high_priority_active", DefaultYieldDuration)
	}
}

// EnterHighPriority begins a high-priority section.
// Multiple calls are supported through reference counting.
func EnterHighPriority() {
	HighPriorityCount.Add(1)
	traceYieldEvent("enter_high_priority", 0)
}

// ExitHighPriority ends a high-priority section.
// If this is the last high-priority section, it will signal any waiting goroutines.
func ExitHighPriority() {
	if HighPriorityCount.Add(-1) == 0 {
		Mu.Lock()
		Cond.Broadcast()
		Mu.Unlock()
		traceYieldEvent("exit_high_priority", 0)
	}
}

// WaitIfActive blocks the current goroutine until no high-priority sections are active.
// This is an efficient blocking operation that uses sync.Cond to avoid busy waiting.
func WaitIfActive() {
	start := time.Now()
	for HighPriorityCount.Load() > 0 {
		Mu.Lock()
		Cond.Wait()
		Mu.Unlock()
	}
	traceYieldEvent("wait_complete", time.Since(start))
}

// IsHighPriorityActive returns true if any high-priority sections are currently active.
func IsHighPriorityActive() bool {
	return HighPriorityCount.Load() > 0
}
