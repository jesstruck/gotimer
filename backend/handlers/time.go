package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"time-tracker-app/backend/database"
	"time-tracker-app/backend/models"

	"github.com/gorilla/mux"
)

const defaultTimerLunchDuration = 30

type apiError struct {
	Error string `json:"error"`
}

type timeEntryPayload struct {
	Date               string `json:"date"`
	StartTimeSnakeCase string `json:"start_time"`
	EndTimeSnakeCase   string `json:"end_time"`
	LunchSnakeCase     *int   `json:"lunch_duration"`
	StartTimeCamelCase string `json:"startTime"`
	EndTimeCamelCase   string `json:"endTime"`
	LunchCamelCase     *int   `json:"lunchDuration"`
}

func (p timeEntryPayload) startRaw() string {
	if p.StartTimeSnakeCase != "" {
		return p.StartTimeSnakeCase
	}
	return p.StartTimeCamelCase
}

func (p timeEntryPayload) endRaw() string {
	if p.EndTimeSnakeCase != "" {
		return p.EndTimeSnakeCase
	}
	return p.EndTimeCamelCase
}

func (p timeEntryPayload) lunchMinutes(defaultValue int) int {
	if p.LunchSnakeCase != nil {
		return *p.LunchSnakeCase
	}
	if p.LunchCamelCase != nil {
		return *p.LunchCamelCase
	}
	return defaultValue
}

func SetCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func HandleOptions(w http.ResponseWriter, _ *http.Request) {
	SetCORSHeaders(w)
	w.WriteHeader(http.StatusNoContent)
}

func CreateTimeEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		HandleOptions(w, r)
		return
	}

	SetCORSHeaders(w)

	var payload timeEntryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	entry, err := payloadToTimeEntry(payload, "manual")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.CreateTimeEntry(&entry); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, entry.ToResponse())
}

func GetTimeEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		HandleOptions(w, r)
		return
	}
	SetCORSHeaders(w)

	entries, err := database.GetTimeEntries()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := make([]models.TimeEntryResponse, 0, len(entries))
	for _, entry := range entries {
		resp = append(resp, entry.ToResponse())
	}
	writeJSON(w, http.StatusOK, resp)
}

func UpdateTimeEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		HandleOptions(w, r)
		return
	}
	SetCORSHeaders(w)

	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var payload timeEntryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	entry, err := payloadToTimeEntry(payload, "manual")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	entry.ID = id

	if err := database.UpdateTimeEntry(&entry); err != nil {
		handleDBError(w, err)
		return
	}

	updated, err := database.GetTimeEntryByID(id)
	if err != nil {
		handleDBError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated.ToResponse())
}

func DeleteTimeEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		HandleOptions(w, r)
		return
	}
	SetCORSHeaders(w)

	id, err := parseIDParam(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DeleteTimeEntry(id); err != nil {
		handleDBError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func StartTimer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		HandleOptions(w, r)
		return
	}
	SetCORSHeaders(w)

	state, err := database.StartTimer(time.Now().UTC())
	if err != nil {
		handleDBError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, state)
}

func StopTimer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		HandleOptions(w, r)
		return
	}
	SetCORSHeaders(w)

	entry, err := database.StopTimer(time.Now().UTC(), defaultTimerLunchDuration)
	if err != nil {
		handleDBError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, entry.ToResponse())
}

func GetActiveTimer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		HandleOptions(w, r)
		return
	}
	SetCORSHeaders(w)

	state, err := database.GetActiveTimer()
	if err != nil {
		handleDBError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func GetWeeklySummary(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		HandleOptions(w, r)
		return
	}
	SetCORSHeaders(w)

	anchor, offset, err := parseAnchorAndOffset(r, "week_offset")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	summary, err := database.GetWeeklySummary(anchor, offset)
	if err != nil {
		handleDBError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func GetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		HandleOptions(w, r)
		return
	}
	SetCORSHeaders(w)

	anchor, offset, err := parseAnchorAndOffset(r, "month_offset")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	summary, err := database.GetMonthlySummary(anchor, offset)
	if err != nil {
		handleDBError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

// GetSummary keeps backward compatibility with the previous endpoint path.
func GetSummary(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "month" || period == "monthly" {
		GetMonthlySummary(w, r)
		return
	}
	GetWeeklySummary(w, r)
}

func payloadToTimeEntry(payload timeEntryPayload, source string) (models.TimeEntry, error) {
	dateRef, err := parseDateOrToday(payload.Date)
	if err != nil {
		return models.TimeEntry{}, err
	}

	startRaw := payload.startRaw()
	endRaw := payload.endRaw()
	if startRaw == "" {
		return models.TimeEntry{}, errors.New("start_time is required")
	}

	startTime, err := parseFlexibleTime(startRaw, dateRef)
	if err != nil {
		return models.TimeEntry{}, fmt.Errorf("invalid start_time: %w", err)
	}

	isOpen := endRaw == ""
	endTime := startTime
	if !isOpen {
		endTime, err = parseFlexibleTime(endRaw, dateRef)
		if err != nil {
			return models.TimeEntry{}, fmt.Errorf("invalid end_time: %w", err)
		}
	}

	return models.TimeEntry{
		StartTime:     startTime.UTC(),
		EndTime:       endTime.UTC(),
		IsOpen:        isOpen,
		LunchDuration: payload.lunchMinutes(0),
		Source:        source,
	}, nil
}

func parseFlexibleTime(raw string, dateRef time.Time) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02T15:04",
	}
	for _, layout := range formats {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed, nil
		}
	}

	timeOnlyFormats := []string{"15:04", "15:04:05"}
	for _, layout := range timeOnlyFormats {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return time.Date(
				dateRef.Year(),
				dateRef.Month(),
				dateRef.Day(),
				parsed.Hour(),
				parsed.Minute(),
				parsed.Second(),
				0,
				time.UTC,
			), nil
		}
	}

	return time.Time{}, errors.New("unsupported time format")
}

func parseDateOrToday(raw string) (time.Time, error) {
	if raw == "" {
		now := time.Now().UTC()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, errors.New("date must use YYYY-MM-DD")
	}
	return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC), nil
}

func parseAnchorAndOffset(r *http.Request, offsetKey string) (time.Time, int, error) {
	anchorRaw := r.URL.Query().Get("anchor_date")
	anchor, err := parseDateOrToday(anchorRaw)
	if err != nil {
		return time.Time{}, 0, err
	}

	offsetRaw := r.URL.Query().Get(offsetKey)
	if offsetRaw == "" {
		return anchor, 0, nil
	}
	offset, err := strconv.Atoi(offsetRaw)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("%s must be an integer", offsetKey)
	}
	return anchor, offset, nil
}

func parseIDParam(r *http.Request) (int, error) {
	raw := mux.Vars(r)["id"]
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, apiError{Error: message})
}

func handleDBError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, database.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, database.ErrTimerAlreadyActive):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, database.ErrNoActiveTimer):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
