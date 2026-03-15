package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"time-tracker-app/backend/models"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var (
	ErrNotFound           = errors.New("not found")
	ErrTimerAlreadyActive = errors.New("timer already active")
	ErrNoActiveTimer      = errors.New("no active timer")
)

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return err
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS time_entries (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        start_time TEXT NOT NULL,
        end_time TEXT NOT NULL,
        lunch_duration INTEGER NOT NULL DEFAULT 0,
        source TEXT NOT NULL DEFAULT 'manual',
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
    );`

	if _, err := db.Exec(createTableQuery); err != nil {
		return err
	}

	createActiveTimerTable := `CREATE TABLE IF NOT EXISTS active_timer (
		id INTEGER PRIMARY KEY CHECK(id = 1),
		started_at TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`
	if _, err := db.Exec(createActiveTimerTable); err != nil {
		return err
	}

	// Keep compatibility if the DB was created with an earlier schema.
	if err := ensureColumn("time_entries", "source", "TEXT NOT NULL DEFAULT 'manual'"); err != nil {
		return err
	}
	if err := ensureColumn("time_entries", "created_at", "TEXT"); err != nil {
		return err
	}
	if err := ensureColumn("time_entries", "updated_at", "TEXT"); err != nil {
		return err
	}
	if _, err := db.Exec(`UPDATE time_entries SET source = 'manual' WHERE source IS NULL OR source = ''`); err != nil {
		return err
	}
	if _, err := db.Exec(`UPDATE time_entries SET created_at = COALESCE(NULLIF(created_at, ''), datetime('now')) WHERE created_at IS NULL OR created_at = ''`); err != nil {
		return err
	}
	if _, err := db.Exec(`UPDATE time_entries SET updated_at = COALESCE(NULLIF(updated_at, ''), created_at, datetime('now')) WHERE updated_at IS NULL OR updated_at = ''`); err != nil {
		return err
	}

	return nil
}

func Connect() (*sql.DB, error) {
	dsn := strings.TrimSpace(os.Getenv("TIME_ENTRIES_DB_PATH"))
	if dsn == "" {
		dsn = "time_entries.db"
	}
	if err := InitDB(dsn); err != nil {
		return nil, err
	}
	return db, nil
}

func GetDB() *sql.DB {
	return db
}

func CreateTimeEntry(entry *models.TimeEntry) error {
	if db == nil {
		return errors.New("database not initialized")
	}
	if entry == nil {
		return errors.New("nil entry")
	}
	if entry.Source == "" {
		entry.Source = "manual"
	}
	now := time.Now().UTC()
	entry.CreatedAt = now
	entry.UpdatedAt = now

	res, err := db.Exec(
		`INSERT INTO time_entries(start_time, end_time, lunch_duration, source, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)`,
		formatTime(entry.StartTime),
		formatTime(entry.EndTime),
		entry.LunchDuration,
		entry.Source,
		formatTime(entry.CreatedAt),
		formatTime(entry.UpdatedAt),
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		entry.ID = int(id)
	}
	return nil
}

func GetTimeEntries() ([]models.TimeEntry, error) {
	if db == nil {
		return []models.TimeEntry{}, errors.New("database not initialized")
	}

	rows, err := db.Query(`SELECT id, start_time, end_time, lunch_duration, source, created_at, updated_at FROM time_entries ORDER BY start_time DESC`)
	if err != nil {
		return []models.TimeEntry{}, err
	}
	defer rows.Close()

	entries := []models.TimeEntry{}
	for rows.Next() {
		entry, err := scanTimeEntry(rows)
		if err != nil {
			return []models.TimeEntry{}, err
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return []models.TimeEntry{}, err
	}

	return entries, nil
}

func GetTimeEntryByID(id int) (models.TimeEntry, error) {
	if db == nil {
		return models.TimeEntry{}, errors.New("database not initialized")
	}

	row := db.QueryRow(`SELECT id, start_time, end_time, lunch_duration, source, created_at, updated_at FROM time_entries WHERE id = ?`, id)

	var (
		entry                                    models.TimeEntry
		startRaw, endRaw, createdRaw, updatedRaw string
	)
	if err := row.Scan(
		&entry.ID,
		&startRaw,
		&endRaw,
		&entry.LunchDuration,
		&entry.Source,
		&createdRaw,
		&updatedRaw,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.TimeEntry{}, ErrNotFound
		}
		return models.TimeEntry{}, err
	}

	start, err := parseTime(startRaw)
	if err != nil {
		return models.TimeEntry{}, err
	}
	end, err := parseTime(endRaw)
	if err != nil {
		return models.TimeEntry{}, err
	}
	created, err := parseTime(createdRaw)
	if err != nil {
		return models.TimeEntry{}, err
	}
	updated, err := parseTime(updatedRaw)
	if err != nil {
		return models.TimeEntry{}, err
	}

	entry.StartTime = start
	entry.EndTime = end
	entry.CreatedAt = created
	entry.UpdatedAt = updated

	return entry, nil
}

func UpdateTimeEntry(entry *models.TimeEntry) error {
	if db == nil {
		return errors.New("database not initialized")
	}
	if entry == nil {
		return errors.New("nil entry")
	}

	existing, err := GetTimeEntryByID(entry.ID)
	if err != nil {
		return err
	}

	if entry.Source == "" {
		entry.Source = existing.Source
	}
	entry.CreatedAt = existing.CreatedAt
	entry.UpdatedAt = time.Now().UTC()

	res, err := db.Exec(
		`UPDATE time_entries SET start_time = ?, end_time = ?, lunch_duration = ?, source = ?, updated_at = ? WHERE id = ?`,
		formatTime(entry.StartTime),
		formatTime(entry.EndTime),
		entry.LunchDuration,
		entry.Source,
		formatTime(entry.UpdatedAt),
		entry.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func DeleteTimeEntry(id int) error {
	if db == nil {
		return errors.New("database not initialized")
	}
	res, err := db.Exec(`DELETE FROM time_entries WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func StartTimer(startedAt time.Time) (models.TimerState, error) {
	if db == nil {
		return models.TimerState{}, errors.New("database not initialized")
	}

	active, err := GetActiveTimer()
	if err != nil {
		return models.TimerState{}, err
	}
	if active != nil && active.Active {
		return models.TimerState{}, ErrTimerAlreadyActive
	}

	_, err = db.Exec(
		`INSERT OR REPLACE INTO active_timer(id, started_at, created_at) VALUES(1, ?, ?)`,
		formatTime(startedAt.UTC()),
		formatTime(time.Now().UTC()),
	)
	if err != nil {
		return models.TimerState{}, err
	}

	return models.TimerState{Active: true, StartedAt: startedAt.UTC()}, nil
}

func GetActiveTimer() (*models.TimerState, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	var startedRaw string
	err := db.QueryRow(`SELECT started_at FROM active_timer WHERE id = 1`).Scan(&startedRaw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.TimerState{Active: false}, nil
		}
		return nil, err
	}

	startedAt, err := parseTime(startedRaw)
	if err != nil {
		return nil, err
	}

	return &models.TimerState{
		Active:    true,
		StartedAt: startedAt,
	}, nil
}

