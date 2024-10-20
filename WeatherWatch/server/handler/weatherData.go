package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"y/models"

	_ "github.com/lib/pq"
)

var cities = []string{"Delhi", "Mumbai", "Chennai", "Bangalore", "Kolkata", "Hyderabad"}

// FetchAndStoreWeatherData fetches weather data for multiple cities and stores it in the database every 5 minutes
func FetchAndStoreWeatherData(db *sql.DB) {
	for {
		for _, city := range cities {
			data, err := fetchWeatherData(city)
			if err != nil {
				log.Printf("Error fetching data for city %s: %v", city, err)
				continue
			}

			err = storeWeatherData(db, data)
			if err != nil {
				log.Printf("Error storing data for city %s: %v", city, err)
				continue
			}
		}

		time.Sleep(5 * time.Minute) // Wait for 5 minutes before fetching data again
	}
}

// fetchWeatherData retrieves weather data from the OpenWeatherMap API for a specific city
func fetchWeatherData(city string) (models.WeatherData, error) {
	var weatherData models.WeatherData

	apiKey := os.Getenv("apiKey")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return weatherData, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return weatherData, err
	}

	main := result["main"].(map[string]interface{})
	weather := result["weather"].([]interface{})[0].(map[string]interface{})

	weatherData.CityName = city
	weatherData.Timestamp = time.Unix(int64(result["dt"].(float64)), 0).Format(time.RFC3339)
	weatherData.Temperature = main["temp"].(float64) - 273.15 // Convert from Kelvin to Celsius
	weatherData.FeelsLike = main["feels_like"].(float64) - 273.15
	weatherData.WeatherMain = weather["main"].(string)

	return weatherData, nil
}

// storeWeatherData inserts the weather data into the database
func storeWeatherData(db *sql.DB, data models.WeatherData) error {
	sqlStatement := `
        INSERT INTO weather_data (city_name, timestamp, temperature, feels_like, weather_main)
        VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(sqlStatement, data.CityName, data.Timestamp, data.Temperature, data.FeelsLike, data.WeatherMain)
	if err != nil {
		return err
	}

	// After storing data, check recent data for alert conditions
	err = CheckRecentWeatherData(db, data.CityName)
	if err != nil {
		log.Printf("Error checking recent weather data: %v", err)
	}

	return nil
}
