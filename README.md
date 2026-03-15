# Time Tracker Application

This project is a time tracking application that allows users to input their work hours, including start and end times, as well as lunch durations. It calculates the total days of work for the current week or month.

## Project Structure

```
time-tracker-app
├── backend
│   ├── main.go          # Entry point for the backend application
│   ├── go.mod           # Go module definition
│   ├── go.sum           # Dependency checksums
│   ├── handlers          # HTTP request handlers
│   │   └── time.go      # Functions for handling time entries
│   ├── models           # Data models
│   │   └── entry.go     # Time entry data structure
│   ├── database         # Database management
│   │   └── sqlite.go    # SQLite database connection and queries
│   └── README.md        # Backend documentation
├── frontend
│   ├── src              # Frontend source code
│   │   ├── App.tsx      # Main application component
│   │   ├── components    # UI components
│   │   │   ├── TimeEntryForm.tsx # Form for time entry input
│   │   │   └── Summary.tsx        # Summary of total work days
│   │   └── types        # TypeScript types
│   │       └── index.ts # Type definitions
│   ├── package.json      # Frontend npm configuration
│   ├── tsconfig.json     # TypeScript configuration
│   └── README.md         # Frontend documentation
└── README.md             # Overall project documentation
```

## Getting Started

### Prerequisites

- Go (version 1.16 or later)
- Node.js (version 14 or later)
- SQLite

### Backend Setup

1. Navigate to the `backend` directory:
   ```
   cd backend
   ```

2. Install Go dependencies:
   ```
   go mod tidy
   ```

3. Run the backend server:
   ```
   go run main.go
   ```

### Frontend Setup

1. Navigate to the `frontend` directory:
   ```
   cd frontend
   ```

2. Install npm dependencies:
   ```
   npm install
   ```

3. Start the frontend application:
   ```
   npm start
   ```

## Docker Deployment

Run the full stack (frontend + backend + persisted SQLite volume):

```bash
docker compose up --build -d
```

Open:
- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080`

Stop:

```bash
docker compose down
```

Stop and remove database volume:

```bash
docker compose down -v
```

## CI/CD to GHCR (GitHub Actions)

Workflow files:
- `.github/workflows/pr-checks.yml`
- `.github/workflows/main-release.yml`

Behavior:
- Pull requests to `main`: lint commit messages (Conventional Commits) and run backend/frontend tests + validation.
- Push to `main`: lint commit messages in push range, run CI, run semantic-release, then build/push Docker images to GHCR.

Semantic-release:
- Runs only on `main`.
- Creates GitHub release notes and tags (`vX.Y.Z`) from Conventional Commits.
- No release is created when commits do not trigger a semantic version bump.

Commit linting:
- Config file: `commitlint.config.cjs`
- Accepted types: `feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `chore`, `ci`, `build`, `style`, `revert`

Published image names:
- `ghcr.io/<owner>/<repo>-backend`
- `ghcr.io/<owner>/<repo>-frontend`

Common tags:
- `latest` (every successful push to `main`)
- `vX.Y.Z` (when semantic-release publishes a new release)
- `sha-<commit>`

Deploy by pulling from GHCR:

```bash
export GHCR_OWNER=<github-owner-lowercase>
export GHCR_REPO=<github-repository-name-lowercase>
export IMAGE_TAG=latest

docker compose -f docker-compose.ghcr.yml up -d
```

For private packages, login first:

```bash
echo <github_pat_with_read_packages> | docker login ghcr.io -u <github-username> --password-stdin
```

## Usage

- Use the frontend form to input your work hours.
- The application will calculate and display the total days of work for the current week or month based on your entries.

## License

This project is licensed under the MIT License.
