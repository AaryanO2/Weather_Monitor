CREATE TABLE weather_data (
    id SERIAL PRIMARY KEY,
    city_name TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    temperature DECIMAL NOT NULL,
    feels_like DECIMAL NOT NULL,
    weather_main TEXT NOT NULL
);


CREATE TABLE IF NOT EXISTS daily_weather_summaries (
    id SERIAL PRIMARY KEY,
    city VARCHAR(50) NOT NULL,
    date DATE NOT NULL,
    average_temperature DECIMAL NOT NULL,
    max_temperature DECIMAL NOT NULL,
    min_temperature DECIMAL NOT NULL,
    dominant_condition VARCHAR(50) NOT NULL
);
