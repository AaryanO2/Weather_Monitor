package models
// type WeatherData struct {
// 	Weather []struct {
// 		Main string `json:"main"`
// 	} `json:"weather"`
// 	Main struct {
// 		Temp      float64 `json:"temp"`
// 		FeelsLike float64 `json:"feels_like"`
// 	} `json:"main"`
// 	Dt   int64  `json:"dt"`
// 	Name string `json:"name"`
// }

type WeatherData struct {
	ID          int     `json:"id"`
	CityName    string  `json:"city_name"`
	Timestamp   string  `json:"timestamp"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"feels_like"`
	WeatherMain string  `json:"weather_main"`
}

type DailyWeatherSummary struct {
	City               string  `json:"city"`
	Date               string  `json:"date"`
	AverageTemperature float64 `json:"average_temperature"`
	MaxTemperature     float64 `json:"max_temperature"`
	MinTemperature     float64 `json:"min_temperature"`
	DominantCondition  string  `json:"dominant_condition"`
}
