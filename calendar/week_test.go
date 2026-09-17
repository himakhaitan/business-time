package calendar

import (
	"testing"
	"time"
)

func TestCalendar_Weekday(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		value    time.Time
		want     time.Weekday
	}{
		{
			name:     "UTC Monday",
			location: "UTC",
			value:    time.Date(2026, time.March, 16, 12, 0, 0, 0, time.UTC),
			want:     time.Monday,
		},
		{
			name:     "UTC Sunday",
			location: "UTC",
			value:    time.Date(2026, time.March, 15, 12, 0, 0, 0, time.UTC),
			want:     time.Sunday,
		},
		{
			name:     "Kolkata crosses UTC date boundary",
			location: "Asia/Kolkata",
			value:    time.Date(2026, time.March, 15, 20, 0, 0, 0, time.UTC),
			want:     time.Monday,
		},
		{
			name:     "Tokyo crosses UTC date boundary",
			location: "Asia/Tokyo",
			value:    time.Date(2026, time.March, 15, 16, 0, 0, 0, time.UTC),
			want:     time.Monday,
		},
		{
			name:     "New York remains previous local date",
			location: "America/New_York",
			value:    time.Date(2026, time.March, 16, 2, 0, 0, 0, time.UTC),
			want:     time.Sunday,
		},
		{
			name:     "Los Angeles remains previous local date",
			location: "America/Los_Angeles",
			value:    time.Date(2026, time.March, 16, 5, 0, 0, 0, time.UTC),
			want:     time.Sunday,
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

			got := cal.Weekday(tt.value)

			if got != tt.want {
				t.Fatalf("Weekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_Weekday_UsesCalendarLocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		value    time.Time
		want     time.Weekday
	}{
		{
			name:     "UTC to Kolkata",
			location: "Asia/Kolkata",
			value:    time.Date(2026, time.March, 15, 20, 0, 0, 0, time.UTC),
			want:     time.Monday,
		},
		{
			name:     "UTC to New York",
			location: "America/New_York",
			value:    time.Date(2026, time.March, 16, 2, 0, 0, 0, time.UTC),
			want:     time.Sunday,
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

			got := cal.Weekday(tt.value)

			if got != tt.want {
				t.Fatalf("Weekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_StartOfWeek(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		location  string
		weekStart time.Weekday
		value     time.Time
		want      time.Time
	}{
		{
			name:      "Monday start from Monday",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.March, 16, 14, 30, 0, 0, time.UTC),
			want:      time.Date(2026, time.March, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Monday start from Wednesday",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.March, 18, 14, 30, 0, 0, time.UTC),
			want:      time.Date(2026, time.March, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Monday start from Sunday",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.March, 22, 14, 30, 0, 0, time.UTC),
			want:      time.Date(2026, time.March, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Sunday start from Sunday",
			location:  "UTC",
			weekStart: time.Sunday,
			value:     time.Date(2026, time.March, 22, 14, 30, 0, 0, time.UTC),
			want:      time.Date(2026, time.March, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Sunday start from Wednesday",
			location:  "UTC",
			weekStart: time.Sunday,
			value:     time.Date(2026, time.March, 25, 14, 30, 0, 0, time.UTC),
			want:      time.Date(2026, time.March, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Saturday start from Friday",
			location:  "UTC",
			weekStart: time.Saturday,
			value:     time.Date(2026, time.March, 27, 14, 30, 0, 0, time.UTC),
			want:      time.Date(2026, time.March, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Saturday start from Saturday",
			location:  "UTC",
			weekStart: time.Saturday,
			value:     time.Date(2026, time.March, 28, 14, 30, 0, 0, time.UTC),
			want:      time.Date(2026, time.March, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "month boundary",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.April, 1, 12, 0, 0, 0, time.UTC),
			want:      time.Date(2026, time.March, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "year boundary",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2027, time.January, 1, 12, 0, 0, 0, time.UTC),
			want:      time.Date(2026, time.December, 28, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       tt.weekStart,
				FiscalYearStart: time.January,
			})

			got := cal.StartOfWeek(tt.value)

			if !got.Equal(tt.want) {
				t.Fatalf("StartOfWeek() = %v, want %v", got, tt.want)
			}

			if got.Location() != tt.want.Location() {
				t.Fatalf(
					"StartOfWeek() location = %v, want %v",
					got.Location(),
					tt.want.Location(),
				)
			}
		})
	}
}

func TestCalendar_StartOfWeek_AllWeekStarts(t *testing.T) {
	t.Parallel()

	// Wednesday, March 18, 2026.
	value := time.Date(2026, time.March, 18, 15, 30, 0, 0, time.UTC)

	tests := []struct {
		weekStart time.Weekday
		wantDay   int
		want      time.Weekday
	}{
		{
			weekStart: time.Sunday,
			wantDay:   15,
			want:      time.Sunday,
		},
		{
			weekStart: time.Monday,
			wantDay:   16,
			want:      time.Monday,
		},
		{
			weekStart: time.Tuesday,
			wantDay:   17,
			want:      time.Tuesday,
		},
		{
			weekStart: time.Wednesday,
			wantDay:   18,
			want:      time.Wednesday,
		},
		{
			weekStart: time.Thursday,
			wantDay:   12,
			want:      time.Thursday,
		},
		{
			weekStart: time.Friday,
			wantDay:   13,
			want:      time.Friday,
		},
		{
			weekStart: time.Saturday,
			wantDay:   14,
			want:      time.Saturday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.weekStart.String(), func(t *testing.T) {
			t.Parallel()

			cal := mustNewCalendar(t, Config{
				Location:        time.UTC,
				WeekStart:       tt.weekStart,
				FiscalYearStart: time.January,
			})

			got := cal.StartOfWeek(value)

			if got.Day() != tt.wantDay {
				t.Fatalf("StartOfWeek() day = %d, want %d", got.Day(), tt.wantDay)
			}

			if got.Weekday() != tt.want {
				t.Fatalf(
					"StartOfWeek() weekday = %v, want %v",
					got.Weekday(),
					tt.want,
				)
			}

			if got.Hour() != 0 ||
				got.Minute() != 0 ||
				got.Second() != 0 ||
				got.Nanosecond() != 0 {
				t.Fatalf("StartOfWeek() = %v, want midnight", got)
			}
		})
	}
}

func TestCalendar_StartOfWeek_UsesCalendarLocation(t *testing.T) {
	t.Parallel()

	location := mustLoadLocation(t, "Asia/Kolkata")

	cal := mustNewCalendar(t, Config{
		Location:        location,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	// Sunday 20:00 UTC = Monday 01:30 IST.
	value := time.Date(2026, time.March, 15, 20, 0, 0, 0, time.UTC)

	got := cal.StartOfWeek(value)

	want := time.Date(
		2026,
		time.March,
		16,
		0, 0, 0, 0,
		location,
	)

	if !got.Equal(want) {
		t.Fatalf("StartOfWeek() = %v, want %v", got, want)
	}
}

func TestCalendar_StartOfWeek_IsMidnight(t *testing.T) {
	t.Parallel()

	locations := []string{
		"UTC",
		"Asia/Kolkata",
		"Asia/Tokyo",
		"Europe/London",
		"Europe/Berlin",
		"America/New_York",
		"America/Los_Angeles",
		"Australia/Sydney",
		"Pacific/Auckland",
	}

	value := time.Date(2026, time.July, 15, 18, 42, 33, 123456789, time.UTC)

	for _, name := range locations {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, name)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			got := cal.StartOfWeek(value)

			if got.Hour() != 0 ||
				got.Minute() != 0 ||
				got.Second() != 0 ||
				got.Nanosecond() != 0 {
				t.Fatalf("StartOfWeek() = %v, want midnight", got)
			}

			if got.Location() != location {
				t.Fatalf(
					"StartOfWeek() location = %v, want %v",
					got.Location(),
					location,
				)
			}
		})
	}
}

func TestCalendar_EndOfWeek(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		location  string
		weekStart time.Weekday
		value     time.Time
	}{
		{
			name:      "normal week",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.March, 18, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "month boundary",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.March, 30, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "year boundary",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2027, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "DST spring forward",
			location:  "America/New_York",
			weekStart: time.Monday,
			value:     time.Date(2026, time.March, 11, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "DST fall back",
			location:  "America/New_York",
			weekStart: time.Monday,
			value:     time.Date(2026, time.November, 4, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       tt.weekStart,
				FiscalYearStart: time.January,
			})

			got := cal.EndOfWeek(tt.value)
			next := cal.StartOfNextWeek(tt.value)

			want := next.Add(-time.Nanosecond)

			if !got.Equal(want) {
				t.Fatalf("EndOfWeek() = %v, want %v", got, want)
			}

			if !got.Before(next) {
				t.Fatalf("EndOfWeek() = %v, want before %v", got, next)
			}

			if !got.Add(time.Nanosecond).Equal(next) {
				t.Fatalf(
					"EndOfWeek() + 1ns = %v, want %v",
					got.Add(time.Nanosecond),
					next,
				)
			}
		})
	}
}

func TestCalendar_StartOfNextWeek(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		location  string
		weekStart time.Weekday
		value     time.Time
	}{
		{
			name:      "normal week",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.March, 18, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "month transition",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.March, 30, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "year transition",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.December, 30, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       tt.weekStart,
				FiscalYearStart: time.January,
			})

			got := cal.StartOfNextWeek(tt.value)
			current := cal.StartOfWeek(tt.value)

			want := current.AddDate(0, 0, 7)

			if !got.Equal(want) {
				t.Fatalf("StartOfNextWeek() = %v, want %v", got, want)
			}

			if !got.After(current) {
				t.Fatalf(
					"StartOfNextWeek() = %v, want after %v",
					got,
					current,
				)
			}

			if got.Weekday() != tt.weekStart {
				t.Fatalf(
					"StartOfNextWeek() weekday = %v, want %v",
					got.Weekday(),
					tt.weekStart,
				)
			}
		})
	}
}

func TestCalendar_StartOfPreviousWeek(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		location  string
		weekStart time.Weekday
		value     time.Time
	}{
		{
			name:      "normal week",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.March, 18, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "month transition",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2026, time.April, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:      "year transition",
			location:  "UTC",
			weekStart: time.Monday,
			value:     time.Date(2027, time.January, 2, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, tt.location)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       tt.weekStart,
				FiscalYearStart: time.January,
			})

			got := cal.StartOfPreviousWeek(tt.value)
			current := cal.StartOfWeek(tt.value)

			want := current.AddDate(0, 0, -7)

			if !got.Equal(want) {
				t.Fatalf("StartOfPreviousWeek() = %v, want %v", got, want)
			}

			if !got.Before(current) {
				t.Fatalf(
					"StartOfPreviousWeek() = %v, want before %v",
					got,
					current,
				)
			}

			if got.Weekday() != tt.weekStart {
				t.Fatalf(
					"StartOfPreviousWeek() weekday = %v, want %v",
					got.Weekday(),
					tt.weekStart,
				)
			}
		})
	}
}

