package clock

import (
	"sync"
	"testing"
	"time"
)

func TestSystemNow(t *testing.T) {
	t.Parallel()

	before := time.Now().UTC()
	now := NewSystem().Now()
	after := time.Now().UTC()

	if now.Before(before) {
		t.Fatalf("expected Now() to be after or equal to before, got %v", now)
	}

	if now.After(after) {
		t.Fatalf("expected Now() to be before or equal to after, got %v", now)
	}

	if now.Location() != time.UTC {
		t.Fatalf("expected UTC location, got %v", now.Location())
	}
}

func TestStaticNow(t *testing.T) {
	t.Parallel()

	input := time.Date(
		2026,
		time.August,
		31,
		12,
		30,
		45,
		123456789,
		time.FixedZone("IST", 5*60*60+30*60),
	)

	clock := NewStatic(input)

	got := clock.Now()
	want := input.UTC()

	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}

	if got.Location() != time.UTC {
		t.Fatalf("expected UTC location, got %v", got.Location())
	}
}

func TestStaticNowReturnsSameTime(t *testing.T) {
	t.Parallel()

	input := time.Date(
		2026,
		time.August,
		31,
		12,
		30,
		45,
		123456789,
		time.UTC,
	)

	c := NewStatic(input)

	first := c.Now()
	second := c.Now()

	if !first.Equal(second) {
		t.Fatalf("expected repeated Now() calls to return the same time, got %v and %v", first, second)
	}
}

func TestControlledNow(t *testing.T) {
	t.Parallel()

	input := time.Date(
		2026,
		time.August,
		31,
		12,
		30,
		0,
		0,
		time.UTC,
	)

	c := NewControlled(input)

	if got := c.Now(); !got.Equal(input) {
		t.Fatalf("expected %v, got %v", input, got)
	}

	if c.Now().Location() != time.UTC {
		t.Fatalf("expected UTC location, got %v", c.Now().Location())
	}
}

func TestControlledSet(t *testing.T) {
	t.Parallel()

	initial := time.Date(
		2026,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	updated := time.Date(
		2026,
		time.August,
		31,
		12,
		30,
		0,
		0,
		time.FixedZone("IST", 5*60*60+30*60),
	)

	c := NewControlled(initial)
	c.Set(updated)

	got := c.Now()
	want := updated.UTC()

	if !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}

	if got.Location() != time.UTC {
		t.Fatalf("expected UTC location, got %v", got.Location())
	}
}

func TestControlledAdvance(t *testing.T) {
	t.Parallel()

	initial := time.Date(
		2026,
		time.August,
		31,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	c := NewControlled(initial)

	c.Advance(2 * time.Hour)

	want := initial.Add(2 * time.Hour)

	if got := c.Now(); !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestControlledSequentialOperations(t *testing.T) {
	t.Parallel()

	initial := time.Date(
		2026,
		time.January,
		1,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	c := NewControlled(initial)

	c.Advance(2 * time.Hour)
	c.Set(initial.Add(5 * time.Hour))
	c.Advance(-30 * time.Minute)

	want := initial.Add(4*time.Hour + 30*time.Minute)

	if got := c.Now(); !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestControlledPreservesNanosecondPrecision(t *testing.T) {
	t.Parallel()

	input := time.Date(
		2026,
		time.August,
		31,
		12,
		30,
		45,
		123456789,
		time.UTC,
	)

	c := NewControlled(input)

	if got := c.Now(); !got.Equal(input) {
		t.Fatalf("expected %v, got %v", input, got)
	}
}

func TestControlledAdvanceZero(t *testing.T) {
	t.Parallel()

	initial := time.Date(
		2026,
		time.August,
		31,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	c := NewControlled(initial)

	c.Advance(0)

	if got := c.Now(); !got.Equal(initial) {
		t.Fatalf("expected %v, got %v", initial, got)
	}
}

func TestControlledAdvanceNegative(t *testing.T) {
	t.Parallel()

	initial := time.Date(
		2026,
		time.August,
		31,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	c := NewControlled(initial)

	c.Advance(-30 * time.Minute)

	want := initial.Add(-30 * time.Minute)

	if got := c.Now(); !got.Equal(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestControlledConcurrentAccess(t *testing.T) {
	t.Parallel()

	c := NewControlled(time.Date(
		2026,
		time.August,
		31,
		12,
		0,
		0,
		0,
		time.UTC,
	))

	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()

		for i := 0; i < iterations; i++ {
			_ = c.Now()
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < iterations; i++ {
			c.Set(time.Date(
				2026,
				time.August,
				31,
				12,
				0,
				i,
				0,
				time.UTC,
			))
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < iterations; i++ {
			c.Advance(time.Second)
		}
	}()

	wg.Wait()
}

func TestClockImplementations(t *testing.T) {
	t.Parallel()

	var _ Clock = System{}
	var _ Clock = Static{}
	var _ Clock = (*Controlled)(nil)
}
