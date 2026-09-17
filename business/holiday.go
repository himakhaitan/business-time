package business

import "time"

// HolidayCalendar determines whether a given date is a holiday.
//
// A HolidayCalendar is used by BusinessCalendar to identify dates that are
// not considered business days because they are holidays.
//
// Implementations may derive holidays from any source or set of rules,
// including:
//
//   - fixed dates
//   - regional or national holidays
//   - company-specific holidays
//   - calculated holidays
//   - one-off holidays
//   - combinations of multiple holiday calendars
//
// BusinessCalendar is responsible for interpreting time.Time values using
// its configured timezone before applying holiday rules. Implementations
// therefore should treat t as already being in the relevant business
// calendar's location.
//
// Holiday rules should generally depend on the calendar date represented by
// t, rather than its time of day. For example, a holiday should normally
// apply to the entire local calendar day regardless of whether t represents
// midnight, noon, or another time during that day.
//
// A nil HolidayCalendar may be supplied to Config when no holidays need to
// be observed.
type HolidayCalendar interface {
	// IsHoliday reports whether the calendar date represented by t is a
	// holiday.
	//
	// The date is interpreted according to the timezone of the enclosing
	// BusinessCalendar. Implementations should generally use the local
	// calendar date represented by t rather than treating t as an absolute
	// instant.
	//
	// The result should normally be independent of the time of day. If t
	// represents any time within a holiday's local calendar day, IsHoliday
	// should return the same result.
	IsHoliday(t time.Time) bool
}
