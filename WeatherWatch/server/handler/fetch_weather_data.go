package handler

import (
    "database/sql"
    "time"
    "y/models"
)

// Function to fetch weather data at 5-minute intervals
func FetchWeatherDataAtInterval(db *sql.DB, city string, start, end time.Time) ([]models.WeatherData, error) {
    rows, err := db.Query(`
        SELECT id, city_name, timestamp, temperature, feels_like, weather_main
        FROM weather_data
        WHERE city_name = $1 AND timestamp BETWEEN $2 AND $3
        ORDER BY timestamp ASC`,
        city, start, end)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var data []models.WeatherData
    for rows.Next() {
        var wd models.WeatherData
        if err := rows.Scan(&wd.ID, &wd.CityName, &wd.Timestamp, &wd.Temperature, &wd.FeelsLike, &wd.WeatherMain); err != nil {
            return nil, err
        }
        data = append(data, wd)
    }
    return data, nil
}
