package handlers

import (
	"testing"
	"time"
)

func TestPayloadToTimeEntryStartOnlyCreatesOpenEntry(t *testing.T) {
	entry, err := payloadToTimeEntry(timeEntryPayload{
		Date:               "2026-03-17",
		StartTimeSnakeCase: "08:30",
		LunchSnakeCase:     intPtr(30),
	}, "manual")
	if err != nil {
		t.Fatalf("payloadToTimeEntry returned error: %v", err)
	}

	expected := time.Date(2026, 3, 17, 8, 30, 0, 0, time.UTC)
	if !entry.StartTime.Equal(expected) {
		t.Fatalf("expected start %s, got %s", expected, entry.StartTime)
	}
	if !entry.EndTime.Equal(expected) {
		t.Fatalf("expected end to match start %s, got %s", expected, entry.EndTime)
	}
	if !entry.IsOpen {
		t.Fatalf("expected entry to be open")
	}
	if entry.LunchDuration != 30 {
		t.Fatalf("expected lunch duration 30, got %d", entry.LunchDuration)
	}
}

func TestPayloadToTimeEntryWithEndCreatesClosedEntry(t *testing.T) {
	entry, err := payloadToTimeEntry(timeEntryPayload{
		Date:               "2026-03-17",
		StartTimeSnakeCase: "08:30",
		EndTimeSnakeCase:   "16:45",
		LunchSnakeCase:     intPtr(30),
	}, "manual")
	if err != nil {
		t.Fatalf("payloadToTimeEntry returned error: %v", err)
	}

	if entry.IsOpen {
		t.Fatalf("expected entry to be closed")
	}
	expectedEnd := time.Date(2026, 3, 17, 16, 45, 0, 0, time.UTC)
	if !entry.EndTime.Equal(expectedEnd) {
		t.Fatalf("expected end %s, got %s", expectedEnd, entry.EndTime)
	}
}

func TestPayloadToTimeEntryRequiresStartTime(t *testing.T) {
	_, err := payloadToTimeEntry(timeEntryPayload{
		Date:             "2026-03-17",
		EndTimeSnakeCase: "16:45",
	}, "manual")
	if err == nil {
		t.Fatalf("expected error when start time is missing")
	}
	if err.Error() != "start_time is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func intPtr(v int) *int {
	return &v
}
