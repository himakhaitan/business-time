package calendar

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Parallel()

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}

	cfg := Config{
		Location:        loc,
		WeekStart:       time.Monday,
		FiscalYearStart: time.April,
	}

	cal, err := New(cfg)
	if err != nil {
		t.Fatalf("New() returned unexpected error: %v", err)
	}

	if cal == nil {
		t.Fatal("New() returned nil Calendar")
	}

	if got := cal.Location(); got != loc {
		t.Errorf("Location() = %v, want %v", got, loc)
	}

	if got := cal.WeekStart(); got != time.Monday {
		t.Errorf("WeekStart() = %v, want %v", got, time.Monday)
	}

	if got := cal.FiscalYearStart(); got != time.April {
		t.Errorf("FiscalYearStart() = %v, want %v", got, time.April)
	}
}

func TestNewRequiresLocation(t *testing.T) {
	t.Parallel()

	cfg := Config{
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	}

	cal, err := New(cfg)

	if err == nil {
		t.Fatal("New() expected error for nil location, got nil")
	}

	if cal != nil {
		t.Fatalf("New() returned Calendar with invalid config: %#v", cal)
	}

	const want = "calendar: location is required"
	if err.Error() != want {
		t.Errorf("New() error = %q, want %q", err.Error(), want)
	}
}

func TestNewRejectsInvalidWeekStart(t *testing.T) {
	t.Parallel()

	loc := time.UTC

	tests := []struct {
		name      string
		weekStart time.Weekday
	}{
		{
			name:      "below Sunday",
			weekStart: time.Weekday(-1),
		},
		{
			name:      "above Saturday",
			weekStart: time.Weekday(7),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cal, err := New(Config{
				Location:        loc,
				WeekStart:       tt.weekStart,
				FiscalYearStart: time.January,
			})

			if err == nil {
				t.Fatal("New() expected error, got nil")
			}

			if cal != nil {
				t.Fatalf("New() returned Calendar with invalid config: %#v", cal)
			}
		})
	}
}

func TestNewRejectsInvalidFiscalYearStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		fiscalYearStart time.Month
	}{
		{
			name:            "below January",
			fiscalYearStart: time.Month(0),
		},
		{
			name:            "above December",
			fiscalYearStart: time.Month(13),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cal, err := New(Config{
				Location:        time.UTC,
				WeekStart:       time.Monday,
				FiscalYearStart: tt.fiscalYearStart,
			})

			if err == nil {
				t.Fatal("New() expected error, got nil")
			}

			if cal != nil {
				t.Fatalf("New() returned Calendar with invalid config: %#v", cal)
			}
		})
	}
}

func TestNewAcceptsAllValidWeekdays(t *testing.T) {
	t.Parallel()

	for day := time.Sunday; day <= time.Saturday; day++ {
		day := day

		t.Run(day.String(), func(t *testing.T) {
			t.Parallel()

			cal, err := New(Config{
				Location:        time.UTC,
				WeekStart:       day,
				FiscalYearStart: time.January,
			})

			if err != nil {
				t.Fatalf("New() returned unexpected error for %v: %v", day, err)
			}

			if got := cal.WeekStart(); got != day {
				t.Errorf("WeekStart() = %v, want %v", got, day)
			}
		})
	}
}

func TestNewAcceptsAllValidFiscalYearStartMonths(t *testing.T) {
	t.Parallel()

	for month := time.January; month <= time.December; month++ {
		month := month

		t.Run(month.String(), func(t *testing.T) {
			t.Parallel()

			cal, err := New(Config{
				Location:        time.UTC,
				WeekStart:       time.Monday,
				FiscalYearStart: month,
			})

			if err != nil {
				t.Fatalf(
					"New() returned unexpected error for %v: %v",
					month,
					err,
				)
			}

			if got := cal.FiscalYearStart(); got != month {
				t.Errorf(
					"FiscalYearStart() = %v, want %v",
					got,
					month,
				)
			}
		})
	}
}

func TestNewAcceptsDifferentLocations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		loc  *time.Location
	}{
		{
			name: "UTC",
			loc:  time.UTC,
		},
		{
			name: "Local",
			loc:  time.Local,
		},
	}

	locations := []string{
		"Asia/Kolkata",
		"America/New_York",
		"America/Los_Angeles",
		"Europe/London",
		"Europe/Berlin",
		"Australia/Sydney",
		"Pacific/Auckland",
	}

	for _, name := range locations {
		loc, err := time.LoadLocation(name)
		if err != nil {
			t.Fatalf("load location %q: %v", name, err)
		}

		tests = append(tests, struct {
			name string
			loc  *time.Location
		}{
			name: name,
			loc:  loc,
		})
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cal, err := New(Config{
				Location:        tt.loc,
				WeekStart:       time.Monday,
				FiscalYearStart: time.January,
			})

			if err != nil {
				t.Fatalf("New() returned unexpected error: %v", err)
			}

			if got := cal.Location(); got != tt.loc {
				t.Errorf("Location() = %v, want %v", got, tt.loc)
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

func mustNewCalendar(t *testing.T, cfg Config) *Calendar {
	t.Helper()

	cal, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	return cal
}
