package database

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"time-tracker-app/backend/models"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	if err := InitDB(dbPath); err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
}

func TestInitDBCreatesDatabaseFileIfMissing(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "missing", "nested", "time_entries.db")

	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		t.Fatalf("expected DB file to be missing before init, got err=%v", err)
	}

	if err := InitDB(dbPath); err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	info, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("expected DB file to exist after init: %v", err)
	}
	if info.IsDir() {
		t.Fatalf("expected DB path to be a file, got directory")
	}
}

func TestCRUDFlow(t *testing.T) {
	setupTestDB(t)

	entry := &models.TimeEntry{
		StartTime:     time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC),
		EndTime:       time.Date(2026, 3, 16, 17, 0, 0, 0, time.UTC),
		LunchDuration: 30,
		Source:        "manual",
	}

	if err := CreateTimeEntry(entry); err != nil {
		t.Fatalf("CreateTimeEntry failed: %v", err)
	}
	if entry.ID <= 0 {
		t.Fatalf("expected created ID, got %d", entry.ID)
	}

	entries, err := GetTimeEntries()
	if err != nil {
		t.Fatalf("GetTimeEntries failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	entry.EndTime = time.Date(2026, 3, 16, 18, 0, 0, 0, time.UTC)
	entry.LunchDuration = 45
	if err := UpdateTimeEntry(entry); err != nil {
		t.Fatalf("UpdateTimeEntry failed: %v", err)
	}

	updated, err := GetTimeEntryByID(entry.ID)
	if err != nil {
		t.Fatalf("GetTimeEntryByID failed: %v", err)
	}
	if updated.LunchDuration != 45 {
		t.Fatalf("expected lunch 45, got %d", updated.LunchDuration)
	}
	if !updated.EndTime.Equal(entry.EndTime) {
		t.Fatalf("expected end time %s, got %s", entry.EndTime, updated.EndTime)
	}

	if err := DeleteTimeEntry(entry.ID); err != nil {
		t.Fatalf("DeleteTimeEntry failed: %v", err)
	}
	if _, err := GetTimeEntryByID(entry.ID); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestTimerFlow(t *testing.T) {
	setupTestDB(t)

	start := time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC)
	stop := time.Date(2026, 3, 17, 17, 0, 0, 0, time.UTC)

	state, err := StartTimer(start)
	if err != nil {
		t.Fatalf("StartTimer failed: %v", err)
	}
	if !state.Active {
		t.Fatalf("expected active timer")
	}

	if _, err := StartTimer(start.Add(10 * time.Minute)); err != ErrTimerAlreadyActive {
		t.Fatalf("expected ErrTimerAlreadyActive, got %v", err)
	}

	entry, err := StopTimer(stop, 30)
	if err != nil {
		t.Fatalf("StopTimer failed: %v", err)
	}
	if entry.Source != "timer" {
		t.Fatalf("expected source timer, got %s", entry.Source)
	}
	if entry.LunchDuration != 30 {
		t.Fatalf("expected lunch 30, got %d", entry.LunchDuration)
	}

	active, err := GetActiveTimer()
	if err != nil {
		t.Fatalf("GetActiveTimer failed: %v", err)
	}
	if active == nil || active.Active {
		t.Fatalf("expected inactive timer after stop")
	}
}

func TestWeeklyAndMonthlySummaryWithRounding(t *testing.T) {
	setupTestDB(t)

	// Monday: 09:00-17:04 with 30m lunch => 454 worked => rounded 450
	e1 := &models.TimeEntry{
		StartTime:     time.Date(2026, 3, 16, 9, 0, 0, 0, time.UTC),
		EndTime:       time.Date(2026, 3, 16, 17, 4, 0, 0, time.UTC),
		LunchDuration: 30,
		Source:        "manual",
	}
	// Tuesday: 09:00-17:05 with 30m lunch => 455 worked => rounded 480
	e2 := &models.TimeEntry{
		StartTime:     time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
		EndTime:       time.Date(2026, 3, 17, 17, 5, 0, 0, time.UTC),
		LunchDuration: 30,
		Source:        "manual",
	}
	// April entry should not be included for March summary.
	e3 := &models.TimeEntry{
		StartTime:     time.Date(2026, 4, 1, 9, 0, 0, 0, time.UTC),
		EndTime:       time.Date(2026, 4, 1, 17, 0, 0, 0, time.UTC),
		LunchDuration: 30,
		Source:        "manual",
	}

	for _, entry := range []*models.TimeEntry{e1, e2, e3} {
		if err := CreateTimeEntry(entry); err != nil {
			t.Fatalf("CreateTimeEntry failed: %v", err)
		}
	}

	anchor := time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC)
	weekly, err := GetWeeklySummary(anchor, 0)
	if err != nil {
		t.Fatalf("GetWeeklySummary failed: %v", err)
	}
	if weekly.TotalEntries != 2 {
		t.Fatalf("expected weekly total entries=2, got %d", weekly.TotalEntries)
	}
	if weekly.TotalRoundedMins != 930 {
		t.Fatalf("expected weekly rounded minutes=930, got %d", weekly.TotalRoundedMins)
	}

	monthly, err := GetMonthlySummary(anchor, 0)
	if err != nil {
		t.Fatalf("GetMonthlySummary failed: %v", err)
	}
	if monthly.TotalEntries != 2 {
		t.Fatalf("expected monthly total entries=2, got %d", monthly.TotalEntries)
	}
	if monthly.TotalRoundedMins != 930 {
		t.Fatalf("expected monthly rounded minutes=930, got %d", monthly.TotalRoundedMins)
	}
}

func TestPreferenceRoundTrip(t *testing.T) {
	setupTestDB(t)

	if _, err := GetPreference("ui_theme"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing preference, got %v", err)
	}

	if err := SetPreference("ui_theme", "blue-whale"); err != nil {
		t.Fatalf("SetPreference failed: %v", err)
	}

	value, err := GetPreference("ui_theme")
	if err != nil {
		t.Fatalf("GetPreference failed: %v", err)
	}
	if value != "blue-whale" {
		t.Fatalf("expected blue-whale, got %s", value)
	}

	if err := SetPreference("ui_theme", "sea-turtle"); err != nil {
		t.Fatalf("SetPreference update failed: %v", err)
	}

	updated, err := GetPreference("ui_theme")
	if err != nil {
		t.Fatalf("GetPreference after update failed: %v", err)
	}
	if updated != "sea-turtle" {
		t.Fatalf("expected sea-turtle, got %s", updated)
	}
}
