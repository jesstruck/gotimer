# Frontend Time Tracker Application

This frontend is built with React + TypeScript and uses a single week-planner UI.

## Getting Started

1. Install dependencies:
   ```bash
   npm install
   ```
2. Run the dev server:
   ```bash
   npm start
   ```
3. Open `http://localhost:3000`.

The frontend proxies API calls to `http://localhost:8080` (see `package.json`).

## Current UI (Week Planner Only)

- Exactly seven rows are shown: Monday-Sunday.
- One-click week navigation controls:
  - Previous Week
  - Next Week
  - Current Week
- Per day:
  - start time
  - end time
  - lunch minutes
  - start/stop buttons
  - Stop is shown when a valid start time exists and end time is empty; otherwise Start is shown
  - start/end auto-save when valid start time exists (open entries are supported until end time is set)
  - delete button
  - total display (when start/end exist)
- Bottom-left totals block includes:
  - one line per ISO week in current month (`<weekNumber>: <hours>h`)
  - `Month total: <hours>h`
  - weeks that only contain weekend days in that month are excluded
- No extra summary/list/filter pages are shown in the frontend.

## Key Source Files

- `src/App.tsx`: mounts the planner as the only UI.
- `src/components/WeekPlanner.tsx`: full week planner behavior.
- `src/types/index.ts`: frontend API data models.
- `src/utils/datetime.ts`: shared UTC/date/time helpers.

## Testing

Run frontend tests:

```bash
CI=true npm test -- --watch=false
```

## Docker

Frontend is built as static files and served via Nginx, with `/api/*` proxied to backend service.

Build frontend image:

```bash
docker build -t time-tracker-frontend .
```
