package business

import (
	"testing"
	"time"
)

func TestHolidayCalendar_IsHoliday(t *testing.T) {
	holidays := &testHolidayCalendar{
		dates: map[string]bool{
			"2026-01-01": true,
			"2026-01-26": true,
			"2026-08-15": true,
		},
	}

	tests := []struct {
		name string
		date time.Time
		want bool
	}{
		{
			name: "holiday",
			date: time.Date(
				2026,
				time.January,
				1,
				0,
				0,
				0,
				0,
				time.UTC,
			),
			want: true,
		},
		{
			name: "another holiday",
			date: time.Date(
				2026,
				time.January,
				26,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			want: true,
		},
		{
			name: "non-holiday",
			date: time.Date(
				2026,
				time.January,
				2,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			want: false,
		},
		{
			name: "another non-holiday",
			date: time.Date(
				2026,
				time.January,
				27,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := holidays.IsHoliday(tt.date)

			if got != tt.want {
				t.Errorf("IsHoliday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHolidayCalendar_IsHoliday_IndependentOfTimeOfDay(t *testing.T) {
	holidays := &testHolidayCalendar{
		dates: map[string]bool{
			"2026-01-01": true,
		},
	}

	tests := []struct {
		name string
		date time.Time
	}{
		{
			name: "midnight",
			date: time.Date(
				2026,
				time.January,
				1,
				0,
				0,
				0,
				0,
				time.UTC,
			),
		},
		{
			name: "morning",
			date: time.Date(
				2026,
				time.January,
				1,
				9,
				30,
				0,
				0,
				time.UTC,
			),
		},
		{
			name: "afternoon",
			date: time.Date(
				2026,
				time.January,
				1,
				14,
				45,
				0,
				0,
				time.UTC,
			),
		},
		{
			name: "evening",
			date: time.Date(
				2026,
				time.January,
				1,
				20,
				15,
				0,
				0,
				time.UTC,
			),
		},
		{
			name: "end of day",
			date: time.Date(
				2026,
				time.January,
				1,
				23,
				59,
				59,
				999999999,
				time.UTC,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := holidays.IsHoliday(tt.date); !got {
				t.Errorf("IsHoliday() = false, want true for %v", tt.date)
			}
		})
	}
}

func TestHolidayCalendar_IsHoliday_DifferentDates(t *testing.T) {
	holidays := &testHolidayCalendar{
		dates: map[string]bool{
			"2026-01-01": true,
		},
	}

	tests := []struct {
		name string
		date time.Time
		want bool
	}{
		{
			name: "day before",
			date: time.Date(
				2025,
				time.December,
				31,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			want: false,
		},
		{
			name: "holiday",
			date: time.Date(
				2026,
				time.January,
				1,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			want: true,
		},
		{
			name: "day after",
			date: time.Date(
				2026,
				time.January,
				2,
				12,
				0,
				0,
				0,
				time.UTC,
			),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := holidays.IsHoliday(tt.date)

			if got != tt.want {
				t.Errorf("IsHoliday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHolidayCalendar_IsHoliday_DifferentLocations(t *testing.T) {
	holidays := &testHolidayCalendar{
		dates: map[string]bool{
			"2026-01-01": true,
		},
	}

	locations := []struct {
		name string
		loc  *time.Location
	}{
		{
			name: "UTC",
			loc:  time.UTC,
		},
		{
			name: "India",
			loc:  time.FixedZone("IST", 5*60*60+30*60),
		},
		{
			name: "New York",
			loc:  time.FixedZone("EST", -5*60*60),
		},
	}

	for _, location := range locations {
		t.Run(location.name, func(t *testing.T) {
			date := time.Date(
				2026,
				time.January,
				1,
				12,
				0,
				0,
				0,
				location.loc,
			)

			if got := holidays.IsHoliday(date); !got {
				t.Errorf("IsHoliday() = false, want true")
			}
		})
	}
}

func TestHolidayCalendar_ImplementsInterface(t *testing.T) {
	var _ HolidayCalendar = (*testHolidayCalendar)(nil)
}
