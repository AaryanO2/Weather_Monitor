package handler

import (
	"database/sql"
	"fmt"
	"y/alert"
	"y/config"
	"y/models"
)

// checkRecentWeatherData retrieves the most recent two weather data entries for a city
// and checks if they meet the configured conditions. If they do, it triggers an alert.
func CheckRecentWeatherData(db *sql.DB, city string) error {
	// Query to get the most recent two weather data entries for the specified city
	sqlStatement := `
		SELECT city_name, timestamp, temperature, feels_like, weather_main
		FROM weather_data
		WHERE city_name = $1
		ORDER BY timestamp DESC
		LIMIT 2`

	rows, err := db.Query(sqlStatement, city)
	if err != nil {
		return fmt.Errorf("unable to execute the query: %v", err)
	}
	defer rows.Close()

	var recentData []models.WeatherData

	// Scan the results into the recentData slice
	for rows.Next() {
		var data models.WeatherData
		if err := rows.Scan(&data.CityName, &data.Timestamp, &data.Temperature, &data.FeelsLike, &data.WeatherMain); err != nil {
			return fmt.Errorf("unable to scan the row: %v", err)
		}
		recentData = append(recentData, data)
	}

	// Ensure we have exactly two entries to compare
	if len(recentData) < 2 {
		return fmt.Errorf("not enough data to compare for city: %s", city)
	}

	// Compare the recent two entries against the configured thresholds
	exceeds := false
	for _, data := range recentData {
		if data.Temperature > config.Thresholds.Temperature {
			exceeds = true
		}

		if config.Thresholds.Condition != "" && data.WeatherMain != config.Thresholds.Condition {
			exceeds = false
			break
		}
	}

	// Trigger an alert if the conditions are met
	if exceeds {
		alert.TriggerAlert(recentData[0])
	}

	return nil
}
