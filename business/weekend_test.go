package business

import (
	"testing"
	"time"
)

func TestNewFixedWeekend(t *testing.T) {
	tests := []struct {
		name  string
		days  []time.Weekday
		check time.Weekday
		want  bool
	}{
		{
			name:  "single weekday",
			days:  []time.Weekday{time.Sunday},
			check: time.Sunday,
			want:  true,
		},
		{
			name:  "multiple weekdays",
			days:  []time.Weekday{time.Saturday, time.Sunday},
			check: time.Saturday,
			want:  true,
		},
		{
			name:  "weekday not configured",
			days:  []time.Weekday{time.Saturday, time.Sunday},
			check: time.Monday,
			want:  false,
		},
		{
			name:  "no weekdays",
			days:  nil,
			check: time.Sunday,
			want:  false,
		},
		{
			name: "all weekdays",
			days: []time.Weekday{
				time.Sunday,
				time.Monday,
				time.Tuesday,
				time.Wednesday,
				time.Thursday,
				time.Friday,
				time.Saturday,
			},
			check: time.Wednesday,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weekend := NewFixedWeekend(tt.days...)

			got := weekend.IsWeekend(dateOnWeekday(tt.check))
			if got != tt.want {
				t.Errorf("IsWeekend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixedWeekend_IsWeekend(t *testing.T) {
	weekend := NewFixedWeekend(
		time.Saturday,
		time.Sunday,
	)

	tests := []struct {
		name string
		date time.Time
		want bool
	}{
		{
			name: "saturday",
			date: time.Date(2026, time.January, 3, 0, 0, 0, 0, time.UTC),
			want: true,
		},
		{
			name: "sunday",
			date: time.Date(2026, time.January, 4, 0, 0, 0, 0, time.UTC),
			want: true,
		},
		{
			name: "monday",
			date: time.Date(2026, time.January, 5, 0, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "tuesday",
			date: time.Date(2026, time.January, 6, 0, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "wednesday",
			date: time.Date(2026, time.January, 7, 0, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "thursday",
			date: time.Date(2026, time.January, 8, 0, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "friday",
			date: time.Date(2026, time.January, 9, 0, 0, 0, 0, time.UTC),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := weekend.IsWeekend(tt.date)
			if got != tt.want {
				t.Errorf("IsWeekend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixedWeekend_IsWeekend_IgnoresTimeOfDay(t *testing.T) {
	weekend := NewFixedWeekend(
		time.Saturday,
		time.Sunday,
	)

	tests := []struct {
		name string
		date time.Time
	}{
		{
			name: "start of day",
			date: time.Date(2026, time.January, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "morning",
			date: time.Date(2026, time.January, 3, 9, 30, 0, 0, time.UTC),
		},
		{
			name: "afternoon",
			date: time.Date(2026, time.January, 3, 14, 45, 30, 0, time.UTC),
		},
		{
			name: "evening",
			date: time.Date(2026, time.January, 3, 20, 15, 59, 0, time.UTC),
		},
		{
			name: "end of day",
			date: time.Date(2026, time.January, 3, 23, 59, 59, 999999999, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := weekend.IsWeekend(tt.date); !got {
				t.Errorf("IsWeekend() = false, want true for %v", tt.date)
			}
		})
	}
}

func TestFixedWeekend_DuplicateWeekdays(t *testing.T) {
	weekend := NewFixedWeekend(
		time.Saturday,
		time.Saturday,
		time.Sunday,
		time.Sunday,
	)

	tests := []struct {
		name string
		day  time.Weekday
		want bool
	}{
		{
			name: "saturday",
			day:  time.Saturday,
			want: true,
		},
		{
			name: "sunday",
			day:  time.Sunday,
			want: true,
		},
		{
			name: "monday",
			day:  time.Monday,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := weekend.IsWeekend(dateOnWeekday(tt.day))
			if got != tt.want {
				t.Errorf("IsWeekend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixedWeekend_NoDays(t *testing.T) {
	weekend := NewFixedWeekend()

	for day := time.Sunday; day <= time.Saturday; day++ {
		t.Run(day.String(), func(t *testing.T) {
			if got := weekend.IsWeekend(dateOnWeekday(day)); got {
				t.Errorf("IsWeekend() = true for %v, want false", day)
			}
		})
	}
}

func TestFixedWeekend_AllDays(t *testing.T) {
	weekend := NewFixedWeekend(
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	)

	for day := time.Sunday; day <= time.Saturday; day++ {
		t.Run(day.String(), func(t *testing.T) {
			if got := weekend.IsWeekend(dateOnWeekday(day)); !got {
				t.Errorf("IsWeekend() = false for %v, want true", day)
			}
		})
	}
}

func TestFixedWeekend_DifferentLocations(t *testing.T) {
	weekend := NewFixedWeekend(
		time.Saturday,
		time.Sunday,
	)

	locations := []struct {
		name string
		loc  *time.Location
	}{
		{
			name: "UTC",
			loc:  time.UTC,
		},
		{
			name: "New York",
			loc:  time.FixedZone("EST", -5*60*60),
		},
		{
			name: "India",
			loc:  time.FixedZone("IST", 5*60*60+30*60),
		},
	}

	for _, location := range locations {
		t.Run(location.name, func(t *testing.T) {
			date := time.Date(
				2026,
				time.January,
				3,
				15,
				30,
				0,
				0,
				location.loc,
			)

			if got := weekend.IsWeekend(date); !got {
				t.Errorf("IsWeekend() = false, want true")
			}
		})
	}
}

func TestStandardWeekend(t *testing.T) {
	weekend := StandardWeekend()

	tests := []struct {
		day  time.Weekday
		want bool
	}{
		{time.Sunday, true},
		{time.Monday, false},
		{time.Tuesday, false},
		{time.Wednesday, false},
		{time.Thursday, false},
		{time.Friday, false},
		{time.Saturday, true},
	}

	for _, tt := range tests {
		t.Run(tt.day.String(), func(t *testing.T) {
			got := weekend.IsWeekend(dateOnWeekday(tt.day))
			if got != tt.want {
				t.Errorf("IsWeekend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSundayWeekend(t *testing.T) {
	weekend := SundayWeekend()

	tests := []struct {
		day  time.Weekday
		want bool
	}{
		{time.Sunday, true},
		{time.Monday, false},
		{time.Tuesday, false},
		{time.Wednesday, false},
		{time.Thursday, false},
		{time.Friday, false},
		{time.Saturday, false},
	}

	for _, tt := range tests {
		t.Run(tt.day.String(), func(t *testing.T) {
			got := weekend.IsWeekend(dateOnWeekday(tt.day))
			if got != tt.want {
				t.Errorf("IsWeekend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeekendRules_ImplementWeekendRule(t *testing.T) {
	var _ WeekendRule = FixedWeekend{}
	var _ WeekendRule = NewFixedWeekend(time.Saturday, time.Sunday)
	var _ WeekendRule = StandardWeekend()
	var _ WeekendRule = SundayWeekend()
}
