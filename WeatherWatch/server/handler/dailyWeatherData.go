package handler

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"y/models"
)

// Function to fetch weather data for a specific city and date from the database
func FetchWeatherDataForCity(db *sql.DB, city string, date time.Time) []models.WeatherData {
	query := `
    SELECT temperature, weather_main, timestamp
    FROM weather_data
    WHERE city_name = $1 AND DATE(timestamp) = $2`

	today := time.Now().Format("2006-01-02")

	rows, err := db.Query(query, city, today)
	if err != nil {
		log.Fatalf("Failed to fetch weather data for %s: %v", city, err)
	}
	defer rows.Close()

	var weatherData []models.WeatherData
	for rows.Next() {
		var data models.WeatherData
		if err := rows.Scan(&data.Temperature, &data.WeatherMain, &data.Timestamp); err != nil {
			log.Fatalf("Failed to scan weather data: %v", err)
		}
		weatherData = append(weatherData, data)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error in weather data rows: %v", err)
	}
	return weatherData
}

func CalculateDailyWeatherSummary(db *sql.DB, city string, date time.Time) models.DailyWeatherSummary {
	weatherData := FetchWeatherDataForCity(db, city, date)
	fmt.Println(weatherData)
	var totalTemp, maxTemp, minTemp float64
	weatherConditionCount := make(map[string]int)
	count := 0

	for _, data := range weatherData {
		totalTemp += data.Temperature
		if data.Temperature > maxTemp || count == 0 {
			maxTemp = data.Temperature
		}
		if minTemp == 0 || data.Temperature < minTemp {
			minTemp = data.Temperature
		}
		weatherConditionCount[data.WeatherMain]++
		count++
	}

	// Determine dominant weather condition
	dominantCondition := ""
	maxCount := 0
	for condition, cnt := range weatherConditionCount {
		if cnt > maxCount {
			maxCount = cnt
			dominantCondition = condition
		}
	}

	// Calculate the average temperature
	averageTemp := totalTemp / float64(count)

	// Create the summary
	summary := models.DailyWeatherSummary{
		City:               city,
		Date:               date.Format("2006-01-02"),
		AverageTemperature: averageTemp,
		MaxTemperature:     maxTemp,
		MinTemperature:     minTemp,
		DominantCondition:  dominantCondition,
	}

	// Store the daily weather summary in the database
	fmt.Println(summary)
	storeDailyWeatherSummary(db, summary)
	return summary
}

// Function to store the daily weather summary in the database
func storeDailyWeatherSummary(db *sql.DB, summary models.DailyWeatherSummary) {
	sqlStatement := `
    INSERT INTO daily_weather_summaries (city, date, average_temperature, max_temperature, min_temperature, dominant_condition)
    VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(sqlStatement, summary.City, summary.Date, summary.AverageTemperature, summary.MaxTemperature, summary.MinTemperature, summary.DominantCondition)

	if err != nil {
		log.Fatalf("Unable to store daily weather summary: %v", err)
	}
}
