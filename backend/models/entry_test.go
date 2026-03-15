package models

import "testing"

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
