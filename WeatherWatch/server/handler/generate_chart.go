package handler

import (
	"os"
	"time"

	"y/models"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func GenerateLineChart(data []models.WeatherData, city string) (string, error) {
	var xValues []time.Time
	var tempValues, feelsLikeValues []float64

	for _, wd := range data {
		// Ensure that the timestamp is correctly parsed
		timestamp, err := time.Parse(time.RFC3339, wd.Timestamp)
		if err != nil {
			return "", err
		}
		xValues = append(xValues, timestamp)
		tempValues = append(tempValues, wd.Temperature)
		feelsLikeValues = append(feelsLikeValues, wd.FeelsLike)
	}

	// Calculate maximum values for tempValues and feelsLikeValues
	maxTemp := maxFloat64(tempValues)
	maxFeelsLike := maxFloat64(feelsLikeValues)

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:           "Time",
			ValueFormatter: chart.TimeValueFormatterWithFormat("2006-01-02 15:04:05"), // Update format if needed
		},
		YAxis: chart.YAxis{
			Name: "Temperature (°C)",
			Range: &chart.ContinuousRange{
				Min: 0,
				Max: max(maxTemp, maxFeelsLike),
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "Temperature",
				XValues: xValues,
				YValues: tempValues,
				Style: chart.Style{
					StrokeColor: drawing.Color{
						R: 0,
						G: 0,
						B: 255,
						A: 255,
					}, // Blue color for the line
					StrokeWidth: 2.0,
				},
			},
			chart.TimeSeries{
				Name:    "Feels Like",
				XValues: xValues,
				YValues: feelsLikeValues,
				Style: chart.Style{
					StrokeColor: drawing.Color{
						R: 0,
						G: 255,
						B: 0,
						A: 255,
					}, // Green color for the line
					StrokeWidth: 2.0,
				},
			},
		},
	}

	filename := city + "_weather_chart.png"
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = graph.Render(chart.PNG, file)
	if err != nil {
		return "", err
	}

	return filename, nil
}

// Helper function to find the maximum value in a slice of float64
func maxFloat64(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

// Helper function to find the maximum value between two float64
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
