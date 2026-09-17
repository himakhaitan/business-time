package business

import (
	"errors"
	"testing"
	"time"

	"github.com/himakhaitan/business-time/calendar"
)

func testCalendar(t *testing.T) *calendar.Calendar {
	t.Helper()

	cal, err := calendar.New(calendar.Config{
		Location:        time.UTC,
		WeekStart:       time.Monday,
		FiscalYearStart: time.January,
	})
	if err != nil {
		t.Fatalf("calendar.New() error = %v", err)
	}

	return cal
}

func TestNew(t *testing.T) {
	cal := testCalendar(t)
	weekend := NewFixedWeekend(time.Friday, time.Saturday)

	holidays := &testHolidayCalendar{
		dates: map[string]bool{
			"2026-01-26": true,
		},
	}

	businessCal, err := New(Config{
		Calendar: cal,
		Weekend:  weekend,
		Holidays: holidays,
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if businessCal == nil {
		t.Fatal("New() returned nil BusinessCalendar")
	}

	if businessCal.calendar != cal {
		t.Error("BusinessCalendar.calendar does not match configured calendar")
	}

	if businessCal.weekend == nil {
		t.Fatal("BusinessCalendar.weekend is nil")
	}

	friday := time.Date(2026, time.January, 2, 12, 0, 0, 0, time.UTC)
	saturday := time.Date(2026, time.January, 3, 12, 0, 0, 0, time.UTC)
	thursday := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)

	if !businessCal.weekend.IsWeekend(friday) {
		t.Error("BusinessCalendar.weekend does not recognize Friday as a weekend")
	}

	if !businessCal.weekend.IsWeekend(saturday) {
		t.Error("BusinessCalendar.weekend does not recognize Saturday as a weekend")
	}

	if businessCal.weekend.IsWeekend(thursday) {
		t.Error("BusinessCalendar.weekend incorrectly recognizes Thursday as a weekend")
	}

	if businessCal.holidays != holidays {
		t.Error("BusinessCalendar.holidays does not match configured holidays")
	}
}

func TestNew_NilCalendar(t *testing.T) {
	_, err := New(Config{})

	if !errors.Is(err, ErrNilCalendar) {
		t.Fatalf("New() error = %v, want %v", err, ErrNilCalendar)
	}
}

func TestNew_NilCalendarWithWeekend(t *testing.T) {
	weekend := StandardWeekend()

	_, err := New(Config{
		Calendar: nil,
		Weekend:  weekend,
	})

	if !errors.Is(err, ErrNilCalendar) {
		t.Fatalf("New() error = %v, want %v", err, ErrNilCalendar)
	}
}

func TestNew_NilWeekend(t *testing.T) {
	cal := testCalendar(t)

	_, err := New(Config{
		Calendar: cal,
	})

	if !errors.Is(err, ErrNilWeekend) {
		t.Fatalf("New() error = %v, want %v", err, ErrNilWeekend)
	}
}

func TestNew_CustomWeekend(t *testing.T) {
	cal := testCalendar(t)

	weekend := NewFixedWeekend(
		time.Friday,
		time.Saturday,
	)

	businessCal, err := New(Config{
		Calendar: cal,
		Weekend:  weekend,
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name string
		day  time.Weekday
		want bool
	}{
		{
			name: "friday is weekend",
			day:  time.Friday,
			want: true,
		},
		{
			name: "saturday is weekend",
			day:  time.Saturday,
			want: true,
		},
		{
			name: "sunday is working day",
			day:  time.Sunday,
			want: false,
		},
		{
			name: "monday is working day",
			day:  time.Monday,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := businessCal.Weekend().IsWeekend(dateOnWeekday(tt.day))

			if got != tt.want {
				t.Errorf(
					"Weekend().IsWeekend() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestNew_NilHolidays(t *testing.T) {
	cal := testCalendar(t)

	businessCal, err := New(Config{
		Calendar: cal,
		Weekend:  StandardWeekend(),
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if businessCal.Holidays() != nil {
		t.Error("Holidays() != nil, want nil")
	}
}

func TestNew_CustomHolidays(t *testing.T) {
	cal := testCalendar(t)

	holidays := &testHolidayCalendar{
		dates: map[string]bool{
			"2026-01-26": true,
			"2026-08-15": true,
		},
	}

	businessCal, err := New(Config{
		Calendar: cal,
		Weekend:  StandardWeekend(),
		Holidays: holidays,
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if businessCal.Holidays() == nil {
		t.Fatal("Holidays() returned nil")
	}

	if businessCal.Holidays() != holidays {
		t.Error("Holidays() does not return configured holiday calendar")
	}

	if !businessCal.Holidays().IsHoliday(
		time.Date(2026, time.January, 26, 12, 0, 0, 0, time.UTC),
	) {
		t.Error("configured holiday was not returned as a holiday")
	}
}

func TestNew_PreservesCalendarConfiguration(t *testing.T) {
	cal, err := calendar.New(calendar.Config{
		Location:        time.FixedZone("TEST", 5*60*60+30*60),
		WeekStart:       time.Sunday,
		FiscalYearStart: time.April,
	})
	if err != nil {
		t.Fatalf("calendar.New() error = %v", err)
	}

	businessCal, err := New(Config{
		Calendar: cal,
		Weekend:  StandardWeekend(),
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if businessCal.Calendar() != cal {
		t.Fatal("Calendar() does not return the configured calendar")
	}

	// The BusinessCalendar should retain the exact calendar instance rather
	// than creating or replacing it during construction.
	if businessCal.calendar != cal {
		t.Error("underlying calendar instance was replaced")
	}
}

func TestCalendar(t *testing.T) {
	cal := testCalendar(t)

	businessCal, err := New(Config{
		Calendar: cal,
		Weekend:  StandardWeekend(),
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got := businessCal.Calendar()

	if got == nil {
		t.Fatal("Calendar() returned nil")
	}

	if got != cal {
		t.Error("Calendar() returned a different calendar")
	}
}

func TestWeekend(t *testing.T) {
	cal := testCalendar(t)
	weekend := SundayWeekend()

	businessCal, err := New(Config{
		Calendar: cal,
		Weekend:  weekend,
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got := businessCal.Weekend()

	if got == nil {
		t.Fatal("Weekend() returned nil")
	}

	sunday := time.Date(
		2026, time.January, 4,
		12, 0, 0, 0,
		time.UTC,
	)

	monday := time.Date(
		2026, time.January, 5,
		12, 0, 0, 0,
		time.UTC,
	)

	if !got.IsWeekend(sunday) {
		t.Error("Weekend() returned a rule that does not recognize Sunday as a weekend")
	}

	if got.IsWeekend(monday) {
		t.Error("Weekend() returned a rule that incorrectly recognizes Monday as a weekend")
	}
}

func TestHolidays(t *testing.T) {
	cal := testCalendar(t)

	holidays := &testHolidayCalendar{
		dates: map[string]bool{
			"2026-01-26": true,
		},
	}

	businessCal, err := New(Config{
		Calendar: cal,
		Holidays: holidays,
		Weekend:  StandardWeekend(),
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got := businessCal.Holidays()

	if got == nil {
		t.Fatal("Holidays() returned nil")
	}

	if got != holidays {
		t.Error("Holidays() returned a different holiday calendar")
	}
}

func TestNew_WeekendWithHoliday(t *testing.T) {
	cal := testCalendar(t)

	holidays := testHolidayCalendar{
		dates: map[string]bool{
			"2026-01-05": true,
		},
	}

	businessCal, err := New(Config{
		Calendar: cal,
		Weekend:  StandardWeekend(),
		Holidays: holidays,
	})

	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	saturday := time.Date(
		2026, time.January, 3,
		12, 0, 0, 0,
		time.UTC,
	)

	sunday := time.Date(
		2026, time.January, 4,
		12, 0, 0, 0,
		time.UTC,
	)

	monday := time.Date(
		2026, time.January, 5,
		12, 0, 0, 0,
		time.UTC,
	)

	if !businessCal.Weekend().IsWeekend(saturday) {
		t.Error("Saturday should be a weekend")
	}

	if !businessCal.Weekend().IsWeekend(sunday) {
		t.Error("Sunday should be a weekend")
	}

	if businessCal.Weekend().IsWeekend(monday) {
		t.Error("Monday should not be a weekend")
	}

	if !businessCal.Holidays().IsHoliday(monday) {
		t.Error("Monday should be a holiday")
	}
}

func TestNew_DoesNotMutateConfig(t *testing.T) {
	cal := testCalendar(t)
	weekend := SundayWeekend()
	holidays := &testHolidayCalendar{
		dates: map[string]bool{
			"2026-01-01": true,
		},
	}

	config := Config{
		Calendar: cal,
		Weekend:  weekend,
		Holidays: holidays,
	}

	_, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if config.Calendar != cal {
		t.Error("New() mutated Config.Calendar")
	}

	if config.Weekend == nil {
		t.Fatal("New() mutated Config.Weekend to nil")
	}

	sunday := time.Date(
		2026, time.January, 4,
		12, 0, 0, 0,
		time.UTC,
	)

	monday := time.Date(
		2026, time.January, 5,
		12, 0, 0, 0,
		time.UTC,
	)

	if !config.Weekend.IsWeekend(sunday) {
		t.Error("New() mutated Config.Weekend behavior for Sunday")
	}

	if config.Weekend.IsWeekend(monday) {
		t.Error("New() mutated Config.Weekend behavior for Monday")
	}

	if config.Holidays != holidays {
		t.Error("New() mutated Config.Holidays")
	}
}

func TestNew_MultipleBusinessCalendars(t *testing.T) {
	cal := testCalendar(t)

	standardWeekend := StandardWeekend()
	sundayWeekend := SundayWeekend()

	first, err := New(Config{
		Calendar: cal,
		Weekend:  standardWeekend,
	})
	if err != nil {
		t.Fatalf("first New() error = %v", err)
	}

	second, err := New(Config{
		Calendar: cal,
		Weekend:  sundayWeekend,
	})
	if err != nil {
		t.Fatalf("second New() error = %v", err)
	}

	if first == second {
		t.Fatal("New() returned the same BusinessCalendar instance")
	}

	if first.Calendar() != second.Calendar() {
		t.Error("calendars should reference the configured calendar")
	}

	saturday := dateOnWeekday(time.Saturday)
	sunday := dateOnWeekday(time.Sunday)

	if !first.Weekend().IsWeekend(saturday) {
		t.Error("standard weekend should include Saturday")
	}

	if first.Weekend().IsWeekend(sunday) == false {
		t.Error("standard weekend should include Sunday")
	}

	if second.Weekend().IsWeekend(saturday) {
		t.Error("Sunday-only weekend should not include Saturday")
	}

	if !second.Weekend().IsWeekend(sunday) {
		t.Error("Sunday-only weekend should include Sunday")
	}
}

type testHolidayCalendar struct {
	dates map[string]bool
}

func (h testHolidayCalendar) IsHoliday(t time.Time) bool {
	return h.dates[t.Format("2006-01-02")]
}

func dateOnWeekday(day time.Weekday) time.Time {
	// 2026-01-04 is a Sunday.
	sunday := time.Date(
		2026,
		time.January,
		4,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	offset := (int(day) - int(time.Sunday) + 7) % 7

	return sunday.AddDate(0, 0, offset)
}
