package mqttpackage

import (
	"fmt"
	"log"

	"hiliriset_ecoprint_golang/config"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
    client paho.Client
}

func NewMQTTClient(handler *MQTTHandler) *MQTTClient {
    opts := paho.NewClientOptions().
        AddBroker(fmt.Sprintf("tcp://%s:%s", config.APPConfig.MQTTHost, config.APPConfig.MQTTPort)).
        SetClientID("hiliriset_ecoprint_backend").
        SetCleanSession(true).
        SetOnConnectHandler(func(c paho.Client) {
            log.Println("MQTT connected")
            subscribeTopics(c, handler)
        }).
        SetConnectionLostHandler(func(c paho.Client, err error) {
            log.Printf("MQTT connection lost: %v", err)
        })


    client := paho.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatalf("MQTT failed to connect: %v", token.Error())
    }

    return &MQTTClient{client: client}
}

func subscribeTopics(client paho.Client, handler *MQTTHandler) {
    // Wildcard — catches all ESPs: esp/any-public-id/telemetry
    token := client.Subscribe("esp/+/telemetry", 1, handler.handleTelemetry)
    token.Wait()
    if token.Error() != nil {
        log.Printf("MQTT failed to subscribe: %v", token.Error())
    }
    log.Println("MQTT subscribed to esp/+/telemetry")
}

func (m *MQTTClient) Publish(topic string, payload string) error {
    token := m.client.Publish(topic, 1, false, payload)
    token.Wait()
    return token.Error()
}

func (m *MQTTClient) Disconnect() {
    m.client.Disconnect(250)
}
