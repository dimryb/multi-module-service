package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const ( // TODO: конфигурацию брать из конфиг файла и получать по mqtt
	Version           = "2.0.0"
	DefaultLatitude   = "56.4977"
	DefaultLongitude  = "84.9744"
	WeatherAPIBaseURL = "https://api.open-meteo.com/v1/forecast"
	WeatherAPITimeout = 5 * time.Second
	MQTTTopicOut      = "temperature/CurrentOutdoor"
	MQTTTopicIn       = "temperature/CurrentIndoor"
	MQTTBrokerHost    = "localhost" // Укажите свой адрес и порт брокера
	MQTTBrokerPort    = 1883
	PublishPeriod     = 5 * time.Second
)

type WeatherMqtt struct {
	latitude         string
	longitude        string
	defaultLatitude  string
	defaultLongitude string
	mqttClient       mqtt.Client
}

type WeatherResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
	} `json:"current_weather"`
}

func NewWeatherMqtt(broker string, port int, defaultLat string, defaultLon string) *WeatherMqtt {
	mqttOptions := mqtt.NewClientOptions()
	mqttOptions.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))

	client := mqtt.NewClient(mqttOptions)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Ошибка подключения к MQTT-брокеру: %v", token.Error())
	}

	return &WeatherMqtt{
		latitude:         defaultLat,
		longitude:        defaultLon,
		defaultLatitude:  defaultLat,
		defaultLongitude: defaultLon,
		mqttClient:       client,
	}
}

func (w *WeatherMqtt) onConnect(client mqtt.Client) {
	log.Println("Подключено к MQTT-брокеру")
	client.Subscribe("Gnss/Latitude", 0, w.onMessage) // TODO: топики в конфиг и не хардкодить
	client.Subscribe("Gnss/Longitude", 0, w.onMessage)
}

func (w *WeatherMqtt) onMessage(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())
	if topic == "Gnss/Latitude" && isValidCoord(payload) { // TODO: топики в конфиг и не хардкодить
		w.latitude = payload
	} else if topic == "Gnss/Longitude" && isValidCoord(payload) {
		w.longitude = payload
	}
}

func (w *WeatherMqtt) getWeatherData() (float64, error) {
	lat := w.latitude
	lon := w.longitude
	apiURL := fmt.Sprintf("%s?latitude=%s&longitude=%s&current_weather=true", WeatherAPIBaseURL, lat, lon)

	resp, err := http.Get(apiURL)

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

func isValidCoord(coord string) bool {
	_, err := strconv.ParseFloat(coord, 64)
	return err == nil
}

func (w *WeatherMqtt) sendToMqtt(temperature float64) {
	w.mqttClient.Publish("temperature/CurrentOutdoor", 0, true, fmt.Sprintf("%.2f", temperature)) // TODO: топики в конфиг и не хардкодить
	log.Printf("Опубликовано: температура за бортом = %.2f", temperature)
}

func (w *WeatherMqtt) Run() {
	w.onConnect(w.mqttClient)
	ticker := time.NewTicker(PublishPeriod)
	defer ticker.Stop()

	w.requestAndPublishWeatherData()

	for range ticker.C {
		w.requestAndPublishWeatherData()
	}
}

func (w *WeatherMqtt) requestAndPublishWeatherData() {
	temp, err := w.getWeatherData()
	if err != nil {
		log.Printf("ошибка получения данных о погоде: %v", err)
		return
	}
	w.sendToMqtt(temp)
}

func checkFlags() {
	versionFlag := flag.Bool("version", false, "Показать версию программы")
	debugFlag := flag.Bool("debug", false, "Включить дебаг-режим")
	configFileFlag := flag.String("config", "", "Путь к файлу конфигурации")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("go-weather-mqtt версия %s\n", Version)
		os.Exit(0)
	}

	if *debugFlag {
		fmt.Println("Дебаг-режим включен")
		// TODO: включить дополнительные логи или отладочную информацию
	}

	if *configFileFlag != "" {
		fmt.Printf("Используется файл конфигурации: %s\n", *configFileFlag)
		// TODO: добавить логику для загрузки иной конфигурации
	}
}

func main() {
	checkFlags()

	weatherMqtt := NewWeatherMqtt(MQTTBrokerHost, MQTTBrokerPort, DefaultLatitude, DefaultLongitude)
	weatherMqtt.Run()
}
