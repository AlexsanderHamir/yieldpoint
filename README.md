# yieldpoint

A Go package that enables cooperative goroutine yielding based on priority-aware scheduling.

## Overview

`yieldpoint` provides a simple yet powerful mechanism for implementing cooperative multitasking in Go applications. It allows goroutines to voluntarily yield execution when high-priority tasks are active.

## Features

- **Priority-based Yielding**: Goroutines can yield execution when high-priority tasks are active
- **Efficient Blocking**: Uses `sync.Cond` for efficient blocking without busy waiting
- **Context Support**: Context-aware variants of all operations
- **High Priority Support**: Simple boolean flag for high-priority tasks
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

        // higher performance
        yieldpoint.MaybeYieldFast()

        // higher performance
        yieldpoint.WaitIfActiveFast()
    }()
}
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

## API Reference

### Core Functions

- `MaybeYield()`: Voluntarily yields if high-priority is active
- `EnterHighPriority()`: Begins a high-priority section
- `ExitHighPriority()`: Ends a high-priority section
- `WaitIfActive()`: Blocks until high-priority section ends
- `IsHighPriorityActive()`: Checks if any high-priority sections are active

### Configuration Functions

- `SetSpinWaitIterations(iterations int)`: Sets the number of spin iterations before falling back to blocking wait. This setting applies to the fast variant `WaitIfActiveFast`.
- `SetDefaultYieldDuration(duration time.Duration)`: Sets the default duration for yielding operations. This setting applies to the standard `MaybeYield` to ensure it yields.

## High Performance Functions

- `MaybeYieldFast()`: High performance version
- `WaitIfActiveFast()`: High performance version


### Context-aware Functions

- `MaybeYieldWithContext(ctx context.Context) error`
- `WaitIfActiveWithContext(ctx context.Context) error`

## Contributions

Share your talents and ideas !!

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
