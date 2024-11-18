package main

import (
	"sync"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	duration := 500 * time.Millisecond
	tick := New(duration)

	var ticks []time.Time
	var mu sync.Mutex
	done := make(chan struct{})

	go func() {
		for tickTime := range tick.C {
			mu.Lock()
			ticks = append(ticks, tickTime)
			mu.Unlock()
		}
		close(done)
	}()

	time.Sleep(2 * time.Second)

	tick.Close()

	<-done

	mu.Lock()
	defer mu.Unlock()

	// Ensure we received at least 3 ticks (500ms * 4 intervals = 2 seconds)
	if len(ticks) < 3 {
		t.Errorf("expected at least 3 ticks, got %d", len(ticks))
	}

	// Ensure the ticks are spaced approximately `duration` apart
	for i := 1; i < len(ticks); i++ {
		diff := ticks[i].Sub(ticks[i-1])
		if diff < duration-50*time.Millisecond || diff > duration+50*time.Millisecond {
			t.Errorf("ticks not spaced correctly: got %v, expected ~%v", diff, duration)
		}
	}
}

func TestTickerClose(t *testing.T) {
	// Create a new ticker with a short duration
	tick := New(100 * time.Millisecond)

	// Close the ticker immediately
	tick.Close()

	// Wait briefly to ensure closure
	time.Sleep(200 * time.Millisecond)

	// Attempt to read from the channel (should not receive any values)
	select {
	case _, ok := <-tick.C:
		if ok {
			t.Errorf("expected ticker channel to be closed")
		}
	default:
		// Channel is empty and closed, pass the test
	}
}
