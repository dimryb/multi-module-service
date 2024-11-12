# Модуль MQTTClient

## Обзор

Модуль `mqttclient` предоставляет упрощенную обертку над [библиотекой Paho MQTT для Go](https://github.com/eclipse/paho.mqtt.golang), позволяя легко публиковать сообщения и подписываться на топики MQTT. Модуль предоставляет интерфейс для взаимодействия с MQTT-брокерами с помощью простых методов.

## Функциональность

- **Подключение** к MQTT-брокеру.
- **Публикация** сообщений в указанные топики MQTT.
- **Подписка** на топики с обработкой входящих сообщений.
- **Отключение** от брокера с завершением всех операций.

## Установка

Для использования модуля добавьте его в свой проект. Убедитесь, что библиотека `paho.mqtt.golang` установлена. Если она не установлена, выполните команду:

```bash
go get github.com/eclipse/paho.mqtt.golang
```

После этого подключите модуль в файлах вашего проекта.

## Использование

### Подключение модуля

```go
import "path/to/your/module/mqttclient"
```

### Создание нового MQTT клиента

Для создания нового клиента используйте функцию `NewClient`, указав адрес и порт брокера:

```go
client, err := mqttclient.NewClient("broker.example.com", 1883)
if err != nil {
    log.Fatalf("Не удалось подключиться к MQTT брокеру: %v", err)
}
defer client.Disconnect(250) // Отключение с ожиданием 250 мс для завершения операций
```

Функция возвращает реализацию интерфейса `MQTTClientInterface`, включающего методы `Publish`, `Subscribe` и `Disconnect`.

### Публикация сообщения

Для публикации сообщения в определенную топик используйте метод `Publish`, указав топик, уровень QoS, флаг сохранения и данные сообщения:

```go
err := client.Publish("test/topic", 0, false, "Привет, MQTT!")
if err != nil {
    log.Printf("Не удалось отправить сообщение: %v", err)
}
```

- **topic**: Название топики для публикации сообщения.
- **qos**: Уровень качества обслуживания (QoS) (0, 1 или 2).
- **retained**: Если true, сообщение сохраняется брокером.
- **payload**: Данные сообщения.

### Подписка на топик

Для подписки на топик и обработки входящих сообщений используйте метод `Subscribe`:

```go
callback := func(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Получено сообщение в теме %s: %s\n", msg.Topic(), string(msg.Payload()))
}

err := client.Subscribe("test/topic", 0, callback)
if err != nil {
    log.Printf("Не удалось подписаться на топик: %v", err)
}
```

- **topic**: Название топика для подписки.
- **qos**: Уровень QoS (0, 1 или 2).
- **callback**: Функция-обработчик сообщений, должна иметь сигнатуру `func(mqtt.Client, mqtt.Message)`.

### Отключение от брокера

Для корректного завершения работы клиента используйте метод `Disconnect`:

```go
client.Disconnect(250)
```

Параметр `quiesce` указывает время ожидания (в мс) перед отключением, позволяя клиенту завершить операции.

## Пример

Полный пример использования модуля `mqttclient`:

```go
package main

import (
    "fmt"
    "log"
    "time"
    "path/to/your/module/mqttclient"
)

func main() {
    // Инициализация MQTT клиента
    client, err := mqttclient.NewClient("localhost", 1883)
    if err != nil {
        log.Fatalf("Не удалось подключиться к MQTT брокеру: %v", err)
    }
    defer client.Disconnect(250)

    // Подписка на топик
    callback := func(client mqtt.Client, msg mqtt.Message) {
        fmt.Printf("Получено сообщение в теме %s: %s\n", msg.Topic(), string(msg.Payload()))
    }
    err = client.Subscribe("test/topic", 0, callback)
    if err != nil {
        log.Printf("Не удалось подписаться на топик: %v", err)
    }

    // Публикация сообщения
    err = client.Publish("test/topic", 0, false, "Привет, MQTT!")
    if err != nil {
        log.Printf("Не удалось отправить сообщение: %v", err)
    }

    // Время для обработки сообщений
    time.Sleep(2 * time.Second)
}
```

## Справка по API

### `NewClient`

```go
func NewClient(broker string, port int) (MQTTClientInterface, error)
```

Создает и подключает MQTT клиента к брокеру.

- **broker**: IP или hostname брокера.
- **port**: Порт брокера.
- **Возвращает**: Экземпляр `MQTTClientInterface` и ошибку при неудачном подключении.

### `Publish`

```go
func (m *mqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error
```

Публикует сообщение в указанную топик.

- **topic**: Название топики.
- **qos**: Уровень QoS.
- **retained**: Флаг сохранения сообщения.
- **payload**: Данные сообщения.

### `Subscribe`

```go
func (m *mqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error
```

Подписывается на топик с обработкой входящих сообщений.

- **topic**: Название топика.
- **qos**: Уровень QoS.
- **callback**: Функция-обработчик сообщений.

### `Disconnect`

```go
func (m *mqttClient) Disconnect(quiesce uint)
```

Отключает клиента с ожиданием завершения операций.

- **quiesce**: Время ожидания перед отключением в миллисекундах.

## Зависимости

Этот модуль зависит от следующей библиотеки:
- **Paho MQTT client** (`github.com/eclipse/paho.mqtt.golang`)