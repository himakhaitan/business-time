# business-time

> Business time handling for Go, without making developers think about time.

`business-time` is a collection of small, composable Go packages for building software around **business time**.

Calendars. Business days. Working hours. Holidays. Schedules. Durations. Deadlines. Reporting periods.

The goal is simple:

> **When you're building business software, time should be infrastructure — not something you have to keep solving.**

`business-time` builds on Go's standard `time` package rather than replacing it. You keep using `time.Time`, `time.Duration`, and `time.Location`, while `business-time` provides the higher-level semantics that business applications repeatedly need.

---

## Why business-time?

Time looks simple until software needs to answer questions like:

* What day is this timestamp in the customer's timezone?
* When does the business day start?
* Is this timestamp inside the same week as another timestamp?
* When does the fiscal quarter end?
* Is today a working day?
* Is Monday a holiday?
* How many business hours are between two timestamps?
* When is an SLA deadline due?
* What happens when a deadline falls on a weekend?
* What are the working hours for this organization?
* What timezone defines the business calendar?
* Does daylight saving time change the result?
* Which fiscal year does this date belong to?
* How should a reporting period be calculated?

These are not application-specific problems.

They appear in:

* payments
* banking
* accounting
* payroll
* invoicing
* subscriptions
* logistics
* CRM systems
* HR systems
* ecommerce
* marketplaces
* customer support
* enterprise software
* reporting
* analytics
* workflow systems
* scheduling
* SLAs
* notifications
* billing
* compliance systems

Most applications end up implementing some version of these concepts themselves.

`business-time` aims to provide those concepts as reusable, well-defined building blocks.

---

## Design philosophy

### 1. Built on Go's `time` package

`business-time` does not introduce another datetime type.

You continue to use:

```go
time.Time
time.Duration
time.Location
time.Weekday
time.Month
```

The library adds business and calendar semantics around those primitives.

This keeps the API familiar and makes interoperability with the Go standard library straightforward.

---

### 2. Explicit timezone semantics

Timezones are not an implementation detail.

A calendar needs to know which location defines its days, weeks, months, and years.

For that reason, calendar configuration is explicit:

```go
cal, err := calendar.New(calendar.Config{
    Location:        location,
    WeekStart:       time.Monday,
    FiscalYearStart: time.January,
})
```

There is no implicit dependence on the machine's local timezone for calendar calculations.

This makes behavior predictable across:

* developer machines
* CI environments
* containers
* cloud infrastructure
* databases
* distributed services
* multiple geographic regions

---

### 3. Calendar time is not duration time

A calendar day is not necessarily 24 elapsed hours.

A calendar week is not necessarily 168 elapsed hours.

A month does not have a fixed duration.

A year does not have a fixed number of days.

This distinction becomes particularly important around daylight-saving transitions.

For example, adding one calendar day is different from adding:

```go
24 * time.Hour
```

`business-time` uses calendar arithmetic where calendar semantics are intended.

---

### 4. Half-open intervals by default

Where possible, periods are naturally represented as:

```text
[start, end)
```

For example:

```text
[StartOfDay, StartOfNextDay)
```

rather than:

```text
[StartOfDay, EndOfDay]
```

This avoids the need to reason about the "last possible instant" of a period and works naturally with comparisons, databases, ranges, and interval arithmetic.

`EndOfDay`, `EndOfWeek`, `EndOfMonth`, etc. are provided where a closed boundary is useful, but the beginning of the next period is the preferred exclusive boundary for interval calculations.

---

### 5. Small packages, focused responsibilities

`business-time` is intentionally not one giant utility package.

Each package should answer a particular category of time-related questions.

```text
business-time/
│
├── calendar/     Calendar dates and periods
├── clock/        Controllable current-time source
│
├── business/     Business-day and holiday rules
├── schedule/     Working hours and recurring schedules
├── duration/     Business-aware duration calculations
├── deadline/     Business deadlines and SLA calculations
└── period/       Fiscal and reporting periods
```

Only the packages that solve your problem need to be imported.

---

# Packages

## `calendar`

Calendar semantics and calendar arithmetic.

