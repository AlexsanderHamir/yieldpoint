# yieldpoint

**Status: Early Development**

`yieldpoint` is an experimental Go library for cooperative goroutine yielding and priority-aware scheduling. It lets developers mark *high-priority* code paths, prompting other goroutines to yield or pause. This reduces contention and improves responsiveness in latency-sensitive applications.

## What is it?

Go’s scheduler doesn’t expose fine-grained control over goroutine priority or preemption. That can make it hard to ensure critical sections run smoothly.

`yieldpoint` fills that gap by providing:

- A way to mark high-priority sections
- Cooperative yielding from lower-priority goroutines
- Simple APIs to manage priority-aware scheduling

## Why use it?

Scenarios where `yieldpoint` can help:

- 🎮 Game engines needing smooth frame updates  
- 🎧 Real-time audio/video pipelines  
- 🛠 Background jobs mixed with critical tasks  
- 🖱 Interactive apps that prioritize user input  

## Example

```go
// High-priority section
yieldpoint.EnterHighPriority()
// Critical work here
yieldpoint.ExitHighPriority()

// In background goroutines
yieldpoint.MaybeYield()


