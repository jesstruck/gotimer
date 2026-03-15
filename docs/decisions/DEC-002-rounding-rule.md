# DEC-002: Time Rounding Rule

- Status: Accepted
- Date: 2026-03-14
- Type: Technical/Domain Logic

## Decision
Rounding uses half-hour boundaries and is applied to final worked duration only:
- If time is less than 5 minutes past an hour/half-hour boundary, round down.
- If time is 5 minutes or more past an hour/half-hour boundary, round up to the next half-hour.

Examples:
- `09:04` -> `09:00`
- `09:05` -> `09:30`
- `09:34` -> `09:30`
- `09:35` -> `10:00`

## Impact
This directly affects total hour calculations and test cases for weekly/monthly summaries.
