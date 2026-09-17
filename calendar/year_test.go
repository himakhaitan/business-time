package calendar

import (
	"testing"
	"time"
)

func TestCalendar_StartOfYear(t *testing.T) {
	tests := []struct {
		name            string
		fiscalYearStart time.Month
		location        string
		value           time.Time
		want            time.Time
	}{
		{
			name:            "January fiscal year at beginning",
			fiscalYearStart: time.January,
			location:        "UTC",
			value:           time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "January fiscal year at end",
			fiscalYearStart: time.January,
			location:        "UTC",
			value:           time.Date(2026, time.December, 15, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "April fiscal year before fiscal start",
			fiscalYearStart: time.April,
			location:        "UTC",
			value:           time.Date(2026, time.March, 31, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "April fiscal year on fiscal start",
			fiscalYearStart: time.April,
			location:        "UTC",
			value:           time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			want:            time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "April fiscal year after fiscal start",
			fiscalYearStart: time.April,
			location:        "UTC",
			value:           time.Date(2026, time.April, 15, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "April fiscal year at year end",
			fiscalYearStart: time.April,
			location:        "UTC",
			value:           time.Date(2027, time.March, 31, 23, 59, 59, 0, time.UTC),
			want:            time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "July fiscal year before fiscal start",
			fiscalYearStart: time.July,
			location:        "UTC",
			value:           time.Date(2026, time.June, 30, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "July fiscal year after fiscal start",
			fiscalYearStart: time.July,
			location:        "UTC",
			value:           time.Date(2026, time.August, 15, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "October fiscal year before fiscal start",
			fiscalYearStart: time.October,
			location:        "UTC",
			value:           time.Date(2026, time.September, 30, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2025, time.October, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "October fiscal year after fiscal start",
			fiscalYearStart: time.October,
			location:        "UTC",
			value:           time.Date(2026, time.October, 1, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)

			cal := newTestCalendar(t, location)
			cal.fiscalYearStart = tt.fiscalYearStart

			got := cal.StartOfYear(tt.value)

			if !got.Equal(tt.want) {
				t.Errorf("StartOfYear() = %v, want %v", got, tt.want)
			}

			if got.Location() != location {
				t.Errorf(
					"StartOfYear() location = %v, want %v",
					got.Location(),
					location,
				)
			}

			if got.Month() != tt.fiscalYearStart {
				t.Errorf(
					"StartOfYear() month = %v, want fiscal start month %v",
					got.Month(),
					tt.fiscalYearStart,
				)
			}

			if got.Day() != 1 {
				t.Errorf(
					"StartOfYear() day = %d, want 1",
					got.Day(),
				)
			}

			if got.Hour() != 0 ||
				got.Minute() != 0 ||
				got.Second() != 0 ||
				got.Nanosecond() != 0 {
				t.Errorf(
					"StartOfYear() = %v, want local midnight",
					got,
				)
			}
		})
	}
}

func TestCalendar_StartOfYear_UsesCalendarLocation(t *testing.T) {
	calendarLocation := mustLoadLocation(t, "Asia/Kolkata")

	cal := newTestCalendar(t, calendarLocation)
	cal.fiscalYearStart = time.April

	value := time.Date(2026, time.March, 31, 18, 0, 0, 0, time.UTC)
	got := cal.StartOfYear(value)

	want := time.Date(
		2025,
		time.April,
		1,
		0,
		0,
		0,
		0,
		calendarLocation,
	)

	if !got.Equal(want) {
		t.Errorf("StartOfYear() = %v, want %v", got, want)
	}
}

func TestCalendar_StartOfYear_ExactlyAtFiscalBoundary(t *testing.T) {
	tests := []struct {
		name            string
		fiscalYearStart time.Month
		value           time.Time
	}{
		{
			name:            "January",
			fiscalYearStart: time.January,
			value:           time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "April",
			fiscalYearStart: time.April,
			value:           time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "July",
			fiscalYearStart: time.July,
			value:           time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "October",
			fiscalYearStart: time.October,
			value:           time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := newTestCalendar(t, time.UTC)
			cal.fiscalYearStart = tt.fiscalYearStart

			got := cal.StartOfYear(tt.value)

			if !got.Equal(tt.value) {
				t.Errorf(
					"StartOfYear() = %v, want %v",
					got,
					tt.value,
				)
			}
		})
	}
}

func TestCalendar_StartOfNextYear(t *testing.T) {
	tests := []struct {
		name            string
		fiscalYearStart time.Month
		value           time.Time
		want            time.Time
	}{
		{
			name:            "January fiscal year",
			fiscalYearStart: time.January,
			value:           time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "April fiscal year",
			fiscalYearStart: time.April,
			value:           time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2027, time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "October fiscal year",
			fiscalYearStart: time.October,
			value:           time.Date(2026, time.February, 15, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := mustNewCalendar(t, Config{
				Location:        time.UTC,
				WeekStart:       time.Monday,
				FiscalYearStart: tt.fiscalYearStart,
			})

			got := cal.StartOfNextYear(tt.value)

			if !got.Equal(tt.want) {
				t.Errorf(
					"StartOfNextYear() = %v, want %v",
					got,
					tt.want,
				)
			}

			start := cal.StartOfYear(tt.value)

			if !got.After(start) {
				t.Errorf(
					"StartOfNextYear() = %v, want after StartOfYear() = %v",
					got,
					start,
				)
			}
		})
	}
}

func TestCalendar_StartOfPreviousYear(t *testing.T) {
	tests := []struct {
		name            string
		fiscalYearStart time.Month
		value           time.Time
		want            time.Time
	}{
		{
			name:            "January fiscal year",
			fiscalYearStart: time.January,
			value:           time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "April fiscal year",
			fiscalYearStart: time.April,
			value:           time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:            "October fiscal year",
			fiscalYearStart: time.October,
			value:           time.Date(2026, time.February, 15, 12, 0, 0, 0, time.UTC),
			want:            time.Date(2024, time.October, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := mustNewCalendar(t, Config{
				Location:        time.UTC,
				WeekStart:       time.Monday,
				FiscalYearStart: tt.fiscalYearStart,
			})

			got := cal.StartOfPreviousYear(tt.value)

			if !got.Equal(tt.want) {
				t.Errorf(
					"StartOfPreviousYear() = %v, want %v",
					got,
					tt.want,
				)
			}

			start := cal.StartOfYear(tt.value)

			if !got.Before(start) {
				t.Errorf(
					"StartOfPreviousYear() = %v, want before StartOfYear() = %v",
					got,
					start,
				)
			}
		})
	}
}

func TestCalendar_EndOfYear(t *testing.T) {
	tests := []struct {
		name            string
		fiscalYearStart time.Month
		value           time.Time
	}{
		{
			name:            "January fiscal year",
			fiscalYearStart: time.January,
			value:           time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:            "April fiscal year",
			fiscalYearStart: time.April,
			value:           time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:            "October fiscal year",
			fiscalYearStart: time.October,
			value:           time.Date(2026, time.February, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := newTestCalendar(t, time.UTC)
			cal.fiscalYearStart = tt.fiscalYearStart

			start := cal.StartOfYear(tt.value)
			next := cal.StartOfNextYear(tt.value)
			end := cal.EndOfYear(tt.value)

			want := next.Add(-time.Nanosecond)

			if !end.Equal(want) {
				t.Errorf("EndOfYear() = %v, want %v", end, want)
			}

			if !start.Before(end) {
				t.Errorf(
					"StartOfYear() = %v should be before EndOfYear() = %v",
					start,
					end,
				)
			}

			if !end.Before(next) {
				t.Errorf(
					"EndOfYear() = %v should be before StartOfNextYear() = %v",
					end,
					next,
				)
			}

			if !end.Add(time.Nanosecond).Equal(next) {
				t.Errorf(
					"EndOfYear() + 1ns = %v, want %v",
					end.Add(time.Nanosecond),
					next,
				)
			}
		})
	}
}

func TestCalendar_EndOfYear_IsSameFiscalYear(t *testing.T) {
	cal := newTestCalendar(t, time.UTC)
	cal.fiscalYearStart = time.April

	value := time.Date(
		2026,
		time.June,
		15,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	end := cal.EndOfYear(value)
	next := cal.StartOfNextYear(value)

	if !cal.IsSameYear(end, value) {
		t.Errorf(
			"EndOfYear() = %v should belong to same fiscal year as %v",
			end,
			value,
		)
	}

	if cal.IsSameYear(next, value) {
		t.Errorf(
			"StartOfNextYear() = %v should not belong to same fiscal year as %v",
			next,
			value,
		)
	}
}

func TestCalendar_IsSameYear(t *testing.T) {
	tests := []struct {
		name            string
		fiscalYearStart time.Month
		a               time.Time
		b               time.Time
		want            bool
	}{
		{
			name:            "same calendar year",
			fiscalYearStart: time.January,
			a:               time.Date(2026, time.March, 1, 12, 0, 0, 0, time.UTC),
			b:               time.Date(2026, time.November, 1, 12, 0, 0, 0, time.UTC),
			want:            true,
		},
		{
			name:            "different calendar years",
			fiscalYearStart: time.January,
			a:               time.Date(2025, time.December, 31, 12, 0, 0, 0, time.UTC),
			b:               time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC),
			want:            false,
		},
		{
			name:            "same April fiscal year across calendar years",
			fiscalYearStart: time.April,
			a:               time.Date(2026, time.April, 1, 12, 0, 0, 0, time.UTC),
			b:               time.Date(2027, time.March, 31, 12, 0, 0, 0, time.UTC),
			want:            true,
		},
		{
			name:            "different April fiscal years",
			fiscalYearStart: time.April,
			a:               time.Date(2026, time.March, 31, 12, 0, 0, 0, time.UTC),
			b:               time.Date(2026, time.April, 1, 12, 0, 0, 0, time.UTC),
			want:            false,
		},
		{
			name:            "same October fiscal year across calendar years",
			fiscalYearStart: time.October,
			a:               time.Date(2026, time.October, 1, 12, 0, 0, 0, time.UTC),
			b:               time.Date(2027, time.September, 30, 12, 0, 0, 0, time.UTC),
			want:            true,
		},
		{
			name:            "different October fiscal years",
			fiscalYearStart: time.October,
			a:               time.Date(2026, time.September, 30, 12, 0, 0, 0, time.UTC),
			b:               time.Date(2026, time.October, 1, 12, 0, 0, 0, time.UTC),
			want:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := newTestCalendar(t, time.UTC)
			cal.fiscalYearStart = tt.fiscalYearStart

			got := cal.IsSameYear(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("IsSameYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_IsSameYear_UsesCalendarLocation(t *testing.T) {
	location := mustLoadLocation(t, "Asia/Kolkata")

	cal := mustNewCalendar(t, Config{
		Location:        location,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	a := time.Date(
		2025,
		time.December,
		31,
		18,
		0,
		0,
		0,
		time.UTC,
	)

	b := time.Date(
		2026,
		time.January,
		1,
		18,
		0,
		0,
		0,
		time.UTC,
	)

	if cal.IsSameYear(a, b) {
		t.Errorf("IsSameYear() = true, want false")
	}
}

func TestCalendar_YearBoundaries_FormHalfOpenInterval(t *testing.T) {
	tests := []struct {
		name            string
		fiscalYearStart time.Month
		value           time.Time
	}{
		{
			name:            "January fiscal year",
			fiscalYearStart: time.January,
			value:           time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:            "April fiscal year",
			fiscalYearStart: time.April,
			value:           time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:            "July fiscal year",
			fiscalYearStart: time.July,
			value:           time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:            "October fiscal year",
			fiscalYearStart: time.October,
			value:           time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := newTestCalendar(t, time.UTC)
			cal.fiscalYearStart = tt.fiscalYearStart

			start := cal.StartOfYear(tt.value)
			next := cal.StartOfNextYear(tt.value)
			end := cal.EndOfYear(tt.value)

			if !start.Before(next) {
				t.Fatalf(
					"expected start < next, got %v >= %v",
					start,
					next,
				)
			}

			if !end.Before(next) {
				t.Fatalf(
					"expected end < next, got %v >= %v",
					end,
					next,
				)
			}

			if !end.Add(time.Nanosecond).Equal(next) {
				t.Errorf(
					"EndOfYear() + 1ns = %v, want %v",
					end.Add(time.Nanosecond),
					next,
				)
			}

			if !cal.IsSameYear(start, tt.value) {
				t.Errorf(
					"StartOfYear() = %v does not belong to value's fiscal year",
					start,
				)
			}

			if cal.IsSameYear(next, tt.value) {
				t.Errorf(
					"StartOfNextYear() = %v belongs to value's fiscal year",
					next,
				)
			}
		})
	}
}

func TestCalendar_YearBoundaries_DST(t *testing.T) {
	tests := []struct {
		name            string
		location        string
		fiscalYearStart time.Month
		value           time.Time
	}{
		{
			name:            "New York January fiscal year",
			location:        "America/New_York",
			fiscalYearStart: time.January,
			value:           time.Date(2026, time.July, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:            "New York April fiscal year",
			location:        "America/New_York",
			fiscalYearStart: time.April,
			value:           time.Date(2026, time.July, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:            "Los Angeles October fiscal year",
			location:        "America/Los_Angeles",
			fiscalYearStart: time.October,
			value:           time.Date(2027, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:            "Sydney July fiscal year",
			location:        "Australia/Sydney",
			fiscalYearStart: time.July,
			value:           time.Date(2026, time.December, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location := mustLoadLocation(t, tt.location)

			cal := newTestCalendar(t, location)
			cal.fiscalYearStart = tt.fiscalYearStart

			start := cal.StartOfYear(tt.value)
			next := cal.StartOfNextYear(tt.value)

			if start.Location() != location {
				t.Errorf(
					"StartOfYear() location = %v, want %v",
					start.Location(),
					location,
				)
			}

			if next.Location() != location {
				t.Errorf(
					"StartOfNextYear() location = %v, want %v",
					next.Location(),
					location,
				)
			}

			if start.Hour() != 0 ||
				start.Minute() != 0 ||
				start.Second() != 0 ||
				start.Nanosecond() != 0 {
				t.Errorf(
					"StartOfYear() = %v, want local midnight",
					start,
				)
			}

			if next.Hour() != 0 ||
				next.Minute() != 0 ||
				next.Second() != 0 ||
				next.Nanosecond() != 0 {
				t.Errorf(
					"StartOfNextYear() = %v, want local midnight",
					next,
				)
			}

			if !start.Before(next) {
				t.Errorf(
					"StartOfYear() = %v should be before StartOfNextYear() = %v",
					start,
					next,
				)
			}
		})
	}
}

func TestAddYears(t *testing.T) {
	tests := []struct {
		name  string
		value time.Time
		years int
		want  time.Time
	}{
		{
			name:  "add zero years",
			value: time.Date(2026, time.January, 15, 12, 30, 0, 0, time.UTC),
			years: 0,
			want:  time.Date(2026, time.January, 15, 12, 30, 0, 0, time.UTC),
		},
		{
			name:  "add one year",
			value: time.Date(2026, time.January, 15, 12, 30, 0, 0, time.UTC),
			years: 1,
			want:  time.Date(2027, time.January, 15, 12, 30, 0, 0, time.UTC),
		},
		{
			name:  "subtract one year",
			value: time.Date(2026, time.January, 15, 12, 30, 0, 0, time.UTC),
			years: -1,
			want:  time.Date(2025, time.January, 15, 12, 30, 0, 0, time.UTC),
		},
		{
			name:  "add multiple years",
			value: time.Date(2026, time.January, 15, 12, 30, 0, 0, time.UTC),
			years: 10,
			want:  time.Date(2036, time.January, 15, 12, 30, 0, 0, time.UTC),
		},
		{
			name:  "leap day follows AddDate semantics",
			value: time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC),
			years: 1,
			want:  time.Date(2025, time.March, 1, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddYears(tt.value, tt.years)

			if !got.Equal(tt.want) {
				t.Errorf("AddYears() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddYears_PreservesLocation(t *testing.T) {
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
		12,
		30,
		45,
		123456789,
		time.UTC,
	)

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)
			localizedValue := value.In(location)

			got := AddYears(localizedValue, 1)

			if got.Location() != location {
				t.Errorf(
					"AddYears() location = %v, want %v",
					got.Location(),
					location,
				)
			}
		})
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year int
		want bool
	}{
		{year: 2024, want: true},
		{year: 2028, want: true},
		{year: 2000, want: true},
		{year: 1900, want: false},
		{year: 2100, want: false},
		{year: 2200, want: false},
		{year: 2400, want: true},
		{year: 2023, want: false},
		{year: 2025, want: false},
		{year: 2026, want: false},
		{year: 2027, want: false},
	}

	for _, tt := range tests {
		t.Run(
			string(rune(tt.year)),
			func(t *testing.T) {
				got := IsLeapYear(tt.year)

				if got != tt.want {
					t.Errorf(
						"IsLeapYear(%d) = %v, want %v",
						tt.year,
						got,
						tt.want,
					)
				}
			},
		)
	}
}

func TestDaysInYear(t *testing.T) {
	tests := []struct {
		year int
		want int
	}{
		{year: 2024, want: 366},
		{year: 2028, want: 366},
		{year: 2000, want: 366},
		{year: 1900, want: 365},
		{year: 2100, want: 365},
		{year: 2200, want: 365},
		{year: 2400, want: 366},
		{year: 2023, want: 365},
		{year: 2025, want: 365},
		{year: 2026, want: 365},
	}

	for _, tt := range tests {
		t.Run(
			string(rune(tt.year)),
			func(t *testing.T) {
				got := DaysInYear(tt.year)

				if got != tt.want {
					t.Errorf(
						"DaysInYear(%d) = %d, want %d",
						tt.year,
						got,
						tt.want,
					)
				}
			},
		)
	}
}

func TestDaysInYear_MatchesIsLeapYear(t *testing.T) {
	for year := 1; year <= 3000; year++ {
		got := DaysInYear(year)

		if IsLeapYear(year) && got != 366 {
			t.Errorf(
				"DaysInYear(%d) = %d, want 366 for leap year",
				year,
				got,
			)
		}

		if !IsLeapYear(year) && got != 365 {
			t.Errorf(
				"DaysInYear(%d) = %d, want 365 for common year",
				year,
				got,
			)
		}
	}
}

func TestCalendar_FiscalYearBoundariesAcrossAllStartMonths(t *testing.T) {
	for fiscalStart := time.January; fiscalStart <= time.December; fiscalStart++ {
		t.Run(fiscalStart.String(), func(t *testing.T) {
			cal := newTestCalendar(t, time.UTC)
			cal.fiscalYearStart = fiscalStart

			value := time.Date(
				2026,
				fiscalStart,
				15,
				12,
				0,
				0,
				0,
				time.UTC,
			)

			start := cal.StartOfYear(value)
			next := cal.StartOfNextYear(value)
			previous := cal.StartOfPreviousYear(value)

			if start.Year() != 2026 {
				t.Errorf(
					"StartOfYear() year = %d, want 2026",
					start.Year(),
				)
			}

			if start.Month() != fiscalStart {
				t.Errorf(
					"StartOfYear() month = %v, want %v",
					start.Month(),
					fiscalStart,
				)
			}

			if next.Year() != 2027 {
				t.Errorf(
					"StartOfNextYear() year = %d, want 2027",
					next.Year(),
				)
			}

			if next.Month() != fiscalStart {
				t.Errorf(
					"StartOfNextYear() month = %v, want %v",
					next.Month(),
					fiscalStart,
				)
			}

			if previous.Year() != 2025 {
				t.Errorf(
					"StartOfPreviousYear() year = %d, want 2025",
					previous.Year(),
				)
			}

			if previous.Month() != fiscalStart {
				t.Errorf(
					"StartOfPreviousYear() month = %v, want %v",
					previous.Month(),
					fiscalStart,
				)
			}
		})
	}
}
