package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"time-tracker-app/backend/database"
)

func setupHandlerDB(t *testing.T) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	if err := database.InitDB(dbPath); err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
}

func TestGetThemePreferenceDefaultsToSeaTurtle(t *testing.T) {
	setupHandlerDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/preferences/theme", nil)
	rec := httptest.NewRecorder()

	GetThemePreference(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var payload themePreferencePayload
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.Theme != defaultThemeID {
		t.Fatalf("expected theme %s, got %s", defaultThemeID, payload.Theme)
	}
}

func TestUpdateThemePreferencePersistsValue(t *testing.T) {
	setupHandlerDB(t)

	body := bytes.NewBufferString(`{"theme":"blue-whale"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/preferences/theme", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	UpdateThemePreference(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/preferences/theme", nil)
	getRec := httptest.NewRecorder()
	GetThemePreference(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", getRec.Code)
	}

	var payload themePreferencePayload
	if err := json.NewDecoder(getRec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.Theme != "blue-whale" {
		t.Fatalf("expected blue-whale, got %s", payload.Theme)
	}
}
