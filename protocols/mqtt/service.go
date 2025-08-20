package mqtt

import "time"

func subscribeToMqtt(url string, handler func(message []byte)) error {
	for {
		// Simulate subscribing to MQTT messages
		// In a real implementation, you would use an MQTT client library to connect and subscribe
		message := []byte("{\"test\":\"test message\"}") // Replace with actual message from MQTT broker
		handler(message)

		// Simulate waiting for the next message
		// In a real implementation, you would wait for the next message from the MQTT broker
		time.Sleep(1 * time.Second)
	}
}