The `calendar` package provides:

* calendar days
* weeks
* months
* quarters
* fiscal years
* period boundaries
* weekday calculations
* day/month/year calculations
* calendar arithmetic

Example:

```go
package main

import (
    "fmt"
    "time"

    "github.com/himakhaitan/business-time/calendar"
)

func main() {
    cal, err := calendar.New(calendar.Config{
        Location:        time.UTC,
        WeekStart:       time.Monday,
        FiscalYearStart: time.January,
    })
    if err != nil {
        panic(err)
    }

    now := time.Now()

    fmt.Println(cal.StartOfDay(now))
    fmt.Println(cal.StartOfWeek(now))
    fmt.Println(cal.StartOfMonth(now))
    fmt.Println(cal.StartOfQuarter(now))
    fmt.Println(cal.StartOfYear(now))
}
```

`calendar` deliberately does **not** know what a business day is.

It does not contain:

* holidays
* working hours
* weekend policies
* employee schedules
* office hours
* business closures

Those concepts belong to higher-level packages.

---

## `clock`

A small abstraction for code that needs the current time.

The purpose of `clock` is to make code that depends on "now" easier to test and control.

For example, production code can use a real clock while tests can provide a deterministic clock.

This package is intentionally independent from the rest of the library.

Calendar calculations should receive their timestamps explicitly rather than secretly calling `time.Now()`.

---

# Coming soon

The initial release focuses on the foundational calendar layer.

The long-term goal is to cover the major categories of time semantics commonly required by business software.

## `business`

Business-day semantics.

Planned capabilities include:

* business-day detection
* configurable weekends
* holidays
* holiday calendars
* working-day calculations
* next business day
* previous business day
* adding business days
* counting business days
* business-day ranges

Example API direction:

```go
business.IsBusinessDay(t)
business.NextBusinessDay(t)
business.PreviousBusinessDay(t)
business.AddBusinessDays(t, 5)
```

The exact API will evolve as the package is designed.

---

## `schedule`

Working-time and recurring schedules.

Planned capabilities include:

* working hours
* multiple working intervals per day
* weekday schedules
* breaks
* recurring schedules
* timezone-aware schedules
* schedule boundaries
* open/closed checks

For example, a schedule might represent:

```text
Monday    09:00–12:00, 13:00–18:00
Tuesday   09:00–12:00, 13:00–18:00
Wednesday 09:00–12:00, 13:00–18:00
Thursday  09:00–12:00, 13:00–18:00
Friday    09:00–12:00, 13:00–17:00
Saturday  Closed
Sunday    Closed
```

Schedules and business-day rules are intentionally separate concepts.

A day can be a business day without every instant of that day being a working instant.

---

## `duration`

Business-aware duration calculations.

Planned capabilities include calculations such as:

```text
elapsed duration
calendar duration
business duration
working duration
```

Examples:

* How many business days are between two dates?
* How many working hours elapsed?
* How much scheduled working time remains?
* How much business time passed between two events?

The package will distinguish elapsed time from calendar time and business/working time rather than treating them as interchangeable.

---

## `deadline`

Business deadlines and SLA calculations.

Planned capabilities include:

* business deadlines
* working-hour deadlines
* SLA expiration
* deadline adjustment
* deadline checks
* remaining business time
* next valid deadline
* weekend/holiday adjustment

A future API might allow concepts such as:

```text
Start:
    Monday 16:00

SLA:
    8 working hours

Schedule:
    Mon–Fri, 09:00–18:00

Result:
    Tuesday 16:00
```

The exact API will be designed around composability with `business` and `schedule`.

---

## `period`

Fiscal and reporting periods.

Planned capabilities include:

* fiscal periods
* fiscal quarters
* reporting periods
* accounting periods
* period labels
* period comparisons
* period boundaries
* custom period definitions

The existing `calendar` package provides the fundamental fiscal-year and quarter semantics that this package can build upon.

---

# Timezone and DST behavior

Timezone correctness is one of the core goals of this project.

Consider:

```go
start := cal.StartOfDay(value)
next := cal.StartOfNextDay(value)
```

