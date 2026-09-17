# business-time

> Business time handling for Go — business days, working hours, holidays, calendars, schedules, deadlines, and date calculations.

`business-time` is a collection of small, composable Go packages for building software around **business time**.

It provides higher-level time semantics commonly needed in business applications, including:

* **Business days** and weekdays
* **Working hours** and business periods
* **Holidays** and holiday calendars
* **Schedules** and recurring periods
* **Business durations** and elapsed time
* **Deadlines** and date calculations
* **Reporting periods** and calendar-aware calculations

`business-time` builds on Go's standard [`time`](https://pkg.go.dev/time) package rather than replacing it. You keep using `time.Time`, `time.Duration`, and `time.Location` while `business-time` provides the business rules and calculations around them.

## Why business-time?

Time looks simple until software needs to reason about **business rules, calendars, time zones, and working periods**.

Sooner or later, business applications need to answer questions like:

* What day is this timestamp in the customer's timezone?
* Is this a business day?
* Is Monday a holiday?
* When does the business day start and end?
* What are this organization's working hours?
* How many business hours are between two timestamps?
* When is an SLA deadline due?
* What happens when a deadline falls on a weekend or holiday?
* Which timezone defines the business calendar?
* Does daylight saving time affect the calculation?
* Which week, month, quarter, or fiscal year does this date belong to?
* When does a fiscal quarter end?
* How should a reporting period be calculated?

These aren't problems unique to one type of application. They appear across:

* Payments and banking
* Accounting and payroll
* Invoicing and billing
* Subscriptions
* Logistics and ecommerce
* CRM and customer support
* HR and enterprise software
* Reporting and analytics
* Workflow and scheduling
* SLAs and notifications
* Compliance systems
* Marketplaces

Yet many applications end up implementing their own versions of business days, holidays, working hours, deadlines, calendars, and reporting periods — often with slightly different rules and edge-case behavior.

`business-time` aims to make these concepts **reusable, explicit, and composable building blocks** for Go applications.

Instead of repeatedly solving the same class of time problems, applications can define their business-time rules once and build on top of them.

> Business time should be infrastructure — not something you have to keep solving.

## Design philosophy

### 1. Built on Go's `time` package

`business-time` does not introduce another datetime type or replace Go's standard time primitives.

You continue to use:

```go
time.Time
time.Duration
time.Location
time.Weekday
time.Month
```

The library adds **business and calendar semantics** around these existing types.

This keeps the API familiar, preserves interoperability with the standard library and existing Go libraries, and avoids forcing developers to learn or convert between competing time representations.

### 2. Explicit timezone semantics

Timezones are part of the business rules, not an implementation detail.

A calendar needs to know which location defines its **days, weeks, months, years, and working periods**.

For that reason, timezone configuration is explicit:

```go
cal, err := calendar.New(calendar.Config{
    Location:        location,
    WeekStart:       time.Monday,
    FiscalYearStart: time.January,
})
```

Calendar calculations never implicitly depend on the machine's local timezone.

This makes behavior **deterministic and predictable** across:

* Developer machines
* CI environments
* Containers
* Cloud infrastructure
* Databases
* Distributed services
* Multiple geographic regions

The same calendar configuration should produce the same business-time semantics regardless of where the code runs.

### 3. Calendar time is not duration time

A **calendar interval** and an **elapsed duration** are not the same thing.

A calendar day is not necessarily 24 elapsed hours.

A calendar week is not necessarily 168 elapsed hours.

A month has no fixed duration.

A year does not have a fixed number of days.

This distinction matters whenever calculations cross timezone or daylight-saving transitions.

For example, these two operations express different intent:

```go
// Elapsed duration
t.Add(24 * time.Hour)

// Calendar arithmetic
// "the same local time on the next calendar day"
```

Adding `24 * time.Hour` means **24 elapsed hours**.

Adding one calendar day means **move to the next calendar date**, preserving the relevant calendar semantics.

`business-time` keeps these concepts distinct and uses **calendar arithmetic when calendar semantics are intended**, rather than treating every time calculation as a duration calculation.

### 4. Half-open intervals by default

Where possible, periods are represented using **half-open intervals**:

```text
[start, end)
```

The start is inclusive; the end is exclusive.

For example, a calendar day is naturally represented as:

```text
[StartOfDay, StartOfNextDay)
```

rather than:

```text
[StartOfDay, EndOfDay]
```

This avoids having to reason about the **"last possible instant"** of a period and makes interval calculations easier to compose.

It also works naturally with comparisons, database queries, ranges, and adjacent periods:

```text
[Monday, Tuesday)
[Tuesday, Wednesday)
```

