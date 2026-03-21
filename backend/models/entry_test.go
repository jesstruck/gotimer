package models

import (
	"testing"
	"time"
)

func TestRoundWorkedMinutes(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{name: "negative to zero", input: -10, expected: 0},
		{name: "zero remains zero", input: 0, expected: 0},
		{name: "exact bucket", input: 60, expected: 60},
		{name: "below threshold rounds down", input: 454, expected: 450},
		{name: "at threshold rounds up", input: 455, expected: 480},
		{name: "small value below threshold", input: 4, expected: 0},
		{name: "small value at threshold", input: 5, expected: 30},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := RoundWorkedMinutes(tt.input)
			if got != tt.expected {
				t.Fatalf("RoundWorkedMinutes(%d) = %d, expected %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestWorkedMinutesReturnsZeroForOpenEntry(t *testing.T) {
	entry := TimeEntry{
		StartTime:     time.Date(2026, 3, 17, 9, 0, 0, 0, time.UTC),
		EndTime:       time.Date(2026, 3, 17, 17, 0, 0, 0, time.UTC),
		LunchDuration: 30,
		IsOpen:        true,
	}

	if got := entry.WorkedMinutes(); got != 0 {
		t.Fatalf("expected open entry worked minutes 0, got %d", got)
	}
}
