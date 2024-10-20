package router

import (
	"database/sql"
	"net/http"
	"y/handler"

	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router(db *sql.DB) *mux.Router {

	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*") // Change "*" to your frontend URL in production
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			next.ServeHTTP(w, r)
		})
	})
	// router.HandleFunc("/api/v1/weather", handler.FetchCurrentWeatherForAllCitiesHandler()).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/weather", handler.GetWeatherDataHandler(db)).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/v1/weather/chart", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeWeatherChart(db, w, r)
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/weather/summary", func(w http.ResponseWriter, r *http.Request) {
		handler.DailyWeatherSummaryHandler(w, r, db)
	}).Methods("GET", "OPTIONS")

	return router
}
