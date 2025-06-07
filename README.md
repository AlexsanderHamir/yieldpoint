# yieldpoint

A Go package for cooperative goroutine yielding with priority-aware scheduling.

## Overview

`yieldpoint` enables goroutines to voluntarily yield execution when high-priority tasks are active, using atomic operations and condition variables for efficient synchronization.

## Features

- **Priority-based Yielding**: Voluntary yielding when high-priority tasks are active
- **Efficient Blocking**: Uses `sync.Cond` for non-busy waiting
- **Context Support**: Timeout-aware operations with context
- **Thread Safety**: Atomic operations for high-priority counting
- **Performance Optimizations**: Fast variants with spin-wait strategies
- **Configurable**: Adjustable spin-wait iterations and yield durations

## Installation

```bash
go get github.com/AlexsanderHamir/yieldpoint
```

## Usage

### Basic Usage

```go
package main

import "github.com/AlexsanderHamir/yieldpoint"

func main() {
    // High-priority section
    yieldpoint.EnterHighPriority()
    defer yieldpoint.ExitHighPriority()

    go func() {
        // Standard variants
        yieldpoint.MaybeYield()      // Quick yield if high-priority active
        yieldpoint.WaitIfActive()    // Block until high-priority ends

        // Fast variant for performance-critical paths
        yieldpoint.WaitIfActiveFast() // Spin-wait before blocking
    }()
}
```

### Context Support

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

// Non-blocking yield with timeout
if err := yieldpoint.MaybeYieldWithContext(ctx); err != nil {
    // Handle timeout/cancellation
}

// Blocking wait with timeout
if err := yieldpoint.WaitIfActiveWithContext(ctx); err != nil {
    // Handle timeout/cancellation
}
```

### Nested High-Priority

```go
// Reference counting for nested sections
yieldpoint.EnterHighPriority() // Count = 1
yieldpoint.EnterHighPriority() // Count = 2
defer yieldpoint.ExitHighPriority() // Count = 1
defer yieldpoint.ExitHighPriority() // Count = 0, signals waiters
```

## API Reference

### Core Functions

- `MaybeYield()`: Yields if high-priority active
- `EnterHighPriority()`: Begins high-priority section (reference counted)
- `ExitHighPriority()`: Ends high-priority section, signals if last
- `WaitIfActive()`: Blocks until high-priority ends
- `IsHighPriorityActive()`: Checks high-priority status

### Performance Variants

- `WaitIfActiveFast()`: Spin-wait strategy for short waits
  - Configurable via `SetSpinWaitIterations`
  - Falls back to mutex-based waiting

### Context Functions

- `MaybeYieldWithContext(ctx)`: Non-blocking yield with timeout
- `WaitIfActiveWithContext(ctx)`: Blocking wait with timeout

### Configuration

- `SetSpinWaitIterations(n int)`: Configure spin-wait behavior

## Performance

- Use fast variants in performance-critical paths
- Tune `SpinWaitIterations` based on wait duration:
  - Higher: Better for very short waits
  - Lower: Better for longer waits
- Set appropriate timeouts for context operations

## Thread Safety

- Atomic high-priority counting
- Mutex and condition variable for blocking
- Safe for concurrent use

## Contribution

Share your talents and ideas!!

## License

MIT License - see [LICENSE](LICENSE)
