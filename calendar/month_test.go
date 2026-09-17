package calendar

import (
	"testing"
	"time"
)

func TestCalendar_StartOfMonth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		input    time.Time
		want     time.Time
	}{
		{
			name:     "UTC beginning of January",
			location: "UTC",
			input:    time.Date(2025, 1, 15, 14, 30, 45, 123, time.UTC),
			want:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "UTC beginning of December",
			location: "UTC",
			input:    time.Date(2025, 12, 31, 23, 59, 59, 999, time.UTC),
			want:     time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "India timezone",
			location: "Asia/Kolkata",
			input:    time.Date(2025, 6, 20, 18, 45, 30, 123, time.UTC),
			want:     time.Date(2025, 6, 1, 0, 0, 0, 0, mustLoadLocation(t, "Asia/Kolkata")),
		},
		{
			name:     "New York timezone",
			location: "America/New_York",
			input:    time.Date(2025, 3, 15, 18, 30, 0, 0, time.UTC),
			want:     time.Date(2025, 3, 1, 0, 0, 0, 0, mustLoadLocation(t, "America/New_York")),
		},
		{
			name:     "month containing leap day",
			location: "UTC",
			input:    time.Date(2024, 2, 29, 23, 59, 59, 0, time.UTC),
			want:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "last second of month",
			location: "UTC",
			input:    time.Date(2025, 4, 30, 23, 59, 59, 999999999, time.UTC),
			want:     time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			got := cal.StartOfMonth(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("StartOfMonth() = %v, want %v", got, tt.want)
			}

			if got.Location() != location {
				t.Fatalf(
					"StartOfMonth() location = %v, want %v",
					got.Location(),
					location,
				)
			}
		})
	}
}

