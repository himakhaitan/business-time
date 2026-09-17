package business

import "time"

// WeekendRule determines whether a given date is considered a weekend.
//
// A WeekendRule allows BusinessCalendar to support both conventional and
// organization-specific working-week definitions. The rule may represent
// simple weekly weekends, alternate Saturdays, regional working-week
// patterns, or any other recurring or custom weekend logic.
//
// BusinessCalendar is responsible for interpreting time.Time values using
// its configured timezone before applying the weekend rule. Implementations
// therefore receive t in the relevant business calendar's location.
//
// The rule should generally depend on the calendar date represented by t,
// rather than its time of day.
type WeekendRule interface {
	// IsWeekend reports whether the calendar date represented by t is a
	// weekend.
	//
	// The date is interpreted according to the timezone of the enclosing
	// BusinessCalendar. Implementations should generally use the local
	// calendar date represented by t rather than treating t as an absolute
	// instant.
	//
	// The result should normally be independent of the time of day. If t
	// represents any time within the same local calendar day, IsWeekend
	// should return the same result.
	IsWeekend(t time.Time) bool
}

// FixedWeekend defines a weekend using a fixed set of weekdays.
//
// It is suitable for conventional working-week definitions where the same
// weekdays are considered weekends every week. For example, a Saturday and
// Sunday weekend can be created with:
//
//	weekend := NewFixedWeekend(time.Saturday, time.Sunday)
//
// For calendars with more complex rules, such as alternate Saturdays or
// organization-specific working Saturdays, implement WeekendRule directly.
type FixedWeekend struct {
	days map[time.Weekday]struct{}
}

// NewFixedWeekend creates a WeekendRule that considers the supplied
// weekdays to be weekends.
//
// Duplicate weekdays are ignored.
//
// Example:
//
//	weekend := business.NewFixedWeekend(
//		time.Saturday,
//		time.Sunday,
//	)
func NewFixedWeekend(days ...time.Weekday) WeekendRule {
	weekend := FixedWeekend{
		days: make(map[time.Weekday]struct{}, len(days)),
	}

	for _, day := range days {
		weekend.days[day] = struct{}{}
	}

	return weekend
}

// IsWeekend reports whether the weekday represented by t is configured as
// a weekend.
//
// The time of day and date are otherwise ignored because FixedWeekend
// defines weekends solely by the day of the week.
func (w FixedWeekend) IsWeekend(t time.Time) bool {
	_, ok := w.days[t.Weekday()]
	return ok
}

// StandardWeekend is the conventional Saturday and Sunday weekend.
//
// It can be used directly when constructing a BusinessCalendar:
//
//	businessCal, err := business.New(business.Config{
//		Calendar: cal,
//		Weekend:  business.StandardWeekend(),
//	})
func StandardWeekend() WeekendRule {
	return NewFixedWeekend(
		time.Saturday,
		time.Sunday,
	)
}

// SundayWeekend returns a WeekendRule that considers Sunday a non-working
// day and all other weekdays as working days.
//
// It is useful for calendars that observe a six-day working week with
// Sunday as the sole weekly weekend day.
func SundayWeekend() WeekendRule {
	return NewFixedWeekend(time.Sunday)
}
