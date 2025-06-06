# yieldpoint

A Go package that enables cooperative goroutine yielding based on priority-aware scheduling.

## Overview

`yieldpoint` provides a simple yet powerful mechanism for implementing cooperative multitasking in Go applications. It allows goroutines to voluntarily yield execution when high-priority tasks are active, making it ideal for:

- Game engines
- Real-time systems
- Job schedulers
- Any application requiring fine-grained control over goroutine scheduling

## Features

- **Priority-based Yielding**: Goroutines can yield execution when high-priority tasks are active
- **Efficient Blocking**: Uses `sync.Cond` for efficient blocking without busy waiting
- **Context Support**: Context-aware variants of all operations
- **High Priority Support**: Simple boolean flag for high-priority tasks
- **Tracing**: Optional instrumentation for observing yield events
- **Thread Safety**: All operations are thread-safe
- **Nesting Support**: High-priority sections can be nested

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
    // Start a high-priority section
    yieldpoint.EnterHighPriority()
    defer yieldpoint.ExitHighPriority()

    // In another goroutine
    go func() {
        // This will yield if high-priority is active
        yieldpoint.MaybeYield()

        // Or block until high-priority ends
        yieldpoint.WaitIfActive()
    }()
}
```

### Priority Levels

```go
// Set goroutine priority
yieldpoint.SetPriority(yieldpoint.PriorityHigh)

// Get current priority
priority := yieldpoint.GetPriority()
```

### Context Support

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

// Yield with context
err := yieldpoint.MaybeYieldWithContext(ctx)
if err != nil {
    // Handle context cancellation
}

// Wait with context
err = yieldpoint.WaitIfActiveWithContext(ctx)
if err != nil {
    // Handle context cancellation
}
```

### Tracing

```go
yieldpoint.SetTraceFunc(func(e yieldpoint.YieldEvent) {
    fmt.Printf("Goroutine %d yielded at %v (duration: %v)\n",
        e.GoroutineID, e.Timestamp, e.Duration)
})
```

## API Reference

### Core Functions

- `MaybeYield()`: Voluntarily yields if high-priority is active
- `EnterHighPriority()`: Begins a high-priority section
- `ExitHighPriority()`: Ends a high-priority section
- `WaitIfActive()`: Blocks until high-priority section ends
- `IsHighPriorityActive()`: Checks if any high-priority sections are active

### Priority Functions

- `SetHighPriority(high bool)`: Sets whether the current goroutine has high priority
- `GetHighPriority() bool`: Gets whether the current goroutine has high priority

### Context-aware Functions

- `MaybeYieldWithContext(ctx context.Context) error`
- `WaitIfActiveWithContext(ctx context.Context) error`

### Tracing

- `SetTraceFunc(func(YieldEvent))`: Sets a callback for yield events

## Status

This package is experimental. Feedback and contributions are welcome!

## License

MIT License - see LICENSE file for details
