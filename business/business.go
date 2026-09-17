// Package business provides business-day calendar functionality built on top
// of Go's standard time package and the calendar package.
//
// A BusinessCalendar combines a base calendar with configurable weekend and
// holiday rules. It provides the context required for determining business
// days and performing business-day calculations.
//
// Weekend and holiday rules are extensible. The package provides common
// implementations for typical business calendars, while applications can
// provide their own rules for organization-specific working weeks, alternate
// Saturdays, regional calendars, or other requirements.
//
// BusinessCalendar uses calendar.Calendar for fundamental calendar semantics,
// including timezone, week boundaries, and fiscal-year configuration.
//
// Example:
//
//	cal, err := calendar.New(calendar.Config{
//		Location:        time.UTC,
//		WeekStart:       time.Monday,
//		FiscalYearStart: time.January,
//	})
//	if err != nil {
//		return err
//	}
//
//	businessCal, err := business.New(business.Config{
//		Calendar: cal,
//		Weekend:  business.StandardWeekend(),
//	})
//	if err != nil {
//		return err
//	}
//
//	_ = businessCal
package business

import (
	"errors"

	"github.com/himakhaitan/business-time/calendar"
)

var (
	// ErrNilCalendar indicates that a nil calendar was provided when
	// constructing a BusinessCalendar.
	//
	// A BusinessCalendar requires a calendar.Calendar to provide the
	// underlying calendar semantics, including location and calendar
	// boundaries.
	ErrNilCalendar = errors.New("business: calendar is required")

	// ErrNilWeekend indicates that no weekend rule was provided when
	// constructing a BusinessCalendar.
	//
	// A WeekendRule is required because the definition of a business day
	// depends on knowing which days are considered weekends.
	ErrNilWeekend = errors.New("business: weekend rule is required")
)

// Config configures a BusinessCalendar.
//
// Calendar is required and provides the underlying calendar semantics.
// Weekend controls which dates are considered weekends. If Weekend is nil,
// StandardWeekend is used.
//
// Holidays is optional. A nil HolidayCalendar means that no holidays are
// observed beyond the configured weekend rule.
type Config struct {
	// Calendar defines the underlying calendar semantics, including
	// location, week boundaries, and fiscal-year configuration.
	Calendar *calendar.Calendar

	// Weekend defines which dates are considered weekends.
	//
	// This can represent a standard Saturday/Sunday weekend, an
	// organization-specific working week, alternate Saturdays, or any
	// other custom weekend rule.
	Weekend WeekendRule

	// Holidays defines which dates are considered holidays.
	//
	// A nil HolidayCalendar means that no additional holidays are
	// observed.
	Holidays HolidayCalendar
}

// BusinessCalendar represents a configured business calendar.
//
// A BusinessCalendar combines a base calendar with weekend and holiday rules.
// It provides the context required for business-day calculations, navigation,
// arithmetic, and adjustments.
//
// A BusinessCalendar should be created with New and treated as immutable after
// construction.
type BusinessCalendar struct {
	// calendar defines the underlying calendar semantics.
	calendar *calendar.Calendar

	// weekend determines whether a given date is considered a weekend.
	weekend WeekendRule

	// holidays determines whether a given date is a holiday.
	holidays HolidayCalendar
}

// New creates a BusinessCalendar from the supplied configuration.
//
// Calendar and Weekend are required. Holidays are optional.
//
// Example:
//
//	businessCal, err := business.New(business.Config{
//		Calendar: cal,
//		Weekend:  business.StandardWeekend(),
//	})
func New(config Config) (*BusinessCalendar, error) {
	if config.Calendar == nil {
		return nil, ErrNilCalendar
	}

	if config.Weekend == nil {
		return nil, ErrNilWeekend
	}

	return &BusinessCalendar{
		calendar: config.Calendar,
		weekend:  config.Weekend,
		holidays: config.Holidays,
	}, nil
}

// Calendar returns the underlying calendar.Calendar.
//
// The returned calendar defines the timezone and other calendar semantics
// used by the BusinessCalendar.
func (b *BusinessCalendar) Calendar() *calendar.Calendar {
	return b.calendar
}

// Weekend returns the WeekendRule used by the BusinessCalendar.
func (b *BusinessCalendar) Weekend() WeekendRule {
	return b.weekend
}

// Holidays returns the HolidayCalendar used by the BusinessCalendar.
//
// It returns nil when no holiday calendar was configured.
func (b *BusinessCalendar) Holidays() HolidayCalendar {
	return b.holidays
}
