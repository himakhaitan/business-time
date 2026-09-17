# Changelog

All notable changes to `business-time` are documented in this file.

The project follows [Semantic Versioning](https://semver.org/).

## [1.0.0] - 2026-09-17

### Added

* Introduced the `calendar` package for timezone-aware calendar operations.
* Added calendar-day operations:

  * `StartOfDay`
  * `EndOfDay`
  * `StartOfNextDay`
  * `StartOfPreviousDay`
  * `IsSameDay`
  * `DayOfYear`
* Added calendar-week operations:

  * `Weekday`
  * `StartOfWeek`
  * `EndOfWeek`
  * `StartOfNextWeek`
  * `StartOfPreviousWeek`
  * `IsSameWeek`
  * `AddWeeks`
* Added calendar-month operations:

  * `StartOfMonth`
  * `EndOfMonth`
  * `StartOfNextMonth`
  * `StartOfPreviousMonth`
  * `IsSameMonth`
  * `DaysInMonth`
  * `AddMonths`
* Added calendar-quarter operations:

  * `StartOfQuarter`
  * `EndOfQuarter`
  * `StartOfNextQuarter`
  * `StartOfPreviousQuarter`
  * `IsSameQuarter`
  * `Quarter`
  * `AddQuarters`
* Added fiscal-year-aware calendar operations:

  * `StartOfYear`
  * `EndOfYear`
  * `StartOfNextYear`
  * `StartOfPreviousYear`
  * `IsSameYear`
  * `AddYears`
* Added leap-year and year-length helpers:

  * `IsLeapYear`
  * `DaysInYear`
* Added configurable calendar settings:

  * timezone/location
  * week start
  * fiscal year start
* Added the `clock` package for injectable time sources.
* Added timezone-aware calendar boundary calculations.
* Added DST-aware calendar arithmetic using calendar operations rather than fixed durations.
* Added comprehensive tests covering timezone boundaries, DST transitions, leap years, fiscal years, and calendar boundaries.

[1.0.0]: https://github.com/himakhaitan/business-time/releases/tag/v1.0.0
