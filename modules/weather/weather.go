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
	latitude      string
	longitude     string
	mqttClient    mqttclient.MQTTClientInterface
	publishPeriod time.Duration
}

func NewWeatherService(client mqttclient.MQTTClientInterface, latitude, longitude string, period time.Duration) *WeatherService {
	return &WeatherService{
		latitude:      latitude,
		longitude:     longitude,
		mqttClient:    client,
		publishPeriod: period,
	}
}

func (w *WeatherService) Run() {
	log.Printf("Модуль WeatherService запущен. Период публикации: %v", w.publishPeriod)

	w.fetchAndPublishWeatherData()

	ticker := time.NewTicker(w.publishPeriod)	
	defer ticker.Stop()	

	for range ticker.C {
		w.fetchAndPublishWeatherData()
	}
}

func (w *WeatherService) fetchAndPublishWeatherData() {
	temperature, err := w.getWeatherData(w.latitude, w.longitude)
	if err != nil {
		log.Printf("Ошибка получения данных о погоде: %v", err)
		return
	}
	log.Printf("Отправка в MQTT температуры: %.2f", temperature)
	w.mqttClient.Publish("temperature/CurrentOutdoor", 0, false, fmt.Sprintf("%.2f", temperature))
}

func (w *WeatherService) getWeatherData(latitude, longitude string) (float64, error) {
	apiURL := fmt.Sprintf("%s?latitude=%s&longitude=%s&current_weather=true", WeatherAPIBaseURL, latitude, longitude)
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
