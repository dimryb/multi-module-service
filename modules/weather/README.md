# Модуль WeatherService

## Обзор

Модуль `WeatherService` предоставляет сервис мониторинга температуры на основе данных из [API Open-Meteo](https://open-meteo.com/). Он периодически запрашивает текущую температуру по заданным координатам и публикует данные в MQTT-топик.

## Функциональность

- **Получение данных о погоде** через HTTP-запрос к API.
- **Публикация температуры** в MQTT-топик с определенным интервалом.

## Установка и использование

Для использования модуля `WeatherService` необходимо:
1. Установить библиотеку `paho.mqtt.golang`, если она не была установлена ранее.
2. Подключить модуль `weather` и модуль `mqttclient` в ваш проект.

```bash
go get github.com/eclipse/paho.mqtt.golang
```

### Подключение модуля в проекте

```go
import "path/to/your/module/weather"
```

## Использование

### Создание сервиса WeatherService

Для создания нового экземпляра `WeatherService`, выполните инициализацию с указанием MQTT-клиента, координат и интервала публикации:

```go
mqttClient, _ := mqttclient.NewClient("broker.example.com", 1883)
weatherService := weather.NewWeatherService(mqttClient, "55.7558", "37.6173", 10*time.Minute)
```

- **mqttClient**: MQTT-клиент, реализующий интерфейс `MQTTClientInterface`.
- **latitude** и **longitude**: Широта и долгота для получения данных о погоде.
- **period**: Период публикации температуры в MQTT.

### Запуск модуля WeatherService

Для запуска службы вызовите метод `Run()`, который будет циклически получать данные о погоде и публиковать их в MQTT:

```go
weatherService.Run()
```

Сервис будет автоматически публиковать температуру в указанном топике с заданным интервалом.

## Пример кода

```go
package main

import (
    "log"
    "time"
    "path/to/your/module/weather"
    "path/to/your/module/mqttclient"
)

func main() {
    // Инициализация MQTT-клиента
    mqttClient, err := mqttclient.NewClient("broker.example.com", 1883)
    if err != nil {
        log.Fatalf("Не удалось подключиться к MQTT брокеру: %v", err)
    }
    defer mqttClient.Disconnect(250)

    // Инициализация WeatherService
    weatherService := weather.NewWeatherService(mqttClient, "55.7558", "37.6173", 10*time.Minute)

    // Запуск WeatherService
    weatherService.Run()
}
```

## Справка по API

### `NewWeatherService`

```go
func NewWeatherService(client mqttclient.MQTTClientInterface, latitude, longitude string, period time.Duration) *WeatherService
```

Создает новый экземпляр `WeatherService`.

- **client**: MQTT-клиент.
- **latitude**: Широта местоположения.
- **longitude**: Долгота местоположения.
- **period**: Интервал публикации данных.

### `Run`

```go
func (w *WeatherService) Run()
```

Запускает цикл публикации данных о температуре в MQTT. Период публикации указан при инициализации.

### `fetchAndPublishWeatherData`

```go
func (w *WeatherService) fetchAndPublishWeatherData()
```

Запрашивает данные о температуре через API и публикует их в MQTT.

### `getWeatherData`

```go
func (w *WeatherService) getWeatherData(latitude, longitude string) (float64, error)
```

Выполняет HTTP-запрос к API Open-Meteo для получения текущей температуры.

- **latitude** и **longitude**: Координаты местоположения.
- **Возвращает**: Текущую температуру и ошибку, если запрос не удался.

## Зависимости

Модуль зависит от:
- **paho.mqtt.golang** для работы с MQTT (`github.com/eclipse/paho.mqtt.golang`)
- **encoding/json** для обработки JSON-ответов