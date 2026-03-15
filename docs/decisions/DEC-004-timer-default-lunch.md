# DEC-004: Timer Flow Default Lunch Duration

- Status: Accepted
- Date: 2026-03-14
- Type: UX/Product

## Decision
When an entry is created through timer start/stop flow, lunch duration defaults to `30` minutes.

## Rationale
This keeps timer flow quick while reflecting a common default workday pattern.

## Impact
- Timer-created entries use `30` lunch minutes unless later edited by the user.
- Tests should verify this default behavior for timer-created entries.