These intervals meet at `Tuesday` without overlapping.

`EndOfDay`, `EndOfWeek`, `EndOfMonth`, and similar boundaries may still be useful when explicitly requested. However, for interval calculations, **the beginning of the next period is the preferred exclusive boundary**.

### 5. Small packages, focused responsibilities

`business-time` is intentionally not one giant utility package.

Each package focuses on a specific category of time-related problems and provides its own **configuration options** for adapting those semantics to different business requirements.

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

Each package can be configured independently, with more advanced configuration available when the default behavior is not sufficient for a particular business calendar or workflow.

This keeps responsibilities clear while allowing the individual building blocks to handle real-world business-time requirements.

Only the packages that solve your problem need to be imported.

# Packages

## `calendar`

**Calendar semantics and calendar arithmetic.**

The `calendar` package provides the foundational calendar operations used throughout `business-time`.

It handles:

* Calendar days
* Weeks
* Months
* Quarters
* Fiscal years
* Period boundaries
* Weekday calculations
* Day, month, and year calculations
* Calendar arithmetic

### Example

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

### What `calendar` does not handle

`calendar` deliberately focuses on **calendar semantics**, not business rules.

It does not define:

* Holidays
* Working hours
* Weekend policies
* Employee schedules
* Office hours
* Business closures

Those concepts belong to higher-level packages such as `business`, `schedule`, and `deadline`.

This separation keeps `calendar` predictable and reusable: it answers questions about **calendar time** without making assumptions about what constitutes a **business day**.

## `clock`

A small abstraction for code that needs access to the **current time**.

The `clock` package makes time-dependent code easier to **test, control, and reason about**.

Production code can use a real clock, while tests can provide a deterministic clock with a known point in time.

This makes it possible to test behavior that depends on "now" without relying on the system clock or introducing sleeps and timing-sensitive assertions.

The package is intentionally independent from the rest of `business-time`.

Calendar and business-time calculations should receive their timestamps explicitly rather than implicitly calling `time.Now()`.

This keeps time as an explicit dependency and makes calculations **deterministic, reproducible, and easier to test**.

# Coming soon

The initial release focuses on the **foundational calendar layer**.

The long-term goal is to provide a cohesive set of small, composable packages covering the major categories of time semantics commonly required by business software.

## `business`

**Business-day semantics.**

Planned capabilities include:

* Business-day detection
* Configurable weekends
* Holidays and holiday calendars
* Working-day calculations
* Next and previous business day
* Adding and counting business days
* Business-day ranges

The package will build on calendar semantics while introducing rules that define what constitutes a **business day**.

Possible API direction:

```go
business.IsBusinessDay(t)
business.NextBusinessDay(t)
business.PreviousBusinessDay(t)
business.AddBusinessDays(t, 5)
```

The exact API will evolve as the package is designed and tested.

## `schedule`

**Working-time and recurring schedules.**

Planned capabilities include:

* Working hours
* Multiple working intervals per day
* Weekday schedules
* Breaks
* Recurring schedules
* Timezone-aware schedules
* Schedule boundaries
* Open/closed checks

For example, a schedule might represent:

```text
Monday     09:00–12:00, 13:00–18:00
Tuesday    09:00–12:00, 13:00–18:00
Wednesday  09:00–12:00, 13:00–18:00
Thursday   09:00–12:00, 13:00–18:00
Friday     09:00–12:00, 13:00–17:00
Saturday   Closed
Sunday     Closed
```

Schedules and business-day rules are intentionally separate concepts.

A day can be a business day without every instant within that day being a working instant.

## `duration`

**Business-aware duration calculations.**

Planned capabilities include distinguishing between:

```text
Elapsed duration
Calendar duration
Business duration
Working duration
```

Potential use cases include:

* Counting business days between dates
* Calculating elapsed working hours
* Calculating scheduled working time remaining
* Measuring business time between events

The package will keep **elapsed time, calendar time, and business/working time** as distinct concepts rather than treating them as interchangeable.

## `deadline`

**Business deadlines and SLA calculations.**

Planned capabilities include:

* Business deadlines
* Working-hour deadlines
* SLA expiration
* Deadline adjustment
* Deadline checks
* Remaining business time
* Next valid deadline
* Weekend and holiday adjustments

For example:

```text
Start:
    Monday 16:00

SLA:
    8 working hours

Schedule:
    Monday–Friday, 09:00–18:00

Result:
    Tuesday 16:00
```

The package is intended to compose with `business` and `schedule`, allowing deadlines to be calculated against explicit business calendars and working schedules.

The exact API will be designed around these composable building blocks.

## `period`

**Fiscal and reporting periods.**

