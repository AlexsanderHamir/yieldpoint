package yieldpoint

import (
	"context"
	"sync/atomic"
	"time"
)

// IsHighPriority stores whether the current goroutine has high priority
// This is implemented using a goroutine-local storage pattern
var IsHighPriority atomic.Value

// SetHighPriority sets the high priority flag for the current goroutine
func SetHighPriority(high bool) {
	IsHighPriority.Store(high)
}

// GetHighPriority returns whether the current goroutine has high priority
func GetHighPriority() bool {
	if high, ok := IsHighPriority.Load().(bool); ok {
		return high
	}
	return false // Default to normal priority
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