func StopTimer(stoppedAt time.Time, defaultLunchDuration int) (models.TimeEntry, error) {
	if db == nil {
		return models.TimeEntry{}, errors.New("database not initialized")
	}

	tx, err := db.Begin()
	if err != nil {
		return models.TimeEntry{}, err
	}
	defer tx.Rollback()

	var startedRaw string
	err = tx.QueryRow(`SELECT started_at FROM active_timer WHERE id = 1`).Scan(&startedRaw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.TimeEntry{}, ErrNoActiveTimer
		}
		return models.TimeEntry{}, err
	}

	startedAt, err := parseTime(startedRaw)
	if err != nil {
		return models.TimeEntry{}, err
	}

	if stoppedAt.Before(startedAt) {
		return models.TimeEntry{}, errors.New("stop time is before start time")
	}

	now := time.Now().UTC()
	res, err := tx.Exec(
		`INSERT INTO time_entries(start_time, end_time, lunch_duration, source, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)`,
		formatTime(startedAt),
		formatTime(stoppedAt.UTC()),
		defaultLunchDuration,
		"timer",
		formatTime(now),
		formatTime(now),
	)
	if err != nil {
		return models.TimeEntry{}, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return models.TimeEntry{}, err
	}

	if _, err := tx.Exec(`DELETE FROM active_timer WHERE id = 1`); err != nil {
		return models.TimeEntry{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.TimeEntry{}, err
	}

	return GetTimeEntryByID(int(lastID))
}