The difference between these two timestamps is not necessarily:

```go
24 * time.Hour
```

A local calendar day can be shorter or longer because of timezone transitions.

That is why `business-time` prefers operations such as:

```go
AddDate(0, 0, 1)
```

for calendar-day arithmetic instead of:

```go
Add(24 * time.Hour)
```

The same principle applies to weeks, months, quarters, and years.

---

# Configuration

Calendar configuration is intentionally small.

```go
type Config struct {
    Location        *time.Location
    WeekStart       time.Weekday
    FiscalYearStart time.Month
}
```

### Location

Defines the timezone used to interpret calendar boundaries.

### WeekStart

Defines which weekday begins a calendar week.

For example:

```go
WeekStart: time.Monday
```

or:

```go
WeekStart: time.Sunday
```

### FiscalYearStart

Defines the first month of the fiscal year.

For example:

```go
FiscalYearStart: time.April
```

means:

```text
Q1 = April–June
Q2 = July–September
Q3 = October–December
Q4 = January–March
```

---

# What business-time does not try to do

`business-time` is intentionally focused on time semantics.

It is not intended to become:

* a database ORM
* a scheduling UI
* a cron replacement
* a calendar application
* a timezone database
* a payroll system
* an accounting system
* a booking system
* a workflow engine

The library provides the time primitives those systems can build upon.

---

# Compatibility with Go

`business-time` is designed to work naturally with Go's standard library.

You can freely combine it with:

```go
time.Time
time.Duration
time.Location
time.Weekday
time.Month
```

No custom datetime representation is required.

For example:

```go
start := cal.StartOfMonth(order.CreatedAt)
end := cal.StartOfNextMonth(order.CreatedAt)

rows, err := db.QueryContext(
    ctx,
    query,
    start,
    end,
)
```

The resulting boundaries can be passed directly to APIs, databases, queues, jobs, and other Go libraries that already understand `time.Time`.

---

# Installation

Install the module with:

```bash
go get github.com/himakhaitan/business-time
```

Then import the package you need:

```go
import "github.com/himakhaitan/business-time/calendar"
```

---

# Requirements

* Go 1.XX or later

The project uses the Go standard library and is designed to keep external dependencies to a minimum.

---

# Project status

`business-time` is being developed incrementally.

The first release establishes the foundational calendar layer.

Current:

* [x] Calendar configuration
* [x] Timezone-aware day operations
* [x] Week operations
* [x] Month operations
* [x] Fiscal quarter operations
* [x] Fiscal year operations
* [x] Calendar arithmetic
* [x] DST-aware calendar boundaries
* [x] Explicit timezone semantics
* [x] Extensive calendar tests

Coming soon:

* [ ] Business days
* [ ] Holiday calendars
* [ ] Working schedules
* [ ] Business durations
* [ ] SLA/deadline calculations
* [ ] Reporting periods
* [ ] More business-time primitives

The API is expected to grow around focused packages rather than a single catch-all API.

---

# Versioning

`business-time` follows semantic versioning.

Once the project reaches a stable `v1.x` API, backwards-incompatible API changes will require a new major version.

The project aims to keep individual packages small and their public APIs deliberate so that applications can depend on only the functionality they need.

---

# Philosophy

Business software should not need to reinvent time.

A developer working on:

```text
orders
payments
subscriptions
invoices
support tickets
SLAs
reports
payroll
workflows
logistics
```

should not have to repeatedly answer:

> "What exactly does a business day mean here?"

or:

> "What happens when daylight saving time changes?"

or:

> "Which timezone defines this date?"

Those decisions should be explicit, reusable, and tested.

That's what `business-time` is trying to provide.

> **Build your business logic. Let business-time handle the time.**

---

# Contributing

Contributions are welcome.

Before opening a pull request:

```bash
gofmt -w .
go test ./...
go test -race ./...
go vet ./...
```

Changes to exported APIs should include appropriate documentation and tests.

Time-related behavior should include tests for relevant timezone, calendar-boundary, DST, and leap-year cases where applicable.

---

# License

See [LICENSE](LICENSE).
