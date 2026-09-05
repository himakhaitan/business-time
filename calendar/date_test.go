package calendar

import (
	"testing"
	"time"
)

func newTestCalendar(t *testing.T, location *time.Location) *Calendar {
	t.Helper()

	cal, err := New(Config{
		Location:        location,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})
	if err != nil {
		t.Fatalf("failed to create calendar: %v", err)
	}

	return cal
}

func TestCalendar_StartOfDay(t *testing.T) {
	tests := []struct {
		name     string
		location *time.Location
		value    time.Time
		want     string
	}{
		{
			name:     "UTC",
			location: time.UTC,
			value:    time.Date(2026, time.January, 15, 14, 30, 45, 123456789, time.UTC),
			want:     "2026-01-15T00:00:00Z",
		},
		{
			name:     "Asia/Kolkata",
			location: mustLoadLocation(t, "Asia/Kolkata"),
			value:    time.Date(2026, time.January, 15, 23, 59, 59, 999999999, mustLoadLocation(t, "Asia/Kolkata")),
			want:     "2026-01-15T00:00:00+05:30",
		},
		{
			name:     "America/New_York",
			location: mustLoadLocation(t, "America/New_York"),
			value:    time.Date(2026, time.July, 15, 18, 30, 0, 0, time.UTC),
			want:     "2026-07-15T00:00:00-04:00",
		},
		{
			name:     "different input location",
			location: mustLoadLocation(t, "Asia/Kolkata"),
			value:    time.Date(2026, time.January, 15, 23, 30, 0, 0, time.UTC),
			want:     "2026-01-16T00:00:00+05:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := newTestCalendar(t, tt.location)

			got := cal.StartOfDay(tt.value)

			if got.Format(time.RFC3339Nano) != tt.want {
				t.Errorf(
					"StartOfDay() = %s, want %s",
					got.Format(time.RFC3339Nano),
					tt.want,
				)
			}

			if got.Location() != tt.location {
				t.Errorf(
					"StartOfDay() location = %v, want %v",
					got.Location(),
					tt.location,
				)
			}
		})
	}

}

func TestCalendar_StartOfDay_PreservesCalendarDateAcrossInputLocations(t *testing.T) {
	location := mustLoadLocation(t, "America/New_York")
	cal := newTestCalendar(t, location)

	// This is January 16 in Kolkata but January 15 in New York.
	value := time.Date(
		2026,
		time.January,
		16,
		4,
		30,
		0,
		0,
		mustLoadLocation(t, "Asia/Kolkata"),
	)

	got := cal.StartOfDay(value)

	want := time.Date(
		2026,
		time.January,
		15,
		0,
		0,
		0,
		0,
		location,
	)

	if !got.Equal(want) {
		t.Errorf("StartOfDay() = %v, want %v", got, want)
	}

}

func TestCalendar_EndOfDay(t *testing.T) {
	tests := []struct {
		name     string
		location *time.Location
		value    time.Time
	}{
		{
			name:     "UTC",
			location: time.UTC,
			value:    time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Asia/Kolkata",
			location: mustLoadLocation(t, "Asia/Kolkata"),
			value:    time.Date(2026, time.January, 15, 12, 0, 0, 0, mustLoadLocation(t, "Asia/Kolkata")),
		},
		{
			name:     "America/New_York",
			location: mustLoadLocation(t, "America/New_York"),
			value:    time.Date(2026, time.July, 15, 12, 0, 0, 0, mustLoadLocation(t, "America/New_York")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := newTestCalendar(t, tt.location)

			start := cal.StartOfDay(tt.value)
			next := cal.StartOfNextDay(tt.value)
			end := cal.EndOfDay(tt.value)

			// EndOfDay is the final representable nanosecond before
			// the beginning of the next calendar day.
			want := next.Add(-time.Nanosecond)

			if !end.Equal(want) {
				t.Errorf("EndOfDay() = %v, want %v", end, want)
			}

			if !end.Before(next) {
				t.Errorf("EndOfDay() = %v, should be before StartOfNextDay() = %v", end, next)
			}

			if end.Before(start) {
				t.Errorf("EndOfDay() = %v, should not be before StartOfDay() = %v", end, start)
			}

			if end.Location() != tt.location {
				t.Errorf(
					"EndOfDay() location = %v, want %v",
					end.Location(),
					tt.location,
				)
			}
		})
	}

}