func TestCalendar_IsSameWeek(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		location  string
		weekStart time.Weekday
		a         time.Time
		b         time.Time
		want      bool
	}{
		{
			name:      "same instant",
			location:  "UTC",
			weekStart: time.Monday,
			a:         time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC),
			b:         time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC),
			want:      true,
		},
		{
			name:      "same week different days",
			location:  "UTC",
			weekStart: time.Monday,
			a:         time.Date(2026, time.March, 16, 10, 0, 0, 0, time.UTC),
			b:         time.Date(2026, time.March, 22, 20, 0, 0, 0, time.UTC),
			want:      true,
		},
		{
			name:      "different weeks",
			location:  "UTC",
			weekStart: time.Monday,
			a:         time.Date(2026, time.March, 22, 23, 59, 59, 0, time.UTC),
			b:         time.Date(2026, time.March, 23, 0, 0, 0, 0, time.UTC),
			want:      false,
		},
		{
			name:      "Sunday start makes Sunday same week as Monday",
			location:  "UTC",
			weekStart: time.Sunday,
			a:         time.Date(2026, time.March, 22, 10, 0, 0, 0, time.UTC),
			b:         time.Date(2026, time.March, 23, 10, 0, 0, 0, time.UTC),
			want:      true,
		},
		{
			name:      "Monday start makes Sunday different from Monday",
			location:  "UTC",
			weekStart: time.Monday,
			a:         time.Date(2026, time.March, 22, 10, 0, 0, 0, time.UTC),
			b:         time.Date(2026, time.March, 23, 10, 0, 0, 0, time.UTC),
			want:      false,
		},
		{
			name:      "same instant represented in different locations",
			location:  "Asia/Kolkata",
			weekStart: time.Monday,
			a:         time.Date(2026, time.March, 15, 20, 0, 0, 0, time.UTC),
			b: time.Date(
				2026, time.March, 16, 1, 30, 0, 0,
				mustLoadLocation(t, "Asia/Kolkata"),
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
				WeekStart:       tt.weekStart,
				FiscalYearStart: time.January,
			})

			got := cal.IsSameWeek(tt.a, tt.b)

			if got != tt.want {
				t.Fatalf("IsSameWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalendar_IsSameWeek_UsesCalendarLocation(t *testing.T) {
	t.Parallel()

	instant := time.Date(2026, time.March, 15, 23, 30, 0, 0, time.UTC)

	kolkata := mustLoadLocation(t, "Asia/Kolkata")
	newYork := mustLoadLocation(t, "America/New_York")

	tests := []struct {
		name     string
		location *time.Location
	}{
		{
			name:     "Kolkata",
			location: kolkata,
		},
		{
			name:     "New York",
			location: newYork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cal := mustNewCalendar(t, Config{
				Location:        tt.location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			a := instant
			b := instant.In(tt.location)

			if !cal.IsSameWeek(a, b) {
				t.Fatalf("IsSameWeek() = false, want true for same instant")
			}
		})
	}
}

func TestCalendar_WeekBoundaries(t *testing.T) {
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

	for _, name := range locations {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, name)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			value := time.Date(
				2026,
				time.July,
				15,
				17, 42, 33,
				123456789,
				time.UTC,
			)

			start := cal.StartOfWeek(value)
			next := cal.StartOfNextWeek(value)
			end := cal.EndOfWeek(value)
			previous := cal.StartOfPreviousWeek(value)

			if !start.Before(next) {
				t.Fatalf("start %v should be before next %v", start, next)
			}

			if !previous.Before(start) {
				t.Fatalf(
					"previous %v should be before start %v",
					previous,
					start,
				)
			}

			if !end.Before(next) {
				t.Fatalf("end %v should be before next %v", end, next)
			}

			if !end.Add(time.Nanosecond).Equal(next) {
				t.Fatalf(
					"end + 1ns = %v, want %v",
					end.Add(time.Nanosecond),
					next,
				)
			}

			if start.Weekday() != time.Monday {
				t.Fatalf(
					"start weekday = %v, want Monday",
					start.Weekday(),
				)
			}

			if next.Weekday() != time.Monday {
				t.Fatalf(
					"next weekday = %v, want Monday",
					next.Weekday(),
				)
			}

			if start.Hour() != 0 ||
				start.Minute() != 0 ||
				start.Second() != 0 ||
				start.Nanosecond() != 0 {
				t.Fatalf("start %v is not midnight", start)
			}
		})
	}
}