Planned capabilities include:

* Fiscal periods
* Fiscal quarters
* Reporting periods
* Accounting periods
* Period labels
* Period comparisons
* Period boundaries
* Custom period definitions

The existing `calendar` package provides the foundational calendar and fiscal-year semantics that this package can build upon.

# Timezone and DST behavior

Timezone correctness is a core goal of `business-time`.

A local calendar day represents a **calendar boundary**, not a fixed number of elapsed hours.

For example:

```go
start := cal.StartOfDay(value)
next := cal.StartOfNextDay(value)
```

The duration between these timestamps is not necessarily:

```go
24 * time.Hour
```

A local day can be shorter or longer when a timezone transition changes the UTC offset.

For the same reason, calendar-day arithmetic should use calendar operations when calendar semantics are intended:

```go
t.AddDate(0, 0, 1)
```

rather than treating a calendar day as:

```go
t.Add(24 * time.Hour)
```

`AddDate` expresses **"the same local time on the next calendar day"**, while `Add(24 * time.Hour)` expresses **"24 elapsed hours later."**

The same distinction applies to weeks, months, quarters, and years: **calendar periods should be calculated using calendar semantics, not assumed durations.**

# What business-time does not try to do

`business-time` is intentionally focused on **time semantics and reusable time primitives**.

It is not intended to become:

* A database ORM
* A scheduling UI
* A cron replacement
* A calendar application
* A timezone database
* A payroll system
* An accounting system
* A booking system
* A workflow engine

Instead, the library provides the **time and calendar primitives** that these systems can build upon.

This keeps the scope focused and prevents `business-time` from becoming an application framework.

# Compatibility with Go

`business-time` is designed to work naturally with Go's standard library.

It builds directly on familiar types:

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

The resulting `time.Time` values can be passed directly to databases, APIs, queues, jobs, and other Go libraries that already understand Go's standard time types.

# Installation

Install the module with:

```bash
go get github.com/himakhaitan/business-time
```

Then import the package you need:

```go
import "github.com/himakhaitan/business-time/calendar"
```

# Requirements

* Go 1.XX or later

`business-time` primarily uses the Go standard library and aims to keep external dependencies to a minimum.

# Project status

`business-time` is being developed incrementally.

The first release focuses on establishing the **foundational calendar layer**.

### Current

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

### Coming soon

* [ ] Business days
* [ ] Holiday calendars
* [ ] Working schedules
* [ ] Business durations
* [ ] SLA and deadline calculations
* [ ] Reporting periods
* [ ] More business-time primitives

The API is expected to grow through **focused packages with clear responsibilities**, rather than a single catch-all API.

# Versioning

`business-time` follows [Semantic Versioning](https://semver.org/).

Once the project reaches a stable `v1.x` API, backwards-incompatible API changes will require a new major version.

The project aims to keep individual packages small and their public APIs deliberate, allowing applications to depend only on the functionality they need.

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

should not have to repeatedly answer questions like:

> What exactly does a business day mean here?

> What happens when daylight saving time changes?

> Which timezone defines this date?

These decisions should be **explicit, reusable, and tested**.

That's what `business-time` aims to provide.

> **Build your business logic. Let `business-time` handle the time.**

# Contributing

Contributions are welcome — especially if you've ever had to fight with a calendar bug that should have been simple.

Found a bug, have an idea, or want to make business time a little less painful?

See the CONTRIBUTING guide to get started.

# License

See [LICENSE](LICENSE).

# Made By

Hi, I'm **Himanshu Khaitan** — a developer who enjoys building things, obsessing over details, and occasionally turning "there has to be a better way" into an open-source project.

`business-time` started while I was building a financial system myself.

As the system grew, I kept running into the same deceptively difficult problems — business days, calendar boundaries, timezones, deadlines, working time, and all the little edge cases that come with handling time in financial software.

I found myself solving the same class of problems again and again.

So I started pulling those solutions out into small, reusable building blocks.

That eventually became `business-time`.

The idea is simple:

> **Developers shouldn't have to reinvent business time every time they build business software.**

**Find me around the internet:**

* Instagram — [@hey.hima_](https://www.instagram.com/hey.hima_)
* GitHub — [@himakhaitan](https://github.com/himakhaitan)
* LinkedIn — [Himanshu Khaitan](https://www.linkedin.com/in/himakhaitan)

And if `business-time` saved you from another timezone-induced headache, you can always:
**[Buy me a coffee](https://buymeacoffee.com/himakhaitan)**
Every coffee helps fund more open-source experiments, questionable ideas, and hopefully useful software.

> *Built with Go, curiosity, and an unreasonable amount of respect for `time.Time`.*
