package mqttclient

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClientInterface interface {
    Publish(topic string, qos byte, retained bool, payload interface{}) error
    Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error
    Disconnect(quiesce uint)
}

type mqttClient  struct {
	client mqtt.Client
}

func NewClient(broker string, port int) (MQTTClientInterface, error) {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &mqttClient {client: client}, nil
}

func (m *mqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
    token := m.client.Publish(topic, qos, retained, payload)
    token.Wait()
    return token.Error()
}

func (m *mqttClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
    token := m.client.Subscribe(topic, qos, callback)
    token.Wait()
    return token.Error()
}

func (m *mqttClient) Disconnect(quiesce uint) {
    m.client.Disconnect(quiesce)
}