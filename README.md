# WeatherWatch
This project is a Go application for fetching, storing, and querying weather data using the OpenWeatherMap API and a PostgreSQL database. The application includes functionality for:

## Features
- Weather Data Fetching: Fetches current weather data from the OpenWeatherMap API.
- Data Storage: Stores weather data in a PostgreSQL database.
- Daily Weather Summaries: Computes and stores daily summaries of weather data.
- REST API: Provides endpoints to access weather data and summaries.
- Retry Logic: Handles database connection retries to ensure reliability.
  
## Prerequisites
1. Go 1.18 or higher
2. Docker and Docker Compose
3. PostgreSQL
4. OpenWeatherMap API key
## Setup
1. Clone the Repository
    Copy code
  ```bash
  - git clone https://github.com/AaryanO2/Weather_Monitor.git
  - cd WeatherWatch
```
2. Configure Environment Variables
  - Create a .env file in the root directory and add your environment variables:
  - eg: 
      env
      Copy code
      POSTGRES_USER=postgres
      POSTGRES_PASSWORD=yourPassword
      POSTGRES_DB=codedb
      DATABASE_URL=postgres://postgres:yourPassword@db:5432/codedb?sslmode=disable
  -Update the Docker-compose file replace YOUR_API_KEY with OpenWeatherMap API key
3. Build and Run the Application
4. Using Docker Compose, build and start the services:
  ```bash
  - cd server
  - docker-compose up --build
```
This command will set up the PostgreSQL database and the Go server.

## Accessing the Client
1. Open the client folder.
2. Click on the index.html file to access the website.
