# Custom Ticker Implementation in Go

This program demonstrates a custom ticker implementation in Go. It is similar to Go's built-in `time.Ticker` but provides a greater degree of control over its lifecycle. The ticker sends timestamps to a channel at regular intervals until it is stopped. Below is an explanation of the components and their usage.

---

### **Struct Definition**
```go
type ticker struct {
	C         chan time.Time
	isRunning atomic.Uint32
}
```
- **`ticker` struct**:
  - **`C`**: A channel that emits the current time at specified intervals.
  - **`isRunning`**: An atomic flag (0 or 1) to indicate whether the ticker is running.

---

### **Ticker Initialization**
```go
func New(every time.Duration) *ticker {
	t := &ticker{}
	t.isRunning.Store(1)
	t.tick(every)
	return t
}
```
- **`New` Function**:
  - Accepts a `time.Duration` for the interval between ticks.
  - Initializes a `ticker` instance and sets the `isRunning` flag to 1 (active).
  - Starts the tick generation by invoking the `tick` method.

---

### **Tick Generation**
```go
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
```
- **`tick` Method**:
  - Creates a new channel (`ch`) to emit timestamps.
  - Spawns a goroutine that continuously:
    - Checks if the ticker is running (`isRunning` flag).
    - Sends the current time to the channel at intervals defined by `every`.
    - Updates the `nextTick` to maintain the interval.
  - Assigns the channel to the `C` field of the ticker.

---

### **Stopping the Ticker**
```go
func (t *ticker) Close() {
	if t.isRunning.Swap(0) == 1 {
		close(t.C)
	}
}
```
- **`Close` Method**:
  - Stops the ticker by setting `isRunning` to 0.
  - Closes the channel to signal that no more ticks will be sent.

---

### **6. Example Function**
```go
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
```

## **Example Output**
```plaintext
hello world
2024/11/17 12:00:01 val is  2024-11-17 12:00:01.123456789 +0000 UTC
2024/11/17 12:00:02 val is  2024-11-17 12:00:02.123456789 +0000 UTC
...
2024/11/17 12:00:10 val is  2024-11-17 12:00:10.123456789 +0000 UTC
2024/11/17 12:00:10 exited ticker
```

---