package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"y/handler"
	"y/router"

	_ "github.com/lib/pq"
)

func main() {
	var db *sql.DB
	var err error
	// var databaseURL string

	// Get the DATABASE_URL environment variable
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Retry connecting to the database
	for i := 0; i < 10; i++ {
		// Open the database connection
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			log.Printf("Error opening database connection: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Check if the database is reachable
		if err := db.Ping(); err != nil {
			log.Printf("Error pinging database: %v", err)
			db.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// Successfully connected
		log.Println("Successfully connected to the database!")
		break
	}

	if err != nil {
		log.Fatalf("Could not connect to the database after retries: %v", err)
	}
	defer db.Close()

	go handler.FetchAndStoreWeatherData(db)
	go handler.ScheduleDailyWeatherSummary(db) // Schedule daily summary

	r := router.Router(db)
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
