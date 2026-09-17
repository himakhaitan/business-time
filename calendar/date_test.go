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
		want     time.Time
	}{
		{
			name:     "UTC",
			location: time.UTC,
			value:    time.Date(2026, time.January, 15, 14, 30, 45, 123456789, time.UTC),
			want:     time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Asia/Kolkata",
			location: mustLoadLocation(t, "Asia/Kolkata"),
			value: time.Date(
				2026, time.January, 15,
				23, 59, 59, 999999999,
				mustLoadLocation(t, "Asia/Kolkata"),
			),
			want: time.Date(
				2026, time.January, 15,
				0, 0, 0, 0,
				mustLoadLocation(t, "Asia/Kolkata"),
			),
		},
		{
			name:     "America/New_York",
			location: mustLoadLocation(t, "America/New_York"),
			value:    time.Date(2026, time.July, 15, 18, 30, 0, 0, time.UTC),
			want: time.Date(
				2026, time.July, 15,
				0, 0, 0, 0,
				mustLoadLocation(t, "America/New_York"),
			),
		},
		{
			name:     "input location differs from calendar location",
			location: mustLoadLocation(t, "Asia/Kolkata"),
			value:    time.Date(2026, time.January, 15, 23, 30, 0, 0, time.UTC),
			want: time.Date(
				2026, time.January, 16,
				0, 0, 0, 0,
				mustLoadLocation(t, "Asia/Kolkata"),
			),
		},
		{
			name:     "exact midnight",
			location: time.UTC,
			value:    time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC),
			want:     time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "last nanosecond of day",
			location: time.UTC,
			value:    time.Date(2026, time.January, 15, 23, 59, 59, 999999999, time.UTC),
			want:     time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "month boundary",
			location: time.UTC,
			value:    time.Date(2026, time.February, 1, 12, 0, 0, 0, time.UTC),
			want:     time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "year boundary",
			location: time.UTC,
			value:    time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC),
			want:     time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := newTestCalendar(t, tt.location)

			got := cal.StartOfDay(tt.value)

			if !got.Equal(tt.want) {
				t.Errorf("StartOfDay() = %v, want %v", got, tt.want)
			}

			if got.Location() != tt.location {
				t.Errorf(
					"StartOfDay() location = %v, want %v",
					got.Location(),
					tt.location,
				)
			}

			if got.Hour() != 0 ||
				got.Minute() != 0 ||
				got.Second() != 0 ||
				got.Nanosecond() != 0 {
				t.Errorf("StartOfDay() = %v, want local midnight", got)
			}
		})
	}
}

func TestCalendar_StartOfDay_UsesCalendarLocation(t *testing.T) {
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
		name     string
		location string
		wantDate time.Time
	}{
		{
			name:     "UTC",
			location: "UTC",
			wantDate: time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Asia/Kolkata",
			location: "Asia/Kolkata",
			wantDate: time.Date(
				2026,
				time.January,
				16,
				0,
				0,
				0,
				0,
				mustLoadLocation(t, "Asia/Kolkata"),
			),
		},
		{
			name:     "America/New_York",
			location: "America/New_York",
			wantDate: time.Date(
				2026,
				time.January,
				15,
				0,
				0,
				0,
				0,
				mustLoadLocation(t, "America/New_York"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			got := cal.StartOfDay(value)

			if !got.Equal(tt.wantDate) {
				t.Errorf("StartOfDay() = %v, want %v", got, tt.wantDate)
			}
		})
	}
}

func TestCalendar_StartOfDay_PreservesCalendarDateAcrossInputLocations(t *testing.T) {
	calendarLocation := mustLoadLocation(t, "America/New_York")
	inputLocation := mustLoadLocation(t, "Asia/Kolkata")

	// This instant is January 16 in Kolkata but January 15 in New York.
	value := time.Date(
		2026,
		time.January,
		16,
		4,
		30,
		0,
		0,
		inputLocation,
	)

	cal := newTestCalendar(t, calendarLocation)

	got := cal.StartOfDay(value)

	want := time.Date(
		2026,
		time.January,
		15,
		0,
		0,
		0,
		0,
		calendarLocation,
	)

	if !got.Equal(want) {
		t.Errorf("StartOfDay() = %v, want %v", got, want)
	}
}

