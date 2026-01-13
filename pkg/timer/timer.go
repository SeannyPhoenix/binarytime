package timer

import (
	"log"
	"time"
)

type Timer struct {
	running  bool
	elapsed  time.Duration
	duration time.Duration
	last     time.Time
}

// Set clears the timer and sets the
// duration to the given value, even if
// it has already started
// Typical use:
// var t Timer
// t.Set(10 * time.Second)
func (t *Timer) Set(d time.Duration) {
	*t = Timer{duration: d}
}

// Reset clears the timer and keeps the
// duration. The current timer is cleared,
// even if it has already started
func (t *Timer) Reset() {
	*t = Timer{duration: t.duration}
}

// Start begins or resumes the timer.
// If the timer is already running or
// has finished, start is a no-op
func (t *Timer) Start(now time.Time) {
	if !t.running && !t.Finished() {
		t.running = true
		t.Tick(now)
	}
}

// Stop pauses the timer. If the timer
// is not running, Stop is a no-op
func (t *Timer) Stop(now time.Time) {
	if t.running {
		t.Tick(now)
		t.running = false
	}
}

// Toggle starts the timer if it is not running,
// and stops the timer if it is. If the timer
// is finished, Toggle is a no-op
func (t *Timer) Toggle(now time.Time) {
	if t.running {
		t.Stop(now)
	} else {
		t.Start(now)
	}
}

// Tick increments the timer if it is running.
// If the time since the last tick would
// exceed the timer's duration, it marks the
// timer complete and stops it. If the timer
// is not running, which includes a timer that
// has finished, Tick is a no-op
func (t *Timer) Tick(now time.Time) {
	if t.running {
		if t.last.IsZero() {
			log.Println(now.Sub(t.last))
		}
		delta := now.Sub(t.last)
		t.elapsed = min(t.duration, t.elapsed+delta)
		t.last = now
		if t.elapsed >= t.duration {
			t.running = false
		}
	}
}

// Duration returns the diration set for the timer
func (t *Timer) Duration() time.Duration {
	return t.duration
}

// Elapsed returns the duration the timer
// has already run without updating it
func (t *Timer) Elapsed() time.Duration {
	return t.elapsed
}

// Progress returns the percent elapsed
// as a float32 between 0 and 1 without
// updating the timer
func (t *Timer) Progress() float32 {
	if t.duration == 0 {
		return 0
	}
	return float32(t.elapsed) / float32(t.duration)
}

// Running returns the running state of the timer
func (t *Timer) Running() bool {
	return t.running
}

// Started returns if the timer has been started.
func (t *Timer) Started() bool {
	return !t.last.IsZero()
}

// Finished returns if the timer has completed.
func (t *Timer) Finished() bool {
	return t.elapsed != 0 && t.elapsed == t.duration
}
