package models

import "time"

// WeatherData represents the structure for weather data in MongoDB
type WeatherData struct {
	City        string    `bson:"city" json:"city"`
	Description string    `bson:"description" json:"description"`
	Temp        float64   `bson:"temp" json:"temp"`
	LastUpdated time.Time `bson:"last_updated" json:"last_updated"`
}

// WeatherJSON represents the response structure from the OpenWeather API
type WeatherJSON struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`

	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`

	Name string `json:"name"`
}