func TestCalendar_EndOfDay(t *testing.T) {
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
			name:     "Asia/Kolkata",
			location: "Asia/Kolkata",
			value: time.Date(
				2026,
				time.January,
				15,
				12,
				0,
				0,
				0,
				mustLoadLocation(t, "Asia/Kolkata"),
			),
		},
		{
			name:     "America/New_York",
			location: "America/New_York",
			value: time.Date(
				2026,
				time.July,
				15,
				12,
				0,
				0,
				0,
				mustLoadLocation(t, "America/New_York"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			start := cal.StartOfDay(tt.value)
			next := cal.StartOfNextDay(tt.value)
			end := cal.EndOfDay(tt.value)

			want := next.Add(-time.Nanosecond)

			if !end.Equal(want) {
				t.Errorf("EndOfDay() = %v, want %v", end, want)
			}

			if !end.Before(next) {
				t.Errorf(
					"EndOfDay() = %v, should be before StartOfNextDay() = %v",
					end,
					next,
				)
			}

			if end.Before(start) {
				t.Errorf(
					"EndOfDay() = %v, should not be before StartOfDay() = %v",
					end,
					start,
				)
			}

			if end.Location() != location {
				t.Errorf(
					"EndOfDay() location = %v, want %v",
					end.Location(),
					location,
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

	if !end.Add(time.Nanosecond).Equal(next) {
		t.Errorf(
			"EndOfDay() + 1ns = %v, want %v",
			end.Add(time.Nanosecond),
			next,
		)
	}

	if !cal.IsSameDay(end, value) {
		t.Errorf("EndOfDay() = %v, want same calendar day as %v", end, value)
	}

	if cal.IsSameDay(next, value) {
		t.Errorf("StartOfNextDay() = %v, should be next calendar day", next)
	}
}

func TestCalendar_StartOfNextDay(t *testing.T) {
	tests := []struct {
		name     string
		location string
		value    time.Time
		want     time.Time
	}{
		{
			name:     "normal day",
			location: "UTC",
			value:    time.Date(2026, time.January, 15, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2026, time.January, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "month boundary",
			location: "UTC",
			value:    time.Date(2026, time.January, 31, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "year boundary",
			location: "UTC",
			value:    time.Date(2026, time.December, 31, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			got := cal.StartOfNextDay(tt.value)

			if !got.Equal(tt.want) {
				t.Errorf("StartOfNextDay() = %v, want %v", got, tt.want)
			}

			if got.Location() != location {
				t.Errorf(
					"StartOfNextDay() location = %v, want %v",
					got.Location(),
					location,
				)
			}
		})
	}
}

func TestCalendar_StartOfPreviousDay(t *testing.T) {
	tests := []struct {
		name     string
		location string
		value    time.Time
		want     time.Time
	}{
		{
			name:     "normal day",
			location: "UTC",
			value:    time.Date(2026, time.January, 15, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2026, time.January, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "month boundary",
			location: "UTC",
			value:    time.Date(2026, time.February, 1, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "year boundary",
			location: "UTC",
			value:    time.Date(2026, time.January, 1, 15, 30, 0, 0, time.UTC),
			want:     time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			got := cal.StartOfPreviousDay(tt.value)

			if !got.Equal(tt.want) {
				t.Errorf("StartOfPreviousDay() = %v, want %v", got, tt.want)
			}

			if got.Location() != location {
				t.Errorf(
					"StartOfPreviousDay() location = %v, want %v",
					got.Location(),
					location,
				)
			}
		})
	}
}

func TestCalendar_StartOfNextDay_DST(t *testing.T) {
	tests := []struct {
		name         string
		location     string
		value        time.Time
		wantDuration time.Duration
	}{
		{
			name:         "New York spring forward",
			location:     "America/New_York",
			value:        time.Date(2026, time.March, 8, 12, 0, 0, 0, time.UTC),
			wantDuration: 23 * time.Hour,
		},
		{
			name:         "New York fall back",
			location:     "America/New_York",
			value:        time.Date(2026, time.November, 1, 12, 0, 0, 0, time.UTC),
			wantDuration: 25 * time.Hour,
		},
		{
			name:         "Los Angeles spring forward",
			location:     "America/Los_Angeles",
			value:        time.Date(2026, time.March, 8, 12, 0, 0, 0, time.UTC),
			wantDuration: 23 * time.Hour,
		},
		{
			name:         "Los Angeles fall back",
			location:     "America/Los_Angeles",
			value:        time.Date(2026, time.November, 1, 12, 0, 0, 0, time.UTC),
			wantDuration: 25 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			start := cal.StartOfDay(tt.value)
			next := cal.StartOfNextDay(tt.value)

			if next.Hour() != 0 ||
				next.Minute() != 0 ||
				next.Second() != 0 ||
				next.Nanosecond() != 0 {
				t.Fatalf("StartOfNextDay() = %v, want local midnight", next)
			}

			if !cal.IsSameDay(start, tt.value) {
				t.Fatalf("StartOfDay() = %v, not same day as value %v", start, tt.value)
			}

			if cal.IsSameDay(next, tt.value) {
				t.Fatalf("StartOfNextDay() = %v, still belongs to value's day", next)
			}

			if got := next.Sub(start); got != tt.wantDuration {
				t.Errorf(
					"calendar day duration = %v, want %v",
					got,
					tt.wantDuration,
				)
			}
		})
	}
}

func TestCalendar_StartOfPreviousDay_DST(t *testing.T) {
	tests := []struct {
		name     string
		location string
		value    time.Time
	}{
		{
			name:     "day after spring forward",
			location: "America/New_York",
			value:    time.Date(2026, time.March, 9, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "day after fall back",
			location: "America/New_York",
			value:    time.Date(2026, time.November, 2, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "day after spring forward Los Angeles",
			location: "America/Los_Angeles",
			value:    time.Date(2026, time.March, 9, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "day after fall back Los Angeles",
			location: "America/Los_Angeles",
			value:    time.Date(2026, time.November, 2, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			current := cal.StartOfDay(tt.value)
			previous := cal.StartOfPreviousDay(tt.value)

			if !previous.Before(current) {
				t.Fatalf(
					"StartOfPreviousDay() = %v, want before StartOfDay() = %v",
					previous,
					current,
				)
			}

			if previous.Hour() != 0 ||
				previous.Minute() != 0 ||
				previous.Second() != 0 ||
				previous.Nanosecond() != 0 {
				t.Errorf("StartOfPreviousDay() = %v, want local midnight", previous)
			}

			if cal.IsSameDay(previous, current) {
				t.Errorf(
					"previous day %v should not be same calendar day as current %v",
					previous,
					current,
				)
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
			name: "Kolkata date boundary",
			a:    time.Date(2026, time.January, 15, 18, 29, 59, 0, time.UTC),
			b:    time.Date(2026, time.January, 15, 18, 30, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "same instant different input locations",
			a:    time.Date(2026, time.January, 15, 18, 30, 0, 0, time.UTC),
			b: time.Date(
				2026,
				time.January,
				16,
				0,
				0,
				0,
				0,
				location,
			),
			want: true,
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
			name:     "different days in UTC",
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
		{
			name:     "same instant different zones",
			location: "America/New_York",
			a:        time.Date(2026, time.January, 15, 18, 0, 0, 0, time.UTC),
			b: time.Date(
				2026,
				time.January,
				15,
				13,
				0,
				0,
				0,
				mustLoadLocation(t, "America/New_York"),
			),
			want: true,
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
		{
			location: "America/New_York",
			want:     15,
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
				t.Fatalf(
					"expected start < end, got %v >= %v",
					start,
					end,
				)
			}

			if !cal.IsSameDay(start, tt.value) {
				t.Errorf("StartOfDay() is not on value's calendar day")
			}

			if cal.IsSameDay(end, tt.value) {
				t.Errorf("StartOfNextDay() belongs to value's calendar day")
			}

			if !cal.EndOfDay(tt.value).Add(time.Nanosecond).Equal(end) {
				t.Errorf(
					"EndOfDay() does not immediately precede StartOfNextDay()",
				)
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
				t.Errorf(
					"StartOfNextDay() location = %v, want %v",
					next.Location(),
					location,
				)
			}

			if previous.Location() != location {
				t.Errorf(
					"StartOfPreviousDay() location = %v, want %v",
					previous.Location(),
					location,
				)
			}

			if end.Location() != location {
				t.Errorf(
					"EndOfDay() location = %v, want %v",
					end.Location(),
					location,
				)
			}

			if !previous.Before(start) {
				t.Errorf(
					"previous day %v should be before start %v",
					previous,
					start,
				)
			}

			if !start.Before(next) {
				t.Errorf(
					"start %v should be before next %v",
					start,
					next,
				)
			}

			if !end.Before(next) {
				t.Errorf(
					"end %v should be before next %v",
					end,
					next,
				)
			}

			if !end.Add(time.Nanosecond).Equal(next) {
				t.Errorf(
					"end + 1ns = %v, want %v",
					end.Add(time.Nanosecond),
					next,
				)
			}

			if !cal.IsSameDay(start, value) {
				t.Errorf(
					"StartOfDay() is not on value's calendar day",
				)
			}

			if cal.IsSameDay(next, value) {
				t.Errorf(
					"StartOfNextDay() is still on value's calendar day",
				)
			}
		})
	}
}

func TestCalendar_DayOperations_AroundMidnight(t *testing.T) {
	location := mustLoadLocation(t, "Asia/Kolkata")
	cal := newTestCalendar(t, location)

	tests := []struct {
		name  string
		value time.Time
		want  time.Time
	}{
		{
			name: "one nanosecond before midnight",
			value: time.Date(
				2026,
				time.January,
				15,
				23,
				59,
				59,
				999999999,
				location,
			),
			want: time.Date(
				2026,
				time.January,
				15,
				0,
				0,
				0,
				0,
				location,
			),
		},
		{
			name: "exact midnight",
			value: time.Date(
				2026,
				time.January,
				16,
				0,
				0,
				0,
				0,
				location,
			),
			want: time.Date(
				2026,
				time.January,
				16,
				0,
				0,
				0,
				0,
				location,
			),
		},
		{
			name: "one nanosecond after midnight",
			value: time.Date(
				2026,
				time.January,
				16,
				0,
				0,
				0,
				1,
				location,
			),
			want: time.Date(
				2026,
				time.January,
				16,
				0,
				0,
				0,
				0,
				location,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cal.StartOfDay(tt.value)

			if !got.Equal(tt.want) {
				t.Errorf("StartOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_DayOperations_AroundMonthAndYearBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		location string
		value    time.Time
	}{
		{
			name:     "last day of January",
			location: "UTC",
			value:    time.Date(2026, time.January, 31, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "first day of February",
			location: "UTC",
			value:    time.Date(2026, time.February, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "last day of December",
			location: "UTC",
			value:    time.Date(2026, time.December, 31, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "first day of January",
			location: "UTC",
			value:    time.Date(2027, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "leap day",
			location: "UTC",
			value:    time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)
			cal := newTestCalendar(t, location)

			start := cal.StartOfDay(tt.value)
			previous := cal.StartOfPreviousDay(tt.value)
			next := cal.StartOfNextDay(tt.value)

			if !previous.Before(start) {
				t.Errorf("previous = %v, want before start = %v", previous, start)
			}

			if !start.Before(next) {
				t.Errorf("start = %v, want before next = %v", start, next)
			}

			if !cal.IsSameDay(start, tt.value) {
				t.Errorf("start = %v, not same day as value = %v", start, tt.value)
			}

			if cal.IsSameDay(previous, tt.value) {
				t.Errorf("previous = %v, unexpectedly same day as value", previous)
			}

			if cal.IsSameDay(next, tt.value) {
				t.Errorf("next = %v, unexpectedly same day as value", next)
			}
		})
	}
}

func TestCalendar_DayOperations_PreserveCalendarLocation(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	value := time.Date(
		2026,
		time.July,
		15,
		18,
		45,
		30,
		123456789,
		time.UTC,
	)

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)
			cal := newTestCalendar(t, location)

			results := []time.Time{
				cal.StartOfDay(value),
				cal.EndOfDay(value),
				cal.StartOfNextDay(value),
				cal.StartOfPreviousDay(value),
			}

			for _, result := range results {
				if result.Location() != location {
					t.Errorf(
						"result location = %v, want %v",
						result.Location(),
						location,
					)
				}
			}
		})
	}
}

func TestCalendar_DayOperations_DifferentInputLocations(t *testing.T) {
	calendarLocation := mustLoadLocation(t, "America/New_York")
	cal := newTestCalendar(t, calendarLocation)

	instant := time.Date(
		2026,
		time.January,
		16,
		2,
		0,
		0,
		0,
		time.UTC,
	)

	inputLocations := []string{
		"UTC",
		"Asia/Kolkata",
		"Asia/Tokyo",
		"Europe/London",
		"America/Los_Angeles",
		"Australia/Sydney",
	}

	var wantStart time.Time
	for i, name := range inputLocations {
		location := mustLoadLocation(t, name)
		value := instant.In(location)

		start := cal.StartOfDay(value)

		if i == 0 {
			wantStart = start
			continue
		}

		if !start.Equal(wantStart) {
			t.Errorf(
				"StartOfDay() for input location %s = %v, want %v",
				name,
				start,
				wantStart,
			)
		}
	}
}

func TestCalendar_DayOperations_DaySequenceInvariant(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"America/Los_Angeles",
		"Europe/Berlin",
		"Australia/Sydney",
	}

	dates := []time.Time{
		time.Date(2024, time.February, 28, 12, 0, 0, 0, time.UTC),
		time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC),
		time.Date(2024, time.March, 10, 12, 0, 0, 0, time.UTC),
		time.Date(2024, time.November, 3, 12, 0, 0, 0, time.UTC),
		time.Date(2025, time.December, 31, 12, 0, 0, 0, time.UTC),
		time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC),
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)
			cal := newTestCalendar(t, location)

			for _, value := range dates {
				start := cal.StartOfDay(value)
				next := cal.StartOfNextDay(value)
				nextNext := cal.StartOfNextDay(next)
				previous := cal.StartOfPreviousDay(value)

				if !next.After(start) {
					t.Errorf(
						"next day %v is not after start %v",
						next,
						start,
					)
				}

				if !nextNext.After(next) {
					t.Errorf(
						"next-next day %v is not after next %v",
						nextNext,
						next,
					)
				}

				if !previous.Before(start) {
					t.Errorf(
						"previous day %v is not before start %v",
						previous,
						start,
					)
				}

				if cal.IsSameDay(next, value) {
					t.Errorf(
						"next day %v should not be same day as value %v",
						next,
						value,
					)
				}

				if cal.IsSameDay(previous, value) {
					t.Errorf(
						"previous day %v should not be same day as value %v",
						previous,
						value,
					)
				}
			}
		})
	}
}

func TestCalendar_EndOfDay_DoesNotCrossIntoNextDay(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"America/Los_Angeles",
		"Europe/Berlin",
		"Australia/Sydney",
	}

	values := []time.Time{
		time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
		time.Date(2026, time.March, 8, 12, 0, 0, 0, time.UTC),
		time.Date(2026, time.November, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2026, time.December, 31, 12, 0, 0, 0, time.UTC),
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)
			cal := newTestCalendar(t, location)

			for _, value := range values {
				end := cal.EndOfDay(value)
				next := cal.StartOfNextDay(value)

				if !cal.IsSameDay(end, value) {
					t.Errorf(
						"EndOfDay() = %v does not belong to value's day %v",
						end,
						value,
					)
				}

				if cal.IsSameDay(end, next) {
					t.Errorf(
						"EndOfDay() = %v unexpectedly belongs to next day %v",
						end,
						next,
					)
				}

				if !end.Add(time.Nanosecond).Equal(next) {
					t.Errorf(
						"EndOfDay() + 1ns = %v, want %v",
						end.Add(time.Nanosecond),
						next,
					)
				}
			}
		})
	}
}