func TestCalendar_WeekBoundary_Invariants(t *testing.T) {
	t.Parallel()

	location := mustLoadLocation(t, "America/New_York")

	cal := mustNewCalendar(t, Config{
		Location:        location,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	tests := []time.Time{
		time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, time.February, 28, 23, 59, 59, 999999999, time.UTC),
		time.Date(2026, time.March, 8, 6, 59, 59, 0, time.UTC),
		time.Date(2026, time.March, 8, 7, 0, 0, 0, time.UTC),
		time.Date(2026, time.March, 15, 12, 0, 0, 0, time.UTC),
		time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC),
		time.Date(2026, time.November, 1, 5, 59, 59, 0, time.UTC),
		time.Date(2026, time.November, 1, 6, 0, 0, 0, time.UTC),
		time.Date(2026, time.December, 31, 23, 59, 59, 999999999, time.UTC),
	}

	for _, value := range tests {
		value := value

		t.Run(value.Format(time.RFC3339Nano), func(t *testing.T) {
			t.Parallel()

			start := cal.StartOfWeek(value)
			next := cal.StartOfNextWeek(value)
			end := cal.EndOfWeek(value)

			if !start.Before(next) {
				t.Fatalf("start %v is not before next %v", start, next)
			}

			if !start.Equal(cal.StartOfWeek(start)) {
				t.Fatalf("StartOfWeek(start) changed start")
			}

			if !next.Equal(cal.StartOfNextWeek(start)) {
				t.Fatalf("StartOfNextWeek(start) is inconsistent")
			}

			if !end.Equal(next.Add(-time.Nanosecond)) {
				t.Fatalf("EndOfWeek() invariant failed")
			}

			if !value.Before(start) && !value.Before(next) {
				// This is intentionally not used as the actual invariant below;
				// time.Time comparison uses absolute instants, which is correct.
			}

			if value.Before(start) || !value.Before(next) {
				t.Fatalf(
					"value %v is not inside [%v, %v)",
					value,
					start,
					next,
				)
			}
		})
	}
}