func TestCalendar_StartOfMonth_UsesConfiguredLocation(t *testing.T) {
	t.Parallel()

	location := mustLoadLocation(t, "Asia/Kolkata")

	cal := mustNewCalendar(t, Config{
		Location:        location,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	// This instant is March 1 in India but February 28 in New York.
	value := time.Date(
		2025,
		2,
		28,
		19,
		0,
		0,
		0,
		time.UTC,
	)

	got := cal.StartOfMonth(value)

	want := time.Date(
		2025,
		3,
		1,
		0,
		0,
		0,
		0,
		location,
	)

	if !got.Equal(want) {
		t.Fatalf("StartOfMonth() = %v, want %v", got, want)
	}
}

func TestCalendar_EndOfMonth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		input    time.Time
	}{
		{
			name:     "January UTC",
			location: "UTC",
			input:    time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "February leap year UTC",
			location: "UTC",
			input:    time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "February common year UTC",
			location: "UTC",
			input:    time.Date(2025, 2, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "DST month New York",
			location: "America/New_York",
			input:    time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "DST month Los Angeles",
			location: "America/Los_Angeles",
			input:    time.Date(2025, 11, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			got := cal.EndOfMonth(tt.input)
			next := cal.StartOfNextMonth(tt.input)

			want := next.Add(-time.Nanosecond)

			if !got.Equal(want) {
				t.Fatalf("EndOfMonth() = %v, want %v", got, want)
			}

			if !got.Before(next) {
				t.Fatalf(
					"EndOfMonth() = %v, expected it to be before next month %v",
					got,
					next,
				)
			}

			if got.Add(time.Nanosecond) != next {
				t.Fatalf(
					"EndOfMonth() + 1ns = %v, want %v",
					got.Add(time.Nanosecond),
					next,
				)
			}
		})
	}
}

func TestCalendar_StartOfNextMonth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		input    time.Time
		want     time.Time
	}{
		{
			name:     "January to February",
			location: "UTC",
			input:    time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			want:     time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "December to January",
			location: "UTC",
			input:    time.Date(2025, 12, 15, 12, 0, 0, 0, time.UTC),
			want:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "February leap year",
			location: "UTC",
			input:    time.Date(2024, 2, 29, 23, 59, 59, 0, time.UTC),
			want:     time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "February common year",
			location: "UTC",
			input:    time.Date(2025, 2, 28, 23, 59, 59, 0, time.UTC),
			want:     time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "DST transition month",
			location: "America/New_York",
			input:    time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC),
			want: time.Date(
				2025,
				4,
				1,
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
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			got := cal.StartOfNextMonth(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("StartOfNextMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_StartOfPreviousMonth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		input    time.Time
		want     time.Time
	}{
		{
			name:     "February to January",
			location: "UTC",
			input:    time.Date(2025, 2, 15, 12, 0, 0, 0, time.UTC),
			want:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "January to previous December",
			location: "UTC",
			input:    time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			want:     time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "March leap year",
			location: "UTC",
			input:    time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
			want:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "DST timezone",
			location: "America/New_York",
			input:    time.Date(2025, 4, 15, 12, 0, 0, 0, time.UTC),
			want: time.Date(
				2025,
				3,
				1,
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
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			got := cal.StartOfPreviousMonth(tt.input)

			if !got.Equal(tt.want) {
				t.Fatalf("StartOfPreviousMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_IsSameMonth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		a        time.Time
		b        time.Time
		want     bool
	}{
		{
			name:     "same instant",
			location: "UTC",
			a:        time.Date(2025, 5, 10, 12, 0, 0, 0, time.UTC),
			b:        time.Date(2025, 5, 10, 12, 0, 0, 0, time.UTC),
			want:     true,
		},
		{
			name:     "different days same month",
			location: "UTC",
			a:        time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			b:        time.Date(2025, 5, 31, 23, 59, 59, 0, time.UTC),
			want:     true,
		},
		{
			name:     "different months",
			location: "UTC",
			a:        time.Date(2025, 5, 31, 23, 59, 59, 0, time.UTC),
			b:        time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			want:     false,
		},
		{
			name:     "same month different years",
			location: "UTC",
			a:        time.Date(2025, 5, 15, 12, 0, 0, 0, time.UTC),
			b:        time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC),
			want:     false,
		},
		{
			name:     "different UTC dates same India month",
			location: "Asia/Kolkata",
			a:        time.Date(2025, 3, 31, 20, 0, 0, 0, time.UTC),
			b:        time.Date(2025, 3, 31, 23, 59, 0, 0, time.UTC),
			want:     true,
		},
		{
			name:     "different UTC months same India month",
			location: "Asia/Kolkata",
			a:        time.Date(2025, 2, 28, 23, 0, 0, 0, time.UTC),
			b:        time.Date(2025, 3, 1, 1, 0, 0, 0, time.UTC),
			want:     true,
		},
		{
			name:     "same instant represented in different locations",
			location: "UTC",
			a:        time.Date(2025, 5, 1, 0, 30, 0, 0, time.UTC),
			b: time.Date(
				2025,
				4,
				30,
				20,
				30,
				0,
				0,
				mustLoadLocation(t, "America/New_York"),
			),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			got := cal.IsSameMonth(tt.a, tt.b)

			if got != tt.want {
				t.Fatalf("IsSameMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_DaysInMonth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		input    time.Time
		want     int
	}{
		{
			name:     "January",
			location: "UTC",
			input:    time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			want:     31,
		},
		{
			name:     "February common year",
			location: "UTC",
			input:    time.Date(2025, 2, 15, 12, 0, 0, 0, time.UTC),
			want:     28,
		},
		{
			name:     "February leap year",
			location: "UTC",
			input:    time.Date(2024, 2, 15, 12, 0, 0, 0, time.UTC),
			want:     29,
		},
		{
			name:     "April",
			location: "UTC",
			input:    time.Date(2025, 4, 15, 12, 0, 0, 0, time.UTC),
			want:     30,
		},
		{
			name:     "May",
			location: "UTC",
			input:    time.Date(2025, 5, 15, 12, 0, 0, 0, time.UTC),
			want:     31,
		},
		{
			name:     "June",
			location: "UTC",
			input:    time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC),
			want:     30,
		},
		{
			name:     "September",
			location: "UTC",
			input:    time.Date(2025, 9, 15, 12, 0, 0, 0, time.UTC),
			want:     30,
		},
		{
			name:     "December",
			location: "UTC",
			input:    time.Date(2025, 12, 15, 12, 0, 0, 0, time.UTC),
			want:     31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			got := cal.DaysInMonth(tt.input)

			if got != tt.want {
				t.Fatalf("DaysInMonth() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCalendar_DaysInMonth_AllMonths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		month time.Month
		days  int
	}{
		{time.January, 31},
		{time.February, 28},
		{time.March, 31},
		{time.April, 30},
		{time.May, 31},
		{time.June, 30},
		{time.July, 31},
		{time.August, 31},
		{time.September, 30},
		{time.October, 31},
		{time.November, 30},
		{time.December, 31},
	}

	cal := mustNewCalendar(t, Config{
		Location:        time.UTC,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	for _, tt := range tests {
		t.Run(tt.month.String(), func(t *testing.T) {
			got := cal.DaysInMonth(
				time.Date(2025, tt.month, 15, 12, 0, 0, 0, time.UTC),
			)

			if got != tt.days {
				t.Fatalf(
					"DaysInMonth(%v) = %d, want %d",
					tt.month,
					got,
					tt.days,
				)
			}
		})
	}
}

func TestCalendar_DaysInMonth_LeapYears(t *testing.T) {
	t.Parallel()

	tests := []struct {
		year int
		want int
	}{
		{year: 2000, want: 29},
		{year: 2004, want: 29},
		{year: 2020, want: 29},
		{year: 2024, want: 29},
		{year: 2100, want: 28},
		{year: 2023, want: 28},
		{year: 2025, want: 28},
	}

	cal := mustNewCalendar(t, Config{
		Location:        time.UTC,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	for _, tt := range tests {
		t.Run(time.Date(tt.year, time.February, 1, 0, 0, 0, 0, time.UTC).Format("2006"), func(t *testing.T) {
			got := cal.DaysInMonth(
				time.Date(tt.year, time.February, 15, 12, 0, 0, 0, time.UTC),
			)

			if got != tt.want {
				t.Fatalf(
					"DaysInMonth(%d) = %d, want %d",
					tt.year,
					got,
					tt.want,
				)
			}
		})
	}
}

func TestCalendar_DaysInMonth_DST(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		input    time.Time
		want     int
	}{
		{
			name:     "New York March",
			location: "America/New_York",
			input:    time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC),
			want:     31,
		},
		{
			name:     "New York November",
			location: "America/New_York",
			input:    time.Date(2025, 11, 15, 12, 0, 0, 0, time.UTC),
			want:     30,
		},
		{
			name:     "Los Angeles March",
			location: "America/Los_Angeles",
			input:    time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC),
			want:     31,
		},
		{
			name:     "Los Angeles November",
			location: "America/Los_Angeles",
			input:    time.Date(2025, 11, 15, 12, 0, 0, 0, time.UTC),
			want:     30,
		},
		{
			name:     "India has no DST",
			location: "Asia/Kolkata",
			input:    time.Date(2025, 3, 15, 12, 0, 0, 0, time.UTC),
			want:     31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			got := cal.DaysInMonth(tt.input)

			if got != tt.want {
				t.Fatalf(
					"DaysInMonth() = %d, want %d",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestCalendar_MonthBoundariesAreConsistent(t *testing.T) {
	t.Parallel()

	locations := []string{
		"UTC",
		"Asia/Kolkata",
		"Asia/Tokyo",
		"Europe/London",
		"Europe/Berlin",
		"America/New_York",
		"America/Los_Angeles",
		"America/Chicago",
		"Australia/Sydney",
		"Pacific/Auckland",
	}

	for _, locationName := range locations {
		t.Run(locationName, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, locationName)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			for _, month := range []time.Month{
				time.January,
				time.February,
				time.March,
				time.April,
				time.May,
				time.June,
				time.July,
				time.August,
				time.September,
				time.October,
				time.November,
				time.December,
			} {
				value := time.Date(
					2025,
					month,
					15,
					12,
					30,
					0,
					0,
					location,
				)

				start := cal.StartOfMonth(value)
				next := cal.StartOfNextMonth(value)
				previous := cal.StartOfPreviousMonth(value)
				end := cal.EndOfMonth(value)

				if start.Day() != 1 {
					t.Fatalf(
						"StartOfMonth() day = %d, want 1",
						start.Day(),
					)
				}

				if start.Hour() != 0 ||
					start.Minute() != 0 ||
					start.Second() != 0 ||
					start.Nanosecond() != 0 {
					t.Fatalf(
						"StartOfMonth() = %v, want local midnight",
						start,
					)
				}

				if next.Day() != 1 {
					t.Fatalf(
						"StartOfNextMonth() day = %d, want 1",
						next.Day(),
					)
				}

				if previous.Day() != 1 {
					t.Fatalf(
						"StartOfPreviousMonth() day = %d, want 1",
						previous.Day(),
					)
				}

				if !start.Before(next) {
					t.Fatalf(
						"StartOfMonth() = %v is not before StartOfNextMonth() = %v",
						start,
						next,
					)
				}

				if !previous.Before(start) {
					t.Fatalf(
						"StartOfPreviousMonth() = %v is not before StartOfMonth() = %v",
						previous,
						start,
					)
				}

				if end.Add(time.Nanosecond) != next {
					t.Fatalf(
						"EndOfMonth() + 1ns = %v, want %v",
						end.Add(time.Nanosecond),
						next,
					)
				}
			}
		})
	}
}

func TestAddMonths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		value  time.Time
		months int
		want   time.Time
	}{
		{
			name:   "add zero months",
			value:  time.Date(2025, 5, 15, 12, 30, 45, 0, time.UTC),
			months: 0,
			want:   time.Date(2025, 5, 15, 12, 30, 45, 0, time.UTC),
		},
		{
			name:   "add one month",
			value:  time.Date(2025, 5, 15, 12, 30, 45, 0, time.UTC),
			months: 1,
			want:   time.Date(2025, 6, 15, 12, 30, 45, 0, time.UTC),
		},
		{
			name:   "subtract one month",
			value:  time.Date(2025, 5, 15, 12, 30, 45, 0, time.UTC),
			months: -1,
			want:   time.Date(2025, 4, 15, 12, 30, 45, 0, time.UTC),
		},
		{
			name:   "cross year forward",
			value:  time.Date(2025, 11, 15, 12, 0, 0, 0, time.UTC),
			months: 2,
			want:   time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:   "cross year backward",
			value:  time.Date(2025, 2, 15, 12, 0, 0, 0, time.UTC),
			months: -3,
			want:   time.Date(2024, 11, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:   "add twelve months",
			value:  time.Date(2025, 5, 15, 12, 0, 0, 0, time.UTC),
			months: 12,
			want:   time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:   "subtract twelve months",
			value:  time.Date(2025, 5, 15, 12, 0, 0, 0, time.UTC),
			months: -12,
			want:   time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:   "preserves wall clock time",
			value:  time.Date(2025, 5, 15, 23, 59, 58, 123456789, time.UTC),
			months: 1,
			want:   time.Date(2025, 6, 15, 23, 59, 58, 123456789, time.UTC),
		},
		{
			name: "preserves timezone",
			value: time.Date(
				2025,
				5,
				15,
				12,
				0,
				0,
				0,
				mustLoadLocation(t, "Asia/Kolkata"),
			),
			months: 1,
			want: time.Date(
				2025,
				6,
				15,
				12,
				0,
				0,
				0,
				mustLoadLocation(t, "Asia/Kolkata"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := AddMonths(tt.value, tt.months)

			if !got.Equal(tt.want) {
				t.Fatalf("AddMonths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddMonths_CalendarMonthArithmetic(t *testing.T) {
	t.Parallel()

	// These cases intentionally use dates near month ends.
	// AddDate follows Go's calendar arithmetic semantics rather than
	// clamping to the final day of the destination month.
	tests := []struct {
		name   string
		value  time.Time
		months int
	}{
		{
			name:   "January 31 plus one month",
			value:  time.Date(2025, 1, 31, 12, 0, 0, 0, time.UTC),
			months: 1,
		},
		{
			name:   "March 31 minus one month",
			value:  time.Date(2025, 3, 31, 12, 0, 0, 0, time.UTC),
			months: -1,
		},
		{
			name:   "leap day plus one year",
			value:  time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			months: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := AddMonths(tt.value, tt.months)
			want := tt.value.AddDate(0, tt.months, 0)

			if !got.Equal(want) {
				t.Fatalf("AddMonths() = %v, want %v", got, want)
			}
		})
	}
}

func TestCalendar_MonthOperations_PreserveConfiguredLocation(t *testing.T) {
	t.Parallel()

	locations := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, locationName := range locations {
		t.Run(locationName, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, locationName)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			value := time.Date(
				2025,
				6,
				15,
				12,
				30,
				0,
				0,
				time.UTC,
			)

			start := cal.StartOfMonth(value)
			next := cal.StartOfNextMonth(value)
			previous := cal.StartOfPreviousMonth(value)
			end := cal.EndOfMonth(value)

			for name, got := range map[string]time.Time{
				"StartOfMonth":         start,
				"StartOfNextMonth":     next,
				"StartOfPreviousMonth": previous,
				"EndOfMonth":           end,
			} {
				if got.Location() != location {
					t.Errorf(
						"%s location = %v, want %v",
						name,
						got.Location(),
						location,
					)
				}
			}
		})
	}
}

func TestCalendar_MonthHalfOpenInterval(t *testing.T) {
	t.Parallel()

	locations := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"America/Los_Angeles",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, locationName := range locations {
		t.Run(locationName, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, locationName)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			value := time.Date(
				2025,
				6,
				15,
				12,
				0,
				0,
				0,
				location,
			)

			start := cal.StartOfMonth(value)
			next := cal.StartOfNextMonth(value)
			end := cal.EndOfMonth(value)

			if !start.Before(next) {
				t.Fatalf(
					"month interval is invalid: start=%v next=%v",
					start,
					next,
				)
			}

			if !end.Before(next) {
				t.Fatalf(
					"EndOfMonth() = %v should be before StartOfNextMonth() = %v",
					end,
					next,
				)
			}

			if end.Add(time.Nanosecond) != next {
				t.Fatalf(
					"EndOfMonth() + 1ns = %v, want %v",
					end.Add(time.Nanosecond),
					next,
				)
			}
		})
	}
}

func TestCalendar_MonthBoundaryAtExactTransition(t *testing.T) {
	t.Parallel()

	location := mustLoadLocation(t, "Asia/Kolkata")

	cal := mustNewCalendar(t, Config{
		Location:        location,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	tests := []struct {
		name  string
		input time.Time
	}{
		{
			name:  "exact beginning of month",
			input: time.Date(2025, 5, 1, 0, 0, 0, 0, location),
		},
		{
			name: "one nanosecond before month end",
			input: time.Date(
				2025,
				5,
				31,
				23,
				59,
				59,
				999999999,
				location,
			),
		},
		{
			name:  "exact beginning of next month",
			input: time.Date(2025, 6, 1, 0, 0, 0, 0, location),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			start := cal.StartOfMonth(tt.input)
			next := cal.StartOfNextMonth(tt.input)
			end := cal.EndOfMonth(tt.input)

			if !start.Before(next) {
				t.Fatalf(
					"StartOfMonth() = %v, StartOfNextMonth() = %v",
					start,
					next,
				)
			}

			if end.Add(time.Nanosecond) != next {
				t.Fatalf(
					"EndOfMonth() + 1ns = %v, want %v",
					end.Add(time.Nanosecond),
					next,
				)
			}
		})
	}
}

func TestCalendar_MonthOperations_AcrossYearBoundary(t *testing.T) {
	t.Parallel()

	cal := mustNewCalendar(t, Config{
		Location:        time.UTC,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	tests := []struct {
		name     string
		input    time.Time
		wantNext time.Time
		wantPrev time.Time
	}{
		{
			name:     "January",
			input:    time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC),
			wantNext: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			wantPrev: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "December",
			input:    time.Date(2025, 12, 15, 12, 0, 0, 0, time.UTC),
			wantNext: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			wantPrev: time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			next := cal.StartOfNextMonth(tt.input)
			prev := cal.StartOfPreviousMonth(tt.input)

			if !next.Equal(tt.wantNext) {
				t.Fatalf(
					"StartOfNextMonth() = %v, want %v",
					next,
					tt.wantNext,
				)
			}

			if !prev.Equal(tt.wantPrev) {
				t.Fatalf(
					"StartOfPreviousMonth() = %v, want %v",
					prev,
					tt.wantPrev,
				)
			}
		})
	}
}

func TestCalendar_MonthOperations_AtDifferentInputLocations(t *testing.T) {
	t.Parallel()

	location := mustLoadLocation(t, "Asia/Kolkata")

	cal := mustNewCalendar(t, Config{
		Location:        location,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	inputs := []time.Time{
		time.Date(2025, 5, 15, 12, 0, 0, 0, time.UTC),
		time.Date(
			2025,
			5,
			15,
			12,
			0,
			0,
			0,
			mustLoadLocation(t, "America/New_York"),
		),
		time.Date(
			2025,
			5,
			15,
			12,
			0,
			0,
			0,
			mustLoadLocation(t, "Asia/Tokyo"),
		),
	}

	for _, input := range inputs {
		t.Run(input.Location().String(), func(t *testing.T) {
			t.Parallel()

			got := cal.StartOfMonth(input)

			if got.Location() != location {
				t.Fatalf(
					"StartOfMonth() location = %v, want %v",
					got.Location(),
					location,
				)
			}

			if got.Day() != 1 {
				t.Fatalf(
					"StartOfMonth() day = %d, want 1",
					got.Day(),
				)
			}

			if got.Hour() != 0 ||
				got.Minute() != 0 ||
				got.Second() != 0 ||
				got.Nanosecond() != 0 {
				t.Fatalf(
					"StartOfMonth() = %v, want midnight",
					got,
				)
			}
		})
	}
}
