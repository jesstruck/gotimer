# Week Planner UI Specification

## Intent
The frontend must present a single week-planner view and nothing else.

## View Rules
- Show exactly seven day rows: Monday through Sunday.
- Each row represents one day in the selected ISO week.
- Include one-click week navigation buttons:
  - Previous Week
  - Next Week
  - Current Week
- A compact totals block is shown in the bottom-left area of the planner.
- No separate summary dashboard, filter panel, or full entry list page.

## Row Interactions
Each day row includes:
- Start time input
- End time input
- Lunch duration input (minutes)
- Start button
- Stop button
- Save button
- Delete button
- Total display

## Totals
- A total is displayed for a day only when start and end are present.
- Total uses the MVP rounding rule:
  - round final worked duration only
  - half-hour buckets
  - `< 5` minutes past bucket boundary rounds down
  - `>= 5` rounds up
- Bottom-left totals block includes:
  - A line per ISO week in the current month (format: `<weekNumber>: <hours>h`)
  - A month total line (format: `Month total: <hours>h`)
  - ISO weeks that contain only weekend days in the month are excluded

## Data Behavior
- Save persists the row to backend (create/update).
- Existing entries for days in the shown week populate corresponding rows.
- Timer buttons are day-row controls and feed the row time fields.

## Out of Scope in This UI
- Monthly summary card/panel.
- Week/month filters.
- Multi-page entry management UI.