func TestCalendar_StartOfWeek_Idempotent(t *testing.T) {
	t.Parallel()

	locations := []string{
		"UTC",
		"Asia/Kolkata",
		"Europe/London",
		"America/New_York",
		"Australia/Sydney",
	}

	for _, name := range locations {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, name)

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			values := []time.Time{
				time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC),
				time.Date(2026, time.March, 8, 12, 0, 0, 0, time.UTC),
				time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
				time.Date(2026, time.November, 1, 12, 0, 0, 0, time.UTC),
				time.Date(2026, time.December, 31, 12, 0, 0, 0, time.UTC),
			}

			for _, value := range values {
				start := cal.StartOfWeek(value)
				got := cal.StartOfWeek(start)

				if !got.Equal(start) {
					t.Fatalf(
						"StartOfWeek(StartOfWeek(%v)) = %v, want %v",
						value,
						got,
						start,
					)
				}
			}
		})
	}
}

func TestCalendar_WeekStartRoundTrip(t *testing.T) {
	t.Parallel()

	location := time.UTC

	for _, weekStart := range []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	} {
		weekStart := weekStart

		t.Run(weekStart.String(), func(t *testing.T) {
			t.Parallel()

			cal := mustNewCalendar(t, Config{
				Location:        location,
				WeekStart:       weekStart,
				FiscalYearStart: time.January,
			})

			for day := time.Sunday; day <= time.Saturday; day++ {
				value := time.Date(
					2026,
					time.March,
					15+int(day),
					12, 0, 0, 0,
					location,
				)

				start := cal.StartOfWeek(value)

				if start.Weekday() != weekStart {
					t.Fatalf(
						"StartOfWeek(%v) = %v, weekday = %v, want %v",
						value,
						start,
						start.Weekday(),
						weekStart,
					)
				}

				next := cal.StartOfNextWeek(value)

				if next.Weekday() != weekStart {
					t.Fatalf(
						"StartOfNextWeek(%v) weekday = %v, want %v",
						value,
						next.Weekday(),
						weekStart,
					)
				}
			}
		})
	}
}

