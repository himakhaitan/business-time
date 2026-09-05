package calendar

import (
	"time"
)

// StartOfDay returns the first instant of the calendar day containing value.
//
// The day is determined using the Calendar's configured location rather than
// the location carried by value. This makes the result consistent with the
// calendar's definition of a day, regardless of where value originated.
//
// The returned time is midnight at the beginning of that local calendar day.
// Calendar boundaries are constructed in the configured location so that
// daylight-saving transitions and other timezone rules are respected.
//
// For example, if a Calendar is configured for America/New_York and value
// represents an instant that falls on March 10 in New York, StartOfDay returns
// midnight on March 10 in New York, even if value itself is represented in UTC.
//
// StartOfDay is useful when building day-based queries, grouping events by
// local date, or defining the beginning of a calendar-day interval.
func (c *Calendar) StartOfDay(value time.Time) time.Time {
	local := value.In(c.location)

	return time.Date(
		local.Year(),
		local.Month(),
		local.Day(),
		0, 0, 0, 0,
		c.location,
	)

}

// EndOfDay returns the final representable instant of the calendar day
// containing value.
//
// The day is determined using the Calendar's configured location, and the
// returned time is expressed in that same location. The result represents
// an inclusive, closed boundary: the returned instant is still part of the
// calendar day.
//
// The calculation is based on the beginning of the following calendar day
// rather than assuming that a calendar day is exactly 24 hours long. This
// ensures that calendar boundaries remain correct across daylight-saving
// transitions and other timezone offset changes.
//
// When defining ranges or intervals, prefer a half-open interval using
// StartOfDay and StartOfNextDay:
//
//	[StartOfDay(value), StartOfNextDay(value))
//
// This convention includes the entire calendar day without requiring a
// "last instant" calculation and is generally safer for comparisons,
// database queries, and interval arithmetic.
func (c *Calendar) EndOfDay(value time.Time) time.Time {
	return c.StartOfNextDay(value).Add(-time.Nanosecond)
}

// StartOfNextDay returns the first instant of the calendar day immediately
// following the day containing value.
//
// The calculation is performed according to the Calendar's configured
// location and uses calendar-day arithmetic rather than adding a fixed
// 24-hour duration. As a result, it correctly handles days that are shorter
// or longer than 24 hours because of daylight-saving transitions.
//
// The returned time is midnight at the beginning of the next local calendar
// day and is expressed in the Calendar's configured location.
//
// This method is particularly useful for constructing half-open intervals,
// such as [StartOfDay(value), StartOfNextDay(value)), where the end boundary
// is exclusive.
func (c *Calendar) StartOfNextDay(value time.Time) time.Time {
	start := c.StartOfDay(value)
	return start.AddDate(0, 0, 1)
}

// StartOfPreviousDay returns the first instant of the calendar day immediately
// preceding the day containing value.
//
// The calculation is performed according to the Calendar's configured
// location and uses calendar-day arithmetic rather than subtracting a fixed
// 24-hour duration. This ensures that the result remains correct across
// daylight-saving transitions and other timezone offset changes.
//
// The returned time is midnight at the beginning of the previous local
// calendar day and is expressed in the Calendar's configured location.
func (c *Calendar) StartOfPreviousDay(value time.Time) time.Time {
	start := c.StartOfDay(value)
	return start.AddDate(0, 0, -1)
}

// IsSameDay reports whether a and b belong to the same calendar day according
// to the Calendar's configured location.
//
// The comparison is based on the local year, month, and day after converting
// both values into the Calendar's location. The absolute instants represented
// by a and b may therefore be different, and the locations carried by the
// input values do not affect the result.
//
// This is useful when comparing timestamps from different systems, services,
// or time zones while applying a single, well-defined business or application
// calendar.
//
// For example, two timestamps that fall on different UTC dates may still be
// considered the same day when viewed in a timezone where both occur on the
// same local calendar date.
func (c *Calendar) IsSameDay(a, b time.Time) bool {
	a = a.In(c.location)
	b = b.In(c.location)

	return a.Year() == b.Year() &&
		a.Month() == b.Month() &&
		a.Day() == b.Day()

}

// DayOfYear returns the ordinal position of the calendar day within its year,
// according to the Calendar's configured location.
//
// The result ranges from 1 through 365 in a common year and from 1 through
// 366 in a leap year. The calendar day is determined after converting value
// into the Calendar's configured location, so the result is based on the
// local date rather than the date represented by value's original location.
//
// For example, an instant near midnight may belong to one calendar day in UTC
// and another calendar day in the configured location. DayOfYear always uses
// the latter.
func (c *Calendar) DayOfYear(value time.Time) int {
	return value.In(c.location).YearDay()
}
