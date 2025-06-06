# yieldpoint Examples

This directory contains various examples demonstrating different features of the yieldpoint package.

## Examples

1. **Basic Usage** (`basic_usage/`)

   - Demonstrates basic high-priority task handling
   - Shows how background tasks yield to high-priority tasks
   - Run with: `go run examples/basic_usage/main.go`

2. **Context Support** (`context_support/`)

   - Demonstrates using yieldpoint with context for timeout and cancellation
   - Shows how to handle context cancellation in yielding operations
   - Run with: `go run examples/context_support/main.go`

3. **Nested Priority** (`nested_priority/`)

   - Demonstrates nested high-priority sections
   - Shows how multiple priority levels can be active simultaneously
   - Run with: `go run examples/nested_priority/main.go`

4. **Concurrent Yielding** (`concurrent_yielding/`)

   - Demonstrates how multiple goroutines can yield simultaneously
   - Shows coordination between concurrent yielding operations
   - Run with: `go run examples/concurrent_yielding/main.go`

5. **Fast** (`fast/`)

   - Demonstrates high-performance yielding patterns
   - Shows optimized usage of yieldpoint for latency-sensitive applications
   - Run with: `go run examples/fast/main.go`

6. **Wait If Active** (`wait_if_active/`)

   - Demonstrates the use of WaitIfActive functionality
   - Shows how to wait for high-priority tasks to complete
   - Compares WaitIfActive with MaybeYield behavior
   - Run with: `go run examples/wait_if_active/main.go`

## Running the Examples

Each example is in its own directory and can be run independently. To run an example:

```bash
go run examples/<example_name>/main.go
```

For example, to run the basic usage example:

```bash
go run examples/basic_usage/main.go
```

## Notes

- Each example is self-contained and demonstrates a specific feature or use case
- The examples use different sleep durations to simulate work and make the output more readable
- You may need to adjust the sleep durations based on your system's performance
- All examples use the local yieldpoint package, so make sure you're running them from the root directory of the project