func GetWeeklySummary(anchorDate time.Time, weekOffset int) (models.SummaryResponse, error) {
	anchor := anchorDate.UTC().AddDate(0, 0, 7*weekOffset)
	start, end := isoWeekRange(anchor)
	return buildSummary("weekly", anchor, start, end)
}

func GetMonthlySummary(anchorDate time.Time, monthOffset int) (models.SummaryResponse, error) {
	anchor := anchorDate.UTC().AddDate(0, monthOffset, 0)
	start := time.Date(anchor.Year(), anchor.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return buildSummary("monthly", anchor, start, end)
}

func buildSummary(period string, anchorDate time.Time, rangeStart, rangeEnd time.Time) (models.SummaryResponse, error) {
	entries, err := GetTimeEntries()
	if err != nil {
		return models.SummaryResponse{}, err
	}

	totalWorked := 0
	totalRounded := 0
	totalEntries := 0

	for _, entry := range entries {
		if entry.StartTime.Before(rangeStart) || entry.StartTime.After(rangeEnd) {
			continue
		}
		totalEntries++
		worked := entry.WorkedMinutes()
		rounded := models.RoundWorkedMinutes(worked)
		totalWorked += worked
		totalRounded += rounded
	}

	return models.SummaryResponse{
		Period:             period,
		AnchorDate:         anchorDate.Format("2006-01-02"),
		RangeStart:         rangeStart,
		RangeEnd:           rangeEnd,
		TotalEntries:       totalEntries,
		TotalWorkedMinutes: totalWorked,
		TotalRoundedMins:   totalRounded,
		TotalRoundedHours:  float64(totalRounded) / 60.0,
	}, nil
}

func ensureColumn(tableName, columnName, columnDefinition string) error {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid        int
			name       string
			colType    string
			notNull    int
			defaultVal sql.NullString
			pk         int
		)
		if err := rows.Scan(&cid, &name, &colType, &notNull, &defaultVal, &pk); err != nil {
			return err
		}
		if strings.EqualFold(name, columnName) {
			return nil
		}
	}

	alter := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, columnName, columnDefinition)
	_, err = db.Exec(alter)
	return err
}

func scanTimeEntry(scanner interface {
	Scan(dest ...any) error
}) (models.TimeEntry, error) {
	var (
		entry                                    models.TimeEntry
		startRaw, endRaw, createdRaw, updatedRaw string
	)

	if err := scanner.Scan(
		&entry.ID,
		&startRaw,
		&endRaw,
		&entry.LunchDuration,
		&entry.Source,
		&createdRaw,
		&updatedRaw,
	); err != nil {
		return models.TimeEntry{}, err
	}

	start, err := parseTime(startRaw)
	if err != nil {
		return models.TimeEntry{}, err
	}
	end, err := parseTime(endRaw)
	if err != nil {
		return models.TimeEntry{}, err
	}
	created, err := parseTime(createdRaw)
	if err != nil {
		return models.TimeEntry{}, err
	}
	updated, err := parseTime(updatedRaw)
	if err != nil {
		return models.TimeEntry{}, err
	}

	entry.StartTime = start
	entry.EndTime = end
	entry.CreatedAt = created
	entry.UpdatedAt = updated
	if entry.Source == "" {
		entry.Source = "manual"
	}
	return entry, nil
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

func parseTime(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, errors.New("empty time value")
	}

	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}
	for _, format := range formats {
		if parsed, err := time.Parse(format, raw); err == nil {
			return parsed.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %s", raw)
}

func isoWeekRange(anchor time.Time) (time.Time, time.Time) {
	weekday := int(anchor.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start := time.Date(anchor.Year(), anchor.Month(), anchor.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -(weekday - 1))
	end := start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return start, end
}
