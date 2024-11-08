package weather

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"multi-module-service/modules/mqttclient"
)

const (
	WeatherAPIBaseURL = "https://api.open-meteo.com/v1/forecast"
	WeatherAPITimeout = 5 * time.Second
)

type WeatherService struct {
	latitude  string
	longitude string
	client    *mqttclient.MQTTClient
	period int
}

func NewWeatherService(latitude, longitude string, client *mqttclient.MQTTClient, period int) *WeatherService {
	return &WeatherService{
		latitude:  latitude,
		longitude: longitude,
		client:    client,
		period: period,
	}
}

func (w *WeatherService) StartWeatherCycle() {
	go func() {
		period := 5 * time.Second
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		for range ticker.C {
			temp, err := w.getWeatherData()
			if err != nil {
				log.Printf("ошибка получения данных о погоде: %v", err)
				continue
			}
			w.client.PublishTemperature(temp)
		}
	}()
}

func (w *WeatherService) getWeatherData() (float64, error) {
	apiURL := fmt.Sprintf("%s?latitude=%s&longitude=%s&current_weather=true", WeatherAPIBaseURL, w.latitude, w.longitude)
	client := &http.Client{Timeout: WeatherAPITimeout}

	resp, err := client.Get(apiURL)
	if err != nil {
		return 0, fmt.Errorf("ошибка запроса к API погоды: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		CurrentWeather struct {
			Temperature float64 `json:"temperature"`
		} `json:"current_weather"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return result.CurrentWeather.Temperature, nil
}
