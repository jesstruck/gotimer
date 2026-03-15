# Implementation Plan: Time Tracker MVP (Finalized Scope)

## Final MVP Scope
- Support both manual entry and start/stop timer entry.
- Support editing and deleting existing entries.
- Frontend is a week-planner layout showing only Monday-Sunday for a selected ISO week.
- Week switching is done with one-click Previous Week / Next Week / Current Week buttons.
- For each day, user can enter start/end/lunch manually or use start/stop buttons.
- For each day with start and end entered, show the total rounded duration.
- In the bottom-left of the planner, show:
  - all ISO weeks in current month with per-week totals (e.g. `10: 37h`)
  - a month total line (e.g. `Month total: 163h`)
- No additional views (monthly cards, generic entry list pages, or filter panels).
- Apply rounding rule:
  - Rounding is applied to final worked duration only.
  - Duration rounds to half-hour steps.
  - Less than 5 minutes past an hour/half-hour boundary rounds down.
  - 5 minutes or more past an hour/half-hour boundary rounds up to the next half-hour.
- In timer flow, default lunch duration is 30 minutes.
- Week boundaries follow ISO weeks (Monday-Sunday).
- Automated tests are required for MVP completion.

## Out of Scope for MVP
- Project/tag support.
- Additional validation rules beyond basic technical correctness.
- Timezone logic.
- Authentication/authorization.

## Phase 1: Product and UX Lock
- Finalize UX for both entry modes:
  - manual start/end/lunch input
  - timer start/stop flow
- Finalize weekly navigation UX:
  - one-click previous week
  - one-click next week
- Finalize summary UX:
  - weekly totals
  - monthly totals

Exit criteria:
- Wireframes/flows approved for manual entry, timer flow, week navigation, and weekly/monthly summary views.

Status:
- Completed.

Phase 1 artifacts:
- `docs/PHASE1.md`
- `docs/decisions/DEC-001-mvp-scope.md`
- `docs/decisions/DEC-002-rounding-rule.md`
- `docs/decisions/DEC-003-summary-navigation.md`

## Phase 2: Core Backend Delivery
- Implement/confirm time-entry persistence model (start, end, lunch, timestamps, IDs).
- Implement backend APIs for:
  - create entry
  - list entries
  - update entry
  - delete entry
  - timer start/stop
  - weekly summary
  - monthly summary
- Implement rounding behavior in backend calculation path so totals are consistent.

Exit criteria:
- API supports full MVP lifecycle and returns weekly/monthly aggregates correctly.

Status:
- Completed.

## Phase 3: Core Frontend Delivery
- Build a week-planner UI that renders only Monday-Sunday.
- For each day row:
  - manual start/end/lunch entry fields
  - start/stop buttons
  - save action
  - delete action
  - total display when start and end exist
- Remove separate summary/list/filter UI surfaces from frontend.

Exit criteria:
- User can complete end-to-end flows from the week planner only:
  - enter day values manually and save
  - use start/stop on a day row
  - update existing day entry through the row and save
  - see per-day totals for days with start/end

Status:
- Completed.

## Phase 4: QA and MVP Release
- Add required automated tests for backend and frontend MVP flows.
- Add/complete test coverage for:
  - rounding behavior
  - weekly and monthly summary correctness
  - edit/delete behavior
  - timer flow behavior
- Run full sanity check against this scope before release.

Exit criteria:
- All required tests pass and MVP behavior matches finalized scope.

Status:
- Completed.

Phase 4 verification:
- Backend tests: `GOCACHE=$(pwd)/.cache/go-build go test ./...` (in `backend/`)
- Frontend tests: `CI=true npm test -- --watch=false` (in `frontend/`)