func TestCalendar_StartOfNextWeek_AcrossDST(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		value    time.Time
	}{
		{
			name:     "New York spring transition",
			location: "America/New_York",
			value:    time.Date(2026, time.March, 11, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "New York fall transition",
			location: "America/New_York",
			value:    time.Date(2026, time.November, 4, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Los Angeles spring transition",
			location: "America/Los_Angeles",
			value:    time.Date(2026, time.March, 11, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Los Angeles fall transition",
			location: "America/Los_Angeles",
			value:    time.Date(2026, time.November, 4, 12, 0, 0, 0, time.UTC),
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

			start := cal.StartOfWeek(tt.value)
			next := cal.StartOfNextWeek(tt.value)

			want := start.AddDate(0, 0, 7)

			if !next.Equal(want) {
				t.Fatalf(
					"StartOfNextWeek() = %v, want %v",
					next,
					want,
				)
			}

			if next.Weekday() != time.Monday {
				t.Fatalf(
					"next weekday = %v, want Monday",
					next.Weekday(),
				)
			}
		})
	}
}

func TestCalendar_WeekDuration_CanDifferFrom168Hours(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		location  string
		value     time.Time
		wantHours int64
	}{
		{
			name:      "New York spring DST",
			location:  "America/New_York",
			value:     time.Date(2026, time.March, 4, 12, 0, 0, 0, time.UTC),
			wantHours: 167,
		},
		{
			name:      "New York fall DST",
			location:  "America/New_York",
			value:     time.Date(2026, time.October, 28, 12, 0, 0, 0, time.UTC),
			wantHours: 169,
		},
		{
			name:      "UTC normal week",
			location:  "UTC",
			value:     time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
			wantHours: 168,
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

			start := cal.StartOfWeek(tt.value)
			next := cal.StartOfNextWeek(tt.value)

			gotHours := int64(next.Sub(start) / time.Hour)

			if gotHours != tt.wantHours {
				t.Fatalf(
					"week duration = %d hours, want %d",
					gotHours,
					tt.wantHours,
				)
			}
		})
	}
}

