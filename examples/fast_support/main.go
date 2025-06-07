package main

import (
	"testing"
	"time"

	"github.com/AlexsanderHamir/yieldpoint"
)

func TestWaitIfActive(t *testing.T) {
	yieldpoint.EnterHighPriority()

	done := make(chan struct{})
	go func() {
		yieldpoint.WaitIfActiveFast()
		close(done)
	}()

	time.Sleep(10 * time.Millisecond)

	yieldpoint.ExitHighPriority()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("WaitIfActiveFast timed out")
	}
}

