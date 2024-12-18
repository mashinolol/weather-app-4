package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"weather-app-3/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var WeatherCollection *mongo.Collection

func GetWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var weather models.WeatherData
	err := WeatherCollection.FindOne(ctx, bson.M{"city": city}).Decode(&weather)
	if err != nil {
		http.Error(w, "Weather data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func PutWeatherHandler(w http.ResponseWriter, r *http.Request, baseURL, apiKey string) {
	var requestBody struct {
		City string `json:"city"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	city := requestBody.City
	if city == "" {
		http.Error(w, "City is required", http.StatusBadRequest)
		return
	}

	// Fetch weather data from OpenWeather API
	searchURL := fmt.Sprintf("%v?appid=%s&q=%s", baseURL, apiKey, city)
	response, err := http.Get(searchURL)
	if err != nil {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch weather data from API", http.StatusInternalServerError)
		return
	}

	weatherBytes, _ := io.ReadAll(response.Body)
	var weatherAPIResponse models.WeatherJSON
	if err := json.Unmarshal(weatherBytes, &weatherAPIResponse); err != nil {
		http.Error(w, "Failed to parse weather data", http.StatusInternalServerError)
		return
	}

	// Prepare the data for MongoDB
	weatherData := models.WeatherData{
		City:        weatherAPIResponse.Name,
		Description: weatherAPIResponse.Weather[0].Description,
		Temp:        weatherAPIResponse.Main.Temp - 273.15,
		LastUpdated: time.Now(),
	}

	// Upsert data into MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"city": weatherData.City}
	update := bson.M{"$set": weatherData}
	opts := options.Update().SetUpsert(true)

	_, err = WeatherCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		http.Error(w, "Failed to update weather data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weatherData)
}
