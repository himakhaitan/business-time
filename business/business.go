// Package business provides business-day calendar functionality built on top
// of Go's standard time package and the calendar package.
//
// A BusinessCalendar combines a base calendar with configurable weekend and
// holiday rules to determine business days and perform business-day
// calculations.
//
// The package provides common weekend and holiday implementations while
// allowing applications to define custom rules for organization-specific
// working weeks, regional holidays, alternate Saturdays, and other
// requirements.
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
	// ErrNilCalendar indicates that no calendar was provided when creating a
	// BusinessCalendar.
	ErrNilCalendar = errors.New("business: calendar is required")

	// ErrNilWeekend indicates that no weekend rule was provided when creating
	// a BusinessCalendar.
	ErrNilWeekend = errors.New("business: weekend rule is required")
)

// Config defines the configuration for a BusinessCalendar.
//
// Calendar is required and provides the underlying calendar semantics.
// Weekend defines which weekdays are considered weekends. Holidays is
// optional and can be omitted when the calendar does not observe holidays.
type Config struct {
	// Calendar provides the underlying calendar semantics, including
	// location, week boundaries, and fiscal-year configuration.
	Calendar *calendar.Calendar

	// Weekend defines which dates are considered weekends.
	//
	// If nil, New uses StandardWeekend.
	Weekend WeekendRule

	// Holidays defines which dates are considered holidays.
	//
	// A nil value means that no additional holidays are observed.
	Holidays HolidayCalendar
}

// BusinessCalendar represents a configured business calendar.
//
// It combines a base calendar with weekend and holiday rules to determine
// business days and perform business-day calculations.
//
// BusinessCalendar should be created with New and treated as immutable after
// construction.
type BusinessCalendar struct {
	calendar *calendar.Calendar
	weekend  WeekendRule
	holidays HolidayCalendar
}

// New creates a BusinessCalendar from the supplied configuration.
//
// Calendar is required. If Weekend is nil, StandardWeekend is used.
// Holidays is optional.
//
// Example:
//
//	businessCal, err := business.New(business.Config{
//		Calendar: cal,
//	})
func New(config Config) (*BusinessCalendar, error) {
	if config.Calendar == nil {
		return nil, ErrNilCalendar
	}

	weekend := config.Weekend
	if weekend == nil {
		weekend = StandardWeekend()
	}

	return &BusinessCalendar{
		calendar: config.Calendar,
		weekend:  weekend,
		holidays: config.Holidays,
	}, nil
}

// Calendar returns the underlying calendar.Calendar.
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