func TestAddWeeks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value time.Time
		weeks int
		want  time.Time
	}{
		{
			name:  "zero",
			value: time.Date(2026, time.March, 18, 12, 30, 45, 123, time.UTC),
			weeks: 0,
			want:  time.Date(2026, time.March, 18, 12, 30, 45, 123, time.UTC),
		},
		{
			name:  "one week forward",
			value: time.Date(2026, time.March, 18, 12, 30, 45, 123, time.UTC),
			weeks: 1,
			want:  time.Date(2026, time.March, 25, 12, 30, 45, 123, time.UTC),
		},
		{
			name:  "one week backward",
			value: time.Date(2026, time.March, 18, 12, 30, 45, 123, time.UTC),
			weeks: -1,
			want:  time.Date(2026, time.March, 11, 12, 30, 45, 123, time.UTC),
		},
		{
			name:  "multiple weeks forward",
			value: time.Date(2026, time.March, 18, 12, 30, 45, 123, time.UTC),
			weeks: 10,
			want:  time.Date(2026, time.May, 27, 12, 30, 45, 123, time.UTC),
		},
		{
			name:  "multiple weeks backward",
			value: time.Date(2026, time.March, 18, 12, 30, 45, 123, time.UTC),
			weeks: -10,
			want:  time.Date(2026, time.January, 7, 12, 30, 45, 123, time.UTC),
		},
		{
			name:  "cross month",
			value: time.Date(2026, time.March, 30, 12, 0, 0, 0, time.UTC),
			weeks: 1,
			want:  time.Date(2026, time.April, 6, 12, 0, 0, 0, time.UTC),
		},
		{
			name:  "cross year",
			value: time.Date(2026, time.December, 28, 12, 0, 0, 0, time.UTC),
			weeks: 1,
			want:  time.Date(2027, time.January, 4, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := AddWeeks(tt.value, tt.weeks)

			if !got.Equal(tt.want) {
				t.Fatalf("AddWeeks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddWeeks_PreservesCalendarLocation(t *testing.T) {
	t.Parallel()

	locations := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, name := range locations {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			location := mustLoadLocation(t, name)

			value := time.Date(
				2026,
				time.July,
				15,
				14, 30, 0, 0,
				location,
			)

			got := AddWeeks(value, 3)

			if got.Location() != location {
				t.Fatalf(
					"AddWeeks() location = %v, want %v",
					got.Location(),
					location,
				)
			}

			if got.Hour() != value.Hour() ||
				got.Minute() != value.Minute() ||
				got.Second() != value.Second() ||
				got.Nanosecond() != value.Nanosecond() {
				t.Fatalf(
					"AddWeeks() changed wall-clock time: got %v, want time %v",
					got,
					value,
				)
			}
		})
	}
}

func TestAddWeeks_AcrossDST(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		value    time.Time
		weeks    int
	}{
		{
			name:     "spring forward",
			location: "America/New_York",
			value:    time.Date(2026, time.March, 8, 12, 30, 0, 0, mustLoadLocation(t, "America/New_York")),
			weeks:    1,
		},
		{
			name:     "fall back",
			location: "America/New_York",
			value:    time.Date(2026, time.November, 1, 12, 30, 0, 0, mustLoadLocation(t, "America/New_York")),
			weeks:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := AddWeeks(tt.value, tt.weeks)

			want := tt.value.AddDate(0, 0, tt.weeks*7)

			if !got.Equal(want) {
				t.Fatalf(
					"AddWeeks() = %v, want %v",
					got,
					want,
				)
			}
		})
	}
}