func TestCalendar_EndOfDay_IsInclusiveBoundary(t *testing.T) {
	location := mustLoadLocation(t, "Asia/Kolkata")
	cal := newTestCalendar(t, location)

	value := time.Date(
		2026,
		time.January,
		15,
		12,
		0,
		0,
		0,
		location,
	)

	end := cal.EndOfDay(value)
	next := cal.StartOfNextDay(value)

	// The instant immediately after EndOfDay is the beginning
	// of the next calendar day.
	if !end.Add(time.Nanosecond).Equal(next) {
		t.Errorf(
			"EndOfDay() + 1ns = %v, want %v",
			end.Add(time.Nanosecond),
			next,
		)
	}

}

func TestCalendar_StartOfNextDay(t *testing.T) {
	tests := []struct {
		name     string
		location *time.Location
		value    time.Time
		want     time.Time
	}{
		{
			name:     "normal day",
			location: time.UTC,
			value:    time.Date(2026, time.January, 15, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2026, time.January, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "month boundary",
			location: time.UTC,
			value:    time.Date(2026, time.January, 31, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "year boundary",
			location: time.UTC,
			value:    time.Date(2026, time.December, 31, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := newTestCalendar(t, tt.location)

			got := cal.StartOfNextDay(tt.value)

			if !got.Equal(tt.want) {
				t.Errorf("StartOfNextDay() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestCalendar_StartOfPreviousDay(t *testing.T) {
	tests := []struct {
		name     string
		location *time.Location
		value    time.Time
		want     time.Time
	}{
		{
			name:     "normal day",
			location: time.UTC,
			value:    time.Date(2026, time.January, 15, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2026, time.January, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "month boundary",
			location: time.UTC,
			value:    time.Date(2026, time.February, 1, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "year boundary",
			location: time.UTC,
			value:    time.Date(2026, time.January, 1, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := newTestCalendar(t, tt.location)

			got := cal.StartOfPreviousDay(tt.value)

			if !got.Equal(tt.want) {
				t.Errorf("StartOfPreviousDay() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestCalendar_StartOfNextDay_DST(t *testing.T) {
	tests := []struct {
		name     string
		location string
		date     time.Time
	}{
		{
			name:     "spring forward",
			location: "America/New_York",
			date:     time.Date(2026, time.March, 8, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "fall back",
			location: "America/New_York",
			date:     time.Date(2026, time.November, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "spring forward Los Angeles",
			location: "America/Los_Angeles",
			date:     time.Date(2026, time.March, 8, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "fall back Los Angeles",
			location: "America/Los_Angeles",
			date:     time.Date(2026, time.November, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			start := cal.StartOfDay(tt.date)
			next := cal.StartOfNextDay(tt.date)

			if next.Hour() != 0 || next.Minute() != 0 || next.Second() != 0 || next.Nanosecond() != 0 {
				t.Fatalf("StartOfNextDay() = %v, want local midnight", next)
			}

			if !next.After(start) {
				t.Fatalf("StartOfNextDay() = %v, want after StartOfDay() = %v", next, start)
			}

			localStart := start.In(location)
			localNext := next.In(location)

			if localNext.YearDay() != localStart.YearDay()+1 {
				// Handle year boundary separately.
				if !(localStart.Month() == time.December &&
					localStart.Day() == 31 &&
					localNext.Month() == time.January &&
					localNext.Day() == 1) {
					t.Errorf(
						"StartOfNextDay() date = %v, want the following local calendar day after %v",
						localNext,
						localStart,
					)
				}
			}

			// DST days are not necessarily 24 elapsed hours.
			elapsed := next.Sub(start)

			if elapsed != 23*time.Hour && elapsed != 24*time.Hour && elapsed != 25*time.Hour {
				t.Errorf(
					"unexpected calendar-day duration: %v",
					elapsed,
				)
			}
		})
	}

}

func TestCalendar_StartOfPreviousDay_DST(t *testing.T) {
	tests := []struct {
		name     string
		location string
		date     time.Time
	}{
		{
			name:     "day after spring forward",
			location: "America/New_York",
			date:     time.Date(2026, time.March, 9, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "day after fall back",
			location: "America/New_York",
			date:     time.Date(2026, time.November, 2, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			current := cal.StartOfDay(tt.date)
			previous := cal.StartOfPreviousDay(tt.date)

			if !previous.Before(current) {
				t.Fatalf(
					"StartOfPreviousDay() = %v, want before StartOfDay() = %v",
					previous,
					current,
				)
			}

			if previous.Hour() != 0 || previous.Minute() != 0 || previous.Second() != 0 || previous.Nanosecond() != 0 {
				t.Errorf("StartOfPreviousDay() = %v, want local midnight", previous)
			}
		})
	}

}

func TestCalendar_IsSameDay(t *testing.T) {
	location := mustLoadLocation(t, "Asia/Kolkata")
	cal := newTestCalendar(t, location)

	tests := []struct {
		name string
		a    time.Time
		b    time.Time
		want bool
	}{
		{
			name: "same instant",
			a:    time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC),
			b:    time.Date(2026, time.January, 15, 10, 0, 0, 0, time.UTC),
			want: true,
		},
		{
			name: "different times same local day",
			a:    time.Date(2026, time.January, 15, 0, 1, 0, 0, location),
			b:    time.Date(2026, time.January, 15, 23, 59, 59, 0, location),
			want: true,
		},
		{
			name: "different local days",
			a:    time.Date(2026, time.January, 15, 23, 59, 59, 0, location),
			b:    time.Date(2026, time.January, 16, 0, 0, 0, 0, location),
			want: false,
		},
		{
			name: "different UTC dates but same Kolkata day",
			a:    time.Date(2026, time.January, 15, 18, 30, 0, 0, time.UTC),
			b:    time.Date(2026, time.January, 15, 19, 30, 0, 0, time.UTC),
			want: true,
		},
		{
			name: "different UTC dates and different Kolkata days",
			a:    time.Date(2026, time.January, 15, 17, 0, 0, 0, time.UTC),
			b:    time.Date(2026, time.January, 15, 18, 31, 0, 0, time.UTC),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cal.IsSameDay(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("IsSameDay() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestCalendar_IsSameDay_UsesCalendarLocation(t *testing.T) {
	tests := []struct {
		name     string
		location string
		a        time.Time
		b        time.Time
		want     bool
	}{
		{
			name:     "same day in UTC",
			location: "UTC",
			a:        time.Date(2026, time.January, 15, 23, 30, 0, 0, time.UTC),
			b:        time.Date(2026, time.January, 16, 0, 30, 0, 0, time.UTC),
			want:     false,
		},
		{
			name:     "same day in Kolkata",
			location: "Asia/Kolkata",
			a:        time.Date(2026, time.January, 15, 23, 30, 0, 0, time.UTC),
			b:        time.Date(2026, time.January, 16, 0, 30, 0, 0, time.UTC),
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			got := cal.IsSameDay(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("IsSameDay() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestCalendar_DayOfYear(t *testing.T) {
	tests := []struct {
		name     string
		location string
		value    time.Time
		want     int
	}{
		{
			name:     "January first",
			location: "UTC",
			value:    time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC),
			want:     1,
		},
		{
			name:     "January thirty first",
			location: "UTC",
			value:    time.Date(2026, time.January, 31, 12, 0, 0, 0, time.UTC),
			want:     31,
		},
		{
			name:     "December thirty first common year",
			location: "UTC",
			value:    time.Date(2025, time.December, 31, 12, 0, 0, 0, time.UTC),
			want:     365,
		},
		{
			name:     "December thirty first leap year",
			location: "UTC",
			value:    time.Date(2024, time.December, 31, 12, 0, 0, 0, time.UTC),
			want:     366,
		},
		{
			name:     "leap day",
			location: "UTC",
			value:    time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC),
			want:     60,
		},
		{
			name:     "Kolkata local date",
			location: "Asia/Kolkata",
			value:    time.Date(2026, time.January, 15, 0, 30, 0, 0, time.UTC),
			want:     15,
		},
		{
			name:     "timezone changes local date",
			location: "Asia/Kolkata",
			value:    time.Date(2026, time.January, 15, 20, 0, 0, 0, time.UTC),
			want:     16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			got := cal.DayOfYear(tt.value)

			if got != tt.want {
				t.Errorf("DayOfYear() = %d, want %d", got, tt.want)
			}
		})
	}

}

func TestCalendar_DayOfYear_UsesCalendarLocation(t *testing.T) {
	value := time.Date(
		2026,
		time.January,
		15,
		23,
		30,
		0,
		0,
		time.UTC,
	)

	tests := []struct {
		location string
		want     int
	}{
		{
			location: "UTC",
			want:     15,
		},
		{
			location: "Asia/Kolkata",
			want:     16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.location, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			got := cal.DayOfYear(value)

			if got != tt.want {
				t.Errorf("DayOfYear() = %d, want %d", got, tt.want)
			}
		})
	}

}

func TestCalendar_DayBoundaries_FormCompleteHalfOpenInterval(t *testing.T) {
	tests := []struct {
		name     string
		location string
		value    time.Time
	}{
		{
			name:     "UTC",
			location: "UTC",
			value:    time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Kolkata",
			location: "Asia/Kolkata",
			value:    time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "New York spring forward",
			location: "America/New_York",
			value:    time.Date(2026, time.March, 8, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "New York fall back",
			location: "America/New_York",
			value:    time.Date(2026, time.November, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			start := cal.StartOfDay(tt.value)
			end := cal.StartOfNextDay(tt.value)

			if !start.Before(end) {
				t.Fatalf("expected start < end, got %v >= %v", start, end)
			}

			if !cal.IsSameDay(start, tt.value) {
				t.Errorf("StartOfDay() is not on the same calendar day as value")
			}

			// The exclusive end belongs to the next calendar day.
			if cal.IsSameDay(end, tt.value) {
				t.Errorf("StartOfNextDay() should belong to the next calendar day")
			}

			// EndOfDay is exactly the final representable instant before
			// the exclusive end boundary.
			if !cal.EndOfDay(tt.value).Add(time.Nanosecond).Equal(end) {
				t.Errorf("EndOfDay() does not immediately precede StartOfNextDay()")
			}
		})
	}

}

func TestCalendar_DayOperations_AreStableAcrossManyTimezones(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"Asia/Tokyo",
		"Asia/Singapore",
		"Europe/London",
		"Europe/Berlin",
		"America/New_York",
		"America/Los_Angeles",
		"America/Chicago",
		"Australia/Sydney",
		"Pacific/Auckland",
	}

	value := time.Date(
		2026,
		time.July,
		15,
		12,
		34,
		56,
		789000000,
		time.UTC,
	)

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)
			cal := newTestCalendar(t, location)

			start := cal.StartOfDay(value)
			next := cal.StartOfNextDay(value)
			previous := cal.StartOfPreviousDay(value)
			end := cal.EndOfDay(value)

			if start.Location() != location {
				t.Errorf("StartOfDay() location = %v, want %v", start.Location(), location)
			}

			if next.Location() != location {
				t.Errorf("StartOfNextDay() location = %v, want %v", next.Location(), location)
			}

			if previous.Location() != location {
				t.Errorf("StartOfPreviousDay() location = %v, want %v", previous.Location(), location)
			}

			if end.Location() != location {
				t.Errorf("EndOfDay() location = %v, want %v", end.Location(), location)
			}

			if !previous.Before(start) {
				t.Errorf("previous day %v should be before start %v", previous, start)
			}

			if !start.Before(next) {
				t.Errorf("start %v should be before next %v", start, next)
			}

			if !end.Before(next) {
				t.Errorf("end %v should be before next %v", end, next)
			}

			if !end.Add(time.Nanosecond).Equal(next) {
				t.Errorf("end + 1ns = %v, want %v", end.Add(time.Nanosecond), next)
			}

			if !cal.IsSameDay(start, value) {
				t.Errorf("StartOfDay() is not on value's calendar day")
			}

			if cal.IsSameDay(next, value) {
				t.Errorf("StartOfNextDay() is still on value's calendar day")
			}
		})
	}

}

func mustLoadLocation(t *testing.T, name string) *time.Location {
	t.Helper()

	location, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("failed to load timezone %q: %v", name, err)
	}

	return location

}
