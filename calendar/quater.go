package calendar

import "time"

// StartOfQuarter returns the first instant of the fiscal quarter containing
// value.
//
// The fiscal quarter is determined using the Calendar's configured
// FiscalYearStart. The Calendar's location determines the local date
// used to identify the quarter.
//
// For example, when FiscalYearStart is January:
//
//	January–March   = Q1
//	April–June      = Q2
//	July–September  = Q3
//	October–December = Q4
//
// When FiscalYearStart is April:
//
//	April–June      = Q1
//	July–September  = Q2
//	October–December = Q3
//	January–March   = Q4
func (c *Calendar) StartOfQuarter(value time.Time) time.Time {
	fiscalYearStart := c.StartOfYear(value)

	local := value.In(c.location)

	monthOffset := (int(local.Month()) - int(c.fiscalYearStart) + 12) % 12
	quarterOffset := monthOffset / 3

	return fiscalYearStart.AddDate(0, quarterOffset*3, 0)
}

// EndOfQuarter returns the final representable instant of the fiscal quarter
// containing value.
//
// For interval calculations, prefer:
//
//	[StartOfQuarter(value), StartOfNextQuarter(value))
func (c *Calendar) EndOfQuarter(value time.Time) time.Time {
	return c.StartOfNextQuarter(value).Add(-time.Nanosecond)
}

// StartOfNextQuarter returns the first instant of the fiscal quarter
// immediately following the quarter containing value.
func (c *Calendar) StartOfNextQuarter(value time.Time) time.Time {
	return c.StartOfQuarter(value).AddDate(0, 3, 0)
}

// StartOfPreviousQuarter returns the first instant of the fiscal quarter
// immediately preceding the quarter containing value.
func (c *Calendar) StartOfPreviousQuarter(value time.Time) time.Time {
	return c.StartOfQuarter(value).AddDate(0, -3, 0)
}

// IsSameQuarter reports whether a and b belong to the same fiscal quarter
// according to the Calendar's configured location and fiscal-year start.
func (c *Calendar) IsSameQuarter(a, b time.Time) bool {
	return c.StartOfQuarter(a).Equal(c.StartOfQuarter(b))
}

// Quarter returns the fiscal quarter number containing value.
//
// The result is between 1 and 4. The quarter is determined using the
// Calendar's configured FiscalYearStart and Location.
func (c *Calendar) Quarter(value time.Time) int {
	local := value.In(c.location)

	monthOffset := (int(local.Month()) - int(c.fiscalYearStart) + 12) % 12

	return monthOffset/3 + 1
}

// AddQuarters returns value shifted by the specified number of calendar
// quarters.
//
// A quarter is treated as three calendar months rather than a fixed duration.
// The calculation follows time.Time.AddDate semantics.
//
// AddQuarters is independent of Calendar configuration and therefore does not
// apply fiscal-year or timezone-specific quarter semantics.
func AddQuarters(value time.Time, quarters int) time.Time {
	return value.AddDate(0, quarters*3, 0)
}
