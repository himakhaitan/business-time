# Contributing to business-time

Thank you for contributing to `business-time`.

`business-time` is a Go library designed to make business-related time handling predictable, explicit, and easy to use. Time-related code is deceptively difficult: timezone boundaries, daylight-saving transitions, leap years, calendar arithmetic, fiscal periods, and ambiguous semantics can all produce subtle bugs.

Because of that, contributions to this project are held to a high standard of correctness and clarity.

This document explains how to contribute code, tests, documentation, and ideas to the project.

---

## Table of Contents

* [Before You Start](#before-you-start)
* [Project Philosophy](#project-philosophy)
* [Repository Structure](#repository-structure)
* [Development Requirements](#development-requirements)
* [Getting Started](#getting-started)
* [Making Changes](#making-changes)
* [API Design Guidelines](#api-design-guidelines)
* [Time and Timezone Guidelines](#time-and-timezone-guidelines)
* [Daylight-Saving Time](#daylight-saving-time)
* [Calendar Arithmetic](#calendar-arithmetic)
* [Testing Requirements](#testing-requirements)
* [Documentation Requirements](#documentation-requirements)
* [Commit Messages](#commit-messages)
* [Pull Requests](#pull-requests)
* [Backward Compatibility](#backward-compatibility)
* [What Usually Does Not Belong in This Project](#what-usually-does-not-belong-in-this-project)
* [Before Opening a Pull Request](#before-opening-a-pull-request)
* [Release Process](#release-process)
* [Questions and Discussions](#questions-and-discussions)

---

## Before You Start

For small fixes, you can usually open a pull request directly.

For larger changes, especially changes involving:

* new public APIs
* new packages
* changes to existing semantics
* timezone behavior
* calendar arithmetic
* business-day rules
* scheduling
* duration semantics
* breaking changes

please open an issue or discussion first.

The goal is to agree on the behavior and API before investing significant effort in implementation.

A good proposal should explain:

1. What problem are you trying to solve?
2. Why does the functionality belong in `business-time`?
3. Which package should own it?
4. What should the public API look like?
5. What are the timezone and DST semantics?
6. What happens at boundaries?
7. What alternatives were considered?

---

# Project Philosophy

The primary goal of `business-time` is:

> **Business time handling for Go, without making developers think about time.**

The library should provide small, predictable abstractions that allow developers to use only the aspect of time they need.

For example:

```go
calendar.StartOfMonth(...)
```

should have clearly defined semantics without requiring users to understand the implementation.

At the same time, the library must not hide important decisions.

In particular:

* timezone behavior must be explicit
* calendar boundaries must be well-defined
* duration and calendar arithmetic must not be confused
* APIs should have predictable behavior around DST
* configuration should represent real semantic decisions
* public APIs should remain small
* standard library types should be preferred where appropriate

`business-time` complements Go's `time` package. It does not attempt to replace it.

---

# Repository Structure

The project is organized into focused packages.

```text
business-time/
├── calendar/
├── clock/
├── go.mod
├── README.md
├── CHANGELOG.md
├── CONTRIBUTING.md
├── LICENSE
└── .github/
    └── workflows/
```

### `calendar`

Provides generic calendar operations such as:

* days
* weeks
* months
* quarters
* years
* calendar boundaries
* calendar arithmetic

`calendar` should not contain business-specific policies such as holidays or working hours.

### `clock`

Provides time-source abstractions for code that needs the current time.

The rest of the library should not unnecessarily depend on `clock`.

### Future packages

The project may eventually include packages such as:

```text
business/
schedule/
duration/
deadline/
period/
```

These should only be introduced when there is a clear, well-defined need.

Do not create packages simply because they appear on the roadmap.

---

# Development Requirements

You should have:

* Go installed
* Git installed

The project should be developed using a currently supported Go version.

Check the repository's `go.mod` for the minimum supported Go version.

Verify your installation:

```bash
go version
```

---

# Getting Started

Clone the repository:

```bash
git clone <repository-url>
cd business-time
```

Run the complete test suite:

```bash
go test ./...
```

Run the race detector:

```bash
go test -race ./...
```

Run `go vet`:

```bash
go vet ./...
```

Format the code:

```bash
gofmt -w .
```

Before committing, make sure formatting does not produce additional changes:

```bash
gofmt -l .
```

The command should produce no output.

---

# Making Changes

Before making a change, identify which package owns the behavior.

For example:

* timezone-aware calendar boundaries → `calendar`
* current-time abstraction → `clock`
* holidays/business days → `business`
* working intervals → `schedule`
* elapsed business duration → `duration`
* SLA/deadline calculations → `deadline`
* fiscal/reporting periods → `period`

Avoid placing unrelated functionality into an existing package simply because it is convenient.

## Prefer small changes

A pull request should ideally solve one problem.

Good:

```text
fix(calendar): handle fiscal quarter year boundary
```

Less desirable:

```text
add quarter support, refactor calendar, change timezone behavior,
rewrite tests, update README, and add a new package
```

If several changes are logically independent, split them into separate pull requests when practical.

---

# API Design Guidelines

Public APIs are long-term contracts.

Before adding an exported function, type, field, constant, or variable, ask:

1. Is this functionality necessary?
2. Is this the correct package?
3. Is the API obvious to a new user?
4. Can the behavior be expressed using existing Go types?
5. Does the API expose implementation details unnecessarily?
6. Is the behavior unambiguous around timezone boundaries?
7. Can the API be supported long-term without breaking users?

## Prefer Go standard library types

Use:

```go
time.Time
time.Duration
time.Location
time.Weekday
time.Month
```

where they naturally fit.

Avoid introducing custom date/time types unless there is a strong reason.

For example, do not create a custom `Date` type simply to represent a calendar date if `time.Time` is sufficient.

## Avoid unnecessary convenience APIs

Do not add several functions that provide slightly different spellings of the same operation.

Every exported API increases the project's long-term maintenance and compatibility requirements.

Prefer a small, coherent API surface.

---

# Configuration

Configuration should represent meaningful semantic decisions.

For example, a calendar may require:

```go
type Config struct {
    Location        *time.Location
    WeekStart       time.Weekday
    FiscalYearStart time.Month
}
```

These values affect the meaning of calendar operations and therefore belong in the calendar configuration.

Do not silently introduce environment-dependent behavior when an explicit configuration would make the semantics clearer.

In particular, avoid silently using:

```go
time.Local
```

when the timezone is an important part of the operation.

---

# Time and Timezone Guidelines

Timezone behavior is one of the most important aspects of this project.

When implementing a calendar operation, always ask:

> **Which timezone defines the calendar day?**

A `time.Time` value carries a location, but the calendar may have its own configured location.

For example:

```go
calendar, _ := calendar.New(calendar.Config{
    Location: time.UTC,
})
```

The calendar's configured location defines how calendar boundaries should be interpreted.

When appropriate, convert values explicitly:

```go
local := value.In(c.location)
```

Then perform calendar calculations using that local representation.

## Do not assume UTC is always the calendar timezone

UTC is an excellent representation for instants, but calendar operations are often defined by a local timezone.

For example:

```text
2026-03-10 00:30 UTC
```

may belong to the previous local calendar day in another timezone.

Calendar APIs must account for this.

---

# Daylight-Saving Time

Never assume that a local calendar day is exactly 24 hours.

For example, a DST transition can produce a local day that is:

```text
23 hours
```

or:

```text
25 hours
```

Therefore, avoid implementations such as:

```go
start.Add(24 * time.Hour)
```

when the intention is:

> the beginning of the next local calendar day

Prefer calendar arithmetic:

```go
start.AddDate(0, 0, 1)
```

Similarly, week and month boundaries should use calendar arithmetic rather than fixed-duration arithmetic.

## Canonical intervals

When representing a calendar range, prefer half-open intervals:

```text
[start, end)
```

For example:

```go
[StartOfDay(value), StartOfNextDay(value))
```

This avoids needing to define an artificial "last nanosecond" for most interval operations.

If an API exposes an inclusive `EndOf...` function, its semantics must be documented clearly.

---

# Calendar Arithmetic

Calendar arithmetic and duration arithmetic are different concepts.

These are not equivalent:

```go
t.Add(24 * time.Hour)
```

and:

```go
t.AddDate(0, 0, 1)
```

The first adds an elapsed duration.

The second moves to the corresponding calendar date.

When implementing calendar functionality, determine which semantic is intended before choosing the operation.

## Follow standard library semantics

Package-level helpers such as:

```go
AddDays
AddWeeks
AddMonths
AddQuarters
AddYears
```

should generally follow the semantics of:

```go
time.Time.AddDate
```

unless the API explicitly defines different behavior.

Do not silently introduce "smart" date normalization without documenting and justifying it.

For example, changing the behavior of `AddMonths` to clamp month-end dates differently from `time.Time.AddDate` would be a semantic decision that needs careful consideration.

---

# Fiscal Calendars

Fiscal-year and fiscal-quarter behavior must be treated separately from ordinary Gregorian calendar boundaries.

For example, if:

```go
FiscalYearStart: time.October
```

then:

```text
October → Q1
January → Q2
April   → Q3
July    → Q4
```

When implementing fiscal operations, reuse existing fiscal-year logic where possible.

For example, quarter calculations should derive their year from the calendar's fiscal-year calculation rather than independently reconstructing the year.

This prevents subtle inconsistencies between:

```go
StartOfYear()
```

and:

```go
StartOfQuarter()
```

---

# Testing Requirements

Time-related code requires more than ordinary happy-path tests.

Every change involving dates or times should be evaluated for:

* timezone differences
* DST transitions
* leap years
* month boundaries
* year boundaries
* fiscal-year boundaries
* week boundaries
* midnight boundaries
* values represented in different locations

Tests should verify **semantics**, not implementation details.

---

# Required Boundary Testing

When adding or modifying calendar operations, consider tests around:

### Day

* beginning of day
* end of day
* previous day
* next day
* midnight
* timezone conversion
* DST transition

### Week

* configured week start
* week boundary
* previous week
* next week
* DST transition inside a week
* different week-start configurations

### Month

* first day
* last day
* February
* leap year February
* previous month
* next month
* year boundary

### Quarter

* quarter boundaries
* fiscal-year boundary
* fiscal year beginning in different months
* quarter crossing a calendar year
* previous/next quarter

### Year

* calendar year boundary
* fiscal year boundary
* leap year
* non-leap year
* previous/next year

---

# Timezone Testing

Tests should use multiple locations where timezone behavior matters.

At minimum, consider:

```text
UTC
America/New_York
Europe/London
Asia/Kolkata
```

Additional locations may be appropriate when testing specific timezone rules.

Do not assume that a test passing in UTC proves that calendar behavior is correct.

---

# DST Testing

DST tests should explicitly target actual transition dates.

For example, in `America/New_York`, tests should cover both:

* the spring transition
* the autumn transition

A useful assertion for a DST-sensitive week is that the elapsed duration of the week may differ from 168 hours.

For example:

```text
spring-transition week → 167 hours
autumn-transition week → 169 hours
```

The exact dates should be documented in the test so future readers understand why the assertion exists.

---

# Leap-Year Testing

Leap-year behavior must include:

```text
divisible by 4
divisible by 100
divisible by 400
```

For example:

```text
2024 → leap year
2100 → not a leap year
2000 → leap year
```

This is important because a simplistic:

```go
year%4 == 0
```

implementation is incorrect.

---

# Table-Driven Tests

Prefer table-driven tests for related cases.

Example:

```go
func TestIsLeapYear(t *testing.T) {
    tests := []struct {
        year int
        want bool
    }{
        {2023, false},
        {2024, true},
        {2100, false},
        {2000, true},
    }

    for _, tt := range tests {
        t.Run(strconv.Itoa(tt.year), func(t *testing.T) {
            if got := IsLeapYear(tt.year); got != tt.want {
                t.Fatalf("IsLeapYear(%d) = %v, want %v", tt.year, got, tt.want)
            }
        })
    }
}
```

Tests should clearly communicate what behavior is being protected.

---

# Test Names

Test names should describe behavior.

Prefer:

```text
TestCalendar_StartOfDay_UsesConfiguredLocation
TestCalendar_StartOfWeek_RespectsWeekStart
TestCalendar_StartOfQuarter_HandlesFiscalYearBoundary
```

over:

```text
TestStartOfDay1
TestQuarter
TestCalendarStuff
```

When a test exists because of a subtle bug, make that reason obvious.

---

# Documentation Requirements

Every exported Go symbol must have a GoDoc comment.

The comment should explain the behavior and semantics, not merely repeat the function name.

Weak:

```go
// StartOfDay returns the start of the day.
```

Better:

```go
// StartOfDay returns the first instant of the calendar day containing value.
//
// The day is determined using the Calendar's configured location rather than
// the location carried by value.
```

For time-related APIs, documentation should explain important semantic decisions such as:

* timezone used
* inclusive/exclusive boundaries
* DST behavior
* fiscal semantics
* calendar vs duration arithmetic
* configuration dependencies

If behavior is surprising, document it explicitly.

---

# README and Documentation

Changes to public behavior should generally include documentation updates.

If you add:

* a package
* an exported API
* a major concept
* configuration
* user-visible behavior

consider whether the README needs to explain it.

Documentation should help users understand **when to use an API**, not only what the API technically does.

---

# Commit Messages

The project uses Conventional Commit-style messages.

Use:

```text
type(scope): description
```

Examples:

```text
feat(calendar): add fiscal quarter calculations
fix(calendar): handle fiscal quarter year boundary
docs(calendar): clarify timezone semantics
test(calendar): add DST transition coverage
refactor(calendar): simplify month boundary calculations
ci: add Go version matrix
chore: update development tooling
```

## Commit Types

### `feat`

A new user-facing capability.

```text
feat(calendar): add StartOfPreviousMonth
```

### `fix`

A bug fix.

```text
fix(calendar): correct fiscal quarter year calculation
```

### `docs`

Documentation-only changes.

```text
docs: improve contributing guide
```

### `test`

Tests without a production-code change.

```text
test(calendar): add leap-year coverage
```

### `refactor`

Internal restructuring without changing intended behavior.

```text
refactor(calendar): centralize fiscal year calculation
```

### `perf`

Performance improvements.

```text
perf(calendar): avoid repeated timezone conversion
```

### `ci`

Continuous integration changes.

```text
ci: add race detector to workflow
```

### `chore`

Maintenance work that does not affect the library API.

```text
chore: update repository metadata
```

Keep commit subjects:

* short
* imperative
* specific
* free of unnecessary punctuation

---

# Pull Requests

A pull request should explain:

1. What changed?
2. Why was it needed?
3. What package does it affect?
4. What behavior changed?
5. How was it tested?
6. Are there timezone/DST implications?
7. Does it introduce or modify a public API?

For example:

```markdown
## Summary

Adds fiscal quarter calculations to the calendar package.

## Motivation

Users with non-January fiscal years need quarter boundaries
based on their configured fiscal year.

## Behavior

Quarter numbering is relative to `FiscalYearStart`.

## Testing

- Added all fiscal-year start months.
- Added fiscal year boundary tests.
- Added year-wrap tests.
- Added timezone tests.

## Compatibility

No existing public APIs were changed.
```

---

# Pull Request Checklist

Before requesting review, verify:

* [ ] The change solves a clearly defined problem.
* [ ] The functionality belongs in the chosen package.
* [ ] Public APIs are minimal and intentional.
* [ ] All exported symbols have GoDoc comments.
* [ ] Tests cover normal behavior.
* [ ] Tests cover relevant boundary conditions.
* [ ] Timezone behavior is explicitly tested where relevant.
* [ ] DST behavior is tested where relevant.
* [ ] Leap-year behavior is tested where relevant.
* [ ] Existing tests still pass.
* [ ] Code has been formatted with `gofmt`.
* [ ] `go vet ./...` passes.
* [ ] `go test ./...` passes.
* [ ] `go test -race ./...` passes.
* [ ] Documentation has been updated where necessary.
* [ ] Commit messages follow the project's conventions.

---

# Backward Compatibility

Once an API is released as part of a stable version, treat it as a long-term contract.

Avoid breaking changes unless there is a compelling reason.

In particular, be cautious about changing:

* function signatures
* exported types
* exported fields
* configuration semantics
* timezone semantics
* boundary definitions
* return values
* error behavior

A seemingly small change such as changing whether an end boundary is inclusive or exclusive can materially affect users.

When behavior needs to change, consider whether a new API can provide the desired semantics without breaking the existing one.

---

# What Usually Does Not Belong in This Project

`business-time` should remain focused.

The following generally do not belong in the core library unless there is a strong, well-defined reason:

* database-specific date handling
* ORM integrations
* HTTP framework integrations
* UI formatting
* locale-specific presentation formatting
* arbitrary string parsing
* application-specific business rules
* domain-specific billing logic
* database models
* custom datetime types that duplicate `time.Time`

Integrations can potentially live in separate projects or packages.

The core library should provide reusable time semantics rather than application-specific behavior.

---

# Avoid Premature Abstraction

Do not introduce abstractions merely because they might be useful someday.

For example, avoid creating an interface when there is only one implementation and no clear need for substitution.

Prefer:

```go
type Calendar struct {
    ...
}
```

over introducing multiple interfaces and implementations without a demonstrated use case.

Similarly, do not add a new package simply because it appears on the project's future roadmap.

The current API should solve today's well-defined problem cleanly.

---

# Performance

Correctness comes before micro-optimization.

Calendar operations should be simple and efficient, but avoid sacrificing readability or semantic correctness for speculative performance improvements.

If a performance-sensitive change is proposed:

1. identify the actual bottleneck
2. benchmark the existing behavior
3. implement the change
4. benchmark the new behavior
5. verify that semantics remain unchanged

Use Go benchmarks when appropriate.

---

# Error Handling

Errors should be returned when invalid configuration or input cannot be represented safely.

Do not silently accept invalid configuration if doing so could produce surprising behavior.

For example, invalid calendar configuration should be rejected during construction rather than causing unpredictable behavior later.

Error messages should be:

* clear
* actionable
* stable enough for humans to understand
* free from unnecessary implementation details

---

# Security

Although `business-time` is primarily a time-handling library, contributions should still consider security implications.

Avoid:

* unsafe parsing
* unnecessary reflection
* executing external commands
* accepting unbounded input without reason
* introducing unnecessary dependencies

Dependencies should have a clear justification.

Prefer the Go standard library when it provides the required functionality.

---

# Dependency Policy

New dependencies should be justified.

Before adding a dependency, consider:

* Can the standard library solve the problem?
* Does the dependency provide substantial value?
* Is it actively maintained?
* Does it introduce transitive dependencies?
* Does it affect the library's users?
* Is the functionality appropriate for the core library?

A small foundational library should avoid dependency bloat.

---

# Release Process

Releases are versioned according to semantic versioning.

The project aims to automate changelog and release management using Conventional Commit-style messages.

For contributors, the important rule is:

> Write commit messages that accurately describe the user-visible impact of your change.

Examples:

```text
feat(calendar): add fiscal quarter calculations
```

may contribute to a feature release, while:

```text
fix(calendar): correct DST boundary calculation
```

represents a bug fix.

The exact release process is maintained by the project maintainers and may evolve as the repository matures.

---

# Before Opening a Pull Request

Run the complete local validation:

```bash
gofmt -w .
go mod tidy
go test ./...
go test -race ./...
go vet ./...
```

Then inspect the repository:

```bash
git status
git diff
```

Make sure you are not accidentally committing:

* generated files
* editor files
* local configuration
* credentials
* build artifacts
* unrelated changes

For public API changes, also inspect the generated documentation:

```bash
go doc ./calendar
go doc ./clock
```

Ask yourself:

> Would I understand this API if I encountered it for the first time?

---

# Questions and Discussions

If you are unsure about an implementation, package boundary, API design, or semantic decision, start a discussion before making a large change.

In particular, please discuss changes involving:

* new packages
* new public APIs
* breaking changes
* timezone semantics
* DST behavior
* fiscal calendar behavior
* business-day definitions
* schedule semantics
* duration semantics

For a time library, getting the semantics right is more important than getting the implementation written quickly.

---

# Final Principle

`business-time` exists to make time-related business logic easier to reason about.

Every contribution should move the project toward that goal.

Prefer APIs that are:

* explicit
* predictable
* timezone-aware
* composable
* idiomatic Go
* easy to test
* difficult to misuse
* small enough to understand
* stable enough to depend on

When in doubt, choose **clear semantics over clever abstractions**.
