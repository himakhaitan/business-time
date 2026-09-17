package calendar

import (
	"time"
)

// Weekday returns the weekday of value according to the Calendar's configured
// location.
//
// The location carried by value is ignored for the purpose of determining the
// weekday. Instead, value is first converted to the Calendar's location and
// the weekday is calculated from that local date.
//
// This matters when the same instant falls on different calendar dates in
// different time zones.
//
// Weekday itself is intrinsic to a date, while the Calendar's location
// determines which local date the instant belongs to.
//
// For example, an instant that is Sunday in UTC may already be Monday in
// Asia/Kolkata. A Calendar configured for Asia/Kolkata will therefore return
// Monday.
func (c *Calendar) Weekday(value time.Time) time.Weekday {
	return value.In(c.location).Weekday()
}

// StartOfWeek returns the first instant of the calendar week containing value.
//
// The beginning of the week is determined by the Calendar's configured
// WeekStart. For example, a Calendar configured with Monday as WeekStart will
// consider Monday the beginning of every calendar week.
//
// The day containing value is determined using the Calendar's configured
// location rather than the location carried by value.
//
// The returned time is midnight at the beginning of the configured week,
// expressed in the Calendar's configured location.
//
// Week boundaries are constructed using calendar-day arithmetic rather than
// fixed durations so they remain correct across daylight-saving transitions
// and other timezone offset changes.
//
// For example, if WeekStart is Monday and value falls on Wednesday,
// StartOfWeek returns the preceding Monday at local midnight.
//
// StartOfWeek is useful when grouping records by week, constructing weekly
// reporting ranges, or defining week-based intervals.
func (c *Calendar) StartOfWeek(value time.Time) time.Time {
	local := value.In(c.location)

	// Calculate how many days we need to move backwards from the current
	// weekday to reach the configured start of the week.
	offset := (int(local.Weekday()) - int(c.weekStart) + 7) % 7

	startOfDay := c.StartOfDay(local)

	return startOfDay.AddDate(0, 0, -offset)
}

// EndOfWeek returns the final representable instant of the calendar week
// containing value.
//
// The week is determined using the Calendar's configured location and
// WeekStart. The returned time is expressed in the Calendar's configured
// location.
//
// The calculation is based on the beginning of the following calendar week
// rather than assuming that a week is exactly seven 24-hour periods. This
// ensures that the calendar boundary remains correct across daylight-saving
// transitions and other timezone offset changes.
//
// The result represents an inclusive, closed boundary: the final nanosecond
// belonging to the calendar week.
//
// When defining ranges or intervals, prefer a half-open interval using
// StartOfWeek and StartOfNextWeek:
//
//	[StartOfWeek(value), StartOfNextWeek(value))
//
// This convention includes the entire week without requiring a "last instant"
// calculation and is generally safer for comparisons, database queries, and
// interval arithmetic.
func (c *Calendar) EndOfWeek(value time.Time) time.Time {
	return c.StartOfNextWeek(value).Add(-time.Nanosecond)
}

// StartOfNextWeek returns the first instant of the calendar week immediately
// following the week containing value.
//
// The calculation respects the Calendar's configured WeekStart and location.
//
// Calendar arithmetic is used instead of adding a fixed seven-day duration.
// This is important because a week containing a daylight-saving transition
// may contain fewer or more than 168 elapsed hours.
//
// The returned time is midnight at the beginning of the next calendar week
// and is expressed in the Calendar's configured location.
//
// This method is particularly useful for constructing half-open intervals:
//
//	[StartOfWeek(value), StartOfNextWeek(value))
func (c *Calendar) StartOfNextWeek(value time.Time) time.Time {
	start := c.StartOfWeek(value)
	return start.AddDate(0, 0, 7)
}

// StartOfPreviousWeek returns the first instant of the calendar week
// immediately preceding the week containing value.
//
// The calculation respects the Calendar's configured WeekStart and location.
//
// Calendar arithmetic is used instead of subtracting a fixed seven-day
// duration so the result remains correct across daylight-saving transitions
// and other timezone offset changes.
//
// The returned time is midnight at the beginning of the previous calendar week
// and is expressed in the Calendar's configured location.
func (c *Calendar) StartOfPreviousWeek(value time.Time) time.Time {
	start := c.StartOfWeek(value)
	return start.AddDate(0, 0, -7)
}

// IsSameWeek reports whether a and b belong to the same calendar week
// according to the Calendar's configured location and WeekStart.
//
// The comparison is based on the local calendar week containing each instant.
// The locations carried by a and b do not affect the result because both
// values are interpreted using the Calendar's configured location.
//
// This is intentionally different from comparing elapsed durations. Two
// timestamps can be separated by several days while still belonging to the
// same calendar week, and two timestamps separated by only a short duration
// can belong to different calendar weeks when they fall across a week
// boundary.
//
// For example, a Sunday evening and the following Monday morning may belong
// to different weeks when WeekStart is Monday, even if only a few hours apart.
func (c *Calendar) IsSameWeek(a, b time.Time) bool {
	return c.StartOfWeek(a).Equal(c.StartOfWeek(b))
}

// AddWeeks returns value shifted by the specified number of calendar weeks.
//
// A calendar week is treated as seven calendar days rather than a fixed
// duration of 168 hours. This distinction is important around daylight-saving
// transitions, where a local week may contain fewer or more than 168 elapsed
// hours.
//
// The returned time preserves the location and wall-clock time semantics of
// value as defined by time.Time's AddDate behavior.
//
// A positive number moves value forward; a negative number moves it backward.
// Zero returns a value representing the same instant.
//
// AddWeeks is independent of Calendar configuration and is therefore provided
// as a package-level function rather than a Calendar method.
func AddWeeks(value time.Time, weeks int) time.Time {
	return value.AddDate(0, 0, weeks*7)
}
