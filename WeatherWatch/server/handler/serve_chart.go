package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ServeWeatherChart handles requests to generate and serve a weather chart.
func ServeWeatherChart(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// Parse query parameters
	city := r.URL.Query().Get("city")
	start, err := time.Parse(time.RFC3339, r.URL.Query().Get("start"))
	if err != nil {
		http.Error(w, "Invalid start time format. Expected RFC3339 format.", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(time.RFC3339, r.URL.Query().Get("end"))
	if err != nil {
		http.Error(w, "Invalid end time format. Expected RFC3339 format.", http.StatusBadRequest)
		return
	}

	// Fetch weather data at the specified intervals
	data, err := FetchWeatherDataAtInterval(db, city, start, end)
	if err != nil {
		log.Printf("Error fetching weather data: %v", err)
		http.Error(w, "Error fetching weather data", http.StatusInternalServerError)
		return
	}

	// Generate line chart for the fetched data
	fmt.Println(data)
	fmt.Println("he")
	chartPath, err := GenerateLineChart(data, city)
	if err != nil {
		log.Printf("Error generating chart: %v", err)
		http.Error(w, "Error generating chart", http.StatusInternalServerError)
		return
	}

	// Serve the generated chart file
	http.ServeFile(w, r, chartPath)
}
