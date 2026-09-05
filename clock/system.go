package clock

import "time"

// System is a Clock implementation backed by the system wall clock.
//
// It provides the current time from the host operating system and normalizes
// the result to UTC. System is intended for production use where the
// application needs the actual current time.
type System struct{}

// NewSystem returns a production Clock backed by the system wall clock.
//
// The returned Clock provides the current time in UTC.
func NewSystem() Clock {
	return System{}
}

// Now returns the current system time as a UTC instant.
//
// The value is obtained from the system wall clock at the time of the call.
// Successive calls may return different values as time advances.
func (System) Now() time.Time {
	return time.Now().UTC()
}
