package mqtt

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tim0-12432/simple-test-server/config"
)

func subscribeToMqtt(url string, handler func(message []byte)) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + url)
	opts.SetClientID("simple-test-server-mqtt_client")
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		data := struct {
			Topic   string `json:"topic"`
			Payload string `json:"payload"`
		}{
			Topic:   msg.Topic(),
			Payload: string(msg.Payload()),
		}

		// Marshal to JSON
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshaling MQTT message: %v", err)
			return
		}
		if config.EnvConfig.Env == "DEV" {
			log.Printf("Received MQTT message: %s", jsonBytes)
		}

		handler(jsonBytes)
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if token := client.Subscribe("#", 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
