package mqttclient

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client mqtt.Client
}

func NewClient(broker string, port int) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Ошибка подключения к MQTT-брокеру: %v", token.Error())
	}

	return &MQTTClient{client: client}
}

func (m *MQTTClient) PublishTemperature(temperature float64) {
    topic := "temperature/CurrentOutdoor"
    m.client.Publish(topic, 0, false, fmt.Sprintf("%.2f", temperature))
    log.Printf("Опубликовано: температура за бортом = %.2f", temperature)
}