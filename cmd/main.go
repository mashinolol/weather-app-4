package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"weather-app-3/internal/db"
	"weather-app-3/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	BASE_URL := os.Getenv("BASE_URL")
	API_KEY := os.Getenv("API_KEY")

	// Connect to MongoDB
	client, err := db.Connect(MONGO_URI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal("Failed to disconnect MongoDB:", err)
		}
	}()

	// Initialize weather collection
	weatherCollection := client.Database("weatherdb").Collection("weather")
	handlers.WeatherCollection = weatherCollection

	// Setup routes
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetWeatherHandler(w, r)
		case http.MethodPut:
			handlers.PutWeatherHandler(w, r, BASE_URL, API_KEY)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
