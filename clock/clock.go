// Package clock provides abstractions for obtaining the current time.
//
// The package decouples time-dependent application logic from the system
// wall clock. This allows components to depend on a small, explicit time
// interface rather than calling time.Now directly.
//
// A Clock represents a source of the current time. Production code can use
// the system clock, while tests can provide deterministic implementations
// that return a known point in time.
//
// All times returned by the Clock interface are represented in UTC.
// Applications that need calendar-specific behavior should interpret these
// instants using an explicit time zone.
package clock

import "time"

// Clock represents a source of the current time.
//
// Clock is intentionally small so that it can be easily implemented,
// substituted, and injected into components that depend on the current
// time.
//
// Implementations should return the current time as an absolute instant.
// The returned time should use UTC to provide a consistent representation
// across applications, services, and infrastructure.
type Clock interface {
	// Now returns the current time as a UTC instant.
	Now() time.Time
}
