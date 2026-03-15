# Phase 1: Product and UX Lock

This document captures the finalized Phase 1 scope, user flows, and low-fidelity wireframes for MVP.

## 1) Scope Lock

### In Scope (MVP)
- Manual entry (start time, end time, lunch duration).
- Timer entry (start/stop).
- Edit and delete existing entries.
- Weekly summary and monthly summary.
- One-click week navigation (previous week / next week).
- Rounding rule based on half-hour boundaries, applied to final worked duration only.
- Timer flow default lunch duration set to 30 minutes.
- ISO week boundaries (Monday-Sunday).

### Out of Scope (MVP)
- Project/tag support.
- Generic filtering controls.
- Timezone logic.
- Authentication/authorization.
- Extra business validation rules beyond baseline technical correctness.

## 2) User Flows

### Flow A: Manual Entry
1. User opens app.
2. User enters start time, end time, lunch duration.
3. User saves entry.
4. Entry appears in list.
5. Weekly/monthly totals update.

### Flow B: Timer Entry
1. User clicks `Start Timer`.
2. Timer runs.
3. User clicks `Stop Timer`.
4. Lunch duration defaults to 30 minutes.
5. Entry is saved and appears in list.
6. Weekly/monthly totals update.

### Flow C: Edit Entry
1. User selects an existing entry.
2. User edits values.
3. User saves changes.
4. Entry and totals refresh.

### Flow D: Delete Entry
1. User selects an existing entry.
2. User clicks delete.
3. User confirms delete.
4. Entry is removed and totals refresh.

### Flow E: Summary Navigation
1. User views weekly summary by default.
2. User clicks `Previous Week` or `Next Week` once per navigation step.
3. User switches between `Weekly` and `Monthly` summary views.

## 3) Low-Fidelity Wireframes (Text)

### Main Screen
```text
+---------------------------------------------------------------+
| Time Tracker                                                  |
+---------------------------------------------------------------+
| [Manual Entry] [Timer Entry]                                 |
|                                                               |
| Manual Entry Form                                             |
| Start: [ 09:00 ]  End: [ 17:00 ]  Lunch(min): [ 30 ] [Save]  |
|                                                               |
| Timer Entry                                                   |
| [Start Timer] [Stop Timer]   Running: 00:42:10               |
|                                                               |
| Summary                                                       |
| [Weekly] [Monthly]   [< Previous Week] [Next Week >]         |
| Weekly Total: 37.5h                                           |
| Monthly Total: 142.0h                                         |
|                                                               |
| Entries                                                       |
| Date       Start   End     Lunch   Rounded Hours   Actions   |
| 2026-03-16 09:00   17:00   30      7.5            [Edit][X] |
| 2026-03-17 08:30   16:30   30      7.5            [Edit][X] |
+---------------------------------------------------------------+
```

### Edit Entry Modal
```text
+----------------------------------------+
| Edit Entry                             |
+----------------------------------------+
| Date:   [ 2026-03-16 ]                 |
| Start:  [ 09:00 ]                      |
| End:    [ 17:00 ]                      |
| Lunch:  [ 30 ]                         |
|                                        |
| [Cancel]                     [Save]    |
+----------------------------------------+
```

### Delete Confirmation
```text
+----------------------------------------+
| Delete this entry?                     |
| This action cannot be undone.          |
|                                        |
| [Cancel]                   [Delete]    |
+----------------------------------------+
```

## 4) Acceptance Criteria for Phase 1
- Scope is locked and aligned with MVP in `docs/PLAN.md`.
- Core flows are documented (manual, timer, edit, delete, summary navigation).
- Low-fidelity wireframes are documented.
- Key UX and technical decisions are recorded in `docs/decisions`.

## 5) Finalized Decision Summary
- Rounding is applied to final worked duration only.
- Timer flow lunch duration defaults to 30 minutes.
- Week boundaries follow ISO week (Monday-Sunday).
