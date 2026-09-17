package calendar

import (
	"time"
)

// StartOfMonth returns the first instant of the calendar month containing
// value.
//
// The month is determined using the Calendar's configured location rather
// than the location carried by value. This ensures that the result is
// consistent regardless of where the input timestamp originated.
//
// The returned time is midnight on the first day of the month, expressed in
// the Calendar's configured location.
//
// The boundary is constructed as a calendar date rather than by truncating a
// fixed duration. This is important because time.Duration-based truncation
// does not represent calendar boundaries correctly across timezone changes.
//
// For example, if a Calendar is configured for America/New_York and value
// represents an instant that falls on March 15 in New York, StartOfMonth
// returns midnight on March 1 in New York.
func (c *Calendar) StartOfMonth(value time.Time) time.Time {
	local := value.In(c.location)

	return time.Date(
		local.Year(),
		local.Month(),
		1,
		0, 0, 0, 0,
		c.location,
	)
}

// EndOfMonth returns the final representable instant of the calendar month
// containing value.
//
// The month is determined using the Calendar's configured location, and the
// returned time is expressed in that same location.
//
// The calculation is based on the beginning of the following calendar month
// rather than assuming that a month has a fixed number of days. This correctly
// handles months with different lengths, leap years, daylight-saving
// transitions, and other timezone offset changes.
//
// The result represents an inclusive, closed boundary: the final nanosecond
// belonging to the calendar month.
//
// When defining ranges or intervals, prefer a half-open interval using
// StartOfMonth and StartOfNextMonth:
//
//	[StartOfMonth(value), StartOfNextMonth(value))
//
// This convention avoids relying on a "last instant" calculation and is
// generally safer for comparisons, database queries, and interval arithmetic.
func (c *Calendar) EndOfMonth(value time.Time) time.Time {
	return c.StartOfNextMonth(value).Add(-time.Nanosecond)
}

// StartOfNextMonth returns the first instant of the calendar month immediately
// following the month containing value.
//
// The calculation is performed according to the Calendar's configured
// location and uses calendar-month arithmetic rather than adding a fixed
// duration.
//
// The returned time is midnight on the first day of the next calendar month,
// expressed in the Calendar's configured location.
//
// This method is useful for constructing half-open monthly intervals:
//
//	[StartOfMonth(value), StartOfNextMonth(value))
//
// Using the beginning of the next month as the exclusive boundary avoids
// assumptions about the number of days or hours in the current month.
func (c *Calendar) StartOfNextMonth(value time.Time) time.Time {
	start := c.StartOfMonth(value)
	return start.AddDate(0, 1, 0)
}

// StartOfPreviousMonth returns the first instant of the calendar month
// immediately preceding the month containing value.
//
// The calculation is performed according to the Calendar's configured
// location and uses calendar-month arithmetic rather than subtracting a
// fixed duration.
//
// The returned time is midnight on the first day of the previous calendar
// month, expressed in the Calendar's configured location.
//
// Calendar arithmetic ensures that month boundaries remain correct across
// different month lengths, leap years, daylight-saving transitions, and
// year boundaries.
func (c *Calendar) StartOfPreviousMonth(value time.Time) time.Time {
	start := c.StartOfMonth(value)
	return start.AddDate(0, -1, 0)
}

// IsSameMonth reports whether a and b belong to the same calendar month
// according to the Calendar's configured location.
//
// The comparison is based on the local year and month after converting both
// values into the Calendar's location. The locations carried by a and b do
// not affect the result.
//
// The year is part of the comparison because January 2025 and January 2026
// are different calendar months.
//
// This method is useful when grouping timestamps by reporting month,
// determining whether two events belong to the same monthly period, or
// comparing timestamps received from systems operating in different
// time zones.
func (c *Calendar) IsSameMonth(a, b time.Time) bool {
	a = a.In(c.location)
	b = b.In(c.location)

	return a.Year() == b.Year() &&
		a.Month() == b.Month()
}

// DaysInMonth returns the number of calendar days in the month containing
// value according to the Calendar's configured location.
//
// The result accounts for the different lengths of calendar months and leap
// years. It ranges from 28 to 31.
//
// The month is determined after converting value into the Calendar's
// configured location, so an instant near a month boundary may produce a
// different result depending on the Calendar's timezone.
//
// The calculation uses the first day of the following month. The day of the
// month of that boundary is exactly the number of days in the current month.
// This avoids measuring elapsed time, which can be incorrect around
// daylight-saving transitions where a calendar day may not be exactly
// 24 hours.
func (c *Calendar) DaysInMonth(value time.Time) int {
	next := c.StartOfNextMonth(value)

	return next.AddDate(0, 0, -1).Day()
}

// AddMonths returns value shifted by the specified number of calendar months.
//
// A calendar month is not treated as a fixed duration because months contain
// different numbers of days. The calculation follows time.Time.AddDate
// semantics.
//
// AddMonths is independent of Calendar configuration. It does not apply
// timezone-specific calendar boundaries or fiscal-month semantics.
func AddMonths(value time.Time, months int) time.Time {
	return value.AddDate(0, months, 0)
}
