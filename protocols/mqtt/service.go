package mqtt

import (
	"context"
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tim0-12432/simple-test-server/config"
)

// startMqttSubscriber starts an MQTT client connected to the given broker URL and
// invokes handler for each received message. It returns a stop function which
// unsubscribes and disconnects the client.
func startMqttSubscriber(ctx context.Context, url string, handler func(message []byte)) (func(), error) {
	opts := mqtt.NewClientOptions()
	if config.EnvConfig.Env == "DEV" {
		log.Printf("Connecting to MQTT broker at %s", url)
	}
	opts.AddBroker("tcp://" + url)
	opts.SetClientID("simple-test-server-mqtt_client")
	opts.SetAutoReconnect(true)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		data := struct {
			Topic   string `json:"topic"`
			Payload string `json:"payload"`
		}{
			Topic:   msg.Topic(),
			Payload: string(msg.Payload()),
		}

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
		return nil, token.Error()
	}

	if token := client.Subscribe("#", 0, nil); token.Wait() && token.Error() != nil {
		client.Disconnect(250)
		return nil, token.Error()
	}

	stop := func() {
		// try to unsubscribe and disconnect gracefully
		if token := client.Unsubscribe("#"); token != nil {
			token.Wait()
		}
		client.Disconnect(250)
	}

	// monitor context cancellation and stop client when cancelled
	go func() {
		<-ctx.Done()
		stop()
	}()

	// keep function non-blocking
	return stop, nil
}
