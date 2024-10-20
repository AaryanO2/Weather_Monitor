package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"y/models"
)

func DailyWeatherSummaryHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Set the header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Get the city name from the query parameters
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City name is required", http.StatusBadRequest)
		return
	}

	// Prepare to hold the summaries
	var summaries []models.DailyWeatherSummary

	// Iterate over the past 5 days
	for i := 0; i < 5; i++ {
		date := time.Now().AddDate(0, 0, -i)                    // Get the date for i days ago
		summary := CalculateDailyWeatherSummary(db, city, date) // Call the function to get the summary
		summaries = append(summaries, summary)                  // Append the summary to the slice
	}

	// Encode the summaries to JSON and write to the response
	if err := json.NewEncoder(w).Encode(summaries); err != nil {
		http.Error(w, "Failed to encode summary", http.StatusInternalServerError)
		log.Println("Error encoding summary:", err)
		return
	}
}