func TestCalendar_WeekMethods_AtExactBoundaries(t *testing.T) {
	t.Parallel()

	location := time.UTC

	cal := mustNewCalendar(t, Config{
		Location:        location,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	start := time.Date(
		2026,
		time.March,
		16,
		0, 0, 0, 0,
		location,
	)

	end := time.Date(
		2026,
		time.March,
		22,
		23, 59, 59, 999999999,
		location,
	)

	next := time.Date(
		2026,
		time.March,
		23,
		0, 0, 0, 0,
		location,
	)

	tests := []struct {
		name  string
		value time.Time
	}{
		{
			name:  "exact start",
			value: start,
		},
		{
			name:  "exact end",
			value: end,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotStart := cal.StartOfWeek(tt.value)
			gotNext := cal.StartOfNextWeek(tt.value)
			gotEnd := cal.EndOfWeek(tt.value)

			if !gotStart.Equal(start) {
				t.Fatalf("StartOfWeek() = %v, want %v", gotStart, start)
			}

			if !gotNext.Equal(next) {
				t.Fatalf(
					"StartOfNextWeek() = %v, want %v",
					gotNext,
					next,
				)
			}

			if !gotEnd.Equal(end) {
				t.Fatalf("EndOfWeek() = %v, want %v", gotEnd, end)
			}
		})
	}
}

func TestCalendar_WeekMethods_AtExclusiveBoundary(t *testing.T) {
	t.Parallel()

	location := time.UTC

	cal := mustNewCalendar(t, Config{
		Location:        location,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	nextWeek := time.Date(
		2026,
		time.March,
		23,
		0, 0, 0, 0,
		location,
	)

	got := cal.StartOfWeek(nextWeek)

	if !got.Equal(nextWeek) {
		t.Fatalf("StartOfWeek() = %v, want %v", got, nextWeek)
	}

	previous := cal.StartOfPreviousWeek(nextWeek)

	wantPrevious := time.Date(
		2026,
		time.March,
		16,
		0, 0, 0, 0,
		location,
	)

	if !previous.Equal(wantPrevious) {
		t.Fatalf(
			"StartOfPreviousWeek() = %v, want %v",
			previous,
			wantPrevious,
		)
	}
}

func TestCalendar_WeekMethods_AcrossMonthBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     time.Time
		weekStart time.Weekday
		wantStart time.Time
		wantNext  time.Time
	}{
		{
			name:      "week starts in previous month",
			value:     time.Date(2026, time.April, 1, 12, 0, 0, 0, time.UTC),
			weekStart: time.Monday,
			wantStart: time.Date(2026, time.March, 30, 0, 0, 0, 0, time.UTC),
			wantNext:  time.Date(2026, time.April, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "week ends in next month",
			value:     time.Date(2026, time.April, 30, 12, 0, 0, 0, time.UTC),
			weekStart: time.Monday,
			wantStart: time.Date(2026, time.April, 27, 0, 0, 0, 0, time.UTC),
			wantNext:  time.Date(2026, time.May, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Sunday based week crossing month",
			value:     time.Date(2026, time.May, 1, 12, 0, 0, 0, time.UTC),
			weekStart: time.Sunday,
			wantStart: time.Date(2026, time.April, 26, 0, 0, 0, 0, time.UTC),
			wantNext:  time.Date(2026, time.May, 3, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cal := mustNewCalendar(t, Config{
				Location:        time.UTC,
				WeekStart:       tt.weekStart,
				FiscalYearStart: time.January,
			})

			gotStart := cal.StartOfWeek(tt.value)
			gotNext := cal.StartOfNextWeek(tt.value)

			if !gotStart.Equal(tt.wantStart) {
				t.Fatalf(
					"StartOfWeek() = %v, want %v",
					gotStart,
					tt.wantStart,
				)
			}

			if !gotNext.Equal(tt.wantNext) {
				t.Fatalf(
					"StartOfNextWeek() = %v, want %v",
					gotNext,
					tt.wantNext,
				)
			}
		})
	}
}

func TestCalendar_WeekMethods_AcrossYearBoundaries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     time.Time
		weekStart time.Weekday
		wantStart time.Time
		wantNext  time.Time
	}{
		{
			name:      "week begins in previous year",
			value:     time.Date(2027, time.January, 1, 12, 0, 0, 0, time.UTC),
			weekStart: time.Monday,
			wantStart: time.Date(2026, time.December, 28, 0, 0, 0, 0, time.UTC),
			wantNext:  time.Date(2027, time.January, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "week crosses year boundary",
			value:     time.Date(2026, time.December, 31, 12, 0, 0, 0, time.UTC),
			weekStart: time.Sunday,
			wantStart: time.Date(2026, time.December, 27, 0, 0, 0, 0, time.UTC),
			wantNext:  time.Date(2027, time.January, 3, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cal := mustNewCalendar(t, Config{
				Location:        time.UTC,
				WeekStart:       tt.weekStart,
				FiscalYearStart: time.January,
			})

			gotStart := cal.StartOfWeek(tt.value)
			gotNext := cal.StartOfNextWeek(tt.value)

			if !gotStart.Equal(tt.wantStart) {
				t.Fatalf(
					"StartOfWeek() = %v, want %v",
					gotStart,
					tt.wantStart,
				)
			}

			if !gotNext.Equal(tt.wantNext) {
				t.Fatalf(
					"StartOfNextWeek() = %v, want %v",
					gotNext,
					tt.wantNext,
				)
			}
		})
	}
}

func TestCalendar_WeekMethods_MultipleTimezones(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location string
		value    time.Time
	}{
		{
			name:     "UTC",
			location: "UTC",
			value:    time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "India",
			location: "Asia/Kolkata",
			value:    time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Japan",
			location: "Asia/Tokyo",
			value:    time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "London",
			location: "Europe/London",
			value:    time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Berlin",
			location: "Europe/Berlin",
			value:    time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "New York",
			location: "America/New_York",
			value:    time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Los Angeles",
			location: "America/Los_Angeles",
			value:    time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Sydney",
			location: "Australia/Sydney",
			value:    time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Auckland",
			location: "Pacific/Auckland",
			value:    time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
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

			start := cal.StartOfWeek(tt.value)
			next := cal.StartOfNextWeek(tt.value)
			end := cal.EndOfWeek(tt.value)

			if start.Location() != location {
				t.Fatalf("start location = %v, want %v", start.Location(), location)
			}

			if next.Location() != location {
				t.Fatalf("next location = %v, want %v", next.Location(), location)
			}

			if end.Location() != location {
				t.Fatalf("end location = %v, want %v", end.Location(), location)
			}

			if start.Weekday() != time.Monday {
				t.Fatalf("start weekday = %v, want Monday", start.Weekday())
			}

			if next.Weekday() != time.Monday {
				t.Fatalf("next weekday = %v, want Monday", next.Weekday())
			}

			if !start.Before(next) {
				t.Fatalf("start %v should be before next %v", start, next)
			}

			if !end.Equal(next.Add(-time.Nanosecond)) {
				t.Fatalf("end %v != next - 1ns %v", end, next.Add(-time.Nanosecond))
			}
		})
	}
}

