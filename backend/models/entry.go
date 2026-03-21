package models

import (
	"math"
	"time"
)

type TimeEntry struct {
	ID            int       `json:"id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	IsOpen        bool      `json:"is_open"`
	LunchDuration int       `json:"lunch_duration"`
	Source        string    `json:"source"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TimeEntryResponse struct {
	ID                 int       `json:"id"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	IsOpen             bool      `json:"is_open"`
	LunchDuration      int       `json:"lunch_duration"`
	Source             string    `json:"source"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	WorkedMinutes      int       `json:"worked_minutes"`
	RoundedWorkedMins  int       `json:"rounded_worked_minutes"`
	RoundedWorkedHours float64   `json:"rounded_worked_hours"`
}

type TimerState struct {
	Active    bool      `json:"active"`
	StartedAt time.Time `json:"started_at,omitempty"`
}

type SummaryResponse struct {
	Period             string    `json:"period"`
	AnchorDate         string    `json:"anchor_date"`
	RangeStart         time.Time `json:"range_start"`
	RangeEnd           time.Time `json:"range_end"`
	TotalEntries       int       `json:"total_entries"`
	TotalWorkedMinutes int       `json:"total_worked_minutes"`
	TotalRoundedMins   int       `json:"total_rounded_minutes"`
	TotalRoundedHours  float64   `json:"total_rounded_hours"`
}

func (te TimeEntry) WorkedMinutes() int {
	if te.IsOpen {
		return 0
	}

	diffMinutes := int(math.Round(te.EndTime.Sub(te.StartTime).Minutes()))
	worked := diffMinutes - te.LunchDuration
	if worked < 0 {
		return 0
	}
	return worked
}

// RoundWorkedMinutes applies the agreed MVP rounding rule to final worked
// duration only using 30-minute buckets with a 5-minute threshold.
func RoundWorkedMinutes(minutes int) int {
	if minutes <= 0 {
		return 0
	}

	remainder := minutes % 30
	if remainder == 0 {
		return minutes
	}
	if remainder < 5 {
		return minutes - remainder
	}
	return minutes + (30 - remainder)
}

func (te TimeEntry) ToResponse() TimeEntryResponse {
	worked := te.WorkedMinutes()
	rounded := RoundWorkedMinutes(worked)

	return TimeEntryResponse{
		ID:                 te.ID,
		StartTime:          te.StartTime,
		EndTime:            te.EndTime,
		IsOpen:             te.IsOpen,
		LunchDuration:      te.LunchDuration,
		Source:             te.Source,
		CreatedAt:          te.CreatedAt,
		UpdatedAt:          te.UpdatedAt,
		WorkedMinutes:      worked,
		RoundedWorkedMins:  rounded,
		RoundedWorkedHours: float64(rounded) / 60.0,
	}
}
