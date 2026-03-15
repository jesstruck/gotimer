# DEC-005: Week Planner Only UI

- Status: Accepted
- Date: 2026-03-14
- Type: UX/Product

## Decision
Frontend UI is constrained to a single week-planner view:
- Only Monday-Sunday rows are shown.
- One-click week navigation is included (Previous, Next, Current Week).
- Each day row supports start/end/lunch entry and start/stop + save actions.
- Per-day total is shown when start and end are entered.
- A bottom-left totals block displays:
  - All ISO weeks in current month with per-week totals
  - ISO weeks that only contain weekend days in the month are excluded
  - Month total
- Other UI surfaces are removed from the frontend experience.

## Rationale
This aligns the product with a planner-style workflow and reduces UI complexity.

## Impact
- Frontend structure centers around one planner component.
- Prior summary/list-oriented UI components are removed.
