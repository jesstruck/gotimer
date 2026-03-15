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

## Usage

- Use the frontend form to input your work hours.
- The application will calculate and display the total days of work for the current week or month based on your entries.

## License

This project is licensed under the MIT License.
