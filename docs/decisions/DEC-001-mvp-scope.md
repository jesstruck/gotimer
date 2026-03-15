# DEC-001: MVP Scope Lock

- Status: Accepted
- Date: 2026-03-14
- Type: Product/Scope

## Decision
The MVP includes:
- Manual entry and timer entry.
- Edit and delete entries.
- Weekly and monthly summaries.
- One-click week navigation (previous/next week).
- Rounding behavior (see DEC-002).
- Automated tests as a release requirement.

The MVP excludes:
- Project/tag support.
- Generic filtering controls.
- Timezone logic.
- Authentication/authorization.
- Additional business validation rules beyond baseline technical correctness.

## Rationale
This scope keeps the MVP focused on core time tracking behavior while limiting implementation and release risk.

## Impact
- Backend and frontend implementation focus on core entry lifecycle and summary calculations.
- Validation, auth, and timezone work are deferred.
