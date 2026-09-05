package clock

import "time"

// Static is a Clock that always returns the same point in time.
//
// It is primarily useful for tests and other deterministic scenarios where
// application behavior needs to be evaluated against a known point in time.
// The configured time is normalized to UTC when the clock is created.
//
// Static does not advance with the passage of time. Every call to Now returns
// the same instant.
type Static struct {
	current time.Time
}

// NewStatic returns a Clock that always returns t.
//
// The provided time is normalized to UTC so that Static follows the Clock
// contract of returning times in a consistent UTC representation.
//
// NewStatic is useful for making time-dependent tests deterministic by
// replacing the system clock with a known point in time.
func NewStatic(t time.Time) Clock {
	return Static{
		current: t.UTC(),
	}
}

// Now returns the configured static time in UTC.
//
// Every call returns the same point in time that was provided to NewStatic.
func (s Static) Now() time.Time {
	return s.current
}
