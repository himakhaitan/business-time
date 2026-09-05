package clock

import (
	"sync"
	"time"
)

// Controlled is a manually controlled Clock implementation intended for
// deterministic tests.
//
// Unlike Static, which always returns the same instant, Controlled allows
// tests to change the current time explicitly using Set or to move it
// forward or backward using Advance.
//
// Controlled is safe for concurrent use. Calls to Now can safely run
// concurrently with Set and Advance.
type Controlled struct {
	mu      sync.RWMutex
	current time.Time
}

// NewControlled returns a new Controlled clock initialized to t.
//
// The provided time is normalized to UTC so that Controlled follows the
// Clock contract of returning times in a consistent UTC representation.
func NewControlled(t time.Time) *Controlled {
	return &Controlled{
		current: t.UTC(),
	}
}

// Now returns the current time of the controlled clock in UTC.
//
// Unlike System.Now, this value only changes when Set or Advance is called.
func (c *Controlled) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.current
}

// Set changes the current time of the controlled clock.
//
// The provided time is normalized to UTC. Subsequent calls to Now return
// the new time until the clock is changed again using Set or Advance.
func (c *Controlled) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.current = t.UTC()
}

// Advance moves the current time of the controlled clock by d.
//
// Both positive and negative durations are supported. A positive duration
// moves the clock forward, while a negative duration moves it backward.
//
// The operation changes the clock by an elapsed duration and does not depend
// on calendar or time-zone rules.
func (c *Controlled) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.current = c.current.Add(d)
}
