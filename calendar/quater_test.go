package calendar

import (
	"strconv"
	"testing"
	"time"
)

func TestStartOfQuarter_CalendarQuarters(t *testing.T) {
	tests := []struct {
		name        string
		value       time.Time
		wantYear    int
		wantMonth   time.Month
		wantQuarter int
	}{
		{
			name:        "January",
			value:       time.Date(2026, time.January, 15, 12, 30, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.January,
			wantQuarter: 1,
		},
		{
			name:        "March",
			value:       time.Date(2026, time.March, 31, 23, 59, 59, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.January,
			wantQuarter: 1,
		},
		{
			name:        "April",
			value:       time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.April,
			wantQuarter: 2,
		},
		{
			name:        "June",
			value:       time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.April,
			wantQuarter: 2,
		},
		{
			name:        "July",
			value:       time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.July,
			wantQuarter: 3,
		},
		{
			name:        "September",
			value:       time.Date(2026, time.September, 30, 23, 59, 59, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.July,
			wantQuarter: 3,
		},
		{
			name:        "October",
			value:       time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.October,
			wantQuarter: 4,
		},
		{
			name:        "December",
			value:       time.Date(2026, time.December, 31, 23, 59, 59, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.October,
			wantQuarter: 4,
		},
	}

	cal := newTestCalendar(t, time.UTC)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cal.StartOfQuarter(tt.value)

			want := time.Date(
				tt.wantYear,
				tt.wantMonth,
				1,
				0, 0, 0, 0,
				time.UTC,
			)

			if !got.Equal(want) {
				t.Errorf("StartOfQuarter() = %v, want %v", got, want)
			}

			if got.Location() != time.UTC {
				t.Errorf(
					"StartOfQuarter() location = %v, want %v",
					got.Location(),
					time.UTC,
				)
			}

			if gotQuarter := cal.Quarter(tt.value); gotQuarter != tt.wantQuarter {
				t.Errorf(
					"Quarter() = %d, want %d",
					gotQuarter,
					tt.wantQuarter,
				)
			}
		})
	}
}

func TestStartOfQuarter_FiscalYearStarts(t *testing.T) {
	tests := []struct {
		fiscalStart time.Month
		value       time.Time
		wantYear    int
		wantMonth   time.Month
		wantQuarter int
	}{
		// Fiscal year starts in January.
		{
			fiscalStart: time.January,
			value:       time.Date(2026, time.February, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.January,
			wantQuarter: 1,
		},

		// Fiscal year starts in April.
		{
			fiscalStart: time.April,
			value:       time.Date(2026, time.May, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.April,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.April,
			value:       time.Date(2026, time.August, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.July,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.April,
			value:       time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.October,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.April,
			value:       time.Date(2027, time.February, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.January,
			wantQuarter: 4,
		},

		// Fiscal year starts in July.
		{
			fiscalStart: time.July,
			value:       time.Date(2026, time.August, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.July,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.July,
			value:       time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.October,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.July,
			value:       time.Date(2027, time.February, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.January,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.July,
			value:       time.Date(2027, time.May, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.April,
			wantQuarter: 4,
		},

		// Fiscal year starts in October.
		{
			fiscalStart: time.October,
			value:       time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.October,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.October,
			value:       time.Date(2027, time.February, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.January,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.October,
			value:       time.Date(2027, time.May, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.April,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.October,
			value:       time.Date(2027, time.August, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.July,
			wantQuarter: 4,
		},
	}

	for _, tt := range tests {
		name := tt.fiscalStart.String() + "_" + tt.value.Format("2006-01-02")

		t.Run(name, func(t *testing.T) {
			cal := mustNewCalendar(
				t,
				Config{
					Location:        time.UTC,
					WeekStart:       time.Monday,
					FiscalYearStart: tt.fiscalStart,
				},
			)

			got := cal.StartOfQuarter(tt.value)

			want := time.Date(
				tt.wantYear,
				tt.wantMonth,
				1,
				0, 0, 0, 0,
				time.UTC,
			)

			if !got.Equal(want) {
				t.Errorf("StartOfQuarter() = %v, want %v", got, want)
			}

			if gotQuarter := cal.Quarter(tt.value); gotQuarter != tt.wantQuarter {
				t.Errorf(
					"Quarter() = %d, want %d",
					gotQuarter,
					tt.wantQuarter,
				)
			}
		})
	}
}

func TestStartOfQuarter_AllFiscalStartMonths(t *testing.T) {
	tests := []struct {
		fiscalStart time.Month
		value       time.Time
		wantYear    int
		wantMonth   time.Month
		wantQuarter int
	}{
		// January fiscal year:
		//
		// Q1: Jan-Mar
		// Q2: Apr-Jun
		// Q3: Jul-Sep
		// Q4: Oct-Dec
		{
			fiscalStart: time.January,
			value:       time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.January,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.January,
			value:       time.Date(2026, time.April, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.April,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.January,
			value:       time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.July,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.January,
			value:       time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.October,
			wantQuarter: 4,
		},

		// February fiscal year:
		//
		// Q1: Feb-Apr
		// Q2: May-Jul
		// Q3: Aug-Oct
		// Q4: Nov-Jan
		{
			fiscalStart: time.February,
			value:       time.Date(2026, time.February, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.February,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.February,
			value:       time.Date(2026, time.May, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.May,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.February,
			value:       time.Date(2026, time.August, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.August,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.February,
			value:       time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.November,
			wantQuarter: 4,
		},

		// March fiscal year:
		//
		// Q1: Mar-May
		// Q2: Jun-Aug
		// Q3: Sep-Nov
		// Q4: Dec-Feb
		{
			fiscalStart: time.March,
			value:       time.Date(2026, time.March, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.March,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.March,
			value:       time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.June,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.March,
			value:       time.Date(2026, time.September, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.September,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.March,
			value:       time.Date(2026, time.December, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.December,
			wantQuarter: 4,
		},

		// April fiscal year:
		//
		// Q1: Apr-Jun
		// Q2: Jul-Sep
		// Q3: Oct-Dec
		// Q4: Jan-Mar
		{
			fiscalStart: time.April,
			value:       time.Date(2026, time.April, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.April,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.April,
			value:       time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.July,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.April,
			value:       time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.October,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.April,
			value:       time.Date(2027, time.January, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.January,
			wantQuarter: 4,
		},

		// May fiscal year:
		//
		// Q1: May-Jul
		// Q2: Aug-Oct
		// Q3: Nov-Jan
		// Q4: Feb-Apr
		{
			fiscalStart: time.May,
			value:       time.Date(2026, time.May, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.May,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.May,
			value:       time.Date(2026, time.August, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.August,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.May,
			value:       time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.November,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.May,
			value:       time.Date(2027, time.February, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.February,
			wantQuarter: 4,
		},

		// June fiscal year.
		{
			fiscalStart: time.June,
			value:       time.Date(2026, time.June, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.June,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.June,
			value:       time.Date(2026, time.September, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.September,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.June,
			value:       time.Date(2026, time.December, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.December,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.June,
			value:       time.Date(2027, time.March, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.March,
			wantQuarter: 4,
		},

		// July fiscal year.
		{
			fiscalStart: time.July,
			value:       time.Date(2026, time.July, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.July,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.July,
			value:       time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.October,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.July,
			value:       time.Date(2027, time.January, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.January,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.July,
			value:       time.Date(2027, time.April, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.April,
			wantQuarter: 4,
		},

		// August fiscal year.
		{
			fiscalStart: time.August,
			value:       time.Date(2026, time.August, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.August,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.August,
			value:       time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.November,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.August,
			value:       time.Date(2027, time.February, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.February,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.August,
			value:       time.Date(2027, time.May, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.May,
			wantQuarter: 4,
		},

		// September fiscal year.
		{
			fiscalStart: time.September,
			value:       time.Date(2026, time.September, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.September,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.September,
			value:       time.Date(2026, time.December, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.December,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.September,
			value:       time.Date(2027, time.March, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.March,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.September,
			value:       time.Date(2027, time.June, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.June,
			wantQuarter: 4,
		},

		// October fiscal year.
		{
			fiscalStart: time.October,
			value:       time.Date(2026, time.October, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.October,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.October,
			value:       time.Date(2027, time.January, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.January,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.October,
			value:       time.Date(2027, time.April, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.April,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.October,
			value:       time.Date(2027, time.July, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.July,
			wantQuarter: 4,
		},

		// November fiscal year.
		{
			fiscalStart: time.November,
			value:       time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.November,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.November,
			value:       time.Date(2027, time.February, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.February,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.November,
			value:       time.Date(2027, time.May, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.May,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.November,
			value:       time.Date(2027, time.August, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.August,
			wantQuarter: 4,
		},

		// December fiscal year.
		{
			fiscalStart: time.December,
			value:       time.Date(2026, time.December, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2026,
			wantMonth:   time.December,
			wantQuarter: 1,
		},
		{
			fiscalStart: time.December,
			value:       time.Date(2027, time.March, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.March,
			wantQuarter: 2,
		},
		{
			fiscalStart: time.December,
			value:       time.Date(2027, time.June, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.June,
			wantQuarter: 3,
		},
		{
			fiscalStart: time.December,
			value:       time.Date(2027, time.September, 15, 12, 0, 0, 0, time.UTC),
			wantYear:    2027,
			wantMonth:   time.September,
			wantQuarter: 4,
		},
	}

	for _, tt := range tests {
		name := tt.fiscalStart.String()

		t.Run(name, func(t *testing.T) {
			cal := mustNewCalendar(
				t,
				Config{
					Location:        time.UTC,
					WeekStart:       time.Monday,
					FiscalYearStart: tt.fiscalStart,
				},
			)

			got := cal.StartOfQuarter(tt.value)

			want := time.Date(
				tt.wantYear,
				tt.wantMonth,
				1,
				0, 0, 0, 0,
				time.UTC,
			)

			if !got.Equal(want) {
				t.Errorf(
					"StartOfQuarter() = %v, want %v",
					got,
					want,
				)
			}

			if gotQuarter := cal.Quarter(tt.value); gotQuarter != tt.wantQuarter {
				t.Errorf(
					"Quarter() = %d, want %d",
					gotQuarter,
					tt.wantQuarter,
				)
			}
		})
	}
}

func TestStartOfQuarter_ExactlyAtBoundaries(t *testing.T) {
	for fiscalStart := time.January; fiscalStart <= time.December; fiscalStart++ {
		t.Run(fiscalStart.String(), func(t *testing.T) {
			cal := mustNewCalendar(
				t,
				Config{
					Location:        time.UTC,
					WeekStart:       time.Monday,
					FiscalYearStart: fiscalStart,
				},
			)

			for quarter := 0; quarter < 4; quarter++ {
				month := time.Month(
					(int(fiscalStart)-1+quarter*3)%12 + 1,
				)

				year := 2026

				if month < fiscalStart {
					year++
				}

				value := time.Date(
					year,
					month,
					1,
					0, 0, 0, 0,
					time.UTC,
				)

				got := cal.StartOfQuarter(value)

				if !got.Equal(value) {
					t.Errorf(
						"StartOfQuarter(%v) = %v, want %v",
						value,
						got,
						value,
					)
				}

				if gotQuarter := cal.Quarter(value); gotQuarter != quarter+1 {
					t.Errorf(
						"Quarter(%v) = %d, want %d",
						value,
						gotQuarter,
						quarter+1,
					)
				}
			}
		})
	}
}

func TestStartOfQuarter_QuarterEndBoundaries(t *testing.T) {
	cal := mustNewCalendar(
		t,
		Config{
			Location:        time.UTC,
			WeekStart:       time.Monday,
			FiscalYearStart: time.April,
		},
	)

	tests := []struct {
		name  string
		value time.Time
	}{
		{
			name:  "June last instant",
			value: time.Date(2026, time.June, 30, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:  "September last instant",
			value: time.Date(2026, time.September, 30, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:  "December last instant",
			value: time.Date(2026, time.December, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:  "March last instant",
			value: time.Date(2027, time.March, 31, 23, 59, 59, 999999999, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cal.StartOfQuarter(tt.value)

			quarter := cal.Quarter(tt.value)

			if quarter < 1 || quarter > 4 {
				t.Fatalf("Quarter() = %d, want 1..4", quarter)
			}

			next := cal.StartOfNextQuarter(tt.value)

			if !got.Before(tt.value) {
				t.Errorf(
					"StartOfQuarter() = %v, want before value %v",
					got,
					tt.value,
				)
			}

			if !next.After(tt.value) {
				t.Errorf(
					"StartOfNextQuarter() = %v, want after value %v",
					next,
					tt.value,
				)
			}
		})
	}
}

func TestEndOfQuarter(t *testing.T) {
	cal := mustNewCalendar(
		t,
		Config{
			Location:        time.UTC,
			WeekStart:       time.Monday,
			FiscalYearStart: time.April,
		},
	)

	tests := []struct {
		name  string
		value time.Time
	}{
		{
			name:  "Q1",
			value: time.Date(2026, time.May, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:  "Q2",
			value: time.Date(2026, time.August, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:  "Q3",
			value: time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
		},
		{
			name:  "Q4",
			value: time.Date(2027, time.February, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cal.EndOfQuarter(tt.value)
			next := cal.StartOfNextQuarter(tt.value)

			want := next.Add(-time.Nanosecond)

			if !got.Equal(want) {
				t.Errorf("EndOfQuarter() = %v, want %v", got, want)
			}

			if got.Location() != next.Location() {
				t.Errorf(
					"EndOfQuarter() location = %v, want %v",
					got.Location(),
					next.Location(),
				)
			}

			if got.After(next) {
				t.Errorf(
					"EndOfQuarter() = %v, must not be after next quarter %v",
					got,
					next,
				)
			}

			if !got.Add(time.Nanosecond).Equal(next) {
				t.Errorf(
					"EndOfQuarter() + 1ns = %v, want %v",
					got.Add(time.Nanosecond),
					next,
				)
			}
		})
	}
}

func TestStartOfNextQuarter(t *testing.T) {
	cal := mustNewCalendar(
		t,
		Config{
			Location:        time.UTC,
			WeekStart:       time.Monday,
			FiscalYearStart: time.April,
		},
	)

	tests := []struct {
		value     time.Time
		wantYear  int
		wantMonth time.Month
	}{
		{
			value:     time.Date(2026, time.May, 15, 12, 0, 0, 0, time.UTC),
			wantYear:  2026,
			wantMonth: time.July,
		},
		{
			value:     time.Date(2026, time.August, 15, 12, 0, 0, 0, time.UTC),
			wantYear:  2026,
			wantMonth: time.October,
		},
		{
			value:     time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
			wantYear:  2027,
			wantMonth: time.January,
		},
		{
			value:     time.Date(2027, time.February, 15, 12, 0, 0, 0, time.UTC),
			wantYear:  2027,
			wantMonth: time.April,
		},
	}

	for _, tt := range tests {
		t.Run(tt.value.Format("2006-01-02"), func(t *testing.T) {
			got := cal.StartOfNextQuarter(tt.value)

			want := time.Date(
				tt.wantYear,
				tt.wantMonth,
				1,
				0, 0, 0, 0,
				time.UTC,
			)

			if !got.Equal(want) {
				t.Errorf("StartOfNextQuarter() = %v, want %v", got, want)
			}
		})
	}
}

func TestStartOfPreviousQuarter(t *testing.T) {
	cal := mustNewCalendar(
		t,
		Config{
			Location:        time.UTC,
			WeekStart:       time.Monday,
			FiscalYearStart: time.April,
		},
	)

	tests := []struct {
		value     time.Time
		wantYear  int
		wantMonth time.Month
	}{
		{
			value:     time.Date(2026, time.May, 15, 12, 0, 0, 0, time.UTC),
			wantYear:  2026,
			wantMonth: time.January,
		},
		{
			value:     time.Date(2026, time.August, 15, 12, 0, 0, 0, time.UTC),
			wantYear:  2026,
			wantMonth: time.April,
		},
		{
			value:     time.Date(2026, time.November, 15, 12, 0, 0, 0, time.UTC),
			wantYear:  2026,
			wantMonth: time.July,
		},
		{
			value:     time.Date(2027, time.February, 15, 12, 0, 0, 0, time.UTC),
			wantYear:  2026,
			wantMonth: time.October,
		},
	}

	for _, tt := range tests {
		t.Run(tt.value.Format("2006-01-02"), func(t *testing.T) {
			got := cal.StartOfPreviousQuarter(tt.value)

			want := time.Date(
				tt.wantYear,
				tt.wantMonth,
				1,
				0, 0, 0, 0,
				time.UTC,
			)

			if !got.Equal(want) {
				t.Errorf("StartOfPreviousQuarter() = %v, want %v", got, want)
			}
		})
	}
}

func TestIsSameQuarter(t *testing.T) {
	cal := mustNewCalendar(
		t,
		Config{
			Location:        time.UTC,
			WeekStart:       time.Monday,
			FiscalYearStart: time.April,
		},
	)

	tests := []struct {
		name string
		a    time.Time
		b    time.Time
		want bool
	}{
		{
			name: "same quarter",
			a:    time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC),
			want: true,
		},
		{
			name: "different quarters",
			a:    time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC),
			b:    time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "different fiscal years",
			a:    time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			b:    time.Date(2027, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cal.IsSameQuarter(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("IsSameQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSameQuarter_DifferentFiscalYears(t *testing.T) {
	cal := mustNewCalendar(
		t,
		Config{
			Location:        time.UTC,
			WeekStart:       time.Monday,
			FiscalYearStart: time.April,
		},
	)

	tests := []struct {
		name string
		a    time.Time
		b    time.Time
		want bool
	}{
		{
			name: "Q1 2026 vs Q1 2027",
			a:    time.Date(2026, time.April, 15, 12, 0, 0, 0, time.UTC),
			b:    time.Date(2027, time.April, 15, 12, 0, 0, 0, time.UTC),
			want: false,
		},
		{
			name: "Q4 2026 fiscal year vs Q4 2027 fiscal year",
			a:    time.Date(2027, time.January, 15, 12, 0, 0, 0, time.UTC),
			b:    time.Date(2028, time.January, 15, 12, 0, 0, 0, time.UTC),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cal.IsSameQuarter(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("IsSameQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuarter_ConfiguredLocationControlsSemantics(t *testing.T) {
	location := mustLoadLocation(t, "Asia/Kolkata")

	cal := mustNewCalendar(
		t,
		Config{
			Location:        location,
			WeekStart:       time.Monday,
			FiscalYearStart: time.April,
		},
	)

	// 2026-03-31 18:45 UTC = 2026-04-01 00:15 IST.
	value := time.Date(
		2026,
		time.March,
		31,
		18, 45, 0, 0,
		time.UTC,
	)

	got := cal.StartOfQuarter(value)

	want := time.Date(
		2026,
		time.April,
		1,
		0, 0, 0, 0,
		location,
	)

	if !got.Equal(want) {
		t.Errorf("StartOfQuarter() = %v, want %v", got, want)
	}

	if gotQuarter := cal.Quarter(value); gotQuarter != 1 {
		t.Errorf("Quarter() = %d, want 1", gotQuarter)
	}
}

func TestQuarter_IgnoresInputLocation(t *testing.T) {
	cal := mustNewCalendar(
		t,
		Config{
			Location:        time.UTC,
			WeekStart:       time.Monday,
			FiscalYearStart: time.April,
		},
	)

	_ = cal
}

func TestQuarter_ConfiguredLocationAcrossTimezones(t *testing.T) {
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

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)

			cal := mustNewCalendar(
				t,
				Config{
					Location:        location,
					WeekStart:       time.Monday,
					FiscalYearStart: time.April,
				},
			)

			value := time.Date(
				2026,
				time.August,
				15,
				12, 30, 45, 123456789,
				time.UTC,
			)

			start := cal.StartOfQuarter(value)

			want := time.Date(
				2026,
				time.July,
				1,
				0, 0, 0, 0,
				location,
			)

			if !start.Equal(want) {
				t.Errorf(
					"StartOfQuarter() = %v, want %v",
					start,
					want,
				)
			}

			if start.Location() != location {
				t.Errorf(
					"StartOfQuarter() location = %v, want %v",
					start.Location(),
					location,
				)
			}

			if got := cal.Quarter(value); got != 2 {
				t.Errorf("Quarter() = %d, want 2", got)
			}
		})
	}
}

func TestQuarter_InputLocationDoesNotAffectResult(t *testing.T) {
	cal := mustNewCalendar(t, Config{
		Location:        time.UTC,
		WeekStart:       time.Monday,
		FiscalYearStart: time.April,
	})

	locations := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, timezone := range locations {
		t.Run(timezone, func(t *testing.T) {
			inputLocation := mustLoadLocation(t, timezone)

			value := time.Date(
				2026,
				time.May,
				15,
				12, 30, 45, 123456789,
				inputLocation,
			)

			got := cal.StartOfQuarter(value)

			want := time.Date(
				2026,
				time.April,
				1,
				0, 0, 0, 0,
				time.UTC,
			)

			if !got.Equal(want) {
				t.Errorf(
					"StartOfQuarter() = %v, want %v",
					got,
					want,
				)
			}

			if got.Location() != time.UTC {
				t.Errorf(
					"StartOfQuarter() location = %v, want UTC",
					got.Location(),
				)
			}
		})
	}
}

func TestStartOfQuarter_PreservesConfiguredLocation(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)

			cal := mustNewCalendar(
				t,
				Config{
					Location:        location,
					WeekStart:       time.Monday,
					FiscalYearStart: time.April,
				},
			)

			value := time.Date(
				2026,
				time.August,
				15,
				12, 30, 45, 123456789,
				location,
			)

			got := cal.StartOfQuarter(value)

			if got.Location() != location {
				t.Errorf(
					"StartOfQuarter() location = %v, want %v",
					got.Location(),
					location,
				)
			}
		})
	}
}

func TestStartOfNextQuarter_PreservesLocation(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)

			cal := mustNewCalendar(
				t,
				Config{
					Location:        location,
					WeekStart:       time.Monday,
					FiscalYearStart: time.April,
				},
			)

			value := time.Date(
				2026,
				time.August,
				15,
				12, 30, 45, 123456789,
				location,
			)

			got := cal.StartOfNextQuarter(value)

			if got.Location() != location {
				t.Errorf(
					"StartOfNextQuarter() location = %v, want %v",
					got.Location(),
					location,
				)
			}
		})
	}
}

func TestStartOfPreviousQuarter_PreservesLocation(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)

			cal := mustNewCalendar(
				t,
				Config{
					Location:        location,
					WeekStart:       time.Monday,
					FiscalYearStart: time.April,
				},
			)

			value := time.Date(
				2026,
				time.August,
				15,
				12, 30, 45, 123456789,
				location,
			)

			got := cal.StartOfPreviousQuarter(value)

			if got.Location() != location {
				t.Errorf(
					"StartOfPreviousQuarter() location = %v, want %v",
					got.Location(),
					location,
				)
			}
		})
	}
}

func TestEndOfQuarter_PreservesLocation(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)

			cal := mustNewCalendar(
				t,
				Config{
					Location:        location,
					WeekStart:       time.Monday,
					FiscalYearStart: time.April,
				},
			)

			value := time.Date(
				2026,
				time.August,
				15,
				12, 30, 45, 123456789,
				location,
			)

			got := cal.EndOfQuarter(value)

			if got.Location() != location {
				t.Errorf(
					"EndOfQuarter() location = %v, want %v",
					got.Location(),
					location,
				)
			}
		})
	}
}

func TestQuarter_DST(t *testing.T) {
	timezones := []string{
		"America/New_York",
		"America/Los_Angeles",
		"Europe/London",
		"Europe/Berlin",
		"Australia/Sydney",
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)

			cal := mustNewCalendar(
				t,
				Config{
					Location:        location,
					WeekStart:       time.Monday,
					FiscalYearStart: time.January,
				},
			)

			// Q1 contains the spring DST transition in the selected
			// locations. The quarter boundary itself must still occur
			// at local midnight on April 1.
			value := time.Date(
				2026,
				time.March,
				15,
				12, 0, 0, 0,
				location,
			)

			start := cal.StartOfQuarter(value)
			next := cal.StartOfNextQuarter(value)

			wantStart := time.Date(
				2026,
				time.January,
				1,
				0, 0, 0, 0,
				location,
			)

			wantNext := time.Date(
				2026,
				time.April,
				1,
				0, 0, 0, 0,
				location,
			)

			if !start.Equal(wantStart) {
				t.Errorf(
					"StartOfQuarter() = %v, want %v",
					start,
					wantStart,
				)
			}

			if !next.Equal(wantNext) {
				t.Errorf(
					"StartOfNextQuarter() = %v, want %v",
					next,
					wantNext,
				)
			}
		})
	}
}

func TestQuarter_DST_FallBack(t *testing.T) {
	timezones := []string{
		"America/New_York",
		"America/Los_Angeles",
		"Europe/London",
		"Europe/Berlin",
		"Australia/Sydney",
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)

			cal := mustNewCalendar(
				t,
				Config{
					Location:        location,
					WeekStart:       time.Monday,
					FiscalYearStart: time.January,
				},
			)

			value := time.Date(
				2026,
				time.November,
				1,
				12, 0, 0, 0,
				location,
			)

			start := cal.StartOfQuarter(value)

			want := time.Date(
				2026,
				time.October,
				1,
				0, 0, 0, 0,
				location,
			)

			if !start.Equal(want) {
				t.Errorf(
					"StartOfQuarter() = %v, want %v",
					start,
					want,
				)
			}

			next := cal.StartOfNextQuarter(value)

			wantNext := time.Date(
				2027,
				time.January,
				1,
				0, 0, 0, 0,
				location,
			)

			if !next.Equal(wantNext) {
				t.Errorf(
					"StartOfNextQuarter() = %v, want %v",
					next,
					wantNext,
				)
			}
		})
	}
}

func TestQuarter_HalfOpenInterval(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)

			cal := mustNewCalendar(
				t,
				Config{
					Location:        location,
					WeekStart:       time.Monday,
					FiscalYearStart: time.April,
				},
			)

			value := time.Date(
				2026,
				time.August,
				15,
				12, 0, 0, 0,
				location,
			)

			start := cal.StartOfQuarter(value)
			next := cal.StartOfNextQuarter(value)
			end := cal.EndOfQuarter(value)

			if !start.Before(next) {
				t.Fatalf(
					"StartOfQuarter() = %v, StartOfNextQuarter() = %v",
					start,
					next,
				)
			}

			if !end.Before(next) {
				t.Fatalf(
					"EndOfQuarter() = %v, StartOfNextQuarter() = %v",
					end,
					next,
				)
			}

			if !end.Add(time.Nanosecond).Equal(next) {
				t.Errorf(
					"EndOfQuarter() + 1ns = %v, want %v",
					end.Add(time.Nanosecond),
					next,
				)
			}

			if !start.Before(value) || !next.After(value) {
				t.Errorf(
					"value %v is not inside [%v, %v)",
					value,
					start,
					next,
				)
			}
		})
	}
}

func TestQuarter_ConsecutiveBoundaries(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)

			for fiscalStart := time.January; fiscalStart <= time.December; fiscalStart++ {
				cal := mustNewCalendar(
					t,
					Config{
						Location:        location,
						WeekStart:       time.Monday,
						FiscalYearStart: fiscalStart,
					},
				)

				value := time.Date(
					2026,
					time.August,
					15,
					12, 0, 0, 0,
					location,
				)

				start := cal.StartOfQuarter(value)
				next := cal.StartOfNextQuarter(value)
				previous := cal.StartOfPreviousQuarter(value)

				if !previous.Before(start) {
					t.Errorf(
						"previous quarter %v is not before current quarter %v",
						previous,
						start,
					)
				}

				if !start.Before(next) {
					t.Errorf(
						"current quarter %v is not before next quarter %v",
						start,
						next,
					)
				}

				if !cal.StartOfNextQuarter(
					time.Date(
						start.Year(),
						start.Month(),
						start.Day(),
						12, 0, 0, 0,
						location,
					),
				).Equal(next) {
					t.Errorf(
						"next quarter boundary is inconsistent: got %v, want %v",
						cal.StartOfNextQuarter(value),
						next,
					)
				}
			}
		})
	}
}

func TestQuarter_QuarterAlwaysBetweenOneAndFour(t *testing.T) {
	timezones := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)

			for fiscalStart := time.January; fiscalStart <= time.December; fiscalStart++ {
				cal := mustNewCalendar(
					t,
					Config{
						Location:        location,
						WeekStart:       time.Monday,
						FiscalYearStart: fiscalStart,
					},
				)

				for month := time.January; month <= time.December; month++ {
					value := time.Date(
						2026,
						month,
						15,
						12, 0, 0, 0,
						location,
					)

					quarter := cal.Quarter(value)

					if quarter < 1 || quarter > 4 {
						t.Errorf(
							"Quarter(%v) = %d, want value in [1,4]",
							value,
							quarter,
						)
					}
				}
			}
		})
	}
}

func TestQuarter_AllMonthsMapToExactlyOneQuarter(t *testing.T) {
	cal := mustNewCalendar(t, Config{
		Location:        time.UTC,
		WeekStart:       time.Monday,
		FiscalYearStart: time.April,
	})

	for month := time.January; month <= time.December; month++ {
		value := time.Date(
			2026,
			month,
			15,
			12, 0, 0, 0,
			time.UTC,
		)

		quarter := cal.Quarter(value)

		expected := ((int(month)-int(time.April)+12)%12)/3 + 1

		if quarter != expected {
			t.Errorf(
				"Quarter(%v) = %d, want %d",
				value,
				quarter,
				expected,
			)
		}
	}
}

func TestQuarter_SameQuarterAcrossDifferentInputLocations(t *testing.T) {
	cal := mustNewCalendar(t, Config{
		Location:        time.UTC,
		WeekStart:       time.Monday,
		FiscalYearStart: time.April,
	})

	locations := []string{
		"UTC",
		"Asia/Kolkata",
		"America/New_York",
		"Europe/London",
		"Australia/Sydney",
	}

	base := time.Date(
		2026,
		time.May,
		15,
		12, 0, 0, 0,
		time.UTC,
	)

	wantStart := cal.StartOfQuarter(base)

	for _, timezone := range locations {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)
			value := base.In(location)

			got := cal.StartOfQuarter(value)

			if !got.Equal(wantStart) {
				t.Errorf(
					"StartOfQuarter() = %v, want %v",
					got,
					wantStart,
				)
			}
		})
	}
}

func TestAddQuarters(t *testing.T) {
	tests := []struct {
		name     string
		value    time.Time
		quarters int
		want     time.Time
	}{
		{
			name: "zero",
			value: time.Date(
				2026, time.May, 15,
				12, 30, 45, 123456789,
				time.UTC,
			),
			quarters: 0,
			want: time.Date(
				2026, time.May, 15,
				12, 30, 45, 123456789,
				time.UTC,
			),
		},
		{
			name: "positive one",
			value: time.Date(
				2026, time.January, 15,
				12, 30, 45, 123456789,
				time.UTC,
			),
			quarters: 1,
			want: time.Date(
				2026, time.April, 15,
				12, 30, 45, 123456789,
				time.UTC,
			),
		},
		{
			name: "positive multiple",
			value: time.Date(
				2026, time.January, 15,
				12, 30, 45, 123456789,
				time.UTC,
			),
			quarters: 4,
			want: time.Date(
				2027, time.January, 15,
				12, 30, 45, 123456789,
				time.UTC,
			),
		},
		{
			name: "negative one",
			value: time.Date(
				2026, time.January, 15,
				12, 30, 45, 123456789,
				time.UTC,
			),
			quarters: -1,
			want: time.Date(
				2025, time.October, 15,
				12, 30, 45, 123456789,
				time.UTC,
			),
		},
		{
			name: "negative multiple",
			value: time.Date(
				2026, time.January, 15,
				12, 30, 45, 123456789,
				time.UTC,
			),
			quarters: -4,
			want: time.Date(
				2025, time.January, 15,
				12, 30, 45, 123456789,
				time.UTC,
			),
		},
		{
			name: "year boundary forward",
			value: time.Date(
				2026, time.November, 15,
				12, 0, 0, 0,
				time.UTC,
			),
			quarters: 1,
			want: time.Date(
				2027, time.February, 15,
				12, 0, 0, 0,
				time.UTC,
			),
		},
		{
			name: "year boundary backward",
			value: time.Date(
				2026, time.February, 15,
				12, 0, 0, 0,
				time.UTC,
			),
			quarters: -1,
			want: time.Date(
				2025, time.November, 15,
				12, 0, 0, 0,
				time.UTC,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddQuarters(tt.value, tt.quarters)

			if !got.Equal(tt.want) {
				t.Errorf("AddQuarters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddQuarters_MatchesAddDate(t *testing.T) {
	values := []time.Time{
		time.Date(2026, time.January, 31, 12, 30, 0, 0, time.UTC),
		time.Date(2026, time.March, 31, 23, 59, 59, 0, time.UTC),
		time.Date(2026, time.May, 31, 8, 15, 0, 0, time.UTC),
		time.Date(2026, time.August, 31, 18, 45, 0, 0, time.UTC),
		time.Date(2026, time.October, 31, 4, 20, 0, 0, time.UTC),
		time.Date(2026, time.December, 31, 23, 59, 59, 999999999, time.UTC),
	}

	quarters := []int{-8, -4, -1, 0, 1, 2, 4, 8}

	for _, value := range values {
		for _, q := range quarters {
			t.Run(
				value.Format("2006-01-02")+"_"+strconv.Itoa(q),
				func(t *testing.T) {
					got := AddQuarters(value, q)
					want := value.AddDate(0, q*3, 0)

					if !got.Equal(want) {
						t.Errorf(
							"AddQuarters(%v, %d) = %v, want %v",
							value,
							q,
							got,
							want,
						)
					}
				},
			)
		}
	}
}

func TestAddQuarters_PreservesLocation(t *testing.T) {
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
		12, 30, 45, 123456789,
		time.UTC,
	)

	for _, timezone := range timezones {
		t.Run(timezone, func(t *testing.T) {
			location := mustLoadLocation(t, timezone)
			localizedValue := value.In(location)

			got := AddQuarters(localizedValue, 1)

			if got.Location() != location {
				t.Errorf(
					"AddQuarters() location = %v, want %v",
					got.Location(),
					location,
				)
			}
		})
	}
}

func TestAddQuarters_PreservesTimeComponents(t *testing.T) {
	value := time.Date(
		2026,
		time.May,
		17,
		14, 23, 51,
		987654321,
		time.UTC,
	)

	got := AddQuarters(value, 1)

	if got.Hour() != value.Hour() ||
		got.Minute() != value.Minute() ||
		got.Second() != value.Second() ||
		got.Nanosecond() != value.Nanosecond() {
		t.Errorf(
			"AddQuarters() changed time components: got %v, want hour=%d minute=%d second=%d nanosecond=%d",
			got,
			value.Hour(),
			value.Minute(),
			value.Second(),
			value.Nanosecond(),
		)
	}
}

func TestQuarter_YearBoundary(t *testing.T) {
	cal := mustNewCalendar(t, Config{
		Location:        time.UTC,
		WeekStart:       time.Monday,
		FiscalYearStart: time.April,
	})

	tests := []struct {
		name        string
		value       time.Time
		wantStart   time.Time
		wantNext    time.Time
		wantQuarter int
	}{
		{
			name: "Q3 ending December",
			value: time.Date(
				2026, time.December, 31,
				23, 59, 59, 999999999,
				time.UTC,
			),
			wantStart: time.Date(
				2026, time.October, 1,
				0, 0, 0, 0,
				time.UTC,
			),
			wantNext: time.Date(
				2027, time.January, 1,
				0, 0, 0, 0,
				time.UTC,
			),
			wantQuarter: 3,
		},
		{
			name: "Q4 ending March",
			value: time.Date(
				2027, time.March, 31,
				23, 59, 59, 999999999,
				time.UTC,
			),
			wantStart: time.Date(
				2027, time.January, 1,
				0, 0, 0, 0,
				time.UTC,
			),
			wantNext: time.Date(
				2027, time.April, 1,
				0, 0, 0, 0,
				time.UTC,
			),
			wantQuarter: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cal.StartOfQuarter(tt.value); !got.Equal(tt.wantStart) {
				t.Errorf(
					"StartOfQuarter() = %v, want %v",
					got,
					tt.wantStart,
				)
			}

			if got := cal.StartOfNextQuarter(tt.value); !got.Equal(tt.wantNext) {
				t.Errorf(
					"StartOfNextQuarter() = %v, want %v",
					got,
					tt.wantNext,
				)
			}

			if got := cal.Quarter(tt.value); got != tt.wantQuarter {
				t.Errorf(
					"Quarter() = %d, want %d",
					got,
					tt.wantQuarter,
				)
			}
		})
	}
}

func TestQuarter_AllFiscalStartMonths_BoundaryInvariant(t *testing.T) {
	for fiscalStart := time.January; fiscalStart <= time.December; fiscalStart++ {
		t.Run(fiscalStart.String(), func(t *testing.T) {
			cal := mustNewCalendar(t, Config{
				Location:        time.UTC,
				WeekStart:       time.Monday,
				FiscalYearStart: fiscalStart,
			})

			for quarter := 1; quarter <= 4; quarter++ {
				// Calculate the first month of this fiscal quarter.
				startMonth := int(fiscalStart) + (quarter-1)*3
				year := 2026

				if startMonth > 12 {
					startMonth -= 12
					year++
				}

				start := time.Date(
					year,
					time.Month(startMonth),
					1,
					0, 0, 0, 0,
					time.UTC,
				)

				t.Run("Q"+strconv.Itoa(quarter), func(t *testing.T) {
					// The calculated boundary must be recognized as
					// the start of the same quarter.
					if got := cal.StartOfQuarter(start); !got.Equal(start) {
						t.Errorf(
							"StartOfQuarter() = %v, want %v",
							got,
							start,
						)
					}

					// The quarter number must match the expected
					// fiscal quarter.
					if got := cal.Quarter(start); got != quarter {
						t.Errorf(
							"Quarter() = %d, want %d",
							got,
							quarter,
						)
					}

					next := cal.StartOfNextQuarter(start)

					// The next quarter must begin after the current
					// quarter.
					if !next.After(start) {
						t.Errorf(
							"StartOfNextQuarter() = %v, want after %v",
							next,
							start,
						)
					}

					// Moving backwards from the next quarter must
					// return exactly to the current quarter.
					if got := cal.StartOfPreviousQuarter(next); !got.Equal(start) {
						t.Errorf(
							"StartOfPreviousQuarter(StartOfNextQuarter()) = %v, want %v",
							got,
							start,
						)
					}

					// The end of the quarter must be exactly one
					// nanosecond before the next quarter.
					end := cal.EndOfQuarter(start)

					if !end.Add(time.Nanosecond).Equal(next) {
						t.Errorf(
							"EndOfQuarter() + 1ns = %v, want %v",
							end.Add(time.Nanosecond),
							next,
						)
					}

					// The next quarter must be recognized as the
					// following fiscal quarter.
					nextQuarter := quarter + 1
					if nextQuarter > 4 {
						nextQuarter = 1
					}

					if got := cal.Quarter(next); got != nextQuarter {
						t.Errorf(
							"Quarter(next) = %d, want %d",
							got,
							nextQuarter,
						)
					}
				})
			}
		})
	}
}
