## Summary

<!--
Briefly describe what this pull request changes and why.
-->


## Motivation

<!--
What problem does this solve?
Why is this change needed?
-->


## Changes

<!--
List the important changes made by this PR.
-->

- 
- 
- 

## Package

<!--
Which package(s) does this change affect?
-->

- [ ] `calendar`
- [ ] `clock`
- [ ] `business`
- [ ] `schedule`
- [ ] `duration`
- [ ] `deadline`
- [ ] `period`
- [ ] Repository / CI / documentation
- [ ] Other

## API Changes

<!--
Does this PR introduce, remove, or modify a public API?

If yes, describe the API and why it is designed this way.
-->

- [ ] No public API changes
- [ ] Adds a public API
- [ ] Changes an existing public API
- [ ] Removes a public API

## Time Semantics

<!--
Complete this section when the PR affects dates, times, calendars,
timezones, durations, schedules, or deadlines.
-->

### Timezone

<!--
Which timezone defines the behavior?
-->


### Calendar vs Duration

<!--
Is this calendar arithmetic or elapsed-duration arithmetic?
Explain if relevant.
-->


### DST

<!--
Does daylight-saving time affect this behavior?
If yes, describe how it is handled and which transitions are tested.
-->

- [ ] Not applicable
- [ ] Considered and tested
- [ ] Requires additional discussion

### Boundaries

<!--
Describe boundary semantics where relevant.

Examples:
- [start, end)
- inclusive end
- start of local day
- fiscal-year boundary
-->


## Testing

<!--
Describe what has been tested.
-->

- [ ] Existing tests pass
- [ ] Added or updated unit tests
- [ ] Added boundary tests
- [ ] Added timezone tests
- [ ] Added DST tests
- [ ] Added leap-year tests
- [ ] Added fiscal-calendar tests
- [ ] Added race-detector coverage where relevant

## Documentation

- [ ] Added/updated GoDoc
- [ ] Updated README
- [ ] Updated other documentation
- [ ] No documentation changes required

## Compatibility

<!--
Does this change affect backward compatibility?
-->

- [ ] No breaking changes
- [ ] Potentially breaking change — discussed in an issue
- [ ] Breaking change — requires maintainer review

## Validation

<!--
Confirm that these commands pass locally.
-->

```bash
gofmt -l .
go test ./...
go test -race ./...
go vet ./...