package main

import (
	"log"
	"net/http"
	"time"

	"time-tracker-app/backend/database"
	"time-tracker-app/backend/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/time-entries", handlers.CreateTimeEntry).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/time-entries", handlers.GetTimeEntries).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/time-entries/{id:[0-9]+}", handlers.UpdateTimeEntry).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/time-entries/{id:[0-9]+}", handlers.DeleteTimeEntry).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/api/timer/start", handlers.StartTimer).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/timer/stop", handlers.StopTimer).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/timer/active", handlers.GetActiveTimer).Methods("GET", "OPTIONS")

	r.HandleFunc("/api/summaries/weekly", handlers.GetWeeklySummary).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/summaries/monthly", handlers.GetMonthlySummary).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/preferences/theme", handlers.GetThemePreference).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/preferences/theme", handlers.UpdateThemePreference).Methods("PUT", "OPTIONS")

	// Backward-compatible alias for old frontend usage.
	r.HandleFunc("/api/time-entries/summary", handlers.GetSummary).Methods("GET", "OPTIONS")
	r.PathPrefix("/").HandlerFunc(handlers.HandleOptions).Methods("OPTIONS")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server is running on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