func TestCalendar_WeekMethods_LeapYear(t *testing.T) {
	t.Parallel()

	tests := []time.Time{
		time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC),
		time.Date(2024, time.December, 31, 12, 0, 0, 0, time.UTC),
		time.Date(2028, time.February, 29, 12, 0, 0, 0, time.UTC),
		time.Date(2032, time.February, 29, 12, 0, 0, 0, time.UTC),
	}

	cal := mustNewCalendar(t, Config{
		Location:        time.UTC,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	for _, value := range tests {
		t.Run(value.Format("2006-01-02"), func(t *testing.T) {
			t.Parallel()

			start := cal.StartOfWeek(value)
			next := cal.StartOfNextWeek(value)

			if !start.Before(next) {
				t.Fatalf("start %v should be before next %v", start, next)
			}

			if !cal.IsSameWeek(value, start) {
				t.Fatalf("value and start should belong to same week")
			}

			if cal.IsSameWeek(value, next) {
				t.Fatalf("value and next should belong to different weeks")
			}
		})
	}
}

func TestAddWeeks_AlgebraicProperties(t *testing.T) {
	t.Parallel()

	location := mustLoadLocation(t, "America/New_York")

	values := []time.Time{
		time.Date(2026, time.January, 1, 12, 0, 0, 0, location),
		time.Date(2026, time.March, 8, 12, 0, 0, 0, location),
		time.Date(2026, time.July, 15, 12, 0, 0, 0, location),
		time.Date(2026, time.November, 1, 12, 0, 0, 0, location),
		time.Date(2026, time.December, 31, 12, 0, 0, 0, location),
	}

	tests := []struct {
		name  string
		value time.Time
		a     int
		b     int
	}{
		{
			name:  "positive values",
			value: values[0],
			a:     2,
			b:     3,
		},
		{
			name:  "negative values",
			value: values[1],
			a:     -2,
			b:     -3,
		},
		{
			name:  "cross DST",
			value: values[2],
			a:     10,
			b:     -4,
		},
		{
			name:  "zero",
			value: values[3],
			a:     0,
			b:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			left := AddWeeks(AddWeeks(tt.value, tt.a), tt.b)
			right := AddWeeks(tt.value, tt.a+tt.b)

			if !left.Equal(right) {
				t.Fatalf(
					"AddWeeks(AddWeeks(value, %d), %d) = %v, want %v",
					tt.a,
					tt.b,
					left,
					right,
				)
			}
		})
	}
}

func TestCalendar_WeekMethods_DoNotDependOnInputLocation(t *testing.T) {
	t.Parallel()

	location := mustLoadLocation(t, "Asia/Kolkata")

	cal := mustNewCalendar(t, Config{
		Location:        location,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})

	instant := time.Date(2026, time.March, 15, 20, 0, 0, 0, time.UTC)

	inputs := []time.Time{
		instant,
		instant.In(time.UTC),
		instant.In(mustLoadLocation(t, "Asia/Kolkata")),
		instant.In(mustLoadLocation(t, "America/New_York")),
		instant.In(mustLoadLocation(t, "Asia/Tokyo")),
	}

	want := cal.StartOfWeek(instant)

	for i, value := range inputs {
		got := cal.StartOfWeek(value)

		if !got.Equal(want) {
			t.Fatalf(
				"input %d: StartOfWeek() = %v, want %v",
				i,
				got,
				want,
			)
		}
	}
}
