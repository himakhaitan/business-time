package calendar

import "time"

// StartOfYear returns the first instant of the fiscal year containing value.
//
// The fiscal year begins in the Calendar's configured
// FiscalYearStart. The Calendar's configured location determines the
// local date used to identify the fiscal year.
//
// When FiscalYearStart is January, this corresponds to the ordinary
// calendar year.
func (c *Calendar) StartOfYear(value time.Time) time.Time {
	local := value.In(c.location)

	year := local.Year()

	if local.Month() < c.fiscalYearStart {
		year--
	}

	return time.Date(
		year,
		c.fiscalYearStart,
		1,
		0, 0, 0, 0,
		c.location,
	)
}

// EndOfYear returns the final representable instant of the fiscal year
// containing value.
//
// For interval calculations, prefer:
//
//	[StartOfYear(value), StartOfNextYear(value))
func (c *Calendar) EndOfYear(value time.Time) time.Time {
	return c.StartOfNextYear(value).Add(-time.Nanosecond)
}

// StartOfNextYear returns the first instant of the fiscal year immediately
// following the fiscal year containing value.
func (c *Calendar) StartOfNextYear(value time.Time) time.Time {
	return c.StartOfYear(value).AddDate(1, 0, 0)
}

// StartOfPreviousYear returns the first instant of the fiscal year immediately
// preceding the fiscal year containing value.
func (c *Calendar) StartOfPreviousYear(value time.Time) time.Time {
	return c.StartOfYear(value).AddDate(-1, 0, 0)
}

// IsSameYear reports whether a and b belong to the same fiscal year according
// to the Calendar's configured location and fiscal-year start.
func (c *Calendar) IsSameYear(a, b time.Time) bool {
	return c.StartOfYear(a).Equal(c.StartOfYear(b))
}

// AddYears returns value shifted by the specified number of calendar years.
//
// Calendar-year arithmetic follows time.Time's AddDate semantics.
func AddYears(value time.Time, years int) time.Time {
	return value.AddDate(years, 0, 0)
}

// IsLeapYear reports whether year is a leap year in the Gregorian calendar.
func IsLeapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%100 != 0)
}

// DaysInYear returns the number of days in the Gregorian calendar year.
//
// The result is 365 for a common year and 366 for a leap year.
func DaysInYear(year int) int {
	if IsLeapYear(year) {
		return 366
	}

	return 365
}
