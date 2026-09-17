// Package calendar provides timezone-aware calendar operations for Go.
//
// The calendar package builds on Go's standard time package and provides
// calendar semantics for days, weeks, months, quarters, and years.
//
// Unlike duration-based arithmetic, calendar operations account for the
// configured timezone and calendar boundaries. This is important when
// working across daylight-saving transitions, leap years, fiscal years,
// and other calendar boundaries.
//
// A Calendar is configured with a location, week start, and fiscal year
// start. These settings are used consistently by calendar operations.
//
// Example:
//
//	cal, err := calendar.New(calendar.Config{
//	   Location: time.UTC,
//	   WeekStart: time.Monday,
//	   FiscalYearStart: time.January,
//	})
//	if err != nil {
//	   return err
//	}
//	start := cal.StartOfMonth(time.Now())
//	end := cal.StartOfNextMonth(time.Now())
//
//	fmt.Println(start)
//	fmt.Println(end)
package calendar

import (
	"fmt"
	"time"
)

// Config defines the calendar context used by Calendar.
//
// A calendar is not just a collection of date calculations. Operations such as
// the start of a day, the start of a week, and the quarter containing a date
// depend on the calendar context in which the date is interpreted.
//
// Location determines the timezone used for calendar boundaries.
// WeekStart determines which weekday is considered the first day of a week.
// FiscalYearStart determines which month begins the fiscal year.
//
// All fields are required. business-time intentionally does not provide
// implicit defaults because timezone and calendar conventions are important
// business decisions that should be explicit in production applications.
type Config struct {
	// Location specifies the timezone in which calendar calculations are
	// performed.
	//
	// Calendar boundaries such as the start of a day, week, month, quarter,
	// or year are interpreted in this location. The location also determines
	// how daylight-saving-time transitions affect calendar arithmetic.
	//
	// Location must not be nil.
	Location *time.Location

	// WeekStart specifies the weekday on which a calendar week begins.
	//
	// For example, use time.Monday for calendars where Monday is the first
	// day of the week, or time.Sunday for calendars where Sunday is the first
	// day of the week.
	//
	// WeekStart must be a valid time.Weekday value.
	WeekStart time.Weekday

	// FiscalYearStart specifies the month in which the fiscal year begins.
	//
	// For a calendar year that follows January through December, use
	// time.January. For example, a fiscal year beginning in April should use
	// time.April.
	//
	// FiscalYearStart must be a valid time.Month value.
	FiscalYearStart time.Month
}

// Calendar represents a configured calendar context.
//
// Calendar provides calendar-aware operations such as determining the start
// of a day, week, month, quarter, or year. The configured location, week
// start, and fiscal-year start are fixed when the Calendar is created.
//
// A Calendar is safe to share between goroutines because its configuration
// cannot be modified after construction.
type Calendar struct {
	location        *time.Location
	weekStart       time.Weekday
	fiscalYearStart time.Month
}

// New creates a Calendar from the supplied configuration.
//
// New validates the calendar configuration before creating the Calendar.
// Configuration is immutable after construction, so a Calendar can be safely
// reused and shared throughout an application.
//
// The caller must explicitly provide all calendar configuration. In
// particular, business-time does not silently fall back to time.Local because
// timezone selection can materially change the result of calendar
// calculations.
//
// New returns an error when the configuration is incomplete or contains an
// invalid week start or fiscal-year start.
func New(cfg Config) (*Calendar, error) {
	if cfg.Location == nil {
		return nil, fmt.Errorf("calendar: location is required")
	}

	if cfg.WeekStart < time.Sunday || cfg.WeekStart > time.Saturday {
		return nil, fmt.Errorf("calendar: invalid week start: %v", cfg.WeekStart)
	}

	if cfg.FiscalYearStart < time.January || cfg.FiscalYearStart > time.December {
		return nil, fmt.Errorf(
			"calendar: invalid fiscal year start: %v",
			cfg.FiscalYearStart,
		)
	}

	return &Calendar{
		location:        cfg.Location,
		weekStart:       cfg.WeekStart,
		fiscalYearStart: cfg.FiscalYearStart,
	}, nil
}

// Location returns the timezone used by the Calendar.
//
// Calendar operations that depend on local calendar boundaries, such as
// StartOfDay and StartOfWeek, use this location to interpret the supplied
// time.Time value.
func (c *Calendar) Location() *time.Location {
	return c.location
}

// WeekStart returns the weekday configured as the first day of the week.
func (c *Calendar) WeekStart() time.Weekday {
	return c.weekStart
}

// FiscalYearStart returns the month configured as the first month of the
// fiscal year.
func (c *Calendar) FiscalYearStart() time.Month {
	return c.fiscalYearStart
}
