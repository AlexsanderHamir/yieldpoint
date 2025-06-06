package yieldpoint

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

// YieldEvent represents a single yield event in the system
type YieldEvent struct {
	// GoroutineID is the ID of the goroutine that yielded
	GoroutineID uint64
	// Timestamp is when the yield occurred
	Timestamp time.Time
	// Duration is how long the yield lasted (if applicable)
	Duration time.Duration
	// Reason describes why the yield occurred
	Reason string
	// IsHighPriority indicates if the yielding goroutine has high priority
	IsHighPriority bool
}

var (
	// traceFunc is the callback function for yield events
	traceFunc atomic.Value
)

// SetTraceFunc sets a callback function that will be called for each yield event.
// The function can be nil to disable tracing.
func SetTraceFunc(fn func(YieldEvent)) {
	traceFunc.Store(fn)
}

// traceYieldEvent records a yield event if tracing is enabled
func traceYieldEvent(reason string, duration time.Duration) {
	if fn, ok := traceFunc.Load().(func(YieldEvent)); ok && fn != nil {
		fn(YieldEvent{
			GoroutineID:    getGoroutineID(),
			Timestamp:      time.Now(),
			Duration:       duration,
			Reason:         reason,
			IsHighPriority: GetHighPriority(),
		})
	}
}

// getGoroutineID returns the current goroutine's ID
func getGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	var id uint64
	// Parse the goroutine ID from the stack trace
	// Format: "goroutine N [status]:"
	_, err := fmt.Sscanf(string(b), "goroutine %d", &id)
	if err != nil {
		return 0
	}
	return id
}
