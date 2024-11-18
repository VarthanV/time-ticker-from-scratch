package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

type ticker struct {
	C         chan time.Time
	isRunning atomic.Uint32
}

func New(every time.Duration) *ticker {
	t := &ticker{}
	t.isRunning.Store(1)
	t.tick(every)
	return t
}

func (t *ticker) tick(every time.Duration) {

	ch := make(chan time.Time)

	go func() {
		nextTick := time.Now().Add(every)

		for {
			if t.isRunning.Load() == 0 {
				return
			}

			if time.Now().Equal(nextTick) ||
				time.Now().After(nextTick) {
				ch <- time.Now()
				nextTick = time.Now().Add(every)
			}
		}
	}()
	t.C = ch
}

func (t *ticker) Close() {
	if t.isRunning.Swap(0) == 1 {
		close(t.C)
	}
}

func main() {
	fmt.Println("hello world")

	tick := New(1 * time.Second)

	go func() {
		time.Sleep(10 * time.Second)
		tick.Close()
	}()

	for val := range tick.C {
		log.Println("val is ", val.String())
	}
	log.Println("exited ticker")
}
