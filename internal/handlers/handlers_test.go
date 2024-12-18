package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"weather-app-3/internal/db"
	"weather-app-3/internal/models"
)

func TestGetWeatherHandler(t *testing.T) {
	// Подключаемся к тестовой базе данных MongoDB
	client, err := db.Connect("mongodb+srv://admin:ahu0vpdTyJyXLp6I@cluster0.poo5d.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0") // Убедитесь, что у вас есть работающий локальный MongoDB сервер
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			t.Fatalf("Failed to disconnect MongoDB: %v", err)
		}
	}()

	// Создаем временную коллекцию для тестов
	weatherCollection := client.Database("testdb").Collection("weather")
	WeatherCollection = weatherCollection // Инициализируем глобальную переменную

	// Подготовим тестовые данные
	weatherData := models.WeatherData{
		City:        "TestCity",
		Description: "Clear sky",
		Temp:        20.0,
		LastUpdated: time.Now(),
	}

	// Вставляем тестовые данные в коллекцию
	_, err = weatherCollection.InsertOne(context.Background(), weatherData)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Теперь можно тестировать обработчик
	req, err := http.NewRequest("GET", "/weather?city=TestCity", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Здесь мы будем использовать тестовый HTTP-обработчик
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetWeatherHandler)
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, rr.Code)
	}

	// Проверяем, что тело ответа содержит правильные данные
	var weatherResponse models.WeatherData
	err = json.NewDecoder(rr.Body).Decode(&weatherResponse)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if weatherResponse.City != weatherData.City {
		t.Errorf("Expected city %v, got %v", weatherData.City, weatherResponse.City)
	}
}
