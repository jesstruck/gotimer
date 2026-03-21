# Time Tracker Backend

This is the backend component of the Time Tracker application, built with Go, Gorilla Mux, and SQLite.

Phase 2 delivers:
- Full time-entry CRUD endpoints
- Timer start/stop endpoints
- Weekly and monthly summary endpoints
- Server-side worked-duration rounding logic

## Project Structure

- **main.go**: Entry point of the application. Initializes the server and sets up routes.
- **go.mod**: Defines the module and dependencies for the Go application.
- **go.sum**: Contains checksums for the dependencies to ensure consistent builds.
- **handlers/time.go**: Contains functions for handling HTTP requests related to time entries.
- **models/entry.go**: Defines the data structure for a time entry and methods for database interaction.
- **database/sqlite.go**: Manages the SQLite database connection and query execution.

## Setup Instructions

1. Ensure you have Go installed on your machine.
2. Go to the backend directory:
   ```
   cd backend
   ```
3. Install dependencies (if needed):
   ```
   go mod tidy
   ```
4. Run the application:
   ```
   go run main.go
   ```

The API server runs on `http://localhost:8080`.

Database path:
- Defaults to `time_entries.db` in current working directory.
- Can be overridden with `TIME_ENTRIES_DB_PATH` (used by Docker Compose).

## Testing

Run backend tests:

```bash
GOCACHE=$(pwd)/.cache/go-build go test ./...
```

## Docker

Build backend image:

```bash
docker build -t time-tracker-backend .
```

## API Routes

All routes are under `/api`.

### Time Entries
- **POST `/api/time-entries`**: Create a new manual entry.
- **GET `/api/time-entries`**: List all entries.
- **PUT `/api/time-entries/{id}`**: Update an entry.
- **DELETE `/api/time-entries/{id}`**: Delete an entry.

Example request body:

```json
{
  "date": "2026-03-14",
  "start_time": "09:00",
  "end_time": "17:00",
  "lunch_duration": 30
}
```

Accepted time formats:
- RFC3339 (for example `2026-03-14T09:00:00Z`)
- `YYYY-MM-DD HH:MM[:SS]`
- `HH:MM[:SS]` (combined with `date`, or current UTC day if `date` is omitted)

### Timer
- **POST `/api/timer/start`**: Start active timer.
- **POST `/api/timer/stop`**: Stop active timer and create entry.
  - Uses default lunch duration of `30` minutes.
- **GET `/api/timer/active`**: Read current timer state.

### Summaries
- **GET `/api/summaries/weekly`**: Weekly summary
  - Query params:
    - `anchor_date=YYYY-MM-DD` (optional, defaults to today UTC)
    - `week_offset=<int>` (optional)
- **GET `/api/summaries/monthly`**: Monthly summary
  - Query params:
    - `anchor_date=YYYY-MM-DD` (optional, defaults to today UTC)
    - `month_offset=<int>` (optional)

Backward-compatible alias:
- **GET `/api/time-entries/summary`**
  - Weekly by default
  - Monthly if `period=month`

### Preferences
- **GET `/api/preferences/theme`**: Read persisted UI theme preference.
- **PUT `/api/preferences/theme`**: Save UI theme preference.

Example request body:

```json
{
  "theme": "blue-whale"
}
```

## Rounding Rule

Rounding is applied to final worked duration only:
- `< 5` minutes past a 30-minute boundary: round down
- `>= 5` minutes past a 30-minute boundary: round up

## License

This project is licensed under the MIT License. See the LICENSE file for details.
