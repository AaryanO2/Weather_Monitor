package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"y/models"
)

// GetWeatherDataHandler retrieves weather data for a specific city and date from the database
func GetCityWeatherDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		city := queryParams.Get("city")
		date := queryParams.Get("date") // Expected format: "YYYY-MM-DD"

		if city == "" || date == "" {
			http.Error(w, "City and date parameters are required", http.StatusBadRequest)
			return
		}

		sqlStatement := `
			SELECT id, city_name, timestamp, temperature, feels_like, weather_main
			FROM weather_data
			WHERE city_name = $1 AND timestamp = (
				SELECT MAX(timestamp) 
				FROM weather_data 
				WHERE city_name = $1)`

		rows, err := db.Query(sqlStatement, city, date)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to execute the query: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var weatherData []models.WeatherData

		for rows.Next() {
			var data models.WeatherData
			err := rows.Scan(&data.ID, &data.CityName, &data.Timestamp, &data.Temperature, &data.FeelsLike, &data.WeatherMain)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to scan the row: %v", err), http.StatusInternalServerError)
				return
			}
			weatherData = append(weatherData, data)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, fmt.Sprintf("Row iteration error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(weatherData)
	}
}

// Returns Weather for all cities from database
func GetWeatherDataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// List of cities for which you want to fetch the weather data
		cities := []string{"Delhi", "Mumbai", "Chennai", "Bangalore", "Kolkata"}

		// Prepare a slice to hold the weather data for all cities
		var weatherData []models.WeatherData

		for _, city := range cities {
			// Adjust the SQL statement to filter by current city and get the most recent entry
			sqlStatement := `
                SELECT id, city_name, timestamp, temperature, feels_like, weather_main
                FROM weather_data
                WHERE city_name = $1 
                AND timestamp = (
                    SELECT MAX(timestamp) 
                    FROM weather_data 
                    WHERE city_name = $1
                )`

			row := db.QueryRow(sqlStatement, city)

			var data models.WeatherData
			err := row.Scan(&data.ID, &data.CityName, &data.Timestamp, &data.Temperature, &data.FeelsLike, &data.WeatherMain)
			if err != nil {
				if err == sql.ErrNoRows {
					// If no data found for a specific city, you can log it or skip it
					continue // Skip to the next city
				} else {
					http.Error(w, fmt.Sprintf("Unable to scan the row for city %s: %v", city, err), http.StatusInternalServerError)
					return
				}
			}

			weatherData = append(weatherData, data)
		}

		// Check if any weather data was retrieved
		if len(weatherData) == 0 {
			http.Error(w, "No weather data found for the current time", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(weatherData) // Encode the slice of weather data
	}
}

// Returns REALTIME Weather for all cities from API
func FetchCurrentWeatherForAllCitiesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var allWeatherData []models.WeatherData

		// Loop through each city and fetch its current weather data
		for _, city := range cities {
			weatherData, err := fetchWeatherData(city)
			if err != nil {
				log.Printf("Error fetching weather data for city %s: %v", city, err)
				http.Error(w, fmt.Sprintf("Unable to fetch weather data for %s", city), http.StatusInternalServerError)
				return
			}
			allWeatherData = append(allWeatherData, weatherData)
		}

		// Set content type as JSON and encode the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allWeatherData)
	}
}
