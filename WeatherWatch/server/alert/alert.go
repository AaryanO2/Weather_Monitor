package alert

import (
	"fmt"
	"y/models"
)

func TriggerAlert(data models.WeatherData) {
	fmt.Printf("ALERT! City: %s, Temperature: %.2f°C, Condition: %s at %s\n",
		data.CityName, data.Temperature, data.WeatherMain, data.Timestamp)

	// Optionally send an email alert
	// sendEmailAlert(data)
}
